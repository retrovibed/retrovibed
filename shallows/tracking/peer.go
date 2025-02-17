package tracking

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/deeppool/internal/x/langx"
	"github.com/james-lawrence/deeppool/internal/x/md5x"
	"github.com/james-lawrence/deeppool/internal/x/sqlx"
	"github.com/james-lawrence/deeppool/internal/x/squirrelx"
	"github.com/james-lawrence/torrent/dht/v2/krpc"
)

func PeerOptionBEP51(available uint64, ttl uint16) func(*Peer) {
	return func(p *Peer) {
		p.Bep51 = true
		p.Bep51Available = available
		p.Bep51TTL = ttl
	}
}

func NewPeer(md krpc.NodeInfo, options ...func(*Peer)) (m Peer) {
	return langx.Clone(Peer{
		ID:      md5x.Digest(md.ID[:]),
		Peer:    md.ID[:],
		Network: md.Addr.UDP().Network(),
		IP:      md.Addr.IP.String(),
		Port:    md.Addr.UDP().AddrPort().Port(),
	}, options...)
}

func PeerQueryNeedsCheck() squirrel.Sqlizer {
	return squirrel.Expr("torrents_peers.next_check < NOW()")
}

func PeerQueryHasInfoHashes() squirrel.Sqlizer {
	return squirrel.Expr("torrents_peers.bep51_available > 0")
}

func PeerSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) PeerScanner {
	return NewPeerScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func PeerSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(PeerScannerStaticColumns)...).From("torrents_peers")
}
