//go:build !darwin

package acoustics

import (
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// writeWAV synthesizes a mono 16-bit PCM WAV sine wave at the requested rate and
// duration, writes it to a temp file, and returns the path. FFmpeg decodes the
// canonical WAV layout natively, which lets these tests exercise the real
// decoder without committing a binary fixture.
func writeWAV(t *testing.T, sampleRate int, durationSec, freq float64) string {
	t.Helper()

	n := int(float64(sampleRate) * durationSec)
	dataSize := n * 2

	le := binary.LittleEndian
	b := make([]byte, 0, 44+dataSize)
	b = append(b, "RIFF"...)
	b = le.AppendUint32(b, uint32(36+dataSize))
	b = append(b, "WAVE"...)
	b = append(b, "fmt "...)
	b = le.AppendUint32(b, 16)                   // fmt chunk size
	b = le.AppendUint16(b, 1)                    // PCM
	b = le.AppendUint16(b, 1)                    // mono
	b = le.AppendUint32(b, uint32(sampleRate))   // sample rate
	b = le.AppendUint32(b, uint32(sampleRate*2)) // byte rate
	b = le.AppendUint16(b, 2)                    // block align
	b = le.AppendUint16(b, 16)                   // bits per sample
	b = append(b, "data"...)
	b = le.AppendUint32(b, uint32(dataSize))
	for i := range n {
		s := math.Sin(2.0 * math.Pi * freq * float64(i) / float64(sampleRate))
		b = le.AppendUint16(b, uint16(int16(s*32767)))
	}

	path := filepath.Join(t.TempDir(), "tone.wav")
	require.NoError(t, os.WriteFile(path, b, 0o600))
	return path
}

func TestProbeDuration(t *testing.T) {
	t.Run("returns the duration of a valid file", func(t *testing.T) {
		dur, err := ProbeDuration(writeWAV(t, 44100, 12.0, 440))
		require.NoError(t, err)
		require.True(t, dur >= 11.0 && dur <= 13.0, "duration %.3fs, expected ~12s", dur)
	})

	t.Run("errors on a missing file", func(t *testing.T) {
		_, err := ProbeDuration(filepath.Join(t.TempDir(), "missing.wav"))
		require.Error(t, err)
	})
}

func TestDecodePCM(t *testing.T) {
	t.Run("decodes every segment to mono float32 at SampleRate", func(t *testing.T) {
		segments := SegmentsForDuration(12.0)

		pcm, err := DecodePCM(t.Context(), writeWAV(t, 44100, 12.0, 440), segments)
		require.NoError(t, err)
		require.Len(t, pcm, NumSegments)

		expected := int(segments[0].DurationSec * SampleRate)
		for i, seg := range pcm {
			require.True(t, len(seg) >= expected/2 && len(seg) <= expected*2,
				"segment %d: %d samples, expected ~%d at %dHz", i, len(seg), expected, SampleRate)
			for _, v := range seg {
				// resampling overshoots the source's [-1, 1] slightly; this only
				// guards against garbage, not exact normalization.
				require.True(t, v >= -1.5 && v <= 1.5,
					"segment %d: sample %f far outside audio range", i, v)
			}
		}
	})

	t.Run("errors on a missing file", func(t *testing.T) {
		_, err := DecodePCM(t.Context(), filepath.Join(t.TempDir(), "missing.wav"), SegmentsForDuration(12.0))
		require.Error(t, err)
	})

	t.Run("errors on a non-media file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "notaudio.txt")
		require.NoError(t, os.WriteFile(path, []byte("this is not media"), 0o600))
		_, err := DecodePCM(t.Context(), path, SegmentsForDuration(12.0))
		require.Error(t, err)
	})

	t.Run("succeeds on concatenated MP3 with data inside the segment time range", func(t *testing.T) {
		// .fixtures/decode.corrupted.example.1.mp3 is two identical 12s MP3s
		// joined byte-for-byte. Generated with:
		//   ffmpeg -f lavfi -i "sine=frequency=440:duration=12" -ar 44100 -ac 1 tone.wav
		//   ffmpeg -i tone.wav -b:a 128k tone.mp3
		//   cat tone.mp3 tone.mp3 > decode.corrupted.example.1.mp3
		//
		// FFmpeg emits "invalid concatenated file detected" on open. The segment
		// ends well before the seam (6s < 12s) so decoding succeeds cleanly.
		seg := Segment{OffsetSec: 0, DurationSec: 6.0}
		buf, err := decodeSegment(t.Context(), ".fixtures/decode.corrupted.example.1.mp3", seg)
		require.NoError(t, err)
		require.NotEmpty(t, buf)
	})

	t.Run("succeeds on concatenated MP3 when corrupt region is beyond endTs", func(t *testing.T) {
		// Same fixture as below. DurationSec=12 means endTs=12s, which is right at
		// the seam; the frame callback sets done=true before Demux errors on the
		// null region, so the error is suppressed and we get valid audio back.
		seg := Segment{OffsetSec: 0.0, DurationSec: 12.0}
		buf, err := decodeSegment(t.Context(), ".fixtures/decode.corrupted.example.2.mp3", seg)
		require.NoError(t, err)
		require.NotEmpty(t, buf)
	})

	t.Run("fails on concatenated MP3 with corrupt region within segment", func(t *testing.T) {
		// .fixtures/decode.corrupted.example.2.mp3 is a 12s MP3 followed by an
		// equal-sized block of null bytes. Generated with:
		//   ffmpeg -f lavfi -i "sine=frequency=440:duration=12" -ar 44100 -ac 1 tone.wav
		//   ffmpeg -i tone.wav -b:a 128k tone.mp3
		//   SIZE=$(wc -c < tone.mp3)
		//   cp tone.mp3 decode.corrupted.example.2.mp3
		//   dd if=/dev/zero bs=$SIZE count=1 >> decode.corrupted.example.2.mp3
		//
		// FFmpeg estimates ~24s (bitrate × file size). endTs=16s extends into the
		// null region; Demux hits AVERROR_INVALIDDATA before the frame callback
		// ever fires, so done stays false and the error propagates. This documents
		// the desired behavior (return available data rather than fail); it will
		// fail until fixed.
		seg := Segment{OffsetSec: 0.0, DurationSec: 16.0}
		_, err := decodeSegment(t.Context(), ".fixtures/decode.corrupted.example.2.mp3", seg)
		require.ErrorContains(t, err, "Invalid data found when processing input")
	})
}

func TestAnalyzeFile(t *testing.T) {
	t.Run("produces a non-zero vector from a real file", func(t *testing.T) {
		fv, err := AnalyzeFile(t.Context(), writeWAV(t, 44100, 12.0, 440))
		require.NoError(t, err)

		nonZero := 0
		for _, v := range fv {
			if v != 0 {
				nonZero++
			}
		}
		require.Positive(t, nonZero, "expected a non-zero feature vector")
	})

	t.Run("rejects a track shorter than the minimum", func(t *testing.T) {
		_, err := AnalyzeFile(t.Context(), writeWAV(t, 44100, 5.0, 440))
		require.Error(t, err)
	})
}
