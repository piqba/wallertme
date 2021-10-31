package exporters

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-redis/redis/v8"
	"github.com/piqba/wallertme/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

const (
	// JSONFILE ...
	JSONFILE = "json"
	// REDIS ...
	REDIS = "redis"
	// nameRedisClient is the Tracer nameRedisClient used to identify this instrumentation library.
	nameRedisClient = "exporter.redis.client"
	// nameExporterJson is the Tracer nameExporterJson used to identify this instrumentation library.
	nameExporterJson = "exporter.jsonfile"
	// nameExporterRedis is the Tracer nameExporterRedis used to identify this instrumentation library.
	nameExporterRedis = "exporter.redis"
)

// Exporter ...
type Exporter interface {
	ExportData(ctx context.Context, data interface{}) error
}

// ExportToJSON ...
func ExportToJSON(ctx context.Context, path, filename string, value interface{}) error {
	_, span := otel.Tracer(nameExporterJson).Start(ctx, "ExportToJSON")
	defer span.End()
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	file, err := os.OpenFile(filepath.Join(path, filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.Errorf("failed creating file: %v", err)
	}

	datawriter := bufio.NewWriter(file)

	switch line := value.(type) {
	case string:
		_, err = datawriter.WriteString(line + "\n")
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return errors.Errorf("failed to writer string: %v", err)
		}
	}

	datawriter.Flush()
	file.Close()
	span.SetAttributes(attribute.String("exporter.jsonfile", "Success"))

	return nil
}

// ExportToRedisStream ...
func ExportToRedisStream(ctx context.Context, rdb *redis.Client, key, address string, value map[string]interface{}) error {
	_, span := otel.Tracer(nameExporterRedis).Start(ctx, "ExportToRedisStream")
	defer span.End()
	err := rdb.XAdd(context.TODO(), &redis.XAddArgs{
		Stream: fmt.Sprintf("%s::%s", key, address),
		Values: value,
	}).Err()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if err.Error() == ErrRedisXADDStreamID.Error() {
			return ErrRedisXADDStreamID
		}
		return err
	}
	span.SetAttributes(attribute.String("exporter.redis.stream", "Success"))

	return nil
}
