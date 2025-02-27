package main

import (
	"context"
	"eg/compute/fractal"
	"eg/compute/release"
	"log"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
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
			eg.Module(ctx, deb, fractal.Generate),
		),
		eg.Parallel(
			eg.Module(ctx, deb, fractal.Build),
		),
		eg.Parallel(
			eg.Module(ctx, deb, fractal.Tests),
			eg.Module(ctx, deb, fractal.Linting),
		),
		eg.Parallel(
			eg.Module(ctx, deb, release.Flatpak),
		),
		// eg.Module(ctx, deb.OptionLiteral("--publish", "3000:3000"), www.Build, www.Webserver),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
