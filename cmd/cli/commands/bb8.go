package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	domain "github.com/piqba/wallertme/internal/bb8/domain"

	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/exporters"
	"github.com/piqba/wallertme/pkg/logger"
	"github.com/piqba/wallertme/pkg/web3"
	"github.com/spf13/cobra"
)

const (
	TxSender   = "sender"
	TxReceiver = "receiver"
)

var walletsJsonPath, _ = os.Getwd()

var producerCmd = &cobra.Command{
	Use:   "bb8",
	Short: "Publish Txs data from (SOLANA|CARDANO) blockchains to (REDIS)",
	Run: func(cmd *cobra.Command, args []string) {
		exporter, err := cmd.Flags().GetString(flagExporter)
		if err != nil {
			logger.LogError(errors.Errorf("bb8: %v", err).Error())
		}
		timer, err := cmd.Flags().GetString(flagTimer)
		if err != nil {
			logger.LogError(errors.Errorf("bb8: %v", err).Error())
		}
		watcher, err := cmd.Flags().GetBool(flagWatcher)
		if err != nil {
			logger.LogError(errors.Errorf("bb8: %v", err).Error())
		}

		walletsPath, err := cmd.Flags().GetString(flagWalletsPath)
		if err != nil {
			logger.LogError(errors.Errorf("bb8: %v", err).Error())
		}
		walletsName, err := cmd.Flags().GetString(flagWalletsName)
		if err != nil {
			logger.LogError(errors.Errorf("bb8: %v", err).Error())
		}

		wallets, err := readWalletsJsonFile(walletsPath, walletsName)
		if err != nil {
			logger.LogError(errors.Errorf("bb8: %v", err).Error())
		}

		quit := make(chan os.Signal, 1)

		signal.Notify(quit, os.Interrupt)

		exporterType := exporter

		var repo domain.TxRepository
		switch exporterType {
		case exporters.REDIS:
			rdb := exporters.GetRedisDbClient()
			repo = domain.NewTxRepository(exporters.REDIS, rdb)
		case exporters.JSONFILE:
			repo = domain.NewTxRepository(exporters.JSONFILE, nil)

		}

		run := true

		ds, err := time.ParseDuration(timer)
		if err != nil {
			logger.LogError(errors.Errorf("bb8: %v", err).Error())
		}
		for run {
			select {
			case sig := <-quit:
				logger.LogInfo(fmt.Sprintf("bb8: app is shutting down %v", sig.String()))
				_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				run = false
				defer cancel()
			case <-time.Tick(ds):
				for _, w := range wallets {
					Exce(repo, w, exporterType)
				}
				if !watcher {
					logger.LogInfo(fmt.Sprintf("bb8: app is shutting down %v", "closeapp"))
					run = false
				}
			}
		}
	},
}

func init() {
	producerCmd.Flags().String(flagExporter, exporters.REDIS, "select a exporter to send data")
	producerCmd.Flags().String(flagTimer, "1s", "select a time duration to watch all txs")
	producerCmd.Flags().Bool(flagWatcher, false, "select true|false if you want to run this task periodicaly")
	producerCmd.Flags().String(flagWalletsPath, walletsJsonPath, "select the path of wallet.json file")
	producerCmd.Flags().String(flagWalletsName, "wallets.json", "select the name of wallet.json file")
	rootCmd.AddCommand(producerCmd)

}

func readWalletsJsonFile(path, filename string) ([]map[string]interface{}, error) {
	jsonFile, err := os.Open(fmt.Sprintf("%s/%s", path, filename))
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Exce(repo domain.TxRepository, wallet map[string]interface{}, exporterType string) {
	symbol := wallet["symbol"].(string)
	switch symbol {
	case "ADA":
		addrInfo := lastTXForADA(wallet["networkType"].(string), wallet["address"].(string))

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
				logger.LogError(errors.Errorf("main:%s", err).Error())
			}
			if tx.CtbID != lastTXFromRdb {

				if wallet["address"].(string) == tx.FromAddr {
					tx.TypeTx = TxSender
				} else if wallet["address"].(string) == tx.ToAddr {
					tx.TypeTx = TxReceiver
				}

				err := repo.ExportData(tx, symbol)
				if err != nil {
					if errors.ErrorIs(err, exporters.ErrRedisXADDStreamID) {
						logger.LogWarn(fmt.Sprintf("This ID exist, NOT new TX for %s", tx.TruncateAddress(tx.Addr)))
						return
					}
					logger.LogError(errors.Errorf("main:%s", err).Error())
				}
				err = repo.Set(context.TODO(), wallet["address"].(string), tx.CtbID, 0)
				if err != nil {
					logger.LogError(errors.Errorf("main:%s", err).Error())
				}

			}
		case exporters.JSONFILE:
			err := repo.ExportData(tx, symbol)
			if err != nil {
				logger.LogError(errors.Errorf("main:%s", err).Error())
			}
		}
	case "SOL":
		addrInfo := lastTXForSOL(wallet["networkType"].(string), wallet["address"].(string))
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
				logger.LogError(errors.Errorf("main:%s", err).Error())
			}
			if tx.TxID != lastTXFromRdb {

				if wallet["address"].(string) == tx.FromAddr {
					tx.TypeTx = TxSender
				} else if wallet["address"].(string) == tx.ToAddr {
					tx.TypeTx = TxReceiver
				}

				err := repo.ExportData(tx, symbol)
				if err != nil {
					if errors.ErrorIs(err, exporters.ErrRedisXADDStreamID) {
						logger.LogWarn(errors.Errorf("This ID exist, NOT new TX for %s", tx.TruncateAddress(tx.Addr), err).Error())
						return
					}
					logger.LogError(errors.Errorf("main:%s", err).Error())
				}
				err = repo.Set(context.TODO(), wallet["address"].(string), tx.TxID, 0)
				if err != nil {
					logger.LogError(errors.Errorf("main:%s", err).Error())
				}

			}
		case exporters.JSONFILE:
			err := repo.ExportData(tx, symbol)
			if err != nil {
				logger.LogError(errors.Errorf("main:%s", err).Error())
			}
		}
	}

}

func lastTXForADA(networkType, address string) web3.AddrSumary {
	cardano, err := web3.NewAPICardanoClient(web3.APIClientOptions{
		NetworkType: networkType,
	})
	if err != nil {
		logger.LogError(errors.Errorf("main:%s", err).Error())
	}

	sumary, err := cardano.InfoByAddress(context.TODO(), address)
	if err != nil {
		logger.LogError(errors.Errorf("main:%s", err).Error())
	}

	return sumary
}
func lastTXForSOL(networkType, address string) web3.TxInfo {
	solanaApi, err := web3.NewAPISolanaClient(web3.APIClientOptions{
		NetworkType: networkType,
	})
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
