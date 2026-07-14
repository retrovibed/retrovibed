package searchplugin

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/retrovibed/retrovibed/retroapi/internal/errorsx"
)

// watch scans dir for pre-existing *.wasm files (loading each into r), then
// watches dir for further Create/Remove events, loading/unloading plugins
// as they come and go, for the lifetime of ctx.
func watch(ctx context.Context, r *Registry, dir string) error {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return errorsx.Wrapf(err, "unable to create search plugin directory: %s", dir)
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err := w.Add(dir); err != nil {
		return errorsx.Wrapf(err, "unable to watch search plugin directory: %s", dir)
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
		return errorsx.Wrap(err, "unable to load existing search plugins")
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
					errorsx.Log(errorsx.Wrapf(r.Load(ctx, evt.Name), "unable to load search plugin: %s", evt.Name))
				case evt.Op&fsnotify.Remove != 0, evt.Op&fsnotify.Rename != 0:
					r.Unload(evt.Name)
				}
			case werr, ok := <-w.Errors:
				if !ok {
					return
				}
				log.Println("search plugin directory watch error", werr)
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
