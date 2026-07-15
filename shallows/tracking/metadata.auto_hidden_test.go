package tracking

import (
	"testing"
	"time"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/stretchr/testify/require"
)

func TestMetadataOptionAutoHidden(t *testing.T) {
	t.Run("media archive sets hidden", func(t *testing.T) {
		md := NewMetadata(
			new(int160.Random()),
			MetadataOptionMimetype(mimex.RetrovibedMediaArchive),
			MetadataOptionAutoHidden,
		)

		require.NotEqual(t, timex.Inf(), md.HiddenAt)
		require.WithinDuration(t, time.Now(), md.HiddenAt, time.Second)
	})

	t.Run("neural sets hidden", func(t *testing.T) {
		md := NewMetadata(
			new(int160.Random()),
			MetadataOptionMimetype(mimex.RetrovibedNeural),
			MetadataOptionAutoHidden,
		)

		require.NotEqual(t, timex.Inf(), md.HiddenAt)
		require.WithinDuration(t, time.Now(), md.HiddenAt, time.Second)
	})

	t.Run("discovery search module sets hidden", func(t *testing.T) {
		md := NewMetadata(
			new(int160.Random()),
			MetadataOptionMimetype(mimex.RetrovibedDiscoverySearch),
			MetadataOptionAutoHidden,
		)

		require.NotEqual(t, timex.Inf(), md.HiddenAt)
		require.WithinDuration(t, time.Now(), md.HiddenAt, time.Second)
	})

	t.Run("bittorrent does not set hidden", func(t *testing.T) {
		md := NewMetadata(
			new(int160.Random()),
			MetadataOptionMimetype(mimex.Bittorrent),
			MetadataOptionAutoHidden,
		)

		require.Equal(t, timex.Inf(), md.HiddenAt)
	})

	t.Run("unknown mimetype does not set hidden", func(t *testing.T) {
		md := NewMetadata(
			new(int160.Random()),
			MetadataOptionMimetype("application/octet-stream"),
			MetadataOptionAutoHidden,
		)

		require.Equal(t, timex.Inf(), md.HiddenAt)
	})
}
