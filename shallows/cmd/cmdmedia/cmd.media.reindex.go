package cmdmedia

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/unsafepretty"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type reindex struct {
	Database       string `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
	Unindexed      bool   `flag:"" name:"unindexed" help:"only run against records that havent been indexed" default:"false"`
	Identification bool   `flag:"" name:"reidentify" help:"mark the records for reidentification"`
	DryRun         bool   `flag:"" name:"dry-run" help:"dont actually update" negatable:"" default:"true"`
}

func (t reindex) Run(gctx *cmdopts.Global) (err error) {
	var db *sql.DB

	if db, err = cmdopts.DatabaseMetaCustom(gctx.Context, t.Database); err != nil {
		return err
	}
	defer db.Close()

	cleaner := library.NewQueryerCleanerAuto()

	return t.run(gctx.Context, db, cleaner, fsx.DirVirtual(env.MediaDir()))
}

func (t reindex) run(ctx context.Context, db *sql.DB, c library.QueryCleaner, mediastore fsx.Virtual) (err error) {
	var (
		missing squirrel.Sqlizer = squirrelx.Noop{}
	)

	if t.Unindexed {
		missing = library.MetadataQueryNotIndexed()
	}

	query := library.MetadataSearchBuilder().Where(
		squirrel.And{
			squirrel.Expr("1=1"),
			library.MetadataQueryIsDirectory(false),
			missing,
		},
	)

	log.Println("records", errorsx.Zero(sqlx.Count(ctx, db, "SELECT COUNT(*) FROM library_metadata")))

	kid := ""
	if t.Identification {
		kid = uuid.Max.String()
	}

	s := sqlx.Scan(library.MetadataSearch(ctx, db, query))
	for md := range s.Iter() {
		if uuid.FromStringOrNil(md.TorrentID).IsNil() {
			log.Println("unexpected nil torrent id on library.Metadata", md.ID)
			continue
		}

		var tmd tracking.Metadata

		if err = tracking.MetadataFindByID(ctx, db, md.TorrentID).Scan(&tmd); sqlx.ErrNoRows(err) != nil {
			continue
		} else if err != nil {
			return err
		}

		resolved, err := os.Readlink(mediastore.Path(md.ID))
		if err != nil {
			log.Println("failed to read link", err)
			continue
		}

		finfo, err := tracking.FileInfoFromOffset(resolved, md.DiskOffset)
		if err != nil {
			log.Println("failed to read file info", err)
			continue
		}

		o, desc, auto := tracking.GenerateDescription(finfo.Path, &tmd)

		log.Println("---------------------------------------------")
		log.Println("unmodified", o)
		log.Println("resetting description", md.ID, md.Description, "->", desc)
		log.Println("resetting autodescription", md.ID, md.AutoDescription, "->", auto)
		log.Println("neural result", o, "->", unsafepretty.Print(errorsx.Zero(c.Clean(ctx, o)), unsafepretty.OptionNewlineRunes()))
		log.Println("---------------------------------------------")

		if t.DryRun {
			continue
		}

		if err = library.MetadataUpdateReindexByID(ctx, sqlx.Debug(db), md.ID, desc, library.NormalizedDescription(md.Description), stringsx.FirstNonBlank(kid, md.KnownMediaID)).Scan(&md); err != nil {
			return err
		}
	}

	if t.DryRun {
		log.Println("dry run - no records changed")
	}

	return s.Err()
}
