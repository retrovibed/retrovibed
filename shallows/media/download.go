package media

import (
	"github.com/retrovibed/retrovibed/internal/x/langx"
	"github.com/retrovibed/retrovibed/tracking"
)

type DownloadOption func(*Download)

func DownloadOptionFromTorrentMetadata(cc tracking.Metadata) DownloadOption {
	return func(c *Download) {
		c.Media = langx.Autoptr(langx.Clone(Media{}, MediaOptionFromTorrentMetadata(cc)))
		c.Bytes = cc.Bytes
		c.Downloaded = cc.Downloaded
		c.Peers = uint32(cc.Peers)
	}
}
