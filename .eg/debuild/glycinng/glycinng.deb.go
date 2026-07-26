package glycinng

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
	"github.com/egdaemon/eg/runtime/x/wasi/egdebuild"
)

//go:embed .debskel
var debskel embed.FS

const container = "retrovibed.deb.glycinng"

var (
	gcfg egdebuild.Config
)

func init() {
	c := eggit.EnvCommit()
	gcfg = egdebuild.New(
		"glycin-ng",
		"",
		egenv.CacheDirectory("glycinng"),
		egdebuild.Option.Maintainer(maintainer.Name, maintainer.Email),
		egdebuild.Option.SigningKeyID(maintainer.GPGFingerprint),
		egdebuild.Option.ChangeLogDate(c.Committer.When),
		egdebuild.Option.Version(fmt.Sprintf("%s.:autopatch:", version)),
		egdebuild.Option.Description("glycin-ng", "in-process, sandboxed (Landlock/seccomp) image decoder; drop-in replacement for libglycin-2-0/glycin-loaders that avoids their bwrap/subprocess-spawn model"),
		egdebuild.Option.Debian(errorsx.Must(fs.Sub(debskel, ".debskel"))),
		egdebuild.Option.DependsBuild("cargo", "rustc", "lld"),
		egdebuild.Option.Envvar("PACKAGE_VERSION", version),
		egdebuild.Option.Envvar("GIT_COMMIT_HASH", c.Hash.String()),
	)
}

// Prepare only builds the debuild container image. Download runs later, as a
// module step inside that container (via Runner()), because it needs the
// cargo/rustc that image provides for `cargo vendor` -- the outer/default eg
// container Prepare itself runs in does not have a Rust toolchain at all.
func Prepare(ctx context.Context, o eg.Op) error {
	return egdebuild.Prepare(Runner(), errorsx.Must(fs.Sub(debskel, ".debskel")))(ctx, o)
}

// container for this package.
func Runner() eg.ContainerRunner {
	return eg.Container(container)
}

// Only resolute/questing are targeted: glycin-ng needs edition 2024 / rustc
// >= 1.88 (see Cargo.toml's workspace.package.rust-version), which is newer
// than the stock "cargo"/"rustc" apt packages on noble or jammy can satisfy
// as a Build-Depends on Launchpad's own builders for those distros.
//
// Each distro gets a source-only build (what actually gets uploaded --
// Launchpad rebuilds from source) plus a local binary build, so a broken
// build fails here instead of silently on Launchpad after upload.
func Build(ctx context.Context, o eg.Op) error {
	return eg.Sequential(
		egdebuild.Build(gcfg, egdebuild.Option.Distro(egdebuild.UbuntuLatestCodename), egdebuild.Option.BuildBinary(20*time.Minute), egdebuild.Option.NoLint()),
		egdebuild.Build(gcfg, egdebuild.Option.Distro(egdebuild.UbuntuLatestCodename), egdebuild.Option.NoLint()),
		egdebuild.Build(gcfg, egdebuild.Option.Distro("questing"), egdebuild.Option.NoLint()),
	)(ctx, o)
}

func Upload(ctx context.Context, o eg.Op) error {
	return egdebuild.UploadDPut(gcfg, errorsx.Must(fs.Sub(debskel, ".debskel")))(ctx, o)
}
