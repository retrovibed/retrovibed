//go:build genieql.generate
// +build genieql.generate

package ddisc

import (
	"context"

	"github.com/retrovibed/retrovibed/internal/sqlx"
)

func Discovered(
	gql genieql.Structure,
) {
	gql.From(
		gql.Table("ddisc_media"),
	)
}

func DiscoveredScanner(
	gql genieql.Scanner,
	pattern func(i Discovered),
) {
	gql.ColumnNamePrefix("ddisc_media.")
}

func DiscoveredInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Discovered) NewDiscoveredScannerStaticRow,
) {
	gql.Into("ddisc_media").Default("created_at", "updated_at", "tombstoned_at", "released_at", "next_check_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT")
}

func DiscoveredFindByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewDiscoveredScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + DiscoveredScannerStaticColumns + ` FROM ddisc_media WHERE "id" = {id}`)
}

func DiscoveredDeleteByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewDiscoveredScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM ddisc_media WHERE "id" = {id} RETURNING ` + DiscoveredScannerStaticColumns)
}

func DiscoveredFindByNextCheckAt(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewDiscoveredScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + DiscoveredScannerStaticColumns + ` ddisc_media ORDER BY "next_check_at" ASC LIMIT 1`)
}

func DiscoveredIndexed(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, disc Discovered) NewDiscoveredScannerStaticRow,
) {
	gql = gql.Query(`UPDATE ddisc_media SET updated_at = NOW(), next_check_at = 'infinity', known_media_id = {disc.KnownMediaID}, mimetype = {disc.Mimetype} WHERE id = {disc.ID} RETURNING ` + DiscoveredScannerStaticColumns)
}

func DiscoveredCooldown(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Discovered) NewDiscoveredScannerStaticRow,
) {
	gql.Into("ddisc_media").Default("created_at", "updated_at", "tombstoned_at", "released_at", "next_check_at").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, attempts = EXCLUDED.attempts + 1, next_check_at = NOW() + least(to_minutes(CAST(EXCLUDED.attempts AS INT)*2), to_hours(24))")
}

func DiscoveredSinceSync(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, sync string) NewDiscoveredScannerStatic,
) {
	gql = gql.Query(`SELECT ` + DiscoveredScannerStaticColumns + ` FROM ddisc_media WHERE sync_uid > {sync} AND known_media_id IN ('FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF'::uuid, '00000000-0000-0000-0000-000000000000'::uuid) ORDER BY "sync_uid" ASC`)
}

func DiscoveredPartitionSync(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, sync, partition string) NewDiscoveredScannerStatic,
) {
	gql = gql.Query(`SELECT ` + DiscoveredScannerStaticColumns + ` FROM ddisc_media WHERE sync_uid > {sync} AND partition = {partition} ORDER BY "sync_uid" ASC`)
}

func DiscoveredByKnownID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, kid string) NewDiscoveredScannerStatic,
) {
	gql = gql.Query(`SELECT ` + DiscoveredScannerStaticColumns + ` FROM ddisc_media WHERE known_media_id = {kid}`)
}
