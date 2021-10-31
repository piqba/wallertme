package commands

import (
	"os"

	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/logger"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wallertmectl",
	Short: "A cli to send tx data from (SOLANA|CARDANO) blockchain to DISCORD|TELEGRAM|SMTP",
	Long: `Wallertme ctl is a tool focused on: 
	Send tx data from (SOLANA|CARDANO) blockchain to a queue like (REDIS) streams and then send this information
	to DISCORD|TELEGRAM|SMTP
	`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {

		logger.LogError(errors.Errorf("walletctl: %v", err).Error())
		os.Exit(1)
	}

}
func init() {

}
