package cmdlibrary

import (
	"context"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/alecthomas/kong"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/media"
)

type cmdAutoimportCreate struct {
	Path             string        `arg:"" name:"path" help:"directory to monitor, only files directly within it are imported" required:"true"`
	Description      string        `flag:"" name:"description" help:"description of the directory" default:""`
	Debounce         time.Duration `flag:"" name:"debounce" help:"how long a file must go unmodified before it is imported" default:"1h"`
	Mode             string        `flag:"" name:"mode" help:"copy keeps the original file, move removes it after import" enum:"copy,move" default:"move"`
	LibraryDirectory string        `flag:"" name:"library-directory" help:"id of the library directory imported files are placed into, defaults to the library root" default:""`
}

func (t cmdAutoimportCreate) Run(kctx *kong.Context, gctx *cmdopts.Global, tls *cmdopts.TLSConfig, id *cmdopts.SSHID, daemon *cmdopts.Endpoint) (err error) {
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

func (t cmdAutoimportCreate) run(ctx context.Context, w io.Writer, endpoint string, c *http.Client) (err error) {
	path, err := filepath.Abs(t.Path)
	if err != nil {
		return errorsx.Wrapf(err, "unable to resolve path: %s", t.Path)
	}

	mode, err := autoimportModeFromString(t.Mode)
	if err != nil {
		return err
	}

	resp, err := media.AutoimportCreate(ctx, c, endpoint, &media.AutoimportDirectoryCreateRequest{
		Directory: &media.AutoimportDirectory{
			Path:               path,
			Description:        t.Description,
			Debounce:           uint64(t.Debounce / time.Second),
			Mode:               mode,
			LibraryDirectoryId: t.LibraryDirectory,
		},
	})
	if err != nil {
		return err
	}

	return printAutoimportDirectory(w, resp.Directory)
}
