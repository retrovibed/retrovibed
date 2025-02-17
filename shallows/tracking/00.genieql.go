//go:build genieql.generate
// +build genieql.generate

package tracking

import (
	"context"

	"github.com/james-lawrence/deeppool/internal/x/sqlx"
	genieql "github.com/james-lawrence/genieql/ginterp"
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
	gql.Into("torrents_metadata").Default("created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT")
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
	gql.Into("torrents_unknown_infohashes").Default("created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, attempts = EXCLUDED.attempts + 1, next_check = NOW() + least(to_hours(CAST(EXCLUDED.attempts AS INT)), to_hours(24))")
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
	gql.Into("torrents_feed_rss").Default("created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT")
}

func RSSCooldown(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a RSS) NewRSSScannerStaticRow,
) {
	gql.Into("torrents_feed_rss").Default("created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, next_check = NOW() + to_hours(24)")
}
