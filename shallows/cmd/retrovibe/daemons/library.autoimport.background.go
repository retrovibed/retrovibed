package daemons

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
)

// AutoImport periodically scans the monitored directories and uploads files into the default
// daemon's library once they've gone unmodified for their directory's debounce.
func AutoImport(ctx context.Context, q sqlx.Queryer, tlscfg *tls.Config, async *asyncx.Wakeup, enabled bool) error {
	if !enabled {
		log.Println("automatic library import is disabled")
		return nil
	}

	s := backoffx.New(
		backoffx.Constant(5*time.Minute),
		backoffx.Jitter(0.1),
	)

	go asyncx.Periodic(ctx, async, s, "automatic library import initiated - next")
	contextx.Run(ctx, func() {
		errorsx.Log(asyncx.Run(ctx, async, func(ctx context.Context) error {
			// a failed pass is retried on the next wakeup, returning the error would stop the loop.
			errorsx.Log(errorsx.Wrap(autoimport(ctx, q, tlscfg), "automatic library import failed"))
			return nil
		}))
	})

	return nil
}

// the default daemon is resolved every pass since it can be changed at any time.
func autoimport(ctx context.Context, q sqlx.Queryer, tlscfg *tls.Config) error {
	var (
		d meta.Daemon
	)

	if err := meta.DaemonFindDefault(ctx, q).Scan(&d); sqlx.ErrNoRows(err) != nil {
		log.Println("automatic library import skipped - no default daemon")
		return nil
	} else if err != nil {
		return errorsx.Wrap(err, "unable to lookup the default daemon")
	}

	endpoint := fmt.Sprintf("https://%s", d.Hostname)
	c := authn.AutoOauth2Client(ctx, tlscfg, authn.EndpointSSHAuth(endpoint))

	return library.NewAutoimport(ctx, q, media.NewUploader(c, endpoint))
}
