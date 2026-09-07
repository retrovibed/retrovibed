package cmdcommunity

import (
	"os"

	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

type cmdCommunityDelete struct {
	Force bool   `flag:"" name:"force" help:"confirm deletion of the community"`
	Name  string `arg:"" name:"name" help:"name or id of the community to delete" required:"true"`
}

func (t cmdCommunityDelete) Run(gctx *cmdopts.Global, dpc cmdopts.DeeppoolClient) (err error) {
	if !t.Force {
		return errorsx.String("--force flag is required to confirm deletion")
	}

	c, err := dpc.HTTPClient(gctx.Context)
	if err != nil {
		return err
	}

	commresp, err := communityapi.CommunityDelete(gctx.Context, c, t.Name)
	if err != nil {
		return errorsx.Wrap(err, "failed to delete community")
	}

	if err = jsonx.MarshalWrite(os.Stdout, commresp.Community); err != nil {
		return errorsx.Wrap(err, "unable to write to encoder")
	}

	return nil
}
