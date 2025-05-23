package cmd

import (
	"github.com/spf13/cobra"
)

// The `run` command.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts the dns01proxy server in the foreground",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
