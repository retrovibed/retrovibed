package library

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/localex"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

type RecommendationSource string

const (
	RecommendationSourceRandom       = "random"
	RecommendationSourceGenerative   = "generative"
	RecommendationSourceDiscovered   = "discovered"
	RecommendationSourceSearchPlugin = "searchplugin"
)

const (
	RecommendationTTL = 30 * 24 * time.Hour
)

type RecommendationOption func(*Recommendation)

// backfills metadata information from a known record.
// used when a recommendation source has identified itself.
// and we want to use the data.
func RecommendationOptionAutoKnown(v Known) RecommendationOption {
	return func(r *Recommendation) {
		r.KnownMediaID = v.UID
		r.Language = localex.FirstDefined(r.Language, v.OriginalLanguage)
		r.Adult = langx.FirstNonZero(r.Adult, v.Adult)
		r.Image = langx.FirstNonZero(v.PosterPath, v.BackdropPath)
		r.Title = v.Title
		r.Overview = v.Overview
		r.Popularity = v.Popularity
		r.Released = v.Released
	}
}

func RecommendationOptionSourceDiscovered(r *Recommendation) {
	r.Source = md5x.String(RecommendationSourceDiscovered)
}

func RecommendationOptionSourceRandom(r *Recommendation) {
	r.Source = md5x.String(RecommendationSourceRandom)
}

func RecommendationOptionSourceSearchPlugin(r *Recommendation) {
	r.Source = md5x.String(RecommendationSourceSearchPlugin)
}

func RecommendationOptionRecommendationTTL(r *Recommendation) {
	r.TombstoneAt = time.Now().Add(RecommendationTTL)
}

func RecommendationSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) RecommendationScanner {
	return NewRecommendationScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func RecommendationSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(RecommendationScannerStaticColumns)...).
		From("library_recommendations")
}

func RecommendationOptionTestDefaults(r *Recommendation) {
	r.ID = uuid.Nil.String()
	r.Source = md5x.String(RecommendationSourceRandom)
	r.ContentID = uuid.Nil.String()
	r.KnownMediaID = uuid.Nil.String()
	r.TombstoneAt = timex.Inf()
	r.Mimetype = mimex.Binary
}

func RecommendationSourceString(uid string) string {
	switch uid {
	case md5x.String(RecommendationSourceRandom):
		return RecommendationSourceRandom
	case md5x.String(RecommendationSourceGenerative):
		return RecommendationSourceGenerative
	case md5x.String(RecommendationSourceDiscovered):
		return RecommendationSourceDiscovered
	case md5x.String(RecommendationSourceSearchPlugin):
		return RecommendationSourceSearchPlugin
	default:
		return ""
	}
}

func RecommendationOptionID(id string) RecommendationOption {
	return func(r *Recommendation) {
		r.ID = id
	}
}

func RecommendationOptionContentID(kid string) RecommendationOption {
	return func(r *Recommendation) {
		r.ContentID = errorsx.Must(uuid.FromString(kid)).String()
	}
}

func RecommendationFromRandomKnown(ctx context.Context, q sqlx.Queryer, mimetype string, lang string, adult bool) (rec Recommendation, err error) {
	var known Known

	scanner := KnownSearch(ctx, q, KnownSearchBuilder().Where(squirrel.And{
		KnownQueryWithPoster(),
		KnownQueryExplicit(adult),
		KnownQueryLanguage(lang),
		KnownQueryMimetype(mimetype),
	}).OrderBy("random()").Limit(1))
	if known, err = sqlx.ScanOne(scanner); err != nil {
		return rec, err
	}

	if err = RecommendationInsertWithDefaults(ctx, q, Recommendation{
		Source:       md5x.String(RecommendationSourceRandom),
		ContentID:    known.UID,
		KnownMediaID: known.UID,
		TombstoneAt:  time.Now().Add(RecommendationTTL),
		Mimetype:     known.Mimetype,
		Adult:        known.Adult,
		Title:        known.Title,
		Overview:     known.Overview,
		Image:        stringsx.FirstNonBlank(known.PosterPath, known.BackdropPath),
		Popularity:   known.Popularity,
		Released:     known.Released,
	}).Scan(&rec); err != nil {
		return rec, err
	}

	return rec, nil
}

func RecommendationLastGeneratedAt(ctx context.Context, q sqlx.Queryer, source string) (time.Time, error) {
	return sqlx.Timestamp(ctx, q,
		`SELECT COALESCE(MAX(updated_at), '-infinity')::TIMESTAMPTZ FROM library_recommendations WHERE source = $1`,
		md5x.String(source),
	)
}

func RecommendationQueryNotTombstoned() squirrel.Sqlizer {
	return squirrel.Expr("library_recommendations.tombstone_at > NOW()")
}

func RecommendationQueryByContentID(kid string) squirrel.Sqlizer {
	return squirrel.Eq{"library_recommendations.content_id": kid}
}

func RecommendationQueryMimetype(v string) squirrel.Sqlizer {
	if stringsx.Blank(v) {
		return squirrelx.Noop{}
	}

	return squirrel.Expr("library_recommendations.mimetype = ?", v)
}
