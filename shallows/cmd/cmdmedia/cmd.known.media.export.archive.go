package cmdmedia

import (
	"archive/tar"
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/retrovibed/retrovibed/retroapi/asynccompute"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/tarx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type tarchiveexport struct {
	Directory string `flag:"" name:"directory" help:"work directory for the command, defaults to current working directory, usually shouldnt be needed"`
	Pattern   string `flag:"" name:"pattern" help:"name of the archive directory to import" default:"retrovibed.media.archive.d"`
}

func (t tarchiveexport) Run(gctx *cmdopts.Global) (err error) {
	encoder := jsonl.NewEncoder(os.Stdout)

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
				v.AutoDescription = stringsx.Join("\n", v.Title, v.OriginalTitle, v.Overview)
				return insert.Run(ctx, v)
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

		if err := pool.Run(gctx.Context, path); err != nil {
			return err
		}
	}

	if err := w.Err(); err != nil {
		return errorsx.Wrap(err, "unable to walk directory")
	}

	if err := asynccompute.Shutdown(gctx.Context, pool); err != nil {
		return err
	}

	return asynccompute.Shutdown(gctx.Context, insert)
}
