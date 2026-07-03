package cmdddisc

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

type cmdDiscoveryList struct {
	Endpoint   string   `flag:"" name:"library" help:"http address for the library you want to connect to" default:"localhost:9998"`
	NeedsCheck bool     `flag:"" name:"needs-check" help:"only show entries due for another check"`
	ID         []string `flag:"" name:"id" help:"only show entries matching the given id(s)"`
}

func (t cmdDiscoveryList) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	var (
		encoded url.Values
		req     *http.Request
		resp    *http.Response
		result  ddiscapi.DiscoverySearchResponse
	)

	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	if encoded, err = formx.NewEncoder().Encode(&ddiscapi.DiscoverySearchRequest{
		NeedsCheck: t.NeedsCheck,
		Id:         t.ID,
		Limit:      100,
	}); err != nil {
		return errorsx.Wrap(err, "unable to encode request")
	}

	if req, err = http.NewRequestWithContext(gctx.Context, http.MethodGet, fmt.Sprintf("https://%s/ddisc/discovery/?"+encoded.Encode(), t.Endpoint), nil); err != nil {
		return errorsx.Wrap(err, "unable to create http request")
	}

	if resp, err = httpx.AsError(cc.Do(req)); err != nil {
		return errorsx.Wrap(err, "http request failed")
	}

	if err = httpx.DecodeJSON(resp, &result); err != nil {
		return errorsx.Wrap(err, "unable to decode response")
	}

	for _, d := range result.Items {
		if _, err = fmt.Printf("id=%s infohash=%x attempts=%s next_check=%s\n", d.GetId(), d.GetInfohash(), d.GetAttempts(), d.GetNextCheck()); err != nil {
			return err
		}
	}

	return nil
}
