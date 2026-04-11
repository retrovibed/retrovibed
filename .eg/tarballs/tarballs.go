package tarballs

import (
	"strings"
)

type Build struct {
	Arch string
	OS   string
}

// Replace patterns in a string with their values from information from go env.
// %goenv.arch%            -> machine architecture. amd, amd64, arm64, ..., etc.
// %goenv.os%              -> operating system. linux, darwin, windows, ..., etc.
func (c Build) StringReplace(pattern string) string {
	s := strings.ReplaceAll(pattern, "%goenv.arch%", c.Arch)
	s = strings.ReplaceAll(s, "%goenv.os%", c.OS)

	return s
}

func Retrovibed(b *Build) string {
	return b.StringReplace("retrovibed.%goenv.os%.%goenv.arch%")
}

func Flatpak(b *Build) string {
	return b.StringReplace("space.retrovibe.Console.yml")
}
