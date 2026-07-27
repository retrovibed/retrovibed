package daemons

import (
	"context"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/lucenex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

func IdentifyTorrentMedia(ctx context.Context, db sqlx.Queryer, mc library.QueryCleaner) error {
	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryNeedsKnownMediaID(),
			tracking.MetadataQueryNotMediaArchive(),
			tracking.MetadataQueryNotNeural(),
		},
	)

	iter := sqlx.Scan(tracking.MetadataSearch(ctx, db, q))

	log.Println("attempting to locate unidentified media initiated")
	defer log.Println("attempting to locate unidentified media completed")

	ts := time.Now()

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
			log.Println("unable to clean media for torrent - detected messy description ended up with blank", md.ID, md.Description, "|", "''", "|", err)
			if err = tracking.MetadataAssignKnownMediaID(ctx, db, md.ID, uuid.Nil.String()).Scan(&md); err != nil {
				log.Println("unable to mark torrent as unidentifiable", md.ID, err)
			}
			continue
		}

		cleaned = lucenex.Clean(cleaned)
		title, _, _ := library.ParseReleaseEpisode(cleaned)

		if known, err = library.DetectKnownMedia(ctx, db, mimex.Category(md.Mimetype), title, library.KnownMatchCutoff); err != nil {
			log.Println("unable to detect media for torrent", md.ID, md.Description, "|", md.Description, "|", err)
			continue
		}

		if uuid.FromStringOrNil(known.UID).IsNil() {
			log.Println("unable to detect media for torrent", md.ID, md.Description, "|", md.Description)
			if err = tracking.MetadataAssignKnownMediaID(ctx, db, md.ID, uuid.Nil.String()).Scan(&md); err != nil {
				log.Println("unable to mark torrent as unidentifiable", md.ID, err)
			}
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
