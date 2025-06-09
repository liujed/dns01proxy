package cmd

import (
	"fmt"
	"os"

	"github.com/caddyserver/caddy/v2"
	caddycmd "github.com/caddyserver/caddy/v2/cmd"
	"github.com/liujed/caddy-dns01proxy/flags"
	_ "github.com/liujed/dns01proxy/internal/caddyimports"
	"github.com/liujed/goutil/optionals"
	"github.com/spf13/cobra"
)

// Flag definitions.
var (
	flgConfig = flags.Flag[string]{
		Name:         "config",
		ShortName:    optionals.Some('c'),
		UsageMsg:     "read configuration from `FILE`",
		Required:     true,
		FilenameExts: optionals.Some([]string{"json", "toml"}),
	}
)

// The `run` command.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Starts the dns01proxy server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Log an error if there aren't any DNS providers available.
		if len(caddy.GetModules("dns.providers")) == 0 {
			caddy.Log().Error(
				`No DNS provider modules found. dns01proxy will fail to start. Please make sure that your copy of dns01proxy is compiled with a DNS provider. Download a release at https://github.com/liujed/dns01proxy/releases`,
			)
		}

		flags := caddycmd.Flags{FlagSet: cmd.Flags()}
		configFile, err := flags.GetString(flgConfig.Name)
		if err != nil {
			return fmt.Errorf("unable to read %s flag: %w", flgConfig.Name, err)
		}

		debug, err := flags.GetBool(flgDebug.Name)
		if err != nil {
			return fmt.Errorf("unable to read %s flag: %w", flgDebug.Name, err)
		}

		os.Args = []string{os.Args[0], "dns01proxy", "--config", configFile}
		if debug {
			os.Args = append(os.Args, "--debug")
		}

		caddycmd.Main()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	flags.AddStringFlag(runCmd, flgConfig)
}
