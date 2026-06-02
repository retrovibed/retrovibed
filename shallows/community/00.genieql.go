//go:build genieql.generate
// +build genieql.generate

package community

import (
	"context"
	"time"

	genieql "github.com/james-lawrence/genieql/ginterp"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

func PublishedContent(gql genieql.Structure) {
	gql.From(
		gql.Table("published_content"),
	)
}

func PublishedContentScanner(gql genieql.Scanner, pattern func(i PublishedContent)) {
	gql.ColumnNamePrefix("published_content.")
}

func IDScanner(gql genieql.Scanner, pattern func(i string)) {}

func PublishedContentInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a PublishedContent) NewPublishedContentScannerStaticRow,
) {
	gql.Into("published_content").Default("id", "created_at", "updated_at", "tombstoned_at").Conflict("ON CONFLICT (community_id, library_id) DO UPDATE SET magnet_uri = EXCLUDED.magnet_uri, known_media_id = EXCLUDED.known_media_id, publish_mode = EXCLUDED.publish_mode, oauth_google_id = EXCLUDED.oauth_google_id, bytes = EXCLUDED.bytes, published_at = COALESCE(NULLIF(published_at, 'infinity'), EXCLUDED.published_at), updated_at = DEFAULT")
}

func PublishedContentFindByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewPublishedContentScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + PublishedContentScannerStaticColumns + ` FROM published_content WHERE "id" = {id}`)
}

func PublishedContentFindByMagnetURI(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, magnetURI string) NewPublishedContentScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + PublishedContentScannerStaticColumns + ` FROM published_content WHERE "magnet_uri" = {magnetURI}`)
}

func PublishedContentSearchByMagnetPrefix(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, magnetPrefix string) NewPublishedContentScannerStatic,
) {
	gql = gql.Query(`SELECT ` + PublishedContentScannerStaticColumns + ` FROM published_content WHERE "magnet_uri" LIKE {magnetPrefix}`)
}

func PublishedContentDeleteByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewPublishedContentScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM published_content WHERE "id" = {id} RETURNING ` + PublishedContentScannerStaticColumns)
}

func PublishedContentFindByLibraryID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, libraryID string) NewPublishedContentScannerStatic,
) {
	gql = gql.Query(`SELECT ` + PublishedContentScannerStaticColumns + ` FROM published_content WHERE "library_id" = {libraryID}`)
}

func PublishedContentUpdatePublishedAt(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, publishedAt time.Time) NewPublishedContentScannerStaticRow,
) {
	gql = gql.Query(`UPDATE published_content SET published_at = {publishedAt}, updated_at = now() WHERE id = {id} RETURNING ` + PublishedContentScannerStaticColumns)
}

func PublishedContentTombstone(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewPublishedContentScannerStaticRow,
) {
	gql = gql.Query(`UPDATE published_content SET tombstoned_at = now(), updated_at = now() WHERE id = {id} RETURNING ` + PublishedContentScannerStaticColumns)
}

func PublishedContentFindByPendingSync(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewPublishedContentScannerStatic,
) {
	gql = gql.Query(`SELECT ` + PublishedContentScannerStaticColumns + ` FROM published_content WHERE published_at >= 'infinity' AND tombstoned_at = 'infinity' AND (publish_mode > 0 OR oauth_google_id != '00000000-0000-0000-0000-000000000000')`)
}

func PublishedContentFindByCommunityIDForFeed(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, communityID string) NewPublishedContentScannerStatic,
) {
	gql = gql.Query(`SELECT ` + PublishedContentScannerStaticColumns + ` FROM published_content WHERE "community_id" = {communityID} AND publish_mode > 0 ORDER BY published_at DESC`)
}

func PublishedContentUpdateMagnetURI(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string, magnetURI string) NewPublishedContentScannerStaticRow,
) {
	gql = gql.Query(`UPDATE published_content SET magnet_uri = {magnetURI}, updated_at = now() WHERE id = {id} RETURNING ` + PublishedContentScannerStaticColumns)
}

func CommunityMetric(gql genieql.Structure) {
	gql.From(
		gql.Table("community_metrics"),
	)
}

func CommunityMetricScanner(gql genieql.Scanner, pattern func(i CommunityMetric)) {
	gql.ColumnNamePrefix("community_metrics.")
}

func CommunityMetricInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a CommunityMetric) NewCommunityMetricScannerStaticRow,
) {
	gql.Into("community_metrics").Default("id", "created_at", "updated_at").Conflict("ON CONFLICT (community_id, period_start) DO UPDATE SET updated_at = DEFAULT, subscribers = GREATEST(community_metrics.subscribers, EXCLUDED.subscribers)")
}

func CommunityMetricFindByCommunityIDAndPeriod(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, communityID string, periodStart, periodEnd time.Time) NewCommunityMetricScannerStatic,
) {
	gql = gql.Query(`SELECT ` + CommunityMetricScannerStaticColumns + ` FROM community_metrics WHERE community_id = {communityID} AND period_start >= {periodStart} AND period_start <= {periodEnd} ORDER BY period_start DESC`)
}

func PublishedCASMetric(gql genieql.Structure) {
	gql.From(
		gql.Table("published_cas_metrics"),
	)
}

func PublishedCASMetricScanner(gql genieql.Scanner, pattern func(i PublishedCASMetric)) {
	gql.ColumnNamePrefix("published_cas_metrics.")
}

func PublishedCASMetricInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a PublishedCASMetric) NewPublishedCASMetricScannerStaticRow,
) {
	gql.Into("published_cas_metrics").Default("id", "created_at", "updated_at").Conflict("ON CONFLICT (published_content_id, period_start) DO UPDATE SET updated_at = DEFAULT, archivers = EXCLUDED.archivers, bytes = EXCLUDED.bytes, revenue = EXCLUDED.revenue")
}

func PublishedCASMetricFindByCommunityIDAndPeriod(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, communityID string, periodStart, periodEnd time.Time) NewPublishedCASMetricScannerStatic,
) {
	gql = gql.Query(`SELECT ` + PublishedCASMetricScannerStaticColumns + ` FROM published_cas_metrics JOIN published_content ON published_cas_metrics.published_content_id = published_content.id WHERE published_content.community_id = {communityID} AND published_cas_metrics.period_start >= {periodStart} AND published_cas_metrics.period_start <= {periodEnd} ORDER BY published_cas_metrics.period_start DESC`)
}

func CommunitySyncState(gql genieql.Structure) {
	gql.From(
		gql.Table("community_sync_state"),
	)
}

func CommunitySyncStateScanner(gql genieql.Scanner, pattern func(i CommunitySyncState)) {
	gql.ColumnNamePrefix("community_sync_state.")
}

func CommunitySyncStateInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a CommunitySyncState) NewCommunitySyncStateScannerStaticRow,
) {
	gql.Into("community_sync_state").Default("id", "created_at", "updated_at", "sync_feed_at").Conflict("ON CONFLICT (community_id) DO UPDATE SET updated_at = DEFAULT, last_sync_at = EXCLUDED.last_sync_at")
}

func CommunitySyncStateFindByCommunityID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, communityID string) NewCommunitySyncStateScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + CommunitySyncStateScannerStaticColumns + ` FROM community_sync_state WHERE community_id = {communityID}`)
}

func CommunitySyncStateRequestFeedSync(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a CommunitySyncState) NewCommunitySyncStateScannerStaticRow,
) {
	gql.Into("community_sync_state").Default("id", "created_at", "updated_at", "last_sync_at").Conflict("ON CONFLICT (community_id) DO UPDATE SET updated_at = DEFAULT, sync_feed_at = EXCLUDED.sync_feed_at")
}

func CommunitySyncStateLookupFeedSyncRequests(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewCommunitySyncStateScannerStatic,
) {
	gql = gql.Query(`SELECT ` + CommunitySyncStateScannerStaticColumns + ` FROM community_sync_state WHERE sync_feed_at < NOW()`)
}

func CommunitySyncStateRequestFeedSyncCompleted(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, cid string) NewCommunitySyncStateScannerStaticRow,
) {
	gql = gql.Query(`UPDATE community_sync_state SET sync_feed_at = DEFAULT, updated_at = NOW() WHERE community_id = {cid} RETURNING ` + CommunitySubscriptionScannerStaticColumns)
}

func CommunitySubscription(gql genieql.Structure) {
	gql.From(
		gql.Table("community_subscription"),
	)
}

func CommunitySubscriptionScanner(gql genieql.Scanner, pattern func(i CommunitySubscription)) {
	gql.ColumnNamePrefix("community_subscription.")
}

func CommunitySubscriptionInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a CommunitySubscription) NewCommunitySubscriptionScannerStaticRow,
) {
	gql.Into("community_subscription").Default("id", "created_at", "updated_at").Conflict("ON CONFLICT (community_id) DO UPDATE SET updated_at = DEFAULT, auto_download = EXCLUDED.auto_download")
}

func CommunitySubscriptionFindByCommunityID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, communityID string) NewCommunitySubscriptionScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + CommunitySubscriptionScannerStaticColumns + ` FROM community_subscription WHERE community_id = {communityID}`)
}

func CommunitySubscriptionDeleteByCommunityID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, communityID string) NewCommunitySubscriptionScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM community_subscription WHERE community_id = {communityID} RETURNING ` + CommunitySubscriptionScannerStaticColumns)
}

