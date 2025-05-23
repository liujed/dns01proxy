package cmd

import "github.com/spf13/cobra"

// The `hash-password` command.
var hashPasswordCmd = &cobra.Command{
	Use:   "hash-password",
	Short: "Hashes a password and prints the result to stdout in base64",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(hashPasswordCmd)
}
