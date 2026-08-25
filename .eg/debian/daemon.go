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
	return egenv.WorkspaceDirectory(".git", "retrovibed")
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
		egdebuild.Option.Environ("GNUPGHOME=/home/egd/.gnupg"),
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
func Runner() eg.ContainerRunner {
	return eg.Container("retrovibe.debuild.ubuntu")
}

func Build(ctx context.Context, o eg.Op) error {
	return eg.Sequential(
		// override eg compute local's setting of the home directory.
		// ideally we should default to this except when performing local
		// compute workloads.
		eggpg.Seed(eggpg.Options().Home("/home/egd/.gnupg")...),
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

func Copy(cachedir string) eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		debroot := egenv.EphemeralDirectory("deb.retrovibed")
		return shell.Op(
			shell.Newf("tree -L 2 %s", debroot),
			shell.Newf("tree -L 1 %s", cachedir),
			shell.Newf(`rsync -av --include 'retrovibed_*.deb' --exclude '*' --mkpath %s/ %s/`, debroot, cachedir).Debug(),
		)(ctx, o)
	}
}
