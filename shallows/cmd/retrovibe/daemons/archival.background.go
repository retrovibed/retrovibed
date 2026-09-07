package daemons

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/linxGnu/pqueue"
	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/shallows/backups"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/pqueuex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

func AutoArchival(ctx context.Context, q sqlx.Queryer, c *http.Client, mediastore fsx.Virtual, async *asyncx.Wakeup, archive bool) error {
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

func AutoBackup(ctx context.Context, db *sql.DB, c *http.Client, async *asyncx.Wakeup, wq pqueue.Queue, device string, enabled bool) error {
	if !enabled {
		log.Println("automatic backup is disabled")
		return nil
	}

	s := backoffx.New(
		backoffx.Constant(time.Hour),
		backoffx.Jitter(0.1),
	)

	go contextx.RunContext(ctx, pqueuex.NewWorker(wq, backups.NewWorker(c, db)).Consume)
	go asyncx.Periodic(ctx, async, s, "automatic backup initiated - next")
	contextx.Run(ctx, func() {
		errorsx.Log(asyncx.Run(ctx, async, func(ctx context.Context) error {
			if err := backups.Enqueue(ctx, wq, device); err != nil {
				log.Println(errorsx.Wrap(err, "backup request failed"))
			}

			return nil
		}))
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
