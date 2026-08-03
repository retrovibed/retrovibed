package release

import (
	"context"
	"eg/compute/maintainer"
	"eg/compute/tarballs"
	"os"
	"path/filepath"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egdmg"
	"github.com/egdaemon/eg/runtime/x/wasi/eggithub"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

func Release(b *tarballs.Build) eg.OpFn {
	return eggithub.Release(
		egtarball.Archive(tarballs.RetrovibedSource()),
		egtarball.Archive(tarballs.Retrovibed(b)),
		egenv.CacheDirectory(tarballs.Flatpak(b)),
		egenv.WorkspaceDirectory(tarballs.AppImage(b)),
		egenv.WorkspaceDirectory(tarballs.AppImageZsync(b)),
	)
}

func AppImageBuild(b *tarballs.Build) eg.OpFn {
	// packaging metainfo.
	// https://docs.appimage.org/packaging-guide/optional/appstream.html
	// runtime to look into.
	// https://github.com/pkgforge-dev/Anylinux-AppImages
	appimage := tarballs.AppImage(b)
	builddir := egenv.CacheDirectory(tarballs.AppImageBuild(b))
	return shell.Op(
		shell.Newf("mkdir -p %s", builddir),
		shell.Newf(
			"appimage-builder --recipe %s --skip-tests --build-dir %s",
			egenv.WorkingDirectory(".dist", "AppImageBuilder.yml"),
			builddir,
		).Directory(egenv.WorkspaceDirectory()).
			Attempts(3).
			Environ("APPDIR", egtarball.Path(tarballs.Retrovibed(b))).
			Environ("VERSION", tarballs.Version()).
			Environ("APT_ARCH", b.Arch).
			Environ("APPIMAGE_ARCH", tarballs.ArchGoToMachine(b.Arch)).
			Environ("APPIMAGE_FILE_NAME", appimage).
			Environ("GPG_ID", maintainer.GPGID),
		// copy is only needed temporarily while we flesh out app image support.
		shell.Newf("cp %s* %s", egenv.WorkspaceDirectory(appimage), egenv.CacheDirectory()),
	)
}

// SmokeTest builds every Containerfile in .dist/distrobuilds (each becomes
// retrovibed.<filename>) in parallel, then runs test inside each one in
// parallel. Used to verify build output (e.g. the AppImage) behaves
// correctly in clean environments rather than ones shaped by whatever the
// build container happens to already have installed.
func SmokeTest(test eg.OpFn) eg.OpFn {
	const dir = ".dist/distrobuilds"

	return func(ctx context.Context, o eg.Op) error {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		containers := make([]eg.ContainerRunner, 0, len(entries))
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			containers = append(containers, eg.Container("retrovibed."+entry.Name()).BuildFromFile(filepath.Join(dir, entry.Name())))
		}

		builds := make([]eg.OpFn, 0, len(containers))
		runs := make([]eg.OpFn, 0, len(containers))
		for _, container := range containers {
			builds = append(builds, eg.Build(container))
			runs = append(runs, eg.Module(ctx, container, test))
		}

		return eg.Sequential(
			eg.Parallel(builds...),
			eg.Parallel(runs...),
		)(ctx, o)
	}
}

// AppImageSmokeTest runs the built AppImage inside every container defined in
// .dist/distrobuilds, none of which have mpv/ffmpeg/libmpv preinstalled -- so
// this can't pass by accident off libraries the AppImage should be bundling
// itself, it only proves the AppImage is self-contained.
func AppImageSmokeTest(b *tarballs.Build) eg.OpFn {
	appimage := egenv.CacheDirectory(tarballs.AppImage(b))
	return SmokeTest(
		shell.Op(shell.Newf("wlheadless-run -- %s --appimage-extract-and-run --smoke", appimage)),
	)
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
