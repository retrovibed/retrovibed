package cmdddisc

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

type cmdDiscoveryDelete struct {
	Endpoint string `flag:"" name:"library" help:"http address for the library you want to connect to" default:"localhost:9998"`
	ID       string `flag:"" name:"id" help:"id of the discovery entry to remove" required:"true"`
}

func (t cmdDiscoveryDelete) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	mrsp, err := ddiscapi.DiscoveryDelete(gctx.Context, cc, t.Endpoint, t.ID)
	if err != nil {
		return err
	}

	log.Println("discovery entry removed", spew.Sdump(mrsp.Discovery))

	return nil
}
