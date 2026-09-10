package daemons

import (
	"context"

	"golang.org/x/crypto/ssh"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

func AutoRegistration(ctx context.Context, signer ssh.Signer, options ...httpx.ClientOption) {
	contextx.Run(ctx, func() {
		if _, err := authn.AutoRegistration(ctx, signer, options...); err != nil {
			errorsx.Log(errorsx.Wrap(err, "unable to register with archival service"))
		}
	})
}
