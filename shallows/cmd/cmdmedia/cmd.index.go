package cmdmedia

import (
	"database/sql"
	"log"
	"os"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/env"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/squirrelx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type reindex struct {
	Database  string `flag:"" name:"database" help:"database to read" default:"${vars_user_configuration_directory}/meta.db"`
	Unindexed bool   `flag:"" name:"unindexed" help:"only run against records that havent been indexed" default:"false"`
	DryRun    bool   `flag:"" name:"dry-run" help:"dont actually update" negatable:"" default:"true"`
}

func (t reindex) Run(gctx *cmdopts.Global) (err error) {
	var (
		db      *sql.DB
		missing squirrel.Sqlizer = squirrelx.Noop{}
	)

	if db, err = cmdopts.DatabaseCustom(gctx.Context, t.Database); err != nil {
		return err
	}
	defer db.Close()

	mediastore := fsx.DirVirtual(env.MediaDir())

	if t.Unindexed {
		missing = library.MetadataQueryNotIndexed()
	}

	query := library.MetadataSearchBuilder().Where(
		squirrel.And{
			squirrel.Expr("1=1"),
			missing,
		},
	)

	log.Println("records", errorsx.Zero(sqlx.Count(gctx.Context, db, "SELECT COUNT(*) FROM library_metadata")))

	s := sqlx.Scan(library.MetadataSearch(gctx.Context, db, query))
	for md := range s.Iter() {
		if uuid.FromStringOrNil(md.TorrentID).IsNil() {
			log.Println("unexpected nil torrent id on library.Metadata", md.ID)
			continue
		}

		var tmd tracking.Metadata

		if err = tracking.MetadataFindByID(gctx.Context, db, md.TorrentID).Scan(&tmd); sqlx.ErrNoRows(err) != nil {
			continue
		} else if err != nil {
			return err
		}

		resolved, err := os.Readlink(mediastore.Path(md.ID))
		if err != nil {
			log.Println("failed to read link", err)
			continue
		}

		_, desc, auto := tracking.GenerateDescription(resolved, &tmd)
		log.Println("resetting description", md.ID, md.Description, "->", desc)
		log.Println("resetting autodescription", md.ID, md.AutoDescription, "->", auto)
		if t.DryRun {
			continue
		}

		if err = library.MetadataUpdateDescriptionByID(gctx.Context, db, md.ID, desc).Scan(&md); err != nil {
			return err
		}

		if err = library.MetadataUpdateAutodescriptionByID(gctx.Context, db, md.ID, library.NormalizedDescription(md.Description)).Scan(&md); err != nil {
			return err
		}
	}

	if t.DryRun {
		log.Println("dry run - no records changed")
	}

	return s.Err()
}
