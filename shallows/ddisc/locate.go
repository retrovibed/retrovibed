package ddisc

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
)

type LocateOption func(*Locate)

func LocateQueryPending() squirrel.Sqlizer {
	return squirrel.Expr("ddisc_locate.tombstoned_at = 'infinity'::timestamptz")
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
	}, options...)
}

func LocateOptionKnownMedia(id string) LocateOption {
	return func(l *Locate) { l.KnownMediaID = id }
}

func LocateOptionAutoDownload(v bool) LocateOption {
	return func(l *Locate) { l.Autodownload = v }
}

func LocateOptionAdult(v bool) LocateOption {
	return func(l *Locate) { l.Adult = v }
}
