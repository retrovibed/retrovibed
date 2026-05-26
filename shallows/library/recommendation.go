package library

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

type RecommendationSource string

const (
	RecommendationSourceRandom     = "random"
	RecommendationSourceGenerative = "generative"
)

func RecommendationKnownSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) RecommendationKnownScanner {
	return NewRecommendationKnownScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func RecommendationKnownSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(RecommendationScannerStaticColumns, KnownScannerStaticColumns)...).
		From("library_recommendations").
		InnerJoin("library_known_media ON library_known_media.uid = library_recommendations.known_media_id").
		Where(squirrel.Expr("'t'"))
}

func RecommendationOptionTestDefaults(r *Recommendation) {
	r.ID = uuid.Nil.String()
	r.Source = md5x.String(RecommendationSourceRandom)
	r.KnownMediaID = uuid.Nil.String()
	r.TombstoneAt = timex.Inf()
	r.Mimetype = mimex.Binary
}

func RecommendationOptionID(id string) func(*Recommendation) {
	return func(r *Recommendation) {
		r.ID = id
	}
}

func RecommendationOptionKnownMediaID(kid string) func(*Recommendation) {
	return func(r *Recommendation) {
		r.KnownMediaID = errorsx.Must(uuid.FromString(kid)).String()
	}
}

func RecommendationFromRandomKnown(ctx context.Context, q sqlx.Queryer, mimetype string) (rec Recommendation, err error) {
	var known Known

	scanner := KnownSearch(ctx, q, KnownSearchBuilder().Where(squirrel.And{
		KnownQueryExplicit(false),
		KnownQueryWithPoster(),
		KnownQueryMimetype(mimetype),
	}).OrderBy("random()").Limit(1))
	if known, err = sqlx.ScanOne(scanner); err != nil {
		return rec, err
	}

	if err = RecommendationInsertWithDefaults(ctx, q, Recommendation{
		ID:           uuid.Must(uuid.NewV7()).String(),
		Source:       md5x.String(RecommendationSourceRandom),
		KnownMediaID: known.UID,
		TombstoneAt:  time.Now().Add(30 * 24 * time.Hour),
		Mimetype:     known.Mimetype,
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

func RecommendationQueryByKnownMediaID(kid string) squirrel.Sqlizer {
	return squirrel.Eq{"library_recommendations.known_media_id": kid}
}

func RecommendationQueryMimetype(v string) squirrel.Sqlizer {
	if stringsx.Blank(v) {
		return squirrelx.Noop{}
	}

	return squirrel.Expr("library_recommendations.mimetype = ?", v)
}
