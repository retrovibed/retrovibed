package community

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
)

func CommunitySearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) CommunityScanner {
	return NewCommunityScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func CommunitySearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(CommunityScannerStaticColumns)...).From("community")
}

func CommunityQueryNotTombstoned() squirrel.Sqlizer {
	return squirrel.Expr("community.tombstoned_at = 'infinity'")
}
