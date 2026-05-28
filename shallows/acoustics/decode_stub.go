//go:build !darwin && !ffmpeg_enabled

package acoustics

import (
	"context"
	"errors"
)

func DecodePCM(_ context.Context, _ string, _ []Segment) ([][]float32, error) {
	return nil, errors.ErrUnsupported
}

func ProbeDuration(_ string) (float64, error) {
	return 0, errors.ErrUnsupported
}
