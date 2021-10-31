package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	buildTime    string
	version      string
	versionHash  string
	otelNameBb8  = "bb8"
	otelNameR2D2 = "r2d2"
	otelVersion  = "v0.4.0"
	otelNameEnv  = "dev"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of wallertmectl",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Version Git Hash:\t%s\n", shortGitCommit(versionHash))
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	},
}

func shortGitCommit(fullGitCommit string) string {
	shortCommit := ""
	if len(fullGitCommit) >= 6 {
		shortCommit = fullGitCommit[0:6]
	}

	return shortCommit
}
func init() {
	rootCmd.AddCommand(versionCmd)
}
