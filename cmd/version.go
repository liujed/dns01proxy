package cmd

import (
	"fmt"

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
		fmt.Println("Available DNS providers:")
		for idx, caddyModuleInfo := range caddy.GetModules("dns.providers") {
			if idx > 0 {
				fmt.Println()
			}

			goModulePath, version := gomodversions.GetVersionOfValue(caddyModuleInfo.New())
			fmt.Printf("  %s\t%s\n    %s\n",
				caddyModuleInfo.ID.Name(),
				version,
				goModulePath,
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
