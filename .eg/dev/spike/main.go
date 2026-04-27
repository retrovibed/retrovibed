// Package main run just the linters.
package main

import (
	"context"
	"eg/compute/maintainer"
	"eg/compute/tarballs"
	"log"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
	"github.com/egdaemon/eg/runtime/wasi/shell"
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
			eg.Sequential(
				tarballs.Tarchive,
				shell.Op(
					shell.Newf("ls -lha %s", egenv.WorkspaceDirectory()),
					shell.Newf("file %s", tarballs.RetrovibedSource()),
					shell.Newf("tar -tvf %s", tarballs.RetrovibedSource()),
				),
			),
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
