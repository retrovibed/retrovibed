package cmdtorrent

import (
	"database/sql"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type cmdKnownMedia struct {
	Timestamp time.Time `flag:"" name:"timestamp" default:"${vars_timestamp_started}"`
}

func (t cmdKnownMedia) Run(ctx *cmdopts.Global) (err error) {
	var (
		db *sql.DB
	)

	if db, err = cmdopts.DatabaseMeta(ctx.Context); err != nil {
		return err
	}
	defer db.Close()

	q := tracking.MetadataSearchBuilder().Where(squirrel.And{
		tracking.MetadataQueryNeedsKnownMediaID(),
	})

	iter := sqlx.Scan(tracking.MetadataSearch(ctx.Context, db, q))
	for md := range iter.Iter() {
		var (
			known library.Known
		)

		if known, err = library.DetectKnownMedia(ctx.Context, db, mimex.Category(md.Mimetype), md.Description, library.KnownMatchCutoff); err != nil {
			log.Println("failed to detect known media", err)
			continue
		}

		if err = tracking.MetadataAssignKnownMediaID(ctx.Context, db, md.ID, known.UID).Scan(&md); err != nil {
			log.Println("failed to assign known media", md.ID, known.UID)
		}
	}

	if err := sqlx.Discard(sqlx.Scan(library.MetadataTransferKnownMediaIDFromTorrent(ctx.Context, db, t.Timestamp))); err != nil {
		return errorsx.Wrap(iter.Err(), "failed to associate known media with upstream library")
	}

	return err
}
