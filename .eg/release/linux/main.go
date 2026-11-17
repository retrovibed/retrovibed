package main

import (
	"context"
	"eg/compute/console"
	"eg/compute/dist"
	"eg/compute/egbugx"
	"eg/compute/maintainer"
	"eg/compute/release"
	"eg/compute/shallows"
	"eg/compute/tarballs"
	"log"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/x/wasi/egbug"
)

func main() {
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	deb := eg.Container(maintainer.Container)
	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eg.Build(deb.BuildFromFile(".eg/Containerfile")),
		eg.Module(
			ctx,
			deb,
			egbug.Debug(
				eg.Sequential(
					shallows.NeuralsBuild(),
					console.GenerateFlutter,
					console.CompileBinding,
					console.BuildLinux,
					shallows.Compile(),
					dist.Linux(),
					eg.Parallel(
						dist.AppImageBuild(),
						dist.FlatpakManifest(),
						tarballs.Tarchive,
					),
				),
				egbugx.Compute,
			),
		),
		// release.SmokeTest(),
		release.Release(),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
