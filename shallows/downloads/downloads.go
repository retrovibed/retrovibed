package downloads

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
)

func NewDirectoryWatcher(ctx context.Context, c *tls.Config, q sqlx.Queryer) (d Directory, err error) {
	var (
		w *fsnotify.Watcher
	)

	if w, err = fsnotify.NewWatcher(); err != nil {
		return d, err
	}

	d = Directory{
		c: c,
		w: w,
		q: q,
		ignore: func(s string) bool {
			return !strings.HasSuffix(s, ".torrent")
		},
	}
	d.pool = asynccompute.New(d.download)

	return d.background(ctx), nil
}

type Directory struct {
	c      *tls.Config
	pool   *asynccompute.Pool[string]
	w      *fsnotify.Watcher
	q      sqlx.Queryer
	ignore func(s string) bool
}

func (t Directory) Add(path string) (err error) {
	defer func() {
		if err == nil {
			log.Printf("download folder monitoring enabled - %s\n", path)
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

		if t.ignore(path) {
			return nil
		}

		if err = t.pool.Run(context.Background(), filepath.Join(path, name)); err != nil {
			return err
		}

		return nil
	})

	return errorsx.Wrap(err, "unable to find existing torrents")
}

// background download
func (t Directory) download(ctx context.Context, path string) error {
	var (
		err     error
		req     *http.Request
		resp    *http.Response
		decoded media.MediaUploadResponse
		lib     meta.Daemon
	)

	if err := meta.DaemonFindByDownload(ctx, t.q).Scan(&lib); err != nil {
		return errorsx.Wrap(err, "failed to lookup the library to send downloads to")
	}

	c := authn.AutoOauth2Client(ctx, t.c, authn.EndpointSSHAuth(fmt.Sprintf("https://%s", lib.Hostname)))
	log.Println("sending download request to", lib.Hostname)

	mimetype, content, err := httpx.Multipart(func(w *multipart.Writer) error {
		src, err := os.Open(path)
		if err != nil {
			return errorsx.Wrapf(err, "unable to read %s", path)
		}
		defer src.Close()

		part, lerr := w.CreatePart(httpx.NewMultipartHeader(mimex.Binary, "content", filepath.Base(path)))
		if lerr != nil {
			return errorsx.Wrap(lerr, "unable to create archive part")
		}

		if _, lerr = io.Copy(part, src); lerr != nil {
			return errorsx.Wrap(lerr, "unable to copy archive")
		}

		return nil
	})
	if err != nil {
		return errorsx.Wrap(err, "unable to build multipart request")
	}

	if req, err = http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://%s/d/", lib.Hostname), content); err != nil {
		return errorsx.Wrap(err, "unable to create http request")
	}
	req.Header.Add("Content-Type", mimetype)

	if resp, err = httpx.AsError(c.Do(req)); err != nil {
		return errorsx.Wrap(err, "http request failed")
	}

	if err = httpx.DecodeJSON(resp, &decoded); err != nil {
		return errorsx.Wrap(err, "unable to decode response")
	}

	return nil
}

func (t Directory) background(ctx context.Context) Directory {
	go func() {
		freq := time.Second
		flush := time.NewTicker(freq)
		defer flush.Stop()

		pending := make(map[string]struct{})
		for {
			select {
			case evt := <-t.w.Events:
				flush.Reset(freq)
				if t.ignore(evt.Name) {
					continue
				}
				switch evt.Op {
				case fsnotify.Create:
				case fsnotify.Chmod:
				case fsnotify.Write:
					continue // explicitly ignored to reduce noise.
				default:
					log.Println("change ignored", evt.Op)
					continue
				}

				pending[evt.Name] = struct{}{}
			case err := <-t.w.Errors:
				log.Println("watch error", err)
			case <-flush.C:
				for path := range pending {
					errorsx.Log(errorsx.Wrapf(t.pool.Run(ctx, path), "unable to enqueue %s for download", path))
				}
				pending = make(map[string]struct{})
			case <-ctx.Done():
				log.Println("context completed", ctx.Err())
				return
			}
		}
	}()

	return t
}
