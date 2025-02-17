package tracking

import (
	"github.com/james-lawrence/deeppool/internal/x/langx"
	"github.com/james-lawrence/deeppool/internal/x/md5x"
	"github.com/james-lawrence/torrent/metainfo"
)

func HashUID(md *metainfo.Hash) string {
	return md5x.Digest(md.Bytes())
}

func MetadataOptionFromInfo(i *metainfo.Info) func(*Metadata) {
	return func(m *Metadata) {
		m.Description = i.Name
		m.Bytes = uint64(i.TotalLength())
	}
}

func NewMetadata(md *metainfo.Hash, options ...func(*Metadata)) (m Metadata) {
	r := langx.Clone(Metadata{
		ID:       HashUID(md),
		Infohash: md.Bytes(),
	}, options...)
	return r
}
