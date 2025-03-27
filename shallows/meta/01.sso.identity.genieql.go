//go:build genieql.generate
// +build genieql.generate

package meta

import (
	"context"

	genieql "github.com/james-lawrence/genieql/ginterp"
	"github.com/retrovibed/retrovibed/internal/sqlx"
)

//easyjson:json
func IdentitySSH(gql genieql.Structure) {
	gql.From(
		gql.Table("meta_sso_identity_ssh"),
	)
}

func IdentitySSHScanner(gql genieql.Scanner, pattern func(i IdentitySSH)) {
	gql.ColumnNamePrefix("meta_sso_identity_ssh.")
}

func IdentitySSHInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Metadata) NewIdentitySSHScannerStaticRow,
) {
	gql.Into("meta_sso_identity_ssh").Default("created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT")
}
