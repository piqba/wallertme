package main

import (
	"context"
	"fmt"
	"log"

	"github.com/piqba/wallertme/pkg/web3"
)

/*
curl https://explorer-api.testnet.dandelion.link/api/addresses/summary/addr_test1qq6g6s99g9z9w0mlvew28w40lpml9rwfkfgerpkg6g2vpn6dp4cf7k9drrdy0wslarr6hxspcw8ev5ed8lfrmaengneqz34lcx
*/

func main() {
	cardano, err := web3.NewAPICardanoClient(web3.APIClientOptions{})
	if err != nil {
		log.Fatal(err)
	}

	address := "addr_test1qq6g6s99g9z9w0mlvew28w40lpml9rwfkfgerpkg6g2vpn6dp4cf7k9drrdy0wslarr6hxspcw8ev5ed8lfrmaengneqz34lcx"
	sumary, err := cardano.SumaryAddrADA(context.TODO(), address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(sumary.Result.ToJSON())
}
