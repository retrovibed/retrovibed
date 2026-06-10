package acoustics_test

import (
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/acoustics"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestStoreFailed(t *testing.T) {
	t.Run("failed rows are excluded from the indexed count and similarity results", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()
		db := sqltestx.Metadatabase(t)

		var vec acoustics.FeatureVector
		for i := range vec {
			vec[i] = float32(i%7) + 1
		}

		good := uuid.Must(uuid.NewV4())
		require.NoError(t, acoustics.StoreFeatures(ctx, db, good.String(), vec, acoustics.StatsVersion))

		// an undecodable file: zero vector, failed=true.
		bad := uuid.Must(uuid.NewV4())
		require.NoError(t, acoustics.StoreFailed(ctx, db, bad.String(), acoustics.StatsVersion))

		// only the good track counts toward the cold-start threshold.
		count, err := acoustics.IndexedCount(ctx, db, acoustics.StatsVersion)
		require.NoError(t, err)
		require.Equal(t, int64(1), count)

		// similarity never surfaces the failed track.
		ids, err := acoustics.SimilarMediaIDs(ctx, db, vec, nil, 10, -1.0)
		require.NoError(t, err)
		require.Contains(t, ids, good)
		require.NotContains(t, ids, bad)
	})
}
