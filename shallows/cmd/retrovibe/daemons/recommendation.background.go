package daemons

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/backoffx"
	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/ddisc"
	"github.com/retrovibed/retrovibed/shallows/internal/asyncx"
	"github.com/retrovibed/retrovibed/shallows/internal/contextx"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
	"github.com/retrovibed/retrovibed/shallows/library"
)

func RecommendationsBackgroundRun(ctx context.Context, q sqlx.Queryer, plugins searchplugin.R) error {
	if last, err := library.RecommendationLastGeneratedAt(ctx, q, library.RecommendationSourceRandom); err != nil {
		return errorsx.Wrap(err, "recommendations background failed to get last generated at")
	} else if time.Since(last) < 24*time.Hour {
		log.Println("random recommendation last ran", time.Since(last), "ago at", last)
		return nil
	}

	// TODO: recommendation settings to be loaded from database
	lang := userx.LocaleLanguage()
	reclimit := uint(5)

	if _, err := ddisc.Recommendation(ctx, q, mimex.Audio, lang, false); errorsx.Ignore(err, sql.ErrNoRows) != nil {
		return errorsx.Wrap(err, "recommendations audio background failed to generate recommendation")
	} else if err == nil {
		log.Println("recommendations background generated audio recommendation")
	}

	if _, err := ddisc.Recommendation(ctx, q, mimex.Video, lang, false); errorsx.Ignore(err, sql.ErrNoRows) != nil {
		return errorsx.Wrap(err, "recommendations video background failed to generate recommendation")
	} else if err == nil {
		log.Println("recommendations background generated video recommendation")
	}

	for _, mimetype := range []string{mimex.Audio, mimex.Video} {
		if err := ddisc.RecommendationsFromPlugins(ctx, q, plugins, mimetype, reclimit, lang, false); err != nil {
			return errorsx.Wrapf(err, "recommendations background failed to generate %s recommendations from search plugins", mimetype)
		}
		log.Println("recommendations background generated", mimetype, "recommendations from search plugins")
	}

	return nil
}

func RecommendationsBackground(ctx context.Context, q sqlx.Queryer, plugins searchplugin.R, seed string) error {
	wakeup := asyncx.NewWakeup(ctx)
	s := backoffx.New(
		backoffx.Frequency(12*time.Hour, seed),
		backoffx.JitterRandom(time.Minute),
	)

	go asyncx.Periodic(ctx, wakeup, s, "recommendations background")
	contextx.Run(ctx, func() {
		errorsx.Log(asyncx.Run(ctx, wakeup, func(ctx context.Context) error {
			return RecommendationsBackgroundRun(ctx, q, plugins)
		}))
	})

	return nil
}