func CommunitySubscriptionFindAll(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewCommunitySubscriptionScannerStatic,
) {
	gql = gql.Query(`SELECT ` + CommunitySubscriptionScannerStaticColumns + ` FROM community_subscription ORDER BY created_at DESC`)
}

func CommunitySubscriptionUpdateLastSyncAt(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, communityID string, lastSyncAt time.Time) NewCommunitySubscriptionScannerStaticRow,
) {
	gql = gql.Query(`UPDATE community_subscription SET last_sync_at = {lastSyncAt}, updated_at = now() WHERE community_id = {communityID} RETURNING ` + CommunitySubscriptionScannerStaticColumns)
}

func OAuth2Google(gql genieql.Structure) {
	gql.From(
		gql.Table("oauth2_google"),
	)
}

func OAuth2GoogleScanner(gql genieql.Scanner, pattern func(i OAuth2Google)) {
	gql.ColumnNamePrefix("oauth2_google.")
}

func OAuth2GoogleInsertWithDefaults(
	gql genieql.Insert,
	pattern func(ctx context.Context, q sqlx.Queryer, a OAuth2Google) NewOAuth2GoogleScannerStaticRow,
) {
	gql.Into("oauth2_google").Default("id", "created_at", "updated_at").Conflict("ON CONFLICT (id) DO UPDATE SET access_token = EXCLUDED.access_token, refresh_token = EXCLUDED.refresh_token, token_type = EXCLUDED.token_type, expiry = EXCLUDED.expiry, scopes = EXCLUDED.scopes, updated_at = DEFAULT")
}

func OAuth2GoogleFindByID(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer, id string) NewOAuth2GoogleScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + OAuth2GoogleScannerStaticColumns + ` FROM oauth2_google WHERE "id" = {id}`)
}

func OAuth2GoogleFindFirst(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewOAuth2GoogleScannerStaticRow,
) {
	gql = gql.Query(`SELECT ` + OAuth2GoogleScannerStaticColumns + ` FROM oauth2_google LIMIT 1`)
}

func OAuth2GoogleDeleteAll(
	gql genieql.Function,
	pattern func(ctx context.Context, q sqlx.Queryer) NewOAuth2GoogleScannerStaticRow,
) {
	gql = gql.Query(`DELETE FROM oauth2_google RETURNING ` + OAuth2GoogleScannerStaticColumns)
}
