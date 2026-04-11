//go:build genieql.generate
// +build genieql.generate

package identityssh

import (
	"context"

	genieql "github.com/james-lawrence/genieql/ginterp"
	"github.com/retrovibed/retrovibed/internal/sqlx"
)

func Identity(gql genieql.Structure) {
	gql.From(
		gql.Table("meta_sso_identity_ssh"),
	)
}

func IdentityScanner(gql genieql.Scanner, pattern func(i Identity)) {
	gql.ColumnNamePrefix("meta_sso_identity_ssh.")
}

func IdentityFindBy(gql genieql.QueryAutogen, ctx context.Context, q sqlx.Queryer, e Identity) NewIdentityScannerStaticRow {
	gql.From("meta_sso_identity_ssh")
}

func IdentityInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Identity) NewIdentityScannerStaticRow,
) {
	gql.Into("meta_sso_identity_ssh").Default("created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, profile_id = EXCLUDED.profile_id")
}

// used to insert an identity for the oauth2 endpoint
func IdentityInsertZeroWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Identity) NewIdentityScannerStaticRow,
) {
	gql.Into("meta_sso_identity_ssh").Default("created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT")
}
