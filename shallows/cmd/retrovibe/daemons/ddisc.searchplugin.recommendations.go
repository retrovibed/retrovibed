package daemons

import (
	"context"
	"log"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

func SearchPluginRecommendationsRun(ctx context.Context, q sqlx.Queryer, importer tracking.URIImport, plugins searchplugin.T, peertube ddisc.DiscoverStrategy, mc library.QueryCleaner) error {
	log.Println("recommendations from search plugins initiated")
	defer log.Println("recommendations from search plugins completed")

	return nil
}

func SearchPluginRecommendationsBackground(ctx context.Context, q sqlx.Queryer, seed string, importer tracking.URIImport, plugins searchplugin.T, peertube ddisc.DiscoverStrategy, mc library.QueryCleaner) error {
	wakeup := asyncx.NewWakeup(ctx)
	defer wakeup.Broadcast() // kick off an initial drain
	s := backoffx.New(
		backoffx.Exponential(time.Second),
		backoffx.Maximum(time.Hour),
		backoffx.Jitter(0.1),
	)

	// timex.StartOfDay(time.Now())
	// backoffx.DynamicHashWindow()
	go asyncx.Periodic(ctx, wakeup, s, "ddisc search queue drain")
	contextx.Run(ctx, func() {
		errorsx.Log(asyncx.Run(ctx, wakeup, func(ctx context.Context) error {
			return SearchQueueBackgroundRun(ctx, q, importer, plugins, peertube, mc)
		}))
	})

	return nil
}
