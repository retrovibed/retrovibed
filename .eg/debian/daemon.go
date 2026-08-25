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
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
	"github.com/egdaemon/eg/runtime/x/wasi/egdebuild"
	"github.com/egdaemon/eg/runtime/x/wasi/eggolang"
	"github.com/egdaemon/eg/runtime/x/wasi/eggpg"
)

//go:embed .debskel
var debskel embed.FS

func cachedir() string {
	return egenv.WorkspaceDirectory(".git", "retrovibed")
}

func gpgoptions() []eggpg.Option {
	return eggpg.Options().Privileged()
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
		egdebuild.Option.Environ(eggpg.Env(gpgoptions()...)...),
		egdebuild.Option.Runtime(func(r shell.Command) shell.Command {
			return r.Privileged()
		}),
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
		shell.Op(
			// gpg-agent mlocks secret key memory; if RLIMIT_MEMLOCK is capped in
			// this environment (e.g. a hardened shared runner) that can surface as
			// import failures. surface it unconditionally so we can compare against
			// a passing eg compute local run.
			shell.New("ulimit -l"),
			// seed gpg-agent.conf with debug logging before eggpg.Seed launches the
			// agent, so a failed import leaves a log with the agent's actual reason
			// instead of just gpg's generic client-side error.
			shell.Env().EnvironFrom(eggpg.Env(gpgoptions()...)...).Environ("GNUPGHOME", "/home/egd/.gnupg").New(
				`mkdir -p -m 700 "${GNUPGHOME}" && printf 'debug-all\nlog-file %s/gpg-agent.debug.log\n' "${GNUPGHOME}" > "${GNUPGHOME}/gpg-agent.conf"`,
			),
		),
		egbug.DebugFailure(
			eggpg.Seed(gpgoptions()...),
			eg.Sequential(
				eggpg.Debug(gpgoptions()...),
				shell.Op(
					shell.Env().Environ("GNUPGHOME", "/home/egd/.gnupg").New(`cat "${GNUPGHOME}/gpg-agent.debug.log"`).Lenient(true),
				),
			),
		),
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
