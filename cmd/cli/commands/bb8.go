package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	domain "github.com/piqba/wallertme/internal/bb8/domain"

	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/exporters"
	"github.com/piqba/wallertme/pkg/logger"
	"github.com/piqba/wallertme/pkg/storage"
	"github.com/piqba/wallertme/pkg/web3"
	"github.com/spf13/cobra"
)

const (
	TxSender   = "sender"
	TxReceiver = "receiver"
)

var (
	walletsJsonPath, _ = os.Getwd()
)

var producerCmd = &cobra.Command{
	Use:   "bb8",
	Short: "Publish Txs data from (SOLANA|CARDANO) blockchains to (REDIS)",
	Run: func(cmd *cobra.Command, args []string) {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		// flags
		source, err := cmd.Flags().GetString(flagSource)
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
		// END FLAGS

		// load wallets from source migrate to factory pattern
		pgx, err := storage.PostgreSQLConnection()
		if err != nil {
			logger.LogError(errors.Errorf("bb8: %v", err).Error())

		}
		dataSource := storage.NewSource(source, storage.OptionsSource{
			FileName: walletsName,
			PathName: walletsPath,
			Pgx:      pgx,
		})
		wallets, err := dataSource.Wallets()
		if err != nil {
			logger.LogError(errors.Errorf("bb8: %v", err).Error())
		}
		// end Load wallets

		// Define repository to export data
		var repo domain.TxRepository
		rdb := exporters.GetRedisDbClient()
		repo = domain.NewTxRepository(exporters.REDIS, rdb)

		// Logic for watcher periodicaly
		// and crtl+c
		run := true

		// parse timer duration flag
		ds, err := time.ParseDuration(timer)
		if err != nil {
			logger.LogError(errors.Errorf("bb8: %v", err).Error())
		}
		// run periodicaly if watcher flag is true
		for run {
			select {
			case sig := <-quit:

				logger.LogInfo(fmt.Sprintf("bb8: app is shutting down %v", sig.String()))
				_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				run = false
				defer cancel()

			case <-time.Tick(ds):

				var wg sync.WaitGroup
				wg.Add(len(wallets))
				for _, it := range wallets {
					go func(wallet storage.Wallet) {
						Exce(repo, wallet)
						wg.Done()
					}(it)
				}
				wg.Wait()

				if !watcher {
					logger.LogInfo(fmt.Sprintf("bb8: app is shutting down %v", "closeapp"))
					run = false
				}

			}
		}
	},
}

func init() {
	producerCmd.Flags().String(flagSource, "json", "select a wallets data source from (json|db)")
	producerCmd.Flags().String(flagTimer, "1s", "select a time duration to watch all txs")
	producerCmd.Flags().Bool(flagWatcher, false, "select true|false if you want to run this task periodicaly")
	producerCmd.Flags().String(flagWalletsPath, walletsJsonPath, "select the path of wallet.json file")
	producerCmd.Flags().String(flagWalletsName, "wallets.json", "select the name of wallet.json file")
	rootCmd.AddCommand(producerCmd)

}

