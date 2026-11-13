// Package retrokiosk builds a debian package that turns a machine into a
// single-purpose kiosk appliance booting straight into the retrovibed Console.
package retrokiosk

import (
	"context"
	"embed"
	"io/fs"
	"time"

	"eg/compute/errorsx"
	"eg/compute/maintainer"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
	"github.com/egdaemon/eg/runtime/x/wasi/egdebuild"
	"github.com/egdaemon/eg/runtime/x/wasi/eggpg"
)

//go:embed .debskel
var debskel embed.FS

var (
	gcfg egdebuild.Config
)

func cachedir() string {
	return egenv.DirectoryWorkload(".dist", "retrokiosk")
}

func gpgoptions() []eggpg.Option {
	// eggpg.Options() takes GNUPGHOME from the environment, which inside the container is the
	// host-provided /eg.mnt/.gnupg (not writable by us: "unsafe ownership", "Permission denied").
	// pin the isolated per-container home instead.
	return eggpg.Options().IgnoreLocalGNU().Home("/home/egd/.gnupg")
}

func init() {
	c := eggit.EnvCommit()
	gcfg = egdebuild.New(
		"retrokiosk",
		"",
		egenv.DirectoryWorking(".eg", "debuild", "retrokiosk", "rootfs"),
		egdebuild.Option.Maintainer(maintainer.Name, maintainer.Email),
		egdebuild.Option.SigningKeyID(maintainer.GPGFingerprint),
		egdebuild.Option.ChangeLogDate(c.Committer.When),
		egdebuild.Option.Version("0.0.:autopatch:"),
		egdebuild.Option.Description("retrovibed console kiosk appliance", "boots a machine directly into a fullscreen retrovibed Console session, no desktop or login required"),
		egdebuild.Option.Debian(errorsx.Must(fs.Sub(debskel, ".debskel"))),
		egdebuild.Option.DependsBuild("rsync", "tree"),
		egdebuild.Option.Depends("avahi-daemon", "cage", "kbd", "libfuse2", "pipewire-audio", "retrozsync", "ssh"),
		egdebuild.Option.Environ(eggpg.Env(gpgoptions()...)...),
	)
}

func Prepare(ctx context.Context, o eg.Op) error {
	return eg.Parallel(
		// egdebuild.Prepare(Runner(), nil),
		egdebuild.Prepare(Runner(), errorsx.Must(fs.Sub(debskel, ".debskel"))),
	)(ctx, o)
}

func Runner() eg.C8S {
	return egdebuild.Runner()
}

func Build(ctx context.Context, o eg.Op) error {
	return eg.Sequential(
		// override eg compute local's setting of the home directory.
		// ideally we should default to this except when performing local
		// compute workloads.
		egbug.DebugFailure(
			eggpg.Seed(gpgoptions()...),
			eggpg.Debug(gpgoptions()...),
		),
		egdebuild.Build(
			gcfg,
			egdebuild.Option.Distro(egdebuild.UbuntuLatestCodename),
			egdebuild.Option.NoLint(),
		),
		egdebuild.Build(
			gcfg,
			egdebuild.Option.Distro(egdebuild.UbuntuLatestCodename),
			egdebuild.Option.BuildBinary(20*time.Minute),
			egdebuild.Option.NoLint(),
		),
	)(ctx, o)
}

func Upload(ctx context.Context, o eg.Op) error {
	return egdebuild.UploadDPut(
		gcfg,
		errorsx.Must(fs.Sub(debskel, ".debskel")),
		egdebuild.Option.Timeout(20*time.Minute),
	)(ctx, o)
}

// Verify installs the .deb produced by Build (staged by egdebuild under
// EphemeralDirectory("deb.retrokiosk"), matching the root egdebuild.Build
// computes internally) and checks the package sets up cleanly: unit installed
// and preset-enabled, unit file well-formed. It does not start the service —
// that requires a real VT/DRM, which test containers don't have.
func Verify(ctx context.Context, op eg.Op) error {
	debroot := egenv.DirectoryTemp("deb.retrokiosk")
	runtime := shell.Runtime().Timeout(egenv.TTL())

	return eg.Sequential(
		shell.Op(
			runtime.Newf("apt-get install -y $(find %s -maxdepth 2 -name 'retrokiosk_*.deb' | head -n1)", debroot).Privileged(),
		),
		shell.Op(
			runtime.New("test -L /etc/systemd/system/graphical.target.wants/retrokiosk@3.service"),
		),
		shell.Op(
			runtime.New("systemd-analyze verify /usr/lib/systemd/system/retrokiosk@.service").Privileged(),
		),
	)(ctx, op)
}

func Copy(dst string) eg.OpFn {
	return egdebuild.Copy(gcfg, dst)
}
