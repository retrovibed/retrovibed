package cmdddisc

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/ddiscapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

type cmdMediaDelete struct {
	ID string `flag:"" name:"id" help:"discovered media record id" required:"true"`
}

func (t cmdMediaDelete) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(daemon.Endpoint), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, daemon.Endpoint)

	mrsp, err := ddiscapi.MediaDelete(gctx.Context, cc, daemon.Endpoint, t.ID)
	if err != nil {
		return err
	}

	log.Println("media removed", spew.Sdump(mrsp.Media))

	return nil
}
