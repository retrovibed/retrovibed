package debian

import (
	"context"
	"eg/compute/errorsx"
	"eg/compute/maintainer"
	"embed"
	"io/fs"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egdebuild"
)

//go:embed .debskel
var debskel embed.FS

func cachedir() string {
	return egenv.CacheDirectory(".dist", "retrovibed")
}

var (
	gcfg egdebuild.Config
)

func init() {
	c := eggit.EnvCommit()
	gcfg = egdebuild.New(
		"retrovibed",
		"",
		cachedir(),
		egdebuild.Option.Maintainer(maintainer.Name, maintainer.Email),
		egdebuild.Option.SigningKeyID(maintainer.GPGFingerprint),
		egdebuild.Option.ChangeLogDate(c.Committer.When),
		egdebuild.Option.Version("0.0.:autopatch:"),
		egdebuild.Option.Debian(errorsx.Must(fs.Sub(debskel, ".debskel"))),
		egdebuild.Option.DependsBuild("golang-1.24", "dh-make", "debhelper", "duckdb"),
		egdebuild.Option.Depends("duckdb"),
		egdebuild.Option.Environ("VCS_REVISION", c.Hash.String()),
	)
}

func Prepare(ctx context.Context, o eg.Op) error {
	debdir := cachedir()
	sruntime := shell.Runtime()
	return eg.Sequential(
		shell.Op(
			sruntime.Newf("rm -rf %s", debdir),
			sruntime.Newf("mkdir -p %s", debdir),
			sruntime.Newf("git clone --depth 1 file://${PWD}/ %s", debdir),
		),
		egdebuild.Prepare(Runner(), errorsx.Must(fs.Sub(debskel, ".debskel"))),
	)(ctx, o)
}

// container for this package.
func Runner() eg.ContainerRunner {
	return eg.Container("retrovibe.debuild.ubuntu.24.10")
}

func Build(ctx context.Context, o eg.Op) error {
	return eg.Parallel(
		egdebuild.Build(gcfg, egdebuild.Option.Distro("jammy")),
		egdebuild.Build(gcfg, egdebuild.Option.Distro("noble")),
		egdebuild.Build(gcfg, egdebuild.Option.Distro("oracular")),
	)(ctx, o)
}

func Upload(ctx context.Context, o eg.Op) error {
	return egdebuild.UploadDPut(gcfg, errorsx.Must(fs.Sub(debskel, ".debskel")))(ctx, o)
}
