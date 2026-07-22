package cmdddisc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

type cmdPeerList struct {
	Query string `flag:"" name:"query" help:"lucene text search (default field: description)" default:""`
}

func (t cmdPeerList) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(daemon.Endpoint), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, daemon.Endpoint)

	return t.run(gctx.Context, daemon.Endpoint, cc)
}

func (t cmdPeerList) run(ctx context.Context, endpoint string, c *http.Client) (err error) {
	result, err := ddiscapi.PeerSearch(ctx, c, endpoint, &ddiscapi.PeerSearchRequest{
		Query: t.Query,
		Limit: 100,
	})
	if err != nil {
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
