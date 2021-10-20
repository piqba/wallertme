package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	domain "github.com/piqba/wallertme/internal/r2d2/domain"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/piqba/wallertme/pkg/exporters"
	"github.com/piqba/wallertme/pkg/logger"
	"github.com/piqba/wallertme/pkg/notify"
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
	redisDbClient := exporters.GetRedisDbClient()
	tgClientBot := notify.GetTgBot(notify.TgBotOption{
		Debug: false,
		Token: "",
	})

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
	}

	streamADA := exporters.TXS_STREAM_KEY + ":" + "ADA"
	streams := []string{streamADA}
	var ids []string
	consumersGroup := "cardano-group"
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

		switch entries[0].Stream {
		case streamADA:
			Exec(redisDbClient, consumersGroup, entries[0], repo)
		}

	}
}

func Exec(rdb *redis.Client, consumersGroup string, stream redis.XStream, repo domain.TxRepository) {
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
		// sen data to notification provider
		err = repo.SendNotification(string(bytes))
		if err != nil {
			logger.LogError(err.Error())
		}
	}
}
