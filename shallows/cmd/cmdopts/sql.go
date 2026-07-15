package cmdopts

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	_ "github.com/duckdb/duckdb-go/v2"

	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/internal/debugx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/goosex"
)

//go:embed .migrations/*.sql
var embedsqlite embed.FS

func DatabaseMeta(ctx context.Context) (db *sql.DB, err error) {
	return DatabaseCustom(ctx, userx.DefaultConfigDir(userx.DefaultRelRoot(), "meta.db"))
}

func DatabaseCustom(ctx context.Context, path string) (db *sql.DB, err error) {
	log.Println("database path", path)

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}

	if db, err = sql.Open("duckdb", path); err != nil {
		return nil, errorsx.Wrap(err, "unable to open db")
	}
	defer func() {
		if err == nil {
			return
		}
		debugx.Println("closing database due to error during initialization", err)
		errorsx.Log(db.Close())
	}()

	return db, InitializeDatabase(ctx, db)
}

func InitializeDatabase(ctx context.Context, db *sql.DB) (err error) {
	return goosex.InitializeDatabase(ctx, db, errorsx.Must(fs.Sub(embedsqlite, ".migrations")))
}

func Checkpoint(ctx context.Context, db *sql.DB) (err error) {
	log.Println("------------------------------------------------ database checkpoint initiated ------------------------------------------------")
	defer log.Println("------------------------------------------------ database checkpoint completed ------------------------------------------------")

	if _, err := db.ExecContext(ctx, "FORCE CHECKPOINT;"); err != nil {
		return errorsx.Wrap(err, "failed to checkpoint database")
	}

	return nil
}
