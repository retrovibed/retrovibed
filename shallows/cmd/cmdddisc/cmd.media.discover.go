package cmdddisc

import (
	"fmt"
	"log"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// cmdMediaDiscover submits a locate request to a running daemon over its
// HTTP API - it does not run discovery itself.
type cmdMediaDiscover struct {
	Endpoint     string `flag:"" name:"library" help:"http address for the library you want to connect to" default:"localhost:9998"`
	Query        string `arg:"" name:"query" help:"title or free-text search to locate media for"`
	Mimetype     string `flag:"" name:"mimetype" help:"mimetype/category (video, audio, image, text, application) to search for" required:""`
	KnownMediaID string `flag:"" name:"known-media-id" help:"known media id, if already resolved from a catalog search" optional:""`
}

func (t cmdMediaDiscover) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	resp, err := ddiscapi.LocateCreate(gctx.Context, cc, t.Endpoint, &ddiscapi.LocateCreateRequest{
		Locate: &ddiscapi.Locate{
			Query:        t.Query,
			Mimetype:     t.Mimetype,
			KnownMediaId: t.KnownMediaID,
		},
	})
	if err != nil {
		return errorsx.Wrap(err, "unable to submit locate request")
	}

	log.Println("locate request submitted", resp.Locate.Id, resp.Locate.Query)
	return nil
}
