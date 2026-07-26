package main

import (
	"context"
	"log"

	"eg/compute/debuild/glycinng"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
)

func main() {
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	err := eg.Perform(
		ctx,
		glycinng.Prepare,
		eg.Module(
			ctx,
			glycinng.Runner(),
			glycinng.Download,
			glycinng.Build,
			glycinng.Upload,
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
