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
	Endpoint  string `flag:"" name:"endpoint" help:"http address of the retrovibed instance" default:"localhost:9998"`
	ProfileID string `arg:"" name:"profile.id" help:"profile id to revoke access" required:"true"`
}

func (t U12TRevoke) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	ctx, done := context.WithTimeout(gctx.Context, 10*time.Second)
	defer done()

	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	return t.run(ctx, cc)
}

func (t U12TRevoke) run(ctx context.Context, c *http.Client) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("https://%s/meta/authz/%s", t.Endpoint, t.ProfileID), nil)
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
