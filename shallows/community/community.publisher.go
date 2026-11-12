package community

import (
	"github.com/retrovibed/retrovibed/retroapi/uuidx"
)

type CommunityPublisherOption func(*CommunityPublisher)

// MetricPeriodOptionStartDate sets the start date for the metric query.
func CommunityPublisherOptionTestDefaults(c *CommunityPublisher) {
	c.ID = uuidx.NewV7().String()
	c.CommunityID = uuidx.WithSuffix(1)
	c.PublisherID = uuidx.WithSuffix(2)
}
