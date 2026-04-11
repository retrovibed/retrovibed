package media

import (
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/internal/grpcx"
	"github.com/retrovibed/retrovibed/tracking"
)

type PublishedOption func(*Published)

func PublishedOptionFromTorrentMetadata(cc tracking.Metadata) PublishedOption {
	return func(c *Published) {
		c.Id = int160.FromBytes(cc.Infohash).String()
		c.Description = cc.Description
		c.Mimetype = cc.Mimetype
		c.Bytes = cc.Bytes
		c.Entropy = cc.EncryptionSeed
		c.ExpiresAt = grpcx.EncodeTime(cc.ExpiresAt)
	}
}
