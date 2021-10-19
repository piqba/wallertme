package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/piqba/wallertme/pkg/exporters"
	"github.com/piqba/wallertme/pkg/logger"
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

	streams := []string{exporters.TXS_STREAM_KEY}
	var ids []string
	consumersGroup := "cardano-consumer-group"
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
		case "txs":
			printStream(redisDbClient, consumersGroup, entries[0])
		}

	}
}

func printStream(rdb *redis.Client, consumersGroup string, stream redis.XStream) {
	for i := 0; i < len(stream.Messages); i++ {
		messageID := stream.Messages[i].ID
		values := stream.Messages[i].Values
		bytes, err := json.Marshal(values)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(bytes))

		rdb.XAck(
			context.Background(),
			stream.Stream,
			consumersGroup,
			messageID,
		)
	}
}
