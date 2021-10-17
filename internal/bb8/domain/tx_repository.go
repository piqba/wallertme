package bb8

import (
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/piqba/wallertme/pkg/exporters"
)

var (
	cwd, _ = os.Getwd()
)

type TxRepository struct {
	exporterType string
	clientRdb    *redis.Client
	clientPgx    *sqlx.DB
	clientKafka  *kafka.Producer
}

func NewTxRepository(exporterType string, clients ...interface{}) TxRepository {
	var repo TxRepository
	repo.exporterType = exporterType
	for _, c := range clients {
		switch c := c.(type) {
		case *redis.Client:
			repo.clientRdb = c
		case *sqlx.DB:
			repo.clientPgx = c
		case *kafka.Producer:
			repo.clientKafka = c
		default:
			return repo
		}
	}
	return repo
}

func (r *TxRepository) ExportData(data interface{}) error {
	tx := data.(ResultLastTxByAddr)

	switch r.exporterType {
	case exporters.JSONFILE:
		fileName := fmt.Sprintf("export_%d.json", time.Now().UnixNano())

		return exporters.ExportToJSON(cwd, fileName, tx.ToJSON())
	case exporters.REDIS:

		value, err := tx.ToMAP()
		if err != nil {
			return err
		}

		return exporters.ExportToRedisStream(r.clientRdb, exporters.TXS_STREAM_KEY, value)

	case exporters.KAFKA:

		value := tx.ToJSON()

		return exporters.ExportTokafka(r.clientKafka, exporters.TXS_TOPIC_KEY, value)

	case exporters.POSTGRESQL:

		value := tx.ToJSON()
		return exporters.ExportToPostgresql(r.clientPgx, int(tx.CtbTimeIssued), value)

	}

	return nil
}
