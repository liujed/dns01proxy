package dns01proxy

import (
	"fmt"
	"runtime"

	"github.com/liujed/dns01proxy/gomodversions"
)

// Version identifier to use when we're unable to get version information from
// Go's build info.
var FallbackVersion = "unknown"

// Returns the release string, including the application name, version, go-mod
// hash, OS, and architecture. For example, "dns01proxy v0.0.0 (h1:abcd1234=)
// linux/amd64".
func Release() string {
	_, version := gomodversions.GetVersionOfGoPackage(
		"github.com/liujed/dns01proxy",
	)

	// XXX Work around https://github.com/golang/go/issues/29228.
	if version == "unknown" {
		version = FallbackVersion
	}

	return fmt.Sprintf(
		"dns01proxy %s %s/%s",
		version,
		runtime.GOOS,
		runtime.GOARCH,
	)
}
