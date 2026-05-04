package release

import (
	"context"
	"eg/compute/maintainer"
	"eg/compute/tarballs"
	"time"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egdmg"
	"github.com/egdaemon/eg/runtime/x/wasi/eggithub"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

func Release(b *tarballs.Build) eg.OpFn {
	return eggithub.Release(
		egtarball.Archive(tarballs.Retrovibed(b)),
		egtarball.Archive(tarballs.RetrovibedSource()),
		egenv.CacheDirectory(tarballs.Flatpak(b)),
		egenv.CacheDirectory(tarballs.AppImage(b)),
	)
}

func AppImageBuild(b *tarballs.Build) eg.OpFn {
	appimage := tarballs.AppImage(b)
	builddir := egenv.CacheDirectory(tarballs.AppImageBuild(b))

	return shell.Op(
		shell.Newf("echo %s", tarballs.Version()),
		shell.Newf("appimage-builder --help"),
		shell.Newf("mkdir -p %s", builddir),
		shell.Newf(
			"appimage-builder --recipe .dist/AppImageBuilder.yml --skip-tests --build-dir %s",
			builddir,
		).Timeout(20*time.Minute).
			Environ("APPDIR", egtarball.Path(tarballs.Retrovibed(b))).
			Environ("VERSION", tarballs.Version()).
			Environ("APT_ARCH", b.Arch).
			Environ("APPIMAGE_ARCH", tarballs.ArchGoToMachine(b.Arch)).
			Environ("APPIMAGE_FILE_NAME", appimage).
			Environ("GPG_ID", maintainer.GPGID),
		shell.Newf("mv %s* %s", appimage, egenv.CacheDirectory(appimage)).Debug(),
	)
}

func DistroBuilds(ctx context.Context, op eg.Op) error {
	podman := shell.Runtime().Environ("XFG_DATA_HOME", egenv.CacheDirectory())
	return eg.Sequential(
		eg.Parallel(
			shell.Op(
				podman.New("echo ---------------------------------------------------------------"),
				podman.New("podman build -t retrovibed.flatpak.distro.check.ubuntu.noble -f .dist/distrobuilds/ubuntu.noble").Privileged(),
				podman.New("echo ---------------------------------------------------------------"),
			),
			// shell.Op(shell.New("podman build -tag retrovibed.flatpak.distro.check.ubuntu.oracular -f .dist/distrobuilds/ubuntu.oracular")),
		),
		eg.Parallel(
			shell.Op(podman.New("podman run --privileged --rm --volume .eg.cache/flatpak.client.yml:/retrovibed.client.yml:ro retrovibed.flatpak.distro.check.ubuntu.noble cat /retrovibed.client.yml").Privileged()),
			shell.Op(podman.New("podman run --privileged --rm --volume .eg.cache/flatpak.client.yml:/retrovibed.client.yml:ro retrovibed.flatpak.distro.check.ubuntu.noble").Privileged()),
		),
	)(ctx, op)
}

func Tarball(b *tarballs.Build) eg.OpFn {
	archive := tarballs.Retrovibed(b)
	return eg.Sequential(
		egtarball.Pack(archive),
		egtarball.SHA256Op(archive),
	)
}

func DarwinDmg(b1 *tarballs.Build) eg.OpFn {
	path := tarballs.Retrovibed(b1)
	dmg := egdmg.New(path, egdmg.OptionBuildDir(egenv.CacheDirectory()))
	return eg.Sequential(
		egdmg.Build(dmg, egtarball.Path(tarballs.Retrovibed(b1))),
		shell.Op(
			shell.Newf("echo mv %s %s", egdmg.Path(path), egenv.CacheDirectory("retrovibed.darwin.arm64.dmg")),
			shell.Newf("mv %s %s", egdmg.Path(path), egenv.CacheDirectory("retrovibed.darwin.arm64.dmg")),
		),
	)
}
