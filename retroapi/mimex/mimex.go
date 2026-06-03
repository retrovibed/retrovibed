package mimex

import (
	"strings"
)

const (
	Video       = "video"
	Audio       = "audio"
	Image       = "image"
	Text        = "text"
	Application = "application"

	JSON                   = "application/json"
	Binary                 = "application/octet-stream"
	Bittorrent             = "application/x-bittorrent"
	RSS                    = "application/rss+xml"
	ISO9660                = "application/x-iso9660-image"
	RetrovibedMediaArchive = "application/vnd.retrovibed.media.archive"
	RetrovibedNeural       = "application/vnd.retrovibed.neural"
)

func Category(mime string) string {
	if prefix, _, ok := strings.Cut(mime, "/"); ok {
		return prefix
	}

	return mime
}
