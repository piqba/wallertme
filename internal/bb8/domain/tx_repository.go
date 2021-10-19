package bb8

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/piqba/wallertme/pkg/exporters"
	"github.com/piqba/wallertme/pkg/web3"
)

var (
	cwd, _ = os.Getwd()
)

type TxRepository struct {
	exporterType string
	clientRdb    *redis.Client
	clientPgx    *sqlx.DB
	clientKafka  *kafka.Producer
	limitTx      int
}

func NewTx(limit int) TxRepository {
	return TxRepository{limitTx: limit}
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

		return exporters.ExportToRedisStream(
			r.clientRdb,
			exporters.TXS_STREAM_KEY,
			fmt.Sprintf("%d", tx.CtbTimeIssued), // add this for ensuere that it`s unique ID based on TX timestamp
			value,
		)

	case exporters.KAFKA:

		value := tx.ToJSON()

		return exporters.ExportTokafka(r.clientKafka, exporters.TXS_TOPIC_KEY, value)

	case exporters.POSTGRESQL:

		value := tx.ToJSON()
		return exporters.ExportToPostgresql(r.clientPgx, int(tx.CtbTimeIssued), value)

	}

	return nil
}

func (r *TxRepository) InfoByAddress(address string) (ResultInfoByAddr, error) {

	cardano, err := web3.NewAPICardanoClient(web3.APIClientOptions{})
	if err != nil {
		return ResultInfoByAddr{}, err
	}

	sumary, err := cardano.SumaryAddrADA(context.TODO(), address)
	if err != nil {
		return ResultInfoByAddr{}, err
	}

	return ResultInfoByAddr{
		Address:   address,
		Type:      sumary.Result.CAType,
		BlockNO:   sumary.Result.CAChainTip.CTBlockNo,
		BlockHash: sumary.Result.CAChainTip.CTBlockHash,
		TxTotal:   sumary.Result.CATxNum,
		Balance:   sumary.Result.CABalance.GetCoin,
		TotalIn:   sumary.Result.CATotalInput.GetCoin,
		TotalOut:  sumary.Result.CATotalOutput.GetCoin,
		TotalFee:  sumary.Result.CATotalFee.GetCoin,
		TxList:    sumary.Result.CATxList,
	}, nil
}
