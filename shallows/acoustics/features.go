package acoustics

import (
	"math"
	"slices"
)

// WindowFeatures holds the 64-dimensional feature vector for a single window.
type WindowFeatures [FeatureDim]float32

// ExtractWindowFeatures computes the full 64-dim feature vector from a
// magnitude spectrogram (row-major, nFrames × FFTBins).
func ExtractWindowFeatures(mag []float32, nFrames int) WindowFeatures {
	var wf WindowFeatures
	mfccBuf := make([]float32, NumMFCC)
	melBuf := make([]float32, NumMelFilters)
	chromaBuf := make([]float32, NumChroma)

	spectralAcc := newAccumulators(4) // centroid, bandwidth, rolloff, flux
	mfccAcc := newAccumulators(NumMFCC)
	chromaAcc := newAccumulators(NumChroma)
	contrastAcc := newAccumulators(NumContrastBands)
	onsetAcc := newAccumulators(1)
	rmsAcc := newAccumulators(1)

	var prevFrame []float32
	var onsetTimes []int
	medianBuf := make([]float32, 0, nFrames)

	for f := range nFrames {
		frame := mag[f*FFTBins : (f+1)*FFTBins]

		centroid, bandwidth, rolloff := spectralShape(frame)
		flux := float32(0)
		if prevFrame != nil {
			flux = spectralFlux(prevFrame, frame)
		}
		spectralAcc.add(0, centroid)
		spectralAcc.add(1, bandwidth)
		spectralAcc.add(2, rolloff)
		spectralAcc.add(3, flux)

		MelSpectrum(frame, melBuf)
		DCT(melBuf, mfccBuf)
		for i, v := range mfccBuf {
			mfccAcc.add(i, v)
		}

		ExtractChroma(frame, chromaBuf)
		for i, v := range chromaBuf {
			chromaAcc.add(i, v)
		}

		for b := range NumContrastBands {
			contrastAcc.add(b, bandContrast(frame, b))
		}

		rmsAcc.add(0, frameRMS(frame))

		// Onset detection: half-wave rectified flux above running median
		if flux > 0 {
			medianBuf = append(medianBuf, flux)
			if flux > runningMedian(medianBuf)*1.5 {
				if len(onsetTimes) == 0 || f-onsetTimes[len(onsetTimes)-1] > 1 {
					onsetTimes = append(onsetTimes, f)
				}
			}
		}
		onsetAcc.add(0, flux)
		prevFrame = frame
	}

	n := float32(nFrames)
	idx := 0

	for i := range 4 { // spectral shape: centroid, bandwidth, rolloff, flux
		wf[idx], wf[idx+1] = spectralAcc.meanStd(i, n)
		idx += 2
	}
	for i := range NumMFCC {
		wf[idx], wf[idx+1] = mfccAcc.meanStd(i, n)
		idx += 2
	}
	for i := range NumChroma { // mean only
		wf[idx] = chromaAcc.mean(i, n)
		idx++
	}
	for i := range NumContrastBands {
		wf[idx], wf[idx+1] = contrastAcc.meanStd(i, n)
		idx += 2
	}

	wf[idx] = estimateTempo(onsetTimes, nFrames)
	wf[idx+1] = beatRegularity(onsetTimes)
	wf[idx+2], wf[idx+3] = onsetAcc.meanStd(0, n)
	wf[idx+4], wf[idx+5] = rmsAcc.meanStd(0, n)

	return wf
}

func spectralShape(frame []float32) (centroid, bandwidth, rolloff float32) {
	totalEnergy := float32(0)
	weightedSum := float32(0)
	for i, v := range frame {
		totalEnergy += v
		weightedSum += float32(i) * v
	}
	if totalEnergy < 1e-10 {
		return 0, 0, 0
	}

	centroid = weightedSum / totalEnergy

	varSum := float32(0)
	for i, v := range frame {
		d := float32(i) - centroid
		varSum += d * d * v
	}
	bandwidth = float32(math.Sqrt(float64(varSum / totalEnergy)))

	cumulative := float32(0)
	threshold := totalEnergy * 0.85
	for i, v := range frame {
		cumulative += v
		if cumulative >= threshold {
			rolloff = float32(i)
			break
		}
	}

	return centroid, bandwidth, rolloff
}

