package cmdmedia

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"
	"github.com/retrovibed/retrovibed/shallows/internal/asynccompute"
	"github.com/retrovibed/retrovibed/shallows/internal/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/fsx"
	"github.com/retrovibed/retrovibed/shallows/internal/jsonl"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
	"github.com/retrovibed/retrovibed/shallows/internal/tarx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

type knownarchive struct {
	Directory string `flag:"" name:"directory" help:"work directory for the command"`
	Pattern   string `flag:"" name:"pattern" help:"name of the archive directory to create" default:"retrovibed.media.archive.*.d"`
	DryRun    bool   `flag:"" name:"dry-run" help:"do not actually write the data into the archive just log" negatable:"" default:"false"`
}

func (t knownarchive) Run(gctx *cmdopts.Global) (err error) {
	var (
		dir  string
		v    library.Known
		derr error
	)

	d := jsonl.NewDecoder(os.Stdin)

	if strings.Contains(t.Pattern, "*") {
		dir, err = os.MkdirTemp(t.Directory, t.Pattern)
	} else {
		dir = filepath.Join(stringsx.FirstNonBlank(t.Directory, "."), t.Pattern)
		err = os.MkdirAll(dir, 0700)
	}
	if err != nil {
		return errorsx.Wrap(err, "unable to create archive directory")
	}

	unpack := func(path string) error {
		archive := filepath.Join(dir, path)
		dst := strings.TrimPrefix(path, "retrovibed.media.metadata.archive.")
		dst = strings.TrimSuffix(dst, ".tar.gz")
		src, err := os.Open(archive)
		if err != nil {
			return errorsx.Wrap(err, "failed to open archive")
		}
		defer src.Close()

		if err = os.MkdirAll(filepath.Join(dir, dst), 0700); err != nil {
			return errorsx.Wrapf(err, "failed to prepare archive directory: %s", path)
		}

		if err = tarx.Unpack(filepath.Join(dir, dst), src); err != nil {
			return errorsx.Wrapf(err, "failed to unpack archive: %s", path)
		}

		if err = os.Remove(archive); err != nil {
			return errorsx.Wrapf(err, "failed to remove extracted archive")
		}

		return nil
	}

	w := fsx.Walk(os.DirFS(dir))
	for path := range w.Walk() {
		if path == "." {
			continue
		}

		if !strings.HasSuffix(path, ".tar.gz") {
			log.Println("skipping", path)
			continue
		}

		if err = unpack(path); err != nil {
			return err
		}
	}

	if err := w.Err(); err != nil {
		return errorsx.Wrap(err, "unable to walk directory")
	}

	for derr = d.Decode(&v); derr == nil; derr = d.Decode(&v) {
		v.AutoDescription = stringsx.Join("\n", v.Title, v.OriginalTitle, v.Overview)
		encoded, err := json.Marshal(v)
		if err != nil {
			return errorsx.Wrapf(err, "unable to encode record %s", v.UID)
		}

		adir := filepath.Join(dir, fmt.Sprintf("%02x", backoffx.DynamicHashWindow(v.Released.String(), 32)))
		if err := fsx.MkDirs(0700, adir); err != nil {
			return errorsx.Wrapf(err, "unable to make directory: %s %s", v.UID, adir)
		}

		path := filepath.Join(adir, fmt.Sprintf("%x", backoffx.DynamicHashWindow(v.UID, 16)))

		if t.DryRun {
			log.Println(string(encoded))
		} else {
			if err := fsx.AppendTo(path, 0600, encoded, []byte("\n")); err != nil {
				return errorsx.Wrapf(err, "unable to append record %s %s", v.UID, path)
			}
		}
	}

	if errorsx.Ignore(derr, io.EOF) != nil {
		return derr
	}

	compressors := asynccompute.New(func(ctx context.Context, path string) error {
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(filepath.Join(dir, fmt.Sprintf("retrovibed.media.metadata.archive.%s.tar.gz", filepath.Base(path))))
		if err != nil {
			return err
		}

		return errorsx.Compact(tarx.Pack(out, path), out.Close(), os.RemoveAll(path))
	})

	w = fsx.Walk(os.DirFS(dir))

	for path := range w.Walk() {
		if path == "." {
			continue
		}

		if err := compressors.Run(gctx.Context, filepath.Join(dir, path)); err != nil {
			return err
		}
	}

	return errorsx.Compact(w.Err(), asynccompute.Shutdown(gctx.Context, compressors))
}
