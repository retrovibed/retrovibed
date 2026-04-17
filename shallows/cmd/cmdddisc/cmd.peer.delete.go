package cmdddisc

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type cmdPeerDelete struct {
	Endpoint string `flag:"" name:"library" help:"http address for the library you want to connect to" default:"localhost:9998"`
	ID       string `flag:"" name:"peer" help:"hex encoded public peer id"`
}

func (t cmdPeerDelete) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig) (err error) {
	var (
		req  *http.Request
		resp *http.Response
		mrsp ddiscapi.PeerDeleteResponse
	)

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	infohash, err := int160.FromHexEncodedString(t.ID)
	if err != nil {
		return errorsx.Wrap(err, "failed to decode peer id")
	}

	if req, err = http.NewRequestWithContext(gctx.Context, http.MethodDelete, fmt.Sprintf("https://%s/ddisc/%s", t.Endpoint, tracking.PeerUID(infohash)), bytes.NewReader(nil)); err != nil {
		return errorsx.Wrap(err, "unable to create http request")
	}

	if resp, err = httpx.AsError(cc.Do(req)); err != nil {
		return errorsx.Wrap(err, "http request failed")
	}

	if err = httpx.DecodeJSON(resp, &mrsp); err != nil {
		return errorsx.Wrap(err, "unable to decode response")
	}

	log.Println("peer removed", spew.Sdump(mrsp.Peer))

	return nil
}
