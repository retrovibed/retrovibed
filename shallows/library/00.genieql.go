//go:build genieql.generate
// +build genieql.generate

package library

import (
	"context"
	"time"

	genieql "github.com/james-lawrence/genieql/ginterp"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

func Metadata(gql genieql.Structure) {
	gql.From(
		gql.Table("library_metadata"),
	)
}

func MetadataScanner(gql genieql.Scanner, pattern func(i Metadata)) {
	gql.ColumnNamePrefix("library_metadata.")
}

func MetadataInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Metadata) NewMetadataScannerStaticRow,
) {
	gql.Into("library_metadata").Default("created_at", "updated_at", "tombstoned_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, tombstoned_at = DEFAULT, auto_description = EXCLUDED.auto_description, description = EXCLUDED.description, archive_id = CASE WHEN archive_id IN ('ffffffff-ffff-ffff-ffff-ffffffffffff', '00000000-0000-0000-0000-000000000000') THEN EXCLUDED.archive_id ELSE archive_id END, known_media_id = CASE WHEN known_media_id IN ('ffffffff-ffff-ffff-ffff-ffffffffffff', '00000000-0000-0000-0000-000000000000') THEN EXCLUDED.known_media_id ELSE known_media_id END")
}

// directories are created here rather than through MetadataInsertWithDefaults so the
// library's insert never has to reason about them. a directory carries no content, so the
// only thing a conflict can mean is a rename.
func DirectoryUpsert(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Metadata) NewMetadataScannerStaticRow,
) {
	gql.Into("library_metadata").Default("created_at", "updated_at", "tombstoned_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, description = EXCLUDED.description, auto_description = EXCLUDED.auto_description")
}

func MetadataDeleteByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM library_metadata WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataArchivedByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id, aid string, quota uint64) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE library_metadata SET archive_id = {aid}, quota_usage = {quota} WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataHideByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE library_metadata SET hidden_at = NOW(), initiated_at = 'infinity' WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataUpdateAutodescriptionByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, autodescription string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE library_metadata SET updated_at = NOW(), auto_description = {autodescription} WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataUpdateDescriptionByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, description string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE library_metadata SET updated_at = NOW(), description = {description} WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

// reindexed media metadata, resets the known media id.
func MetadataUpdateReindexByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, description string, autodescription string, kid string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE library_metadata SET updated_at = NOW(), description = {description}, auto_description = {autodescription}, known_media_id = {kid} WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataFindByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + MetadataScannerStaticColumns + ` FROM library_metadata WHERE "id" = {id}`)
}

func MetadataTombstoneByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE library_metadata SET tombstoned_at = NOW() WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

// tombstones the folder and everything below it. the migration adds directory_id without a
// cascade, so tombstoning a folder alone leaves its children pointing at a row
// NewTombstonedCleanup is about to hard delete, listed by nothing and counted by
// MetadataDiskStorageUsage.
func MetadataTombstoneSubtreeByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStatic,
) {
	gql = gql.Query(`WITH RECURSIVE subtree(id) AS (SELECT id FROM library_metadata WHERE id = {id} UNION ALL SELECT lm.id FROM library_metadata AS lm INNER JOIN subtree ON lm.directory_id = subtree.id) UPDATE library_metadata SET tombstoned_at = NOW() WHERE library_metadata.id IN (SELECT id FROM subtree) RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataTombstoneByTorrentID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, tid string) NewMetadataScannerStatic,
) {
	gql = gql.Query(`UPDATE library_metadata SET tombstoned_at = NOW() WHERE "torrent_id" = {tid} RETURNING ` + MetadataScannerStaticColumns)
}

func RecentSessionDeleteByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewRecentSessionScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM library_recent_sessions WHERE "id" = {id} RETURNING ` + RecentSessionScannerStaticColumns)
}

