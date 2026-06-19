package tracking

import (
	"testing"
	"time"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/stretchr/testify/require"
)

func TestMetadataOptionExpiresAtTTL(t *testing.T) {
	t.Run("duration > 0", func(t *testing.T) {
		md := NewMetadata(
			new(int160.Random()),
			MetadataOptionExpiresAtTTL(5*time.Minute),
		)

		require.False(t, md.ExpiresAt.IsZero())
		require.WithinDuration(t, md.ExpiresAt, time.Now(), 5*time.Minute+time.Second)
	})

	t.Run("duration == 0 should use the default", func(t *testing.T) {
		md := NewMetadata(
			new(int160.Random()),
			MetadataOptionExpiresAtTTL(0),
		)

		require.Equal(t, timex.Inf(), md.ExpiresAt)
	})
}
