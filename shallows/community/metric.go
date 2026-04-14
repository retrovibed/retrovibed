package community

import (
	"time"

	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

// MetricPeriodOption allows configuring the metric query period.
type MetricPeriodOption func(*timex.Range)

// MetricPeriodOptionStartDate sets the start date for the metric query.
func MetricPeriodOptionStartDate(t time.Time) MetricPeriodOption {
	return func(p *timex.Range) {
		p.Start = t
	}
}

// MetricPeriodOptionEndDate sets the end date for the metric query.
func MetricPeriodOptionEndDate(t time.Time) MetricPeriodOption {
	return func(p *timex.Range) {
		p.End = t
	}
}

// ResolvePeriod resolves the metric period options into start and end times.
// If no options are provided, defaults to last 30 days.
func ResolvePeriod(options ...MetricPeriodOption) (start, end time.Time) {
	p := &timex.Range{}
	for _, opt := range options {
		opt(p)
	}

	p.End = langx.FirstNonZero(p.End, time.Now())
	p.Start = langx.FirstNonZero(p.Start, p.End.AddDate(0, 0, -30))

	return p.Start, p.End
}
