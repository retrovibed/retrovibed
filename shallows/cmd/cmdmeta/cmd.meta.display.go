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

	return t.run(gctx.Context, os.Stdout, signer, gctx.AccountID, cmdopts.ChooseProfile)
}

// run retrieves the session and permissions from the cloud service and prints them. account selects the profile
// when the key is associated with multiple, otherwise choose is asked to select one.
// options customize the underlying http clients, e.g. to point them at a test server.
func (t CloudDisplay) run(ctx context.Context, w io.Writer, id ssh.Signer, account string, choose authn.Chooser, options ...httpx.ClientOption) error {
	session, err := authn.AutoRegistrationOnce(ctx, id, account, false, choose, options...)
	if err != nil {
		return errorsx.Wrap(err, "unable to retrieve cloud session")
	}

	// act as the account the registration resolved.
	authz, err := authn.FetchCloudAuthz(ctx, id, session.Account.GetId(), options...)
	if err != nil {
		return errorsx.Wrap(err, "unable to retrieve cloud permissions")
	}

	if err = authn.PrintIdentity(w, id, session); err != nil {
		return err
	}

	return authn.PrintPermissions(w, authz.Token)
}
