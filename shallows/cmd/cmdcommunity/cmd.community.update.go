package cmdcommunity

import (
	"encoding/json"
	"os"

	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/meta"
	"github.com/retrovibed/retrovibed/metaapi"
)

type cmdCommunityUpdate struct {
	Description string `flag:"" name:"description" help:"description of the community"`
	Mimetype    string `flag:"" name:"mimetype" help:"mimetype for the community, used to specify the general type that will appear in the feed"`
	Name        string `arg:"" name:"name" help:"name or id of the community to update" required:"true"`
}

func (t cmdCommunityUpdate) Run(gctx *cmdopts.Global, dpc cmdopts.DeeppoolClient) (err error) {
	c, err := dpc.HTTPClient(gctx.Context)
	if err != nil {
		return err
	}

	commresp, err := metaapi.CommunityUpdate(gctx.Context, c, t.Name, &meta.CommunityUpdateRequest{
		Community: &meta.Community{
			Description: t.Description,
			Mimetype:    t.Mimetype,
		},
	})
	if err != nil {
		return errorsx.Wrap(err, "failed to update community")
	}

	if err = json.NewEncoder(os.Stdout).Encode(commresp.Community); err != nil {
		return errorsx.Wrap(err, "unable to write to encoder")
	}

	return nil
}
