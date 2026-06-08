package tracking_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/james-lawrence/torrent/dht/int160"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
	"github.com/stretchr/testify/require"
)

func TestImportSymlink(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	tmpdir := t.TempDir()
	id := int160.Random()

	srcvfs := fsx.DirVirtual(filepath.Join(tmpdir, "src"))
	vfs := fsx.DirVirtual(filepath.Join(tmpdir, "media"))
	require.NoError(t, fsx.MkDirs(0700, srcvfs.Path(), vfs.Path()))

	count := 0
	for tx, err := range library.ImportFilesystem(ctx, tracking.ImportSymlink(id, srcvfs, vfs), os.DirFS(testx.Fixture()).(fs.StatFS), ".") {
		require.NoError(t, err)
		uid := md5x.FormatUUID(tx.MD5)
		target, lerr := os.Readlink(vfs.Path(uid))
		require.NoError(t, lerr)
		require.Equal(t, srcvfs.Path(id.String()), target)
		count++
	}

	require.Equal(t, 2, count)
}

func TestImportTorrent(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	tmpdir := t.TempDir()

	q := sqltestx.Metadatabase(t)

	evfs := fsx.DirVirtual(filepath.Join(tmpdir, "examples"))
	mvfs := fsx.DirVirtual(filepath.Join(tmpdir, "media"))
	tvfs := fsx.DirVirtual(filepath.Join(tmpdir, "torrents"))
	require.NoError(t, fsx.MkDirs(0700, mvfs.Path(), evfs.Path(), tvfs.Path()))

	for tx, cause := range library.ImportFilesystem(ctx, library.ImportCopyFile(mvfs), os.DirFS(testx.Fixture()).(fs.StatFS), ".") {
		require.NoError(t, cause)

		lmd := library.NewMetadata(
			md5x.FormatUUID(tx.MD5),
			library.MetadataOptionDescription(filepath.Base(tx.Path)),
			library.MetadataOptionBytes(tx.Bytes),
			library.MetadataOptionMimetype(tx.Mimetype.String()),
		)

		require.NoError(t, library.MetadataInsertWithDefaults(ctx, q, lmd).Scan(&lmd))
	}

	require.Equal(t, 2, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_metadata"))(t))
	info := testx.Must(metainfo.NewFromPath(testx.Fixture()))(t)

	md := metainfo.MetaInfo{
		InfoBytes:    testx.Must(metainfo.Encode(info))(t),
		CreationDate: time.Now().Unix(),
	}

	require.NoError(t, os.WriteFile(evfs.Path("example.torrent"), testx.Must(metainfo.Encode(md))(t), 0600))
	count := 0
	for _, err := range library.ImportFilesystem(ctx, tracking.ImportTorrent(q, mvfs, tvfs), os.DirFS(evfs.Path()).(fs.StatFS), ".") {
		require.NoError(t, err)
		count++
		require.Equal(t, testx.ReadMD5(testx.Fixture("example.1.txt")), testx.ReadMD5(tvfs.Path(md.HashInfoBytes().String(), "example.1.txt")))
		require.Equal(t, testx.ReadMD5(testx.Fixture("example.2.txt")), testx.ReadMD5(tvfs.Path(md.HashInfoBytes().String(), "example.2.txt")))
	}

	require.Equal(t, 1, count)
	require.Equal(t, 0, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_metadata WHERE torrent_id = '00000000-0000-0000-0000-000000000000'"))(t))
}
