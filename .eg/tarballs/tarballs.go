package tarballs

import (
	"context"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

func Retrovibed() string {
	return egtarball.Tarxz(egenv.Pattern("retrovibed.{eg.host.os}.{eg.host.arch}"))
}

func RetrovibedSource() string {
	return egtarball.Targz("retrovibed.source")
}

func Flatpak() string {
	return "space.retrovibe.Console.yml"
}

func AppImage() string {
	return egenv.Pattern("retrovibed.{eg.host.os}.{eg.host.arch}.AppImage")
}

func AppImageZsync() string {
	return AppImage() + ".zsync"
}

func AppImageBuild() string {
	return egenv.Pattern("retrovibed.{eg.host.os}.{eg.host.arch}.AppImageBuild")
}

// Version returns the date-based build version string from git commit metadata.
func Version() string {
	return eggit.EnvCommit().StringReplace("{git.commit.year}.{git.commit.month}.{git.commit.day}")
}

func VersionUnix() string {
	return eggit.EnvCommit().StringReplace("{git.commit.year}.{git.commit.month}.{git.commit.unix}")
}

func Tarchive(ctx context.Context, op eg.Op) error {
	return eg.Sequential(
		shell.Op(
			shell.Newf("git archive --format=tar.gz -o %s %s", egtarball.Archive(RetrovibedSource()), eggit.EnvCommit().StringReplace("{git.hash}")),
		),
	)(ctx, op)
}
