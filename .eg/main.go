package main

import (
	"context"
	"eg/compute/console"
	"eg/compute/maintainer"
	"eg/compute/shallows"
	"log"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
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
		eg.Parallel(
			eg.Module(ctx, deb, console.GenerateBinding, console.Build),
			eg.Module(ctx, deb, shallows.Compile()),
		),
		eg.Parallel(
			eg.Module(ctx, deb, console.Tests),
			eg.Module(ctx, deb, console.Linting),
			eg.Module(ctx, deb, shallows.Test()),
		),
		// eg.Module(ctx, deb.OptionLiteral("--publish", "3000:3000"), www.Build, www.Webserver),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
