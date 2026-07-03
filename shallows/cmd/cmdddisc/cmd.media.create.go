package cmdddisc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
)

type cmdMediaCreate struct {
	Endpoint     string `flag:"" name:"library" help:"http address for the library you want to connect to" default:"localhost:9998"`
	MagnetURI    string `flag:"" name:"magnet" help:"magnet uri of the media" required:"true"`
	Title        string `flag:"" name:"title" help:"title of the media"`
	Description  string `flag:"" name:"description" help:"description of the media"`
	Mimetype     string `flag:"" name:"mimetype" help:"mimetype of the media"`
	KnownMediaID string `flag:"" name:"known-media-id" help:"known media id to associate"`
	Partition    string `flag:"" name:"partition" help:"partition uuid this record belongs to"`
}

func (t cmdMediaCreate) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	var (
		encoded []byte
		req     *http.Request
		resp    *http.Response
		mrsp    ddiscapi.MediaCreateResponse
	)

	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	m, err := metainfo.ParseMagnetURI(t.MagnetURI)
	if err != nil {
		return errorsx.Wrap(err, "failed to parse magnet uri")
	}

	if encoded, err = json.Marshal(&ddiscapi.MediaCreateRequest{
		Media: &ddiscapi.Media{
			Infohash:     m.InfoHash.Bytes(),
			Title:        t.Title,
			Description:  t.Description,
			Mimetype:     t.Mimetype,
			KnownMediaId: t.KnownMediaID,
			Partition:    t.Partition,
		},
	}); err != nil {
		return errorsx.Wrap(err, "unable to encode request")
	}

	if req, err = http.NewRequestWithContext(gctx.Context, http.MethodPost, fmt.Sprintf("https://%s/ddisc/media/", t.Endpoint), bytes.NewReader(encoded)); err != nil {
		return errorsx.Wrap(err, "unable to create http request")
	}

	if resp, err = httpx.AsError(cc.Do(req)); err != nil {
		return errorsx.Wrap(err, "http request failed")
	}

	if err = httpx.DecodeJSON(resp, &mrsp); err != nil {
		return errorsx.Wrap(err, "unable to decode response")
	}

	log.Println("media created", spew.Sdump(mrsp.Media))

	return nil
}
