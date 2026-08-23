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
	"github.com/egdaemon/eg/runtime/wasi/shell"
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
					console.Generate,
					shallows.Generate,
				),
				console.GenerateBinding,
				shell.Op(
					shell.New("git diff > ${PATCH}").Environ("PATCH", egenv.CacheDirectory("codegen.patch")),
				),
			),
		),
		shell.Op(
			shell.New("git apply ${PATCH}").Environ("PATCH", egenv.CacheDirectory("codegen.patch")),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
