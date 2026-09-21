package cmdopts

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/egdaemon/gdx"
	"github.com/retrovibed/retrovibed/retroapi/userx"
)

// Trace records a runtime trace for the entire run of the cli.
type Trace struct {
	Enabled bool `flag:"" name:"trace" help:"record a runtime trace for the entire run to the cache directory" default:"false" env:"${env_trace}"`
	path    string
	stop    context.CancelFunc
	done    chan struct{}
}

func (t *Trace) AfterApply() error {
	if !t.Enabled {
		return nil
	}

	// the trace uses its own context, it must outlive the shutdown of the main context
	// to capture the cleanup of the run. stopped by Close.
	ctx, stop := context.WithCancel(context.Background())
	t.stop = stop
	t.done = make(chan struct{})
	t.path = filepath.Join(
		userx.DefaultCacheDirectory(userx.DefaultRelRoot()),
		"traces",
		fmt.Sprintf("retrovibe.%d.%s.trace", os.Getpid(), time.Now().UTC().Format("20060102T150405Z")),
	)

	log.Println("recording trace", t.path)

	go func(path string, done chan struct{}) {
		defer close(done)
		// gdx captures are capped by a duration (default 30s), effectively disable it.
		// cancellation is how the trace is stopped, not a failure.
		if err := gdx.RecordFile(ctx, path, gdx.Trace(ctx, gdx.ProfileOptions().ProfileDuration(365*24*time.Hour).Compact()...)); err != nil && !errors.Is(err, context.Canceled) {
			log.Println("unable to record trace", err)
		}
	}(t.path, t.done)

	return nil
}

// Close stops the trace and waits for it to be flushed to disk.
func (t *Trace) Close() {
	if t.stop == nil {
		return
	}

	t.stop()
	<-t.done
	t.stop = nil

	log.Println("trace written", t.path)
}
