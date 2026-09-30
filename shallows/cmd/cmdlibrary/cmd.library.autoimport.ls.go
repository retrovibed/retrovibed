package cmdlibrary

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/alecthomas/kong"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/media"
)

type cmdAutoimportLs struct {
	ID     []string `flag:"" name:"id" help:"only show directories matching the given id(s)"`
	Query  string   `flag:"" name:"query" help:"only show directories whose description matches this text (lucene syntax, e.g. path:inbox)"`
	Offset uint64   `flag:"" name:"offset" help:"page offset for pagination (multiplied by the result limit)"`
}

func (t cmdAutoimportLs) Run(kctx *kong.Context, gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
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

func (t cmdAutoimportLs) run(ctx context.Context, w io.Writer, endpoint string, c *http.Client) (err error) {
	result, err := media.AutoimportSearch(ctx, c, endpoint, &media.AutoimportDirectorySearchRequest{
		Id:     t.ID,
		Query:  t.Query,
		Offset: t.Offset,
		Limit:  100,
	})
	if err != nil {
		return err
	}

	for _, d := range result.Items {
		if err = printAutoimportDirectory(w, d); err != nil {
			return err
		}
	}

	return nil
}
