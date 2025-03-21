package tracking

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/squirrelx"
)

func NewUnknownHash(md metainfo.Hash, options ...func(*UnknownHash)) (m UnknownHash) {
	return langx.Clone(UnknownHash{
		ID:       HashUID(&md),
		Infohash: md[:],
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
