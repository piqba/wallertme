package main

import (
	"context"
	"fmt"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/logger"
	"github.com/piqba/wallertme/pkg/storage"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.LogError(errors.Errorf("main:%s", err).Error())
	}
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu) // Try to use all available CPUs.
}

func main() {
	pgx, err := storage.PostgreSQLConnection(context.Background())
	if err != nil {
		logger.LogError(errors.Errorf("bb8: %v", err).Error())

	}
	dataSource := storage.NewSource("db", storage.OptionsSource{
		Pgx: pgx,
	})
	wallets, err := dataSource.Wallets(context.Background())
	if err != nil {
		logger.LogError(errors.Errorf("bb8: %v", err).Error())
	}
	for _, w := range wallets {
		fmt.Println(w.ToJSON())
	}
}
