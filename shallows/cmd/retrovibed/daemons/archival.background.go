package daemons

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/retrovibed/retrovibed/community"
	"github.com/retrovibed/retrovibed/deeppool"
	"github.com/retrovibed/retrovibed/internal/asyncx"
	"github.com/retrovibed/retrovibed/internal/backoffx"
	"github.com/retrovibed/retrovibed/internal/contextx"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/library"
	"github.com/retrovibed/retrovibed/metaapi"
)

func AutoArchival(ctx context.Context, q sqlx.Queryer, mediastore fsx.Virtual, async *asyncx.Wakeup, archive bool) error {
	if archive {
		if _, err := metaapi.Register(ctx); err != nil {
			return errorsx.Wrap(err, "unable to register with archival service")
		}
	}

	c, err := metaapi.AutoJWTClient(ctx)
	if err != nil {
		return errorsx.Wrap(err, "failed to create oauth2 bearer token")
	}

	s := backoffx.New(
		backoffx.Constant(time.Hour),
		backoffx.Jitter(0.1),
	)

	go asyncx.Periodic(ctx, async, s, "automatic archival initiated - next")
	contextx.Run(ctx, func() {
		errorsx.Log(library.NewAutoArchive(ctx, c, mediastore, q, async, archive))
	})

	return nil
}

func PendingSync(ctx context.Context, q sqlx.Queryer, c *http.Client, mvfs, tvfs fsx.Virtual) error {
	metrics := deeppool.NewMetrics(c)
	published := deeppool.NewPublished(c)

	contextx.Run(ctx, func() {
		errorsx.Log(community.NewPendingSync(ctx, q, c, metrics, published, mvfs, tvfs, time.Minute))
	})

	return nil
}

func AutoReclaim(ctx context.Context, q sqlx.Queryer, mediastore fsx.Virtual, async *asyncx.Wakeup, reclaimdisk bool) error {
	s := backoffx.New(
		backoffx.Constant(time.Hour),
		backoffx.Jitter(0.1),
	)

	if !reclaimdisk {
		log.Println("automatic disk reclaim is disabled - enabling dry-run")
	}

	go asyncx.Periodic(ctx, async, s, "automatic disk reclaim initiated - next")
	contextx.Run(ctx, func() {
		errorsx.Log(library.NewSlowDiskReclaim(ctx, mediastore, q, async, 80, reclaimdisk))
	})
	return nil
}
