// Package retrozsync builds a debian package for github.com/cph6/zsync, a Go
// rewrite of zsync that uses net/http and so supports HTTPS out of the box
// (the classic C zsync package in Ubuntu/Debian has no TLS support at all).
// installs to /usr/local/bin since it isn't the distro's zsync package and
// doesn't replace/conflict with it.
package retrozsync

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"time"

	"eg/compute/errorsx"
	"eg/compute/maintainer"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egdebuild"
)

//go:embed .debskel
var debskel embed.FS

const (
	version = "0.8.0"
)

var (
	gcfg egdebuild.Config
)

func init() {
	c := eggit.EnvCommit()
	gcfg = egdebuild.New(
		"retrozsync",
		"",
		egenv.CacheDirectory("retrozsync"),
		egdebuild.Option.Maintainer(maintainer.Name, maintainer.Email),
		egdebuild.Option.SigningKeyID(maintainer.GPGFingerprint),
		egdebuild.Option.ChangeLogDate(c.Committer.When),
		egdebuild.Option.Version(fmt.Sprintf("%s.:autopatch:", version)),
		egdebuild.Option.Description("Go rewrite of zsync with HTTPS support", "installs zsync and zsyncmake built from github.com/cph6/zsync to /usr/local/bin"),
		egdebuild.Option.Debian(errorsx.Must(fs.Sub(debskel, ".debskel"))),
		egdebuild.Option.DependsBuild("golang-1.26", "tree"),
	)
}

func Prepare(ctx context.Context, o eg.Op) error {
	return eg.Parallel(
		Download,
		egdebuild.Prepare(Runner(), nil),
	)(ctx, o)
}

func Runner() eg.ContainerRunner {
	return egdebuild.Runner()
}

func Build(ctx context.Context, o eg.Op) error {
	debroot := egenv.EphemeralDirectory("deb.retrozsync")
	cachedir := egenv.CacheDirectory(".dist")
	runtime := shell.Runtime().Timeout(egenv.TTL())

	return eg.Sequential(
		eg.Parallel(
			egdebuild.Build(gcfg, egdebuild.Option.Distro(egdebuild.UbuntuLatestCodename)),
			egdebuild.Build(
				gcfg,
				egdebuild.Option.Distro(egdebuild.UbuntuLatestCodename),
				egdebuild.Option.BuildBinary(20*time.Minute),
				egdebuild.Option.NoLint(),
			),
		),
		shell.Op(
			runtime.Newf("mkdir -p %s", cachedir),
			runtime.Newf(`find %s -maxdepth 2 -name 'retrozsync_*.deb' -exec cp -v {} %s \;`, debroot, cachedir),
		),
	)(ctx, o)
}

func Upload(ctx context.Context, o eg.Op) error {
	return egdebuild.UploadDPut(gcfg, errorsx.Must(fs.Sub(debskel, ".debskel")), egdebuild.Option.Timeout(20*time.Minute))(ctx, o)
}
