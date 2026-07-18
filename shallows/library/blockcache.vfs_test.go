package library_test

import (
	"context"
	"crypto/md5"
	"hash"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestVStorageFS(t *testing.T) {
	t.Run("Open method", func(t *testing.T) {
		t.Run("read data from storage", func(t *testing.T) {
			ctx, done := testx.Context(t)
			defer done()
			db := sqltestx.Metadatabase(t)
			var (
				md       library.Metadata
				expected = md5.New()
			)

			storage := fsx.DirVirtual(t.TempDir())
			require.NoError(t, testx.Fake(&md, library.MetadataOptionTestDefaults))
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, md).Scan(&md))

			dcache, err := blockcache.NewDirectoryCache(storage.Path(md.ID))
			require.NoError(t, err)

			n, err := io.Copy(io.NewOffsetWriter(dcache, 0), io.TeeReader(io.LimitReader(cryptox.NewChaCha8(t.Name()), int64(md.Bytes)), expected))
			require.NoError(t, err)
			require.Equal(t, md.Bytes, uint64(n))

			fsys := library.New(nil, storage, func(ctx context.Context, s string) (*library.Metadata, error) {
				var (
					md library.Metadata
				)
				return &md, library.MetadataFindByID(ctx, db, strings.TrimPrefix(s, "m/")).Scan(&md)
			})

			file, err := fsys.Open(md.ID)
			require.NoError(t, err)
			require.NotNil(t, file)

			// Verify the returned file's properties via its Stat() method
			fileInfo, err := file.Stat()
			require.NoError(t, err)
			require.NotNil(t, fileInfo)

			require.Equal(t, md.ID, fileInfo.Name())
			require.EqualValues(t, md.Bytes, fileInfo.Size())
			require.Equal(t, md.CreatedAt, fileInfo.ModTime())
			require.Equal(t, fs.FileMode(0600), fileInfo.Mode(), "File.Mode should be 0600 for regular file")
			require.False(t, fileInfo.IsDir())
			require.Nil(t, fileInfo.Sys())

			// ensure we can read the data
			require.Equal(t, md5x.FormatUUID(expected), testx.IOMD5(file))

			// Ensure the file can be closed
			require.NoError(t, file.Close(), "File.Close should not return an error")
		})

		t.Run("read from proper offset for a multifile block cache", func(t *testing.T) {
			ctx, done := testx.Context(t)
			defer done()
			db := sqltestx.Metadatabase(t)
			var (
				md0       library.Metadata
				md1       library.Metadata
				tid       = uuid.Nil.String()
				expected0 = md5.New()
				expected1 = md5.New()
			)

			mediastorage := fsx.DirVirtual(t.TempDir())
			cachestorage := fsx.DirVirtual(t.TempDir())
			require.NoError(t, testx.Fake(
				&md0,
				library.MetadataOptionTestDefaults,
				library.MetadataOptionDescription("first"),
				library.MetadataOptionTorrentID(tid),
			))
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, md0).Scan(&md0))
			require.EqualValues(t, 16*bytesx.KiB, md0.Bytes)

			require.NoError(t, testx.Fake(
				&md1,
				library.MetadataOptionTestDefaults,
				library.MetadataOptionTestID(uuidx.WithSuffix(1)),
				library.MetadataOptionDescription("second"),
				library.MetadataOptionTorrentID(tid),
				library.MetadataOptionOffset(md0.Bytes),
			))
			require.EqualValues(t, md0.Bytes, md1.DiskOffset)
			require.NoError(t, library.MetadataInsertWithDefaults(ctx, db, md1).Scan(&md1))
			require.EqualValues(t, md0.Bytes, md1.DiskOffset)
			require.EqualValues(t, uuidx.WithSuffix(1), md1.ID)

			dcache, err := blockcache.NewDirectoryCache(cachestorage.Path(md0.TorrentID))
			require.NoError(t, err)

			prng := cryptox.NewChaCha8(t.Name())
			n, err := io.Copy(io.NewOffsetWriter(dcache, 0), io.TeeReader(io.LimitReader(prng, int64(md0.Bytes)), expected0))
			require.NoError(t, err)
			require.Equal(t, md0.Bytes, uint64(n))

			n, err = io.Copy(io.NewOffsetWriter(dcache, int64(md1.DiskOffset)), io.TeeReader(io.LimitReader(prng, int64(md1.Bytes)), expected1))
			require.NoError(t, err)
			require.Equal(t, md1.Bytes, uint64(n))

			// symlink the media storage to cache storage
			log.Println("SYMLINKED", cachestorage.Path(md0.TorrentID), mediastorage.Path(md0.ID))
			require.NoError(t, os.Symlink(cachestorage.Path(md0.TorrentID), mediastorage.Path(md0.ID)))
			log.Println("SYMLINKED", cachestorage.Path(md0.TorrentID), mediastorage.Path(md1.ID))
			require.NoError(t, os.Symlink(cachestorage.Path(md0.TorrentID), mediastorage.Path(md1.ID)))

			fsys := library.New(nil, mediastorage, func(ctx context.Context, s string) (*library.Metadata, error) {
				var (
					md library.Metadata
				)
				return &md, library.MetadataFindByID(ctx, db, strings.TrimPrefix(s, "m/")).Scan(&md)
			})

			check := func(lmd library.Metadata, expected hash.Hash) {
				log.Println("WAKA WAKA", lmd.ID)
				file, err := fsys.Open(lmd.ID)
				require.NoError(t, err)
				require.NotNil(t, file)

				// Verify the returned file's properties via its Stat() method
				fileInfo, err := file.Stat()
				require.NoError(t, err)
				require.NotNil(t, fileInfo)

				require.Equal(t, lmd.ID, fileInfo.Name())
				require.EqualValues(t, lmd.Bytes, fileInfo.Size())
				require.Equal(t, lmd.CreatedAt, fileInfo.ModTime())
				require.Equal(t, fs.FileMode(0600), fileInfo.Mode(), "File.Mode should be 0600 for regular file")
				require.False(t, fileInfo.IsDir())
				require.Nil(t, fileInfo.Sys())

				// ensure we can read the data
				require.Equal(t, md5x.FormatUUID(expected), testx.IOMD5(file), "%s md5 mismatch", lmd.Description)

				// Ensure the file can be closed
				require.NoError(t, file.Close(), "File.Close should not return an error")
			}

			check(md0, expected0)
			check(md1, expected1)
		})
	})
}
