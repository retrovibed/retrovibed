package communityapi_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/stretchr/testify/require"
)

func TestSubscribeCommunity(t *testing.T) {
	t.Run("subscribes and registers an rss feed", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		communityID := uuid.Must(uuid.NewV7()).String()

		sub, err := communityapi.SubscribeCommunity(ctx, q, newCommunityMockClient(communityID), communityID)
		require.NoError(t, err)
		require.Equal(t, communityID, sub.ID)
		require.True(t, sub.SubscribedAt.Before(timex.Inf()))
		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_feed_rss"))
	})

	t.Run("is idempotent for an already subscribed community", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)
		communityID := uuid.Must(uuid.NewV7()).String()
		client := newCommunityMockClient(communityID)

		_, err := communityapi.SubscribeCommunity(ctx, q, client, communityID)
		require.NoError(t, err)
		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_feed_rss"))

		sub, err := communityapi.SubscribeCommunity(ctx, q, client, communityID)
		require.NoError(t, err)
		require.Equal(t, communityID, sub.ID)
		require.Equal(t, 1, sqltestx.Count(t, q, "SELECT COUNT(*) FROM torrents_feed_rss"))
	})
}
