package library_test

import (
	"context"
	"crypto/md5"
	"io"
	"io/fs"
	"strings"
	"testing"

	"github.com/retrovibed/retrovibed/blockcache"
	"github.com/retrovibed/retrovibed/internal/cryptox"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/library"
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
	})
}
