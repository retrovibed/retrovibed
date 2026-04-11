package daemons

import (
	"context"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/library"
	"github.com/retrovibed/retrovibed/tracking"
)

func IdentifyTorrentMedia(ctx context.Context, db sqlx.Queryer) error {
	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryNeedsKnownMediaID(),
			tracking.MetadataQueryNotMetadataArchive(),
		},
	)

	iter := sqlx.Scan(tracking.MetadataSearch(ctx, db, q))

	log.Println("attempting to locate unidentified media initiated")
	defer log.Println("attempting to locate unidentified media completed")

	ts := time.Now()

	for md := range iter.Iter() {
		var (
			err   error
			known library.Known
		)

		if known, err = library.DetectKnownMedia(ctx, db, md.Description); err != nil {
			log.Println("unable to detect media for torrent", md.ID, md.Description, "|", md.Description, "|", err)
			continue
		}

		if uuid.FromStringOrNil(known.UID).IsNil() {
			log.Println("unable to detect media for torrent", md.ID, md.Description, "|", md.Description)
			continue
		}

		if err = tracking.MetadataAssignKnownMediaID(ctx, db, md.ID, known.UID).Scan(&md); err != nil {
			log.Println("unable to assign known media id to torrent", md.ID, known.UID, err)
			continue
		}

		log.Println("matched", md.ID, "->", known.UID, md.Description, "->", known.Title)
	}

	if err := iter.Err(); err != nil {
		return errorsx.Wrap(iter.Err(), "failed to mark known media torrents")
	}

	if err := sqlx.Discard(sqlx.Scan(library.MetadataTransferKnownMediaIDFromTorrent(ctx, db, ts))); err != nil {
		return errorsx.Wrap(iter.Err(), "failed to associate known media with upstream library")
	}

	return nil
}
