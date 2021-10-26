package commands

import (
	"os"

	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	flagProvider = "provider"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "walletctl",
	Short: "A cli to send data from redis to DISCORD|TELEGRAM|SMTP",
	Long:  `this is wallertme ctl`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {

		logger.LogError(errors.Errorf("walletctl: %v", err).Error())
		os.Exit(1)
	}

}
func init() {

}
