package daemons

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

func IdentifyLibraryMedia(ctx context.Context, db sqlx.Queryer, mc library.QueryCleaner) error {
	q := library.MetadataSearchBuilder().Where(
		squirrel.And{
			library.MetadataQueryNeedsKnownMediaID(),
			library.MetadataQueryNotTombstoned(),
		},
	)

	iter := sqlx.Scan(library.MetadataSearch(ctx, db, q))

	log.Println("attempting to identify library media initiated")
	defer log.Println("attempting to identify library media completed")

	for md := range iter.Iter() {
		var (
			err     error
			cleaned string
			known   library.Known
		)

		if cleaned, err = mc.Clean(ctx, md.Description); err != nil {
			log.Println("unable to clean media for torrent", md.ID, md.Description, "|", cleaned, "|", err)
			continue
		} else if stringsx.Blank(cleaned) {
			log.Println("unable to clean media for library - detected messy description ended up with blank", md.ID, md.Description, "|", "''", "|", err)
			continue
		}

		if known, err = library.DetectKnownMedia(ctx, db, mimex.Category(md.Mimetype), cleaned); err != nil {
			log.Println("unable to detect known media for library media", md.ID, md.Description, "|", err)
			continue
		}

		if uuid.FromStringOrNil(known.UID).IsNil() {
			log.Println("no match for library media", md.ID, md.Description)
			continue
		}

		updated := md
		updated.KnownMediaID = known.UID
		if err = library.MetadataUpdate(ctx, db, md.ID, updated).Scan(&md); err != nil {
			log.Println("unable to assign known media id to library media", md.ID, known.UID, err)
			continue
		}

		log.Println("matched", md.ID, "->", known.UID, md.Description, "->", known.Title)
	}

	return errorsx.Wrap(iter.Err(), "failed to identify library media")
}
