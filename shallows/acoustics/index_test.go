package acoustics

import (
	"math"
	"testing"
	"time"
)

func TestSegmentsForDuration(t *testing.T) {
	t.Run("too short", func(t *testing.T) {
		if segs := SegmentsForDuration(9 * time.Second); segs != nil {
			t.Fatalf("expected nil for 9s track, got %d segments", len(segs))
		}
	})

	t.Run("short track splits evenly", func(t *testing.T) {
		segs := SegmentsForDuration(15 * time.Second)
		if len(segs) != 3 {
			t.Fatalf("expected 3 segments, got %d", len(segs))
		}
		for i, s := range segs {
			expectedOffset := time.Duration(i) * 5 * time.Second
			if s.Offset != expectedOffset {
				t.Fatalf("segment %d: offset %v, expected %v", i, s.Offset, expectedOffset)
			}
			if s.Duration != 5*time.Second {
				t.Fatalf("segment %d: duration %v, expected 5s", i, s.Duration)
			}
		}
	})

	t.Run("long track caps at 30s per segment", func(t *testing.T) {
		segs := SegmentsForDuration(300 * time.Second)
		if len(segs) != 3 {
			t.Fatalf("expected 3 segments, got %d", len(segs))
		}
		for i, s := range segs {
			if s.Duration != MaxSegmentDuration {
				t.Fatalf("segment %d: duration %v, expected %v", i, s.Duration, MaxSegmentDuration)
			}
		}
		// Offsets should be 0, 100, 200
		if segs[1].Offset != 100*time.Second {
			t.Fatalf("segment 1 offset: %v, expected 100s", segs[1].Offset)
		}
	})

	t.Run("boundary at 10s", func(t *testing.T) {
		if segs := SegmentsForDuration(10 * time.Second); segs == nil {
			t.Fatal("expected segments for exactly 10s track")
		}
	})
}

func TestAnalyzeSamples(t *testing.T) {
	sine := func(hz float64, dur float64) []float32 {
		n := int(dur * SampleRate)
		s := make([]float32, n)
		for i := range s {
			s[i] = float32(math.Sin(2.0 * math.Pi * hz * float64(i) / SampleRate))
		}
		return s
	}

	t.Run("three segments produce non-zero vector", func(t *testing.T) {
		segments := [][]float32{
			sine(440, 5.0),
			sine(440, 5.0),
			sine(440, 5.0),
		}

		fv := AnalyzeSamples(segments)

		nonZero := 0
		for _, v := range fv {
			if v != 0 {
				nonZero++
			}
		}
		if nonZero < FeatureDim/2 {
			t.Fatalf("expected many non-zero dims, got %d/%d", nonZero, VectorDim)
		}
	})

	t.Run("different segments produce non-zero sigma", func(t *testing.T) {
		segments := [][]float32{
			sine(200, 5.0),
			sine(1000, 5.0),
			sine(4000, 5.0),
		}

		fv := AnalyzeSamples(segments)

		// Inter-window σ dimensions (FeatureDim through VectorDim-1)
		sigmaSum := float32(0.0)
		for i := FeatureDim; i < VectorDim; i++ {
			sigmaSum += fv[i]
		}
		if sigmaSum == 0 {
			t.Fatal("inter-window sigma should be non-zero for different frequency segments")
		}
	})

	t.Run("identical segments produce zero sigma", func(t *testing.T) {
		seg := sine(440, 5.0)
		segments := [][]float32{seg, seg, seg}

		fv := AnalyzeSamples(segments)

		for i := FeatureDim; i < VectorDim; i++ {
			if fv[i] != 0 {
				t.Fatalf("dim %d: expected 0 sigma for identical segments, got %f", i, fv[i])
			}
		}
	})

	t.Run("similar tones produce similar vectors", func(t *testing.T) {
		a := [][]float32{sine(440, 5.0), sine(440, 5.0), sine(440, 5.0)}
		b := [][]float32{sine(445, 5.0), sine(445, 5.0), sine(445, 5.0)}
		c := [][]float32{sine(3000, 5.0), sine(3000, 5.0), sine(3000, 5.0)}

		va := AnalyzeSamples(a)
		vb := AnalyzeSamples(b)
		vc := AnalyzeSamples(c)

		simAB := CosineSimilarity(va, vb)
		simAC := CosineSimilarity(va, vc)

		if simAB <= simAC {
			t.Fatalf("440 Hz should be more similar to 445 Hz (%.4f) than to 3000 Hz (%.4f)", simAB, simAC)
		}
	})
}
