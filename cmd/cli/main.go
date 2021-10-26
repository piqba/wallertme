package main

import (
	"runtime"

	"github.com/joho/godotenv"
	"github.com/piqba/wallertme/cmd/cli/commands"
	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/logger"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.LogError(errors.Errorf("walletctl: %v", err).Error())
	}
	numcpu := runtime.NumCPU()
	runtime.GOMAXPROCS(numcpu) // Try to use all available CPUs.
}

func main() {
	commands.Execute()
}
