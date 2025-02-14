package tracking

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/deeppool/internal/x/langx"
	"github.com/james-lawrence/deeppool/internal/x/sqlx"
	"github.com/james-lawrence/deeppool/internal/x/squirrelx"
	"github.com/james-lawrence/torrent/metainfo"
)

func NewUnknownHash(md metainfo.Hash, options ...func(*UnknownHash)) (m UnknownHash) {
	return langx.Clone(UnknownHash{
		ID:       HashUID(&md),
		Infohash: md[:],
	}, options...)
}

func UnknownSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) UnknownHashScanner {
	return NewUnknownHashScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func UnknownSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(UnknownHashScannerStaticColumns)...).From("torrents_unknown_infohashes")
}
