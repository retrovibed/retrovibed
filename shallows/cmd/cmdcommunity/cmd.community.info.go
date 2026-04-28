package cmdcommunity

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type cmdCommunityInfo struct {
	Description string        `flag:"" name:"description" help:"description of the community"`
	Name        string        `arg:"" name:"name" help:"name of the community globally unique. must be valid url subdomain" required:"true"`
	Output      cmdopts.IOOut `flag:"" name:"output" default:"-" help:"output destination; '-' for stdout"`
}

func (t cmdCommunityInfo) Run(gctx *cmdopts.Global, dpc cmdopts.DeeppoolClient) (err error) {
	out, err := t.Output.Open(os.Stdout)
	if err != nil {
		return errorsx.Wrap(err, "unable to open output")
	}

	c, err := dpc.HTTPClient(gctx.Context)
	if err != nil {
		return errorsx.Wrap(err, "unable to create deeppool client")
	}

	var in io.Reader = bytes.NewReader(nil)
	if cmdopts.Readable(os.Stdin) {
		in = os.Stdin
	}

	return t.run(gctx.Context, c, in, out)
}

func (t cmdCommunityInfo) run(ctx context.Context, c *http.Client, in io.Reader, out io.Writer) (err error) {
	commresp, err := metaapi.CommunityInfo(ctx, c, t.Name)
	if err != nil {
		return errorsx.Wrap(err, "failed to locate community")
	}

	if err = json.NewEncoder(out).Encode(commresp.Community); err != nil {
		return errorsx.Wrap(err, "unable to write to encoder")
	}

	if _, err = io.Copy(out, in); err != nil {
		return errorsx.Wrap(err, "failed to copy stdin to stdout")
	}

	return nil
}
