package media

import (
	"context"
	"database/sql"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/retrovibed/retrovibed/retroapi/jsonx"
	"github.com/retrovibed/retrovibed/retroapi/searchplugin"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/sqlx"
)

func NewRecommendationBackgroundWorker(q sqlx.Queryer, plugins searchplugin.R) RecommendationBackgroundWorker {
	return RecommendationBackgroundWorker{
		q:       q,
		plugins: plugins,
	}
}

type RecommendationBackgroundWorker struct {
	q       sqlx.Queryer
	plugins searchplugin.R
}

func (t RecommendationBackgroundWorker) Message(ctx context.Context, m []byte) (err error) {
	var (
		decoded RecommendationRefreshRequest
	)

	if err = jsonx.Unmarshal(m, &decoded); err != nil {
		return err
	}

	log.Println("recommendations", decoded.Mimetype, "generation initiated")
	defer log.Println("recommendations", decoded.Mimetype, "generation completed")

	log.Println("DERP DERP", spew.Sdump(&decoded))
	if _, err := Recommendation(ctx, t.q, decoded.Mimetype, decoded.Language, decoded.Adult); errorsx.Ignore(err, sql.ErrNoRows) != nil {
		return errorsx.Wrap(err, "recommendations audio background failed to generate recommendation")
	}

	if err := RecommendationsFromPlugins(ctx, t.q, t.plugins, decoded.Mimetype, uint(decoded.Limit), decoded.Language, decoded.Adult); err != nil {
		return errorsx.Wrapf(err, "recommendations background failed to generate %s recommendations from search plugins", decoded.Mimetype)
	}

	return nil
}
