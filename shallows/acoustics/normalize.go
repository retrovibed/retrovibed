package acoustics

import "math"

// FeatureVector is the final 128-dimensional acoustic fingerprint.
type FeatureVector [VectorDim]float64

// AggregateWindows computes the mean and inter-window standard deviation
// of three per-window feature vectors, producing the final 128-dim vector.
func AggregateWindows(windows [3]WindowFeatures) FeatureVector {
	var fv FeatureVector

	for d := range FeatureDim {
		a, b, c := float64(windows[0][d]), float64(windows[1][d]), float64(windows[2][d])
		mean := (a + b + c) / 3.0
		fv[d] = mean

		v := ((a-mean)*(a-mean) + (b-mean)*(b-mean) + (c-mean)*(c-mean)) / 3.0
		fv[FeatureDim+d] = math.Sqrt(v)
	}

	return fv
}

// RunningStats tracks incremental count, sum, and sum-of-squares per dimension
// for z-score normalization.
type RunningStats struct {
	Count int64
	Sum   [VectorDim]float64
	SumSq [VectorDim]float64
}

// Update adds a feature vector to the running statistics.
func (s *RunningStats) Update(fv FeatureVector) {
	s.Count++
	for i, v := range fv {
		s.Sum[i] += v
		s.SumSq[i] += v * v
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
		out[i] = (fv[i] - mean) / math.Sqrt(variance)
	}

	return out
}

// NormalizeAll re-normalizes a batch of feature vectors using the current statistics.
func (s *RunningStats) NormalizeAll(vectors []FeatureVector) []FeatureVector {
	result := make([]FeatureVector, len(vectors))
	for i, fv := range vectors {
		result[i] = s.Normalize(fv)
	}
	return result
}

// Recompute rebuilds running statistics from a batch of feature vectors.
func (s *RunningStats) Recompute(vectors []FeatureVector) {
	*s = RunningStats{}
	for _, fv := range vectors {
		s.Update(fv)
	}
}
