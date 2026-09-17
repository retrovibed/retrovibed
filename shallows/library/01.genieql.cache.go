//go:build genieql.generate && retrovibed.cache
// +build genieql.generate,retrovibed.cache

package library

import (
	"context"

	genieql "github.com/james-lawrence/genieql/ginterp"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

func ScoredScanner(gql genieql.Scanner, pattern func(relevance float64)) {
}

func Known(gql genieql.Structure) {
	gql.From(
		gql.Table("library_known_media"),
	)
}

func KnownScanner(gql genieql.Scanner, pattern func(i Known)) {
	gql.ColumnNamePrefix("library_known_media.")
}

func KnownInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Known) NewKnownScannerStaticRow,
) {
	gql.Into("library_known_media").Default("created_at", "tombstoned_at").Conflict(`ON CONFLICT (uid) DO UPDATE SET title = EXCLUDED.title, original_language = EXCLUDED.original_language, original_title = EXCLUDED.original_title, popularity = EXCLUDED.popularity, overview = EXCLUDED.overview, source = EXCLUDED.source, poster_path = EXCLUDED.poster_path, backdrop_path = EXCLUDED.backdrop_path, mimetype = EXCLUDED.mimetype, "collation" = EXCLUDED."collation", subtitle = EXCLUDED.subtitle, parent_uid = EXCLUDED.parent_uid, released = EXCLUDED.released, adult = EXCLUDED.adult, md5 = EXCLUDED.md5, md5_lower = EXCLUDED.md5_lower, auto_description = EXCLUDED.auto_description, duplicates = duplicates + 1`)
}

// KnownInsertWithDefaultsTOFU writes a discovery-pipeline placeholder row.
// tombstoned_at is bound (not defaulted) so the caller can stamp a TTL; the
// conflict clause refreshes it on every rediscovery, giving TOFU rows a
// sliding expiry rather than one fixed at first creation.
func KnownInsertWithDefaultsTOFU(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Known) NewKnownScannerStaticRow,
) {
	gql.Into("library_known_media").Default("created_at").Conflict("ON CONFLICT (uid) DO UPDATE SET uid = EXCLUDED.uid, tombstoned_at = EXCLUDED.tombstoned_at")
}

func KnownBatchInsertWithDefaults(
	gql genieql.InsertBatch,
	pattern func(ctx context.Context, q sqlx.Queryer, p Known) NewKnownScannerStatic,
) {
	gql.Into("library_known_media").Batch(64).Default("created_at", "tombstoned_at").Conflict(`ON CONFLICT (uid) DO UPDATE SET title = EXCLUDED.title, original_language = EXCLUDED.original_language, original_title = EXCLUDED.original_title, popularity = EXCLUDED.popularity, overview = EXCLUDED.overview, source = EXCLUDED.source, poster_path = EXCLUDED.poster_path, backdrop_path = EXCLUDED.backdrop_path, mimetype = EXCLUDED.mimetype, "collation" = EXCLUDED."collation", subtitle = EXCLUDED.subtitle, parent_uid = EXCLUDED.parent_uid, released = EXCLUDED.released, adult = EXCLUDED.adult, md5 = EXCLUDED.md5, md5_lower = EXCLUDED.md5_lower, auto_description = EXCLUDED.auto_description, duplicates = duplicates + 1`)
}

func KnownFindByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewKnownScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + KnownScannerStaticColumns + ` FROM library_known_media WHERE "uid" = {id}`)
}

func KnownFindByMd5(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, md5 string) NewKnownScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + KnownScannerStaticColumns + ` FROM library_known_media WHERE "md5" = {md5}`)
}

func KnownFindByLastCreated(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewKnownScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + KnownScannerStaticColumns + ` FROM library_known_media ORDER BY created_at DESC LIMIT 1`)
}

func KnownFindRandom(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewKnownScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + KnownScannerStaticColumns + ` FROM library_known_media WHERE NOT adult AND (poster_path <> '' OR backdrop_path <> '') USING SAMPLE 1`)
}

func KnownScoreByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, uid string, terms string, cutoff float32) NewScoredScannerStaticRow,
) {
	gql = gql.Query(`SELECT (jaro_winkler_similarity(title, {terms}, {cutoff}) + jaro_similarity(title, {terms}, {cutoff})) / 2 AS relevance FROM library_known_media WHERE uid = {uid}`)
}

func KnownBestMatch(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, mime string, terms string, cutoff float32) NewKnownScannerStaticRow,
) {
	// TEMPORARY: parent_uid = nil uuid excludes episode rows (which share their
	// show's title) from best-match scoring until real episode-aware matching
	// exists; without it, a show and all its episodes tie on title similarity.
	gql = gql.Query(`WITH scored AS (SELECT uid, {terms} as q, (jaro_winkler_similarity(title, q, {cutoff}) + jaro_similarity(title, q, {cutoff})) / 2 AS relevance FROM library_known_media WHERE NOT adult AND ({mime} = '' OR mimetype = {mime}) AND parent_uid = '00000000-0000-0000-0000-000000000000' ORDER BY relevance DESC) SELECT ` + KnownScannerStaticColumns + ` FROM library_known_media INNER JOIN scored ON library_known_media.uid = scored.uid WHERE scored.relevance > {cutoff} ORDER BY scored.relevance DESC`)
}

func KnownTombstoneByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewKnownScannerStaticRow,
) {
	gql = gql.Query(`UPDATE library_known_media SET tombstoned_at = NOW() WHERE "uid" = {id} RETURNING ` + KnownScannerStaticColumns)
}

func KnownDeleteTombstoned(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewKnownScannerStatic,
) {
	gql = gql.Query(`DELETE FROM library_known_media WHERE "tombstoned_at" < NOW() RETURNING ` + KnownScannerStaticColumns)
}
