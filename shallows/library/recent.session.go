package library

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

func RecentSessionSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) RecentSessionScanner {
	return NewRecentSessionScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func RecentSessionLibrarySearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) RecentSessionLibraryScanner {
	return NewRecentSessionLibraryScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func RecentSessionLibrarySearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(RecentSessionScannerStaticColumns, MetadataScannerStaticColumns)...).From("library_metadata").InnerJoin("library_recent_sessions ON library_recent_sessions.media_id = library_metadata.id").Where(squirrel.Expr("'t'"))
}

func RecentSessionQueryCreated(r timex.Range) squirrel.Sqlizer {
	return squirrelx.Between("library_recent_sessions.created_at", r.Start, r.End)
}

type RecentSessionOption func(*RecentSession)

func RecentSessionOptionID(id string) RecentSessionOption {
	return func(s *RecentSession) {
		s.ID = id
	}
}

func RecentSessionOptionMediaID(id string) RecentSessionOption {
	return func(s *RecentSession) {
		s.MediaID = id
	}
}

func RecentSessionOptionTestDefaults(s *RecentSession) {
	s.ID = uuid.Nil.String()
	s.MediaID = uuid.Nil.String()
	s.Query = nil
	s.Position = 0
	s.Duration = 0
}
