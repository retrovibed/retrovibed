package daemons

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

func AutoPublishing(ctx context.Context, q sqlx.Queryer, c *http.Client, mvfs, tvfs fsx.Virtual, async *asyncx.Wakeup) error {
	metrics := communityapi.NewMetrics(c)
	published := communityapi.NewDeeppoolCommunity(c)

	s := backoffx.New(
		backoffx.Constant(time.Hour),
		backoffx.Jitter(0.1),
	)

	go asyncx.Periodic(ctx, async, s, "automatic publishing initiated")
	contextx.Run(ctx, func() {
		errorsx.Log(asyncx.Run(ctx, async, func(ctx context.Context) error {
			if err := communityapi.SyncPendingToDeeppool(ctx, q, c, metrics, published, deeppool.NewArchiver(c), mvfs, tvfs); err != nil {
				log.Println(errorsx.Wrap(err, "publishing sync failed"))
				return nil
			}

			return nil
		}))
	})

	return nil
}

func AutoFeedSync(ctx context.Context, q sqlx.Queryer, c *http.Client, async *asyncx.Wakeup) error {
	published := communityapi.NewDeeppoolCommunity(c)

	contextx.Run(ctx, func() {
		errorsx.Log(asyncx.Run(ctx, async, func(ctx context.Context) error {
			if err := communityapi.SyncFeed(ctx, q, c, published); err != nil {
				log.Println(errorsx.Wrap(err, "feed sync failed"))
			}
			return nil
		}))
	})

	return nil
}
