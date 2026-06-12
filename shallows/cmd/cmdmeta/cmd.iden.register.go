package cmdmeta

import (
	"context"
	"fmt"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
)

type IdenRegister struct {
	Endpoint string `flag:"" name:"endpoint" help:"http address of the retrovibed instance" default:"localhost:9998"`
}

func (t IdenRegister) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID) (err error) {
	signer, err := id.Signer()
	if err != nil {
		return err
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, t.Endpoint)

	return t.run(gctx.Context, cc)
}

func (t IdenRegister) run(ctx context.Context, c *http.Client) (err error) {
	return nil
}
