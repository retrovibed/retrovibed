package library

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
)

func LocateOptionTestDefaults(t *Locate) {
	t.KnownMediaID = errorsx.Must(uuid.NewV4()).String()
}

func LocateQueryPending() squirrel.Sqlizer {
	return squirrel.Expr("library_locate.located_torrent_id = 'FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF'::uuid")
}

func LocateSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) LocateScanner {
	return NewLocateScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func LocateSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(LocateScannerStaticColumns)...).From("library_locate")
}
