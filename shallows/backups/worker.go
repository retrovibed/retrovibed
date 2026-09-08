package backups

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/linxGnu/pqueue"
	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/pqueuex"
)

// Run snapshots the database and uploads it encrypted for the device. the plain copy lives
// in a private temporary directory only as long as the upload, whether or not it succeeds.
func Run(ctx context.Context, c *http.Client, db *sql.DB, device string, key Key) (m *deeppool.Media, err error) {
	dir, err := os.MkdirTemp("", "retrovibed.backup.*")
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to create backup directory")
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "meta.db")
	if err = Snapshot(ctx, db, path); err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, errorsx.Wrap(err, "unable to open snapshot")
	}
	defer f.Close()

	encrypted, err := key.Encrypt(f)
	if err != nil {
		return nil, err
	}

	return deeppool.NewBackups(c).Upload(ctx, device, mimex.RetrovibedMetaBackup, encrypted)
}

// Request is a queued backup of the current database for a device.
type Request struct {
	Device string `json:"device"`
}

// Enqueue requests a backup of the current database for the device.
func Enqueue(ctx context.Context, wq pqueue.Queue, device string) error {
	return pqueuex.Enqueue(ctx, wq, Request{Device: device})
}

func NewWorker(c *http.Client, db *sql.DB) Worker {
	return Worker{c: c, db: db}
}

// Worker processes queued backup requests. the key is resolved per request so a rotated
// seed takes effect without a restart.
type Worker struct {
	c  *http.Client
	db *sql.DB
}

func (t Worker) Message(ctx context.Context, m []byte) (err error) {
	var (
		req      Request
		key      Key
		uploaded *deeppool.Media
	)

	if err = jsonx.Unmarshal(m, &req); err != nil {
		return err
	}

	if key, err = ResolveKey(ctx, t.c); err != nil {
		return err
	}

	if uploaded, err = Run(ctx, t.c, t.db, req.Device, key); err != nil {
		return errorsx.Wrap(err, "backup upload failed")
	}

	log.Println("backup uploaded", uploaded.Id, uploaded.Bytes)
	return nil
}
