package dns01proxy

import (
	"fmt"
	"runtime"

	"github.com/liujed/dns01proxy/gomodversions"
)

// Returns the release string, including the application name, version, go-mod
// hash, OS, and architecture. For example, "dns01proxy v0.0.0 (h1:abcd1234=)
// linux/amd64".
func Release() string {
	_, version := gomodversions.GetVersionOfGoPackage(
		"github.com/liujed/dns01proxy",
	)
	return fmt.Sprintf(
		"dns01proxy %s %s/%s",
		version,
		runtime.GOOS,
		runtime.GOARCH,
	)
}
