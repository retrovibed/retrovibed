package fsx

import (
	"context"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
)

func Watch(ctx context.Context, cond *sync.Cond, paths ...string) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	for _, path := range paths {
		log.Println("added path", path)
		errorsx.Log(errorsx.Wrapf(w.Add(path), "unable to watch %s", path))
	}

	go func() {
		defer log.Println("file watcher done")
		defer w.Close()
		for {
			select {
			case evt := <-w.Events:
				log.Println("received event", evt.Op, evt.Name)
				cond.Broadcast()
			case err := <-w.Errors:
				log.Println("watch error", err)
			case <-ctx.Done():
				log.Println("context completed", ctx.Err())
				return
			}
		}
	}()

	return nil
}
