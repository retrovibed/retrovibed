package daemons

import (
	"context"

	"golang.org/x/crypto/ssh"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
)

func AutoRegistration(ctx context.Context, signer ssh.Signer, account string, options ...httpx.ClientOption) {
	contextx.Run(ctx, func() {
		if _, err := authn.AutoRegistration(ctx, signer, account, options...); err != nil {
			errorsx.Log(errorsx.Wrap(err, "unable to register with archival service"))
		}
	})
}
