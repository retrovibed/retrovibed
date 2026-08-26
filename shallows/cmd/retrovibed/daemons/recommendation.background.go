package daemons

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/backoffx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

func RecommendationsBackgroundRun(ctx context.Context, q sqlx.Queryer) error {
	last, err := library.RecommendationLastGeneratedAt(ctx, q, library.RecommendationSourceRandom)
	if err != nil {
		return errorsx.Wrap(err, "recommendations background failed to get last generated at")
	}

	if time.Since(last) < 24*time.Hour {
		log.Println("random recommendation last ran", time.Since(last), "ago at", last)
		return nil
	}

	lang := userx.LocaleLanguage()

	if _, err = ddisc.Recommendation(ctx, q, mimex.Audio, lang, false); errorsx.Ignore(err, sql.ErrNoRows) != nil {
		return errorsx.Wrap(err, "recommendations audio background failed to generate recommendation")
	} else if err == nil {
		log.Println("recommendations background generated recommendation")
	}

	if _, err = ddisc.Recommendation(ctx, q, mimex.Video, lang, false); errorsx.Ignore(err, sql.ErrNoRows) != nil {
		return errorsx.Wrap(err, "recommendations video background failed to generate recommendation")
	} else if err == nil {
		log.Println("recommendations background generated recommendation")
	}

	return nil
}

func RecommendationsBackground(ctx context.Context, q sqlx.Queryer) error {
	wakeup := asyncx.NewWakeup(ctx)
	s := backoffx.New(
		backoffx.Constant(time.Hour),
		backoffx.Jitter(0.1),
	)

	go asyncx.Periodic(ctx, wakeup, s, "recommendations background")
	contextx.Run(ctx, func() {
		errorsx.Log(asyncx.Run(ctx, wakeup, func(ctx context.Context) error {
			return RecommendationsBackgroundRun(ctx, q)
		}))
	})

	return nil
}
