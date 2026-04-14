package dhtx

import (
	"context"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/torrent/dht"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

// log out the dht statistics every d period.
func Statistics(ctx context.Context, d time.Duration, ds *dht.Server) {
	timex.NowAndEvery(ctx, d, func(ctx context.Context) error {
		log.Println("dht", ds.ID(), ds.AddrPort(), spew.Sdump(ds.Stats()))
		return nil
	})
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
