package community_test

import (
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/stretchr/testify/require"
)

func TestPublishedContentInsertWithDefaults(t *testing.T) {
	t.Run("upsert reactivates a tombstoned row instead of leaving it tombstoned", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		q := sqltestx.Metadatabase(t)

		var pc community.PublishedContent
		require.NoError(t, testx.Fake(&pc, community.PublishedContentOptionTestDefaults))
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))
		require.Equal(t, timex.Inf(), pc.TombstonedAt)

		require.NoError(t, community.PublishedContentTombstone(ctx, q, pc.ID).Scan(&pc))
		require.NotEqual(t, timex.Inf(), pc.TombstonedAt)

		// republishing to the same community_id/library_id pair hits the ON
		// CONFLICT arbiter on the same row; it must reactivate the row rather
		// than silently leaving it tombstoned while its content is updated.
		require.NoError(t, community.PublishedContentInsertWithDefaults(ctx, q, pc).Scan(&pc))
		require.Equal(t, timex.Inf(), pc.TombstonedAt)
	})
}
