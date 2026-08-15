package main

import (
	"context"
	"log"

	"eg/compute/debuild/retrokiosk"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LUTC | log.Ltime)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	c1 := eg.Container("retrovibed.retrokiosk.test")

	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eg.Build(eg.DefaultModule()),
		eg.Build(c1.BuildFromFile(".eg/tests/retrokiosk/Containerfile")),
		eg.Module(
			ctx,
			c1,
			retrokiosk.Build,
			retrokiosk.Verify,
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
