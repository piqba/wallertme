package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of wallertmectl",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:\t%s\n", Version)
		fmt.Printf("Version Git Hash:\t%s\n", shortGitCommit(VersionHash))
		fmt.Printf("Build time:\t%s\n", BuildTime)
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
