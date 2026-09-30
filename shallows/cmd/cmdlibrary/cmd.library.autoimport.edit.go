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

type cmdAutoimportEdit struct {
	ID               string         `arg:"" name:"id" help:"id of the monitored directory to edit" required:"true"`
	Description      *string        `flag:"" name:"description" help:"update the description"`
	Debounce         *time.Duration `flag:"" name:"debounce" help:"update how long a file must go unmodified before it is imported"`
	Mode             *string        `flag:"" name:"mode" help:"update the mode: copy keeps the original file, move removes it after import"`
	LibraryDirectory *string        `flag:"" name:"library-directory" help:"update the id of the library directory imported files are placed into"`
}

func (t cmdAutoimportEdit) Run(kctx *kong.Context, gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
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

// run replaces the full set of editable fields, so the existing record is fetched first and
// carried forward for whichever fields weren't overridden.
func (t cmdAutoimportEdit) run(ctx context.Context, w io.Writer, endpoint string, c *http.Client) (err error) {
	existing, err := media.AutoimportFind(ctx, c, endpoint, t.ID)
	if err != nil {
		return err
	}

	patched := existing.Directory

	if t.Description != nil {
		patched.Description = *t.Description
	}

	if t.Debounce != nil {
		patched.Debounce = uint64(*t.Debounce / time.Second)
	}

	if t.Mode != nil {
		if patched.Mode, err = autoimportModeFromString(*t.Mode); err != nil {
			return err
		}
	}

	if t.LibraryDirectory != nil {
		patched.LibraryDirectoryId = *t.LibraryDirectory
	}

	resp, err := media.AutoimportUpdate(ctx, c, endpoint, t.ID, &media.AutoimportDirectoryUpdateRequest{
		Directory: patched,
	})
	if err != nil {
		return err
	}

	return printAutoimportDirectory(w, resp.Directory)
}
