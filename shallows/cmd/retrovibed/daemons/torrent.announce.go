package daemons

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
	"github.com/retrovibed/retrovibed/internal/sqlx"
	"github.com/retrovibed/retrovibed/internal/sqlxx"
	"github.com/retrovibed/retrovibed/internal/torrentx"
	"github.com/retrovibed/retrovibed/tracking"
)

func AccounceTorrent(ctx context.Context, db sqlx.Queryer, rootstore fsx.Virtual, tclient *torrent.Client, tstore storage.ClientImpl) {
	const defaultInterval = time.Hour
	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQuerySeeding(),
		},
	)

	err := sqlxx.ScanEach(tracking.MetadataSearch(ctx, db, q), func(md *tracking.Metadata) error {
		infopath := rootstore.Path("torrent", fmt.Sprintf("%s.torrent", metainfo.Hash(md.Infohash).HexString()))

		go func(md tracking.Metadata) {
			metadata, err := torrent.New(metainfo.Hash(md.Infohash), torrent.OptionStorage(tstore), torrent.OptionTrackers([]string{md.Tracker}), torrentx.OptionInfoFromFile(infopath))
			if err != nil {
				log.Println(errorsx.Wrapf(err, "unable to create metadata from %s - %s", md.ID, infopath))
				return
			}

			t, _, err := tclient.Start(metadata)
			if err != nil {
				log.Println(errorsx.Wrapf(err, "unable to seed torrent %s - %s", md.ID, infopath))
				return
			}

			interval := time.Minute
			for {
				select {
				case <-time.After(interval):
				case <-ctx.Done():
					return
				}

				resp, err := torrent.TrackerEvent(ctx, t)
				if err != nil {
					log.Println("tracker even failed, sleeping for an hour", md.ID, err)
					interval = defaultInterval
					continue
				}

				log.Println("announced", md.ID, resp.Leechers, resp.Seeders, "next", resp.Interval)
			}
		}(*md)

		log.Println("resumed", md.ID, md.Description)
		return nil
	})

	errorsx.Log(errorsx.Wrap(err, "failed to announce seeds"))
}
