package communityapi

import (
	"context"
	"net/http"

	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/slicesx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

// SubscribeCommunity looks up cid in deeppool, upserts the local community record,
// and — if not already subscribed — marks it subscribed and registers its RSS feed.
// Idempotent: calling it again for an already-subscribed community is a no-op.
func SubscribeCommunity(ctx context.Context, q sqlx.Queryer, httpc *http.Client, cid string) (existing community.Community, err error) {
	com, err := NewDeeppoolCommunity(httpc).Find(ctx, cid)
	if err != nil {
		return existing, errorsx.Wrapf(err, "unable to find community from deeppool - %s", cid)
	}

	if err := community.CommunityInsertWithDefaults(ctx, q, CommunityFromDeeppool(com)).Scan(&existing); err != nil {
		return existing, errorsx.Wrap(err, "unable to upsert community")
	}

	if existing.SubscribedAt.Before(timex.Inf()) {
		return existing, nil
	}

	sub := community.Community{ID: cid}
	if err := community.CommunityUpsertAutoDownload(ctx, q, sub).Scan(&sub); err != nil {
		return existing, errorsx.Wrap(err, "unable to upsert community")
	}

	if err := community.CommunitySubscribe(ctx, q, cid).Scan(&existing); err != nil {
		return existing, errorsx.Wrap(err, "unable to subscribe")
	}

	feed := tracking.NewFeedRSS(
		"",
		tracking.RSSOptionURL(existing.URL),
		tracking.RSSOptionDescription(stringsx.Join(" - ", slicesx.Filter(stringsx.Present, community.CommunityDomainFromURL(com.Url), com.Description)...)),
		tracking.RSSOptionEncryptionSeed(com.Entropy),
		tracking.RSSOptionAutodownload(true),
		tracking.RSSOptionAutoarchive(true),
		tracking.RSSOptionAutoID,
	)

	if err := tracking.RSSInsertWithDefaults(ctx, q, feed).Scan(&feed); err != nil {
		return existing, errorsx.Wrap(err, "unable to register rss feed")
	}

	return existing, nil
}
