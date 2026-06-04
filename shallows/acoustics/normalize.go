package acoustics

import "math"

// FeatureVector is the final 128-dimensional acoustic fingerprint.
// float32 matches the DSP pipeline and the FLOAT[128] storage column;
// running stats and similarity computations widen to float64 inside their
// accumulators where many small values sum.
type FeatureVector [VectorDim]float32

// AggregateWindows computes the mean and inter-window standard deviation
// of three per-window feature vectors, producing the final 128-dim vector.
func AggregateWindows(windows [3]WindowFeatures) FeatureVector {
	var fv FeatureVector

	for d := range FeatureDim {
		a, b, c := float64(windows[0][d]), float64(windows[1][d]), float64(windows[2][d])
		mean := (a + b + c) / 3.0
		fv[d] = float32(mean)

		v := ((a-mean)*(a-mean) + (b-mean)*(b-mean) + (c-mean)*(c-mean)) / 3.0
		fv[FeatureDim+d] = float32(math.Sqrt(v))
	}

	return fv
}

// RunningStats tracks incremental count, sum, and sum-of-squares per dimension
// for z-score normalization. Accumulators are float64 to remain stable across
// large catalogs where sum-of-squares can grow into the 10^9 range.
type RunningStats struct {
	Count int64
	Sum   [VectorDim]float64
	SumSq [VectorDim]float64
}

// Update adds a feature vector to the running statistics.
func (s *RunningStats) Update(fv FeatureVector) {
	s.Count++
	for i, v := range fv {
		f := float64(v)
		s.Sum[i] += f
		s.SumSq[i] += f * f
	}
}

// Normalize applies z-score normalization to fv using the current statistics.
// Returns the normalized vector. If count < 2 or σ is near zero for a dimension,
// that dimension is set to 0.
func (s *RunningStats) Normalize(fv FeatureVector) FeatureVector {
	var out FeatureVector
	if s.Count < 2 {
		return out
	}

	n := float64(s.Count)
	for i := range VectorDim {
		mean := s.Sum[i] / n
		variance := s.SumSq[i]/n - mean*mean
		if variance < 1e-12 {
			continue
		}
		out[i] = float32((float64(fv[i]) - mean) / math.Sqrt(variance))
	}

	return out
}

// Recompute rebuilds running statistics from a batch of feature vectors.
func (s *RunningStats) Recompute(vectors []FeatureVector) {
	*s = RunningStats{}
	for _, fv := range vectors {
		s.Update(fv)
	}
}

// CosineSimilarity returns the cosine of the angle between two vectors.
// Accumulators widen to float64 to keep the dot product stable across 128 terms.
func CosineSimilarity(a, b FeatureVector) float64 {
	dot, normA, normB := 0.0, 0.0, 0.0
	for i := range VectorDim {
		x, y := float64(a[i]), float64(b[i])
		dot += x * y
		normA += x * x
		normB += y * y
	}
	if denom := math.Sqrt(normA) * math.Sqrt(normB); denom > 1e-12 {
		return dot / denom
	}
	return 0
}
