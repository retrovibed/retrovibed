package publishplugin

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
)

// watch scans dir for pre-existing *.wasm files (loading each into r), then
// watches dir for further Create/Remove events, loading/unloading plugins
// as they come and go, for the lifetime of ctx.
func watch(ctx context.Context, r *Registry, dir string) error {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return errorsx.Wrapf(err, "unable to create publish plugin directory: %s", dir)
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err := w.Add(dir); err != nil {
		return errorsx.Wrapf(err, "unable to watch publish plugin directory: %s", dir)
	}

	err = fs.WalkDir(os.DirFS(dir), ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !isWasm(name) {
			return nil
		}
		return r.Load(ctx, filepath.Join(dir, name))
	})
	if err != nil {
		return errorsx.Wrap(err, "unable to load existing publish plugins")
	}

	go func() {
		defer w.Close()
		for {
			select {
			case evt, ok := <-w.Events:
				if !ok {
					return
				}
				if !isWasm(evt.Name) {
					continue
				}
				switch {
				case evt.Op&fsnotify.Create != 0:
					errorsx.Log(errorsx.Wrapf(r.Load(ctx, evt.Name), "unable to load publish plugin: %s", evt.Name))
				case evt.Op&fsnotify.Remove != 0, evt.Op&fsnotify.Rename != 0:
					r.Unload(evt.Name)
				}
			case werr, ok := <-w.Errors:
				if !ok {
					return
				}
				log.Println("publish plugin directory watch error", werr)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func isWasm(name string) bool {
	return strings.HasSuffix(name, ".wasm")
}

// SanitizeName strips any character unsafe for use as a single path
// component under the publish plugin directory - keeping only letters,
// digits, '-', and '_' - so the result cannot contain "..", an absolute
// path, or a path separator regardless of what an HTTP client supplies.
func SanitizeName(name string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			return r
		default:
			return -1
		}
	}, name)
}
