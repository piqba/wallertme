package cli

import (
	"fmt"

	"github.com/piqba/wallertme/pkg/errors"
	"github.com/piqba/wallertme/pkg/logger"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "r2d2",
	Short: "Send data from redis to services providers telegram|discord|smtp",
	Run: func(cmd *cobra.Command, args []string) {
		provider, err := cmd.Flags().GetString(flagProvider)
		if err != nil {
			logger.LogError(errors.Errorf("r2d2ctl: %v", err).Error())
		}

		fmt.Println(provider)

	},
}

func init() {
	fetchCmd.Flags().String(flagProvider, "telegram", "select a provider to send notifications")
	rootCmd.AddCommand(fetchCmd)

}
