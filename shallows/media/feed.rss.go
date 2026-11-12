package media

import (
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/rss"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

func NewTrackingFeedRSSFromFeedRSS(req *rss.FeedCreateRequest) func(*tracking.RSS) {
	return func(t *tracking.RSS) {
		t.Description = req.Feed.Description
		t.URL = req.Feed.Url
		t.Autodownload = req.Feed.Autodownload
		t.Autoarchive = req.Feed.Autoarchive
		t.Contributing = req.Feed.Contributing
		t.EncryptionSeed = req.Feed.EncryptionSeed
		t.Digest = langx.FirstNonZero(req.Feed.Digest, uuidx.Nil().String())
		t.NextCheck = timex.FirstNonZero(
			errorsx.Zero(grpcx.DecodeTime(req.Feed.NextCheck)),
			t.NextCheck,
		)
	}
}
