package media

import "github.com/james-lawrence/deeppool/tracking"

type MediaOption func(*Media)

func MediaOptionFromTorrentMetadata(cc tracking.Metadata) MediaOption {
	return func(c *Media) {
		c.Title = cc.Description
		c.Mimetype = "applications/x-bittorrent"
		// c.CreatedAt = grpcx.EncodeTime(cc.CreatedAt)
		// c.InitiatedAt = grpcx.EncodeTime(cc.InitiatedAt)
		// c.CompletedAt = grpcx.EncodeTime(cc.CompletedAt)
	}
}
