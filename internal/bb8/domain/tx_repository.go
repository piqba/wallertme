package bb8

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/piqba/wallertme/pkg/exporters"
	"github.com/piqba/wallertme/pkg/web3"
)

var (
	cwd, _ = os.Getwd()
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
func (r *TxRepository) ExportData(data interface{}, symbol string) error {

	t := reflect.TypeOf(data)
	if t == reflect.TypeOf(ResultLastTxADA{}) {

		tx := data.(ResultLastTxADA)

		switch r.exporterType {
		case exporters.JSONFILE:
			fileName := fmt.Sprintf("export_%s_%d.json", symbol, time.Now().UnixNano())

			return exporters.ExportToJSON(cwd, fileName, tx.ToJSON())
		case exporters.REDIS:

			value, err := tx.ToMAP()
			if err != nil {
				return err
			}

			return exporters.ExportToRedisStream(
				r.clientRdb,
				exporters.TXS_STREAM_KEY,
				symbol,
				value,
			)
		}
	} else if t == reflect.TypeOf(ResultLastTxSOL{}) {
		tx := data.(ResultLastTxSOL)

		switch r.exporterType {
		case exporters.JSONFILE:
			fileName := fmt.Sprintf("export_%s_%d.json", symbol, time.Now().UnixNano())

			return exporters.ExportToJSON(cwd, fileName, tx.ToJSON())
		case exporters.REDIS:

			value, err := tx.ToMAP()
			if err != nil {
				return err
			}

			return exporters.ExportToRedisStream(
				r.clientRdb,
				exporters.TXS_STREAM_KEY,
				symbol,
				value,
			)

		}
	}

	return nil
}

// InfoByAddress get info by address
func (r *TxRepository) InfoByAddress(address string) (ResultInfoForADA, error) {

	cardano, err := web3.NewAPICardanoClient(web3.APIClientOptions{})
	if err != nil {
		return ResultInfoForADA{}, err
	}

	sumary, err := cardano.InfoByAddress(context.TODO(), address)
	if err != nil {
		return ResultInfoForADA{}, err
	}

	return ResultInfoForADA{
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

// Set an address as a key and the last TX as a value
func (r *TxRepository) Set(ctx context.Context, address, lastTx string, expiration time.Duration) error {
	// err = redisdb.Set("key", "value", 0).Err() never expire
	// err = redisdb.Set("key", "value", time.Hour).Err()
	err := r.clientRdb.Set(ctx, address, lastTx, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get the last tx value  by address as a key
func (r *TxRepository) Get(ctx context.Context, address string) (string, error) {
	lastTx, err := r.clientRdb.Get(ctx, address).Result()
	if err != nil {
		return "", err
	}
	return lastTx, nil
}
