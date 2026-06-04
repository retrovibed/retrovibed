package acoustics

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestFFT(t *testing.T) {
	t.Run("impulse at zero", func(t *testing.T) {
		buf := make([]complex128, WindowSize)
		buf[0] = 1.0

		FFT(buf)

		for i, v := range buf {
			if mag := cmplx.Abs(v); math.Abs(mag-1.0) > 1e-10 {
				t.Fatalf("bin %d: expected magnitude 1.0, got %f", i, mag)
			}
		}
	})

	t.Run("DC signal", func(t *testing.T) {
		// A constant signal of amplitude 1 should have all energy in bin 0.
		buf := make([]complex128, WindowSize)
		for i := range buf {
			buf[i] = 1.0
		}

		FFT(buf)

		dc := cmplx.Abs(buf[0])
		if math.Abs(dc-float64(WindowSize)) > 1e-6 {
			t.Fatalf("DC bin: expected %d, got %f", WindowSize, dc)
		}
		for i := 1; i < WindowSize; i++ {
			mag := cmplx.Abs(buf[i])
			if mag > 1e-10 {
				t.Fatalf("bin %d: expected 0, got %f", i, mag)
			}
		}
	})

	t.Run("pure sine", func(t *testing.T) {
		k := 50
		buf := make([]complex128, WindowSize)
		for i := range buf {
			buf[i] = complex(math.Sin(2.0*math.Pi*float64(k)*float64(i)/float64(WindowSize)), 0)
		}

		FFT(buf)

		peakMag := float64(WindowSize) / 2.0
		for i, v := range buf {
			mag := cmplx.Abs(v)
			if i == k || i == WindowSize-k {
				if math.Abs(mag-peakMag) > 1e-6 {
					t.Fatalf("bin %d: expected %f, got %f", i, peakMag, mag)
				}
			} else if mag > 1e-6 {
				t.Fatalf("bin %d: expected ~0, got %f", i, mag)
			}
		}
	})

	t.Run("Parseval's theorem", func(t *testing.T) {
		// Energy in time domain should equal energy in frequency domain / N.
		buf := make([]complex128, WindowSize)
		for i := range buf {
			buf[i] = complex(float64(i%7)-3.0, 0)
		}

		timeEnergy := 0.0
		for _, v := range buf {
			timeEnergy += real(v) * real(v)
		}

		FFT(buf)

		freqEnergy := 0.0
		for _, v := range buf {
			freqEnergy += real(v)*real(v) + imag(v)*imag(v)
		}
		freqEnergy /= float64(WindowSize)

		if math.Abs(timeEnergy-freqEnergy) > 1e-6 {
			t.Fatalf("Parseval: time energy %f != freq energy %f", timeEnergy, freqEnergy)
		}
	})
}

func TestMagnitude(t *testing.T) {
	t.Run("known values", func(t *testing.T) {
		fftOut := make([]complex128, WindowSize)
		fftOut[0] = complex(3, 4) // magnitude 5
		fftOut[1] = complex(0, 1) // magnitude 1
		fftOut[2] = complex(1, 0) // magnitude 1

		mag := make([]float32, FFTBins)
		Magnitude(fftOut, mag)

		if math.Abs(float64(mag[0])-5.0) > 1e-5 {
			t.Fatalf("bin 0: expected 5.0, got %f", mag[0])
		}
		if math.Abs(float64(mag[1])-1.0) > 1e-5 {
			t.Fatalf("bin 1: expected 1.0, got %f", mag[1])
		}
		if math.Abs(float64(mag[2])-1.0) > 1e-5 {
			t.Fatalf("bin 2: expected 1.0, got %f", mag[2])
		}
	})
}
