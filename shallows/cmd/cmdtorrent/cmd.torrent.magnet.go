package cmdtorrent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/media"
)

type cmdMagnet struct {
	Endpoint string    `flag:"" name:"peer" help:"http address for the library you want to connect to" default:"localhost:9998"`
	Magnets  []url.URL `arg:"" name:"magnet" help:"magnet uri to download" required:"true"`
}

func (t cmdMagnet) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig) error {
	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)))

	for _, uri := range t.Magnets {
		var (
			err     error
			encoded []byte
			req     *http.Request
			resp    *http.Response
			magnet  media.MagnetCreateResponse
		)

		if encoded, err = json.Marshal(&media.MagnetCreateRequest{Uri: uri.String()}); err != nil {
			return errorsx.Wrap(err, "unable to encode magnet request")
		}

		if req, err = http.NewRequestWithContext(gctx.Context, http.MethodPost, fmt.Sprintf("https://%s/d/magnet", t.Endpoint), bytes.NewReader(encoded)); err != nil {
			return errorsx.Wrap(err, "unable to create http request")
		}

		if resp, err = httpx.AsError(c.Do(req)); err != nil {
			return errorsx.Wrap(err, "http request failed")
		}

		if err = httpx.DecodeJSON(resp, &magnet); err != nil {
			return errorsx.Wrap(err, "unable to decode response")
		}

		log.Println("downloading", magnet.Download.Media.Id, magnet.Download.Media.Description)
	}

	return nil
}
