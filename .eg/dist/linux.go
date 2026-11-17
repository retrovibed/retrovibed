package dist

import (
	"eg/compute/console"
	"eg/compute/flathub"
	"eg/compute/release"
	"eg/compute/shallows"
	"eg/compute/tarballs"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egarch"
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

// Linux assembles the linux distribution: the console and binaries are installed into the tarball directory,
// packed into the tarball, and the appimage and flatpak manifest are generated from it.
func Linux() eg.OpFn {
	archive := tarballs.Retrovibed()
	return eg.Sequential(
		eg.Sequential(
			console.Install(),
			shallows.Install(),
			shell.Op(
				shell.Newf("cp --verbose -R .dist/linux/* %s", egtarball.Path(archive)),
				shell.Newf(
					"cat .dist/linux/usr/share/applications/space.retrovibe.Console.desktop | envsubst > %s/usr/share/applications/space.retrovibe.Console.desktop",
					egtarball.Path(archive),
				).
					Environ("VERSION", tarballs.Version()).
					Environ("ARCH", egarch.Host()),
				shell.Newf(
					"cat usr/share/applications/space.retrovibe.Console.desktop",
				).Directory(egtarball.Path(archive)),
			),
			flathub.Metainfo(),
		),
		egbug.Log("tarballing", egtarball.Path(archive), "->", egtarball.Archive(archive)),
		eg.Sequential(
			release.Tarball(),
		),
	)
}
