package cmdddisc

import (
	"bytes"
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

type cmdMediaDelete struct {
	Endpoint string `flag:"" name:"library" help:"http address for the library you want to connect to" default:"localhost:9998"`
	ID       string `flag:"" name:"id" help:"discovered media record id" required:"true"`
}

func (t cmdMediaDelete) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	var (
		req  *http.Request
		resp *http.Response
		mrsp ddiscapi.MediaDeleteResponse
	)

	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	if req, err = http.NewRequestWithContext(gctx.Context, http.MethodDelete, fmt.Sprintf("https://%s/ddisc/media/%s", t.Endpoint, t.ID), bytes.NewReader(nil)); err != nil {
		return errorsx.Wrap(err, "unable to create http request")
	}

	if resp, err = httpx.AsError(cc.Do(req)); err != nil {
		return errorsx.Wrap(err, "http request failed")
	}

	if err = httpx.DecodeJSON(resp, &mrsp); err != nil {
		return errorsx.Wrap(err, "unable to decode response")
	}

	log.Println("media removed", spew.Sdump(mrsp.Media))

	return nil
}
