package cmdopts

import (
	"context"
	"log"
	"runtime/debug"
	"sync"
)

// Global command fields.
type Global struct {
	Context  context.Context    `kong:"-"`
	Shutdown context.CancelFunc `kong:"-"`
	Cleanup  *sync.WaitGroup    `kong:"-"`
}

type Version struct{}

func (t Version) Run(ctx *Global) (err error) {
	var (
		ok   bool
		info *debug.BuildInfo
	)

	if info, ok = debug.ReadBuildInfo(); ok {
		log.Println(info.Main.Path, info.Main.Version)
		return nil
	}

	log.Println("unknown version")
	return nil
}
