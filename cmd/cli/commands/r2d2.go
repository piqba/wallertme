package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	domain "github.com/piqba/wallertme/internal/r2d2/domain"

	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/exporters"
	"github.com/piqba/wallertme/pkg/logger"
	"github.com/piqba/wallertme/pkg/notify"
	"github.com/spf13/cobra"
)

var (
	ADA = exporters.TXS_STREAM_KEY + ":" + "ADA"
	SOL = exporters.TXS_STREAM_KEY + ":" + "SOL"
)
var consumerCmd = &cobra.Command{
	Use:   "r2d2",
	Short: "Subscribe to Txs data topic from (REDIS) and send notifications to this services (telegram|discord|smtp)",
	Run: func(cmd *cobra.Command, args []string) {
		notifier, err := cmd.Flags().GetString(flagNotifier)
		if err != nil {
			logger.LogError(errors.Errorf("r2d2ctl: %v", err).Error())
		}

		redisDbClient := exporters.GetRedisDbClient()
		// telegram client
		tgClientBot := notify.GetTgBot(notify.TgBotOption{
			Debug: false,
			Token: "",
		})

		// discord client
		discordClient, err := notify.NewDiscordClient(notify.DiscordClientOptions{})
		if err != nil {
			logger.LogError(err.Error())
		}

		// smtp client
		smtpClient := notify.NewSender(
			os.Getenv("SMTP_EMAIL_USER"),
			os.Getenv("SMTP_EMAIL_PASSWORD"),
		)

		// notification type
		notificationType := notifier

		var repo domain.TxRepository
		switch notificationType {
		case notify.TELEGRAM:
			repo = domain.NewTxRepository(
				domain.ExternalOptions{
					Type:              notificationType,
					DstNotificationID: 927486129,
				},
				tgClientBot,
			)
		case notify.DISCORD:
			repo = domain.NewTxRepository(
				domain.ExternalOptions{
					Type: notificationType,
				},
				discordClient,
			)
		case notify.SMTP:
			repo = domain.NewTxRepository(domain.ExternalOptions{
				Type:              notificationType,
				DstNotificationID: 0,
			},
				smtpClient,
			)
		}

		streams := []string{ADA, SOL}
		var ids []string
		consumersGroup := "r2d2-consumer"
		for _, v := range streams {
			ids = append(ids, ">")
			err := redisDbClient.XGroupCreate(context.TODO(), v, consumersGroup, "0").Err()
			if err != nil {
				log.Println(err)
			}

		}

		streams = append(streams, ids...) // for each stream it requires an '>' :{"txs", ">"}

		for {
			entries, err := redisDbClient.XReadGroup(context.Background(), &redis.XReadGroupArgs{
				Group:    consumersGroup,
				Consumer: fmt.Sprintf("%d", time.Now().UnixNano()),
				Streams:  streams,
				Count:    2,
				Block:    0,
				NoAck:    false,
			}).Result()
			if err != nil {
				log.Fatal(err)
			}

			for _, it := range entries {
				Exec(redisDbClient, consumersGroup, it, repo)
			}

		}

	},
}

func init() {
	consumerCmd.Flags().String(flagNotifier, notify.TELEGRAM, "select a provider to send notifications")
	rootCmd.AddCommand(consumerCmd)

}

func Exec(
	rdb *redis.Client,
	consumersGroup string,
	stream redis.XStream,
	repo domain.TxRepository,
) {
	for i := 0; i < len(stream.Messages); i++ {
		messageID := stream.Messages[i].ID
		values := stream.Messages[i].Values
		bytes, err := json.Marshal(values)
		if err != nil {
			log.Fatal(err)
		}

		rdb.XAck(
			context.Background(),
			stream.Stream,
			consumersGroup,
			messageID,
		)

		switch stream.Stream {
		case ADA:
			tx := domain.ResultLastTxADA{}
			err := json.Unmarshal(bytes, &tx)
			if err != nil {
				logger.LogError(err.Error())
			}
			// sen data to notification provider
			err = repo.SendNotification(tx)
			if err != nil {
				logger.LogError(err.Error())
			}
		case SOL:
			tx := domain.ResultLastTxSOL{}
			err := json.Unmarshal(bytes, &tx)
			if err != nil {
				logger.LogError(err.Error())
			}
			// sen data to notification provider
			err = repo.SendNotification(tx)
			if err != nil {
				logger.LogError(err.Error())
			}
		}

	}
}
