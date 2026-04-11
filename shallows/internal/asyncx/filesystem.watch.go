package asyncx

import (
	"context"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
)

func FileCreated(evt fsnotify.Event) bool {
	return evt.Op == fsnotify.Create
}

// watch the provided file paths. touching the files if they do not exist initially.
func WatchFiles(ctx context.Context, wu *Wakeup, do func(fsnotify.Event) bool, paths ...string) error {
	if err := fsx.Touch(0600, paths...); err != nil {
		return err
	}

	return watcher(ctx, wu, do, paths...)
}

// watch the provided file paths. making the directories if they do not exist initially.
func WatchDirectories(ctx context.Context, wu *Wakeup, do func(fsnotify.Event) bool, paths ...string) error {
	if err := fsx.MkDirs(0700, paths...); err != nil {
		return err
	}

	return watcher(ctx, wu, do, paths...)
}

func watcher(ctx context.Context, wu *Wakeup, do func(fsnotify.Event) bool, paths ...string) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	addpath := func(path string) {
		if err = w.Add(path); err != nil {
			errorsx.Log(errorsx.Wrapf(err, "unable to watch %s", path))
			return
		}
	}

	for _, path := range paths {
		addpath(path)
	}

	go func() {
		defer w.Close()

		for {
			select {
			case evt := <-w.Events:
				if !do(evt) {
					continue
				}
				wu.Broadcast()
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
