package rss

import (
	"github.com/james-lawrence/deeppool/internal/x/grpcx"
	"github.com/james-lawrence/deeppool/tracking"
)

type FeedOption func(*Feed)

func FeedOptionFromTorrentRSS(cc tracking.RSS) FeedOption {
	return func(c *Feed) {
		c.Id = cc.ID
		c.Description = cc.Description
		c.Url = cc.URL
		c.Autodownload = cc.Autodownload
		c.CreatedAt = grpcx.EncodeTime(cc.CreatedAt)
		c.UpdatedAt = grpcx.EncodeTime(cc.UpdatedAt)
		c.NextCheck = grpcx.EncodeTime(cc.NextCheck)
	}
}
