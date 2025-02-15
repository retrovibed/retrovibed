package tracking

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/deeppool/internal/x/sqlx"
	"github.com/james-lawrence/deeppool/internal/x/squirrelx"
)

func RSSQueryNeedsCheck() squirrel.Sqlizer {
	return squirrel.Expr("torrents_feed_rss.next_check < NOW()")
}

func RSSSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) RSSScanner {
	return NewRSSScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func RSSSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(RSSScannerStaticColumns)...).From("torrents_feed_rss")
}
