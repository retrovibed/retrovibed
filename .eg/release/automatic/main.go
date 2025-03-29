package main

import (
	"context"
	"eg/compute/console"
	"eg/compute/maintainer"
	"eg/compute/release"
	"eg/compute/shallows"
	"log"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/x/wasi/egtarball"
)

func main() {
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	deb := eg.Container(maintainer.Container)
	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eg.Build(deb.BuildFromFile(".eg/Containerfile")),
		eg.Parallel(
			eg.Module(
				ctx,
				deb,
				console.Generate,
			),
			eg.Module(
				ctx,
				deb,
				shallows.Generate,
			),
		),
		eg.Module(ctx, deb, eg.Parallel(
			eg.Sequential(console.GenerateBinding, console.Build),
			shallows.Compile(),
		)),
		eg.Parallel(
			eg.Module(ctx, deb, console.Tests),
			eg.Module(ctx, deb, console.Linting),
			eg.Module(ctx, deb, shallows.Test()),
		),
		egtarball.Clean(
			eg.Module(
				ctx, deb,
				eg.Parallel(
					console.Install,
					shallows.Install,
				),
				release.Tarball,
				eg.Parallel(
					shallows.FlatpakManifest,
					console.FlatpakManifest,
				),
			),
			release.Release,
		),
		eg.Module(
			ctx, deb.OptionLiteral("--privileged"),
			eg.Parallel(
				console.FlatpakBuild,
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
