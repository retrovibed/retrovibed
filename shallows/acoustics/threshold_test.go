package acoustics

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// TestThresholdValidation decodes audio files via the package's FFmpeg-backed
// decoder, computes feature vectors, and evaluates precision/recall at various
// similarity thresholds. Ground truth: tracks in the same subdirectory are
// "similar." Set ACOUSTICS_TEST_DATA to the directory containing artist
// subdirectories. Skips if unset.
func TestThresholdValidation(t *testing.T) {
	dataDir := os.Getenv("ACOUSTICS_TEST_DATA")
	if dataDir == "" {
		t.Skip("ACOUSTICS_TEST_DATA not set")
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
			vec, err := AnalyzeFile(context.Background(), f)
			if err != nil {
				t.Logf("skip %s: %v", f, err)
				continue
			}

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

	var stats RunningStats
	for _, tr := range tracks {
		stats.Update(tr.vec)
	}

	normalized := make([]FeatureVector, len(tracks))
	for i, tr := range tracks {
		normalized[i] = stats.Normalize(tr.vec)
	}

	type pair struct {
		sim     float64
		sameArt bool
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
