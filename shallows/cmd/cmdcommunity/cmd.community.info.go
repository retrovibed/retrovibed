package cmdcommunity

import (
	"encoding/json"
	"os"

	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/metaapi"
)

type cmdCommunityInfo struct {
	Description string `flag:"" name:"description" help:"description of the community"`
	Name        string `arg:"" name:"name" help:"name of the community globally unique. must be valid url subdomain" required:"true"`
}

func (t cmdCommunityInfo) Run(gctx *cmdopts.Global, dpc cmdopts.DeeppoolClient) (err error) {
	c, err := dpc.HTTPClient(gctx.Context)
	if err != nil {
		return err
	}

	commresp, err := metaapi.CommunityInfo(gctx.Context, c, t.Name)
	if err != nil {
		return errorsx.Wrap(err, "failed to locate community")
	}

	if err = json.NewEncoder(os.Stdout).Encode(commresp.Community); err != nil {
		return errorsx.Wrap(err, "unable to write to encoder")
	}

	return nil
}
