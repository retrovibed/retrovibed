package meta

import (
	"log"

	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
)

func NewDateRange(r timex.Range) *DateRange {
	return &DateRange{
		Oldest: grpcx.EncodeTime(timex.RFC3339NanoEncode(r.Start)),
		Newest: grpcx.EncodeTime(timex.RFC3339NanoEncode(r.End)),
	}
}

func TimexRange(r *DateRange, fallback timex.Range) timex.Range {
	if tr, err := timex.NewRangeISO8601(r.Oldest, r.Newest); err == nil {
		return tr
	} else {
		log.Println("invalid date range using fallback", err)
	}

	return fallback
}
