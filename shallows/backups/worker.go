package backups

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/linxGnu/pqueue"
	"github.com/linxGnu/pqueue/entry"
	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/pqueuex"
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

// Request is a queued backup of the current database for a device.
type Request struct {
	Device string `json:"device"`
}

// Enqueue requests a backup unless one is already waiting, so a device that was offline
// for a day uploads one backup when it returns rather than a day of them.
func Enqueue(ctx context.Context, wq pqueue.Queue, device string) error {
	var (
		pending entry.Entry
	)

	if wq.Peek(&pending) {
		return nil
	}

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
		key      string
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
