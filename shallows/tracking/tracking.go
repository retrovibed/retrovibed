package tracking

import (
	"github.com/james-lawrence/deeppool/internal/x/md5x"
	"github.com/james-lawrence/torrent/metainfo"
)

func HashUID(md *metainfo.Hash) string {
	return md5x.Digest(md.Bytes())
}
