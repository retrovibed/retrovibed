package cmdtorrent

import (
	"database/sql"
	"log"

	"github.com/retrovibed/retrovibed/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/library"
	"github.com/retrovibed/retrovibed/tracking"
)

type importFilesystem struct {
	Paths []string `arg:"" name:"paths" help:"files and folders to import" required:"true"`
}

func (t importFilesystem) Run(gctx *cmdopts.Global) (err error) {
	var (
		db *sql.DB
	)

	if db, err = cmdmeta.Database(gctx.Context); err != nil {
		return err
	}
	defer db.Close()

	log.Println("duckdb version", errorsx.Zero(sqlx.String(gctx.Context, db, "SELECT version() AS version")))

	mvfs := fsx.DirVirtual(env.MediaDir())
	tvfs := fsx.DirVirtual(env.TorrentDir())

	if err := fsx.MkDirs(0700, mvfs.Path(), tvfs.Path()); err != nil {
		return err
	}
	op := tracking.ImportTorrent(db, mvfs, tvfs)

	for tx, cause := range library.ImportFilesystem(gctx.Context, op, t.Paths...) {
		if cause != nil {
			log.Println(cause)
			err = errorsx.Compact(err, cause)
			continue
		}

		log.Println("imported", tx.Path)
	}

	log.Println("import complete, torrents will startup once daemon is started")
	return nil
}
