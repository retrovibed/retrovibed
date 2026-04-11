package fsnotifyx

import (
	"context"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/fsx"
)

// runs the action immediately if the file exists and whenever an event occurs.
func OnceAndOnChange(ctx context.Context, path string, action func(ctx context.Context, evt fsnotify.Event) error) error {
	if fsx.Exists(path) {
		if err := action(ctx, fsnotify.Event{Name: path, Op: fsnotify.Chmod}); err != nil {
			return err
		}
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err = w.Add(path); err != nil {
		return errorsx.Wrap(err, "failed to add path")
	}

	go func() {
		defer w.Close()
		for {
			select {
			case evt := <-w.Events:
				if err = action(ctx, evt); err != nil {
					log.Println("file watch failed", evt.Name, err)
				}
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
