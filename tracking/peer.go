package tracking

import (
	"encoding/hex"

	"github.com/james-lawrence/deeppool/internal/x/langx"
	"github.com/james-lawrence/deeppool/internal/x/md5x"
	"github.com/james-lawrence/torrent/metainfo"
)

func OptionFromInfo(i *metainfo.Info) func(*Metadata) {
	return func(m *Metadata) {
		m.Description = i.Name
		m.Bytes = uint64(i.TotalLength())
	}
}

func NewPeer(md *metainfo.Hash, options ...func(*Metadata)) (m Metadata) {
	r := langx.Clone(Metadata{
		ID:       md5x.Digest(md.Bytes()),
		Infohash: hex.EncodeToString(md.Bytes()),
	}, options...)
	return r
}
