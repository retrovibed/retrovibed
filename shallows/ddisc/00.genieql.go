//go:build genieql.generate
// +build genieql.generate

package ddisc

import (
	"context"
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
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
	gql.Into("ddisc_media").Default("created_at", "updated_at", "tombstoned_at", "released_at", "next_check_at", "policy_rank", "policy_rejection").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT")
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
	gql = gql.Query(`UPDATE ddisc_media SET updated_at = NOW(), next_check_at = 'infinity', known_media_id = {disc.KnownMediaID}, mimetype = {disc.Mimetype}, category = {disc.Category} WHERE id = {disc.ID} RETURNING ` + DiscoveredScannerStaticColumns)
}

func DiscoveredCooldown(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a Discovered) NewDiscoveredScannerStaticRow,
) {
	gql.Into("ddisc_media").Default("created_at", "updated_at", "tombstoned_at", "released_at", "next_check_at", "policy_rank", "policy_rejection").Conflict("ON CONFLICT (id) DO UPDATE SET updated_at = DEFAULT, attempts = EXCLUDED.attempts + 1, next_check_at = NOW() + least(to_minutes(CAST(EXCLUDED.attempts AS INT)*2), to_hours(24))")
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

// DiscoveredRank persists the result of running a ranking Policy against a
// Discovered row: its computed health, policy_rank, and policy_rejection.
func DiscoveredRank(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, health uint32, rank uint16, rejection string) NewDiscoveredScannerStaticRow,
) {
	gql = gql.Query(`UPDATE ddisc_media SET updated_at = NOW(), health = {health}, policy_rank = {rank}, policy_rejection = {rejection} WHERE "id" = {id} RETURNING ` + DiscoveredScannerStaticColumns)
}

func SearchQueue(
	gql genieql.Structure,
) {
	gql.From(
		gql.Table("ddisc_search_queue"),
	)
}

func SearchQueueScanner(
	gql genieql.Scanner,
	pattern func(i SearchQueue),
) {
	gql.ColumnNamePrefix("ddisc_search_queue.")
}

// SearchQueueEnqueue records a known_media_id as needing external search
// plugin discovery. Idempotent: re-enqueuing an already-pending id just
// touches updated_at.
func SearchQueueEnqueue(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a SearchQueue) NewSearchQueueScannerStaticRow,
) {
	gql.Into("ddisc_search_queue").Default("id", "created_at", "updated_at", "next_check_at", "attempts").Conflict("ON CONFLICT (known_media_id) DO UPDATE SET updated_at = NOW()")
}

func SearchQueuePending(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewSearchQueueScannerStatic,
) {
	gql = gql.Query(`SELECT ` + SearchQueueScannerStaticColumns + ` FROM ddisc_search_queue WHERE next_check_at <= NOW()`)
}

// SearchQueueResolve removes a known_media_id from the queue once at least
// one Discovered row has been persisted for it.
func SearchQueueResolve(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, kid string) NewSearchQueueScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM ddisc_search_queue WHERE "known_media_id" = {kid} RETURNING ` + SearchQueueScannerStaticColumns)
}

// SearchQueueCooldown bumps attempts and pushes next_check_at out after a
// drain attempt found nothing, same backoff shape as DiscoveredCooldown.
func SearchQueueCooldown(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, kid string) NewSearchQueueScannerStaticRow,
) {
	gql = gql.Query(`UPDATE ddisc_search_queue SET updated_at = NOW(), attempts = attempts + 1, next_check_at = NOW() + least(to_minutes(CAST(attempts AS INT)*2), to_hours(24)) WHERE "known_media_id" = {kid} RETURNING ` + SearchQueueScannerStaticColumns)
}

// SearchQueuePurge deletes entries older than maxAge (measured from
// created_at) that have never resolved, so the queue doesn't grow unbounded
// with known-media-ids no plugin will ever find.
func SearchQueuePurge(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, maxAge time.Duration) NewSearchQueueScannerStatic,
) {
	gql = gql.Query(`DELETE FROM ddisc_search_queue WHERE created_at <= NOW() - to_seconds(CAST({maxAge} AS BIGINT) / 1000000000) RETURNING ` + SearchQueueScannerStaticColumns)
}
