//go:build !(darwin || android)

package acoustics

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/mutablelogic/go-media/frame"
	"github.com/mutablelogic/go-media/reader"
	ff "github.com/mutablelogic/go-media/sys/ffmpeg80"
	"github.com/mutablelogic/go-media/writer"
	"github.com/retrovibed/retrovibed/retroapi/errorsx"
)

// ProbeDuration opens the file with FFmpeg and returns its duration.
func ProbeDuration(path string) (time.Duration, error) {
	r, err := reader.Open(path)
	if err != nil {
		return 0, err
	}
	defer r.Close()
	return r.Duration(), nil
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
			return nil, err
		}
		results[i] = buf
	}
	return results, nil
}

func decodeSegment(ctx context.Context, path string, seg Segment) ([]float32, error) {
	var (
		buf = make([]float32, 0, 330_750)
	)
	r, err := reader.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	best := -1
	for index, stream := range r.Streams() {
		if stream.Par().CodecType() == ff.AVMEDIA_TYPE_AUDIO && (best < 0 || index < best) {
			best = index
		}
	}
	if best < 0 {
		return nil, errorsx.String("no audio stream")
	}

	// mono float32 at SampleRate, in frames of any size.
	par := ff.AVCodec_parameters_alloc()
	if par == nil {
		return nil, errorsx.String("failed to allocate codec parameters")
	}
	defer ff.AVCodec_parameters_free(par)
	par.SetCodecType(ff.AVMEDIA_TYPE_AUDIO)
	par.SetSampleFormat(ff.AV_SAMPLE_FMT_FLT)
	par.SetSampleRate(SampleRate)
	var layout ff.AVChannelLayout
	ff.AVUtil_channel_layout_default(&layout, 1)
	if err = par.SetChannelLayout(layout); err != nil {
		return nil, err
	}
	resampler, err := writer.NewResampler(par, 0, ff.AVRational{})
	if err != nil {
		return nil, err
	}
	defer resampler.Close()

	if err = r.Seek(best, seg.Offset); err != nil {
		return nil, err
	}

	endTs := seg.Offset + seg.Duration

	var done bool
	collect := func(f frame.Frame) error {
		buf = append(buf, f.(*frame.AudioFrame).Float32(0)...)
		return nil
	}

	decoder, err := reader.NewDecoder(func(f frame.Frame) error {
		af := f.(*frame.AudioFrame)

		if time.Duration(ff.AVUtil_rational_rescale_q(af.Pts(), af.TimeBase(), ff.AVUtil_rational(1, int(time.Second)))) >= endTs {
			done = true
			return io.EOF
		}

		return resampler.Process(f, collect)
	})
	if err != nil {
		return nil, err
	}
	if err = decoder.Add(best); err != nil {
		return nil, err
	}

	if err = r.Decode(ctx, nil, decoder); err != nil && !(done && errors.Is(err, io.EOF)) {
		return nil, err
	}

	if !done {
		if err = resampler.Process(nil, collect); err != nil {
			return nil, err
		}
	}

	return buf, nil
}
