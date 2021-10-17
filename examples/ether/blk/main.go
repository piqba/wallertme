package main

import (
	"context"
	"fmt"
	"log"

	"github.com/piqba/wallertme/pkg/web3"
)

/*
curl --location --request POST 'localhost:8545/' \
--header 'Content-Type: application/json' \
--data-raw '{
	"jsonrpc":"2.0",
	"method":"eth_getBlockByNumber",
	"params":[
		"latest",
		true
	],
	"id":1
}'
*/

func main() {

	api, err := web3.NewAPIEthClient(
		web3.APIClientOptions{
			Server: web3.GanacheDevNet,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	payload := web3.PayloadReqEth{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params: []interface{}{
			"latest",
			true,
		},
		ID: 1,
	}

	blc, err := api.BlockByNumberETH(context.TODO(), payload)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", blc.ToJSON())
}
