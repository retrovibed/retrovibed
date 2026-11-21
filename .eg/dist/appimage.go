package dist

import (
	"path/filepath"

	"eg/compute/maintainer"
	"eg/compute/tarballs"

	_eg "github.com/egdaemon/eg"
	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egarch"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

func AppImageBuild() eg.OpFn {
	// packaging metainfo.
	// https://docs.appimage.org/packaging-guide/optional/appstream.html
	// runtime to look into.
	// https://github.com/pkgforge-dev/Anylinux-AppImages
	appimage := tarballs.AppImage()
	builddir := egenv.DirectoryWorkload("build", tarballs.AppImageBuild())
	// appimage-builder modifies its AppDir in place (deploys packages, patches binaries), so it works on
	// a clone of the tarball directory to keep the tarball, and the builds sharing it, untouched.
	appdir := egtarball.Path(tarballs.AppImageBuild())
	// persistent cache for build independent data (apt debs and indexes, AppRun binaries, runtime).
	cachedir := egenv.DirectoryCache(_eg.DefaultModuleDirectory())
	return shell.Op(
		shell.Newf("rsync -a --mkpath %s/ %s/", egtarball.Path(tarballs.Retrovibed()), appdir),
		// appimage-builder runs apt with its own APT_CONFIG (ignoring /etc/apt), but it does read
		// the apt.conf.d inside the build dir. sync the host's apt configuration into it so apt
		// behaves the same as the host (retries for transient failures, e.g. 503s from the
		// launchpad PPA). docker-clean is excluded: it disables the package cache dirs that
		// appimage-builder points at its cache dir.
		shell.Newf(
			"rsync -a --mkpath --exclude=docker-clean /etc/apt/apt.conf.d/ %s/",
			filepath.Join(builddir, "apt", "apt.conf.d"),
		),
		shell.Newf(
			"appimage-builder --recipe %s --skip-tests",
			egenv.DirectoryWorking(".dist", "AppImageBuilder.yml"),
		).Directory(egenv.DirectoryWorkload()).
			Timeout(egenv.TTL()).
			Environ("APPIMAGE_BUILDER_BUILD_DIR", builddir).
			Environ("APPIMAGE_BUILDER_CACHE_DIR", cachedir).
			Environ("APPDIR", appdir).
			Environ("VERSION", tarballs.Version()).
			Environ("APT_ARCH", egarch.Host()).
			Environ("APPIMAGE_ARCH", egarch.POSIX()).
			Environ("APPIMAGE_COMP", "zstd").
			Environ("APPIMAGE_FILE_NAME", appimage).
			Environ("GPG_ID", maintainer.GPGID),
		shell.Newf("cp %s* %s", egenv.DirectoryWorkload(appimage), egenv.DirectoryCache()),
	)
}
