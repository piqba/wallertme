package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	domain "github.com/piqba/wallertme/internal/r2d2/domain"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/piqba/wallertme/pkg/exporters"
	"github.com/piqba/wallertme/pkg/logger"
	"github.com/piqba/wallertme/pkg/notify"
)

var (
	ADA = exporters.TXS_STREAM_KEY + "::" + "addr_test1qq6g6s99g9z9w0mlvew28w40lpml9rwfkfgerpkg6g2vpn6dp4cf7k9drrdy0wslarr6hxspcw8ev5ed8lfrmaengneqz34lcx"
	SOL = exporters.TXS_STREAM_KEY + "::" + "9hZaTvCVMcfbheTzebkeGR6Xi2EzMqTtPasbhGoPB94C"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.LogError(err.Error())
	}
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu) // Try to use all available CPUs.
}

func main() {
	log.Println("Consumer started")
	redisDbClient := exporters.GetRedisDbClient(context.TODO())
	// telegram client
	tgClientBot := notify.GetTgBot(context.TODO(), notify.TgBotOption{
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
	notificationType := notify.TELEGRAM

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
			Exec(context.TODO(), redisDbClient, consumersGroup, it, repo)
		}

	}
}

func Exec(
	ctx context.Context,
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
			err = repo.SendNotification(ctx, tx)
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
			err = repo.SendNotification(ctx, tx)
			if err != nil {
				logger.LogError(err.Error())
			}
		}

	}
}
