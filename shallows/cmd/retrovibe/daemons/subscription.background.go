package daemons

import (
	"context"
	"net/http"
	"time"

	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

func SubscriptionSync(ctx context.Context, q sqlx.Queryer, c *http.Client, async *asyncx.Wakeup) error {
	published := communityapi.NewDeeppoolPublished(c)

	contextx.Run(ctx, func() {
		errorsx.Log(communityapi.NewSubscriptionSync(ctx, q, published, async, 5*time.Minute))
	})

	return nil
}
