//go:build genieql.generate
// +build genieql.generate

package tracking

import (
	"context"
	"time"

	genieql "github.com/james-lawrence/genieql/ginterp"
	"github.com/retrovibed/retrovibed/internal/sqlx"
)

func Metadata(gql genieql.Structure) {
	gql.From(
		gql.Table("torrents_metadata"),
	)
}

func MetadataScanner(gql genieql.Scanner, pattern func(i Metadata)) {
	gql.ColumnNamePrefix("torrents_metadata.")
}

func MetadataInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Metadata) NewMetadataScannerStaticRow,
) {
	gql.Into("torrents_metadata").Default("created_at", "updated_at", "paused_at", "verify_at", "imported_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, expires_at = EXCLUDED.expires_at, tracker = EXCLUDED.tracker")
}

func MetadataInsertOrUpdate(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Metadata) NewMetadataScannerStaticRow,
) {
	gql.Into("torrents_metadata").Default("created_at", "updated_at", "imported_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, paused_at = EXCLUDED.paused_at, verify_at = EXCLUDED.verify_at")
}

func MetadataBatchInsertWithDefaults(
	gql genieql.InsertBatch,
	pattern func(ctx context.Context, q sqlx.Queryer, p Metadata) NewMetadataScannerStatic,
) {
	gql.Into("torrents_metadata").Batch(10).Default("created_at", "updated_at", "hidden_at", "initiated_at", "paused_at", "downloaded", "next_announce_at")
}

func MetadataResetByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_metadata SET updated_at = DEFAULT, initiated_at = DEFAULT, completed_at = DEFAULT, paused_at = DEFAULT, next_announce_at = DEFAULT, seeding = DEFAULT, downloaded = DEFAULT, uploaded = DEFAULT, peers = DEFAULT WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataTombstoneByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_metadata SET tombstoned_at = NOW(), initiated_at = 'infinity', seeding = 'f' WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataDeleteByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM torrents_metadata WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataFindByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + MetadataScannerStaticColumns + ` FROM torrents_metadata WHERE "id" = {id}`)
}

func MetadataFindByInfohash(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + MetadataScannerStaticColumns + ` FROM torrents_metadata WHERE "infohash" = unhex({id})`)
}

func MetadataPausedByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_metadata SET paused_at = NOW(), initiated_at = 'infinity' WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataDownloadByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_metadata SET paused_at = 'infinity', initiated_at = NOW() WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataAutoDownloadByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_metadata SET initiated_at = NOW() WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataProgressByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, peers uint16, bytes uint64, downloaded uint64) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_metadata SET updated_at = NOW(), bytes = {bytes}, downloaded = {downloaded}, peers = {peers}, seeding = (bytes == {downloaded}) WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataImportedByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_metadata SET updated_at = NOW(), imported_at = NOW() WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataVerifyByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, peers uint16, downloaded uint64) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_metadata SET updated_at = NOW(), verify_at = NOW(), downloaded = {downloaded}, peers = {peers}, seeding = (bytes == {downloaded}) WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataVerifiedByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, peers uint16, downloaded uint64) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_metadata SET updated_at = NOW(), verify_at = DEFAULT, downloaded = {downloaded}, peers = {peers}, seeding = (bytes == {downloaded}) WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataCompleteByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, peers uint16, bytes uint64, downloaded uint64, uploaded uint64) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_metadata SET updated_at = NOW(), completed_at = NOW(), bytes = {bytes}, downloaded = {downloaded}, uploaded = {uploaded}, peers = {peers}, seeding = (bytes == {downloaded}), verify_at = 'infinity' WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataUploadedByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id []byte, uploaded uint64) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_metadata SET updated_at = NOW(), uploaded = (uploaded + {uploaded}) WHERE "infohash" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataAnnounced(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, nextts time.Time) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_metadata SET updated_at = NOW(), next_announce_at = {nextts} WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataDisableAnnounced(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_metadata SET updated_at = NOW(), next_announce_at = 'infinity' WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataAssignKnownMediaID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, kid string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_metadata SET updated_at = NOW(), known_media_id = {kid} WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func Peer(gql genieql.Structure) {
	gql.From(
		gql.Table("torrents_peers"),
	)
}

func PeerScanner(gql genieql.Scanner, pattern func(i Peer)) {
	gql.ColumnNamePrefix("torrents_peers.")
}

func PeerInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Peer) NewPeerScannerStaticRow,
) {
	gql.Into("torrents_peers").Default("created_at", "updated_at", "next_check").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, ip = EXCLUDED.ip, port = EXCLUDED.port, description = EXCLUDED.description, bep51_available = EXCLUDED.bep51_available, ddisc = EXCLUDED.ddisc, ddisc_partition = EXCLUDED.ddisc_partition, ddisc_syncoffset = EXCLUDED.ddisc_syncoffset")
}

