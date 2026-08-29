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
)

func SearchPluginRecommendationsRun(ctx context.Context, q sqlx.Queryer, plugins searchplugin.T, peertube ddisc.DiscoverStrategy, mc library.QueryCleaner) error {
	log.Println("recommendations from search plugins initiated")
	defer log.Println("recommendations from search plugins completed")

	return nil
}

func SearchPluginRecommendationsBackground(ctx context.Context, q sqlx.Queryer, seed string, wakeup *asyncx.Wakeup, plugins searchplugin.T, peertube ddisc.DiscoverStrategy, mc library.QueryCleaner) error {
	s := backoffx.New(
		backoffx.Frequency(24*time.Hour, seed),
		backoffx.JitterRandom(time.Minute),
	)

	go asyncx.Periodic(ctx, wakeup, s, "media recommendations")
	contextx.Run(ctx, func() {
		errorsx.Log(asyncx.Run(ctx, wakeup, func(ctx context.Context) error {
			return SearchPluginRecommendationsRun(ctx, q, plugins, peertube, mc)
		}))
	})

	return nil
}
