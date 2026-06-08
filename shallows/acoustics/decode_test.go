//go:build !darwin

package acoustics

import (
	"context"
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"testing"
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
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatalf("write wav: %v", err)
	}
	return path
}

func TestProbeDuration(t *testing.T) {
	t.Run("returns the duration of a valid file", func(t *testing.T) {
		dur, err := ProbeDuration(writeWAV(t, 44100, 12.0, 440))
		if err != nil {
			t.Fatalf("ProbeDuration: %v", err)
		}
		if dur < 11.0 || dur > 13.0 {
			t.Fatalf("duration %.3fs, expected ~12s", dur)
		}
	})

	t.Run("errors on a missing file", func(t *testing.T) {
		if _, err := ProbeDuration(filepath.Join(t.TempDir(), "missing.wav")); err == nil {
			t.Fatal("expected error for missing file")
		}
	})
}

func TestDecodePCM(t *testing.T) {
	t.Run("decodes every segment to mono float32 at SampleRate", func(t *testing.T) {
		segments := SegmentsForDuration(12.0)

		pcm, err := DecodePCM(context.Background(), writeWAV(t, 44100, 12.0, 440), segments)
		if err != nil {
			t.Fatalf("DecodePCM: %v", err)
		}
		defer returnPCMBuffers(pcm)

		if len(pcm) != NumSegments {
			t.Fatalf("got %d segments, expected %d", len(pcm), NumSegments)
		}

		expected := int(segments[0].DurationSec * SampleRate)
		for i, seg := range pcm {
			if len(seg) < expected/2 || len(seg) > expected*2 {
				t.Fatalf("segment %d: %d samples, expected ~%d at %dHz", i, len(seg), expected, SampleRate)
			}
			for _, v := range seg {
				// resampling overshoots the source's [-1, 1] slightly; this only
				// guards against garbage, not exact normalization.
				if v < -1.5 || v > 1.5 {
					t.Fatalf("segment %d: sample %f far outside audio range", i, v)
				}
			}
		}
	})

	t.Run("errors on a missing file", func(t *testing.T) {
		_, err := DecodePCM(context.Background(), filepath.Join(t.TempDir(), "missing.wav"), SegmentsForDuration(12.0))
		if err == nil {
			t.Fatal("expected error for missing file")
		}
	})

	t.Run("errors on a non-media file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "notaudio.txt")
		if err := os.WriteFile(path, []byte("this is not media"), 0o600); err != nil {
			t.Fatalf("write file: %v", err)
		}
		if _, err := DecodePCM(context.Background(), path, SegmentsForDuration(12.0)); err == nil {
			t.Fatal("expected error for non-media file")
		}
	})
}

func TestAnalyzeFile(t *testing.T) {
	t.Run("produces a non-zero vector from a real file", func(t *testing.T) {
		fv, err := AnalyzeFile(context.Background(), writeWAV(t, 44100, 12.0, 440))
		if err != nil {
			t.Fatalf("AnalyzeFile: %v", err)
		}

		nonZero := 0
		for _, v := range fv {
			if v != 0 {
				nonZero++
			}
		}
		if nonZero == 0 {
			t.Fatal("expected a non-zero feature vector")
		}
	})

	t.Run("rejects a track shorter than the minimum", func(t *testing.T) {
		if _, err := AnalyzeFile(context.Background(), writeWAV(t, 44100, 5.0, 440)); err == nil {
			t.Fatal("expected error for sub-minimum track")
		}
	})
}
