package main

import (
	"context"
	"eg/compute/flathub"
	"log"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()
	// example successful submittion: https://github.com/flathub/flathub/pull/8409/changes
	err := eg.Perform(ctx, flathub.Submit)
	if err != nil {
		log.Fatalln(err)
	}
}
