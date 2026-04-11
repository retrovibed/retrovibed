package media_test

import (
	"io"
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"

	"github.com/retrovibed/retrovibed/blockcache"
	"github.com/retrovibed/retrovibed/internal/cryptox"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/internal/timex"
	"github.com/retrovibed/retrovibed/library"
	"github.com/retrovibed/retrovibed/media"
	"github.com/stretchr/testify/require"
)

func TestValidateTorrent(t *testing.T) {
	t.Run("marks torrent as completed and seeding when data is valid", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		db := sqltestx.Metadatabase(t)
		mvfs := fsx.DirVirtual(t.TempDir())
		tvfs := fsx.DirVirtual(t.TempDir())

		var lmd library.Metadata
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, lmd).Scan(&lmd))

		src, err := blockcache.NewDirectoryCache(mvfs.Path(lmd.ID))
		require.NoError(t, err)
		_, err = io.Copy(io.NewOffsetWriter(src, 0), io.LimitReader(cryptox.NewChaCha8(t.Name()), int64(lmd.Bytes)))
		require.NoError(t, err)

		tmd, err := media.GenerateTorrent(ctx, db, mvfs, tvfs, &lmd)
		require.NoError(t, err)
		require.True(t, tmd.Seeding)
		require.NotEqual(t, timex.Inf(), tmd.CompletedAt)

		// reset to simulate an incomplete torrent
		tmd.Seeding = false
		tmd.CompletedAt = timex.Inf()

		err = media.ValidateTorrent(ctx, db, tvfs, &tmd)
		require.NoError(t, err)
		require.True(t, tmd.Seeding)
		require.NotEqual(t, timex.Inf(), tmd.CompletedAt)
	})

	t.Run("returns error when piece hash does not match", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		db := sqltestx.Metadatabase(t)
		mvfs := fsx.DirVirtual(t.TempDir())
		tvfs := fsx.DirVirtual(t.TempDir())

		var lmd library.Metadata
		require.NoError(t, testx.Fake(&lmd, library.MetadataOptionTestDefaults))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, lmd).Scan(&lmd))

		src, err := blockcache.NewDirectoryCache(mvfs.Path(lmd.ID))
		require.NoError(t, err)
		_, err = io.Copy(io.NewOffsetWriter(src, 0), io.LimitReader(cryptox.NewChaCha8(t.Name()), int64(lmd.Bytes)))
		require.NoError(t, err)

		tmd, err := media.GenerateTorrent(ctx, db, mvfs, tvfs, &lmd)
		require.NoError(t, err)

		// corrupt the torrent data directory with different content
		dst, err := blockcache.NewDirectoryCache(tvfs.Path(int160.FromBytes(tmd.Infohash).String()))
		require.NoError(t, err)
		_, err = io.Copy(io.NewOffsetWriter(dst, 0), io.LimitReader(cryptox.NewChaCha8("different-seed"), int64(lmd.Bytes)))
		require.NoError(t, err)

		err = media.ValidateTorrent(ctx, db, tvfs, &tmd)
		require.Error(t, err)
	})
}
