package duckdb

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"eg/compute/errorsx"
	"eg/compute/maintainer"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
	"github.com/egdaemon/eg/runtime/x/wasi/egccache"
	"github.com/egdaemon/eg/runtime/x/wasi/egdebuild"
	"github.com/egdaemon/eg/runtime/x/wasi/eggpg"
)

//go:embed .debskel
var debskel embed.FS

const (
	container = "retrovibed.deb.duckdb"
	version   = "1.5.5"
)

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
	egccache.CacheDirectory()
	c := eggit.EnvCommit()
	gcfg = egdebuild.New(
		"duckdb",
		"",
		egenv.DirectoryCache("duckdb"),
		egdebuild.Option.Maintainer(maintainer.Name, maintainer.Email),
		egdebuild.Option.SigningKeyID(maintainer.GPGFingerprint),
		egdebuild.Option.ChangeLogDate(c.Committer.When),
		egdebuild.Option.Version(fmt.Sprintf("%s.:autopatch:", version)),
		egdebuild.Option.Description("duckdb", "embeddable columnar database"),
		egdebuild.Option.Debian(errorsx.Must(fs.Sub(debskel, ".debskel"))),
		egdebuild.Option.DependsBuild("rsync", "curl", "tree", "ca-certificates", "cmake", "ninja-build", "libssl-dev", "git"),
		egdebuild.Option.Envvar("PACKAGE_VERSION", version),
		egdebuild.Option.Envvar("GIT_COMMIT_HASH", c.Hash.String()),
	)
}

// withgpg adds the gpg environment to the config. resolved at build time instead of in init
// because eggpg.Env panics when the EG_GPG_KEYRING_* variables are unset, which would break
// every module that merely imports this package.
func withgpg(c egdebuild.Config) egdebuild.Config {
	return egdebuild.From(c, egdebuild.Option.Environ(eggpg.Env(gpgoptions()...)...))
}

func Prepare(ctx context.Context, o eg.Op) error {
	return eg.Parallel(
		Download,
		egdebuild.Prepare(Runner(), errorsx.Must(fs.Sub(debskel, ".debskel"))),
	)(ctx, o)
}

// container for this package.
func Runner() eg.ContainerRunner {
	return eg.Container(container)
}

func Build(ctx context.Context, o eg.Op) error {
	gcfg := withgpg(gcfg)
	return eg.Sequential(
		egbug.DebugFailure(
			eggpg.Seed(gpgoptions()...),
			eggpg.Debug(gpgoptions()...),
		),
		eg.Parallel(
			// build the package to improve the chances it'll actually build in within ubuntu launchpad.
			// egdebuild.Build(
			// 	gcfg,
			// 	egdebuild.Option.Distro(latest),
			// 	egdebuild.Option.BuildBinary(20*time.Minute),
			// 	egdebuild.Option.Environ(egccache.Env()...),
			// 	egdebuild.Option.NoLint(),
			// ),
			egdebuild.Build(gcfg, egdebuild.Option.Distro(egdebuild.UbuntuLatestCodename), egdebuild.Option.NoLint()),
			egdebuild.Build(gcfg, egdebuild.Option.Distro("questing"), egdebuild.Option.NoLint()),
			egdebuild.Build(gcfg, egdebuild.Option.Distro("noble"), egdebuild.Option.NoLint()),
			egdebuild.Build(gcfg, egdebuild.Option.Distro("jammy")),
		),
	)(ctx, o)
}

func Upload(ctx context.Context, o eg.Op) error {
	return egdebuild.UploadDPut(withgpg(gcfg), errorsx.Must(fs.Sub(debskel, ".debskel")))(ctx, o)
}
