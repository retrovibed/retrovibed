package downloads

import (
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
	"github.com/james-lawrence/deeppool/internal/x/fsx"
	"github.com/james-lawrence/deeppool/internal/x/sqlx"
	"github.com/james-lawrence/deeppool/internal/x/timex"
	"github.com/james-lawrence/deeppool/internal/x/userx"
	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/storage"
)

type downloader interface {
	Start(t torrent.Metadata) (dl torrent.Torrent, added bool, err error)
}

func NewDirectoryWatcher(ctx context.Context, q sqlx.Queryer, dl downloader) (d Directory, err error) {
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
		s: userx.DefaultDataDirectory(userx.DefaultRelRoot(), "media"),
		q: q,
	}.background(ctx), nil
}

type Directory struct {
	d downloader
	q sqlx.Queryer
	w *fsnotify.Watcher
	c string
	s string
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
	meta, err := torrent.NewFromMetaInfoFile(path, torrent.OptionStorage(storage.NewFile(t.s)))
	if err != nil {
		log.Println("unable to process", path, "ignoring", err)
		return
	}

	var (
		dst *os.File
		src io.ReadCloser
	)

	if err = os.Mkdir(t.c, 0700); fsx.IgnoreIsExist(err) != nil {
		log.Println(errorsx.Wrap(err, "unable to ensure temp directory"))
		return
	}

	if dst, err = os.CreateTemp(t.c, meta.InfoHash.HexString()); err != nil {
		log.Println(errorsx.Wrap(err, "unable to open download destination"))
		return
	}
	defer dst.Close()

	tor, _, err := t.d.Start(meta)
	if err != nil {
		log.Println(errorsx.Wrap(err, "unable to start torrent"))
		return
	}

	go timex.Every(10*time.Second, func() {
		stats := tor.Stats()
		log.Printf("%s: peers(%d:%d:%d) pieces(%d:%d:%d:%d)\n", tor.Metainfo().HashInfoBytes().HexString(), stats.ActivePeers, stats.PendingPeers, stats.TotalPeers, stats.Missing, stats.Outstanding, stats.Unverified, stats.Completed)
	})

	if err = torrent.DownloadInto(ctx, dst, tor); err != nil {
		log.Println(errorsx.Wrap(err, "download failed"))
		return
	}

	if _, err := io.Copy(dst, src); err != nil {
		log.Println("download failed", err)
		return
	}

	if err := os.MkdirAll(t.s, 0700); err != nil {
		log.Println("unable to ensure storage directory", err)
		return
	}

	if err := os.Rename(dst.Name(), filepath.Join(t.s, meta.InfoHash.HexString())); err != nil {
		log.Println("unable rename", dst.Name(), "->", filepath.Join(t.s, meta.InfoHash.HexString()), err)
		return
	}
}

func (t Directory) background(ctx context.Context) Directory {
	go func() {
		for {
			select {
			case evt := <-t.w.Events:
				switch evt.Op {
				case fsnotify.Create:
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
