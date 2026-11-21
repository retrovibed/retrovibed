package release

import (
	"context"
	"eg/compute/tarballs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/egunsafe/ffiegcontainer"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egarch"
	"github.com/egdaemon/eg/runtime/x/wasi/egdmg"
	"github.com/egdaemon/eg/runtime/x/wasi/eggithub"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

// Dist rsyncs files into the persistent .dist cache directory
// (egenv.CacheDirectory(".dist")), creating it if necessary. Useful for
// leaving a locally-inspectable copy of build outputs (e.g. release
// APK/AAB) alongside the other packaging artifacts already dropped there.
func Dist(files ...string) eg.OpFn {
	return shell.Op(
		shell.Newf("rsync --mkpath -av %s %s/", strings.Join(files, " "), egenv.DirectoryCache(".dist")),
	)
}

func Release() eg.OpFn {
	return eg.Sequential(
		eggithub.Draft(
			egtarball.Archive(tarballs.RetrovibedSource()),
			egtarball.Archive(tarballs.Retrovibed()),
			egenv.DirectoryCache(tarballs.Flatpak()),
			egenv.DirectoryWorkload(tarballs.AppImage()),
			egenv.DirectoryWorkload(tarballs.AppImageZsync()),
		),
	)
}

// SmokeTest builds every Containerfile in .dist/distrobuilds (each becomes
// retrovibed.smoke.<filename>) in parallel, then runs the AppImage inside
// each one sequentially (they'd otherwise fight over the same bind-mounted
// host /root's retrovibed config db). None of these containers have
// mpv/ffmpeg/libmpv preinstalled -- so this can't pass by accident off
// libraries the AppImage should be bundling itself, it only proves the
// AppImage is self-contained.
func SmokeTest() eg.OpFn {
	return func(ctx context.Context, o eg.Op) error {
		const dir = ".dist/distrobuilds"
		appimage := egenv.DirectoryCache(tarballs.AppImage())
		if egarch.Host() != "amd64" {
			log.Println("skipping smoke tests on non-amd64 systems")
			return nil
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		builds := make([]eg.OpFn, 0, len(entries))
		runs := make([]eg.OpFn, 0, len(entries))
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			image := "retrovibed.smoke." + entry.Name()
			container := eg.Container(image).BuildFromFile(filepath.Join(dir, entry.Name()))
			builds = append(builds, eg.Build(container))
			runs = append(runs, func(ctx context.Context, _ eg.Op) error {
				return ffiegcontainer.Run(
					ctx,
					image,
					appimage,
					[]string{"wlheadless-run", "--", "/retrovibed.AppImage", "--appimage-extract-and-run", "--smoke"},
					[]string{"--replace", "--volume", appimage + ":/retrovibed.AppImage:ro"},
				)
			})
		}

		return eg.Sequential(
			eg.Parallel(builds...),
			eg.Sequential(runs...),
		)(ctx, o)
	}
}

func Tarball() eg.OpFn {
	archive := tarballs.Retrovibed()
	return eg.Sequential(
		egtarball.Pack(archive),
		egtarball.SHA256Op(archive),
		// local runs only (eg compute local, where cloning is disabled): leave a copy of the packed tarball
		// in the persistent cache directory so its contents can be extracted and reused
		// (e.g. as APPDIR when testing appimage-builder locally).
		eg.When(
			!egenv.CloneEnabled(),
			shell.Op(
				shell.Newf("cp %s %s", egtarball.Archive(archive), egenv.DirectoryCache()),
			),
		),
	)
}

func DarwinDmg() eg.OpFn {
	path := tarballs.Retrovibed()
	dmg := egdmg.New(path, egdmg.OptionBuildDir(egenv.DirectoryCache()))
	return eg.Sequential(
		egdmg.Build(dmg, egtarball.Path(tarballs.Retrovibed())),
		shell.Op(
			shell.Newf("echo mv %s %s", egdmg.Path(path), egenv.DirectoryCache("retrovibed.darwin.arm64.dmg")),
			shell.Newf("mv %s %s", egdmg.Path(path), egenv.DirectoryCache("retrovibed.darwin.arm64.dmg")),
		),
	)
}
