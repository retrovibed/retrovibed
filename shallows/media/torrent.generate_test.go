package media_test

import (
	"crypto/md5"
	"io"
	"os"
	"testing"

	"github.com/james-lawrence/torrent/dht/int160"

	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/tracking"
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
		require.EqualValues(t, lmd.Bytes, tmd.Available)
		require.Zero(t, tmd.Downloaded, "self-published content was never fetched from peers, so downloaded must be 0")
		require.Equal(t, lmd.KnownMediaID, tmd.KnownMediaID)
		require.True(t, tmd.Seeding)
		require.NotEmpty(t, tmd.Infohash)
		require.Equal(t, timex.Inf(), tmd.InitiatedAt)
		require.NotEqual(t, timex.Inf(), tmd.CompletedAt)

		_, err = os.Stat(tvfs.Path(infohashHex + tracking.TorrentSuffix))
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
