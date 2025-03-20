package cmdmedia

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/retrovibed/retrovibed/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/internal/env"
	"github.com/retrovibed/retrovibed/internal/x/errorsx"
	"github.com/retrovibed/retrovibed/internal/x/fsx"
	"github.com/retrovibed/retrovibed/internal/x/md5x"
	"github.com/retrovibed/retrovibed/library"
)

type importFilesystem struct {
	DryRun         bool     `flag:"" name:"dry-run" help:"print what files would be imported but do not actually perform the import" default:"false"`
	DeleteOriginal bool     `flag:"" name:"delete-original" short:"d" help:"after file is copied delete the original file from the disk" default:"false"`
	Paths          []string `arg:"" name:"paths" help:"files and folders to import" required:"true"`
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
		log.Println("DERP DERP", env.MediaDir())
		vfs := fsx.DirVirtual(env.MediaDir())
		op = library.ImportCopyFile(vfs)
	}

	for tx, cause := range library.ImportFilesystem(gctx.Context, op, t.Paths...) {
		log.Println("checkpoint")
		if cause != nil {
			log.Println(cause)
			err = errorsx.Compact(err, cause)
			continue
		}

		if t.DryRun {
			log.Println("imported", tx.Path)
			continue
		}

		log.Println("checkpoint")
		lmd := library.NewMetadata(
			md5x.FormatString(tx.MD5),
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
	DryRun bool   `flag:"" name:"dry-run" help:"print what files would be imported but do not actually perform the import" default:"false"`
	Path   string `arg:"" name:"path" help:"directory to export into" required:"true"`
}

func (t exportFilesystem) Run(gctx *cmdopts.Global) (err error) {
	return errors.ErrUnsupported
}
