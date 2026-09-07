package communityapi

import (
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/stringsx"
)

func NewCommunity(opts ...func(*Community)) *Community {
	var c Community
	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

// CommunityFromDeeppool converts a deeppool Community proto to the local DB model.
// Sync-only fields (AutoDownload, LastSyncAt, SyncFeedAt, SyncCursorPublishedContent)
// are left as zero values — they come from local state, not deeppool.
func CommunityFromDeeppool(c *Community) community.Community {
	return community.Community{
		ID:                 c.Id,
		AccountID:          stringsx.DefaultIfBlank(c.AccountId, uuid.Nil.String()),
		CreatedAt:          errorsx.Zero(grpcx.DecodeTime(c.CreatedAt)),
		UpdatedAt:          errorsx.Zero(grpcx.DecodeTime(c.UpdatedAt)),
		Mimetype:           c.Mimetype,
		Description:        c.Description,
		Entropy:            c.Entropy,
		Bytes:              int64(c.Bytes),
		DefaultPublishMode: int32(c.DefaultPublishMode),
		Hidden:             c.Hidden,
		URL:                c.Url,
		Adult:              c.Adult,
		DefaultTTL:         int64(c.DefaultTtl),
		DefaultLanguage:    c.DefaultLanguage,
	}
}

// CommunityOptionFromDB converts a local DB community to proto options.
// Call site should apply timex.JSONSafeEncodeOption before passing the DB value.
func CommunityOptionFromDB(c community.Community) func(*Community) {
	return func(p *Community) {
		p.Id = c.ID
		p.AccountId = c.AccountID
		p.CreatedAt = grpcx.EncodeTime(c.CreatedAt)
		p.UpdatedAt = grpcx.EncodeTime(c.UpdatedAt)
		p.SubscribedAt = grpcx.EncodeTime(c.SubscribedAt)
		p.Mimetype = c.Mimetype
		p.Description = c.Description
		p.Entropy = c.Entropy
		p.Bytes = uint64(c.Bytes)
		p.DefaultPublishMode = PublishMode(c.DefaultPublishMode)
		p.Hidden = c.Hidden
		p.Url = c.URL
		p.Adult = c.Adult
		p.DefaultTtl = uint64(c.DefaultTTL)
		p.DefaultLanguage = c.DefaultLanguage
		p.LastSyncAt = grpcx.EncodeTime(c.LastSyncAt)
	}
}

// CommunityMetricOptionFromDB converts a database model to proto options.
func CommunityMetricOptionFromDB(cm community.CommunityMetric) func(*CommunityMetric) {
	return func(m *CommunityMetric) {
		m.Id = cm.ID
		m.CommunityId = cm.CommunityID
		m.PeriodStart = grpcx.EncodeTime(cm.PeriodStart)
		m.PeriodEnd = grpcx.EncodeTime(cm.PeriodEnd)
		m.Subscribers = cm.Subscribers
	}
}
