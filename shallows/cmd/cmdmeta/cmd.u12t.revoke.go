package cmdmeta

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type U12TRevoke struct {
	ProfileID string `arg:"" name:"profile.id" help:"profile id to revoke access" required:"true"`
}

func (t U12TRevoke) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
	ctx, done := context.WithTimeout(gctx.Context, 10*time.Second)
	defer done()

	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(daemon.Endpoint), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, daemon.Endpoint)

	return t.run(ctx, daemon.Endpoint, cc)
}

func (t U12TRevoke) run(ctx context.Context, endpoint string, c *http.Client) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/meta/authz/%s", endpoint, t.ProfileID), nil)
	if err != nil {
		return errorsx.Wrap(err, "unable to create request")
	}

	var revoke metaapi.AuthzRevokeResponse
	resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return errorsx.Wrap(err, "authz revoke failed")
	}

	return httpx.DecodeJSON(resp, &revoke)
}
