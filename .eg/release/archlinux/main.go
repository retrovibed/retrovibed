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

	pkgbuild := egenv.DirectoryCache("PKGBUILD")

	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eg.Build(eg.DefaultModule()),
		archlinux.Prepare,
		eg.Module(
			ctx,
			archlinux.Runner(),
			eg.Sequential(
				archlinux.Generate(pkgbuild),
				archlinux.Build(pkgbuild),
				// archlinux.Publish(pkgbuild),
				eggithub.Draft(
					pkgbuild,
					egenv.DirectoryCache("pacman", "*.pkg.tar.*"),
				),
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
