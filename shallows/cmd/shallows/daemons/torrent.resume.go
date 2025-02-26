package daemons

import (
	"context"
	"io"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/sqlx"
	"github.com/james-lawrence/deeppool/internal/x/sqlxx"
	"github.com/james-lawrence/deeppool/tracking"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
)

func ResumeDownloads(ctx context.Context, db sqlx.Queryer, tclient *torrent.Client, tstore storage.ClientImpl) {
	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryInitiated(),
			tracking.MetadataQueryIncomplete(),
			tracking.MetadataQueryNotPaused(),
		},
	)

	err := sqlxx.ScanEach(tracking.MetadataSearch(ctx, db, q), func(md *tracking.Metadata) error {
		metadata, err := torrent.New(metainfo.Hash(md.Infohash), torrent.OptionStorage(tstore), torrent.OptionTrackers([][]string{{md.Tracker}}))
		if err != nil {
			return errorsx.Wrapf(err, "unable to create metadata from metadata %s", md.ID)
		}

		t, _, err := tclient.Start(metadata)
		if err != nil {
			return errorsx.Wrapf(err, "unable to start download %s", md.ID)
		}

		go func(md *tracking.Metadata) {
			var (
				downloaded int64
			)

			pctx, done := context.WithCancel(ctx)
			defer done()

			// update the progress.
			go tracking.DownloadProgress(pctx, db, *md, t)

			// just copying as we receive data to block until done.
			if downloaded, err = torrent.DownloadInto(ctx, io.Discard, t); err != nil {
				log.Println(errorsx.Wrap(err, "download failed"))
				return
			}

			log.Println("download completed", md.ID, md.Description, downloaded)
			if err := tracking.MetadataProgressByID(ctx, db, md.ID, 0, uint64(downloaded)).Scan(md); err != nil {
				log.Println("failed to update progress", err)
			}
		}(md)

		log.Println("resumed", md.ID, md.Description)
		return nil
	})

	errorsx.Log(errorsx.Wrap(err, "failed to resume all downloads"))
}
