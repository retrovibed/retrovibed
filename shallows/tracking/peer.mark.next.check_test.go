package tracking

import (
	"testing"
	"time"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/stretchr/testify/require"
)

func TestPeerMarkNextCheck(t *testing.T) {
	t.Run("insert defaults next_check to negative infinity", func(t *testing.T) {
		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		p := NewPeer(int160.Random(), PeerOptionBEP51(1, 30))

		var inserted Peer
		require.NoError(t, PeerMarkNextCheck(ctx, q, p).Scan(&inserted))

		require.True(t, timex.NegInf().Equal(inserted.NextCheck), "expected -infinity, got %s", inserted.NextCheck)
	})

	t.Run("upsert reschedules next_check from bep51_ttl and updates tombstone", func(t *testing.T) {
		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		id := int160.Random()

		original := NewPeer(id, PeerOptionBEP51(1, 30))
		require.NoError(t, PeerMarkNextCheck(ctx, q, original).Scan(&original))

		const ttl = 120 * time.Second
		updatedTombstone := time.Now().Add(48 * time.Hour).Truncate(time.Millisecond)

		rescheduled := NewPeer(
			id,
			PeerOptionBEP51(1, uint16(ttl.Seconds())),
			PeerOptionTombstone(updatedTombstone),
		)

		before := time.Now()
		var stored Peer
		require.NoError(t, PeerMarkNextCheck(ctx, q, rescheduled).Scan(&stored))

		require.False(t, timex.NegInf().Equal(stored.NextCheck), "next_check should have been rescheduled")
		require.WithinDuration(t, before.Add(ttl), stored.NextCheck, 10*time.Second)
		require.True(t, updatedTombstone.Equal(stored.TombstonedAt), "expected %s, got %s", updatedTombstone, stored.TombstonedAt)
	})
}
