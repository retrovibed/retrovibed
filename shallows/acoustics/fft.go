package acoustics

import "math"

// FFT computes an in-place radix-2 DIT FFT on buf using pre-computed
// twiddle factors and bit-reversal permutation. len(buf) must equal WindowSize.
func FFT(buf []complex128) {
	n := len(buf)

	for i := range n {
		if j := int(bitReversal[i]); i < j {
			buf[i], buf[j] = buf[j], buf[i]
		}
	}

	for size := 2; size <= n; size <<= 1 {
		half := size / 2
		step := n / size
		for start := 0; start < n; start += size {
			for k := range half {
				t := twiddleFactors[k*step] * buf[start+half+k]
				buf[start+half+k] = buf[start+k] - t
				buf[start+k] += t
			}
		}
	}
}

// Magnitude writes the magnitude spectrum of the first FFTBins elements
// of fftOut into mag. Both slices must have length >= FFTBins.
func Magnitude(fftOut []complex128, mag []float32) {
	for i := range FFTBins {
		r, im := real(fftOut[i]), imag(fftOut[i])
		mag[i] = float32(math.Sqrt(r*r + im*im))
	}
}