// spectralFlux: half-wave rectified L2 norm of frame-over-frame increase.
func spectralFlux(prev, curr []float32) float32 {
	sum := float32(0)
	for i := range FFTBins {
		if d := curr[i] - prev[i]; d > 0 {
			sum += d * d
		}
	}
	return float32(math.Sqrt(float64(sum) / float64(FFTBins)))
}

// bandContrast: difference between mean of top 20% and bottom 20% magnitudes.
func bandContrast(frame []float32, band int) float32 {
	lo, hi := contrastEdges[band], min(contrastEdges[band+1], FFTBins)
	if lo >= hi {
		return 0
	}

	sorted := make([]float32, hi-lo)
	copy(sorted, frame[lo:hi])
	slices.Sort(sorted)

	n := len(sorted)
	fifth := max(n/5, 1)

	peak := float32(0)
	for _, v := range sorted[n-fifth:] {
		peak += v
	}
	valley := float32(0)
	for _, v := range sorted[:fifth] {
		valley += v
	}
	return peak/float32(fifth) - valley/float32(fifth)
}

func frameRMS(frame []float32) float32 {
	sum := float32(0)
	for _, v := range frame {
		sum += v * v
	}
	return float32(math.Sqrt(float64(sum) / float64(len(frame))))
}

func estimateTempo(onsetTimes []int, nFrames int) float32 {
	if len(onsetTimes) < 2 {
		return 0
	}

	// Autocorrelation of binary onset signal over 60-200 BPM lag range.
	frameRate := float64(SampleRate) / float64(HopSize)
	minLag := int(frameRate * 60.0 / 200.0)
	maxLag := min(int(frameRate*60.0/60.0), nFrames/2)
	if minLag >= maxLag {
		return 0
	}

	onsetSignal := make([]float32, nFrames)
	for _, t := range onsetTimes {
		onsetSignal[t] = 1
	}

	bestLag, bestCorr := minLag, float32(0)
	for lag := minLag; lag <= maxLag; lag++ {
		corr := float32(0)
		for i := 0; i+lag < nFrames; i++ {
			corr += onsetSignal[i] * onsetSignal[i+lag]
		}
		if corr > bestCorr {
			bestCorr = corr
			bestLag = lag
		}
	}

	if bestCorr == 0 {
		return 0
	}
	return float32(60.0 * frameRate / float64(bestLag))
}

func beatRegularity(onsetTimes []int) float32 {
	if len(onsetTimes) < 3 {
		return 0
	}

	intervals := make([]float32, len(onsetTimes)-1)
	for i := range intervals {
		intervals[i] = float32(onsetTimes[i+1] - onsetTimes[i])
	}

	mean := float32(0)
	for _, v := range intervals {
		mean += v
	}
	mean /= float32(len(intervals))

	variance := float32(0)
	for _, v := range intervals {
		d := v - mean
		variance += d * d
	}
	return float32(math.Sqrt(float64(variance / float32(len(intervals)))))
}

func runningMedian(buf []float32) float32 {
	if len(buf) == 0 {
		return 0
	}
	sorted := make([]float32, len(buf))
	copy(sorted, buf)
	slices.Sort(sorted)
	return sorted[len(sorted)/2]
}

// accumulators tracks running sum and sum-of-squares for mean/σ computation.
type accumulators struct {
	sum   []float64
	sumSq []float64
}

func newAccumulators(n int) accumulators {
	return accumulators{sum: make([]float64, n), sumSq: make([]float64, n)}
}

func (a *accumulators) add(ch int, v float32) {
	fv := float64(v)
	a.sum[ch] += fv
	a.sumSq[ch] += fv * fv
}

func (a *accumulators) mean(ch int, n float32) float32 {
	return float32(a.sum[ch] / float64(n))
}

func (a *accumulators) meanStd(ch int, n float32) (float32, float32) {
	m := a.sum[ch] / float64(n)
	v := max(a.sumSq[ch]/float64(n)-m*m, 0)
	return float32(m), float32(math.Sqrt(v))
}
