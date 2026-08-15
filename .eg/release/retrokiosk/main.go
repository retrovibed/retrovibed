package main

import (
	"context"
	"log"

	"eg/compute/debuild/retrokiosk"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
)

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile | log.LUTC)
}

func main() {
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		retrokiosk.Prepare,
		eg.Module(
			ctx,
			retrokiosk.Runner(),
			retrokiosk.Build,
			retrokiosk.Upload,
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
