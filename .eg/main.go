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
			eg.Sequential(
				eg.Parallel(
					shallows.Generate,
					console.Generate,
					shallows.NeuralsBuild(),
				),
				egbug.Log("GENERATE COMPLETED"),
				eg.Parallel(
					eg.Sequential(console.GenerateBinding, console.BuildLinux),
					shallows.Compile(),
				),
				egbug.Log("BUILD COMPLETED"),
				eg.Parallel(
					console.Tests,
					console.Linting,
					shallows.Test(),
					shallows.Linting,
					// console.RunDev("wlheadless-run -- flutter run --target lib/main.smoke.dart"),
				),
				egbug.Log("TESTS COMPLETED"),
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
