package acoustics

// STFT computes a Short-Time Fourier Transform on mono float32 PCM samples.
// Returns a magnitude spectrogram: nFrames × FFTBins, stored row-major.
func STFT(samples []float32) []float32 {
	if len(samples) < WindowSize {
		return nil
	}

	nFrames := SpectrogramFrames(len(samples))
	mag := make([]float32, nFrames*FFTBins)
	fftBuf := make([]complex128, WindowSize)

	for f := range nFrames {
		offset := f * HopSize
		for i := range WindowSize {
			fftBuf[i] = complex(float64(samples[offset+i])*float64(hannWindow[i]), 0)
		}
		FFT(fftBuf)
		Magnitude(fftBuf, mag[f*FFTBins:])
	}

	return mag
}

func SpectrogramFrames(nSamples int) int {
	if nSamples < WindowSize {
		return 0
	}
	return (nSamples - WindowSize) / HopSize + 1
}
