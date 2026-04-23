package cmdcommunity

import (
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdcommunitylibrary"
)

type Community struct {
	Deeppool Deeppool                    `cmd:"" help:"deeppool community commands"`
	Library  cmdcommunitylibrary.Library `cmd:"" help:"library community commands"`
}

type Deeppool struct {
	Create  cmdCommunityCreate  `cmd:"" help:"create a new community"`
	Delete  cmdCommunityDelete  `cmd:"" help:"delete a community"`
	Update  cmdCommunityUpdate  `cmd:"" help:"update a community"`
	Info    cmdCommunityInfo    `cmd:"" help:"display the community details"`
	Publish cmdCommunityPublish `cmd:"" help:"create and publish a rss feed from a list of torrents"`
	Add     cmdCommunityAdd     `cmd:"" help:"add content to a community in deeppool"`
	List    cmdCommunityList    `cmd:"" help:"list published content for a community"`
	Import  cmdCommunityImport  `cmd:"" help:"import published content to local database from stdin"`
	Sync    cmdCommunitySync    `cmd:"" help:"sync published content from deeppool to local database"`
}
