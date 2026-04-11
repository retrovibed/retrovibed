package media_test

import (
	"crypto/md5"
	"io"
	"os"
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"

	"github.com/retrovibed/retrovibed/blockcache"
	"github.com/retrovibed/retrovibed/internal/cryptox"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/internal/timex"
	"github.com/retrovibed/retrovibed/library"
	"github.com/retrovibed/retrovibed/media"
	"github.com/stretchr/testify/require"
)

func TestGenerateTorrent(t *testing.T) {
	t.Run("generates torrent metadata and filesystem layout", func(t *testing.T) {
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

		infohashHex := int160.FromBytes(tmd.Infohash).String()

		require.Equal(t, 1, testx.Must(sqlx.Count(ctx, db, "SELECT COUNT(*) FROM torrents_metadata"))(t))
		require.Equal(t, tmd.ID, lmd.TorrentID)
		require.EqualValues(t, lmd.Bytes, tmd.Bytes)
		require.EqualValues(t, lmd.Bytes, tmd.Downloaded)
		require.Equal(t, lmd.KnownMediaID, tmd.KnownMediaID)
		require.True(t, tmd.Seeding)
		require.NotEmpty(t, tmd.Infohash)
		require.Equal(t, timex.Inf(), tmd.InitiatedAt)
		require.NotEqual(t, timex.Inf(), tmd.CompletedAt)

		_, err = os.Stat(tvfs.Path(infohashHex + ".torrent"))
		require.NoError(t, err)

		link, err := os.Readlink(mvfs.Path(lmd.ID))
		require.NoError(t, err)
		require.Equal(t, tvfs.Path(infohashHex), link)
	})

	t.Run("content is faithfully copied to the torrent directory", func(t *testing.T) {
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

		expected := md5.New()
		_, err = io.Copy(expected, io.NewSectionReader(src, 0, int64(lmd.Bytes)))
		require.NoError(t, err)

		tmd, err := media.GenerateTorrent(ctx, db, mvfs, tvfs, &lmd)
		require.NoError(t, err)

		dst, err := blockcache.NewDirectoryCache(tvfs.Path(int160.FromBytes(tmd.Infohash).String()))
		require.NoError(t, err)

		actual := md5.New()
		_, err = io.Copy(actual, io.NewSectionReader(dst, 0, int64(lmd.Bytes)))
		require.NoError(t, err)

		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))
	})
}
