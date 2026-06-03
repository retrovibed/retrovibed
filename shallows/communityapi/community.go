package communityapi

import (
	"github.com/retrovibed/retrovibed/shallows/community"
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"google.golang.org/protobuf/encoding/protojson"
)

func (t *PublishContentRequest) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(t)
}

func (t *PublishContentRequest) UnmarshalJSON(b []byte) error {
	return protojson.Unmarshal(b, t)
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
