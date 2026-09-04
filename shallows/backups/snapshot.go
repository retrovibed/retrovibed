package backups

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// Snapshot writes a consistent, encrypted copy of the open database to dst. duckdb does the
// copy and the encryption itself; the statements run on one connection because ATTACH is
// not visible across the pool.
func Snapshot(ctx context.Context, db *sql.DB, dst string, key string) (err error) {
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

	if _, err = conn.ExecContext(ctx, "FORCE CHECKPOINT"); err != nil {
		return errorsx.Wrap(err, "unable to checkpoint")
	}

	if _, err = conn.ExecContext(ctx, fmt.Sprintf("ATTACH '%s' AS backup (ENCRYPTION_KEY '%s')", dst, key)); err != nil {
		return errorsx.Wrap(err, "unable to attach backup")
	}
	defer func() {
		if _, cause := conn.ExecContext(ctx, "DETACH backup"); cause != nil && err == nil {
			err = errorsx.Wrap(cause, "unable to detach backup")
		}
	}()

	if _, err = conn.ExecContext(ctx, fmt.Sprintf("COPY FROM DATABASE %s TO backup", catalog)); err != nil {
		return errorsx.Wrap(err, "unable to copy database")
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
	defer func() {
		if _, cause := conn.ExecContext(ctx, "DETACH backup"); cause != nil && err == nil {
			err = errorsx.Wrap(cause, "unable to detach backup")
		}
	}()

	if _, err = conn.ExecContext(ctx, fmt.Sprintf("COPY FROM DATABASE backup TO %s", catalog)); err != nil {
		return errorsx.Wrap(err, "unable to copy database")
	}

	if _, err = conn.ExecContext(ctx, "FORCE CHECKPOINT"); err != nil {
		return errorsx.Wrap(err, "unable to checkpoint")
	}

	return nil
}