// Exce excecute proccessing functions
func Exce(repo domain.TxRepository, wallet storage.Wallet) {

	symbol := wallet.Symbol
	address := wallet.Address
	networkType := wallet.NetworkType
	isActive := wallet.IsActive

	if isActive {

		switch symbol {
		case "ADA":
			addrInfo := lastTXForADA(networkType, address)

			lastTX := addrInfo.Result.CATxList[len(addrInfo.Result.CATxList)-1]
			tx := domain.ResultLastTxADA{
				Addr:          address,
				CtbID:         lastTX.CtbID,
				CtbTimeIssued: fmt.Sprintf("%d", lastTX.CtbTimeIssued),
				FromAddr:      lastTX.CtbOutputs[0].CtaAddress,
				ToAddr:        lastTX.CtbOutputs[1].CtaAddress,
				Balance:       addrInfo.Result.CABalance.GetCoin,
				Ammount:       lastTX.CtbOutputs[1].CtaAmount.GetCoin,
				Symbol:        symbol,
			}

			lastTXFromRdb, err := repo.Get(context.TODO(), address)
			if err != nil {
				logger.LogError(errors.Errorf("bb8: %s", err).Error())
			}
			if tx.CtbID != lastTXFromRdb {

				if address == tx.FromAddr {
					tx.TypeTx = TxSender
				} else if address == tx.ToAddr {
					tx.TypeTx = TxReceiver
				}

				err := repo.ExportData(tx, symbol)
				if err != nil {
					if errors.ErrorIs(err, exporters.ErrRedisXADDStreamID) {
						logger.LogWarn(fmt.Sprintf("This ID exist, NOT new TX for %s", tx.TruncateAddress(tx.Addr)))
						return
					}
					logger.LogError(errors.Errorf("bb8: %s", err).Error())
				}
				err = repo.Set(context.TODO(), address, tx.CtbID, 0)
				if err != nil {
					logger.LogError(errors.Errorf("bb8: %s", err).Error())
				}

			}

		case "SOL":
			addrInfo := lastTXForSOL(networkType, address)
			tx := domain.ResultLastTxSOL{
				Addr:      address,
				TxID:      addrInfo.Result.Transaction.Signatures[0],
				Timestamp: fmt.Sprintf("%v", addrInfo.Result.BlockTime),
				FromAddr:  addrInfo.Result.Transaction.Message.AccountKeys[0],
				ToAddr:    addrInfo.Result.Transaction.Message.AccountKeys[1],
				Balance:   fmt.Sprintf("%v", addrInfo.Result.Meta.PostBalances[0]),
				Ammount:   fmt.Sprintf("%v", addrInfo.Result.Meta.PreBalances[0]-addrInfo.Result.Meta.PostBalances[0]),
				Symbol:    symbol,
			}
			if address == tx.FromAddr {
				tx.TypeTx = TxSender
			} else if address == tx.ToAddr {
				tx.TypeTx = TxReceiver
			}

			lastTXFromRdb, err := repo.Get(context.TODO(), address)
			if err != nil {
				logger.LogError(errors.Errorf("bb8: %s", err).Error())
			}
			if tx.TxID != lastTXFromRdb {

				if address == tx.FromAddr {
					tx.TypeTx = TxSender
				} else if address == tx.ToAddr {
					tx.TypeTx = TxReceiver
				}

				err := repo.ExportData(tx, symbol)
				if err != nil {
					if errors.ErrorIs(err, exporters.ErrRedisXADDStreamID) {
						logger.LogWarn(errors.Errorf("This ID exist, NOT new TX for %s", tx.TruncateAddress(tx.Addr), err).Error())
						return
					}
					logger.LogError(errors.Errorf("bb8: %s", err).Error())
				}
				err = repo.Set(context.TODO(), address, tx.TxID, 0)
				if err != nil {
					logger.LogError(errors.Errorf("bb8: %s", err).Error())
				}

			}

		}
	} else {
		// isActive false
	}

}

// lastTXForADA return las TX from Cardano
func lastTXForADA(networkType, address string) web3.AddrSumary {
	cardano, err := web3.NewAPICardanoClient(web3.APIClientOptions{
		NetworkType: networkType,
	})
	if err != nil {
		logger.LogError(errors.Errorf("bb8: %s", err).Error())
	}

	sumary, err := cardano.InfoByAddress(context.TODO(), address)
	if err != nil {
		logger.LogError(errors.Errorf("bb8: %s", err).Error())
	}

	return sumary
}

// lastTXForSOL return las TX from Solana
func lastTXForSOL(networkType, address string) web3.TxInfo {
	solanaApi, err := web3.NewAPISolanaClient(web3.APIClientOptions{
		NetworkType: networkType,
	})
	if err != nil {
		logger.LogError(errors.Errorf("bb8: %s", err).Error())
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
		logger.LogError(errors.Errorf("bb8: %s", err).Error())
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
		logger.LogError(errors.Errorf("bb8: %s", err).Error())
	}
	return infoTx
}
