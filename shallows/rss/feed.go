package rss

import (
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type FeedOption func(*Feed)

func FeedOptionFromTorrentRSS(cc tracking.RSS) FeedOption {
	return func(c *Feed) {
		c.Id = cc.ID
		c.Description = cc.Description
		c.Url = cc.URL
		c.Autodownload = cc.Autodownload
		c.Autoarchive = cc.Autoarchive
		c.Contributing = cc.Contributing
		c.EncryptionSeed = cc.EncryptionSeed
		c.Digest = cc.Digest
		c.CreatedAt = grpcx.EncodeTime(cc.CreatedAt)
		c.UpdatedAt = grpcx.EncodeTime(cc.UpdatedAt)
		c.NextCheck = grpcx.EncodeTime(cc.NextCheck)
	}
}
