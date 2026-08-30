package daemons

import (
	"context"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/linxGnu/pqueue"
	"github.com/linxGnu/pqueue/entry"
	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"google.golang.org/protobuf/proto"
)

func SearchPluginRecommendationsRun(ctx context.Context, wq pqueue.Queue, q sqlx.Queryer, plugins searchplugin.T, peertube ddisc.DiscoverStrategy, mc library.QueryCleaner) error {
	log.Println("recommendations from search plugins initiated")
	defer log.Println("recommendations from search plugins completed")

	for ctx.Err() == nil {
		var (
			e      entry.Entry
			recreq media.RecommendationSearchRequest
		)

		if ok := wq.Dequeue(&e); !ok {
			return errorsx.New("failed to dequeue")
		}

		if err := proto.Unmarshal(e, &recreq); err != nil {
			return errorsx.Wrap(err, "derp derp")
		}

		log.Println("recommendation request", spew.Sdump(&recreq))
	}

	return nil
}

func SearchPluginRecommendationsBackground(ctx context.Context, wq pqueue.Queue, q sqlx.Queryer, seed string, importer tracking.URIImport, plugins searchplugin.T, peertube ddisc.DiscoverStrategy, mc library.QueryCleaner) error {
	wakeup := asyncx.NewWakeup(ctx)
	s := backoffx.New(
		backoffx.Frequency(24*time.Hour, seed),
		backoffx.JitterRandom(time.Minute),
	)

	go asyncx.Periodic(ctx, wakeup, s, "ddisc search queue drain")
	contextx.Run(ctx, func() {
		errorsx.Log(asyncx.Run(ctx, wakeup, func(ctx context.Context) error {
			return nil
		}))
	})

	return nil
}
