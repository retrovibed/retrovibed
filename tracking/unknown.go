package tracking

import (
	"github.com/james-lawrence/deeppool/internal/x/langx"
	"github.com/james-lawrence/deeppool/internal/x/md5x"
	"github.com/james-lawrence/torrent/metainfo"
)

func NewUnknownHash(md metainfo.Hash, options ...func(*UnknownHash)) (m UnknownHash) {
	return langx.Clone(UnknownHash{
		ID:       md5x.Digest(md[:]),
		Infohash: md[:],
	}, options...)
}
