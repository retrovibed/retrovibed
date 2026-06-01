package daemons

import (
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"time"

	"github.com/retrovibed/retrovibed/shallows/acoustics"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

// AcousticsBackgroundRun indexes one unindexed audio track per invocation.
func AcousticsBackgroundRun(ctx context.Context, q sqlx.Queryer, media fs.FS) error {
	ids, err := acoustics.UnindexedMediaIDs(ctx, q, 1)
	if err != nil {
		return errorsx.Wrap(err, "acoustics: find unindexed")
	}
	if len(ids) == 0 {
		return nil
	}

	id := ids[0]

	// Copy media from VFS to a temp file for FFmpeg (which requires a seekable path).
	tmpPath, err := copyToTemp(media, id.String())
	if err != nil {
		return errorsx.Wrap(err, "acoustics: copy to temp")
	}
	defer func() {
		errorsx.Log(errorsx.Wrap(os.Remove(tmpPath), "acoustics: cleanup temp"))
	}()

	vec, err := acoustics.AnalyzeFile(ctx, tmpPath)
	if err != nil {
		return errorsx.Wrap(err, "acoustics: analyze file")
	}

	if err = acoustics.StoreFeatures(ctx, q, id, vec, acoustics.StatsVersion); err != nil {
		return errorsx.Wrap(err, "acoustics: store features")
	}

	log.Println("acoustics: indexed", id)
	return nil
}

// AcousticsBackground runs the acoustic indexer on a 1-second interval with jitter.
func AcousticsBackground(ctx context.Context, q sqlx.Queryer, media fs.FS) error {
	wakeup := asyncx.NewWakeup(ctx)
	s := backoffx.New(
		backoffx.Constant(time.Second),
		backoffx.Jitter(0.1),
	)

	go asyncx.Periodic(ctx, wakeup, s, "acoustics background indexer")
	contextx.Run(ctx, func() {
		errorsx.Log(asyncx.Run(ctx, wakeup, func(ctx context.Context) error {
			return AcousticsBackgroundRun(ctx, q, media)
		}))
	})

	return nil
}

func copyToTemp(media fs.FS, name string) (string, error) {
	src, err := media.Open(name)
	if err != nil {
		return "", err
	}
	defer src.Close()

	tmp, err := os.CreateTemp("", "acoustics-*.audio")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tmp, src)
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}

	if err != nil {
		errorsx.Log(errorsx.Wrap(os.Remove(tmp.Name()), "acoustics: cleanup temp"))
		return "", err
	}
	return tmp.Name(), nil
}
