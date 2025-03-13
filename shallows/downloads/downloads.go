package downloads

import (
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/langx"
	"github.com/james-lawrence/deeppool/internal/x/slicesx"
	"github.com/james-lawrence/deeppool/internal/x/sqlx"
	"github.com/james-lawrence/deeppool/internal/x/userx"
	"github.com/james-lawrence/deeppool/tracking"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/storage"
)

type downloader interface {
	Start(t torrent.Metadata) (dl torrent.Torrent, added bool, err error)
}

func NewDirectoryWatcher(ctx context.Context, q sqlx.Queryer, dl downloader, s storage.ClientImpl) (d Directory, err error) {
	var (
		w *fsnotify.Watcher
	)

	if w, err = fsnotify.NewWatcher(); err != nil {
		return d, err
	}

	return Directory{
		d: dl,
		w: w,
		c: userx.DefaultCacheDirectory(userx.DefaultRelRoot()),
		s: s,
		q: q,
	}.background(ctx), nil
}

type Directory struct {
	d downloader
	q sqlx.Queryer
	w *fsnotify.Watcher
	c string
	s storage.ClientImpl
}

func (t Directory) Add(path string) (err error) {
	defer func() {
		if err == nil {
			log.Println("watching", path)
		}
	}()

	if err = errorsx.Wrapf(t.w.Add(path), "unable to watch: %s", path); err != nil {
		return err
	}

	err = fs.WalkDir(os.DirFS(path), ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".torrent") {
			return nil
		}

		go t.download(context.Background(), filepath.Join(path, name))

		return nil
	})

	return errorsx.Wrap(err, "unable to find existing torrents")
}

// background download
func (t Directory) download(ctx context.Context, path string) {
	meta, err := torrent.NewFromMetaInfoFile(path, torrent.OptionStorage(t.s))
	if err != nil {
		log.Println("unable to process", path, "ignoring", err)
		return
	}

	var (
		md         tracking.Metadata
		downloaded int64
	)

	tor, _, err := t.d.Start(meta)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to start torrent"))
		return
	}

	log.Println("wait for torrent info", meta.InfoHash)
	select {
	case <-tor.GotInfo():
	case <-ctx.Done():
		log.Println("failed to retrieve torrent information, manually restart will be required")
		return
	}

	if err = tracking.MetadataInsertWithDefaults(
		ctx,
		t.q,
		tracking.NewMetadata(langx.Autoptr(tor.Metadata().InfoHash),
			tracking.MetadataOptionFromInfo(tor.Info()),
			tracking.MetadataOptionTrackers(slicesx.Flatten(meta.Trackers...)...),
		),
	).Scan(&md); err != nil {
		log.Println(errorsx.Wrap(err, "unable to insert metadata"))
		return
	}

	if err = tracking.MetadataDownloadByID(ctx, t.q, md.ID).Scan(&md); err != nil {
		log.Println(errorsx.Wrap(err, "unable to mark metadata as downloading"))
		return
	}

	pctx, done := context.WithCancel(ctx)
	defer done()

	// update the progress.
	go tracking.DownloadProgress(pctx, t.q, &md, tor)

	// just copying as we receive data to block until done.
	if downloaded, err = torrent.DownloadInto(ctx, io.Discard, tor); err != nil {
		log.Println(errorsx.Wrap(err, "download failed"))
		return
	}

	log.Println("download completed", md.ID, md.Description, downloaded)

	if err := tracking.MetadataProgressByID(ctx, t.q, md.ID, 0, uint64(downloaded)).Scan(&md); err != nil {
		log.Println("failed to update progress", err)
	}
}

func (t Directory) background(ctx context.Context) Directory {
	go func() {
		for {
			select {
			case evt := <-t.w.Events:
				switch evt.Op {
				case fsnotify.Create:
				case fsnotify.Chmod:
					// fallthrough.
				default:
					log.Println("change ignored", evt.Op)
					continue
				}

				go t.download(ctx, evt.Name)
			case err := <-t.w.Errors:
				log.Println("watch error", err)
			case <-ctx.Done():
				log.Println("context completed", ctx.Err())
				return
			}
		}
	}()

	return t
}
