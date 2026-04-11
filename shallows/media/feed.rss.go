package media

import (
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/internal/errorsx"
	"github.com/retrovibed/retrovibed/internal/grpcx"
	"github.com/retrovibed/retrovibed/internal/langx"
	"github.com/retrovibed/retrovibed/internal/timex"
	"github.com/retrovibed/retrovibed/rss"
	"github.com/retrovibed/retrovibed/tracking"
)

func NewTrackingFeedRSSFromFeedRSS(req *rss.FeedCreateRequest) func(*tracking.RSS) {
	return func(t *tracking.RSS) {
		t.Description = req.Feed.Description
		t.URL = req.Feed.Url
		t.Autodownload = req.Feed.Autodownload
		t.Autoarchive = req.Feed.Autoarchive
		t.Contributing = req.Feed.Contributing
		t.EncryptionSeed = req.Feed.EncryptionSeed
		t.Digest = langx.FirstNonZero(req.Feed.Digest, uuid.Nil.String())
		t.NextCheck = timex.FirstNonZero(
			errorsx.Zero(grpcx.DecodeTime(req.Feed.NextCheck)),
			t.NextCheck,
		)
	}
}
