package cmdtorrent

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
	"github.com/james-lawrence/torrent/storage"
	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/asynccompute"
	"github.com/retrovibed/retrovibed/shallows/internal/debugx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/media"
	"github.com/retrovibed/retrovibed/shallows/meta"
)

type importDirectory struct {
	Endpoint  string        `flag:"" name:"peer" help:"http address for the daemon you want to import to, usually you do this on device" default:"localhost:9998"`
	Entropy   string        `flag:"" name:"entropy" help:"encryption entropy value to set for uploads" default:""`
	Mimetype  string        `flag:"" name:"mimetype" help:"mimetype of the media if known and consistent"`
	TTL       time.Duration `flag:"" name:"ttl" help:"when the torrents from this import should be marked for expiration, which means they're no longer valid"`
	Directory string        `arg:"" name:"directory" help:"directory containing the content to import. each immediate file / directory generates a torrent. subdirectories create multi file torrents"`
}

func (t importDirectory) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig) error {
	type Workload struct {
		Path  string
		Entry fs.DirEntry
	}

	var (
		com meta.Community
	)

	encoder := jsonl.NewEncoder(os.Stdout)

	if cmdopts.Readable(os.Stdin) {
		debugx.Println("reading community from stdin")
		if err := json.NewDecoder(os.Stdin).Decode(&com); err != nil {
			return err
		}

		// echo it back to stdout so next stream can read it.
		if err := encoder.Encode(&com); err != nil {
			return err
		}
	} else {
		com = meta.Community{
			Entropy:  t.Entropy,
			Mimetype: t.Mimetype,
		}
	}

	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)))

	publisher := asynccompute.New(func(ctx context.Context, w Workload) (err error) {
		var (
			m = new(media.PublishedUploadResponse)
		)
		defer func() {
			errorsx.Log(errorsx.Wrapf(err, "import failed: %s", w.Path))
		}()

		info, err := metainfo.NewFromPath(w.Path)
		if err != nil {
			return err
		}

		md, err := torrent.NewFromInfo(info, torrent.OptionStorage(storage.NewFile(filepath.Dir(w.Path), storage.FileOptionPathMakerFixed(filepath.Base(w.Path)))), torrent.OptionDisplayName(info.Name))
		if err != nil {
			return err
		}

		mimetype, data, err := media.PublishRequest(ctx, md, &media.PublishedUploadRequest{Ttl: uint64(t.TTL), Entropy: com.Entropy, Mimetype: com.Mimetype})
		if err != nil {
			return err
		}
		defer data.Close()

		req, err := http.NewRequestWithContext(gctx.Context, http.MethodPost, fmt.Sprintf("https://%s/d/publish", t.Endpoint), data)
		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", mimetype)

		resp, err := httpx.AsError(c.Do(req))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(m); err != nil {
			return err
		}

		return encoder.Encode(m.Published)
	})

	w := fsx.WalkDir(os.DirFS(t.Directory))
	for p, e := range w.Walk() {
		if p == "." {
			continue
		}

		if err := publisher.Run(gctx.Context, Workload{Path: filepath.Join(t.Directory, p), Entry: e}); err != nil {
			return errorsx.Wrap(err, "unable to enqueue publish workload")
		}
	}

	if err := w.Err(); err != nil {
		return errorsx.Wrapf(err, "failed to import directory: %s", t.Directory)
	}

	if err := asynccompute.Shutdown(gctx.Context, publisher); err != nil {
		return err
	}

	return nil
}
