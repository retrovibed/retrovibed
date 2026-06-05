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
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

// AcousticsBackgroundRun drains the unindexed audio backlog.
func AcousticsBackgroundRun(ctx context.Context, q sqlx.Queryer, media fs.FS) error {
	v := sqlx.Scan(acoustics.AudioFeaturesUnindexedMediaIDs(ctx, q))

	for m := range v.Iter() {
		log.Println("acousticly indexing initiated", m.ID, m.Description)
		errorsx.Log(acousticsIndexOne(ctx, q, media, m.ID))
		log.Println("acousticly indexing completed", m.ID, m.Description)
	}

	if err := v.Err(); err != nil {
		return errorsx.Wrap(err, "acoustics: find unindexed")
	}

	return nil
}

func acousticsIndexOne(ctx context.Context, q sqlx.Queryer, media fs.FS, id string) error {
	tmpPath, err := copyToTemp(media, id)
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

// AcousticsBackground drains the index, then polls for new media on an
// exponential backoff that maxes out at an hour.
func AcousticsBackground(ctx context.Context, q sqlx.Queryer, media fs.FS) error {
	wakeup := asyncx.NewWakeup(ctx)
	defer wakeup.Broadcast() // kick off an initial indexing process
	s := backoffx.New(
		backoffx.Exponential(time.Second),
		backoffx.Maximum(time.Hour),
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

func copyToTemp(media fs.FS, name string) (_ string, err error) {
	src, err := media.Open(name)
	if err != nil {
		return "", err
	}
	defer src.Close()

	tmp, err := os.CreateTemp("", "acoustics-*.audio")
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			errorsx.Log(errorsx.Wrap(os.Remove(tmp.Name()), "acoustics: cleanup temp"))
		}
	}()

	_, err = io.Copy(tmp, src)
	err = langx.FirstNonZero(err, tmp.Close())
	if err != nil {
		return "", err
	}

	return tmp.Name(), nil
}
