package cmdcommunity

import (
	"encoding/json"
	"os"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/metaapi"
)

type cmdCommunityUpdate struct {
	Description        *string `flag:"" name:"description" help:"description of the community"`
	Mimetype           *string `flag:"" name:"mimetype" help:"mimetype for the community, used to specify the general type that will appear in the feed"`
	DefaultPublishMode *string `flag:"" name:"default-publish-mode" enum:"UNLISTED,LISTED,SYNDICATED" help:"default publish mode for content in the community"`
	DefaultTtl         *uint64 `flag:"" name:"default-ttl" help:"default time-to-live for content in the community (seconds)"`
	DefaultLanguage    *string `flag:"" name:"default-language" help:"default language for content in the community"`
	Hidden             *bool   `flag:"" name:"hidden" help:"hide the community from public listings"`
	Adult              *bool   `flag:"" name:"adult" help:"mark the community as containing adult content"`
	Name               string  `arg:"" name:"name" help:"name or id of the community to update" required:"true"`
}

func (t cmdCommunityUpdate) publishmode(m string, fallback communityapi.PublishMode) *communityapi.PublishMode {
	switch langx.Autoderef(t.DefaultPublishMode) {
	case communityapi.PublishMode_UNLISTED.String():
		return new(communityapi.PublishMode_UNLISTED)
	case communityapi.PublishMode_LISTED.String():
		return new(communityapi.PublishMode_LISTED)
	case communityapi.PublishMode_SYNDICATED.String():
		return new(communityapi.PublishMode_SYNDICATED)
	default:
		return new(fallback)
	}
}

func (t cmdCommunityUpdate) Run(gctx *cmdopts.Global, dpc cmdopts.DeeppoolClient) (err error) {
	c, err := dpc.HTTPClient(gctx.Context)
	if err != nil {
		return err
	}

	inforesp, err := metaapi.CommunityInfo(gctx.Context, c, t.Name)
	if err != nil {
		return errorsx.Wrap(err, "failed to read community")
	}
	current := inforesp.Community

	commresp, err := metaapi.CommunityUpdate(gctx.Context, c, t.Name, &communityapi.CommunityUpdateRequest{
		Community: &communityapi.Community{
			Description:        *langx.FirstNonZero(t.Description, &current.Description),
			Mimetype:           *langx.FirstNonZero(t.Mimetype, &current.Mimetype),
			DefaultPublishMode: *t.publishmode(langx.Zero(t.DefaultPublishMode), current.DefaultPublishMode),
			DefaultTtl:         *langx.FirstNonZero(t.DefaultTtl, &current.DefaultTtl),
			DefaultLanguage:    *langx.FirstNonZero(t.DefaultLanguage, &current.DefaultLanguage),
			Hidden:             *langx.FirstNonZero(t.Hidden, &current.Hidden),
			Adult:              *langx.FirstNonZero(t.Adult, &current.Adult),
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
