package tracking

import (
	"github.com/james-lawrence/deeppool/internal/x/langx"
	"github.com/james-lawrence/deeppool/internal/x/md5x"
	"github.com/james-lawrence/torrent/dht/v2/krpc"
)

func PeerOptionBEP51(available uint64) func(*Peer) {
	return func(p *Peer) {
		p.Bep51 = true
		p.Bep51Available = available
	}
}

func NewPeer(md krpc.NodeInfo, options ...func(*Peer)) (m Peer) {
	return langx.Clone(Peer{
		ID:   md5x.Digest(md.ID[:]),
		IP:   md.Addr.IP.String(),
		Port: md.Addr.UDP().AddrPort().Port(),
	}, options...)
}
