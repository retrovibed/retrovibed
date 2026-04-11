package meta

import (
	"github.com/retrovibed/retrovibed/internal/grpcx"
	"github.com/retrovibed/retrovibed/internal/timex"
)

func NewDateRange(r timex.Range) *DateRange {
	return &DateRange{
		Oldest: grpcx.EncodeTime(timex.RFC3339NanoEncode(r.Start)),
		Newest: grpcx.EncodeTime(timex.RFC3339NanoEncode(r.End)),
	}
}
