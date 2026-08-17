package debian

import (
	"context"
	"eg/compute/errorsx"
	"eg/compute/maintainer"
	"embed"
	"io/fs"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egdebuild"
	"github.com/egdaemon/eg/runtime/x/wasi/eggolang"
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
		egdebuild.Option.Description("media distribution platform", "provides torrenting functionality with builtin media player and cloud storage allowing you to take your content anywhere you go"),
		egdebuild.Option.Debian(errorsx.Must(fs.Sub(debskel, ".debskel"))),
		egdebuild.Option.DependsBuild("golang-1.26", "cargo", "rustc", "tree", "dh-make", "debhelper", "pkg-config", "duckdb", "libavcodec-dev", "libavformat-dev", "libavutil-dev", "libswresample-dev", "libavfilter-dev", "libavdevice-dev", "libswscale-dev"),
		egdebuild.Option.Depends("duckdb", "ffmpeg"),
	)
}

func Prepare(ctx context.Context, o eg.Op) error {
	debdir := cachedir()
	sruntime := shell.Runtime()
	return eg.Sequential(
		shell.Op(
			sruntime.Newf("echo '-----------------------------------------'"),
			sruntime.Newf("eg gpg keyring --name=\"${EG_GPG_KEYRING_NAME}\" --email=\"${EG_GPG_KEYRING_EMAIL}\" --seed=\"${EG_GPG_KEYRING_SEED}\""),
			sruntime.Newf("echo '-----------------------------------------'"),
			sruntime.Newf("rm -rf %s", debdir),
			sruntime.Newf("mkdir -p %s", debdir),
			sruntime.Newf("git clone --depth 1 file://${PWD}/ %s", debdir),
		),
		egdebuild.Prepare(Runner(), errorsx.Must(fs.Sub(debskel, ".debskel"))),
	)(ctx, o)
}

// container for this package.
func Runner() eg.ContainerRunner {
	return eg.Container("retrovibe.debuild.ubuntu")
}

func Build(ctx context.Context, o eg.Op) error {
	return eg.Sequential(
		egdebuild.Build(gcfg, egdebuild.Option.Distro(egdebuild.UbuntuLatestCodename), egdebuild.Option.NoLint()), // resolute
		egdebuild.Build(
			gcfg,
			egdebuild.Option.Distro(egdebuild.UbuntuLatestCodename),
			egdebuild.Option.BuildBinary(20*time.Minute),
			egdebuild.Option.Environ(eggolang.Env()...),
			egdebuild.Option.NoLint(),
		),
	)(ctx, o)
}

func Upload(ctx context.Context, o eg.Op) error {
	return egdebuild.UploadDPut(gcfg, errorsx.Must(fs.Sub(debskel, ".debskel")))(ctx, o)
}

// Release builds and uploads the debian package.
func Release(ctx context.Context, o eg.Op) error {
	deb := eg.Container(maintainer.Container)
	return eg.Perform(
		ctx,
		eg.Parallel(
			eg.Build(deb.BuildFromFile(".eg/Containerfile")),
			Prepare,
		),
		eg.Module(
			ctx,
			deb,
			Build,
			Upload,
		),
	)
}
