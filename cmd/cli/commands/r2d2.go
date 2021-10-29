package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	domain "github.com/piqba/wallertme/internal/r2d2/domain"
	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/exporters"
	"github.com/piqba/wallertme/pkg/logger"
	"github.com/piqba/wallertme/pkg/notify"
	"github.com/piqba/wallertme/pkg/storage"
	"github.com/spf13/cobra"
)

var (
	// seededRand random number
	// #nosec
	seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
)

const (
	// CHARSET of characters
	CHARSET = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789$*-/+"
	// SALT for define the size of the buffer
	SALT = 10
)

// StringWithCharset generate randoms words
func StringWithCharset() string {
	b := make([]byte, SALT)
	for i := range b {
		b[i] = CHARSET[seededRand.Intn(len(CHARSET))]
	}
	return string(b)
}

var consumerCmd = &cobra.Command{
	Use:   "r2d2",
	Short: "Subscribe to Txs data topic from (REDIS) and send notifications to this services (telegram|discord|smtp)",
	Run: func(cmd *cobra.Command, args []string) {
		// flags
		source, err := cmd.Flags().GetString(flagDataSource)
		if err != nil {
			logger.LogError(errors.Errorf("bb8: %v", err).Error())
		}
		groupName, err := cmd.Flags().GetString(flagConsumerGroup)
		if err != nil {
			logger.LogError(errors.Errorf("bb8: %v", err).Error())
		}
		walletsPath, err := cmd.Flags().GetString(flagWalletsPath)
		if err != nil {
			logger.LogError(errors.Errorf("bb8: %v", err).Error())
		}
		walletsName, err := cmd.Flags().GetString(flagWalletsName)
		if err != nil {
			logger.LogError(errors.Errorf("bb8: %v", err).Error())
		}
		// end flags

		// load wallets from source migrate to factory pattern
		dataSource := storage.NewSource(source, storage.OptionsSource{
			FileName: walletsName,
			PathName: walletsPath,
		})
		wallets, err := dataSource.WalletsTONotify()
		if err != nil {
			logger.LogError(errors.Errorf("bb8: %v", err).Error())
		}
		// end Load wallets

		// vars
		redisDbClient := exporters.GetRedisDbClient()

		streams := []string{}
		// create automatically streams key txs::addr
		for _, wallet := range wallets {
			streams = append(streams, fmt.Sprintf("%s::%s", "txs", wallet.Address))
		}

		var ids []string
		if groupName == "" {
			groupName = "r2d2-" + StringWithCharset()
		}
		groupName = "r2d2-" + groupName
		for _, v := range streams {
			ids = append(ids, ">")
			err := redisDbClient.XGroupCreate(context.TODO(), v, groupName, "0").Err()
			if err != nil {
				log.Println(err)
			}

		}

		streams = append(streams, ids...) // for each stream it requires an '>' :{"txs", ">"}

		for {
			entries, err := redisDbClient.XReadGroup(context.Background(), &redis.XReadGroupArgs{
				Group:    groupName,
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
				for _, wallet := range wallets {

					Exec(
						redisDbClient,
						groupName,
						it,
						wallet,
					)

				}
			}

		}

	},
}

func init() {
	consumerCmd.Flags().String(flagDataSource, "json", "select a wallets data source from (json|db)")
	consumerCmd.Flags().String(flagConsumerGroup, "", "select a name for your consumer group")
	consumerCmd.Flags().String(flagWalletsPath, walletsJsonPath, "select the path of wallet.json file")
	consumerCmd.Flags().String(flagWalletsName, "wallets.json", "select the name of wallet.json file")

	rootCmd.AddCommand(consumerCmd)

}

func Exec(
	rdb *redis.Client,
	consumersGroup string,
	stream redis.XStream,
	wallet storage.Wallet,
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
		if strings.Contains(string(bytes), "ADA") {

			tx := domain.ResultLastTxADA{}
			err := json.Unmarshal(bytes, &tx)
			if err != nil {
				logger.LogError(errors.Errorf("r2d2ctl: %v", err).Error())
			}
			if wallet.Address == tx.Addr {
				repo := getNotify(wallet)
				// sen data to notification provider
				err = repo.SendNotification(tx)
				if err != nil {
					logger.LogError(errors.Errorf("r2d2ctl: %v", err).Error())
				}
			}
		} else {
			tx := domain.ResultLastTxSOL{}
			err := json.Unmarshal(bytes, &tx)
			if err != nil {
				logger.LogError(errors.Errorf("r2d2ctl: %v", err).Error())
			}
			if wallet.Address == tx.Addr {
				repo := getNotify(wallet)
				// sen data to notification provider
				err = repo.SendNotification(tx)
				if err != nil {
					logger.LogError(errors.Errorf("r2d2ctl: %v", err).Error())
				}
			}
		}

	}
}

func getNotify(wallet storage.Wallet) domain.TxRepository {
	var repo domain.TxRepository

	// telegram client
	tgClientBot := notify.GetTgBot(notify.TgBotOption{
		Debug: false,
		Token: "",
	})

	// smtp client
	smtpClient := notify.NewSender(
		os.Getenv("SMTP_EMAIL_USER"),
		os.Getenv("SMTP_EMAIL_PASSWORD"),
	)

	for _, n := range wallet.NotifierService {

		switch n.Name {

		case notify.TELEGRAM:
			tgID, err := strconv.Atoi(n.UserID)
			if err != nil {
				logger.LogError(errors.Errorf("r2d2ctl: %v", err).Error())
			}
			repo = domain.NewTxRepository(
				domain.ExternalOptions{
					Type:              n.Name,
					DstNotificationID: int64(tgID),
				},
				tgClientBot,
			)
		case notify.DISCORD:
			// discord client
			discordClient, err := notify.NewDiscordClient(notify.DiscordClientOptions{
				ServerHook: n.UserID,
			})
			if err != nil {
				logger.LogError(errors.Errorf("r2d2ctl: %v", err).Error())
			}
			repo = domain.NewTxRepository(
				domain.ExternalOptions{
					Type: n.Name,
				},
				discordClient,
			)
		case notify.SMTP:
			repo = domain.NewTxRepository(domain.ExternalOptions{
				Type:              n.Name,
				DstNotificationID: 0,
			},
				smtpClient,
			)
		}
	}

	return repo

}
