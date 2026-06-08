package communityapi

import (
	"github.com/retrovibed/retrovibed/shallows/community"
	"net/http"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestSyncFeed(t *testing.T) {
	t.Run("regenerates feed for each community with a pending sync request", func(t *testing.T) {
		var (
			ctx, done    = testx.Context(t)
			q            = sqltestx.Metadatabase(t)
			feeds        = &mockFeedPublisher{}
			communityID1 = uuid.Must(uuid.NewV7()).String()
			communityID2 = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		for _, cid := range []string{communityID1, communityID2} {
			var cs community.Community
			require.NoError(t, community.CommunityRequestFeedSync(ctx, q, community.Community{
				ID:         cid,
				SyncFeedAt: time.Now().Add(-time.Minute),
			}).Scan(&cs))
		}

		require.NoError(t, SyncFeed(ctx, q, http.DefaultClient, feeds))
		require.Contains(t, feeds.feeds, communityID1)
		require.Contains(t, feeds.feeds, communityID2)
	})

	t.Run("no feeds regenerated when no pending sync requests", func(t *testing.T) {
		var (
			ctx, done = testx.Context(t)
			q         = sqltestx.Metadatabase(t)
			feeds     = &mockFeedPublisher{}
		)
		defer done()

		require.NoError(t, SyncFeed(ctx, q, http.DefaultClient, feeds))
		require.Empty(t, feeds.feeds)
	})

	t.Run("does not regenerate feed for sync request scheduled in the future", func(t *testing.T) {
		var (
			ctx, done   = testx.Context(t)
			q           = sqltestx.Metadatabase(t)
			feeds       = &mockFeedPublisher{}
			communityID = uuid.Must(uuid.NewV7()).String()
		)
		defer done()

		var cs community.Community
		require.NoError(t, community.CommunityRequestFeedSync(ctx, q, community.Community{
			ID:         communityID,
			SyncFeedAt: time.Now().Add(time.Hour),
		}).Scan(&cs))

		require.NoError(t, SyncFeed(ctx, q, http.DefaultClient, feeds))
		require.Empty(t, feeds.feeds)
	})
}
