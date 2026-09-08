package cmdmeta

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/deeppool"
	"github.com/retrovibed/retrovibed/shallows/backups"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sshx"
)

type CloudBackupRestore struct {
	Database string `flag:"" name:"database" help:"database to restore into" default:"${vars_user_configuration_directory}/meta.db"`
	Force    bool   `flag:"" name:"force" help:"replace an existing database, keeping a timestamped copy beside it" default:"false"`
}

// fetches the latest encrypted backup and decrypts it into a fresh database file. the
// daemon must not be running against the target: duckdb allows one writer, and a database
// swapped out from under a live process is a corrupted database.
func (t CloudBackupRestore) Run(gctx *cmdopts.Global) (err error) {
	if fsx.Exists(t.Database) && !t.Force {
		return errorsx.String(fmt.Sprintf("%s already exists, pass --force to replace it", t.Database))
	}

	id, err := sshx.Load(env.PrivateKeyPath())
	if err != nil {
		return err
	}

	c, err := authn.AutoJWTClient(gctx.Context, id)
	if err != nil {
		return errorsx.Wrap(err, "unable to authenticate")
	}

	key, err := backups.ResolveKey(gctx.Context, c)
	if err != nil {
		return err
	}

	client := deeppool.NewBackups(c)
	latest, err := client.Latest(gctx.Context)
	if err != nil {
		return errorsx.Wrap(err, "unable to locate latest backup")
	}

	dir, err := os.MkdirTemp(filepath.Dir(t.Database), "retrovibed.restore.*")
	if err != nil {
		return errorsx.Wrap(err, "unable to create restore directory")
	}
	defer os.RemoveAll(dir)

	encrypted := filepath.Join(dir, "backup.db")
	f, err := os.Create(encrypted)
	if err != nil {
		return errorsx.Wrap(err, "unable to create download")
	}
	defer func() {
		err = errors.Join(err, errorsx.Wrap(f.Close(), "unable to close backup file"))
	}()

	if err = client.Download(gctx.Context, latest.Id, f); err != nil {
		return errorsx.Wrap(err, "unable to download backup")
	}

	// the restored database is built beside the target and moved into place only once it
	// opens, so an interrupted restore or a wrong key leaves whatever was there untouched.
	restored := filepath.Join(dir, "meta.db")
	if err = backups.Restore(key, encrypted, restored); err != nil {
		return err
	}

	db, err := sql.Open("duckdb", restored)
	if err != nil {
		return errorsx.Wrap(err, "unable to open restored database")
	}

	if err = db.PingContext(gctx.Context); err != nil {
		errorsx.Log(db.Close())
		return errorsx.Wrap(err, "unable to open restored database")
	}

	if err = db.Close(); err != nil {
		return errorsx.Wrap(err, "unable to close restored database")
	}

	if fsx.Exists(t.Database) {
		kept := fmt.Sprintf("%s.%d", t.Database, time.Now().Unix())
		if err = os.Rename(t.Database, kept); err != nil {
			return errorsx.Wrap(err, "unable to keep existing database")
		}
		log.Println("existing database kept at", kept)
	}

	if err = os.MkdirAll(filepath.Dir(t.Database), 0700); err != nil {
		return errorsx.Wrap(err, "unable to create database directory")
	}

	if err = os.Rename(restored, t.Database); err != nil {
		return errorsx.Wrap(err, "unable to move restored database into place")
	}

	log.Println("restored", latest.Id, "into", t.Database)
	return nil
}
