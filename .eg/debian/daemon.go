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
	"github.com/egdaemon/eg/runtime/x/wasi/eggpg"
)

//go:embed .debskel
var debskel embed.FS

func cachedir() string {
	return egenv.DirectoryWorkload(".git", "retrovibed")
}

func gpgoptions() []eggpg.Option {
	// eggpg.Options() takes GNUPGHOME from the environment, which inside the container is the
	// host-provided /eg.mnt/.gnupg (not writable by us: "unsafe ownership", "Permission denied").
	// pin the isolated per-container home instead.
	return eggpg.Options().Home("/home/egd/.gnupg")
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
		egdebuild.Option.DependsBuild("golang-1.27", "cargo", "rustc", "tree", "dh-make", "debhelper", "pkg-config", "duckdb", "libavcodec-dev", "libavformat-dev", "libavutil-dev", "libswresample-dev", "libavfilter-dev", "libavdevice-dev", "libswscale-dev"),
		egdebuild.Option.Depends("duckdb", "ffmpeg"),
		egdebuild.Option.Environ(eggpg.Env(gpgoptions()...)...),
	)
}

func Prepare(ctx context.Context, o eg.Op) error {
	debdir := cachedir()
	sruntime := shell.Runtime()
	return eg.Sequential(
		shell.Op(
			sruntime.Newf("git worktree add %s HEAD", debdir),
		),
		egdebuild.Prepare(Runner(), errorsx.Must(fs.Sub(debskel, ".debskel"))),
	)(ctx, o)
}

// container for this package.
func Runner() eg.C8S {
	return eg.Container("retrovibe.debuild.ubuntu")
}

func Build(ctx context.Context, o eg.Op) error {
	return eg.Sequential(
		eggpg.Seed(gpgoptions()...),
		// eggpg.Debug(gpgoptions()...),
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

func Copy(dst string) eg.OpFn {
	return egdebuild.Copy(gcfg, dst)
}
