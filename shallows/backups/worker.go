package backups

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// Run takes one encrypted snapshot and uploads it for the device. the snapshot lands in a
// temporary file that is ciphertext from the first byte, and is removed whether or not the
// upload succeeds.
func Run(ctx context.Context, c *http.Client, db *sql.DB, device string, key string) (m *deeppool.Media, err error) {
	dir, err := os.MkdirTemp("", "retrovibed.backup.*")
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create backup directory")
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "meta.db")
	if err = Snapshot(ctx, db, path, key); err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to open snapshot")
	}
	defer f.Close()

	return deeppool.NewBackups(c).Upload(ctx, device, mimex.RetrovibedMetaBackup, f)
}

// NewAutoBackup uploads an encrypted snapshot of the database on every wakeup, shaped like
// library.NewAutoArchive. the key is resolved per pass so a rotated seed takes effect
// without a restart.
func NewAutoBackup(ctx context.Context, c *http.Client, db *sql.DB, async *asyncx.Wakeup, device string, enabled bool) error {
	log.Println("auto backup initiated")
	defer log.Println("auto backup completed")

	backup := func(ctx context.Context) error {
		log.Println("backup initiated")
		defer log.Println("backup completed")

		if !enabled {
			log.Println("dry-run - not backing up")
			return nil
		}

		key, err := ResolveKey(ctx, c)
		if err != nil {
			return err
		}

		m, err := Run(ctx, c, db, device, key)
		if err != nil {
			return errorsx.Wrap(err, "backup upload failed")
		}

		log.Println("backup uploaded", m.Id, m.Bytes)
		return nil
	}

	if err := asyncx.Run(ctx, async, backup); errorsx.Is(err, context.Canceled, context.DeadlineExceeded) {
		return nil
	} else if errorsx.Is(err, asyncx.ErrWakeupClosed) {
		return errorsx.Ignore(backup(ctx), context.Canceled, context.DeadlineExceeded)
	} else {
		return errorsx.Wrap(err, "backup failed")
	}
}
