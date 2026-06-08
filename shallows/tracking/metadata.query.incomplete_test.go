package tracking

import (
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/stretchr/testify/require"
)

func TestMetadataQueryIncomplete(t *testing.T) {
	t.Run("bytes is zero", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		lmd := NewMetadata(
			langx.Autoptr(int160.Random()),
			MetadataOptionBytes(0),
			MetadataOptionDownloaded(0),
		)
		require.NoError(t, MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		sql, args, err := MetadataSearchBuilder().RemoveColumns().Columns("COUNT(*)").Where(MetadataQueryIncomplete()).ToSql()
		require.NoError(t, err)
		require.EqualValues(t, 1, sqltestx.Count(t, q, sql, args...))
	})

	t.Run("downloaded less than bytes", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		lmd := NewMetadata(
			langx.Autoptr(int160.Random()),
			MetadataOptionBytes(1000),
			MetadataOptionDownloaded(500),
		)
		require.NoError(t, MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		sql, args, err := MetadataSearchBuilder().RemoveColumns().Columns("COUNT(*)").Where(MetadataQueryIncomplete()).ToSql()
		require.NoError(t, err)
		require.EqualValues(t, 1, sqltestx.Count(t, q, sql, args...))
	})

	t.Run("downloaded equals bytes", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		lmd := NewMetadata(
			langx.Autoptr(int160.Random()),
			MetadataOptionBytes(1000),
			MetadataOptionDownloaded(1000),
			MetadataOptionCompleted,
		)
		require.NoError(t, MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))

		sql, args, err := MetadataSearchBuilder().RemoveColumns().Columns("COUNT(*)").Where(MetadataQueryIncomplete()).ToSql()
		require.NoError(t, err)
		require.EqualValues(t, 0, sqltestx.Count(t, q, sql, args...))
	})
}
