package exporters

import (
	"context"
	"os"
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
	// ErrRedisXADDStreamID ...
	ErrRedisXADDStreamID = errors.NewError("ERR The ID specified in XADD is equal or smaller than the target stream top item")
)

// GetRedisDbClient ...
func GetRedisDbClient() *redis.Client {

	clientInstance := redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_URI"),  // use default Addr
		Password:     os.Getenv("REDIS_PASS"), // no password set
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
