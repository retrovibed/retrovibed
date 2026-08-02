package main

import (
	"context"
	"log"

	"eg/compute/archlinux"

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
		eg.Build(eg.DefaultModule()),
		archlinux.Prepare,
		eg.Module(
			ctx,
			archlinux.AURRunner(),
			eg.Sequential(
				archlinux.Generate(egenv.CacheDirectory("PKGBUILD")),
				eg.Parallel(
					archlinux.Publish(egenv.CacheDirectory("PKGBUILD")),
					eggithub.Release(egenv.CacheDirectory("PKGBUILD")),
				),
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
