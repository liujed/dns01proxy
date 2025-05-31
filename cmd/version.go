package cmd

import (
	"fmt"
	"strings"

	"github.com/caddyserver/caddy/v2"
	caddydns01proxy "github.com/liujed/caddy-dns01proxy"
	"github.com/liujed/dns01proxy"
	"github.com/liujed/dns01proxy/gomodversions"
	_ "github.com/liujed/dns01proxy/internal/caddyimports"
	"github.com/spf13/cobra"
)

// The `version` command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(dns01proxy.Release())
		fmt.Println()
		fmt.Println("Built with")
		_, caddyFullVersion := caddy.Version()
		fmt.Printf("  caddy %s\n", caddyFullVersion)

		_, version := gomodversions.GetVersionOfType[caddydns01proxy.App]()
		fmt.Printf("  caddy-dns01proxy %s\n", version)

		fmt.Println()
		fmt.Println()

		fmt.Print(getDNSProviderVersions())
	},
}

// Returns details about the DNS providers available, suitable for displaying as
// part of help or version information.
func getDNSProviderVersions() string {
	modules := caddy.GetModules("dns.providers")
	if len(modules) == 0 {
		return `NO DNS PROVIDERS ARE AVAILABLE IN THIS BUILD. Please make sure that your copy of
dns01proxy is compiled with a DNS provider. Download a release at
https://github.com/liujed/dns01proxy/releases
`
	}

	buf := strings.Builder{}
	buf.WriteString("DNS providers in this build:\n")
	for _, caddyModuleInfo := range modules {
		buf.WriteRune('\n')

		goModulePath, version := gomodversions.GetVersionOfValue(caddyModuleInfo.New())
		buf.WriteString(
			fmt.Sprintf(
				"  %s\t%s\n    %s\n",
				caddyModuleInfo.ID.Name(),
				version,
				goModulePath,
			),
		)
	}

	return buf.String()
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
