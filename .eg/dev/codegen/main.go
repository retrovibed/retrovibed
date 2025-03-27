package main

import (
	"context"
	"log"

	"eg/compute/console"
	"eg/compute/maintainer"
	"eg/compute/shallows"

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
			eg.Module(ctx, deb, console.Generate),
			eg.Module(ctx, deb, shallows.Generate),
		),
		eg.Module(ctx, deb, console.GenerateBinding),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
