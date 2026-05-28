package acoustics

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestThresholdValidation decodes MP3s via ffmpeg CLI, computes feature vectors,
// and evaluates precision/recall at various similarity thresholds.
// Ground truth: tracks in the same subdirectory are "similar."
//
// Set ACOUSTICS_TEST_DATA to the directory containing artist subdirectories.
// Skips if unset or ffmpeg is unavailable.
func TestThresholdValidation(t *testing.T) {
	dataDir := os.Getenv("ACOUSTICS_TEST_DATA")
	if dataDir == "" {
		t.Skip("ACOUSTICS_TEST_DATA not set")
	}

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not in PATH")
	}

	type track struct {
		path   string
		artist string
		vec    FeatureVector
	}

	var tracks []track
	artists, err := os.ReadDir(dataDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, artist := range artists {
		if !artist.IsDir() {
			continue
		}

		files, err := filepath.Glob(filepath.Join(dataDir, artist.Name(), "*.mp3"))
		if err != nil {
			t.Fatal(err)
		}

		for _, f := range files {
			pcm, err := decodeWithFFmpeg(f)
			if err != nil {
				t.Logf("skip %s: %v", f, err)
				continue
			}

			dur := float64(len(pcm)) / float64(SampleRate)
			segs := SegmentsForDuration(dur)
			if segs == nil {
				t.Logf("skip %s: too short (%.1fs)", f, dur)
				continue
			}

			segments := splitPCM(pcm, segs)
			vec := AnalyzeSamples(segments)
			tracks = append(tracks, track{
				path:   filepath.Base(f),
				artist: artist.Name(),
				vec:    vec,
			})
		}
	}

	if len(tracks) < 10 {
		t.Fatalf("only %d tracks decoded, need at least 10", len(tracks))
	}
	artistSet := make(map[string]struct{})
	for _, tr := range tracks {
		artistSet[tr.artist] = struct{}{}
	}
	t.Logf("analyzed %d tracks across %d artists", len(tracks), len(artistSet))

	// Normalize all vectors.
	var stats RunningStats
	for _, tr := range tracks {
		stats.Update(tr.vec)
	}

	normalized := make([]FeatureVector, len(tracks))
	for i, tr := range tracks {
		normalized[i] = stats.Normalize(tr.vec)
	}

	// Compute all pairwise similarities.
	type pair struct {
		sim      float64
		sameArt  bool
	}

	var pairs []pair
	for i := 0; i < len(tracks); i++ {
		for j := i + 1; j < len(tracks); j++ {
			pairs = append(pairs, pair{
				sim:     CosineSimilarity(normalized[i], normalized[j]),
				sameArt: tracks[i].artist == tracks[j].artist,
			})
		}
	}

	sameCount := 0
	for _, p := range pairs {
		if p.sameArt {
			sameCount++
		}
	}
	t.Logf("%d pairs total: %d same-artist, %d cross-artist", len(pairs), sameCount, len(pairs)-sameCount)

	// Evaluate thresholds.
	t.Logf("")
	t.Logf("threshold | precision | recall | f1")
	t.Logf("----------|-----------|--------|-------")
	for _, thresh := range []float64{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8} {
		tp, fp, fn := 0, 0, 0
		for _, p := range pairs {
			above := p.sim >= thresh
			if p.sameArt && above {
				tp++
			} else if !p.sameArt && above {
				fp++
			} else if p.sameArt && !above {
				fn++
			}
		}

		precision := 0.0
		if tp+fp > 0 {
			precision = float64(tp) / float64(tp+fp)
		}
		recall := 0.0
		if tp+fn > 0 {
			recall = float64(tp) / float64(tp+fn)
		}
		f1 := 0.0
		if precision+recall > 0 {
			f1 = 2 * precision * recall / (precision + recall)
		}

		t.Logf("  %.1f     |   %.3f   | %.3f  | %.3f", thresh, precision, recall, f1)
	}
}

func decodeWithFFmpeg(path string) ([]float32, error) {
	cmd := exec.Command("ffmpeg",
		"-i", path,
		"-f", "f32le",
		"-ac", "1",
		"-ar", fmt.Sprintf("%d", SampleRate),
		"pipe:1",
	)
	cmd.Stderr = nil

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	nSamples := len(out) / 4
	samples := make([]float32, nSamples)
	for i := range nSamples {
		samples[i] = math.Float32frombits(binary.LittleEndian.Uint32(out[i*4:]))
	}
	return samples, nil
}

func splitPCM(pcm []float32, segs []Segment) [][]float32 {
	result := make([][]float32, len(segs))
	for i, seg := range segs {
		start := int(seg.OffsetSec * SampleRate)
		end := start + int(seg.DurationSec*SampleRate)
		if start >= len(pcm) {
			continue
		}
		end = min(end, len(pcm))
		result[i] = pcm[start:end]
	}
	return result
}

