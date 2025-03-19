package daemons

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/internal/x/errorsx"
	"github.com/retrovibed/retrovibed/internal/x/fsx"
	"github.com/retrovibed/retrovibed/internal/x/sqlx"
	"github.com/retrovibed/retrovibed/internal/x/sqlxx"
	"github.com/retrovibed/retrovibed/tracking"
)

func ResumeDownloads(ctx context.Context, db sqlx.Queryer, rootstore fsx.Virtual, tclient *torrent.Client, tstore storage.ClientImpl) {
	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryInitiated(),
			tracking.MetadataQueryIncomplete(),
			tracking.MetadataQueryNotPaused(),
		},
	)

	err := sqlxx.ScanEach(tracking.MetadataSearch(ctx, db, q), func(md *tracking.Metadata) error {
		metadata, err := torrent.New(metainfo.Hash(md.Infohash), torrent.OptionStorage(tstore), torrent.OptionTrackers([]string{md.Tracker}))
		if err != nil {
			return errorsx.Wrapf(err, "unable to create metadata from metadata %s", md.ID)
		}

		t, _, err := tclient.Start(metadata)
		if err != nil {
			return errorsx.Wrapf(err, "unable to start download %s", md.ID)
		}

		go func(md *tracking.Metadata, dl torrent.Torrent) {
			errorsx.Log(errorsx.Wrap(tracking.Download(ctx, db, rootstore, md, dl), "resume failed"))
		}(md, t)

		log.Println("resumed", md.ID, md.Description)
		return nil
	})

	errorsx.Log(errorsx.Wrap(err, "failed to resume all downloads"))
}
