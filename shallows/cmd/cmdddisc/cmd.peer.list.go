package cmdddisc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

type cmdPeerList struct {
	Endpoint string `flag:"" name:"library" help:"http address for the library you want to connect to" default:"localhost:9998"`
	Query    string `flag:"" name:"query" help:"lucene text search (default field: description)" default:""`
}

func (t cmdPeerList) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	return t.run(gctx.Context, cc)
}

func (t cmdPeerList) run(ctx context.Context, c *http.Client) (err error) {
	req := ddiscapi.PeerSearchRequest{
		Query: t.Query,
		Limit: 100,
	}

	encoded, err := formx.NewEncoder().Encode(&req)
	if err != nil {
		return errorsx.Wrap(err, "unable to encode request")
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/ddisc/?"+encoded.Encode(), t.Endpoint), nil)
	if err != nil {
		return errorsx.Wrap(err, "unable to create request")
	}

	resp, err := httpx.AsError(c.Do(hreq))
	if err != nil {
		return errorsx.Wrap(err, "request failed")
	}

	var result ddiscapi.PeerSearchResponse
	if err = httpx.DecodeJSON(resp, &result); err != nil {
		return err
	}

	for _, p := range result.Items {
		if _, err = fmt.Println(peerLine(p)); err != nil {
			return err
		}
	}

	return nil
}

func peerLine(p *ddiscapi.Peer) string {
	return fmt.Sprintf("id=%s description=%s partition=%s ddisc=%t bep51=%t", p.GetId(), p.GetDescription(), p.GetPartition(), p.GetDdisc(), p.GetBep51())
}
