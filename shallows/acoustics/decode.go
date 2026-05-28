//go:build !darwin && ffmpeg_enabled

package acoustics

import (
	"context"
	"io"

	media "github.com/mutablelogic/go-media"
	"github.com/mutablelogic/go-media/pkg/ffmpeg"
	"github.com/retrovibed/retrovibed/shallows/internal/errorsx"
)

// ProbeDuration opens the file with FFmpeg and returns its duration in seconds.
func ProbeDuration(path string) (float64, error) {
	reader, err := ffmpeg.Open(path)
	if err != nil {
		return 0, err
	}
	defer reader.Close()
	return reader.Duration().Seconds(), nil
}

// DecodePCM opens an audio file and decodes each segment to mono float32 at
// SampleRate Hz. Uses a single ffmpeg.Open + reader.Map, seeking between segments.
func DecodePCM(ctx context.Context, path string, segments []Segment) ([][]float32, error) {
	reader, err := ffmpeg.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	best := reader.BestStream(media.AUDIO)
	if best < 0 {
		return nil, errorsx.String("no audio stream")
	}

	par, err := ffmpeg.NewAudioPar("flt", "mono", SampleRate)
	if err != nil {
		return nil, err
	}

	decoders, err := reader.Map(func(stream int, p *ffmpeg.Par) (*ffmpeg.Par, error) {
		if stream == best {
			return par, nil
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	defer decoders.Close()

	var (
		buf    []float32
		endTs  float64
	)

	results := make([][]float32, len(segments))
	for i, seg := range segments {
		err = reader.Seek(best, seg.OffsetSec)
		if err != nil {
			return nil, err
		}

		buf = pcmPool.Get().([]float32)[:0]
		endTs = seg.OffsetSec + seg.DurationSec

		err = reader.DecodeWithContext(ctx, decoders, func(_ int, frame *ffmpeg.Frame) error {
			if frame.Type() != media.AUDIO {
				return nil
			}
			if ts := frame.Ts(); ts != ffmpeg.TS_UNDEFINED && ts >= endTs {
				return io.EOF
			}
			buf = append(buf, frame.Float32(0)...)
			return nil
		})

		if errorsx.Ignore(err, io.EOF) != nil {
			pcmPool.Put(buf[:0])
			return nil, err
		}

		results[i] = buf
	}

	return results, nil
}
