package cmdmeta

import (
	"context"
	"io"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
)

type CloudDisplay struct{}

func (t CloudDisplay) Run(gctx *cmdopts.Global, id *cmdopts.SSHID) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return err
	}

	return t.run(gctx.Context, os.Stdout, signer)
}

// run retrieves the session and permissions from the cloud service and prints them.
// options customize the underlying http clients, e.g. to point them at a test server.
func (t CloudDisplay) run(ctx context.Context, w io.Writer, id ssh.Signer, options ...httpx.ClientOption) error {
	session, err := authn.AutoRegistration(ctx, id, options...)
	if err != nil {
		return errorsx.Wrap(err, "unable to retrieve cloud session")
	}

	authz, err := authn.FetchCloudAuthz(ctx, id, options...)
	if err != nil {
		return errorsx.Wrap(err, "unable to retrieve cloud permissions")
	}

	if err = authn.PrintIdentity(w, id, session); err != nil {
		return err
	}

	return authn.PrintPermissions(w, authz.Token)
}
