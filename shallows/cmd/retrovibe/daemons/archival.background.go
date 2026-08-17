package daemons

import (
	"context"
	"log"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"golang.org/x/crypto/ssh"
)

func AutoArchival(ctx context.Context, signer ssh.Signer, q sqlx.Queryer, mediastore fsx.Virtual, async *asyncx.Wakeup, archive bool) error {
	c, err := authn.AutoJWTClient(ctx, signer)
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
