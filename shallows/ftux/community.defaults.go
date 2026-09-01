package ftux

import (
	"encoding/json"

	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
)

// PrepareDefaultCommunities returns the curated set of community subscription
// suggestions shown during first-time setup, caching them locally on first run.
func PrepareDefaultCommunities() (_ []*communityapi.Community, err error) {
	encoded, err := fsx.AutoCached(userx.DefaultConfigDir(userx.DefaultRelRoot(), "default.communities.json"), func() ([]byte, error) {
		return json.Marshal([]*communityapi.Community{
			{
				Id:          "00000000-0000-0000-0000-000000000001",
				Description: "Retrovibed - test data",
				Url:         community.CommunityURLFromDomain("vibed"),
			},
			{
				Id:          "00000000-0000-0000-0000-000000000002",
				Description: "Retrovibed - media metadata. posters, ratings, descriptions. (~3 GiB)",
				Url:         community.CommunityURLFromDomain("media"),
			},
			{
				Id:          "00000000-0000-0000-0000-000000000003",
				Description: "Retroneural - enables various small AI driven functionality for retrovibed",
				Url:         community.CommunityURLFromDomain("neurals"),
			},
		})
	})
	if err != nil {
		return nil, err
	}

	var out []*communityapi.Community
	err = json.Unmarshal(encoded, &out)
	return out, err
}
