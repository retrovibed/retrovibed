package media

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/linxGnu/pqueue"
	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/pqueuex"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

func RecommendationsBackgroundRun(ctx context.Context, q sqlx.Queryer, wq pqueue.Queue) error {
	if last, err := library.RecommendationLastGeneratedAt(ctx, q, library.RecommendationSourceRandom); err != nil {
		return errorsx.Wrap(err, "recommendations background failed to get last generated at")
	} else if time.Since(last) < 24*time.Hour {
		log.Println("random recommendation last ran", time.Since(last), "ago at", last)
		return nil
	}

	lang := userx.LocaleLanguage()
	reclimit := uint64(128)

	reqaudio := RecommendationRefreshRequest{
		ProfileId: uuid.Nil.String(),
		Mimetype:  mimex.Audio,
		Adult:     false,
		Language:  lang,
		Limit:     reclimit,
	}
	reqvideo := RecommendationRefreshRequest{
		ProfileId: uuid.Nil.String(),
		Mimetype:  mimex.Video,
		Adult:     false,
		Language:  lang,
		Limit:     reclimit,
	}

	return errors.Join(
		errorsx.Wrap(pqueuex.Enqueue(ctx, wq, &reqaudio), "failed to enqueue audio recommendation request"),
		errorsx.Wrap(pqueuex.Enqueue(ctx, wq, &reqvideo), "failed to enqueue video recommendation request"),
	)
}

func RecommendationsBackground(ctx context.Context, seed string, q sqlx.Queryer, wq pqueue.Queue, p searchplugin.R) error {
	wakeup := asyncx.NewWakeup(ctx)
	s := backoffx.New(
		backoffx.Frequency(12*time.Hour, seed),
		backoffx.JitterRandom(time.Minute),
	)

	go contextx.RunContext(ctx, pqueuex.NewWorker(wq, NewRecommendationBackgroundWorker(q, p)).Consume)
	go asyncx.Periodic(ctx, wakeup, s, "recommendations background")
	contextx.Run(ctx, func() {
		errorsx.Log(asyncx.Run(ctx, wakeup, func(ctx context.Context) error {
			return RecommendationsBackgroundRun(ctx, q, wq)
		}))
	})

	return nil
}
