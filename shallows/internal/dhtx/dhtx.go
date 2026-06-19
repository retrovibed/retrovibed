package dhtx

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/torrent/dht"
	"github.com/retrovibed/retrovibed/shallows/internal/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

func Statistics(ds *dht.Server) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		var b = bytes.NewBuffer(nil)
		if err := dht.Dump(ds, b); err != nil {
			return err
		}

		log.Println("dht", ds.ID(ds.DynamicAddrPort()), ds.AddrPort(ds.DynamicAddrPort()), spew.Sdump(ds.Stats()))
		log.Println(b.String())
		return nil
	}
}

// log out the dht statistics every d period.
func BackgroundStatistics(ctx context.Context, d time.Duration, ds *dht.Server) {
	timex.NowAndEvery(ctx, d, Statistics(ds))
}

func RecordBootstrapNodes(ctx context.Context, d time.Duration, min int, ds *dht.Server, dst string) {
	timex.Every(d, func() {
		current := ds.Nodes()
		if len(current) < min {
			log.Println("ignoring peers", len(current))
			return
		}

		log.Println("saving torrent peers", len(current))
		errorsx.Log(
			errorsx.Wrap(
				dht.WriteNodesToFile(current, dst),
				"unable to persist peers",
			),
		)
	})
}

func PopulateTable(ctx context.Context, d *dht.Server, path string) {
	if peers, err := d.AddNodesFromFile(path); err == nil {
		log.Printf("added peers %p - %d\n", d, peers)
	} else {
		log.Println("unable to read peers", err)
	}
}

func WaitForMinimumNodes(ctx context.Context, min int, dhts *dht.Server, do func()) {
	b := backoffx.New(backoffx.Exponential(time.Second), backoffx.Maximum(time.Minute))
	for attempts := 0; ; attempts++ {
		if dhts.NumNodes() > 32 {
			break
		}

		log.Printf("minimum nodes not available, waiting %p %d\n", dhts, dhts.NumNodes())
		time.Sleep(b.Backoff(attempts))
	}

	do()
}
