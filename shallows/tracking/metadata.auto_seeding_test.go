package tracking

import (
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/stretchr/testify/require"
)

func TestMetadataOptionAutoSeeding(t *testing.T) {
	t.Run("available equal to bytes sets seeding", func(t *testing.T) {
		md := NewMetadata(
			new(int160.Random()),
			MetadataOptionBytes(1000),
			MetadataOptionAvailable(1000),
			MetadataOptionAutoSeeding,
		)

		require.True(t, md.Seeding)
	})

	t.Run("available less than bytes does not set seeding", func(t *testing.T) {
		md := NewMetadata(
			new(int160.Random()),
			MetadataOptionBytes(1000),
			MetadataOptionAvailable(500),
			MetadataOptionAutoSeeding,
		)

		require.False(t, md.Seeding)
	})

	t.Run("downloaded equal to bytes but available less does not set seeding", func(t *testing.T) {
		// self-published or otherwise locally-sourced content must not be
		// marked seeding just because Downloaded happens to equal Bytes -
		// only Available (bytes actually present on disk) governs seeding.
		md := NewMetadata(
			new(int160.Random()),
			MetadataOptionBytes(1000),
			MetadataOptionDownloaded(1000),
			MetadataOptionAvailable(0),
			MetadataOptionAutoSeeding,
		)

		require.False(t, md.Seeding)
	})
}
