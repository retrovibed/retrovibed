package cmdmedia

import (
	"archive/tar"
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/tarx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type tarchiveexport struct {
	Directory string `flag:"" name:"directory" help:"work directory for the command, defaults to current working directory, usually shouldnt be needed"`
	Pattern   string `flag:"" name:"pattern" help:"name of the archive directory to import" default:"retrovibed.media.archive.d"`
}

func (t tarchiveexport) Run(kctx *kong.Context, gctx *cmdopts.Global) (err error) {
	return t.run(gctx.Context, kctx.Stdout)
}

func (t tarchiveexport) run(ctx context.Context, out io.Writer) (err error) {
	encoder := jsonl.NewEncoder(out)

	insert := asynccompute.New(func(ctx context.Context, v library.Known) error {
		return encoder.Encode(v)
	}, asynccompute.Workers[library.Known](1))

	pool := asynccompute.New(func(ctx context.Context, path string) error {
		archive, err := os.Open(path)
		if err != nil {
			return errorsx.Wrap(err, "unable to open read archive")
		}
		defer archive.Close()

		iter, err := tarx.UnpackSeq(archive)
		if err != nil {
			return errorsx.Wrap(err, "unable to open read archive")
		}

		importtarfile := func(_ *tar.Header, content *tar.Reader) error {
			var (
				derr error
				i    uint64
				v    library.Known
			)

			d := jsonl.NewDecoder(content)

			for derr = d.Decode(&v); derr == nil; i, derr = i+1, d.Decode(&v) {
				if err := insert.Run(
					ctx,
					langx.Clone(
						v,
						timex.JSONSafeDecodeOption,
						library.KnownOptionAutoDescription,
					),
				); err != nil {
					return err
				}
			}

			if err := errorsx.Ignore(derr, io.EOF); err != nil {
				return err
			}

			return nil
		}

		for header, content := range iter {
			errorsx.Log(importtarfile(header, content))
		}

		return nil
	})

	w := fsx.WalkDir(os.DirFS(t.Pattern))
	for path := range w.Walk() {
		if path == "." {
			continue
		}

		if !strings.HasSuffix(path, ".tar.gz") {
			log.Println("skipping", path)
			continue
		}

		if err := pool.Run(ctx, path); err != nil {
			return err
		}
	}

	if err := w.Err(); err != nil {
		return errorsx.Wrap(err, "unable to walk directory")
	}

	return asynccompute.Shutdown(ctx, pool, insert)
}
