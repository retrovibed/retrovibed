package tracking

import (
	"context"
	"net/netip"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/internal/backoffx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/squirrelx"
	"github.com/retrovibed/retrovibed/internal/torrentx"
)

type OptionUnknownHash func(*UnknownHash)

func OptionUnknownHashPeer(id int160.T, ip netip.AddrPort) OptionUnknownHash {
	return func(uh *UnknownHash) {
		uh.Peer = id.Bytes()
		uh.IP = netip.IPv6Unspecified()
		if ip := ip.Addr(); ip.IsValid() {
			uh.IP = ip
		}
		uh.Port = ip.Port()
	}
}

func NewUnknownHash(md int160.T, options ...func(*UnknownHash)) (m UnknownHash) {
	return langx.Clone(UnknownHash{
		ID:        torrentx.HashUID(&md),
		Infohash:  md.Bytes(),
		NextCheck: time.Now().Add(backoffx.DynamicHash5m(md.String())),
	}, options...)
}

func UnknownHashQueryNeedsCheck() squirrel.Sqlizer {
	return squirrel.Expr("torrents_unknown_infohashes.next_check < NOW()")
}

func UnknownSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) UnknownHashScanner {
	return NewUnknownHashScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func UnknownSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(UnknownHashScannerStaticColumns)...).From("torrents_unknown_infohashes")
}
