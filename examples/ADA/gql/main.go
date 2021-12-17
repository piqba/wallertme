package main

import (
	"context"
	"fmt"

	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/logger"
	"github.com/piqba/wallertme/pkg/web3"
)

func main() {

	cardanoApi, err := web3.NewAPICardanoClient(web3.APIClientOptions{})
	if err != nil {
		logger.LogError(errors.Errorf("main:%s", err).Error())
	}
	pld := web3.PayloadReqJSONGQL{
		Query: `
		query utxoSetForAddress (
			$address: String!
		){
			utxos(
				order_by: { value: desc }
				where: { address: { _eq: $address }}
				limit :1
			) {
				# address,
			  value,
			  txHash,
			  transaction{
				block{number, hash},
				fee,
				totalOutput,
				includedAt,
				inputs{address,value},
				outputs{address, value}
			  },
			}
		}
		`,
		Variables: map[string]string{
			"address": "addr_test1qq6g6s99g9z9w0mlvew28w40lpml9rwfkfgerpkg6g2vpn6dp4cf7k9drrdy0wslarr6hxspcw8ev5ed8lfrmaengneqz34lcx",
		},
	}

	data, err := cardanoApi.LastTxByAddressADA(context.Background(), pld)
	if err != nil {
		logger.LogError(err.Error())
	}
	fmt.Println(data.ToJSON())
}
