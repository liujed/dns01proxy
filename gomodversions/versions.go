package gomodversions

import (
	"reflect"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/liujed/goutil/maps"
)

// Maps Go module paths to their version strings.
var gomodVersions = maps.NewHashMap[string, string]()

// Returns the path and version information of the Go module that defines the
// given package. Version strings include package version and go-mod hash (e.g.,
// "v0.0.0 (h1:abcd1234=)").
func GetVersionOfGoPackage(pkgPath string) (
	goModulePath string,
	version string,
) {
	sync.OnceFunc(func() {
		// Load version information for all modules.
		buildInfo, ok := debug.ReadBuildInfo()
		if !ok {
			return
		}
		for _, module := range buildInfo.Deps {
			buf := strings.Builder{}
			buf.WriteString(module.Version)
			if module.Sum != "" {
				buf.WriteString(" (")
				buf.WriteString(module.Sum)
				buf.WriteRune(')')
			}
			if module.Replace != nil {
				buf.WriteString(" => ")
				buf.WriteString(module.Replace.Path)
				if module.Replace.Version != "" {
					buf.WriteRune('@')
					buf.WriteString(module.Replace.Version)
				}
				if module.Replace.Sum != "" {
					buf.WriteString(" (")
					buf.WriteString(module.Replace.Sum)
					buf.WriteRune(')')
				}
			}

			gomodVersions.Put(module.Path, buf.String())
		}
	})()

	modPath := pkgPath
	for {
		if version, exists := gomodVersions.Get(modPath).Get(); exists {
			return modPath, version
		}

		idx := strings.LastIndex(modPath, "/")
		if idx == -1 {
			return pkgPath, "unknown"
		}
		modPath = modPath[0:idx]
	}
}

// Returns the path and version of the Go module that defines the type of the
// given value.
func GetVersionOfValue(
	v any,
) (
	goModulePath string,
	version string,
) {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	pkgPath := rv.Type().PkgPath()
	return GetVersionOfGoPackage(pkgPath)
}

// Returns the path and version of the Go module that defines the type T.
func GetVersionOfType[T any]() (
	goModulePath string,
	version string,
) {
	t := reflect.TypeFor[T]()
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	pkgPath := t.PkgPath()

	return GetVersionOfGoPackage(pkgPath)
}
