//go:build darwin || android

package acoustics

import (
	"context"
	"errors"
	"time"
)

func DecodePCM(_ context.Context, _ string, _ []Segment) ([][]float32, error) {
	return nil, errors.ErrUnsupported
}

func ProbeDuration(_ string) (time.Duration, error) {
	return 0, errors.ErrUnsupported
}
