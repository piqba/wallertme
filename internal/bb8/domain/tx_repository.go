package bb8

import (
	"context"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/piqba/wallertme/pkg/exporters"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// TxRepository repository
type TxRepository struct {
	exporterType string
	clientRdb    *redis.Client
	clientPgx    *sqlx.DB
	// clientKafka  *kafka.Producer
	limitTx int
}

// NewTx constructor
func NewTx(limit int) TxRepository {
	return TxRepository{limitTx: limit}
}

// NewTxRepository constructor client
func NewTxRepository(exporterType string, clients ...interface{}) TxRepository {
	var repo TxRepository
	repo.exporterType = exporterType
	for _, c := range clients {
		switch c := c.(type) {
		case *redis.Client:
			repo.clientRdb = c
		case *sqlx.DB:
			repo.clientPgx = c

		default:
			return repo
		}
	}
	return repo
}

// ExportData export data from ADA or SOL symbol
func (r *TxRepository) ExportData(ctx context.Context, data interface{}) error {
	_, span := otel.Tracer(nameBb8).Start(ctx, "ExportData")
	defer span.End()
	t := reflect.TypeOf(data)
	if t == reflect.TypeOf(ResultLastTxADA{}) {

		tx := data.(ResultLastTxADA)

		value, err := tx.ToMAP()
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		span.SetAttributes(attribute.String("bb8.redis.stream.ada", "Success"))

		return exporters.ExportToRedisStream(
			ctx,
			r.clientRdb,
			exporters.TXS_STREAM_KEY,
			tx.Addr,
			value,
		)

	} else if t == reflect.TypeOf(ResultLastTxSOL{}) {
		tx := data.(ResultLastTxSOL)

		value, err := tx.ToMAP()
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		span.SetAttributes(attribute.String("bb8.redis.stream.sol", "Success"))

		return exporters.ExportToRedisStream(
			ctx,
			r.clientRdb,
			exporters.TXS_STREAM_KEY,
			tx.Addr,
			value,
		)

	}

	return nil
}

// Set an address as a key and the last TX as a value
func (r *TxRepository) Set(ctx context.Context, address, lastTx string, expiration time.Duration) error {
	_, span := otel.Tracer(nameBb8).Start(ctx, "Set")
	defer span.End()
	// err = redisdb.Set("key", "value", 0).Err() never expire
	// err = redisdb.Set("key", "value", time.Hour).Err()
	err := r.clientRdb.Set(ctx, address, lastTx, expiration).Err()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	span.SetAttributes(attribute.String("bb8.domain.Set.redis", "Success"))

	return nil
}

// Get the last tx value  by address as a key
func (r *TxRepository) Get(ctx context.Context, address string) (string, error) {
	_, span := otel.Tracer(nameBb8).Start(ctx, "Get")
	defer span.End()
	lastTx, err := r.clientRdb.Get(ctx, address).Result()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}
	span.SetAttributes(attribute.String("bb8.domain.Get.redis", "Success"))

	return lastTx, nil
}
