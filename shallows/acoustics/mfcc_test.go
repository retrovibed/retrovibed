package acoustics

import (
	"math"
	"testing"
)

func TestMelSpectrum(t *testing.T) {
	t.Run("silence produces log floor", func(t *testing.T) {
		magFrame := make([]float32, FFTBins)
		melOut := make([]float32, NumMelFilters)
		MelSpectrum(magFrame, melOut)

		logFloor := float32(math.Log(1e-10))
		for i, v := range melOut {
			if math.Abs(float64(v-logFloor)) > 0.01 {
				t.Fatalf("filter %d: expected %f, got %f", i, logFloor, v)
			}
		}
	})

	t.Run("flat spectrum produces non-zero output", func(t *testing.T) {
		magFrame := make([]float32, FFTBins)
		for i := range magFrame {
			magFrame[i] = 1.0
		}
		melOut := make([]float32, NumMelFilters)
		MelSpectrum(magFrame, melOut)

		for i, v := range melOut {
			if v <= 0 {
				t.Fatalf("filter %d: expected positive log energy, got %f", i, v)
			}
		}
	})

	t.Run("energy in low band", func(t *testing.T) {
		// Energy only in the first few bins should activate low mel filters
		// more than high ones.
		magFrame := make([]float32, FFTBins)
		for i := range 20 {
			magFrame[i] = 10.0
		}
		melOut := make([]float32, NumMelFilters)
		MelSpectrum(magFrame, melOut)

		if melOut[0] <= melOut[NumMelFilters-1] {
			t.Fatalf("low filter (%f) should exceed high filter (%f) for low-band energy",
				melOut[0], melOut[NumMelFilters-1])
		}
	})
}

func TestDCT(t *testing.T) {
	t.Run("output dimension", func(t *testing.T) {
		in := make([]float32, NumMelFilters)
		for i := range in {
			in[i] = float32(i + 1)
		}
		out := make([]float32, NumMFCC)
		DCT(in, out)

		allZero := true
		for _, v := range out {
			if v != 0 {
				allZero = false
				break
			}
		}
		if allZero {
			t.Fatal("DCT produced all zeros for non-zero input")
		}
	})

	t.Run("constant input", func(t *testing.T) {
		// Constant input: C0 would be the sum, but we skip C0.
		// C1+ should reflect the constant via cosine basis.
		in := make([]float32, NumMelFilters)
		for i := range in {
			in[i] = 5.0
		}
		out := make([]float32, NumMFCC)
		DCT(in, out)

		// For constant input, DCT coefficients k>=1 should be near zero
		// because cos(pi*k*(i+0.5)/N) sums to ~0 for k>=1.
		for k, v := range out {
			if math.Abs(float64(v)) > 1e-3 {
				t.Fatalf("C%d: expected ~0 for constant input, got %f", k+1, v)
			}
		}
	})
}

func TestExtractMFCC(t *testing.T) {
	t.Run("produces non-zero coefficients", func(t *testing.T) {
		magFrame := make([]float32, FFTBins)
		for i := range magFrame {
			magFrame[i] = float32(1.0 + math.Sin(float64(i)*0.1))
		}
		mfccOut := make([]float32, NumMFCC)
		ExtractMFCC(magFrame, mfccOut)

		nonZero := 0
		for _, v := range mfccOut {
			if v != 0 {
				nonZero++
			}
		}
		if nonZero < NumMFCC/2 {
			t.Fatalf("expected most MFCCs non-zero, got %d/%d", nonZero, NumMFCC)
		}
	})
}
