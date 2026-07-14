package library

import (
	"context"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/internal/ducktype"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/internal/unicodex"
)

type QueryCleaner interface {
	Clean(ctx context.Context, text string) (string, error)
}

func NewQueryCleanerFn(fn func(text string) string) QueryCleanerFn {
	return QueryCleanerFn(func(_ context.Context, text string) (string, error) {
		return fn(text), nil
	})
}

type QueryCleanerFn func(_ context.Context, text string) (string, error)

func (fn QueryCleanerFn) Clean(ctx context.Context, text string) (string, error) {
	return fn(ctx, text)
}

func QueryCleanerNoop() *NoopQueryCleaner {
	return &NoopQueryCleaner{}
}

type NoopQueryCleaner struct{}

func (NoopQueryCleaner) Clean(_ context.Context, text string) (string, error) {
	return text, nil
}

func KnownOptionRandomID(t *Known) {
	t.ID = errorsx.Must(uuid.NewV4()).String()
}

func KnownOptionReleased(ts time.Time) func(*Known) {
	return func(t *Known) {
		t.Released = ts
	}
}

func KnownOptionMimetype(v string) func(*Known) {
	return func(t *Known) {
		t.Mimetype = v
	}
}

func KnownOptionSource(v string) func(*Known) {
	return func(t *Known) {
		t.Source = v
	}
}

func KnownOptionTestNoPoster(t *Known) {
	t.PosterPath = ""
	t.BackdropPath = ""
}

func KnownOptionTestDefaults(t *Known) {
	t.UID = errorsx.Must(uuid.NewV4()).String()
	t.Md5 = errorsx.Must(uuid.NewV4()).String()
	t.Adult = false
	t.Released = time.Now()
	t.Mimetype = mimex.Application
	t.Duplicates = 0
	t.Popularity = 0
	t.OriginalLanguage = userx.LocaleLanguage()
}

// ImportPrefix is a type constraint for import source prefixes.
type importprefix interface {
	~string
}

// create a unique import id from a uint sequence.
func KnownImportedUintID[P importprefix](prefix P, id uint64) string {
	l := id & 0x0000FFFFFFFFFFFF
	h := id & 0xFFFF000000000000 >> 56
	return fmt.Sprintf("%x-0000-0000-%04x-%012x", fnv.New32().Sum([]byte(prefix))[:4], h, l)
}

// create a unique import id from a uuid by mutating its first 4 bytes with the prefix checksum.
func KnownImportedUUID[P importprefix](prefix P, id uuid.UUID) uuid.UUID {
	copy(id[:4], fnv.New32().Sum([]byte(prefix))[:4])
	return id
}

func Unknown() Known {
	return Known{
		UID: uuid.Nil.String(),
	}
}

func KnownSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) KnownScanner {
	return NewKnownScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func KnownQueryUIDGreaterThan(uid string) squirrel.Sqlizer {
	return squirrel.Expr("library_known_media.uid > ?", uid)
}

// KnownQueryExplicit toggles whether adult content is allowed in results.
// allow=false restricts results to non-adult content; allow=true permits
// both adult and non-adult content (it does not restrict to adult-only).
func KnownQueryExplicit(allow bool) squirrel.Sqlizer {
	if allow {
		return squirrelx.Noop{}
	}

	return squirrel.Expr("library_known_media.adult = ?", false)
}

func KnownQueryUID(ids ...string) squirrel.Sqlizer {
	return squirrelx.In("library_known_media.uid", ids...)
}

func KnownQueryLanguage(v string) squirrel.Sqlizer {
	if stringsx.Blank(v) {
		return squirrelx.Noop{}
	}

	return squirrel.Expr("library_known_media.original_language = ?", v)
}

func KnownQueryMimetype(v string) squirrel.Sqlizer {
	if stringsx.Blank(v) {
		return squirrelx.Noop{}
	}

	return squirrel.Expr("library_known_media.mimetype = ?", v)
}

func KnownQuerySource(sources ...string) squirrel.Sqlizer {
	return squirrelx.In("library_known_media.source", sources...)
}

func KnownQueryExcludeSource(sources ...string) squirrel.Sqlizer {
	return squirrelx.NotIn("library_known_media.source", sources...)
}

func KnownQueryDetectLanguage(v string) squirrel.Sqlizer {
	min, max := unicodex.LowHi(unicodex.ISO639_1(v))
	if langx.FirstNonZero(min, max) == 0 {
		return squirrelx.Noop{}
	}

	return squirrelx.Between("unicode(library_known_media.auto_description)", min, max)
}

func KnownQueryReleased(r timex.Range) squirrel.Sqlizer {
	return squirrelx.Between("library_known_media.released", ducktype.NewNullTime(r.Start), ducktype.NewNullTime(r.End))
}

func KnownQueryWithPoster() squirrel.Sqlizer {
	return squirrel.Expr("(library_known_media.poster_path != '' OR library_known_media.backdrop_path != '')")
}

func KnownQuerySimilarity(q string, cutoff float32) squirrel.Sqlizer {
	return squirrel.Expr("((jaro_winkler_similarity(library_known_media.title, ?, ?) + jaro_similarity(library_known_media.title, ?, ?)) / 2) > 0.5", q, cutoff, q, cutoff)
}

func KnownSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(KnownScannerStaticColumns)...).From("library_known_media")
}

func DetectKnownMedia(ctx context.Context, db sqlx.Queryer, mimecat string, query string) (k Known, err error) {
	k = Unknown()

	if err := KnownBestMatch(ctx, db, mimecat, query, 0.7).Scan(&k); sqlx.IgnoreNoRows(err) != nil {
		return k, errorsx.Wrap(err, "unable to score")
	}

	return k, nil
}
