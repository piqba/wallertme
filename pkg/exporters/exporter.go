package exporters

import (
	"bufio"
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/go-redis/redis/v8"
	"github.com/piqba/wallertme/pkg/logger"
)

const (
	// JSONFILE ...
	JSONFILE = "json"
	// KAFKA ...
	// KAFKA = "kafka"
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
func ExportToRedisStream(rdb *redis.Client, key, symbol string, value map[string]interface{}) error {

	err := rdb.XAdd(context.TODO(), &redis.XAddArgs{
		Stream: key + ":" + symbol,
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

// ExportTokafka ...
// Was commented for developments reasson.
// func ExportTokafka(p *kafka.Producer, topic string, value string) error {
// 	var err error
// 	// Delivery report handler for produced messages
// 	go func() {
// 		for e := range p.Events() {
// 			switch ev := e.(type) {
// 			case *kafka.Message:
// 				if ev.TopicPartition.Error != nil {
// 					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
// 					err = errors.Errorf("Delivery failed: %v\n", ev.TopicPartition)
// 				}
// 			}
// 		}
// 	}()

// 	// Produce messages to topic (asynchronously)
// 	p.Produce(&kafka.Message{
// 		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
// 		Value:          []byte(value),
// 	}, nil)

// 	// Wait for message deliveries before shutting down
// 	p.Flush(15 * 1000)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
