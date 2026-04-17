package library_test

import (
	"context"
	"crypto/md5"
	"io"
	"math/rand/v2"
	"testing"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/storage"
	"github.com/james-lawrence/torrent/torrenttest"
	"github.com/retrovibed/retrovibed/retroapi/blockcache"
	"github.com/retrovibed/retrovibed/shallows/internal/bytesx"
	"github.com/retrovibed/retrovibed/shallows/internal/cryptox"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

type downloadable struct {
	md library.Metadata
	c  io.ReaderAt
}

func (t downloadable) Download(ctx context.Context, id string, offset uint64, n uint64, dst io.Writer) error {
	// log.Println("------------------------------ download initiated", id, offset, n)
	// defer log.Println("------------------------------ download completed", id, offset, n)
	dst, err := cryptox.NewOffsetWriterChaCha20(library.MetadataChaCha8(t.md), dst, uint32(offset))
	if err != nil {
		return err
	}
	rn, err := io.Copy(dst, io.NewSectionReader(t.c, int64(offset), int64(n)))
	if err != nil {
		return err
	}

	if rn < int64(n-offset) {
		return io.ErrShortWrite
	}

	return nil
}

type failAt struct {
	d     io.ReaderAt
	cause error
	n     uint64
}

func (t failAt) ReadAt(p []byte, off int64) (n int, err error) {
	n, err = t.d.ReadAt(p, off)
	if uint64(off+int64(n)) >= t.n {
		return n, t.cause
	}
	return n, err
}

func TestDeeppoolTorrentCache(t *testing.T) {
	t.Run("failed retrievals", func(t *testing.T) {
		const (
			bsize   = bytesx.KiB
			tlength = 128 * bytesx.KiB
			flength = 16 * bytesx.KiB
		)

		db := sqltestx.Metadatabase(t)

		cacheddir := fsx.DirVirtual(t.TempDir())
		vfs := fsx.DirVirtual(t.TempDir())
		localstore := blockcache.NewTorrentFromVirtualFS(vfs, blockcache.OptionTorrentCacheStorageBlockLength(bsize))
		encryption := md5x.FormatUUID(md5x.Digest(t.Name()))

		mi, expected, err := torrenttest.Seeded(cacheddir.Path(), tlength, cryptox.NewChaCha8(t.Name()))
		require.NoError(t, err)
		require.EqualValues(t, tlength, mi.TotalLength())

		md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(cacheddir.Path())), torrent.OptionDisplayName(mi.Name))
		require.NoError(t, err)

		srct, err := storage.NewFile(cacheddir.Path()).OpenTorrent(mi, md.ID)
		require.NoError(t, err)

		actual := md5.New()

		n, err := io.Copy(actual, io.NewSectionReader(srct, 0, mi.TotalLength()))
		require.NoError(t, err)
		require.EqualValues(t, tlength, n)
		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		var (
			lmd library.Metadata
			tmd tracking.Metadata
		)

		err = tracking.MetadataInsertWithDefaults(
			t.Context(),
			db,
			tracking.NewMetadata(
				langx.Autoptr(md.ID),
				tracking.MetadataOptionFromInfo(mi),
			),
		).Scan(&tmd)
		require.NoError(t, err)

		err = library.MetadataInsertWithDefaults(
			t.Context(),
			db,
			library.NewMetadata(
				md5x.FormatUUID(expected),
				library.MetadataOptionBytes(uint64(mi.TotalLength())),
				library.MetadataOptionTorrentID(tmd.ID),
				library.MetadataOptionArchiveID(md5x.FormatUUID(expected)),
				library.MetadataOptionEncryptionSeed(encryption),
			),
		).Scan(&lmd)
		require.NoError(t, err)

		const (
			errFailed = errorsx.String("failed hard")
		)
		// local storage is empty, so all data must come from the downloadable cache
		storage := library.NewTorrentStorage(downloadable{md: lmd, c: failAt{d: srct, cause: errFailed, n: flength}}, db, localstore, bsize)
		// storage := library.NewTorrentStorage(downloadable{md: lmd, c: srct}, db, localstore, bsize)
		c, err := storage.OpenTorrent(mi, md.ID)
		require.NoError(t, err)

		var (
			buf [bytesx.KiB]byte
		)
		actual = md5.New()
		n, err = io.CopyBuffer(actual, io.NewSectionReader(c, 0, mi.TotalLength()), buf[:])
		require.ErrorIs(t, err, errFailed)
		require.EqualValues(t, 15360, n)
		require.NotEqual(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		actual = md5.New()
		n, err = io.CopyBuffer(actual, io.NewSectionReader(c, 0, mi.TotalLength()), buf[:])
		require.ErrorIs(t, err, errFailed)
		require.EqualValues(t, 31744, n) // indicative of a bug in the block cache code. this should be 15360
		require.NotEqual(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))
	})

	t.Run("failed retrievals - unrecoverable", func(t *testing.T) {
		const (
			bsize   = bytesx.KiB
			tlength = 128 * bytesx.KiB
			flength = 16 * bytesx.KiB
		)

		db := sqltestx.Metadatabase(t)

		cacheddir := fsx.DirVirtual(t.TempDir())
		vfs := fsx.DirVirtual(t.TempDir())
		localstore := blockcache.NewTorrentFromVirtualFS(vfs, blockcache.OptionTorrentCacheStorageBlockLength(bsize))
		encryption := md5x.FormatUUID(md5x.Digest(t.Name()))

		mi, expected, err := torrenttest.Seeded(cacheddir.Path(), tlength, cryptox.NewChaCha8(t.Name()))
		require.NoError(t, err)
		require.EqualValues(t, tlength, mi.TotalLength())

		md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(cacheddir.Path())), torrent.OptionDisplayName(mi.Name))
		require.NoError(t, err)

		srct, err := storage.NewFile(cacheddir.Path()).OpenTorrent(mi, md.ID)
		require.NoError(t, err)

		actual := md5.New()

		n, err := io.Copy(actual, io.NewSectionReader(srct, 0, mi.TotalLength()))
		require.NoError(t, err)
		require.EqualValues(t, tlength, n)
		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		var (
			lmd library.Metadata
			tmd tracking.Metadata
		)

		err = tracking.MetadataInsertWithDefaults(
			t.Context(),
			db,
			tracking.NewMetadata(
				langx.Autoptr(md.ID),
				tracking.MetadataOptionFromInfo(mi),
			),
		).Scan(&tmd)
		require.NoError(t, err)

		err = library.MetadataInsertWithDefaults(
			t.Context(),
			db,
			library.NewMetadata(
				md5x.FormatUUID(expected),
				library.MetadataOptionBytes(uint64(mi.TotalLength())),
				library.MetadataOptionTorrentID(tmd.ID),
				library.MetadataOptionArchiveID(md5x.FormatUUID(expected)),
				library.MetadataOptionEncryptionSeed(encryption),
			),
		).Scan(&lmd)
		require.NoError(t, err)

		const (
			errFailed = errorsx.String("failed hard")
		)

		// local storage is empty, so all data must come from the downloadable cache
		tstorage := library.NewTorrentStorage(downloadable{md: lmd, c: failAt{d: srct, cause: errorsx.NewUnrecoverable(errFailed), n: flength}}, db, localstore, bsize)
		c, err := tstorage.OpenTorrent(mi, md.ID)
		require.NoError(t, err)

		var (
			buf [bytesx.KiB]byte
		)

		actual = md5.New()
		n, err = io.CopyBuffer(actual, io.NewSectionReader(c, 0, mi.TotalLength()), buf[:])
		require.ErrorIs(t, err, errFailed)
		require.ErrorIs(t, err, errorsx.NewUnrecoverable(nil)) // it must be unrecoverable so we can properly react.
		require.EqualValues(t, 15360, n)
		require.NotEqual(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		actual = md5.New()
		n, err = io.CopyBuffer(actual, io.NewSectionReader(c, 0, mi.TotalLength()), buf[:])
		require.ErrorIs(t, err, errorsx.NewUnrecoverable(nil)) // it must be unrecoverable so we can properly react.
		require.ErrorIs(t, err, library.ErrCacheUnrecoverable) // it must match the expected error.
		require.EqualValues(t, 31744, n)                       // indicative of a bug in the block cache code. this should be 15360
		require.NotEqual(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))
	})

	t.Run("restore from cache less than block", func(t *testing.T) {
		const (
			bsize   = 64
			tlength = 63
		)

		db := sqltestx.Metadatabase(t)

		cacheddir := fsx.DirVirtual(t.TempDir())
		vfs := fsx.DirVirtual(t.TempDir())
		localstore := blockcache.NewTorrentFromVirtualFS(vfs, blockcache.OptionTorrentCacheStorageBlockLength(bsize))
		encryption := md5x.FormatUUID(md5x.Digest(t.Name()))

		mi, expected, err := torrenttest.Seeded(cacheddir.Path(), tlength, cryptox.NewChaCha8(t.Name()))
		require.NoError(t, err)
		md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(cacheddir.Path())), torrent.OptionDisplayName(mi.Name))
		require.NoError(t, err)

		srct, err := storage.NewFile(cacheddir.Path()).OpenTorrent(mi, md.ID)
		require.NoError(t, err)

		actual := md5.New()

		n, err := io.Copy(actual, io.NewSectionReader(srct, 0, mi.TotalLength()))
		require.NoError(t, err)
		require.EqualValues(t, tlength, n)
		require.EqualValues(t, tlength, mi.TotalLength())
		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		var (
			lmd library.Metadata
			tmd tracking.Metadata
		)

		err = tracking.MetadataInsertWithDefaults(
			t.Context(),
			db,
			tracking.NewMetadata(
				langx.Autoptr(md.ID),
				tracking.MetadataOptionFromInfo(mi),
			),
		).Scan(&tmd)
		require.NoError(t, err)

		err = library.MetadataInsertWithDefaults(
			t.Context(),
			db,
			library.NewMetadata(
				md5x.FormatUUID(expected),
				library.MetadataOptionBytes(uint64(mi.TotalLength())),
				library.MetadataOptionTorrentID(tmd.ID),
				library.MetadataOptionArchiveID(md5x.FormatUUID(expected)),
				library.MetadataOptionEncryptionSeed(encryption),
			),
		).Scan(&lmd)
		require.NoError(t, err)

		// local storage is empty, so all data must come from the downloadable cache
		storage := library.NewTorrentStorage(downloadable{md: lmd, c: srct}, db, localstore, bsize)

		c, err := storage.OpenTorrent(mi, md.ID)
		require.NoError(t, err)

		actual = md5.New()
		n, err = io.Copy(actual, io.NewSectionReader(c, 0, mi.TotalLength()))
		require.NoError(t, err)
		require.EqualValues(t, mi.TotalLength(), n)
		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))
	})

	t.Run("restore from cache", func(t *testing.T) {
		const (
			bsize   = 64
			tlength = bsize*2 + 1
		)
		var (
			b [32]byte
		)
		db := sqltestx.Metadatabase(t)

		cacheddir := fsx.DirVirtual(t.TempDir())

		vfs := fsx.DirVirtual(t.TempDir())
		localstore := blockcache.NewTorrentFromVirtualFS(vfs, blockcache.OptionTorrentCacheStorageBlockLength(bsize))
		encryption := md5x.FormatUUID(md5x.Digest(t.Name()))

		mi, expected, err := torrenttest.Seeded(cacheddir.Path(), tlength, cryptox.NewChaCha8(t.Name()))
		require.NoError(t, err)
		md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(cacheddir.Path())), torrent.OptionDisplayName(mi.Name))
		require.NoError(t, err)

		srct, err := storage.NewFile(cacheddir.Path()).OpenTorrent(mi, md.ID)
		require.NoError(t, err)

		actual := md5.New()

		n, err := io.Copy(actual, io.NewSectionReader(srct, 0, mi.TotalLength()))
		require.NoError(t, err)
		require.EqualValues(t, tlength, n)
		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		var (
			lmd library.Metadata
			tmd tracking.Metadata
		)

		err = tracking.MetadataInsertWithDefaults(
			t.Context(),
			db,
			tracking.NewMetadata(
				langx.Autoptr(md.ID),
				tracking.MetadataOptionFromInfo(mi),
			),
		).Scan(&tmd)
		require.NoError(t, err)
		require.EqualValues(t, tlength, tmd.Bytes)

		err = library.MetadataInsertWithDefaults(
			t.Context(),
			db,
			library.NewMetadata(
				md5x.FormatUUID(expected),
				library.MetadataOptionBytes(uint64(mi.TotalLength())),
				library.MetadataOptionTorrentID(tmd.ID),
				library.MetadataOptionArchiveID(md5x.FormatUUID(expected)),
				library.MetadataOptionEncryptionSeed(encryption),
			),
		).Scan(&lmd)
		require.NoError(t, err)

		// local storage is empty, so all data must come from the downloadable cache
		storage := library.NewTorrentStorage(downloadable{md: lmd, c: srct}, db, localstore, bsize)

		c, err := storage.OpenTorrent(mi, md.ID)
		require.NoError(t, err)

		actual = md5.New()
		n, err = io.CopyBuffer(actual, io.NewSectionReader(c, 0, mi.TotalLength()), b[:])
		require.NoError(t, err)
		require.EqualValues(t, mi.TotalLength(), n)
		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))
	})

	t.Run("restore from cache random", func(t *testing.T) {
		const (
			bsize = 64
		)
		var (
			tlength = rand.Uint64N(128 * bytesx.KiB)
			b       [32]byte
		)
		db := sqltestx.Metadatabase(t)

		cacheddir := fsx.DirVirtual(t.TempDir())

		vfs := fsx.DirVirtual(t.TempDir())
		localstore := blockcache.NewTorrentFromVirtualFS(vfs, blockcache.OptionTorrentCacheStorageBlockLength(bsize))
		encryption := md5x.FormatUUID(md5x.Digest(t.Name()))

		mi, expected, err := torrenttest.Seeded(cacheddir.Path(), tlength, cryptox.NewChaCha8(t.Name()))
		require.NoError(t, err)
		md, err := torrent.NewFromInfo(mi, torrent.OptionStorage(storage.NewFile(cacheddir.Path())), torrent.OptionDisplayName(mi.Name))
		require.NoError(t, err)

		srct, err := storage.NewFile(cacheddir.Path()).OpenTorrent(mi, md.ID)
		require.NoError(t, err)

		actual := md5.New()

		n, err := io.Copy(actual, io.NewSectionReader(srct, 0, mi.TotalLength()))
		require.NoError(t, err)
		require.EqualValues(t, tlength, n)
		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))

		var (
			lmd library.Metadata
			tmd tracking.Metadata
		)

		err = tracking.MetadataInsertWithDefaults(
			t.Context(),
			db,
			tracking.NewMetadata(
				langx.Autoptr(md.ID),
				tracking.MetadataOptionFromInfo(mi),
			),
		).Scan(&tmd)
		require.NoError(t, err)
		require.EqualValues(t, tlength, tmd.Bytes)

		err = library.MetadataInsertWithDefaults(
			t.Context(),
			db,
			library.NewMetadata(
				md5x.FormatUUID(expected),
				library.MetadataOptionBytes(uint64(mi.TotalLength())),
				library.MetadataOptionTorrentID(tmd.ID),
				library.MetadataOptionArchiveID(md5x.FormatUUID(expected)),
				library.MetadataOptionEncryptionSeed(encryption),
			),
		).Scan(&lmd)
		require.NoError(t, err)

		// local storage is empty, so all data must come from the downloadable cache
		storage := library.NewTorrentStorage(downloadable{md: lmd, c: srct}, db, localstore, bsize)

		c, err := storage.OpenTorrent(mi, md.ID)
		require.NoError(t, err)

		actual = md5.New()
		n, err = io.CopyBuffer(actual, io.NewSectionReader(c, 0, mi.TotalLength()), b[:])
		require.NoError(t, err)
		require.EqualValues(t, mi.TotalLength(), n)
		require.Equal(t, md5x.FormatUUID(expected), md5x.FormatUUID(actual))
	})
}
