package cmdlibrary

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alecthomas/kong"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/media"
)

type cmdAutoimportShow struct {
	ID string `arg:"" name:"id" help:"id of the monitored directory to show" required:"true"`
}

func (t cmdAutoimportShow) Run(kctx *kong.Context, gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
	ctx, done := context.WithTimeout(gctx.Context, 10*time.Second)
	defer done()

	signer, err := id.Signer()
	if err != nil {
		return errorsx.Wrap(err, "failed to create signer")
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(daemon.Endpoint), authn.SSHTokenSourceOptionSigner(signer))
	cc := authn.AuthzClientLibrary(tls.Config(), c, daemon.Endpoint)

	return t.run(ctx, kctx.Stdout, daemon.Endpoint, cc)
}

func (t cmdAutoimportShow) run(ctx context.Context, w io.Writer, endpoint string, c *http.Client) (err error) {
	resp, err := media.AutoimportFind(ctx, c, endpoint, t.ID)
	if err != nil {
		return err
	}

	if err = printAutoimportDirectory(w, resp.Directory); err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "pending=%d imported=%d\n", resp.Directory.Pending, resp.Directory.Imported)
	return err
}
