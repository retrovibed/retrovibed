package cmdlibrary

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/asynccompute"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/httpx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/mimex"
	"github.com/retrovibed/retrovibed/shallows/media"
)

type importDirectory struct {
	Endpoint    string `flag:"" name:"peer" help:"http address for the daemon you want to import to" default:"localhost:9998"`
	Concurrency uint16 `flag:"" name:"concurrency" help:"number of files to upload concurrently, defaults to the number of cpus" default:"${vars_cores}"`
	Mimetype    string `flag:"" name:"mimetype" help:"override the mimetype for all uploaded files" optional:""`
	Directory   string `arg:"" name:"directory" help:"directory to import; each immediate file is uploaded to the library"`
}

func (t importDirectory) Run(gctx *cmdopts.Global, tls *cmdopts.TLSConfig) error {
	c := authn.AutoOauth2Client(gctx.Context, tls.Config(), authn.EndpointSSHAuth(fmt.Sprintf("https://%s", t.Endpoint)))
	return t.run(gctx.Context, jsonl.NewEncoder(os.Stdout), c)
}

func (t importDirectory) run(ctx context.Context, enc *jsonl.Encoder, c *http.Client) error {
	type Workload struct {
		Path  string
		Entry fs.DirEntry
	}

	log.Println("directory import initiated")
	defer log.Println("directory import completed")

	encwriter := asynccompute.New(func(ctx context.Context, w *media.Media) (err error) {
		return enc.Encode(w)
	}, asynccompute.Workers[*media.Media](1), asynccompute.Backlog[*media.Media](t.Concurrency*2))

	publisher := asynccompute.New(func(ctx context.Context, w Workload) (err error) {
		var (
			m = new(media.MediaUploadResponse)
		)
		defer func() {
			errorsx.Log(errorsx.Wrapf(err, "import failed: %s", w.Path))
		}()

		f, err := os.Open(w.Path)
		if err != nil {
			return err
		}
		defer f.Close()

		filename := strings.TrimPrefix(w.Path, filepath.Dir(strings.TrimSuffix(t.Directory, string(filepath.Separator))))
		filename = strings.ReplaceAll(filename, string(filepath.Separator), " ")
		mimetype := langx.FirstNonZero(t.Mimetype, mime.TypeByExtension(filepath.Ext(filename)), mimex.Binary)

		contentType, body, err := httpx.Multipart(func(mw *multipart.Writer) error {
			part, lerr := mw.CreatePart(httpx.NewMultipartHeader(mimetype, "content", filename))
			if lerr != nil {
				return lerr
			}
			_, lerr = io.Copy(part, f)
			return lerr
		})
		if err != nil {
			return err
		}
		defer body.Close()

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://%s/m/", t.Endpoint), body)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", contentType)

		resp, err := httpx.AsError(c.Do(req))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(m); err != nil {
			return err
		}

		return encwriter.Run(ctx, m.Media)
	}, asynccompute.Workers[Workload](t.Concurrency))

	w := fsx.Walk(os.DirFS(t.Directory))
	for p, e := range w.Walk() {
		if p == "." || e.IsDir() {
			continue
		}

		if err := publisher.Run(ctx, Workload{Path: filepath.Join(t.Directory, p), Entry: e}); err != nil {
			return errorsx.Wrap(err, "unable to enqueue upload workload")
		}
	}

	if err := w.Err(); err != nil {
		return errorsx.Wrapf(err, "failed to import directory: %s", t.Directory)
	}

	return langx.FirstNonNil(
		asynccompute.Shutdown(ctx, publisher),
		asynccompute.Shutdown(ctx, encwriter),
	)
}
