// Package main run just the linters.
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
	log.SetFlags(log.Lshortfile | log.LUTC | log.Ltime)
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
			eg.Parallel(
				console.Linting,
				shallows.Linting,
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
