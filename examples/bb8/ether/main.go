package main

import (
	"context"
	"log"
	"runtime"

	"github.com/joho/godotenv"
	domain "github.com/piqba/wallertme/internal/bb8/domain"
	"github.com/piqba/wallertme/pkg/constants"
	"github.com/piqba/wallertme/pkg/exporters"
	"github.com/piqba/wallertme/pkg/logger"
	"github.com/piqba/wallertme/pkg/web3"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.LogError(err.Error())
	}
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu) // Try to use all available CPUs.
}

func main() {
	exporterType := exporters.REDIS
	var repo domain.TxRepository
	switch exporterType {
	case exporters.REDIS:
		rdb := exporters.GetRedisDbClient()
		repo = domain.NewTxRepository(exporters.REDIS, rdb)
	case exporters.JSONFILE:
		repo = domain.NewTxRepository(exporters.JSONFILE, nil)
	case exporters.POSTGRESQL:
		pg, err := exporters.PostgreSQLConnection()
		if err != nil {
			logger.LogError(err.Error())
		}
		pg.MustExec(constants.SchemaTXS)
		repo = domain.NewTxRepository(exporters.POSTGRESQL, pg)
	case exporters.KAFKA:
		pk := exporters.GetProducerClientKafka()
		repo = domain.NewTxRepository(exporters.KAFKA, pk)
	}

	block, err := getTxByLatestBlock()
	if err != nil {
		log.Fatal(err)
	}

	for _, txInBlock := range block.Result.Transactions {

		tx := domain.ResultTx{
			Time:          block.Result.Timestamp,
			Txfrom:        txInBlock.From,
			Txto:          txInBlock.To,
			Gas:           txInBlock.Gas,
			Gasprice:      txInBlock.GasPrice,
			Block:         block.Result.Number,
			Txhash:        txInBlock.Hash,
			Value:         txInBlock.Value,
			ContractTo:    "",
			ContractValue: "",
		}

		err := repo.ExportData(tx)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func getTxByLatestBlock() (web3.Block, error) {
	api, err := web3.NewAPIEthClient(
		web3.APIClientOptions{
			Server: web3.GanacheDevNet,
		},
	)
	if err != nil {
		return web3.Block{}, err
	}
	payload := web3.PayloadReq{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params: []interface{}{
			"0x2",
			true,
		},
		ID: 1,
	}

	blc, err := api.BlockByNumber(context.TODO(), payload)
	if err != nil {
		log.Fatal(err)
		return web3.Block{}, err
	}
	return blc, nil
}
