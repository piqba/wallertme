package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	domain "github.com/piqba/wallertme/internal/bb8/domain"
	"github.com/piqba/wallertme/pkg/constants"
	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/exporters"
	"github.com/piqba/wallertme/pkg/logger"
	"github.com/piqba/wallertme/pkg/web3"
)

const (
	TxSender   = "sender"
	TxReceiver = "receiver"
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

	// TODO: get from database with other metadata from wallets register
	wallets := []map[string]interface{}{
		{
			"address": "addr_test1qq6g6s99g9z9w0mlvew28w40lpml9rwfkfgerpkg6g2vpn6dp4cf7k9drrdy0wslarr6hxspcw8ev5ed8lfrmaengneqz34lcx",
			"lastTx":  "5ef1187f5e125090675a3c2d2d2cee359aaf6941df625db598ec996ab1011f55",
			"symbol":  "ADA",
		},
		{
			"address": "addr_test1qq5287luxzj5l4lequrqdp5ln76ver4uls3z0m5ykr5gqsv0vxzrwcq5dmmn9e09rvgttzgrngmpxkguy7220r0u0ljqzuww7g",
			"lastTx":  "5ef1187f5e125090675a3c2d2d2cee359aaf6941df625db598ec996ab1011f55",
			"symbol":  "ADA",
		},
		{
			"address": "9hZaTvCVMcfbheTzebkeGR6Xi2EzMqTtPasbhGoPB94C",
			"lastTx":  "3EDaSfApwCzkHcZdLBnMdDAyo9aVV9KaxCxSdmcMuJoq4sAoedb7ziHwBwBDe2jNxjnzZC5oAb9YFfGiHSs6taGu",
			"symbol":  "SOL",
		},
	}
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

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
			for _, w := range wallets {
				Exce(repo, w, exporterType)
			}
			if !watch {
				run = false
			}
		}
	}

}

func Exce(repo domain.TxRepository, wallet map[string]interface{}, exporterType string) {

	switch wallet["symbol"].(string) {
	case "ADA":
		addrInfo := lastTXForADA(wallet["address"].(string))

		lastTX := addrInfo.Result.CATxList[len(addrInfo.Result.CATxList)-1]
		tx := domain.ResultLastTxADA{
			Addr:          wallet["address"].(string),
			CtbID:         lastTX.CtbID,
			CtbTimeIssued: fmt.Sprintf("%d", lastTX.CtbTimeIssued),
			FromAddr:      lastTX.CtbOutputs[0].CtaAddress,
			ToAddr:        lastTX.CtbOutputs[1].CtaAddress,
			Balance:       addrInfo.Result.CABalance.GetCoin,
			Ammount:       lastTX.CtbOutputs[1].CtaAmount.GetCoin,
		}
		switch exporterType {
		case exporters.REDIS:

			lastTXFromRdb, err := repo.Get(context.TODO(), wallet["address"].(string))
			if err != nil {
				logger.LogError(err.Error())
			}
			if tx.CtbID != lastTXFromRdb {

				if wallet["address"].(string) == tx.FromAddr {
					tx.TypeTx = TxSender
				} else if wallet["address"].(string) == tx.ToAddr {
					tx.TypeTx = TxReceiver
				}

				err := repo.ExportData(tx)
				if err != nil {
					if errors.ErrorIs(err, exporters.ErrRedisXADDStreamID) {
						logger.LogWarn(fmt.Sprintf("This ID exist, NOT new TX for %s", tx.TruncateAddress(tx.Addr)))
						return
					}
					logger.LogError(err.Error())
				}
				err = repo.Set(context.TODO(), wallet["address"].(string), tx.CtbID, 0)
				if err != nil {
					logger.LogError(err.Error())
				}

			}
		case exporters.JSONFILE:
			err := repo.ExportData(tx)
			if err != nil {
				logger.LogError(err.Error())
			}
		}
	case "SOL":
		addrInfo := lastTXForSOL(wallet["address"].(string))
		tx := domain.ResultLastTxSOL{
			Addr:      wallet["address"].(string),
			TxID:      addrInfo.Result.Transaction.Signatures[0],
			Timestamp: fmt.Sprintf("%v", addrInfo.Result.BlockTime),
			FromAddr:  addrInfo.Result.Transaction.Message.AccountKeys[0],
			ToAddr:    addrInfo.Result.Transaction.Message.AccountKeys[1],
			Balance:   fmt.Sprintf("%v", addrInfo.Result.Meta.PostBalances[0]),
			Ammount:   fmt.Sprintf("%v", addrInfo.Result.Meta.PreBalances[0]-addrInfo.Result.Meta.PostBalances[0]),
		}
		if wallet["address"].(string) == tx.FromAddr {
			tx.TypeTx = TxSender
		} else if wallet["address"].(string) == tx.ToAddr {
			tx.TypeTx = TxReceiver
		}
		switch exporterType {
		case exporters.REDIS:

			lastTXFromRdb, err := repo.Get(context.TODO(), wallet["address"].(string))
			if err != nil {
				logger.LogError(err.Error())
			}
			if tx.TxID != lastTXFromRdb {

				if wallet["address"].(string) == tx.FromAddr {
					tx.TypeTx = TxSender
				} else if wallet["address"].(string) == tx.ToAddr {
					tx.TypeTx = TxReceiver
				}

				err := repo.ExportData(tx)
				if err != nil {
					if errors.ErrorIs(err, exporters.ErrRedisXADDStreamID) {
						logger.LogWarn(fmt.Sprintf("This ID exist, NOT new TX for %s", tx.TruncateAddress(tx.Addr)))
						return
					}
					logger.LogError(err.Error())
				}
				err = repo.Set(context.TODO(), wallet["address"].(string), tx.TxID, 0)
				if err != nil {
					logger.LogError(err.Error())
				}

			}
		case exporters.JSONFILE:
			err := repo.ExportData(tx)
			if err != nil {
				logger.LogError(err.Error())
			}
		}
	}

}

func lastTXForADA(address string) web3.AddrSumary {
	cardano, err := web3.NewAPICardanoClient(web3.APIClientOptions{})
	if err != nil {
		log.Fatal(err)
	}

	sumary, err := cardano.InfoByAddress(context.TODO(), address)
	if err != nil {
		log.Fatal(err)
	}

	return sumary
}
func lastTXForSOL(address string) web3.TxInfo {
	solanaApi, err := web3.NewAPISolanaClient(web3.APIClientOptions{})
	if err != nil {
		logger.LogError(errors.Errorf("main:%s", err).Error())
	}

	payloadLastTx := web3.PayloadReqJSONRPC{
		Jsonrpc: "2.0",
		Method:  "getSignaturesForAddress",
		Params: []interface{}{
			address,
			map[string]int{
				"limit": 1,
			},
		},
		ID: 1,
	}
	lastTx, err := solanaApi.LastTxByAddress(context.Background(), payloadLastTx)
	if err != nil {
		logger.LogError(errors.Errorf("main:%s", err).Error())
	}

	// get tx info by last signature
	payloadInfoTx := web3.PayloadReqJSONRPC{
		Jsonrpc: "2.0",
		Method:  "getTransaction",
		Params: []interface{}{
			lastTx.Result[0].Signature,
			"json",
		},
		ID: 1,
	}
	infoTx, err := solanaApi.InfoByTx(context.Background(), payloadInfoTx)
	if err != nil {
		logger.LogError(errors.Errorf("main:%s", err).Error())
	}
	return infoTx
}
