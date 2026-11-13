package cmdmeta

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/httpx"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"golang.org/x/crypto/ssh"
)

type IdenDisplay struct{}

func (t IdenDisplay) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return err
	}

	// the local identity is printed first so it is available even when the endpoint is not.
	if err = authn.PrintIdentity(os.Stdout, signer, nil); err != nil {
		return err
	}

	// the session is optional, an unavailable endpoint is not a failure of displaying the identity.
	if err = t.session(gctx, tls, signer, daemon); errors.Is(err, syscall.ECONNREFUSED) {
		fmt.Fprintf(os.Stderr, "session unavailable: connection refused %s\n", daemon.Endpoint)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "session unavailable: %v\n", err)
	}

	return nil
}

func (t IdenDisplay) session(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, signer ssh.Signer, daemon *cmdopts.Endpoint) (err error) {
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

	if err = authn.PrintSession(os.Stdout, &authn.Session{Profile: current.Profile, Account: current.Account}); err != nil {
		return err
	}

	return authn.PrintPermissions(os.Stdout, authz.Token)
}
