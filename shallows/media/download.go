package media

import (
	"github.com/james-lawrence/deeppool/internal/x/langx"
	"github.com/james-lawrence/deeppool/tracking"
)

type DownloadOption func(*Download)

func DownloadOptionFromTorrentMetadata(cc tracking.Metadata) DownloadOption {
	return func(c *Download) {
		c.Media = langx.Autoptr(langx.Clone(Media{}, MediaOptionFromTorrentMetadata(cc)))
		c.Bytes = cc.Bytes
		c.Downloaded = cc.Downloaded
	}
}
