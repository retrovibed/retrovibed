package cmdddisc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type cmdPeerCreate struct {
	Endpoint  string `flag:"" name:"library" help:"http address for the library you want to connect to" default:"localhost:9998"`
	Name      string `flag:"" name:"name" help:"name you wish to assign to this peer"`
	ID        string `flag:"" name:"peer" help:"hex encoded public peer id"`
	Partition string `flag:"" name:"partition" help:"partition this peer is responsible for" required:"true"`
}

func (t cmdPeerCreate) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig) (err error) {
	var (
		encoded []byte
		req     *http.Request
		resp    *http.Response
		mrsp    ddiscapi.PeerCreateResponse
	)

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)))
	cc := metaapi.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	infohash, err := hex.DecodeString(t.ID)
	if err != nil {
		return errorsx.Wrap(err, "failed to decode peer id")
	}

	if encoded, err = json.Marshal(&ddiscapi.PeerCreateRequest{
		Peer: &ddiscapi.Peer{
			Infohash:    infohash,
			Description: t.Name,
			Partition:   t.Partition,
			Ddisc:       true,
			Bep51:       true,
		},
	}); err != nil {
		return errorsx.Wrap(err, "unable to encode magnet request")
	}

	if req, err = http.NewRequestWithContext(gctx.Context, http.MethodPost, fmt.Sprintf("https://%s/ddisc/", t.Endpoint), bytes.NewReader(encoded)); err != nil {
		return errorsx.Wrap(err, "unable to create http request")
	}

	if resp, err = httpx.AsError(cc.Do(req)); err != nil {
		return errorsx.Wrap(err, "http request failed")
	}

	if err = httpx.DecodeJSON(resp, &mrsp); err != nil {
		return errorsx.Wrap(err, "unable to decode response")
	}

	log.Println("peered with", spew.Sdump(mrsp.Peer))

	return nil
}
