package main

import (
	"context"
	"log"

	"eg/compute/archlinux"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
)

func main() {
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eg.Build(eg.DefaultModule()),
		archlinux.Prepare,
		eg.Module(ctx, archlinux.AURRunner(), archlinux.Publish),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
