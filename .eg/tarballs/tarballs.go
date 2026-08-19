package tarballs

import (
	"context"
	"strings"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
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
	return egtarball.Tarxz(b.StringReplace("retrovibed.%goenv.os%.%goenv.arch%"))
}

func RetrovibedSource() string {
	return egtarball.Targz("retrovibed.source")
}

func Flatpak(b *Build) string {
	return b.StringReplace("space.retrovibe.Console.yml")
}

func AppImage(b *Build) string {
	return b.StringReplace("retrovibed.%goenv.os%.%goenv.arch%.AppImage")
}

func AppImageZsync(b *Build) string {
	return AppImage(b) + ".zsync"
}

func AppImageBuild(b *Build) string {
	return b.StringReplace("retrovibed.%goenv.os%.%goenv.arch%.AppImageBuild")
}

// Version returns the date-based build version string from git commit metadata.
func Version() string {
	return eggit.EnvCommit().StringReplace("%git.commit.year%.%git.commit.month%.%git.commit.day%")
}

func VersionUnix() string {
	return eggit.EnvCommit().StringReplace("%git.commit.year%.%git.commit.month%.%git.commit.unix%")
}

func Tarchive(ctx context.Context, op eg.Op) error {
	return eg.Sequential(
		shell.Op(
			shell.Newf("git archive --format=tar.gz -o %s %s", egtarball.Archive(RetrovibedSource()), eggit.EnvCommit().StringReplace("%git.hash%")),
		),
	)(ctx, op)
}
