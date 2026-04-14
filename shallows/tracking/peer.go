package tracking

import (
	"context"
	"net/netip"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/dht/krpc"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/uuidx"
)

type PeerOption func(*Peer)

func PeerOptionNoop(*Peer) {}

func PeerOptionBEP51(available uint64, ttl uint16) func(*Peer) {
	return func(p *Peer) {
		p.Bep51 = true
		p.Bep51Available = available
		p.Bep51TTL = ttl
	}
}

func PeerOptionDdisc(partition, offset uuid.UUID) func(*Peer) {
	return func(p *Peer) {
		p.Ddisc = !partition.IsZero()
		p.DdiscPartition = partition.String()
		p.DdiscSyncoffset = offset.String()
	}
}

func PeerOptionTimestampClone(p0 Peer) func(*Peer) {
	return func(p1 *Peer) {
		p1.CreatedAt = p0.CreatedAt
		p1.UpdatedAt = p0.UpdatedAt
		p1.NextCheck = p0.NextCheck
		p1.TombstonedAt = p0.TombstonedAt
	}
}

func PeerOptionDescription(s string) func(*Peer) {
	return func(p *Peer) {
		p.Description = s
	}
}

func PeerOptionTombstone(ts time.Time) func(*Peer) {
	return func(p *Peer) {
		p.TombstonedAt = ts
	}
}

func PeerOptionTestDefaults(p *Peer) {
	p.Peer = int160.Random().Bytes()
	p.ID = md5x.FormatUUID(md5x.Digest(p.Peer))
	p.Ddisc = true
	p.Bep51 = true
	p.DdiscPartition = uuidx.WithSuffix(1)
	p.DdiscSyncoffset = uuidx.WithSuffix(1)
}

func PeerOptionIP(n netip.AddrPort) func(*Peer) {
	return func(p *Peer) {
		p.IP = n.Addr()
		p.Port = n.Port()
	}
}

func PeerOptionNetwork(n string) func(*Peer) {
	return func(p *Peer) {
		p.Network = n
	}
}

func NewPeerFromInfo(info krpc.NodeInfo, options ...func(*Peer)) (m Peer) {
	return NewPeer(
		info.ID.Int160(),
		PeerOptionIP(info.Addr.AddrPort),
		PeerOptionNetwork(info.Addr.UDP().Network()),
		langx.Compose(options...),
	)
}

// generate the unique uuid from the public int160 id.
func PeerUID(id int160.T) string {
	return md5x.FormatUUID(md5x.Digest(id.Bytes()))
}

func NewPeer(id int160.T, options ...func(*Peer)) (m Peer) {
	return langx.Clone(Peer{
		ID:              PeerUID(id),
		Peer:            id.Bytes(),
		Network:         "udp", // default to udp
		IP:              netip.IPv4Unspecified(),
		Port:            0,
		TombstonedAt:    time.Now().Add(30 * 24 * time.Hour),
		Ddisc:           false,
		DdiscPartition:  uuid.Nil.String(),
		DdiscSyncoffset: uuid.Nil.String(),
	}, options...)
}

func PeerQueryDdiscEnabled() squirrel.Sqlizer {
	return squirrel.Expr("torrents_peers.ddisc")
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
