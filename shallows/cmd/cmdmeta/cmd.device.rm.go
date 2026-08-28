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

type DeviceRm struct {
	DeviceID string `arg:"" name:"device.id" help:"device id to remove" required:"true"`
}

func (t DeviceRm) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
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

func (t DeviceRm) run(ctx context.Context, endpoint string, c *http.Client) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s/meta/d/%s", endpoint, t.DeviceID), nil)
	if err != nil {
		return errorsx.Wrap(err, "unable to create request")
	}

	resp, err := httpx.AsError(c.Do(req))
	if err != nil {
		return errorsx.Wrap(err, "remove device failed")
	}

	var result metaapi.DaemonDeleteResponse
	if err = httpx.DecodeJSON(resp, &result); err != nil {
		return err
	}

	_, err = fmt.Printf("id=%s hostname=%s removed\n", result.Daemon.Id, result.Daemon.Hostname)
	return err
}
