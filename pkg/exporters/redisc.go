package exporters

import (
	"context"
	"time"

	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/logger"

	"github.com/go-redis/redis/v8"
)

const (
	// TXS_STREAM_KEY ...
	TXS_STREAM_KEY = "txs"
)

var (
	// ErrRedisDbCheckConn ...
	ErrRedisDbCheckConn = errors.NewError("Redis: Fail to check connection")
	// RedisDbClient ...
	RedisDbClient = GetRedisDbClient()
)

// GetRedisDbClient ...
func GetRedisDbClient() *redis.Client {

	clientInstance := redis.NewClient(&redis.Options{
		Addr:         ":6379", // use default Addr
		Password:     "",      // no password set
		DB:           0,
		DialTimeout:  60 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	})

	_, err := clientInstance.Ping(context.TODO()).Result()
	if err != nil {
		logger.LogError(ErrRedisDbCheckConn.Error())
	}
	return clientInstance
}
