//go:build genieql.generate
// +build genieql.generate

package tracking

import (
	"context"

	genieql "github.com/james-lawrence/genieql/ginterp"
	"github.com/retrovibed/retrovibed/internal/sqlx"
)

//easyjson:json
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
	gql.Into("torrents_metadata").Default("created_at", "updated_at", "hidden_at", "initiated_at", "paused_at", "downloaded").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, tracker = EXCLUDED.tracker")
}

func MetadataBatchInsertWithDefaults(
	gql genieql.InsertBatch,
	pattern func(ctx context.Context, q sqlx.Queryer, p Metadata) NewMetadataScannerStatic,
) {
	gql.Into("torrents_metadata").Batch(10).Default("id", "created_at", "updated_at")
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

func MetadataProgressByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, peers uint16, completed uint64) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_metadata SET updated_at = NOW(), downloaded = {completed}, peers = {peers}, seeding = (bytes == {completed}) WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

//easyjson:json
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
	gql.Into("torrents_peers").Default("created_at", "updated_at", "next_check").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, ip = EXCLUDED.ip, port = EXCLUDED.port, bep51_available = EXCLUDED.bep51_available")
}

func PeerMarkNextCheck(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Peer) NewPeerScannerStaticRow,
) {
	gql.Into("torrents_peers").Default("created_at", "updated_at", "next_check").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = NOW(), next_check = NOW() + to_seconds(EXCLUDED.bep51_ttl)")
}

//easyjson:json
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
	gql.Into("torrents_unknown_infohashes").Default("created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT")
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

//easyjson:json
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
	gql.Into("torrents_feed_rss").Default("created_at", "updated_at", "next_check", "disabled_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, autodownload = EXCLUDED.autodownload, autoarchive = EXCLUDED.autoarchive, url = EXCLUDED.url, description = EXCLUDED.description")
}

func RSSCooldown(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a RSS) NewRSSScannerStaticRow,
) {
	gql.Into("torrents_feed_rss").Default("created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, next_check = NOW() + to_minutes({ttl})")
}

func RSSCooldownByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, ttl int) NewRSSScannerStaticRow,
) {
	gql = gql.Query(`UPDATE torrents_feed_rss SET updated_at = DEFAULT, next_check = NOW() + to_minutes({ttl}) WHERE "id" = {id} RETURNING ` + RSSScannerStaticColumns)
}

func RSSDeleteByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewRSSScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM torrents_feed_rss WHERE "id" = {id} RETURNING ` + RSSScannerStaticColumns)
}
