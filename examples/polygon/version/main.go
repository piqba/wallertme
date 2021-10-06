package main

import (
	"context"
	"fmt"
	"log"

	"github.com/piqba/wallertme/pkg/web3"
)

func main() {

	api, err := web3.NewAPIClient(
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
