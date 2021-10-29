package exporters

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-redis/redis/v8"
	"github.com/piqba/wallertme/pkg/logger"
)

const (
	// JSONFILE ...
	JSONFILE = "json"
	// REDIS ...
	REDIS = "redis"
)

// Exporter ...
type Exporter interface {
	ExportData(data interface{}) error
}

// ExportToJSON ...
func ExportToJSON(path, filename string, value interface{}) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	file, err := os.OpenFile(filepath.Join(path, filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
		return err
	}

	datawriter := bufio.NewWriter(file)

	switch line := value.(type) {
	case string:
		_, err = datawriter.WriteString(line + "\n")
		if err != nil {
			logger.LogError("failed to writer string: " + err.Error())
			return err
		}
	}

	datawriter.Flush()
	file.Close()
	return nil
}

// ExportToRedisStream ...
func ExportToRedisStream(rdb *redis.Client, key, address string, value map[string]interface{}) error {

	err := rdb.XAdd(context.TODO(), &redis.XAddArgs{
		Stream: fmt.Sprintf("%s::%s", key, address),
		Values: value,
	}).Err()

	if err != nil {
		if err.Error() == ErrRedisXADDStreamID.Error() {
			return ErrRedisXADDStreamID
		}
		return err
	}

	return nil
}
