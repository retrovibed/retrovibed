package tracking

import (
	"database/sql"
	"testing"
	"time"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestPeerNextCheck(t *testing.T) {
	t.Run("claims peers one at a time without double claiming", func(t *testing.T) {
		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		const total = 3

		expected := map[string]struct{}{}
		for range total {
			p := NewPeer(int160.Random(), PeerOptionBEP51(1, 60))
			require.NoError(t, PeerInsertWithDefaults(ctx, q, p).Scan(&p))
			expected[p.ID] = struct{}{}
		}

		seen := map[string]struct{}{}
		for range total {
			var claimed Peer
			require.NoError(t, PeerNextCheck(ctx, q, time.Now().Add(time.Hour)).Scan(&claimed))

			_, ok := expected[claimed.ID]
			require.True(t, ok, "unexpected peer claimed: %s", claimed.ID)
			_, duplicate := seen[claimed.ID]
			require.False(t, duplicate, "peer %s claimed more than once", claimed.ID)
			seen[claimed.ID] = struct{}{}
		}

		require.Len(t, seen, total, "every eligible peer should be claimed exactly once")
	})

	t.Run("returns no rows once every peer is claimed", func(t *testing.T) {
		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		p := NewPeer(int160.Random(), PeerOptionBEP51(1, 60))
		require.NoError(t, PeerInsertWithDefaults(ctx, q, p).Scan(&p))

		var claimed Peer
		require.NoError(t, PeerNextCheck(ctx, q, time.Now().Add(time.Hour)).Scan(&claimed))

		var exhausted Peer
		err := PeerNextCheck(ctx, q, time.Now().Add(time.Hour)).Scan(&exhausted)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("skips peers without info hashes", func(t *testing.T) {
		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		p := NewPeer(int160.Random())
		require.NoError(t, PeerInsertWithDefaults(ctx, q, p).Scan(&p))

		var claimed Peer
		err := PeerNextCheck(ctx, q, time.Now().Add(time.Hour)).Scan(&claimed)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}
