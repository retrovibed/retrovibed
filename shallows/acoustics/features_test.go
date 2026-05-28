package acoustics

import (
	"math"
	"testing"
)

func TestExtractWindowFeatures(t *testing.T) {
	generateSine := func(hz float64, nSamples int) []float32 {
		samples := make([]float32, nSamples)
		for i := range samples {
			samples[i] = float32(math.Sin(2.0 * math.Pi * hz * float64(i) / float64(SampleRate)))
		}
		return samples
	}

	t.Run("dimension count", func(t *testing.T) {
		samples := generateSine(440, SampleRate) // 1 second at 440 Hz
		mag := STFT(samples)
		nFrames := SpectrogramFrames(len(samples))

		wf := ExtractWindowFeatures(mag, nFrames)
		nonZero := 0
		for _, v := range wf {
			if v != 0 {
				nonZero++
			}
		}

		if nonZero == 0 {
			t.Fatal("all features are zero for sine input")
		}
	})

	t.Run("silence produces near-zero features", func(t *testing.T) {
		samples := make([]float32, SampleRate)
		mag := STFT(samples)
		nFrames := SpectrogramFrames(len(samples))

		wf := ExtractWindowFeatures(mag, nFrames)

		// Spectral shape features (indices 0-7) should all be zero.
		for i := range 8 {
			if wf[i] != 0 {
				t.Fatalf("feature[%d]: expected 0, got %f", i, wf[i])
			}
		}
	})

	t.Run("low tone vs high tone centroid", func(t *testing.T) {
		lowSamples := generateSine(200, SampleRate)
		highSamples := generateSine(2000, SampleRate)

		lowMag := STFT(lowSamples)
		highMag := STFT(highSamples)
		nFrames := SpectrogramFrames(SampleRate)

		lowF := ExtractWindowFeatures(lowMag, nFrames)
		highF := ExtractWindowFeatures(highMag, nFrames)

		// Index 0 is centroid mean. Higher tone should have higher centroid.
		if lowF[0] >= highF[0] {
			t.Fatalf("low centroid (%f) should be less than high centroid (%f)", lowF[0], highF[0])
		}
	})

	t.Run("RMS reflects amplitude", func(t *testing.T) {
		quiet := make([]float32, SampleRate)
		loud := make([]float32, SampleRate)
		for i := range quiet {
			quiet[i] = 0.1 * float32(math.Sin(2.0*math.Pi*440.0*float64(i)/float64(SampleRate)))
			loud[i] = 1.0 * float32(math.Sin(2.0*math.Pi*440.0*float64(i)/float64(SampleRate)))
		}

		quietMag := STFT(quiet)
		loudMag := STFT(loud)
		nFrames := SpectrogramFrames(SampleRate)

		quietF := ExtractWindowFeatures(quietMag, nFrames)
		loudF := ExtractWindowFeatures(loudMag, nFrames)

		// RMS mean is at index 62 (feature dim - 2)
		rmsIdx := FeatureDim - 2
		if quietF[rmsIdx] >= loudF[rmsIdx] {
			t.Fatalf("quiet RMS (%f) should be less than loud RMS (%f)", quietF[rmsIdx], loudF[rmsIdx])
		}
	})
}

func TestSpectralShape(t *testing.T) {
	t.Run("flat spectrum", func(t *testing.T) {
		frame := make([]float32, FFTBins)
		for i := range frame {
			frame[i] = 1.0
		}

		centroid, bandwidth, rolloff := spectralShape(frame)

		expectedCentroid := float32(FFTBins-1) / 2.0
		if math.Abs(float64(centroid-expectedCentroid)) > 1.0 {
			t.Fatalf("centroid: expected ~%f, got %f", expectedCentroid, centroid)
		}

		if bandwidth <= 0 {
			t.Fatalf("bandwidth: expected positive, got %f", bandwidth)
		}

		expectedRolloff := float32(FFTBins) * 0.85
		if math.Abs(float64(rolloff-expectedRolloff)) > 5.0 {
			t.Fatalf("rolloff: expected ~%f, got %f", expectedRolloff, rolloff)
		}
	})

	t.Run("zero frame", func(t *testing.T) {
		frame := make([]float32, FFTBins)
		centroid, bandwidth, rolloff := spectralShape(frame)
		if centroid != 0 || bandwidth != 0 || rolloff != 0 {
			t.Fatalf("zero frame should produce zero shape, got %f %f %f", centroid, bandwidth, rolloff)
		}
	})
}

func TestBandContrast(t *testing.T) {
	t.Run("flat band has zero contrast", func(t *testing.T) {
		frame := make([]float32, FFTBins)
		for i := range frame {
			frame[i] = 5.0
		}

		for b := range NumContrastBands {
			c := bandContrast(frame, b)
			if math.Abs(float64(c)) > 1e-5 {
				t.Fatalf("band %d: expected ~0 contrast for flat spectrum, got %f", b, c)
			}
		}
	})

	t.Run("peaked band has positive contrast", func(t *testing.T) {
		frame := make([]float32, FFTBins)
		lo, hi := contrastEdges[0], contrastEdges[1]
		for i := lo; i < hi; i++ {
			frame[i] = 1.0
		}
		// Add a peak
		if lo < hi {
			frame[lo] = 100.0
		}

		c := bandContrast(frame, 0)
		if c <= 0 {
			t.Fatalf("expected positive contrast, got %f", c)
		}
	})
}

func TestEstimateTempo(t *testing.T) {
	t.Run("regular onsets at 120 BPM", func(t *testing.T) {
		// 120 BPM = 2 beats/sec. Frame rate ~21.5 fps. ~10.75 frames/beat.
		frameRate := float64(SampleRate) / float64(HopSize)
		framesPerBeat := frameRate / 2.0 // 120 BPM = 2 Hz

		nFrames := 500
		var onsets []int
		for f := 0; f < nFrames; f += int(framesPerBeat) {
			onsets = append(onsets, f)
		}

		bpm := estimateTempo(onsets, nFrames)
		if math.Abs(float64(bpm)-120.0) > 10.0 {
			t.Fatalf("expected ~120 BPM, got %f", bpm)
		}
	})

	t.Run("no onsets returns zero", func(t *testing.T) {
		if bpm := estimateTempo(nil, 100); bpm != 0 {
			t.Fatalf("expected 0 BPM, got %f", bpm)
		}
	})
}
