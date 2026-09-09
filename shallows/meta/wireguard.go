package meta

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
)

type WireguardNetworkRing uint32

const (
	WireguardNetworkRingUnspecified uint32 = iota
	WireguardNetworkRingDistribution
	WireguardNetworkRingSocial
	WireguardNetworkRingMax
)

func WireguardSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) WireguardScanner {
	return NewWireguardScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func WireguardSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(WireguardScannerStaticColumns)...).From("meta_wireguard")
}

func WireguardOptionDescription(s string) func(*Wireguard) {
	return func(w *Wireguard) {
		w.Description = s
	}
}

func WireguardOptionAutoID(w *Wireguard) {
	w.ID = uuid.Must(uuid.NewV4()).String()
}

func WireguardOptionDistribution(w *Wireguard) {
	w.Nettype = WireguardNetworkRingDistribution
}

func WireguardOptionDNSRateLimit(n uint32) func(*Wireguard) {
	return func(w *Wireguard) {
		w.RateLimitDNS = n
	}
}

func WireguardOptionOutboundRateLimit(n uint32) func(*Wireguard) {
	return func(w *Wireguard) {
		w.RateLimitOutbound = n
	}
}

func WireguardOptionTestDefauts(w *Wireguard) {
	w.ID = uuid.Must(uuid.NewV4()).String()
	w.Nettype = w.Nettype % (WireguardNetworkRingMax)
}

func NewWireguard(uid string, options ...func(*Wireguard)) Wireguard {
	return langx.Clone(Wireguard{
		ID:          uid,
		Description: "",
		Nettype:     WireguardNetworkRingUnspecified,
	}, options...)
}
