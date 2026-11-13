package cmdopts

import (
	"io/fs"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/retrovibed/retrovibed/retroapi/testx"
	"github.com/stretchr/testify/require"
)

// cacheMigrationsBefore returns the cache migrations that sort strictly before
// the named migration, simulating a cache.db created by an older release.
func cacheMigrationsBefore(t testing.TB, name string) fs.FS {
	src, err := fs.Sub(embedmigrationscache, ".migrations.cache")
	require.NoError(t, err)

	entries, err := fs.ReadDir(src, ".")
	require.NoError(t, err)

	subset := fstest.MapFS{}
	for _, e := range entries {
		if e.Name() >= name {
			continue
		}

		b, err := fs.ReadFile(src, e.Name())
		require.NoError(t, err)
		subset[e.Name()] = &fstest.MapFile{Data: b}
	}

	return subset
}

func TestCacheMigrationSyncIDPopulatedTable(t *testing.T) {
	ctx, done := testx.Context(t)
	defer done()

	path := filepath.Join(t.TempDir(), "cache.db")

	// build a cache.db at the schema prior to sync_id and populate it.
	db, err := InitializeDatabase(ctx, path, cacheMigrationsBefore(t, "20261113120000_library_known_media_add_sync_id.sql"))
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `INSERT INTO library_known_media (uid, md5, md5_lower, title) VALUES
		('11111111-1111-4111-8111-111111111111', '21111111-1111-4111-8111-111111111111', 1, 'one'),
		('11111111-1111-4111-8111-111111111112', '21111111-1111-4111-8111-111111111112', 2, 'two'),
		('11111111-1111-4111-8111-111111111113', '21111111-1111-4111-8111-111111111113', 3, 'three')`)
	require.NoError(t, err)
	require.NoError(t, db.Close())

	// reopen with the full migration set, applying the sync_id migration to
	// the existing rows.
	db, err = migrateCache(ctx, path)
	require.NoError(t, err)
	defer db.Close()

	var count, nils int
	require.NoError(t, db.QueryRowContext(ctx, `SELECT COUNT(*), COUNT(*) FILTER (WHERE sync_id = '00000000-0000-0000-0000-000000000000') FROM library_known_media`).Scan(&count, &nils))
	require.Equal(t, 3, count)
	require.Equal(t, 3, nils, "existing rows should be backfilled with the nil uuid")

	var title string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT title FROM library_known_media WHERE uid = '11111111-1111-4111-8111-111111111112'`).Scan(&title))
	require.Equal(t, "two", title)

	// the column must be NOT NULL after the migration.
	_, err = db.ExecContext(ctx, `UPDATE library_known_media SET sync_id = NULL WHERE uid = '11111111-1111-4111-8111-111111111111'`)
	require.Error(t, err)

	// rows inserted after the migration pick up the default.
	_, err = db.ExecContext(ctx, `INSERT INTO library_known_media (uid, md5, md5_lower, title) VALUES ('11111111-1111-4111-8111-111111111114', '21111111-1111-4111-8111-111111111114', 4, 'four')`)
	require.NoError(t, err)

	var sync string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT sync_id::VARCHAR FROM library_known_media WHERE uid = '11111111-1111-4111-8111-111111111114'`).Scan(&sync))
	require.Equal(t, "00000000-0000-0000-0000-000000000000", sync)
}
