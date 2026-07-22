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

type cmdDiscoveryCreate struct {
	MagnetURI string `flag:"" name:"magnet" help:"magnet uri to start tracking (e.g. magnet:?xt=urn:btih:...)" required:"true"`
}

func (t cmdDiscoveryCreate) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
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

	mrsp, err := ddiscapi.DiscoveryCreate(gctx.Context, cc, daemon.Endpoint, &ddiscapi.DiscoveryCreateRequest{
		Discovery: &ddiscapi.Discovery{
			Infohash: m.InfoHash.Bytes(),
		},
	})
	if err != nil {
		return err
	}

	log.Println("discovery entry created", spew.Sdump(mrsp.Discovery))

	return nil
}
