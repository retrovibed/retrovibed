package dhtx

import (
	"context"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/torrent/dht"
	"github.com/retrovibed/retrovibed/internal/x/errorsx"
	"github.com/retrovibed/retrovibed/internal/x/timex"
)

// log out the dht statistics every d period.
func Statistics(ctx context.Context, d time.Duration, ds *dht.Server) {
	timex.NowAndEvery(ctx, 10*time.Second, func(ctx context.Context) error {
		log.Println("dht", ds.ID(), spew.Sdump(ds.Stats()))
		return nil
	})
}

func RecordBootstrapNodes(ctx context.Context, d time.Duration, ds *dht.Server, dst string) {
	go timex.Every(time.Minute, func() {
		current := ds.Nodes()
		log.Println("saving torrent peers", len(current))
		errorsx.Log(
			errorsx.Wrap(
				dht.WriteNodesToFile(current, dst),
				"unable to persist peers",
			),
		)
	})
}
