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
		eg.Module(
			ctx,
			deb,
			eg.Sequential(
				eg.Parallel(
					shallows.Generate,
					console.Generate,
				),
				eg.Parallel(
					eg.Sequential(console.GenerateBinding, console.BuildLinux),
					shallows.Compile(),
				),
				eg.Parallel(
					console.Tests,
					console.Linting,
					shallows.Test(),
					shallows.Linting,
				),
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
