package meta

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/squirrelx"
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

func WireguardOptionDefault(w *Wireguard) {
	w.Default = true
}

func WireguardOptionDNSRateLimit(n uint32) func(*Wireguard) {
	return func(w *Wireguard) {
		w.DNSRateLimit = n
	}
}

func NewWireguard(uid string, options ...func(*Wireguard)) Wireguard {
	return langx.Clone(Wireguard{
		ID:          uid,
		Description: "",
		Default:     false,
	}, options...)
}
