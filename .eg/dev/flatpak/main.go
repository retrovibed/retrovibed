package main

import (
	"context"
	"eg/compute/console"
	"eg/compute/maintainer"
	"eg/compute/release"
	"eg/compute/shallows"
	"eg/compute/tarballs"
	"log"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

// used to test flatpak build.
func main() {
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	deb := eg.Container(maintainer.Container)
	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eg.Build(deb.BuildFromFile(".eg/Containerfile")),
		eg.Module(
			ctx, deb,
			eg.Sequential(
				eg.Parallel(
					eg.Sequential(console.GenerateBinding, console.BuildLinux),
					shallows.Compile(),
				),
				eg.Parallel(
					console.Install,
					shallows.Install,
					shell.Op(
						shell.Newf("cp --verbose -R .dist/linux/* %s", egtarball.Path(tarballs.Retrovibed())),
					),
				),
				release.Tarball,
				eg.Parallel(
					console.FlatpakBuild,
				),
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
