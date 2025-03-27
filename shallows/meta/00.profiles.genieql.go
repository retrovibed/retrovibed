//go:build genieql.generate
// +build genieql.generate

package meta

import (
	"context"

	genieql "github.com/james-lawrence/genieql/ginterp"
	"github.com/retrovibed/retrovibed/internal/sqlx"
)

//easyjson:json
func Profile(gql genieql.Structure) {
	gql.From(
		gql.Table("meta_profiles"),
	)
}

func ProfileScanner(gql genieql.Scanner, pattern func(i Profile)) {
	gql.ColumnNamePrefix("meta_profiles.")
}

func ProfileInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Metadata) NewProfileScannerStaticRow,
) {
	gql.Into("meta_profiles").Default("created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT")
}
