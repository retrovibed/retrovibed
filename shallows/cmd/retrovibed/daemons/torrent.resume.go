package daemons

import (
	"context"
	"fmt"
	"log"

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

func ResumeDownloads(ctx context.Context, db sqlx.Queryer, rootstore fsx.Virtual, tclient *torrent.Client, tstore storage.ClientImpl) {
	q := tracking.MetadataSearchBuilder().Where(
		squirrel.And{
			tracking.MetadataQueryInitiated(),
			tracking.MetadataQueryIncomplete(),
			tracking.MetadataQueryNotPaused(),
		},
	)

	err := sqlxx.ScanEach(tracking.MetadataSearch(ctx, db, q), func(md *tracking.Metadata) error {
		infopath := rootstore.Path("torrent", fmt.Sprintf("%s.torrent", metainfo.Hash(md.Infohash).HexString()))

		metadata, err := torrent.New(metainfo.Hash(md.Infohash), torrent.OptionStorage(tstore), torrent.OptionTrackers([]string{md.Tracker}), torrentx.OptionInfoFromFile(infopath))
		if err != nil {
			log.Println(errorsx.Wrapf(err, "unable to create metadata from %s", md.ID))
			return nil
		}

		t, _, err := tclient.Start(metadata)
		if err != nil {
			return errorsx.Wrapf(err, "unable to start download %s", md.ID)
		}

		go func(infopath string, md *tracking.Metadata, dl torrent.Torrent) {
			errorsx.Log(errorsx.Wrap(tracking.Download(ctx, db, rootstore, md, dl), "resume failed"))
			torrentx.RecordInfo(infopath, dl.Metadata())
		}(infopath, md, t)

		log.Println("resumed", md.ID, md.Description)
		return nil
	})

	errorsx.Log(errorsx.Wrap(err, "failed to resume all downloads"))
}
