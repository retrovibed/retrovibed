package library_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestAutoimportFileImport(t *testing.T) {
	t.Run("copy mode keeps the original", func(t *testing.T) {
		var (
			dir      library.AutoimportDirectory
			uploaded string
			libdir   string
		)

		ctx, done := testx.Context(t)
		defer done()
		q := sqltestx.Metadatabase(t)
		root := t.TempDir()
		mid := uuid.Must(uuid.NewV4()).String()

		require.NoError(t, testx.Fake(&dir, library.AutoimportDirectoryOptionTestDefaults, library.AutoimportDirectoryOptionPath(root), library.AutoimportDirectoryOptionMode(library.AutoimportModeCopy)))
		require.NoError(t, library.AutoimportDirectoryInsertWithDefaults(ctx, q, dir).Scan(&dir))

		mtime := time.Now().Add(-2 * time.Hour)
		require.NoError(t, os.WriteFile(filepath.Join(root, "example.mkv"), []byte("example"), 0600))
		require.NoError(t, os.Chtimes(filepath.Join(root, "example.mkv"), mtime, mtime))
		require.NoError(t, library.AutoimportScan(ctx, q, dir))

		f, err := sqlx.ScanOne(library.AutoimportFileSearch(ctx, q, library.AutoimportFileSearchBuilder().Where(library.AutoimportFileQueryPending())))
		require.NoError(t, err)

		require.NoError(t, library.AutoimportFileImport(ctx, q, library.AutoimportUploaderFunc(func(ctx context.Context, path string, directoryID string) (string, error) {
			uploaded, libdir = path, directoryID
			return mid, nil
		}), dir, f))

		require.Equal(t, filepath.Join(root, "example.mkv"), uploaded)
		require.Equal(t, dir.LibraryDirectoryID, libdir)
		require.FileExists(t, filepath.Join(root, "example.mkv"))

		f, err = sqlx.ScanOne(library.AutoimportFileSearch(ctx, q, library.AutoimportFileSearchBuilder().Where(library.AutoimportFileQueryImported())))
		require.NoError(t, err)
		require.Equal(t, mid, f.LibraryMetadataID)
	})

	t.Run("move mode removes the original", func(t *testing.T) {
		var (
			dir library.AutoimportDirectory
		)

		ctx, done := testx.Context(t)
		defer done()
		q := sqltestx.Metadatabase(t)
		root := t.TempDir()

		require.NoError(t, testx.Fake(&dir, library.AutoimportDirectoryOptionTestDefaults, library.AutoimportDirectoryOptionPath(root), library.AutoimportDirectoryOptionMode(library.AutoimportModeMove)))
		require.NoError(t, library.AutoimportDirectoryInsertWithDefaults(ctx, q, dir).Scan(&dir))

		mtime := time.Now().Add(-2 * time.Hour)
		require.NoError(t, os.WriteFile(filepath.Join(root, "example.mkv"), []byte("example"), 0600))
		require.NoError(t, os.Chtimes(filepath.Join(root, "example.mkv"), mtime, mtime))
		require.NoError(t, library.AutoimportScan(ctx, q, dir))

		f, err := sqlx.ScanOne(library.AutoimportFileSearch(ctx, q, library.AutoimportFileSearchBuilder().Where(library.AutoimportFileQueryPending())))
		require.NoError(t, err)

		require.NoError(t, library.AutoimportFileImport(ctx, q, library.AutoimportUploaderFunc(func(ctx context.Context, path string, directoryID string) (string, error) {
			return uuid.Must(uuid.NewV4()).String(), nil
		}), dir, f))

		require.NoFileExists(t, filepath.Join(root, "example.mkv"))

		n, err := library.AutoimportFileCount(ctx, q, library.AutoimportFileQueryImported())
		require.NoError(t, err)
		require.Equal(t, uint64(1), n)
	})

	t.Run("file modified since the scan is rescheduled", func(t *testing.T) {
		var (
			dir library.AutoimportDirectory
		)

		ctx, done := testx.Context(t)
		defer done()
		q := sqltestx.Metadatabase(t)
		root := t.TempDir()

		require.NoError(t, testx.Fake(&dir, library.AutoimportDirectoryOptionTestDefaults, library.AutoimportDirectoryOptionPath(root)))
		require.NoError(t, library.AutoimportDirectoryInsertWithDefaults(ctx, q, dir).Scan(&dir))

		mtime := time.Now().Add(-2 * time.Hour)
		require.NoError(t, os.WriteFile(filepath.Join(root, "example.mkv"), []byte("example"), 0600))
		require.NoError(t, os.Chtimes(filepath.Join(root, "example.mkv"), mtime, mtime))
		require.NoError(t, library.AutoimportScan(ctx, q, dir))

		f, err := sqlx.ScanOne(library.AutoimportFileSearch(ctx, q, library.AutoimportFileSearchBuilder().Where(library.AutoimportFileQueryPending())))
		require.NoError(t, err)

		modified := time.Now()
		require.NoError(t, os.Chtimes(filepath.Join(root, "example.mkv"), modified, modified))

		require.NoError(t, library.AutoimportFileImport(ctx, q, library.AutoimportUploaderFunc(func(ctx context.Context, path string, directoryID string) (string, error) {
			return "", errorsx.String("uploader should not be invoked")
		}), dir, f))

		f, err = sqlx.ScanOne(library.AutoimportFileSearch(ctx, q, library.AutoimportFileSearchBuilder().Where(library.AutoimportFileQueryByDirectoryID(dir.ID))))
		require.NoError(t, err)
		require.WithinDuration(t, modified.Add(time.Hour), f.ImportAt, time.Millisecond)
		require.Equal(t, uuid.Nil.String(), f.LibraryMetadataID)
		require.Equal(t, uint32(0), f.Attempts)

		n, err := library.AutoimportFileCount(ctx, q, library.AutoimportFileQueryImported())
		require.NoError(t, err)
		require.Equal(t, uint64(0), n)
	})

	t.Run("missing file removes the record", func(t *testing.T) {
		var (
			dir library.AutoimportDirectory
		)

		ctx, done := testx.Context(t)
		defer done()
		q := sqltestx.Metadatabase(t)
		root := t.TempDir()

		require.NoError(t, testx.Fake(&dir, library.AutoimportDirectoryOptionTestDefaults, library.AutoimportDirectoryOptionPath(root)))
		require.NoError(t, library.AutoimportDirectoryInsertWithDefaults(ctx, q, dir).Scan(&dir))

		mtime := time.Now().Add(-2 * time.Hour)
		require.NoError(t, os.WriteFile(filepath.Join(root, "example.mkv"), []byte("example"), 0600))
		require.NoError(t, os.Chtimes(filepath.Join(root, "example.mkv"), mtime, mtime))
		require.NoError(t, library.AutoimportScan(ctx, q, dir))

		f, err := sqlx.ScanOne(library.AutoimportFileSearch(ctx, q, library.AutoimportFileSearchBuilder().Where(library.AutoimportFileQueryPending())))
		require.NoError(t, err)

		require.NoError(t, os.Remove(filepath.Join(root, "example.mkv")))
		require.NoError(t, library.AutoimportFileImport(ctx, q, library.AutoimportUploaderFunc(func(ctx context.Context, path string, directoryID string) (string, error) {
			return "", errorsx.String("uploader should not be invoked")
		}), dir, f))

		var target = sql.ErrNoRows
		_, err = sqlx.ScanOne(library.AutoimportFileSearch(ctx, q, library.AutoimportFileSearchBuilder().Where(library.AutoimportFileQueryByDirectoryID(dir.ID))))
		require.ErrorAs(t, err, &target)
	})

	t.Run("upload failure records the error, backs off, and keeps the original", func(t *testing.T) {
		var (
			dir library.AutoimportDirectory
		)

		ctx, done := testx.Context(t)
		defer done()
		q := sqltestx.Metadatabase(t)
		root := t.TempDir()

		require.NoError(t, testx.Fake(&dir, library.AutoimportDirectoryOptionTestDefaults, library.AutoimportDirectoryOptionPath(root), library.AutoimportDirectoryOptionMode(library.AutoimportModeMove)))
		require.NoError(t, library.AutoimportDirectoryInsertWithDefaults(ctx, q, dir).Scan(&dir))

		mtime := time.Now().Add(-2 * time.Hour)
		require.NoError(t, os.WriteFile(filepath.Join(root, "example.mkv"), []byte("example"), 0600))
		require.NoError(t, os.Chtimes(filepath.Join(root, "example.mkv"), mtime, mtime))
		require.NoError(t, library.AutoimportScan(ctx, q, dir))

		f, err := sqlx.ScanOne(library.AutoimportFileSearch(ctx, q, library.AutoimportFileSearchBuilder().Where(library.AutoimportFileQueryPending())))
		require.NoError(t, err)

		require.NoError(t, library.AutoimportFileImport(ctx, q, library.AutoimportUploaderFunc(func(ctx context.Context, path string, directoryID string) (string, error) {
			return "", errorsx.String("example upload failure")
		}), dir, f))

		f, err = sqlx.ScanOne(library.AutoimportFileSearch(ctx, q, library.AutoimportFileSearchBuilder().Where(library.AutoimportFileQueryByDirectoryID(dir.ID))))
		require.NoError(t, err)
		require.Equal(t, uint32(1), f.Attempts)
		require.Contains(t, f.Error, "example upload failure")
		require.WithinDuration(t, time.Now().Add(time.Minute), f.ImportAt, 5*time.Second)
		require.FileExists(t, filepath.Join(root, "example.mkv"))

		n, err := library.AutoimportFileCount(ctx, q, library.AutoimportFileQueryPending())
		require.NoError(t, err)
		require.Equal(t, uint64(0), n)
	})
}
