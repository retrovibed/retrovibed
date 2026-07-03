package ddiscapi

import (
	"github.com/retrovibed/retrovibed/shallows/internal/grpcx"
	"github.com/retrovibed/retrovibed/shallows/internal/langx"
	"github.com/retrovibed/retrovibed/shallows/internal/timex"
	"github.com/retrovibed/retrovibed/shallows/tracking"
)

func NewDiscoveryFromTrackingUnknownHash(mu tracking.UnknownHash) (_ *Discovery, err error) {
	var d Discovery
	mu = langx.Clone(mu, timex.JSONSafeEncodeOption, timex.UTCEncodeOption)

	if err = grpcx.JSONDecode(mu, &d); err != nil {
		return nil, err
	}

	return &d, nil
}
