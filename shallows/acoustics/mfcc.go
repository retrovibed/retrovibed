package acoustics

import "math"

// MelSpectrum applies the sparse mel filterbank to a magnitude frame,
// writing log-mel energies into melOut. len(melOut) must be >= NumMelFilters.
func MelSpectrum(magFrame []float32, melOut []float32) {
	for f, filt := range melFilterbank {
		energy := float32(0)
		for i, w := range filt.weights {
			energy += w * magFrame[filt.start+i]
		}
		if energy < 1e-10 {
			energy = 1e-10
		}
		melOut[f] = float32(math.Log(float64(energy)))
	}
}

// DCT computes a Type-II DCT on in, writing coefficients 1 through NumMFCC
// into out (skipping C0). len(in) must be >= NumMelFilters, len(out) >= NumMFCC.
func DCT(in []float32, out []float32) {
	n := NumMelFilters
	for k := range NumMFCC {
		sum := float64(0)
		kk := k + 1 // skip C0
		for i := range n {
			sum += float64(in[i]) * math.Cos(math.Pi*float64(kk)*(float64(i)+0.5)/float64(n))
		}
		out[k] = float32(sum)
	}
}

// ExtractMFCC computes 13 MFCCs (C1-C13) from a magnitude spectrum frame.
// Uses the pre-computed mel filterbank and DCT.
func ExtractMFCC(magFrame []float32, mfccOut []float32) {
	melBuf := make([]float32, NumMelFilters)
	MelSpectrum(magFrame, melBuf)
	DCT(melBuf, mfccOut)
}
