package main

import (
	"context"
	"fmt"
	"log"

	domain "github.com/piqba/wallertme/internal/bb8/domain"
)

const (
	addrSender   = "addr_test1qq6g6s99g9z9w0mlvew28w40lpml9rwfkfgerpkg6g2vpn6dp4cf7k9drrdy0wslarr6hxspcw8ev5ed8lfrmaengneqz34lcx"
	addrReceiver = "addr_test1qq5287luxzj5l4lequrqdp5ln76ver4uls3z0m5ykr5gqsv0vxzrwcq5dmmn9e09rvgttzgrngmpxkguy7220r0u0ljqzuww7g"
)

func main() {
	repo := domain.NewTx(10)
	info, err := repo.InfoByAddress(context.TODO(), addrReceiver)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(info.ToJSON())
}
