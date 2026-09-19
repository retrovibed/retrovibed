package cmdopts

import (
	"io"
	"os"
)

// Stdin is the kong binding commands accept to read their input, allowing tests
// to inject a reader instead of touching the process-global os.Stdin.
type Stdin interface {
	io.Reader
}

func Readable(v *os.File) bool {
	stat, _ := v.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
