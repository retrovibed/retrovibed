package ddisc

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
)

type LocateOption func(*Locate)

func LocateQueryPending() squirrel.Sqlizer {
	return squirrel.And{
		squirrel.Expr("ddisc_locate.tombstoned_at = 'infinity'::timestamptz"),
		squirrel.Expr("ddisc_locate.next_check_at <= NOW()"),
	}
}

func LocateQueryCompleted() squirrel.Sqlizer {
	return squirrel.Expr("ddisc_locate.tombstoned_at <= NOW()")
}

func LocateQueryByIDs(ids ...string) squirrel.Sqlizer {
	if len(ids) == 0 {
		return squirrelx.Noop{}
	}
	return squirrel.Eq{"ddisc_locate.id": ids}
}

func LocateQueryAttemptsRange(min, max uint64) squirrel.Sqlizer {
	return squirrelx.Between("ddisc_locate.attempts", min, max)
}

func LocateQueryText(query string) squirrel.Sqlizer {
	if query == "" {
		return squirrelx.Noop{}
	}
	return lucenex.Query(duckdbx.NewLucene(), query, lucenex.WithDefaultField("query"))
}

func LocateSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) LocateScanner {
	return NewLocateScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func LocateSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(LocateScannerStaticColumns)...).From("ddisc_locate")
}

// NewLocate builds a locate request, deriving a deterministic id from
// (query, mimetype) so re-requesting the same thing is idempotent.
// KnownMediaID defaults to Nil (unresolved) - pass LocateOptionKnownMedia
// when the caller already resolved a catalog entry.
func NewLocate(query, mimetype string, options ...LocateOption) (l Locate) {
	return langx.Clone(Locate{
		ID:           md5x.FormatUUID(md5x.Digest(query, mimetype)),
		Query:        query,
		Mimetype:     mimetype,
		KnownMediaID: uuid.Nil.String(),
		ProfileID:    uuid.Nil.String(),
	}, options...)
}

func LocateOptionKnownMedia(id string) LocateOption {
	return func(l *Locate) { l.KnownMediaID = id }
}

func LocateOptionProfile(id string) LocateOption {
	return func(l *Locate) { l.ProfileID = id }
}

func LocateOptionAutoDownload(v bool) LocateOption {
	return func(l *Locate) { l.Autodownload = v }
}

func LocateOptionAdult(v bool) LocateOption {
	return func(l *Locate) { l.Adult = v }
}
