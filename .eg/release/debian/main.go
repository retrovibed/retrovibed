package main

import (
	"context"
	"eg/compute/debian"
	"log"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/x/wasi/eggithub"
)

func main() {
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eg.Parallel(
			debian.Prepare,
		),
		eg.Module(
			ctx,
			debian.Runner(),
			debian.Build,
			debian.Upload,
			debian.Copy(egenv.DirectoryWorkload(".dist")),
		),
		eggithub.Draft(
			egenv.DirectoryWorkload(".dist", "retrovibed_*.deb"),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
