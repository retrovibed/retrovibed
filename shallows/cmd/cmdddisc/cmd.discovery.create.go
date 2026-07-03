package cmdddisc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

type cmdDiscoveryCreate struct {
	Endpoint string `flag:"" name:"library" help:"http address for the library you want to connect to" default:"localhost:9998"`
	Infohash string `flag:"" name:"infohash" help:"hex encoded infohash to start tracking" required:"true"`
}

func (t cmdDiscoveryCreate) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	var (
		encoded []byte
		req     *http.Request
		resp    *http.Response
		mrsp    ddiscapi.DiscoveryCreateResponse
	)

	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	infohash, err := hex.DecodeString(t.Infohash)
	if err != nil {
		return errorsx.Wrap(err, "failed to decode infohash")
	}

	if encoded, err = json.Marshal(&ddiscapi.DiscoveryCreateRequest{
		Discovery: &ddiscapi.Discovery{
			Infohash: infohash,
		},
	}); err != nil {
		return errorsx.Wrap(err, "unable to encode request")
	}

	if req, err = http.NewRequestWithContext(gctx.Context, http.MethodPost, fmt.Sprintf("https://%s/ddisc/discovery/", t.Endpoint), bytes.NewReader(encoded)); err != nil {
		return errorsx.Wrap(err, "unable to create http request")
	}

	if resp, err = httpx.AsError(cc.Do(req)); err != nil {
		return errorsx.Wrap(err, "http request failed")
	}

	if err = httpx.DecodeJSON(resp, &mrsp); err != nil {
		return errorsx.Wrap(err, "unable to decode response")
	}

	log.Println("discovery entry created", spew.Sdump(mrsp.Discovery))

	return nil
}
