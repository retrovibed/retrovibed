package cmdmeta

import (
	"fmt"
	"net/http"
	"os"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
)

type IdenDisplay struct{}

func (t IdenDisplay) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return err
	}

	c := authn.AutoOauth2Client(
		gctx.Context,
		tls.Config(),
		authn.EndpointSSHAuth(daemon.Endpoint),
		authn.SSHTokenSourceOptionSigner(signer),
	)

	req, err := http.NewRequestWithContext(gctx.Context, http.MethodGet, fmt.Sprintf("%s/sso/", daemon.Endpoint), nil)
	if err != nil {
		return errorsx.Wrap(err, "unable to create http request")
	}

	resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return errorsx.Wrap(err, "unable to retrieve local session")
	}
	defer resp.Body.Close()

	var current authn.Authn
	if err = jsonx.UnmarshalRead(resp.Body, &current); err != nil {
		return errorsx.Wrap(err, "unable to decode local session")
	}

	authz, err := authn.FetchAuthz(gctx.Context, c, fmt.Sprintf("%s/meta/authz/", daemon.Endpoint))
	if err != nil {
		return errorsx.Wrap(err, "unable to retrieve local permissions")
	}

	if err = authn.PrintIdentity(os.Stdout, signer, &authn.Session{Profile: current.Profile, Account: current.Account}); err != nil {
		return err
	}

	return authn.PrintPermissions(os.Stdout, authz.Token)
}
