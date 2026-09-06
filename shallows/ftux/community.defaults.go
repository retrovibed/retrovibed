package ftux

import (
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
)

// PrepareDefaultCommunities returns the curated set of community subscription
// suggestions shown during first-time setup, caching them locally on first run.
func PrepareDefaultCommunities() (_ []*communityapi.Community, err error) {
	encoded, err := fsx.AutoCached(userx.DefaultConfigDir(userx.DefaultRelRoot(), "default.communities.json"), func() ([]byte, error) {
		return jsonx.Marshal([]*communityapi.Community{
			{
				Id:          "e54c2d7b-70f9-41fb-92f7-eeb0968f4be4",
				Description: "Retrovibed - neurals enabling various small AI driven functionality for retrovibed",
				Url:         community.CommunityURLFromDomain("neurals"),
			},
			{
				Id:          "8a2e8383-0e50-40a1-8af6-69d4f8c67cb0",
				Description: "Retrovibed - content publishing support",
				Url:         community.CommunityURLFromDomain("retropublish"),
			},
			{
				Id:          "4b43c380-89e5-44f0-a5a1-c4e0b52a4bef",
				Description: "Retrovibed - media metadata. posters, ratings, descriptions. recommended for desktop/mobile. (~3 GiB)",
				Url:         community.CommunityURLFromDomain("media"),
			},
		})

	})
	if err != nil {
		return nil, err
	}

	var out []*communityapi.Community
	return out, jsonx.Unmarshal(encoded, &out)
}
