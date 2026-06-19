package tracking

import (
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/stretchr/testify/require"
)

func TestMetadataSamplePublic(t *testing.T) {
	t.Run("excludes private metadata", func(t *testing.T) {
		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		pub := NewMetadata(new(int160.Random()), MetadataOptionDescription("public"))
		pub.Private = false
		pub.Seeding = true
		require.NoError(t, MetadataInsertWithDefaults(ctx, q, pub).Scan(&pub))

		priv := NewMetadata(new(int160.Random()), MetadataOptionDescription("private"))
		priv.Private = true
		priv.Seeding = true
		require.NoError(t, MetadataInsertWithDefaults(ctx, q, priv).Scan(&priv))

		var results []Metadata
		require.NoError(t, sqlx.ScanInto(MetadataSamplePublic(ctx, q), &results))

		require.Len(t, results, 1)
		require.Equal(t, pub.ID, results[0].ID)
	})

	t.Run("excludes non-seeding metadata", func(t *testing.T) {
		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		seeding := NewMetadata(new(int160.Random()), MetadataOptionDescription("seeding"))
		seeding.Private = false
		seeding.Seeding = true
		require.NoError(t, MetadataInsertWithDefaults(ctx, q, seeding).Scan(&seeding))

		notSeeding := NewMetadata(new(int160.Random()), MetadataOptionDescription("not seeding"))
		notSeeding.Private = false
		notSeeding.Seeding = false
		require.NoError(t, MetadataInsertWithDefaults(ctx, q, notSeeding).Scan(&notSeeding))

		var results []Metadata
		require.NoError(t, sqlx.ScanInto(MetadataSamplePublic(ctx, q), &results))

		require.Len(t, results, 1)
		require.Equal(t, seeding.ID, results[0].ID)
	})

	t.Run("returns no rows when nothing qualifies", func(t *testing.T) {
		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		priv := NewMetadata(new(int160.Random()), MetadataOptionDescription("private"))
		priv.Private = true
		priv.Seeding = true
		require.NoError(t, MetadataInsertWithDefaults(ctx, q, priv).Scan(&priv))

		notSeeding := NewMetadata(new(int160.Random()), MetadataOptionDescription("not seeding"))
		notSeeding.Private = false
		notSeeding.Seeding = false
		require.NoError(t, MetadataInsertWithDefaults(ctx, q, notSeeding).Scan(&notSeeding))

		var results []Metadata
		require.NoError(t, sqlx.ScanInto(MetadataSamplePublic(ctx, q), &results))
		require.Empty(t, results)
	})

	t.Run("returns every qualifying row when under the sample cap", func(t *testing.T) {
		ctx := t.Context()
		q := sqltestx.Metadatabase(t)

		expected := map[string]struct{}{}
		for i := 0; i < 5; i++ {
			lmd := NewMetadata(new(int160.Random()), MetadataOptionDescription("public"))
			lmd.Private = false
			lmd.Seeding = true
			require.NoError(t, MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))
			expected[lmd.ID] = struct{}{}
		}

		// noise that should never appear in the results.
		priv := NewMetadata(new(int160.Random()), MetadataOptionDescription("private"))
		priv.Private = true
		priv.Seeding = true
		require.NoError(t, MetadataInsertWithDefaults(ctx, q, priv).Scan(&priv))

		notSeeding := NewMetadata(new(int160.Random()), MetadataOptionDescription("not seeding"))
		notSeeding.Private = false
		notSeeding.Seeding = false
		require.NoError(t, MetadataInsertWithDefaults(ctx, q, notSeeding).Scan(&notSeeding))

		var results []Metadata
		require.NoError(t, sqlx.ScanInto(MetadataSamplePublic(ctx, q), &results))
		require.Len(t, results, len(expected))

		for _, r := range results {
			_, ok := expected[r.ID]
			require.True(t, ok, "unexpected metadata %s in results", r.ID)
		}
	})
}
