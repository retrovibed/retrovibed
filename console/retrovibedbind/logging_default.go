//go:build !android

package main

import (
	"fmt"
	"log"
	"os"
)

// RedirectLogs is a no-op on non-android platforms.
func redirectlogs() {
	// Standard behavior remains: Go logs to the default stdout/stderr
	// which works normally on desktop and iOS.
	log.SetFlags(log.Lshortfile | log.LUTC | log.Ltime)
	log.SetPrefix(fmt.Sprintf("%d ", os.Getpid()))
}
