package cmdmeta

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alecthomas/kong"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type DeviceLs struct {
	// Query is accepted for parity with metaapi.DaemonSearchRequest, but the
	// server currently decodes it without applying it as a filter (see
	// HTTPDaemons.search) — this is effectively a paginated list-all today.
	Query string `flag:"" name:"query" help:"text search (currently unfiltered server-side)" default:""`
}

func (t DeviceLs) Run(kctx *kong.Context, gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
	ctx, done := context.WithTimeout(gctx.Context, 10*time.Second)
	defer done()

	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(daemon.Endpoint), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, daemon.Endpoint)

	return t.run(ctx, kctx.Stdout, daemon.Endpoint, cc)
}

func (t DeviceLs) run(ctx context.Context, w io.Writer, endpoint string, c *http.Client) (err error) {
	req := metaapi.DaemonSearchRequest{
		Query: t.Query,
		Limit: 128,
	}

	encoded, err := formx.NewEncoder().Encode(&req)
	if err != nil {
		return errorsx.Wrap(err, "unable to encode request")
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/meta/d/?"+encoded.Encode(), endpoint), nil)
	if err != nil {
		return errorsx.Wrap(err, "unable to create request")
	}

	resp, err := httpx.AsError(c.Do(hreq))
	if err != nil {
		return errorsx.Wrap(err, "request failed")
	}

	var result metaapi.DaemonSearchResponse
	if err = httpx.DecodeJSON(resp, &result); err != nil {
		return err
	}

	for _, d := range result.Items {
		if _, err = fmt.Fprintf(w, "id='%s' hostname='%s' description='%s' default=%t downloads=%t\n", d.Id, d.Hostname, d.Description, d.Default, d.Downloads); err != nil {
			return err
		}
	}

	return nil
}
