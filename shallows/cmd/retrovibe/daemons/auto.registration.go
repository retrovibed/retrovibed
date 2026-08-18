package daemons

import (
	"context"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

func AutoRegistration(ctx context.Context, c *http.Client) {
	contextx.Run(ctx, func() {
		if _, err := authn.Register(ctx, c); err != nil {
			errorsx.Log(errorsx.Wrap(err, "unable to register with archival service"))
		}
	})
}
