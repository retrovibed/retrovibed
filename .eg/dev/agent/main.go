// Package main runs a long lived workload for agent development.
// usage: eg compute serve --ttl infinity dev/agent
package main

import (
	"context"
	"eg/compute/daemon"
	"eg/compute/maintainer"
	"log"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/egenv"
	"github.com/egdaemon/eg/runtime/wasi/eggit"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LUTC | log.Ltime)
	ctx, done := context.WithTimeout(context.Background(), egenv.TTL())
	defer done()

	c1 := eg.Container(maintainer.Container)
	err := eg.Perform(
		ctx,
		eggit.AutoClone,
		eg.Build(
			c1.BuildFromFile(".eg/Containerfile"),
		),
		eg.Module(
			ctx,
			c1,
			daemon.Agent,
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}
