package cmdddisc

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

type cmdMediaCreate struct {
	MagnetURI    string `flag:"" name:"magnet" help:"magnet uri of the media" required:"true"`
	Title        string `flag:"" name:"title" help:"title of the media"`
	Description  string `flag:"" name:"description" help:"description of the media"`
	Mimetype     string `flag:"" name:"mimetype" help:"mimetype of the media"`
	KnownMediaID string `flag:"" name:"known-media-id" help:"known media id to associate"`
	Partition    string `flag:"" name:"partition" help:"partition uuid this record belongs to"`
}

func (t cmdMediaCreate) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(daemon.Endpoint), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, daemon.Endpoint)

	m, err := metainfo.ParseMagnetURI(t.MagnetURI)
	if err != nil {
		return errorsx.Wrap(err, "failed to parse magnet uri")
	}

	mrsp, err := ddiscapi.MediaCreate(gctx.Context, cc, daemon.Endpoint, &ddiscapi.MediaCreateRequest{
		Media: &ddiscapi.Media{
			Infohash:     m.InfoHash.Bytes(),
			Title:        t.Title,
			Description:  t.Description,
			Mimetype:     t.Mimetype,
			KnownMediaId: t.KnownMediaID,
			Partition:    t.Partition,
		},
	})
	if err != nil {
		return err
	}

	log.Println("media created", spew.Sdump(mrsp.Media))

	return nil
}
