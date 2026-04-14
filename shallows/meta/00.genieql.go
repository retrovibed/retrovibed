//go:build genieql.generate
// +build genieql.generate

package meta

import (
	"context"

	genieql "github.com/james-lawrence/genieql/ginterp"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
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

func ProfileEnable(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewProfileScannerStaticRow,
) {
	gql = gql.Query(`UPDATE meta_profiles SET disabled_at = DEFAULT, disabled_manually_at = DEFAULT, disabled_pending_approval_at = 'infinity' WHERE "meta_profiles"."id" = {id} RETURNING ` + ProfileScannerStaticColumns)
}

func ProfileInsertWithID(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Profile) NewProfileScannerStaticRow,
) {
	gql.Into("meta_profiles").Default("session_watermark", "created_at", "updated_at", "disabled_at", "disabled_manually_at", "disabled_pending_approval_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT")
}

func ProfileFindByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewProfileScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + ProfileScannerStaticColumns + ` FROM meta_profiles WHERE "id" = {id}`)
}

func ProfileDisableByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewProfileScannerStaticRow,
) {
	gql = gql.Query(`UPDATE meta_profiles SET updated_at = NOW(), disabled_manually_at = NOW(), disabled_pending_approval_at = 'infinity' WHERE "id" = {id} RETURNING ` + ProfileScannerStaticColumns)
}

func ProfileUpdate(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, p Profile) NewProfileScannerStaticRow,
) {
	gql = gql.Query(`UPDATE meta_profiles SET updated_at = NOW(), disabled_at = {p.DisabledAt}, disabled_manually_at = {p.DisabledManuallyAt}, disabled_pending_approval_at = {p.DisabledPendingApprovalAt}, display = {p.Display} WHERE id = {p.ID} RETURNING ` + ProfileScannerStaticColumns)
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

func DaemonUpdateByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, a Daemon) NewDaemonScannerStaticRow,
) {
	gql = gql.Query(`UPDATE meta_daemons SET updated_at = NOW(), description = {a.Description}, hostname = {a.Hostname} WHERE id = {id} RETURNING ` + DaemonScannerStaticColumns)
}

func DaemonFindByLatestUpdated(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewDaemonScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + DaemonScannerStaticColumns + ` FROM meta_daemons ORDER BY updated_at DESC LIMIT 1`)
}

func DaemonFindDefault(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewDaemonScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + DaemonScannerStaticColumns + ` FROM meta_daemons WHERE "default" LIMIT 1`)
}

func DaemonTouch(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewDaemonScannerStaticRow,
) {
	gql = gql.Query(`UPDATE meta_daemons SET updated_at = NOW(), "default" = 'f' WHERE "default"; UPDATE meta_daemons SET updated_at = NOW(), "default" = 't' WHERE "id" = {id} RETURNING ` + DaemonScannerStaticColumns)
}

func DaemonDownload(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewDaemonScannerStaticRow,
) {
	gql = gql.Query(`UPDATE meta_daemons SET updated_at = NOW(), downloads = 'f' WHERE downloads; UPDATE meta_daemons SET updated_at = NOW(), downloads = 't' WHERE "id" = {id} RETURNING ` + DaemonScannerStaticColumns)
}

func DaemonFindByDownload(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewDaemonScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + DaemonScannerStaticColumns + ` FROM meta_daemons WHERE downloads LIMIT 1`)
}

func DaemonDeleteByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewDaemonScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM meta_daemons WHERE "id" = {id} RETURNING ` + DaemonScannerStaticColumns)
}

func ConsumedToken(gql genieql.Structure) {
	gql.From(
		gql.Table("meta_consumed_tokens"),
	)
}

func ConsumedTokenScanner(gql genieql.Scanner, pattern func(a ConsumedToken)) {
	gql.ColumnNamePrefix("meta_consumed_tokens.")
}

func ConsumedTokenFindBy(gql genieql.QueryAutogen, ctx context.Context, q sqlx.Queryer, e ConsumedToken) NewConsumedTokenScannerStaticRow {
	gql.From("meta_consumed_tokens")
}

func ConsumedTokenInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a ConsumedToken) NewConsumedTokenScannerStaticRow,
) {
	gql.Into("meta_consumed_tokens").Default("created_at")
}

func Authz(gql genieql.Structure) {
	gql.From(
		gql.Table("authz_meta"),
	)
}

func AuthzScanner(gql genieql.Scanner, pattern func(a Authz)) {
	gql.ColumnNamePrefix("authz_meta.")
}

func AuthzFindBy(gql genieql.QueryAutogen, ctx context.Context, q sqlx.Queryer, e Authz) NewAuthzScannerStaticRow {
	gql.From("authz_meta")
}

func AuthzInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Authz) NewAuthzScannerStaticRow,
) {
	gql.Into("authz_meta").Default("id", "created_at")
}

func AuthzInsertWithIDDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Authz) NewAuthzScannerStaticRow,
) {
	gql.Into("authz_meta").Default("created_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT")
}

// upsert a single record with default fields.
func AuthzUpsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, p Authz) NewAuthzScannerStaticRow,
) {
	gql.Into("authz_meta").
		Default("id", "created_at", "updated_at").
		Conflict("ON CONFLICT (profile_id) DO UPDATE SET usermanagement = EXCLUDED.usermanagement, billing_read = EXCLUDED.billing_read, billing_modify = EXCLUDED.billing_modify, community_modify = EXCLUDED.community_modify, library_read = EXCLUDED.library_read, library_modify = EXCLUDED.library_modify, updated_at = DEFAULT")
}

func AuthzDeleteByProfileID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewAuthzScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM authz_meta WHERE profile_id = {id} RETURNING ` + AuthzScannerStaticColumns)
}

func Wireguard(gql genieql.Structure) {
	gql.From(
		gql.Table("meta_wireguard"),
	)
}

func WireguardScanner(gql genieql.Scanner, pattern func(a Wireguard)) {
	gql.ColumnNamePrefix("meta_wireguard.")
}

func WireguardFindByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewWireguardScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + WireguardScannerStaticColumns + ` FROM meta_wireguard WHERE id = {id}`)
}

func WireguardInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Wireguard) NewWireguardScannerStaticRow,
) {
	gql.Into("meta_wireguard").Default("created_at", "updated_at")
}

func WireguardDeleteByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewWireguardScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM meta_wireguard WHERE id = {id} RETURNING ` + WireguardScannerStaticColumns)
}

func WireguardCurrent(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewWireguardScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + WireguardScannerStaticColumns + ` FROM meta_wireguard WHERE "default" LIMIT 1`)
}

func WireguardTouch(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewWireguardScannerStaticRow,
) {
	gql = gql.Query(`UPDATE meta_wireguard SET "default" = (id = {id}) WHERE "default" OR id = {id} RETURNING ` + WireguardScannerStaticColumns)
}

func WireguardUpdate(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, a Wireguard) NewWireguardScannerStaticRow,
) {
	gql = gql.Query(`UPDATE meta_wireguard SET description = {a.Description}, port = {a.Port}, dns_rate_limit = {a.DNSRateLimit} WHERE id = {a.ID} RETURNING ` + WireguardScannerStaticColumns)
}
