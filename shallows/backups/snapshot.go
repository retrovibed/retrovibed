package backups

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// Snapshot writes a consistent, encrypted copy of the open database to dst. duckdb does the
// copy and the encryption itself; the statements run on one connection because ATTACH is
// not visible across the pool.
//
// the copy goes through a plain staging file beside dst because hnsw indexes cannot be
// written to an encrypted database: vss sizes its allocator to the plain block size and
// encrypted blocks are smaller. they are dropped from the staging copy and rebuilt by the
// package that owns them on the next start.
func Snapshot(ctx context.Context, db *sql.DB, dst string, key string) (err error) {
	var (
		catalog string
		staging = dst + ".staging"
	)

	conn, err := db.Conn(ctx)
	if err != nil {
		return errorsx.Wrap(err, "unable to acquire connection")
	}
	defer func() { errorsx.Log(conn.Close()) }()
	defer func() { errorsx.Log(errorsx.Ignore(os.Remove(staging), os.ErrNotExist)) }()

	if err = conn.QueryRowContext(ctx, "SELECT current_database()").Scan(&catalog); err != nil {
		return errorsx.Wrap(err, "unable to resolve catalog")
	}

	if _, err = conn.ExecContext(ctx, "FORCE CHECKPOINT"); err != nil {
		return errorsx.Wrap(err, "unable to checkpoint")
	}

	if _, err = conn.ExecContext(ctx, fmt.Sprintf("ATTACH '%s' AS staging", staging)); err != nil {
		return errorsx.Wrap(err, "unable to attach staging")
	}
	defer detach(ctx, conn, "staging", &err)

	if _, err = conn.ExecContext(ctx, fmt.Sprintf("COPY FROM DATABASE %s TO staging", catalog)); err != nil {
		return errorsx.Wrap(err, "unable to copy database")
	}

	if err = dropHNSW(ctx, conn, "staging"); err != nil {
		return err
	}

	if _, err = conn.ExecContext(ctx, fmt.Sprintf("ATTACH '%s' AS backup (ENCRYPTION_KEY '%s')", dst, key)); err != nil {
		return errorsx.Wrap(err, "unable to attach backup")
	}
	defer detach(ctx, conn, "backup", &err)

	if _, err = conn.ExecContext(ctx, "COPY FROM DATABASE staging TO backup"); err != nil {
		return errorsx.Wrap(err, "unable to encrypt database")
	}

	return nil
}

// Restore copies an encrypted backup into the database open at db, which is expected to be
// empty: a fresh file the caller will move into place once this returns.
func Restore(ctx context.Context, db *sql.DB, src string, key string) (err error) {
	var (
		catalog string
	)

	conn, err := db.Conn(ctx)
	if err != nil {
		return errorsx.Wrap(err, "unable to acquire connection")
	}
	defer func() { errorsx.Log(conn.Close()) }()

	if err = conn.QueryRowContext(ctx, "SELECT current_database()").Scan(&catalog); err != nil {
		return errorsx.Wrap(err, "unable to resolve catalog")
	}

	if _, err = conn.ExecContext(ctx, fmt.Sprintf("ATTACH '%s' AS backup (ENCRYPTION_KEY '%s', READ_ONLY)", src, key)); err != nil {
		return errorsx.Wrap(err, "unable to attach backup")
	}
	defer detach(ctx, conn, "backup", &err)

	if _, err = conn.ExecContext(ctx, fmt.Sprintf("COPY FROM DATABASE backup TO %s", catalog)); err != nil {
		return errorsx.Wrap(err, "unable to copy database")
	}

	if _, err = conn.ExecContext(ctx, "FORCE CHECKPOINT"); err != nil {
		return errorsx.Wrap(err, "unable to checkpoint")
	}

	return nil
}

// detaching on the way out is what lets the next call attach under the same name.
func detach(ctx context.Context, conn *sql.Conn, name string, err *error) {
	if _, cause := conn.ExecContext(ctx, "DETACH "+name); cause != nil && *err == nil {
		*err = errorsx.Wrapf(cause, "unable to detach %s", name)
	}
}

// duckdb_indexes() has no type column, so hnsw indexes are found by their definition.
func dropHNSW(ctx context.Context, conn *sql.Conn, catalog string) (err error) {
	var (
		name  string
		names []string
	)

	rows, err := conn.QueryContext(ctx, "SELECT index_name FROM duckdb_indexes() WHERE database_name = $1 AND sql ILIKE '%USING HNSW%'", catalog)
	if err != nil {
		return errorsx.Wrap(err, "unable to list indexes")
	}

	for rows.Next() {
		if err = rows.Scan(&name); err != nil {
			break
		}

		names = append(names, name)
	}

	if err = errorsx.Compact(err, rows.Err(), rows.Close()); err != nil {
		return errorsx.Wrap(err, "unable to list indexes")
	}

	for _, name = range names {
		if _, err = conn.ExecContext(ctx, fmt.Sprintf(`DROP INDEX %s."%s"`, catalog, name)); err != nil {
			return errorsx.Wrapf(err, "unable to drop index %s", name)
		}
	}

	return nil
}
