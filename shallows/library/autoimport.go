package library

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/duckdbx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

const (
	AutoimportModeCopy uint32 = 0 // keep the original file after import.
	AutoimportModeMove uint32 = 1 // remove the original file after import.
)

const (
	autoimportBackoffBase = time.Minute
	autoimportBackoffMax  = 24 * time.Hour
)

func AutoimportDirectoryOptionTestDefaults(d *AutoimportDirectory) {
	d.ID = uuid.Nil.String()
	d.Debounce = time.Hour
	d.Mode = AutoimportModeCopy
	d.LibraryDirectoryID = uuid.Nil.String()
	d.LastScannedAt = timex.NegInf()
}

func AutoimportDirectoryOptionPath(p string) func(*AutoimportDirectory) {
	return func(d *AutoimportDirectory) {
		d.Path = p
	}
}

func AutoimportDirectoryOptionMode(m uint32) func(*AutoimportDirectory) {
	return func(d *AutoimportDirectory) {
		d.Mode = m
	}
}

func AutoimportFileOptionTestDefaults(f *AutoimportFile) {
	f.ID = uuid.Nil.String()
	f.DirectoryID = uuid.Nil.String()
	f.LibraryMetadataID = uuid.Nil.String()
	f.LastImportedAt = timex.NegInf()
	f.Attempts = 0
	f.Error = ""
}

func AutoimportDirectoryQueryByIDs(ids ...string) squirrel.Sqlizer {
	if len(ids) == 0 {
		return squirrelx.Noop{}
	}
	return squirrel.Eq{"library_autoimport_directories.id": ids}
}

func AutoimportDirectoryQueryText(query string) squirrel.Sqlizer {
	if query == "" {
		return squirrelx.Noop{}
	}
	return lucenex.Query(duckdbx.NewLucene(), query, lucenex.WithDefaultField("description"))
}

func AutoimportDirectorySearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) AutoimportDirectoryScanner {
	return NewAutoimportDirectoryScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func AutoimportDirectorySearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(AutoimportDirectoryScannerStaticColumns)...).From("library_autoimport_directories")
}

func AutoimportFileQueryPending() squirrel.Sqlizer {
	return squirrel.And{
		squirrel.Expr("library_autoimport_files.import_at <= NOW()"),
		squirrel.Expr("library_autoimport_files.last_imported_at < library_autoimport_files.import_at"),
	}
}

func AutoimportFileQueryImported() squirrel.Sqlizer {
	return squirrel.Expr("library_autoimport_files.last_imported_at >= library_autoimport_files.import_at")
}

func AutoimportFileQueryByDirectoryID(id string) squirrel.Sqlizer {
	return squirrel.Eq{"library_autoimport_files.directory_id": id}
}

func AutoimportFileSearch(ctx context.Context, q sqlx.Queryer, b squirrel.SelectBuilder) AutoimportFileScanner {
	return NewAutoimportFileScannerStatic(b.RunWith(q).QueryContext(ctx))
}

func AutoimportFileSearchBuilder() squirrel.SelectBuilder {
	return squirrelx.PSQL.Select(sqlx.Columns(AutoimportFileScannerStaticColumns)...).From("library_autoimport_files")
}

func AutoimportFileCount(ctx context.Context, q sqlx.Queryer, where squirrel.Sqlizer) (n uint64, err error) {
	err = squirrelx.PSQL.Select("COUNT(*)").From("library_autoimport_files").Where(where).RunWith(q).QueryRowContext(ctx).Scan(&n)
	return n, err
}

// files still being written to (browser/torrent partials) or hidden files are never imported.
func autoimportIgnored(name string) bool {
	return strings.HasPrefix(name, ".") || strings.HasSuffix(name, ".part") || strings.HasSuffix(name, ".tmp")
}

// duckdb stores timestamps at microsecond precision, the file's modification time has
// nanosecond precision, truncate so comparisons against the stored import_at are stable.
func autoimportAt(info fs.FileInfo, debounce time.Duration) time.Time {
	return info.ModTime().Add(debounce).Truncate(time.Microsecond)
}

func autoimportBackoff(attempts uint32) time.Duration {
	if attempts >= 11 { // 2^11 minutes exceeds the max.
		return autoimportBackoffMax
	}

	return min(autoimportBackoffBase*time.Duration(1<<attempts), autoimportBackoffMax)
}

