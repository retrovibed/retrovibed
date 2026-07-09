package cmdddisc

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

type cmdPeerCreate struct {
	Endpoint  string `flag:"" name:"library" help:"http address for the library you want to connect to" default:"localhost:9998"`
	Name      string `flag:"" name:"name" help:"name you wish to assign to this peer"`
	ID        string `flag:"" name:"peer" help:"hex encoded public peer id"`
	Partition string `flag:"" name:"partition" help:"partition this peer is responsible for" required:"true"`
}

func (t cmdPeerCreate) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	infohash, err := hex.DecodeString(t.ID)
	if err != nil {
		return errorsx.Wrap(err, "failed to decode peer id")
	}

	mrsp, err := ddiscapi.PeerCreate(gctx.Context, cc, t.Endpoint, &ddiscapi.PeerCreateRequest{
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
