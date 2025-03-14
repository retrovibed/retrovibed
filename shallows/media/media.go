package media

import (
	"github.com/james-lawrence/deeppool/internal/x/grpcx"
	"github.com/james-lawrence/deeppool/library"
	"github.com/james-lawrence/deeppool/tracking"
)

type MediaOption func(*Media)

func MediaOptionFromLibraryMetadata(cc library.Metadata) MediaOption {
	return func(c *Media) {
		c.Id = cc.ID
		c.Description = cc.Description
		c.Mimetype = cc.Mimetype
		c.TorrentId = cc.TorrentID
		c.ArchiveId = cc.ArchiveID
		c.CreatedAt = grpcx.EncodeTime(cc.CreatedAt)
		c.UpdatedAt = grpcx.EncodeTime(cc.UpdatedAt)
	}
}

func MediaOptionFromTorrentMetadata(cc tracking.Metadata) MediaOption {
	return func(c *Media) {
		c.Id = cc.ID
		c.Description = cc.Description
		c.CreatedAt = grpcx.EncodeTime(cc.CreatedAt)
		c.UpdatedAt = grpcx.EncodeTime(cc.UpdatedAt)
		c.Mimetype = "applications/x-bittorrent"
	}
}
