package daemons

import (
	"context"

	"golang.org/x/crypto/ssh"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

func AutoRegistration(ctx context.Context, signer ssh.Signer) {
	contextx.Run(ctx, func() {
		if _, err := authn.AutoRegistration(ctx, signer); err != nil {
			errorsx.Log(errorsx.Wrap(err, "unable to register with archival service"))
		}
	})
}
