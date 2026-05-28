package acoustics

import (
	"math"
	"testing"
)

func TestAggregateWindows(t *testing.T) {
	t.Run("mean and sigma", func(t *testing.T) {
		var w [3]WindowFeatures
		w[0][0] = 1.0
		w[1][0] = 2.0
		w[2][0] = 3.0

		fv := AggregateWindows(w)

		expectedMean := 2.0
		if math.Abs(fv[0]-expectedMean) > 1e-10 {
			t.Fatalf("mean: expected %f, got %f", expectedMean, fv[0])
		}

		// σ = sqrt(((1-2)^2 + (2-2)^2 + (3-2)^2) / 3) = sqrt(2/3)
		expectedSigma := math.Sqrt(2.0 / 3.0)
		if math.Abs(fv[FeatureDim]-expectedSigma) > 1e-10 {
			t.Fatalf("sigma: expected %f, got %f", expectedSigma, fv[FeatureDim])
		}
	})

	t.Run("identical windows produce zero sigma", func(t *testing.T) {
		var w [3]WindowFeatures
		for i := range FeatureDim {
			w[0][i] = float32(i)
			w[1][i] = float32(i)
			w[2][i] = float32(i)
		}

		fv := AggregateWindows(w)

		for d := range FeatureDim {
			if fv[FeatureDim+d] != 0 {
				t.Fatalf("sigma[%d]: expected 0, got %f", d, fv[FeatureDim+d])
			}
		}
	})

	t.Run("all dimensions populated", func(t *testing.T) {
		var w [3]WindowFeatures
		for s := range 3 {
			for d := range FeatureDim {
				w[s][d] = float32(s*10 + d)
			}
		}

		fv := AggregateWindows(w)

		for i, v := range fv {
			if v == 0 && i < FeatureDim {
				t.Fatalf("dimension %d: unexpected zero mean", i)
			}
		}
	})
}

func TestRunningStats(t *testing.T) {
	t.Run("normalize zero-mean unit-variance", func(t *testing.T) {
		var stats RunningStats

		// Feed vectors with known distribution in dim 0: [1, 3]
		var a, b FeatureVector
		a[0] = 1.0
		b[0] = 3.0
		stats.Update(a)
		stats.Update(b)

		// Mean = 2, variance = 1, σ = 1
		// Normalize a: (1-2)/1 = -1
		// Normalize b: (3-2)/1 = 1
		na := stats.Normalize(a)
		nb := stats.Normalize(b)

		if math.Abs(na[0]-(-1.0)) > 1e-10 {
			t.Fatalf("normalized a[0]: expected -1, got %f", na[0])
		}
		if math.Abs(nb[0]-1.0) > 1e-10 {
			t.Fatalf("normalized b[0]: expected 1, got %f", nb[0])
		}
	})

	t.Run("insufficient count returns zero", func(t *testing.T) {
		var stats RunningStats
		var fv FeatureVector
		fv[0] = 5.0
		stats.Update(fv)

		norm := stats.Normalize(fv)
		for i, v := range norm {
			if v != 0 {
				t.Fatalf("dim %d: expected 0 with count=1, got %f", i, v)
			}
		}
	})

	t.Run("zero variance dimension stays zero", func(t *testing.T) {
		var stats RunningStats
		var a, b FeatureVector
		a[0] = 5.0
		b[0] = 5.0 // same value; variance = 0
		stats.Update(a)
		stats.Update(b)

		norm := stats.Normalize(a)
		if norm[0] != 0 {
			t.Fatalf("dim 0: expected 0 for zero variance, got %f", norm[0])
		}
	})

	t.Run("recompute matches incremental", func(t *testing.T) {
		vectors := make([]FeatureVector, 10)
		for i := range vectors {
			vectors[i][0] = float64(i)
			vectors[i][5] = float64(i * 2)
		}

		var incremental RunningStats
		for _, v := range vectors {
			incremental.Update(v)
		}

		var recomputed RunningStats
		recomputed.Recompute(vectors)

		if incremental.Count != recomputed.Count {
			t.Fatalf("count mismatch: %d vs %d", incremental.Count, recomputed.Count)
		}
		for d := range VectorDim {
			if math.Abs(incremental.Sum[d]-recomputed.Sum[d]) > 1e-10 {
				t.Fatalf("sum[%d] mismatch", d)
			}
			if math.Abs(incremental.SumSq[d]-recomputed.SumSq[d]) > 1e-10 {
				t.Fatalf("sumSq[%d] mismatch", d)
			}
		}
	})
}
