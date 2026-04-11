package library_test

import (
	"context"
	"crypto/md5"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/blockcache"
	"github.com/retrovibed/retrovibed/deeppool"
	"github.com/retrovibed/retrovibed/internal/asyncx"
	"github.com/retrovibed/retrovibed/internal/bytesx"
	"github.com/retrovibed/retrovibed/internal/cryptox"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"github.com/retrovibed/retrovibed/internal/sqltestx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/internal/uuidx"
	"github.com/retrovibed/retrovibed/library"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type md5archive struct {
	seed   []byte
	result *deeppool.Media
}

func (t *md5archive) Upload(ctx context.Context, mimetype string, r io.Reader) (*deeppool.Media, error) {
	hasher := md5.New()

	decrypt, err := cryptox.NewWriterChaCha20(cryptox.NewChaCha8(t.seed), hasher)
	if err != nil {
		return nil, err
	}
	n, err := io.Copy(decrypt, r)
	if err != nil {
		return nil, err
	}

	t.result = &deeppool.Media{
		Id:        md5x.FormatUUID(hasher),
		Mimetype:  mimetype,
		Bytes:     uint64(n),
		Usage:     uint64(n),
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	return t.result, nil
}

func TestArchive(t *testing.T) {
	t.Run("should successfully archive media with valid metadata", func(t *testing.T) {
		var (
			expected = md5.New()
			v        library.Metadata
		)
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&v, library.MetadataOptionTestDefaults, func(md *library.Metadata) {
			md.Bytes = 16 * bytesx.KiB
		}))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, v).Scan(&v))

		root := fsx.DirVirtual(filepath.Join(t.TempDir()))

		bcache, err := blockcache.NewDirectoryCache(root.Path(v.ID))
		require.NoError(t, err)

		n, err := io.Copy(
			io.NewOffsetWriter(bcache, 0),
			io.TeeReader(
				io.LimitReader(cryptox.NewChaCha8(v.ID), int64(v.Bytes)),
				expected,
			),
		)
		require.NoError(t, err)
		require.Equal(t, v.Bytes, uint64(n))

		archiver := &md5archive{
			seed: uuidx.FirstNonNil(uuid.FromStringOrNil(v.EncryptionSeed), uuid.FromStringOrNil(v.ID)).Bytes(),
		}

		require.NoError(t, library.Archive(ctx, q, &v, root, archiver))
		assert.Equal(t, md5x.FormatUUID(expected), archiver.result.Id)
		assert.Equal(t, v.Bytes, archiver.result.Bytes)
		assert.Equal(t, v.Mimetype, archiver.result.Mimetype)
	})
}

func TestReclaimEndChunks(t *testing.T) {
	t.Run("should remove one chunk (the last existing chunk)", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tempDir := t.TempDir()
		root := fsx.DirVirtual(tempDir)
		var v library.Metadata
		require.NoError(t, testx.Fake(&v, library.MetadataOptionTestDefaults, func(md *library.Metadata) {
			md.Bytes = 160 * bytesx.MiB
		}))

		storageDir := filepath.Join(tempDir, "storage", v.ID)
		require.NoError(t, os.MkdirAll(storageDir, 0755))
		require.NoError(t, os.Symlink(storageDir, root.Path(v.ID)))

		bcache, err := blockcache.NewDirectoryCache(storageDir)
		require.NoError(t, err)

		totalChunks := int64(math.Ceil(float64(v.Bytes) / float64(blockcache.DefaultBlockLength)))

		for i := int64(0); i < totalChunks; i++ {
			n, err := io.Copy(
				io.NewOffsetWriter(bcache, i*blockcache.DefaultBlockLength),
				io.LimitReader(cryptox.NewChaCha8(v.ID), blockcache.DefaultBlockLength),
			)
			require.NoError(t, err)
			require.True(t, n > 0)
		}

		_, err = library.ReclaimEndChunks(ctx, v, root)
		require.NoError(t, err)

		// Only the last chunk should be removed
		lastChunkPath := filepath.Join(storageDir, strconv.FormatInt(totalChunks-1, 10))
		assert.False(t, fsx.Exists(lastChunkPath), "last chunk should be removed")

		// All other chunks should still exist
		for i := int64(0); i < totalChunks-1; i++ {
			chunkPath := filepath.Join(storageDir, strconv.FormatInt(i, 10))
			assert.True(t, fsx.Exists(chunkPath), "chunk %d should still exist", i)
		}
	})

	t.Run("should remove last chunk even for single-chunk files", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		tempDir := t.TempDir()
		root := fsx.DirVirtual(tempDir)
		var v library.Metadata
		require.NoError(t, testx.Fake(&v, library.MetadataOptionTestDefaults, func(md *library.Metadata) {
			md.Bytes = 16 * bytesx.MiB
		}))

		storageDir := filepath.Join(tempDir, "storage", v.ID)
		require.NoError(t, os.MkdirAll(storageDir, 0755))
		require.NoError(t, os.Symlink(storageDir, root.Path(v.ID)))

		bcache, err := blockcache.NewDirectoryCache(storageDir)
		require.NoError(t, err)

		n, err := io.Copy(
			io.NewOffsetWriter(bcache, 0),
			io.LimitReader(cryptox.NewChaCha8(v.ID), int64(v.Bytes)),
		)
		require.NoError(t, err)
		require.Equal(t, v.Bytes, uint64(n))

		_, err = library.ReclaimEndChunks(ctx, v, root)
		require.NoError(t, err)

		chunkPath := filepath.Join(storageDir, "0")
		assert.False(t, fsx.Exists(chunkPath), "single chunk should be removed")
	})
}

