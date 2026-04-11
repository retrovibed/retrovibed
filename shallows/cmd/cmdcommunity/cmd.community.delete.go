package cmdcommunity

import (
	"encoding/json"
	"os"

	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/metaapi"
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

	commresp, err := metaapi.CommunityDelete(gctx.Context, c, t.Name)
	if err != nil {
		return errorsx.Wrap(err, "failed to delete community")
	}

	if err = json.NewEncoder(os.Stdout).Encode(commresp.Community); err != nil {
		return errorsx.Wrap(err, "unable to write to encoder")
	}

	return nil
}
