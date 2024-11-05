package downloads

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/james-lawrence/deeppool/internal/x/errorsx"
)

func NewDirectoryWatcher() (d Directory, err error) {
	var (
		w *fsnotify.Watcher
	)

	if w, err = fsnotify.NewWatcher(); err != nil {
		return d, err
	}

	return Directory{w: w}, nil
}

type Directory struct {
	w *fsnotify.Watcher
}

func (t Directory) Add(path string) (err error) {
	defer func() {
		if err == nil {
			log.Println("watching", path)
		}
	}()
	return errorsx.Wrapf(t.w.Add(path), "unable to watch: %s", path)
}

func (t Directory) background() Directory {
	go func() {
		for {
			select {
			case evt := <-t.w.Events:
				if evt.Op == fsnotify.Chmod {
					continue
				}

				log.Println("change detected", evt.Op)
			case err := <-t.w.Errors:
				log.Println("watch error", err)
			}
		}
	}()
	return t
}
