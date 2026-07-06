package tracking

import (
	"context"
	"math"
	"net/netip"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/shallows/internal/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/ducktype"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/torrentx"
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

func UnknownHashOptionTestDefaults(uh *UnknownHash) {
	id := int160.Random()
	uh.ID = torrentx.HashUID(&id)
	uh.Infohash = id.Bytes()
	uh.IP = netip.IPv6Unspecified()
	uh.Attempts &= math.MaxInt64
}

func UnknownHashQueryNextCheck(r timex.Range) squirrel.Sqlizer {
	return squirrelx.Between("torrents_unknown_infohashes.next_check", ducktype.NewNullTime(r.Start), ducktype.NewNullTime(r.End))
}

func UnknownHashQueryByIDs(ids ...string) squirrel.Sqlizer {
	if len(ids) == 0 {
		return squirrelx.Noop{}
	}

	return squirrel.Eq{"torrents_unknown_infohashes.id": ids}
}

func UnknownHashQueryAttemptsRange(min, max uint64) squirrel.Sqlizer {
	return squirrelx.Between("torrents_unknown_infohashes.attempts", min, max)
}

func UnknownSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) UnknownHashScanner {
	return NewUnknownHashScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func UnknownSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(UnknownHashScannerStaticColumns)...).From("torrents_unknown_infohashes")
}
