//go:build genieql.generate
// +build genieql.generate

package acoustics

import (
	"context"

	genieql "github.com/james-lawrence/genieql/ginterp"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

func LibraryMetadata(gql genieql.Structure) {
	gql.From(gql.Table("library_metadata"))
}

func LibraryMetadataScanner(gql genieql.Scanner, pattern func(i LibraryMetadata)) {
	gql.ColumnNamePrefix("library_metadata.")
}

func AudioFeatures(gql genieql.Structure) {
	gql.From(gql.Table("audio_features"))
}

func AudioFeaturesScanner(gql genieql.Scanner, pattern func(i AudioFeatures)) {
	gql.ColumnNamePrefix("audio_features.")
}

func AudioFeaturesInsert(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a AudioFeatures) NewAudioFeaturesScannerStaticRow,
) {
	gql.Into("audio_features").Default("indexed_at")
}

func AudioFeaturesDeleteByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewAudioFeaturesScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM audio_features WHERE "media_id" = {id} RETURNING ` + AudioFeaturesScannerStaticColumns)
}

func AudioFeaturesFindByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewAudioFeaturesScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + AudioFeaturesScannerStaticColumns + ` FROM audio_features WHERE "media_id" = {id}`)
}

func AudioFeaturesSimilarByVec(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, vec []float32, exclude []string, n int) NewAudioFeaturesScannerStatic,
) {
	gql = gql.Query(`SELECT ` + AudioFeaturesScannerStaticColumns + ` FROM audio_features WHERE NOT failed AND NOT list_contains({exclude}, media_id::VARCHAR) ORDER BY array_cosine_distance(features, {vec}::FLOAT[128]) LIMIT {n}`)
}

func CountScanner(gql genieql.Scanner, pattern func(count int64)) {}

func AudioFeaturesCountByVersion(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, version uint32) NewCountScannerStaticRow,
) {
	gql = gql.Query(`SELECT COUNT(*) AS count FROM audio_features WHERE "stats_version" = {version} AND NOT failed`)
}

func AudioFeaturesUnindexedMediaIDs(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewLibraryMetadataScannerStatic,
) {
	gql = gql.Query(`SELECT ` + LibraryMetadataScannerStaticColumns + ` FROM library_metadata LEFT JOIN audio_features af ON af.media_id = library_metadata.id WHERE library_metadata.mimetype LIKE 'audio/%' AND library_metadata.tombstoned_at = 'infinity' AND af.media_id IS NULL`)
}
