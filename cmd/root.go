package cmd

import (
	"os"

	"github.com/caddyserver/caddy/v2"
	"github.com/liujed/caddy-dns01proxy/flags"
	"github.com/liujed/dns01proxy"
	"github.com/liujed/goutil/optionals"
	"github.com/spf13/cobra"
)

// Flag definitions.
var (
	flgDebug = flags.Flag[bool]{
		Persistent: true,
		Name:       "debug",
		ShortName:  optionals.Some('v'),
		UsageMsg:   "turn on verbose debug logs",
	}
)

// The base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "dns01proxy",
	Short: "Proxy server for ACME DNS-01 challenges",
	Long: `DNS-01 proxy server

dns01proxy is a server for using DNS-01 challenges to obtain TLS certificates
from Let's Encrypt, or any ACME-compatible certificate authority, without
exposing your DNS credentials to every host that needs a certificate.

It acts as a proxy for DNS-01 challenge requests, allowing hosts to delegate
their DNS record updates during ACME validation. This makes it possible to issue
certificates to internal or private hosts that can't (or shouldn't) have direct
access to your DNS provider or API keys.

Designed to work with:
  * acme.sh's 'acmeproxy' provider,
  * Caddy's 'acmeproxy' DNS provider module, and
  * lego's 'httpreq' DNS provider.


DNS PROVIDER MODULES

dns01proxy is built using Caddy, and uses DNS provider modules that are written
by the Caddy community. To limit the amount of third-party code that you run,
dns01proxy ships binaries that are built with a single DNS module.

Make sure that the binary you are running supports your DNS provider, and always
check that you trust the author of the DNS module. The release notes has details
about the source of the DNS module in each build:

  https://github.com/liujed/dns01proxy/releases

` + getDNSProviderVersions(),
	Version: dns01proxy.Release(),
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(caddy.ExitCodeFailedStartup)
	}
}

func init() {
	rootCmd.SetVersionTemplate("{{.Version}}\n")

	flags.AddBoolFlag(rootCmd, flgDebug)
}