func PeerUpdateDdisc(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, a Peer) NewPeerScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_peers SET updated_at = NOW(), description = {a.Description}, ddisc = {a.Ddisc} WHERE "id" = {id} RETURNING ` + PeerScannerStaticColumns)
}

func PeerFindByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewPeerScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + PeerScannerStaticColumns + ` FROM torrents_peers WHERE "id" = {id}`)
}

func PeerDeleteByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewPeerScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM torrents_peers WHERE "id" = {id} RETURNING ` + PeerScannerStaticColumns)
}

func PeerMarkNextCheck(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Peer) NewPeerScannerStaticRow,
) {
	gql.Into("torrents_peers").Default("created_at", "updated_at", "next_check").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = NOW(), next_check = NOW() + to_seconds(EXCLUDED.bep51_ttl), tombstoned_at = EXCLUDED.tombstoned_at")
}

func PeerClearTombstoned(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, ts time.Time) NewPeerScannerStatic,
) {
	gql = gql.Query(`DELETE FROM torrents_peers WHERE "tombstoned_at" <= {ts} RETURNING ` + PeerScannerStaticColumns)
}

func UnknownHash(gql genieql.Structure) {
	gql.From(
		gql.Table("torrents_unknown_infohashes"),
	)
}

func UnknownHashScanner(gql genieql.Scanner, pattern func(i UnknownHash)) {
	gql.ColumnNamePrefix("torrents_unknown_infohashes.")
}

func UnknownHashInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a UnknownHash) NewUnknownHashScannerStaticRow,
) {
	gql.Into("torrents_unknown_infohashes").Default("created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, peer = EXCLUDED.peer, ip = EXCLUDED.ip, port = EXCLUDED.port")
}

func UnknownHashDeleteByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewUnknownHashScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM torrents_unknown_infohashes WHERE "id" = {id} RETURNING ` + UnknownHashScannerStaticColumns)
}

func UnknownHashCooldown(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a UnknownHash) NewUnknownHashScannerStaticRow,
) {
	gql.Into("torrents_unknown_infohashes").Default("created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, attempts = EXCLUDED.attempts + 1, next_check = NOW() + least(to_minutes(CAST(EXCLUDED.attempts AS INT)*2), to_hours(24))")
}

func RSS(
	gql genieql.Structure,
) {
	gql.From(
		gql.Table("torrents_feed_rss"),
	)
}

func RSSScanner(
	gql genieql.Scanner,
	pattern func(i RSS),
) {
	gql.ColumnNamePrefix("torrents_feed_rss.")
}

func RSSInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a RSS) NewRSSScannerStaticRow,
) {
	gql.Into("torrents_feed_rss").Default("created_at", "updated_at", "disabled_at", "ttl_minimum").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, autodownload = EXCLUDED.autodownload, autoarchive = EXCLUDED.autoarchive, url = EXCLUDED.url, description = EXCLUDED.description, next_check = EXCLUDED.next_check, digest = EXCLUDED.digest")
}

func RSSInsertDefaultFeed(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a RSS) NewRSSScannerStaticRow,
) {
	gql.Into("torrents_feed_rss").Default("created_at", "updated_at", "next_check", "disabled_at", "ttl_minimum").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT")
}

func RSSCooldownByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, ttl int, digest string, lastbuild time.Time) NewRSSScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_feed_rss SET updated_at = DEFAULT, next_check = NOW() + to_minutes({ttl}), last_built_at = {lastbuild}, digest = {digest} WHERE "id" = {id} RETURNING ` + RSSScannerStaticColumns)
}

func RSSDeleteByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewRSSScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM torrents_feed_rss WHERE "id" = {id} RETURNING ` + RSSScannerStaticColumns)
}

func RSSDeleteByURL(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, url string) NewRSSScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM torrents_feed_rss WHERE "url" = {url} RETURNING ` + RSSScannerStaticColumns)
}

func RSSFindByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewRSSScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + RSSScannerStaticColumns + ` FROM torrents_feed_rss WHERE "id" = {id}`)
}

func RSSFindByURL(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, url string) NewRSSScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + RSSScannerStaticColumns + ` FROM torrents_feed_rss WHERE "url" = {url}`)
}
