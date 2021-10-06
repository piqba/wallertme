package main

import (
	"context"
	"fmt"
	"log"

	"github.com/piqba/wallertme/pkg/web3"
)

/*
curl --location --request POST 'https://matic-mumbai.chainstacklabs.com/' \
--header 'Content-Type: application/json' \
--data-raw '{
	"jsonrpc":"2.0",
	"method":"eth_getTransactionByHash",
	"params":[
		"0x0314f52b94f624695e9035df6f76ba7c0209a57462ec6c9ade577523883fb681"
	],
	"id":1
}'
*/

func main() {

	polygon, err := web3.NewAPIEthClient(
		web3.APIClientOptions{
			Server: web3.PolygonTestNet,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := polygon.TransactionByHash(
		context.TODO(),
		"0x0314f52b94f624695e9035df6f76ba7c0209a57462ec6c9ade577523883fb681",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", tx.Result.ToJSON())
}
