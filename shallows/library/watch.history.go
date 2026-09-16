package library

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

func WatchHistorySearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) WatchHistoryScanner {
	return NewWatchHistoryScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func WatchHistorySearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(WatchHistoryScannerStaticColumns)...).From("library_watch_history").Where(squirrel.Expr("'t'"))
}

func WatchHistoryQueryCreated(r timex.Range) squirrel.Sqlizer {
	return squirrelx.Between("library_watch_history.created_at", r.Start, r.End)
}

func WatchHistoryQueryMediaID(v string) squirrel.Sqlizer {
	return squirrel.Expr("library_watch_history.media_id = ?", v)
}

func WatchHistoryQueryProfileID(v string) squirrel.Sqlizer {
	return squirrel.Expr("library_watch_history.profile_id = ?", v)
}
