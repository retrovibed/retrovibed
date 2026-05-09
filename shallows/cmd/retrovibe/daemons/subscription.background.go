package daemons

import (
	"context"
	"net/http"
	"time"

	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/deeppool"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

func SubscriptionSync(ctx context.Context, q sqlx.Queryer, c *http.Client) error {
	published := deeppool.NewPublished(c)

	contextx.Run(ctx, func() {
		errorsx.Log(community.NewSubscriptionSync(ctx, q, published, 5*time.Minute))
	})

	return nil
}
