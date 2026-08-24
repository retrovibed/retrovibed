// Package eggradlex sets up the gradle environment for caching, mirroring
// eggolang's GOCACHE/GOMODCACHE and console.go's PUB_CACHE pattern.
package eggradlex

import (
	"fmt"

	"github.com/egdaemon/eg/runtime/wasi/egenv"
)

// CacheDirectory is the persistent GRADLE_USER_HOME shared across builds, so
// downloaded dependencies survive a clean clone of the working directory.
func CacheDirectory() string {
	return egenv.CacheDirectory(".eg", "gradle")
}

// Env returns the environment variables that point gradle at CacheDirectory.
func Env() []string {
	return []string{
		fmt.Sprintf("GRADLE_USER_HOME=%s", CacheDirectory()),
	}
}
