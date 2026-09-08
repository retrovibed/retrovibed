package community

import (
	"github.com/gofrs/uuid/v5"
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
)

type CommunityPublisherOption func(*CommunityPublisher)

// MetricPeriodOptionStartDate sets the start date for the metric query.
func CommunityPublisherOptionTestDefaults(c *CommunityPublisher) {
	c.ID = uuid.Must(uuid.NewV7()).String()
	c.CommunityID = uuidx.WithSuffix(1)
	c.PublisherID = uuidx.WithSuffix(2)
}
