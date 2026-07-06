package cmdddisc

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/formx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/meta"
)

type cmdMediaLs struct {
	Endpoint     string        `flag:"" name:"library" help:"http address for the library you want to query" default:"localhost:9998"`
	Query        string        `flag:"" name:"query" help:"lucene text search (default field: description)" default:""`
	KnownMediaID string        `flag:"" name:"known-media-id" help:"filter by known media id"`
	NextCheck    time.Duration `flag:"" name:"next-check" help:"only show entries due within the the next window" default:"30m"`
	ID           []string      `flag:"" name:"id" help:"only show entries matching the given id(s)"`
	Offload      bool          `flag:"" name:"offload" help:"only show media marked to be offloaded / not indexed by this node"`
	Indexing     bool          `flag:"" name:"indexing" help:"only show media still pending identification"`
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

	knownMediaID := t.KnownMediaID
	switch {
	case t.Offload:
		knownMediaID = uuid.Nil.String()
	case t.Indexing:
		knownMediaID = uuid.Max.String()
	}

	if encoded, err = formx.NewEncoder().Encode(&ddiscapi.MediaSearchRequest{
		Query:        t.Query,
		KnownMediaId: knownMediaID,
		Id:           t.ID,
		NextCheck:    meta.NewDateRange(timex.NewRangeWithin(t.NextCheck)),
		Limit:        100,
	}); err != nil {
		return errorsx.Wrap(err, "unable to encode request")
	}

	log.Println("DERP DERP", encoded.Encode())

	if req, err = http.NewRequestWithContext(gctx.Context, http.MethodGet, fmt.Sprintf("https://%s/ddisc/media/?%s", t.Endpoint, encoded.Encode()), nil); err != nil {
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
