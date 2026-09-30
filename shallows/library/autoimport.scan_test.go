package library_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqltestx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/stretchr/testify/require"
)

func TestAutoimportScan(t *testing.T) {
	t.Run("file unmodified for longer than the debounce is pending", func(t *testing.T) {
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
		require.Equal(t, "example.mkv", f.Name)
		require.Equal(t, dir.ID, f.DirectoryID)
		require.WithinDuration(t, mtime.Add(time.Hour), f.ImportAt, time.Millisecond)
	})

	t.Run("recently modified file is not pending", func(t *testing.T) {
		var (
			dir library.AutoimportDirectory
		)

		ctx, done := testx.Context(t)
		defer done()
		q := sqltestx.Metadatabase(t)
		root := t.TempDir()

		require.NoError(t, testx.Fake(&dir, library.AutoimportDirectoryOptionTestDefaults, library.AutoimportDirectoryOptionPath(root)))
		require.NoError(t, library.AutoimportDirectoryInsertWithDefaults(ctx, q, dir).Scan(&dir))
		require.NoError(t, os.WriteFile(filepath.Join(root, "example.mkv"), []byte("example"), 0600))

		require.NoError(t, library.AutoimportScan(ctx, q, dir))

		n, err := library.AutoimportFileCount(ctx, q, library.AutoimportFileQueryPending())
		require.NoError(t, err)
		require.Equal(t, uint64(0), n)

		n, err = library.AutoimportFileCount(ctx, q, library.AutoimportFileQueryByDirectoryID(dir.ID))
		require.NoError(t, err)
		require.Equal(t, uint64(1), n)
	})

	t.Run("modifying a file pushes the import out", func(t *testing.T) {
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

		modified := time.Now()
		require.NoError(t, os.Chtimes(filepath.Join(root, "example.mkv"), modified, modified))
		require.NoError(t, library.AutoimportScan(ctx, q, dir))

		f, err := sqlx.ScanOne(library.AutoimportFileSearch(ctx, q, library.AutoimportFileSearchBuilder().Where(library.AutoimportFileQueryByDirectoryID(dir.ID))))
		require.NoError(t, err)
		require.WithinDuration(t, modified.Add(time.Hour), f.ImportAt, time.Millisecond)

		n, err := library.AutoimportFileCount(ctx, q, library.AutoimportFileQueryPending())
		require.NoError(t, err)
		require.Equal(t, uint64(0), n)
	})

	t.Run("unchanged imported file is not rearmed", func(t *testing.T) {
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
		require.NoError(t, library.AutoimportFileCompleted(ctx, q, f.ID, f.LibraryMetadataID).Scan(&f))

		require.NoError(t, library.AutoimportScan(ctx, q, dir))

		n, err := library.AutoimportFileCount(ctx, q, library.AutoimportFileQueryPending())
		require.NoError(t, err)
		require.Equal(t, uint64(0), n)

		n, err = library.AutoimportFileCount(ctx, q, library.AutoimportFileQueryImported())
		require.NoError(t, err)
		require.Equal(t, uint64(1), n)
	})

	t.Run("modifying an imported file rearms it", func(t *testing.T) {
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
		require.NoError(t, library.AutoimportFileCompleted(ctx, q, f.ID, f.LibraryMetadataID).Scan(&f))

		modified := time.Now()
		require.NoError(t, os.Chtimes(filepath.Join(root, "example.mkv"), modified, modified))
		require.NoError(t, library.AutoimportScan(ctx, q, dir))

		f, err = sqlx.ScanOne(library.AutoimportFileSearch(ctx, q, library.AutoimportFileSearchBuilder().Where(library.AutoimportFileQueryByDirectoryID(dir.ID))))
		require.NoError(t, err)
		require.True(t, f.ImportAt.After(f.LastImportedAt))

		n, err := library.AutoimportFileCount(ctx, q, library.AutoimportFileQueryImported())
		require.NoError(t, err)
		require.Equal(t, uint64(0), n)
	})

	t.Run("rescan does not undo a retry backoff", func(t *testing.T) {
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

		backoff := time.Now().Add(time.Hour).Truncate(time.Microsecond)
		require.NoError(t, library.AutoimportFileFailed(ctx, q, f.ID, "example failure", backoff).Scan(&f))

		require.NoError(t, library.AutoimportScan(ctx, q, dir))

		f, err = sqlx.ScanOne(library.AutoimportFileSearch(ctx, q, library.AutoimportFileSearchBuilder().Where(library.AutoimportFileQueryByDirectoryID(dir.ID))))
		require.NoError(t, err)
		require.WithinDuration(t, backoff, f.ImportAt, time.Millisecond)
		require.Equal(t, uint32(1), f.Attempts)
		require.Equal(t, "example failure", f.Error)
	})

	t.Run("hidden, partial, and nested files are ignored", func(t *testing.T) {
		var (
			dir library.AutoimportDirectory
		)

		ctx, done := testx.Context(t)
		defer done()
		q := sqltestx.Metadatabase(t)
		root := t.TempDir()

		require.NoError(t, testx.Fake(&dir, library.AutoimportDirectoryOptionTestDefaults, library.AutoimportDirectoryOptionPath(root)))
		require.NoError(t, library.AutoimportDirectoryInsertWithDefaults(ctx, q, dir).Scan(&dir))

		require.NoError(t, os.WriteFile(filepath.Join(root, ".hidden.mkv"), []byte("example"), 0600))
		require.NoError(t, os.WriteFile(filepath.Join(root, "example.mkv.part"), []byte("example"), 0600))
		require.NoError(t, os.WriteFile(filepath.Join(root, "example.mkv.tmp"), []byte("example"), 0600))
		require.NoError(t, os.MkdirAll(filepath.Join(root, "nested"), 0700))
		require.NoError(t, os.WriteFile(filepath.Join(root, "nested", "example.mkv"), []byte("example"), 0600))

		require.NoError(t, library.AutoimportScan(ctx, q, dir))

		n, err := library.AutoimportFileCount(ctx, q, library.AutoimportFileQueryByDirectoryID(dir.ID))
		require.NoError(t, err)
		require.Equal(t, uint64(0), n)
	})

	t.Run("missing directory errors", func(t *testing.T) {
		var (
			dir library.AutoimportDirectory
		)

		ctx, done := testx.Context(t)
		defer done()
		q := sqltestx.Metadatabase(t)

		require.NoError(t, testx.Fake(&dir, library.AutoimportDirectoryOptionTestDefaults, library.AutoimportDirectoryOptionPath(filepath.Join(t.TempDir(), "missing"))))
		require.NoError(t, library.AutoimportDirectoryInsertWithDefaults(ctx, q, dir).Scan(&dir))

		require.Error(t, library.AutoimportScan(ctx, q, dir))
	})
}
