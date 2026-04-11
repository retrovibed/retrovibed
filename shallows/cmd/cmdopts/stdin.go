package cmdopts

import "os"

func Readable(v *os.File) bool {
	stat, _ := v.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
