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

type cmdMediaLs struct {
	Endpoint     string `flag:"" name:"library" help:"http address for the library you want to query" default:"localhost:9998"`
	Query        string `flag:"" name:"query" help:"lucene text search (default field: description)" default:""`
	KnownMediaID string `flag:"" name:"known-media-id" help:"filter by known media id"`
	NeedsCheck   bool   `flag:"" name:"needs-check" help:"only include media due for a recheck"`
}

func (t cmdMediaLs) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	var (
		encoded url.Values
		req     *http.Request
		resp    *http.Response
		result  ddiscapi.MediaSearchResponse
	)

	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	if encoded, err = formx.NewEncoder().Encode(&ddiscapi.MediaSearchRequest{
		Query:        t.Query,
		KnownMediaId: t.KnownMediaID,
		NeedsCheck:   t.NeedsCheck,
		Limit:        100,
	}); err != nil {
		return errorsx.Wrap(err, "unable to encode request")
	}

	if req, err = http.NewRequestWithContext(gctx.Context, http.MethodGet, fmt.Sprintf("https://%s/ddisc/media/?"+encoded.Encode(), t.Endpoint), nil); err != nil {
		return errorsx.Wrap(err, "unable to create http request")
	}

	if resp, err = httpx.AsError(cc.Do(req)); err != nil {
		return errorsx.Wrap(err, "http request failed")
	}

	if err = httpx.DecodeJSON(resp, &result); err != nil {
		return errorsx.Wrap(err, "unable to decode response")
	}

	for _, m := range result.Items {
		if _, err = fmt.Printf("id=%s title=%s known_media_id=%s mimetype=%s infohash=%x\n", m.GetId(), m.GetTitle(), m.GetKnownMediaId(), m.GetMimetype(), m.GetInfohash()); err != nil {
			return err
		}
	}

	return nil
}