func MetadataDeleteByTorrentID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, tid string) NewMetadataScannerStatic,
) {
	gql = gql.Query(`DELETE FROM library_metadata WHERE "torrent_id" = {tid} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataFindByDescription(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, desc string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + MetadataScannerStaticColumns + ` FROM library_metadata WHERE "description" = {desc}`)
}

func MetadataAssociateTorrent(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, desc, tid string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE library_metadata SET torrent_id = {tid} WHERE "description" = {desc} AND torrent_id = '00000000-0000-0000-0000-000000000000' RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataUpdate(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, md Metadata) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE library_metadata SET description = {md.Description}, known_media_id = {md.KnownMediaID}, archive_id = {md.ArchiveID} WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

// a parent drawn from the row's own subtree builds a directory_id cycle, and every recursive
// descent here then runs until the process is killed. such a move matches no row. the root
// sentinel is never itself a row, so moving to the top level always matches.
func MetadataMoveByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id, directory string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`WITH RECURSIVE subtree(id) AS (SELECT id FROM library_metadata WHERE id = {id} UNION ALL SELECT lm.id FROM library_metadata AS lm INNER JOIN subtree ON lm.directory_id = subtree.id) UPDATE library_metadata SET updated_at = NOW(), directory_id = {directory} WHERE library_metadata.id = {id} AND {directory} NOT IN (SELECT id FROM subtree) RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataSetTorrentID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id, tid string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE library_metadata SET torrent_id = {tid} WHERE "id" = {id} AND torrent_id = '00000000-0000-0000-0000-000000000000' RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataTransferKnownMediaIDFromTorrent(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, ts time.Time) NewMetadataScannerStatic,
) {
	gql = gql.Query(`UPDATE library_metadata SET updated_at = NOW(), known_media_id = t.known_media_id FROM torrents_metadata AS t WHERE t.id = library_metadata.torrent_id AND t."updated_at" >= {ts} AND library_metadata.known_media_id = 'ffffffff-ffff-ffff-ffff-ffffffffffff' AND t.known_media_id NOT IN ('ffffffff-ffff-ffff-ffff-ffffffffffff', '00000000-0000-0000-0000-000000000000') RETURNING ` + MetadataScannerStaticColumns)
}

// used to sync known media idea from a known torrent.metadata to every media with that torrent.metadata.
func MetadataSyncKnownMediaIDFromTorrent(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, tid string) NewMetadataScannerStatic,
) {
	gql = gql.Query(`UPDATE library_metadata SET updated_at = NOW(), known_media_id = torrents_metadata.known_media_id FROM torrents_metadata WHERE torrents_metadata."id" = {tid} AND torrents_metadata."id" = library_metadata.torrent_id RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataForTorrentArchiveRetrieval(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, infohash []byte, offset uint64, length uint64) NewMetadataScannerStatic,
) {
	gql = gql.Query(`SELECT ` + MetadataScannerStaticColumns + ` FROM library_metadata INNER JOIN torrents_metadata AS tmd ON library_metadata.torrent_id = tmd.id WHERE to_hex(tmd.infohash) = to_hex({infohash}) AND {offset} BETWEEN disk_offset AND library_metadata.bytes AND library_metadata.archive_id NOT IN ('ffffffff-ffff-ffff-ffff-ffffffffffff', '00000000-0000-0000-0000-000000000000')`)
}

// the folder and every descendant. MetadataSearchBuilder returns a single table
// SelectBuilder and cannot express the recursion.
func MetadataSubtreeByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStatic,
) {
	gql = gql.Query(`WITH RECURSIVE subtree(id) AS (SELECT id FROM library_metadata WHERE id = {id} UNION ALL SELECT lm.id FROM library_metadata AS lm INNER JOIN subtree ON lm.directory_id = subtree.id) SELECT ` + MetadataScannerStaticColumns + ` FROM library_metadata INNER JOIN subtree ON library_metadata.id = subtree.id`)
}

// the row and its ancestors, root first, which is the order a breadcrumb renders in.
// recursion runs upward on the parent and terminates at the root sentinel, which matches
// no row.
func MetadataAncestorsByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStatic,
) {
	gql = gql.Query(`WITH RECURSIVE ancestors(id, directory_id, depth) AS (SELECT id, directory_id, 0 FROM library_metadata WHERE id = {id} UNION ALL SELECT lm.id, lm.directory_id, ancestors.depth + 1 FROM library_metadata AS lm INNER JOIN ancestors ON ancestors.directory_id = lm.id) SELECT ` + MetadataScannerStaticColumns + ` FROM library_metadata INNER JOIN ancestors ON library_metadata.id = ancestors.id ORDER BY ancestors.depth DESC`)
}

// Known, ScoredScanner, and their supporting genieql declarations live in
// 01.genieql.cache.go (build tag retrovibed.cache) since library_known_media
// lives in cache.db, not meta.db.

func RecentSession(gql genieql.Structure) {
	gql.From(
		gql.Table("library_recent_sessions"),
	)
}

func RecentSessionScanner(gql genieql.Scanner, pattern func(i RecentSession)) {
	gql.ColumnNamePrefix("library_recent_sessions.")
}

func RecentSessionLibraryScanner(gql genieql.Scanner, pattern func(i RecentSession, md Metadata)) {}

func RecentSessionInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a RecentSession) NewRecentSessionScannerStaticRow,
) {
	gql.Into("library_recent_sessions").Default("created_at", "updated_at", "last_played_at").Conflict("ON CONFLICT (profile_id, id) DO UPDATE SET position = EXCLUDED.position, duration = EXCLUDED.duration, query = EXCLUDED.query, updated_at = DEFAULT, last_played_at = DEFAULT")
}

func WatchHistory(gql genieql.Structure) {
	gql.From(
		gql.Table("library_watch_history"),
	)
}

func WatchHistoryScanner(gql genieql.Scanner, pattern func(i WatchHistory)) {
	gql.ColumnNamePrefix("library_watch_history.")
}

func WatchHistoryInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a WatchHistory) NewWatchHistoryScannerStaticRow,
) {
	gql.Into("library_watch_history").Default("created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET duration = EXCLUDED.duration, updated_at = DEFAULT")
}

func Recommendation(gql genieql.Structure) {
	gql.From(
		gql.Table("library_recommendations"),
	)
}

func RecommendationScanner(gql genieql.Scanner, pattern func(i Recommendation)) {
	gql.ColumnNamePrefix("library_recommendations.")
}

func RecommendationInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Recommendation) NewRecommendationScannerStaticRow,
) {
	gql.Into("library_recommendations").Default("id", "created_at", "updated_at").Conflict("ON CONFLICT (content_id) DO UPDATE SET recommendations = recommendations + 1, updated_at = DEFAULT")
}

func RecommendationFindByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewRecommendationScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + RecommendationScannerStaticColumns + ` FROM library_recommendations WHERE "id" = {id}`)
}

func RecommendationFindByContentID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, kid string) NewRecommendationScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + RecommendationScannerStaticColumns + ` FROM library_recommendations WHERE "content_id" = {kid}`)
}

func RecommendationDeleteByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewRecommendationScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM library_recommendations WHERE "id" = {id} RETURNING ` + RecommendationScannerStaticColumns)
}

func RecommendationDeleteTombstoned(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewRecommendationScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM library_recommendations WHERE "tombstoned_at" < NOW() RETURNING ` + RecommendationScannerStaticColumns)
}
