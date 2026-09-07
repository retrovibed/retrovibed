package backups

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"

	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// Snapshot writes a consistent copy of the open database to dst. duckdb does the copy; the
// statements run on one connection because ATTACH is not visible across the pool.
func Snapshot(ctx context.Context, db *sql.DB, dst string) (err error) {
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

	if _, err = conn.ExecContext(ctx, fmt.Sprintf("ATTACH '%s' AS backup", dst)); err != nil {
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

// Restore decrypts a downloaded backup into dst, a fresh file the caller moves into place
// once it has been opened successfully.
func Restore(key Key, src string, dst string) (err error) {
	f, err := os.Open(src)
	if err != nil {
		return errorsx.Wrap(err, "unable to open backup")
	}
	defer func() { errorsx.Log(f.Close()) }()

	plain, err := key.Decrypt(f)
	if err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return errorsx.Wrap(err, "unable to create restored database")
	}

	if _, err = io.Copy(out, plain); err != nil {
		errorsx.Log(out.Close())
		return errorsx.Wrap(err, "unable to restore backup")
	}

	return errorsx.Wrap(out.Close(), "unable to finish restore")
}
