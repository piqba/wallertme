package exporters

import (
	"context"
	"os"
	"time"

	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

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
func GetRedisDbClient(ctx context.Context) *redis.Client {
	_, span := otel.Tracer(nameRedisClient).Start(ctx, "GetRedisDbClient")
	defer span.End()
	clientInstance := redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_URI"),
		Username:     "",
		Password:     os.Getenv("REDIS_PASS"),
		DB:           0,
		DialTimeout:  60 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	})

	_, err := clientInstance.Ping(context.TODO()).Result()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.LogError(ErrRedisDbCheckConn.Error())
	}
	span.SetAttributes(attribute.String("create.redis.client", "Success"))

	return clientInstance
}
