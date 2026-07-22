package cmdddisc

import (
	"encoding/hex"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

type cmdPeerCreate struct {
	Name      string `flag:"" name:"name" help:"name you wish to assign to this peer"`
	ID        string `flag:"" name:"peer" help:"hex encoded public peer id"`
	Partition string `flag:"" name:"partition" help:"partition this peer is responsible for" required:"true"`
}

func (t cmdPeerCreate) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(daemon.Endpoint), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, daemon.Endpoint)

	infohash, err := hex.DecodeString(t.ID)
	if err != nil {
		return errorsx.Wrap(err, "failed to decode peer id")
	}

	mrsp, err := ddiscapi.PeerCreate(gctx.Context, cc, daemon.Endpoint, &ddiscapi.PeerCreateRequest{
		Peer: &ddiscapi.Peer{
			Infohash:    infohash,
			Description: t.Name,
			Partition:   t.Partition,
			Ddisc:       true,
			Bep51:       true,
		},
	})
	if err != nil {
		return err
	}

	log.Println("peered with", spew.Sdump(mrsp.Peer))

	return nil
}
