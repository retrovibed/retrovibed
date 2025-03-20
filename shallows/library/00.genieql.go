//go:build genieql.generate
// +build genieql.generate

package library

import (
	"context"

	genieql "github.com/james-lawrence/genieql/ginterp"
	"github.com/retrovibed/retrovibed/internal/x/sqlx"
)

//easyjson:json
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
	gql.Into("library_metadata").Default("created_at", "updated_at", "hidden_at", "tombstoned_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT")
}

func MetadataDeleteByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM library_metadata WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataTombstoneByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE library_metadata SET tombstoned_at = NOW(), initiated_at = 'infinity' WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataHideByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`UPDATE library_metadata SET hidden_at = NOW(), initiated_at = 'infinity' WHERE "id" = {id} RETURNING ` + MetadataScannerStaticColumns)
}

func MetadataFindByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewMetadataScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + MetadataScannerStaticColumns + ` FROM library_metadata WHERE "id" = {id}`)
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
