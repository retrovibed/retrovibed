package acoustics

import (
	"context"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/errorsx"
)

const (
	MinTrackDuration   = 10 * time.Second // shorter tracks excluded
	MaxSegmentDuration = 30 * time.Second // cap per segment
	NumSegments        = 3
	StatsVersion       = 1
)

type Segment struct {
	Offset   time.Duration
	Duration time.Duration
}

// SegmentsForDuration splits a track into three equal parts, each capped at
// MaxSegmentDuration. Returns nil if the track is shorter than MinTrackDuration.
func SegmentsForDuration(duration time.Duration) []Segment {
	if duration < MinTrackDuration {
		return nil
	}

	segDur := min(duration/NumSegments, MaxSegmentDuration)
	segments := make([]Segment, NumSegments)
	for i := range NumSegments {
		segments[i] = Segment{
			Offset:   time.Duration(i) * duration / NumSegments,
			Duration: segDur,
		}
	}
	return segments
}

// AnalyzeSamples computes a FeatureVector from PCM segments (one per window).
// Each segment is mono float32 at SampleRate. Caller is responsible for decoding.
func AnalyzeSamples(segments [][]float32) FeatureVector {
	var windows [3]WindowFeatures
	for i, samples := range segments {
		if i >= 3 {
			break
		}
		mag := STFT(samples)
		if mag == nil {
			continue
		}
		windows[i] = ExtractWindowFeatures(mag, SpectrogramFrames(len(samples)))
	}
	return AggregateWindows(windows)
}

// AnalyzeFile decodes an audio file and computes its feature vector.
// Returns the raw (un-normalized) vector. Caller normalizes via RunningStats.
func AnalyzeFile(ctx context.Context, path string) (FeatureVector, error) {
	dur, err := ProbeDuration(path)
	if err != nil {
		return FeatureVector{}, errorsx.Wrap(err, "duration detection failed")
	}

	segments := SegmentsForDuration(dur)
	if segments == nil {
		return FeatureVector{}, errorsx.String("track too short")
	}

	pcm, err := DecodePCM(ctx, path, segments)
	if err != nil {
		return FeatureVector{}, errorsx.Wrap(err, "decode pcm failed")
	}

	return AnalyzeSamples(pcm), nil
}
