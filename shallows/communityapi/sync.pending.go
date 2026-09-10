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
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

type Publisher interface {
	Publish(ctx context.Context, req *PublishContentRequest, torrent io.Reader) (*PublishContentResponse, error)
	Delete(ctx context.Context, id string) (*PublishContentDeleteResponse, error)
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

func publishToPlugin(ctx context.Context, mvfs fsx.Virtual, publishers publishplugin.T, pub community.PluginPublisher, req publishplugin.Request) (err error) {
	if req.Directory, err = os.MkdirTemp(userx.DefaultRuntimeDirectory(userx.DefaultRelRoot()), "retrovibed.publish.*"); err != nil {
		return errorsx.Wrap(err, "unable to create temporary media file")
	}
	defer func() {
		errorsx.Log(errorsx.Wrap(fsx.IgnoreIsNotExist(os.RemoveAll(req.Directory)), "unable to publishing directory"))
	}()

	diskpath := mvfs.Path(req.MediaPath)
	if err := os.Symlink(diskpath, filepath.Join(req.Directory, req.MediaPath)); err != nil {
		fsx.PrintPath(diskpath)
		fsx.PrintPath(req.Directory)
		fsx.PrintPath(filepath.Join(req.Directory, req.MediaPath))
		return errorsx.Wrap(err, "failed to symlink media into publishing directory")
	}

	if _, err = publishers.Publish(ctx, pub.Path, req); err != nil {
		return errorsx.Wrapf(err, "failed to publish using %s", pub.Path)
	}

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
