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

	JSON                      = "application/json"
	Binary                    = "application/octet-stream"
	Bittorrent                = "application/x-bittorrent"
	Magnet                    = "application/x-magnet"
	RSS                       = "application/rss+xml"
	ISO9660                   = "application/x-iso9660-image"
	RetrovibedMediaArchive    = "application/vnd.retrovibed.media.archive"
	RetrovibedNeural          = "application/vnd.retrovibed.neural"
	RetrovibedDiscoverySearch = "application/vnd.retrovibed.discovery.search.module" // a discovery search module.
	RetrovibedDiscoveryAudio  = "application/vnd.retrovibed.discovery.audio"         // generic audio only mimetype, contextually similar to music + podcasts + audio only formats.
	RetrovibedDiscoveryVideo  = "application/vnd.retrovibed.discovery.video"         // generic video mimetype, contextually similar to combined movies + tv and other video formats.
	RetrovibedDiscoveryMusic  = "application/vnd.retrovibed.discovery.music"
	RetrovibedDiscoveryMovies = "application/vnd.retrovibed.discovery.movies"
	RetrovibedDiscoveryTV     = "application/vnd.retrovibed.discovery.tv"
)

func Category(mime string) string {
	if prefix, _, ok := strings.Cut(mime, "/"); ok {
		return prefix
	}

	return mime
}
