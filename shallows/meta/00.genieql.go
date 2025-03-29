//go:build genieql.generate
// +build genieql.generate

package meta

import (
	"context"

	genieql "github.com/james-lawrence/genieql/ginterp"
	"github.com/retrovibed/retrovibed/internal/sqlx"
)

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
	pattern func(ctx context.Context, q sqlx.Queryer, a Profile) NewProfileScannerStaticRow,
) {
	gql.Into("meta_profiles").Default("id", "session_watermark", "created_at", "updated_at", "disabled_at", "disabled_manually_at", "disabled_pending_approval_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT")
}

func ProfileFindByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewProfileScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + ProfileScannerStaticColumns + ` FROM meta_profiles WHERE "id" = {id}`)
}

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
	pattern func(ctx context.Context, q sqlx.Queryer, a IdentitySSH) NewIdentitySSHScannerStaticRow,
) {
	gql.Into("meta_sso_identity_ssh").Default("created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT")
}

func Daemon(gql genieql.Structure) {
	gql.From(
		gql.Table("meta_daemons"),
	)
}

func DaemonScanner(gql genieql.Scanner, pattern func(i Daemon)) {
	gql.ColumnNamePrefix("meta_daemons.")
}

func DaemonInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Daemon) NewDaemonScannerStaticRow,
) {
	gql.Into("meta_daemons").Default("created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT")
}

func DaemonFindByLatestUpdated(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewDaemonScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + DaemonScannerStaticColumns + ` FROM meta_daemons ORDER BY updated_at DESC LIMIT 1`)
}
