package acoustics

import (
	"math"
	"testing"
)

func TestExtractChroma(t *testing.T) {
	t.Run("silence produces zero chroma", func(t *testing.T) {
		magFrame := make([]float32, FFTBins)
		chromaOut := make([]float32, NumChroma)
		ExtractChroma(magFrame, chromaOut)

		for i, v := range chromaOut {
			if v != 0 {
				t.Fatalf("class %d: expected 0, got %f", i, v)
			}
		}
	})

	t.Run("single bin concentrates in one class", func(t *testing.T) {
		magFrame := make([]float32, FFTBins)
		bin := 100
		magFrame[bin] = 10.0

		chromaOut := make([]float32, NumChroma)
		ExtractChroma(magFrame, chromaOut)

		expectedClass := chromaBinMap[bin]
		if chromaOut[expectedClass] != 10.0 {
			t.Fatalf("class %d: expected 10.0, got %f", expectedClass, chromaOut[expectedClass])
		}

		for i, v := range chromaOut {
			if i != int(expectedClass) && v != 0 {
				t.Fatalf("class %d: expected 0, got %f", i, v)
			}
		}
	})

	t.Run("sums across octaves", func(t *testing.T) {
		// Two bins mapped to the same pitch class should sum.
		magFrame := make([]float32, FFTBins)

		// Find two bins with the same chroma class.
		targetClass := chromaBinMap[50]
		var secondBin int
		for i := 51; i < FFTBins; i++ {
			if chromaBinMap[i] == targetClass {
				secondBin = i
				break
			}
		}
		if secondBin == 0 {
			t.Skip("could not find two bins with same chroma class")
		}

		magFrame[50] = 3.0
		magFrame[secondBin] = 7.0

		chromaOut := make([]float32, NumChroma)
		ExtractChroma(magFrame, chromaOut)

		if math.Abs(float64(chromaOut[targetClass])-10.0) > 1e-5 {
			t.Fatalf("class %d: expected 10.0, got %f", targetClass, chromaOut[targetClass])
		}
	})
}
