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
			Server: web3.GanacheDevNet,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	payload := web3.PayloadReq{
		Jsonrpc: "2.0",
		Method:  "eth_getTransactionByHash",
		Params: []interface{}{
			"0x4e0cc9af12392335b447a823e50d9260a14db1d4445117e48e5247d26623e15c",
		},
		ID: 1,
	}
	tx, err := polygon.TransactionByHashETH(
		context.TODO(),
		payload,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", tx.Result.ToJSON())

	txr, err := polygon.TransactionReceiptETH(
		context.TODO(),
		payload,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", txr.Result.ToJSON())
}
