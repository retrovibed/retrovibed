package cmdmeta

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
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

	return t.run(gctx.Context, c)
}

func (t IdenRegister) run(ctx context.Context, c *http.Client) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://%s/sso/register", t.Endpoint), nil)
	if err != nil {
		return errorsx.Wrap(err, "unable to create http request")
	}

	resp, err := httpx.AsError(c.Do(req))
	if cause, ok := errors.AsType[*httpx.Error](err); ok && cause.Code == http.StatusConflict {
		return nil
	} else if err != nil {
		return errorsx.Wrap(err, "registration failed")
	}
	defer func() {
		err = langx.FirstNonNil(err, resp.Body.Close())
	}()

	log.Println("registration request successful")

	return nil
}
