package tracking_test

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/james-lawrence/torrent/metainfo"
	"github.com/retrovibed/retrovibed/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/testx"
	"github.com/retrovibed/retrovibed/library"
	"github.com/retrovibed/retrovibed/tracking"
	"github.com/stretchr/testify/require"
)

func TestImportTorrent(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	tmpdir := t.TempDir()

	q, err := sql.Open("duckdb", filepath.Join(tmpdir, "meta.db"))
	require.NoError(t, err)
	defer q.Close()

	log.Println("duckdb version", errorsx.Zero(sqlx.String(ctx, q, "SELECT version() AS version")))

	require.NoError(t, cmdmeta.InitializeDatabase(ctx, q))
	evfs := fsx.DirVirtual(filepath.Join(tmpdir, "examples"))
	mvfs := fsx.DirVirtual(filepath.Join(tmpdir, "media"))
	tvfs := fsx.DirVirtual(filepath.Join(tmpdir, "torrents"))
	require.NoError(t, fsx.MkDirs(0700, mvfs.Path(), evfs.Path(), tvfs.Path()))

	for tx, cause := range library.ImportFilesystem(ctx, library.ImportCopyFile(mvfs), testx.Fixture()) {
		require.NoError(t, cause)

		lmd := library.NewMetadata(
			md5x.FormatString(tx.MD5),
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

	for _, err := range library.ImportFilesystem(ctx, tracking.ImportTorrent(q, mvfs, tvfs), evfs.Path()) {
		require.NoError(t, err)
		count++

		require.Equal(t, testx.ReadMD5(testx.Fixture("example.1.txt")), testx.ReadMD5(tvfs.Path(md.HashInfoBytes().HexString(), "example.1.txt")))
		require.Equal(t, testx.ReadMD5(testx.Fixture("example.2.txt")), testx.ReadMD5(tvfs.Path(md.HashInfoBytes().HexString(), "example.2.txt")))
	}

	require.Equal(t, 1, count)
	require.Equal(t, 0, testx.Must(sqlx.Count(ctx, q, "SELECT COUNT(*) FROM library_metadata WHERE torrent_id = '00000000-0000-0000-0000-000000000000'"))(t))

}
