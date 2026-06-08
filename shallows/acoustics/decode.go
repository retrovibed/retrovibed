//go:build !(darwin || android)

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

// DecodePCM decodes each segment of an audio file to mono float32 at SampleRate.
// The file is reopened per segment: go-media flushes the decoder at the end of
// every decode and exposes no way to reset it, so a reader cannot be reused
// across seeks.
func DecodePCM(ctx context.Context, path string, segments []Segment) ([][]float32, error) {
	results := make([][]float32, len(segments))
	for i, seg := range segments {
		buf, err := decodeSegment(ctx, path, seg)
		if err != nil {
			returnPCMBuffers(results[:i])
			return nil, err
		}
		results[i] = buf
	}
	return results, nil
}

func decodeSegment(ctx context.Context, path string, seg Segment) ([]float32, error) {
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

	if err = reader.Seek(best, seg.OffsetSec); err != nil {
		return nil, err
	}

	buf := pcmPool.Get().([]float32)[:0]
	endTs := seg.OffsetSec + seg.DurationSec

	err = reader.Demux(ctx, func(stream int, p *ffmpeg.Par) (*ffmpeg.Par, error) {
		if stream == best {
			return par, nil
		}
		return nil, nil
	}, func(_ int, frame *ffmpeg.Frame) error {
		if frame.Type() != media.AUDIO {
			return nil
		}
		if ts := frame.Ts(); ts != ffmpeg.TS_UNDEFINED && ts >= endTs {
			return io.EOF
		}
		buf = append(buf, frame.Float32(0)...)
		return nil
	}, nil)

	if err != nil {
		pcmPool.Put(buf[:0]) // nolint: staticcheck
		return nil, err
	}

	return buf, nil
}
