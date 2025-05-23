package cmd

import (
	"fmt"

	"github.com/liujed/dns01proxy"
	"github.com/spf13/cobra"
)

// The `version` command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(dns01proxy.Release())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
