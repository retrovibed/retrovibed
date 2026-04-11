package cmdmedia

import (
	"database/sql"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/retrovibed/retrovibed/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/asynccompute"
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/md5x"
	"github.com/retrovibed/retrovibed/library"
)

type importFilesystem struct {
	DryRun         bool   `flag:"" name:"dry-run" help:"print what files would be imported but do not actually perform the import" default:"false"`
	DeleteOriginal bool   `flag:"" name:"delete-original" short:"d" help:"after file is copied delete the original file from the disk" default:"false"`
	Concurrency    uint16 `flag:"" name:"concurrency" help:"number of files to transfer concurrently, defaults to the number of cpus" default:"${vars_cores}"`
	Path           string `arg:"" name:"path" help:"file or folder to import" required:"true"`
}

func (t importFilesystem) Run(gctx *cmdopts.Global) (err error) {
	var (
		db *sql.DB
	)

	if db, err = cmdmeta.Database(gctx.Context); err != nil {
		return err
	}
	defer db.Close()

	op := library.ImportFileDryRun
	if !t.DryRun {
		if err = os.MkdirAll(env.MediaDir(), 0700); err != nil {
			return err
		}

		vfs := fsx.DirVirtual(env.MediaDir())
		op = library.ImportCopyFile(vfs)
	}

	imp := library.NewImporter(op, os.DirFS(t.Path).(fs.StatFS), asynccompute.Workers[string](t.Concurrency))

	for tx, cause := range imp.Import(gctx.Context, ".") {
		if cause != nil {
			log.Println(cause)
			err = errorsx.Compact(err, cause)
			continue
		}

		if t.DryRun {
			log.Println("imported", tx.Path)
			continue
		}

		lmd := library.NewMetadata(
			md5x.FormatUUID(tx.MD5),
			library.MetadataOptionDescription(filepath.Base(tx.Path)),
			library.MetadataOptionBytes(tx.Bytes),
			library.MetadataOptionMimetype(tx.Mimetype.String()),
		)

		if err := library.MetadataInsertWithDefaults(gctx.Context, db, lmd).Scan(&lmd); err != nil {
			return errorsx.Wrap(err, "unable to record library metadata")
		}

		log.Println("new library content", spew.Sdump(lmd))

		if t.DeleteOriginal {
			log.Println("removing", tx.Path)
			if err := os.Remove(tx.Path); err != nil {
				return errorsx.Wrap(err, "unable to remove original file")
			}
		}
	}

	return nil
}

type exportFilesystem struct {
	DryRun bool   `flag:"" name:"dry-run" help:"print what files would be exported but do not actually perform the export" default:"false"`
	Path   string `arg:"" name:"path" help:"directory to export into" required:"true"`
}

func (t exportFilesystem) Run(gctx *cmdopts.Global) (err error) {
	return errors.ErrUnsupported
}
