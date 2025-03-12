package main

import (
	"context"
	"eg/compute/fractal"
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

	deb := eg.Container("fractal.ubuntu.24.10")
	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eg.Build(deb.BuildFromFile(".eg/Containerfile")),
		eg.Parallel(
			eg.Module(
				ctx,
				deb,
				fractal.Generate,
			),
			eg.Module(
				ctx,
				deb,
				shallows.Generate,
			),
		),
		eg.Parallel(
			eg.Module(ctx, deb, fractal.Build),
			eg.Module(
				ctx,
				deb,
				shallows.Compile(),
			),
		),
		eg.Parallel(
			eg.Module(ctx, deb, fractal.Tests),
			eg.Module(ctx, deb, fractal.Linting),
			eg.Module(ctx, deb, shallows.Test()),
		),
		egtarball.Clean(
			eg.Module(
				ctx, deb,
				eg.Parallel(
					fractal.Install,
					shallows.Install,
				),
				release.Tarball,
				eg.Parallel(
					shallows.FlatpakManifest,
					fractal.FlatpakManifest,
				),
			),
			release.Release,
		),
		eg.Module(
			ctx, deb.OptionLiteral("--privileged"),
			eg.Parallel(
				fractal.FlatpakBuild,
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
