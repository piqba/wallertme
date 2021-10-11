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
	"method":"eth_getBalance",
	"params":[
		"0xEF7c91507AFc7BBe79ede720F980bf39Fe488F94",
		"latest"
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
	payload := web3.PayloadReq{
		Jsonrpc: "2.0",
		Method:  "eth_getBalance",
		Params: []interface{}{
			"0xEF7c91507AFc7BBe79ede720F980bf39Fe488F94",
			"latest",
		},
		ID: 1,
	}

	blc, err := api.Balance(context.TODO(), payload)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", blc.ToJSON())
}
