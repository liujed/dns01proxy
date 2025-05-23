package cmd

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"
	"github.com/liujed/caddy-dns01proxy/flags"
	"github.com/liujed/goutil/optionals"
	"github.com/spf13/cobra"
)

// Flag definitions.
var (
	flgPlaintext = flags.Flag[string]{
		Name:      "plaintext",
		ShortName: optionals.Some('p'),
		UsageMsg:  "the plaintext password",
	}

	flgAlgorithm = flags.Flag[string]{
		Name:         "algorithm",
		ShortName:    optionals.Some('a'),
		DefaultValue: "bcrypt",
		UsageMsg:     "the name of the hash algorithm",
		Hidden:       true,
	}
)

// The `hash-password` command.
var hashPasswordCmd = &cobra.Command{
	Use:   "hash-password",
	Short: "Hashes a password and prints the result to stdout in base64",
	Run: func(cmd *cobra.Command, args []string) {
		caddycmd.Main()
	},
}

func init() {
	rootCmd.AddCommand(hashPasswordCmd)

	flags.AddStringFlag(hashPasswordCmd, flgPlaintext)
	flags.AddStringFlag(hashPasswordCmd, flgAlgorithm)
}
