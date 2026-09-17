package library

import (
	"context"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/internal/ducktype"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
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

func KnownOptionCollation(v uint32) func(*Known) {
	return func(t *Known) {
		t.Collation = v
	}
}

func KnownOptionSubtitle(v string) func(*Known) {
	return func(t *Known) {
		t.Subtitle = v
	}
}

func KnownOptionParentUID(v string) func(*Known) {
	return func(t *Known) {
		t.ParentUID = v
	}
}

func KnownOptionAutoDescription(t *Known) {
	t.AutoDescription = stringsx.Join("\n", t.Title, t.OriginalTitle, t.Overview)
}

// KnownCollationSpecialsSeason marks TMDB "Specials" (season_number == 0)
// in the high 16 bits of Collation, instead of 0, so a specials episode
// never collides with Collation == 0 (the standalone/overall item marker).
const KnownCollationSpecialsSeason uint16 = 0xFFFF

// KnownCollationEpisode packs a season/episode pair into a single ordering
// value: high 16 bits season, low 16 bits episode. Collation == 0 is
// reserved for the standalone/overall item (e.g. a TV show's own row).
// TMDB's "Specials" season (season_number == 0) is remapped to
// KnownCollationSpecialsSeason so it doesn't collide with that marker.
func KnownCollationEpisode(season, episode uint16) uint32 {
	if season == 0 {
		season = KnownCollationSpecialsSeason
	}
	return uint32(season)<<16 | uint32(episode)
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
	t.ParentUID = uuid.Nil.String()
	t.OriginalLanguage = userx.LocaleLanguage()
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
	return squirrel.Expr("cache.library_known_media.uid > ?", uid)
}

// KnownQueryExplicit toggles whether adult content is allowed in results.
// allow=false restricts results to non-adult content; allow=true permits
// both adult and non-adult content (it does not restrict to adult-only).
func KnownQueryExplicit(allow bool) squirrel.Sqlizer {
	if allow {
		return squirrelx.Noop{}
	}

	return squirrel.Expr("cache.library_known_media.adult = ?", false)
}

func KnownQueryUID(ids ...string) squirrel.Sqlizer {
	return squirrelx.In("cache.library_known_media.uid", ids...)
}

func KnownQueryParentUID(uid string) squirrel.Sqlizer {
	return squirrel.Expr("cache.library_known_media.parent_uid = ?", uid)
}

func KnownQueryLanguage(v string) squirrel.Sqlizer {
	if stringsx.Blank(v) {
		return squirrelx.Noop{}
	}

	return squirrel.Expr("cache.library_known_media.original_language = ?", v)
}

func KnownQueryMimetype(v string) squirrel.Sqlizer {
	if stringsx.Blank(v) {
		return squirrelx.Noop{}
	}

	return squirrel.Expr("cache.library_known_media.mimetype = ?", v)
}

func KnownQuerySource(sources ...string) squirrel.Sqlizer {
	return squirrelx.In("cache.library_known_media.source", sources...)
}

func KnownQueryExcludeSource(sources ...string) squirrel.Sqlizer {
	return squirrelx.NotIn("cache.library_known_media.source", sources...)
}

func KnownQueryDetectLanguage(v string) squirrel.Sqlizer {
	min, max := unicodex.LowHi(unicodex.ISO639_1(v))
	if langx.FirstNonZero(min, max) == 0 {
		return squirrelx.Noop{}
	}

	return squirrelx.Between("unicode(cache.library_known_media.auto_description)", min, max)
}

func KnownQueryReleased(r timex.Range) squirrel.Sqlizer {
	return squirrelx.Between("cache.library_known_media.released", ducktype.NewNullTime(r.Start), ducktype.NewNullTime(r.End))
}

func KnownQueryWithPoster() squirrel.Sqlizer {
	return squirrel.Expr("(cache.library_known_media.poster_path != '' OR cache.library_known_media.backdrop_path != '')")
}

// KnownMatchCutoff is the default minimum combined jaro-winkler/jaro
// similarity required to accept a known-media match. Jaro-family metrics
// degrade on long strings, so this is set well above the point where two
// unrelated titles can score deceptively high by chance.
const KnownMatchCutoff float32 = 0.85

func KnownQuerySimilarity(q string, cutoff float32) squirrel.Sqlizer {
	return squirrel.Expr("((jaro_winkler_similarity(cache.library_known_media.title, ?, ?) + jaro_similarity(cache.library_known_media.title, ?, ?)) / 2) > ?", q, cutoff, q, cutoff, cutoff)
}

func KnownSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(KnownScannerStaticColumns)...).From("cache.library_known_media")
}

func KnownQueryNotTombstoned() squirrel.Sqlizer {
	return squirrel.Expr("cache.library_known_media.tombstoned_at = 'infinity'")
}

func KnownQueryTombstoned() squirrel.Sqlizer {
	return squirrel.Expr("cache.library_known_media.tombstoned_at < 'infinity'")
}

func DetectKnownMedia(ctx context.Context, db sqlx.Queryer, mimecat string, query string, similarity float32) (k Known, err error) {
	k = Unknown()

	if err := KnownBestMatch(ctx, db, mimecat, query, similarity).Scan(&k); sqlx.IgnoreNoRows(err) != nil {
		return k, errorsx.Wrap(err, "unable to score")
	}

	return k, nil
}

// NewKnownMediaTombstonedCleanup purges known-media catalog rows (e.g. stale
// TOFU placeholders from the discovery pipeline) once tombstoned.
func NewKnownMediaTombstonedCleanup(ctx context.Context, q sqlx.Queryer) error {
	log.Println("known media tombstoned cleanup initiated")
	defer log.Println("known media tombstoned cleanup completed")

	return sqlx.Discard(sqlx.Scan(KnownDeleteTombstoned(ctx, q)))
}
