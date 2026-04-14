package ddisctorrent

import (
	"time"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
)

const (
	ExtensionName         = "ddisc_2025_12_01"
	MethodMeta            = "ddisc_meta"             // meta endpoint used for identifying peers, provides some metadata about the node.
	MethodDisc            = "ddisc_query_discovered" // exchange discovered unknown media.
	MethodSync            = "ddisc_query_sync"       // request discovered media.
	MethodSearch          = "ddisc_query_search"     // search for a particular known bit of media.
	MethodMedia           = "ddisc_media"            // receive results from sync, discovered, and search.
	MethodRecommendations = "ddisc_recommended"      // receive recommendations.
)

type media struct {
	Peer                   krpc.ID       `bencode:"pid"`   // ID of the sending node.
	Token                  string        `bencode:"token"` // required: authorization token
	KnownMediaID           string        `bencode:"id"`    // required: known media id
	AudioBitrate           uint32        `bencode:"audio_bitrate"`
	AudioDefaultLocale     string        `bencode:"audio_default_locale"`
	Bytes                  uint64        `bencode:"bytes"`
	Collation              string        `bencode:"collation"`
	Description            string        `bencode:"description"`
	Infohash               []byte        `bencode:"infohash"`
	Mimetype               string        `bencode:"mimetype"`
	Partition              string        `bencode:"partition"`
	ReleasedAt             time.Time     `bencode:"released_at"`
	SubtitlesDefaultLocale string        `bencode:"subtitles_default_locale"`
	VideoResolution        string        `bencode:"video_resolution"`
	VideoRuntime           time.Duration `bencode:"video_runtime"`
}

func mediaToDiscovered(m media) func(*ddisc.Discovered) {
	return func(d *ddisc.Discovered) {
		d.AudioBitrate = m.AudioBitrate
		d.AudioDefaultLocale = m.AudioDefaultLocale
		d.Bytes = m.Bytes
		d.Collation = m.Collation
		d.Description = m.Description
		d.Infohash = m.Infohash
		d.Mimetype = m.Mimetype
		d.Partition = m.Partition
		d.ReleasedAt = m.ReleasedAt
		d.SubtitlesDefaultLocale = m.SubtitlesDefaultLocale
		d.VideoResolution = m.VideoResolution
		d.VideoRuntime = m.VideoRuntime
	}
}

func mediaFromDiscovered(token string, peer int160.T, d *ddisc.Discovered) media {
	return media{
		Token:                  token,
		Peer:                   peer.AsByteArray(),
		KnownMediaID:           d.KnownMediaID,
		AudioBitrate:           d.AudioBitrate,
		AudioDefaultLocale:     d.AudioDefaultLocale,
		Bytes:                  d.Bytes,
		Collation:              d.Collation,
		Description:            d.Description,
		Infohash:               d.Infohash,
		Mimetype:               d.Mimetype,
		Partition:              d.Partition,
		ReleasedAt:             d.ReleasedAt,
		SubtitlesDefaultLocale: d.SubtitlesDefaultLocale,
		VideoResolution:        d.VideoResolution,
		VideoRuntime:           d.VideoRuntime,
	}
}
