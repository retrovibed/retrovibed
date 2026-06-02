package community

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
)

func PublishedContentSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) PublishedContentScanner {
	return NewPublishedContentScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func PublishedContentSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(PublishedContentScannerStaticColumns)...).From("published_content")
}

func PublishedContentQueryCommunityID(cid string) squirrel.Sqlizer {
	return squirrel.Expr("published_content.community_id = ?", cid)
}

func PublishedContentQueryNotTombstoned() squirrel.Sqlizer {
	return squirrel.Expr("published_content.tombstoned_at = 'infinity'")
}
