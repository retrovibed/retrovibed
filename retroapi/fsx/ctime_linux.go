package fsx

import (
	"io/fs"
	"syscall"
	"time"
)

// ctime returns the inode status change time, the closest thing linux exposes
// through the standard library to a creation timestamp. false when unavailable.
func ctime(info fs.FileInfo) (time.Time, bool) {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return time.Time{}, false
	}

	return time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec), true
}
