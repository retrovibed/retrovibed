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
	gql.Into("torrents_peers").Default("created_at", "updated_at", "next_check").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, ip = EXCLUDED.ip")
}
