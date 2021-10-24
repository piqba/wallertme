package main

import (
	"context"
	"fmt"

	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/logger"
	"github.com/piqba/wallertme/pkg/web3"
)

func main() {
	solanaApi, err := web3.NewAPISolanaClient(web3.APIClientOptions{})
	if err != nil {
		logger.LogError(errors.Errorf("main:%s", err).Error())
	}

	payloadLastTx := web3.PayloadReqJSONRPC{
		Jsonrpc: "2.0",
		Method:  "getSignaturesForAddress",
		Params: []interface{}{
			"9hZaTvCVMcfbheTzebkeGR6Xi2EzMqTtPasbhGoPB94C",
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
	fmt.Println(lastTx.Result[0].ToJSON())

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
	fmt.Println(infoTx.ToJSON())
}
