package main

import (
	"context"
	"eg/compute/flathub"
	"log"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
)

func main() {
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	err := eg.Perform(ctx, flathub.Submit)
	if err != nil {
		log.Fatalln(err)
	}
}
