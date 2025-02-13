package tracking

// import (
// 	"encoding/hex"

// 	"github.com/james-lawrence/deeppool/internal/x/langx"
// 	"github.com/james-lawrence/deeppool/internal/x/md5x"
// 	"github.com/james-lawrence/torrent/metainfo"
// )

// func NewPeer(md *metainfo.Hash, options ...func(*Peer)) (m Peer) {
// 	r := langx.Clone(Metadata{
// 		ID:       md5x.Digest(md.Bytes()),
// 		Infohash: hex.EncodeToString(md.Bytes()),
// 	}, options...)
// 	return r
// }
