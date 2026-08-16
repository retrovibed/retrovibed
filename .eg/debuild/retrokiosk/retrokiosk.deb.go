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
)

//go:embed .debskel
var debskel embed.FS

var (
	gcfg egdebuild.Config
)

func init() {
	c := eggit.EnvCommit()
	gcfg = egdebuild.New(
		"retrokiosk",
		"",
		egenv.WorkingDirectory(".eg", "debuild", "retrokiosk", "rootfs"),
		egdebuild.Option.Maintainer(maintainer.Name, maintainer.Email),
		egdebuild.Option.SigningKeyID(maintainer.GPGFingerprint),
		egdebuild.Option.ChangeLogDate(c.Committer.When),
		egdebuild.Option.Version("0.0.:autopatch:"),
		egdebuild.Option.Description("retrovibed console kiosk appliance", "boots a machine directly into a fullscreen retrovibed Console session, no desktop or login required"),
		egdebuild.Option.Debian(errorsx.Must(fs.Sub(debskel, ".debskel"))),
		egdebuild.Option.DependsBuild("rsync", "tree", "dpkg-dev", "gettext-base"),
		egdebuild.Option.Depends("labwc", "greetd", "seatd", "libfuse2", "zsync"),
	)
}

func Prepare(ctx context.Context, o eg.Op) error {
	return eg.Parallel(
		egdebuild.Prepare(Runner(), nil),
	)(ctx, o)
}

func Runner() eg.ContainerRunner {
	return egdebuild.Runner()
}

func Build(ctx context.Context, o eg.Op) error {
	return eg.Parallel(
		egdebuild.Build(gcfg, egdebuild.Option.Distro(egdebuild.UbuntuLatestCodename)),
		egdebuild.Build(
			gcfg,
			egdebuild.Option.Distro(egdebuild.UbuntuLatestCodename),
			egdebuild.Option.BuildBinary(20*time.Minute),
			egdebuild.Option.NoLint(),
		),
	)(ctx, o)
}

func Upload(ctx context.Context, o eg.Op) error {
	return egdebuild.UploadDPut(gcfg, errorsx.Must(fs.Sub(debskel, ".debskel")), egdebuild.Option.Timeout(20*time.Minute))(ctx, o)
}

// Verify installs the .deb produced by Build (staged by egdebuild under
// EphemeralDirectory("deb.retrokiosk"), matching the root egdebuild.Build
// computes internally) and checks the services it enables come up. There's
// no real display/seat backend in a container, so this only proves service
// startup, not that pixels reach a screen. It checks enablement of the user
// unit rather than starting it live: that needs a working systemd-logind
// user session, which a plain container doesn't have.
func Verify(ctx context.Context, op eg.Op) error {
	debroot := egenv.EphemeralDirectory("deb.retrokiosk")
	runtime := shell.Runtime().Timeout(egenv.TTL())

	return eg.Sequential(
		shell.Op(
			runtime.Newf("apt-get install -y $(find %s -maxdepth 2 -name 'retrokiosk_*.deb' | head -n1)", debroot).Privileged(),
		),
		shell.Op(
			runtime.New("test -L /etc/systemd/system/display-manager.service"),
			runtime.New("test -L /etc/systemd/user/graphical-session.target.wants/retrokiosk.service"),
		),
		egbug.DebugFailure(
			shell.Op(
				runtime.New("systemctl start retrokiosk-greetd.service").Privileged(),
				runtime.New("systemctl is-active retrokiosk-greetd.service").Privileged(),
			),
			shell.Op(
				runtime.New("journalctl --since -2m -u retrokiosk-greetd.service").Privileged(),
			),
		),
		shell.Op(
			runtime.New("systemctl status").Privileged(),
		),
	)(ctx, op)
}
