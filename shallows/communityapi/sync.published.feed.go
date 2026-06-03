package communityapi

import (
	"context"
	"log"
	"net/http"

	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

// SyncFeed resyncs local state to the public feed.
func SyncFeed(ctx context.Context, q sqlx.Queryer, httpc *http.Client, published FeedPublisher) error {
	pending := sqlx.Scan(community.CommunitySyncStateLookupFeedSyncRequests(ctx, q))

	for cs := range pending.Iter() {
		if err := RegenerateFeed(ctx, q, published, cs.CommunityID); err != nil {
			log.Println(errorsx.Wrapf(err, "feed regeneration failed for community: %s", cs.CommunityID))
			continue
		}
	}

	return errorsx.Wrap(pending.Err(), "unable to publish feeds")
}
