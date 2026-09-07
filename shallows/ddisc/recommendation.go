package ddisc

import (
	"time"

	"github.com/retrovibed/retrovibed/retroapi/mimex"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/md5x"
	"github.com/retrovibed/retrovibed/shallows/library"
)

// RecommendationFromDiscovered maps a located Discovered row into a
// library.Recommendation keyed on d's own ddisc_media row. Pure mapping -
// does not touch the database; callers are responsible for persisting it
// via library.RecommendationInsertWithDefaults.
func RecommendationFromDiscovered(d Discovered, options ...library.RecommendationOption) library.Recommendation {
	return langx.Clone(library.Recommendation{
		Source:       md5x.String(library.RecommendationSourceDiscovered),
		ContentID:    d.ID,
		KnownMediaID: d.KnownMediaID,
		Mimetype:     mimex.Category(d.Mimetype),
		Language:     d.AudioDefaultLocale,
		Adult:        d.Adult,
		TombstoneAt:  time.Now().Add(library.RecommendationTTL),
		Title:        d.Title,
		Overview:     d.Description,
		Image:        d.PosterURI,
		Released:     d.ReleasedAt,
	}, options...)
}
