package main

type TemplateData struct {
	// Identifies the Caddy package in this build. For example, "cloudflare" or
	// "rfc2136" for github.com/caddy-dns/cloudflare and
	// github.com/caddy-dns/rfc2136, respectively.
	PackageName string

	// The Go module path for the plugin's Caddy package. For example,
	// "github.com/caddy-dns/cloudflare".
	ModPath string

	// The version number and go-mod hash for the plugin's Caddy package. For
	// example, "v0.0.0 (h1:abcd1234=)".
	Version string
}
