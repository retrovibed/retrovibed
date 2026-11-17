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
	// must match Dir::Cache::archives in the host apt config.
	aptarchives := egenv.DirectoryCache(_eg.DefaultModuleDirectory(), "apt", "archives")
	return shell.Op(
		// appimage-builder runs apt with its own APT_CONFIG (ignoring /etc/apt), but it does read
		// the apt.conf.d inside the build dir. sync the host's apt configuration into it so apt
		// behaves the same as the host (retries for transient failures, e.g. 503s from the
		// launchpad PPA). docker-clean is excluded: it disables the package cache dirs that
		// appimage-builder points at its build dir.
		shell.Newf(
			"rsync -a --mkpath --exclude=docker-clean /etc/apt/apt.conf.d/ %s/",
			filepath.Join(builddir, "apt", "apt.conf.d"),
		),
		// keep downloaded packages in the persistent workload cache so reruns skip what's already fetched.
		shell.Newf(
			`printf 'Dir::Cache::archives "%s";\n' > %s`,
			aptarchives,
			filepath.Join(builddir, "apt", "apt.conf.d", "99.retro.apt"),
		),
		// appimage-builder hardcodes <builddir>/apt/archives when locating downloaded debs, so
		// symlink that to the cache dir configured above so both agree.
		shell.Newf("mkdir -p %s", aptarchives),
		shell.Newf("rm -rf %s", filepath.Join(builddir, "apt", "archives")),
		shell.Newf("ln -s %s %s", aptarchives, filepath.Join(builddir, "apt", "archives")),
		shell.Newf(
			"appimage-builder --recipe %s --skip-tests --build-dir %s",
			egenv.DirectoryWorking(".dist", "AppImageBuilder.yml"),
			builddir,
		).Directory(egenv.DirectoryWorkload()).
			Timeout(egenv.TTL()).
			Environ("APPDIR", egtarball.Path(tarballs.Retrovibed())).
			Environ("VERSION", tarballs.Version()).
			Environ("APT_ARCH", egarch.Host()).
			Environ("APPIMAGE_ARCH", egarch.POSIX()).
			Environ("APPIMAGE_COMP", "zstd").
			Environ("APPIMAGE_FILE_NAME", appimage).
			Environ("GPG_ID", maintainer.GPGID),
		shell.Newf("cp %s* %s", egenv.DirectoryWorkload(appimage), egenv.DirectoryCache()),
	)
}
