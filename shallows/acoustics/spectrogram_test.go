package acoustics

import (
	"math"
	"testing"
)

func TestSTFT(t *testing.T) {
	t.Run("frame count", func(t *testing.T) {
		nSamples := 11025 // 1 second at 11,025 Hz
		expected := SpectrogramFrames(nSamples)
		got := (nSamples - WindowSize) / HopSize + 1

		if expected != got {
			t.Fatalf("expected %d frames, got %d", expected, got)
		}

		samples := make([]float32, nSamples)
		mag := STFT(samples)
		if len(mag) != expected*FFTBins {
			t.Fatalf("expected %d values, got %d", expected*FFTBins, len(mag))
		}
	})

	t.Run("too short returns nil", func(t *testing.T) {
		samples := make([]float32, WindowSize-1)
		if mag := STFT(samples); mag != nil {
			t.Fatalf("expected nil for short input, got %d values", len(mag))
		}
	})

	t.Run("sine energy in correct bin", func(t *testing.T) {
		// Generate a 1-second sine at 440 Hz.
		// At 11,025 Hz sample rate with 1024-point FFT, bin spacing is ~10.77 Hz.
		// 440 Hz should land near bin 40 (440 / 10.77 ≈ 40.85).
		nSamples := 11025
		samples := make([]float32, nSamples)
		for i := range samples {
			samples[i] = float32(math.Sin(2.0 * math.Pi * 440.0 * float64(i) / float64(SampleRate)))
		}

		mag := STFT(samples)
		nFrames := SpectrogramFrames(nSamples)
		targetBin := 41 // closest integer bin to 440 Hz

		for f := range nFrames {
			frame := mag[f*FFTBins : (f+1)*FFTBins]

			peakBin := 0
			peakVal := float32(0)
			for b, v := range frame {
				if v > peakVal {
					peakVal = v
					peakBin = b
				}
			}

			diff := peakBin - targetBin
			if diff < -2 || diff > 2 {
				t.Fatalf("frame %d: peak at bin %d, expected near %d", f, peakBin, targetBin)
			}
		}
	})

	t.Run("silence produces zero magnitude", func(t *testing.T) {
		samples := make([]float32, 2*WindowSize)
		mag := STFT(samples)

		for i, v := range mag {
			if v != 0 {
				t.Fatalf("index %d: expected 0, got %f", i, v)
			}
		}
	})
}

func TestSpectrogramFrames(t *testing.T) {
	t.Run("exact window", func(t *testing.T) {
		if got := SpectrogramFrames(WindowSize); got != 1 {
			t.Fatalf("expected 1 frame, got %d", got)
		}
	})

	t.Run("below window", func(t *testing.T) {
		if got := SpectrogramFrames(WindowSize - 1); got != 0 {
			t.Fatalf("expected 0 frames, got %d", got)
		}
	})

	t.Run("window plus hop", func(t *testing.T) {
		if got := SpectrogramFrames(WindowSize + HopSize); got != 2 {
			t.Fatalf("expected 2 frames, got %d", got)
		}
	})
}
