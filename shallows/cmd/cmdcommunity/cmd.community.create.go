package cmdcommunity

import (
	"encoding/json"
	"os"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type cmdCommunityCreate struct {
	Name        string `flag:"" name:"name" help:"name of the community globally unique. must be valid url subdomain" required:"true"`
	Description string `flag:"" name:"description" help:"description of the community"`
	Mimetype    string `flag:"" name:"mimetype" help:"mimetype for the community, used to specify the general type that will appear in the feed, defaults to application/octet-stream"`
}

func (t cmdCommunityCreate) Run(gctx *cmdopts.Global, dpc cmdopts.DeeppoolClient) (err error) {
	c, err := dpc.HTTPClient(gctx.Context)
	if err != nil {
		return err
	}

	commresp, err := metaapi.CommunityCreate(gctx.Context, c, &communityapi.CommunityCreateRequest{Community: &communityapi.Community{
		Domain:      t.Name,
		Description: t.Description,
		Mimetype:    langx.FirstNonZero(t.Mimetype, mimex.Binary),
	}})
	if err != nil {
		return errorsx.Wrap(err, "failed to create community")
	}

	if err = json.NewEncoder(os.Stdout).Encode(commresp.Community); err != nil {
		return errorsx.Wrap(err, "unable to write to encoder")
	}

	return nil
}
