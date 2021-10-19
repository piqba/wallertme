package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	domain "github.com/piqba/wallertme/internal/bb8/domain"
	"github.com/piqba/wallertme/pkg/constants"
	"github.com/piqba/wallertme/pkg/exporters"
	"github.com/piqba/wallertme/pkg/logger"
	"github.com/piqba/wallertme/pkg/web3"
)

const (
	addrSender   = "addr_test1qq6g6s99g9z9w0mlvew28w40lpml9rwfkfgerpkg6g2vpn6dp4cf7k9drrdy0wslarr6hxspcw8ev5ed8lfrmaengneqz34lcx"
	addrReceiver = "addr_test1qq5287luxzj5l4lequrqdp5ln76ver4uls3z0m5ykr5gqsv0vxzrwcq5dmmn9e09rvgttzgrngmpxkguy7220r0u0ljqzuww7g"
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

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

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

	run := true
	watch := true
	ds, err := time.ParseDuration("1s")
	if err != nil {
		logger.LogError("Failed to parse durations " + err.Error())
	}
	for run {
		select {
		case sig := <-quit:
			logger.LogInfo(fmt.Sprintf("server is shutting down %v", sig.String()))
			_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			run = false
			defer cancel()
		case <-time.Tick(ds):
			Exce(repo)
			if !watch {
				run = false
			}
		}
	}

}

func Exce(repo domain.TxRepository) {

	addrInfo := getTxByAddress(addrReceiver)
	lastTX := addrInfo.Result.CATxList[len(addrInfo.Result.CATxList)-1]

	tx := domain.ResultLastTxByAddr{
		Addr:          addrReceiver,
		CtbID:         lastTX.CtbID,
		CtbTimeIssued: lastTX.CtbTimeIssued,
		FromAddr:      lastTX.CtbOutputs[0].CtaAddress,
		ToAddr:        lastTX.CtbOutputs[1].CtaAddress,
		Balance:       addrInfo.Result.CABalance.GetCoin,
		Ammount:       lastTX.CtbOutputs[1].CtaAmount.GetCoin,
	}
	if addrReceiver == tx.FromAddr {
		tx.TypeTx = "sender"
	} else if addrReceiver == tx.ToAddr {
		tx.TypeTx = "receiver"
	}

	err := repo.ExportData(tx)
	if err != nil {
		if errors.Is(err, exporters.ErrRedisXADDStreamID) {
			logger.LogWarn(fmt.Sprintf("This ID exist, NOT new TX for %s", tx.TruncateAddress(tx.Addr)))
			return
		}
		logger.LogError(err.Error())
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
