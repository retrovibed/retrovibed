//go:build !linux

package fsx

import (
	"io/fs"
	"time"
)

// ctime is unavailable on this platform.
func ctime(info fs.FileInfo) (time.Time, bool) {
	return time.Time{}, false
}