func TestNewSlowDiskReclaim(t *testing.T) {
	t.Run("should process archived files and perform chunk reclamation", func(t *testing.T) {
		ctx, done := testx.Context(t)
		defer done()

		q := sqltestx.Metadatabase(t)

		tempDir := t.TempDir()
		root := fsx.DirVirtual(tempDir)

		var v library.Metadata
		require.NoError(t, testx.Fake(&v, library.MetadataOptionTestDefaults, func(md *library.Metadata) {
			md.Bytes = 128 * bytesx.MiB
			md.ArchiveID = uuid.Must(uuid.NewV4()).String()
		}))
		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, v).Scan(&v))
		storageDir := filepath.Join(tempDir, "storage", v.ID)
		require.NoError(t, os.MkdirAll(storageDir, 0755))
		require.NoError(t, os.Symlink(storageDir, root.Path(v.ID)))

		bcache, err := blockcache.NewDirectoryCache(storageDir)
		require.NoError(t, err)

		numchunks := int64(math.Ceil(float64(v.Bytes) / float64(blockcache.DefaultBlockLength)))
		for i := range numchunks {
			n, err := io.Copy(
				io.NewOffsetWriter(bcache, i*blockcache.DefaultBlockLength),
				io.LimitReader(cryptox.NewChaCha8(v.ID), blockcache.DefaultBlockLength),
			)
			require.NoError(t, err)
			require.True(t, n > 0)
		}

		async := asyncx.NewWakeup(ctx)

		actx, adone := context.WithCancelCause(ctx)
		go func() {
			adone(library.NewSlowDiskReclaim(ctx, root, q, async, 0, true))
		}()
		async.Broadcast()
		require.NoError(t, async.Close())
		<-actx.Done()
		require.NoError(t, errorsx.Ignore(context.Cause(actx), context.Canceled))

		// Since we remove one chunk per pass and the test runs multiple passes,
		// the last 2 chunks should be removed (chunks 3 and 2 for a 4-chunk file)
		lastChunkPath := filepath.Join(storageDir, strconv.FormatInt(numchunks-1, 10))
		assert.False(t, fsx.Exists(lastChunkPath), "last chunk should be removed")

		secondLastChunkPath := filepath.Join(storageDir, strconv.FormatInt(numchunks-2, 10))
		assert.False(t, fsx.Exists(secondLastChunkPath), "second to last chunk should be removed")

		// The first 2 chunks should still exist
		for i := range int64(2) {
			chunkPath := filepath.Join(storageDir, strconv.FormatInt(i, 10))
			assert.True(t, fsx.Exists(chunkPath), "chunk %d should still exist", i)
		}
	})
}
