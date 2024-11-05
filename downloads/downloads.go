package downloads

import (
	"log"

	"github.com/davecgh/go-spew/spew"
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

	return Directory{w: w}.background(), nil
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
				switch evt.Op {
				case fsnotify.Create:
					// do nothing fallthrough.
				default:
					log.Println("change ignored", evt.Op)
					continue
				}

				log.Println("change detected", spew.Sdump(evt))
			case err := <-t.w.Errors:
				log.Println("watch error", err)
			}
		}
	}()
	return t
}