// AutoimportScan records every top level file within the directory, scheduling its import
// for debounce after it was last modified.
func AutoimportScan(ctx context.Context, q sqlx.Queryer, dir AutoimportDirectory) error {
	entries, err := os.ReadDir(dir.Path)
	if err != nil {
		return errorsx.Wrapf(err, "unable to read directory: %s", dir.Path)
	}

	for _, entry := range entries {
		if err = ctx.Err(); err != nil {
			return err
		}

		if !entry.Type().IsRegular() || autoimportIgnored(entry.Name()) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			log.Println(errorsx.Wrapf(err, "unable to stat: %s", filepath.Join(dir.Path, entry.Name())))
			continue
		}

		f := AutoimportFile{
			DirectoryID: dir.ID,
			Name:        entry.Name(),
			ImportAt:    autoimportAt(info, dir.Debounce),
		}

		// no row is returned when the file is unchanged since import_at only moves forward.
		if err = sqlx.IgnoreNoRows(AutoimportFileUpsert(ctx, q, f).Scan(&f)); err != nil {
			return errorsx.Wrapf(err, "unable to record file: %s", filepath.Join(dir.Path, entry.Name()))
		}
	}

	return errorsx.Wrap(AutoimportDirectoryScanned(ctx, q, dir.ID).Scan(&dir), "unable to record directory scan")
}

// AutoimportUploader sends a file into a library, placing it into the given library directory,
// returning the id of the resulting media.
type AutoimportUploader interface {
	Upload(ctx context.Context, path string, directoryID string) (string, error)
}

type AutoimportUploaderFunc func(ctx context.Context, path string, directoryID string) (string, error)

func (t AutoimportUploaderFunc) Upload(ctx context.Context, path string, directoryID string) (string, error) {
	return t(ctx, path, directoryID)
}

// AutoimportFileImport uploads a single pending file into the library. errors are recorded
// on the file row with a backoff, the returned error is only for failures to record state.
func AutoimportFileImport(ctx context.Context, q sqlx.Queryer, up AutoimportUploader, dir AutoimportDirectory, f AutoimportFile) (err error) {
	path := filepath.Join(dir.Path, f.Name)

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return errorsx.Wrapf(AutoimportFileDeleteByID(ctx, q, f.ID).Scan(&f), "unable to remove missing file: %s", path)
	} else if err != nil {
		return autoimportFailed(ctx, q, f, errorsx.Wrapf(err, "unable to stat: %s", path))
	}

	// modified since the last scan, reschedule instead of importing a file that may still be written to.
	if at := autoimportAt(info, dir.Debounce); at.After(f.ImportAt) {
		f.ImportAt = at
		return errorsx.Wrapf(AutoimportFileUpsert(ctx, q, f).Scan(&f), "unable to reschedule file: %s", path)
	}

	mid, err := up.Upload(ctx, path, dir.LibraryDirectoryID)
	if err != nil {
		return autoimportFailed(ctx, q, f, errorsx.Wrapf(err, "unable to upload: %s", path))
	}

	if err = AutoimportFileCompleted(ctx, q, f.ID, mid).Scan(&f); err != nil {
		return errorsx.Wrapf(err, "unable to record import: %s", path)
	}

	if dir.Mode == AutoimportModeMove {
		if err = os.Remove(path); fsx.IgnoreIsNotExist(err) != nil {
			return errorsx.Wrapf(err, "unable to remove original: %s", path)
		}
	}

	log.Println("autoimport completed", path, mid)

	return nil
}

func autoimportFailed(ctx context.Context, q sqlx.Queryer, f AutoimportFile, cause error) error {
	log.Println(cause)
	return errorsx.Wrap(AutoimportFileFailed(ctx, q, f.ID, cause.Error(), time.Now().Add(autoimportBackoff(f.Attempts)).Truncate(time.Microsecond)).Scan(&f), "unable to record import failure")
}

// NewAutoimport scans every monitored directory and then uploads the pending files.
func NewAutoimport(ctx context.Context, q sqlx.Queryer, up AutoimportUploader) (err error) {
	var (
		dirs = make(map[string]AutoimportDirectory)
	)

	diter := sqlx.Scan(AutoimportDirectorySearch(ctx, q, AutoimportDirectorySearchBuilder()))
	for dir := range diter.Iter() {
		dirs[dir.ID] = dir
	}

	if err = diter.Err(); err != nil {
		return errorsx.Wrap(err, "unable to retrieve autoimport directories")
	}

	for _, dir := range dirs {
		errorsx.Log(AutoimportScan(ctx, q, dir))
	}

	// buffer the pending files, the imports write to the same table being iterated.
	var pending []AutoimportFile
	if err = sqlx.ScanInto(AutoimportFileSearch(ctx, q, AutoimportFileSearchBuilder().Where(AutoimportFileQueryPending()).OrderBy("import_at ASC").Limit(128)), &pending); err != nil {
		return errorsx.Wrap(err, "unable to retrieve pending autoimport files")
	}

	for _, f := range pending {
		if err = ctx.Err(); err != nil {
			return err
		}

		dir, ok := dirs[f.DirectoryID]
		if !ok {
			continue
		}

		errorsx.Log(AutoimportFileImport(ctx, q, up, dir, f))
	}

	return nil
}
