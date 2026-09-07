package communityapi

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid/v5"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/linxGnu/pqueue"
	"github.com/retrovibed/retrovibed/retroapi/publishplugin"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/pqueuex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type MetricsPublisher interface {
	Publish(ctx context.Context, req *PublishContentRequest, torrent io.Reader) (*PublishContentResponse, error)
}

func magnetURI(tmd tracking.Metadata, name string) string {
	return metainfo.NewMagnetFromInfohash(tmd.Infohash, metainfo.MagnetOptionDisplayName(name)).String()
}

func ensureTorrent(ctx context.Context, q sqlx.Queryer, mvfs, tvfs fsx.Virtual, lmd *library.Metadata) (tmd tracking.Metadata, err error) {
	if lmd.TorrentID != uuid.Nil.String() {
		if err = tracking.MetadataFindByID(ctx, q, lmd.TorrentID).Scan(&tmd); err != nil {
			return tmd, err
		}

		if !tmd.Seeding || tmd.CompletedAt.Equal(timex.Inf()) {
			if err = media.ValidateTorrent(ctx, q, tvfs, &tmd); err != nil {
				return tmd, err
			}
		}

		return tmd, err
	}

	return media.GenerateTorrent(ctx, q, mvfs, tvfs, lmd)
}

// publishToPlugin materializes lmd's byte range into a flat temp file - a
// wasm guest has no way to interpret blockcache's internal block-file
// layout directly, so the exact section YouTubeUpload would stream is
// instead copied to disk once - and invokes the plugin installed at
// pub.Path via publishers.Publish, removing the temp file once the call
// returns.
func publishToPlugin(ctx context.Context, mvfs fsx.Virtual, publishers publishplugin.T, pub community.PluginPublisher, pc community.PublishedContent, lmd library.Metadata, known library.Known) error {
	dir, err := os.MkdirTemp(userx.DefaultRuntimeDirectory(userx.DefaultRelRoot()), "retrovibed.publish.*")
	if err != nil {
		return errorsx.Wrap(err, "unable to create temporary media file")
	}
	defer func() {
		errorsx.Log(errorsx.Wrap(fsx.IgnoreIsNotExist(os.RemoveAll(dir)), "unable to publishing directory"))
	}()

	if err := os.Symlink(mvfs.Path(lmd.ID), filepath.Join(dir, lmd.ID)); err != nil {
		fsx.PrintPath(mvfs.Path(lmd.ID))
		fsx.PrintPath(dir)
		fsx.PrintPath(filepath.Join(dir, lmd.ID))
		return errorsx.Wrap(err, "failed to symlink media into publishing directory")
	}

	_, err = publishers.Publish(ctx, pub.Path, publishplugin.Request{
		Directory:   dir,
		MediaPath:   lmd.ID,
		Title:       stringsx.FirstNonBlank(known.Title, lmd.Description),
		Description: known.Overview,
		Mimetype:    stringsx.FirstNonBlank(known.Mimetype, lmd.Mimetype),
		CommunityID: pc.CommunityID,
		Magnet:      pc.MagnetURI,
		Adult:       known.Adult,
	})

	return err
}

// SyncPendingToDeeppool syncs pending published content to deeppool and regenerates affected feeds.
func SyncPendingToDeeppool(ctx context.Context, q sqlx.Queryer, wq pqueue.Queue) error {
	pending := sqlx.Scan(community.PublishedContentFindByPendingSync(ctx, q))

	log.Println("sync pending to deeppool initiated")
	defer log.Println("sync pending to deeppool completed")

	for pc := range pending.Iter() {
		if err := pqueuex.Enqueue(ctx, wq, langx.Clone(pc, timex.JSONSafeEncodeOption)); err != nil {
			log.Println(errorsx.Wrap(err, "unable to queue published content"))
			continue
		}
	}

	return pending.Err()
}
