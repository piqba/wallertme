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
	exporterType := exporters.JSONFILE
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
	address := "addr_test1qq6g6s99g9z9w0mlvew28w40lpml9rwfkfgerpkg6g2vpn6dp4cf7k9drrdy0wslarr6hxspcw8ev5ed8lfrmaengneqz34lcx"
	addrInfo := getTxByAddress(address)
	tx := domain.ResultTxADA{
		Address:   addrInfo.Result.CAAddress,
		Type:      addrInfo.Result.CAType,
		BlockNO:   addrInfo.Result.CAChainTip.CTBlockNo,
		BlockHash: addrInfo.Result.CAChainTip.CTBlockHash,
		TxTotal:   addrInfo.Result.CATxNum,
		Balance:   addrInfo.Result.CABalance.GetCoin,
		TotalIn:   addrInfo.Result.CATotalInput.GetCoin,
		TotalOut:  addrInfo.Result.CATotalOutput.GetCoin,
		TotalFee:  addrInfo.Result.CATotalFee.GetCoin,
		TxList:    addrInfo.Result.CATxList,
	}

	err := repo.ExportData(tx)
	if err != nil {
		log.Fatal(err)
	}
}

func getTxByAddress(address string) web3.AddrSumary {
	cardano, err := web3.NewAPICardanoClient(web3.APIClientOptions{})
	if err != nil {
		log.Fatal(err)
	}

	sumary, err := cardano.SumaryAddrADA(context.TODO(), address)
	if err != nil {
		log.Fatal(err)
	}

	return sumary
}
