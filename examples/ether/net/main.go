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
	"method":"net_version",
	"params":[],
	"id":67
}'
*/

func main() {

	api, err := web3.NewAPIEthClient(
		web3.APIClientOptions{
			Server: web3.PolygonTestNet,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	version, err := api.Version(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", version.ToJSON())
}
