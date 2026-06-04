package acoustics

// ExtractChroma maps a magnitude spectrum frame onto 12 pitch classes
// using the pre-computed chromaBinMap.
func ExtractChroma(magFrame []float32, chromaOut []float32) {
	clear(chromaOut[:NumChroma])
	for i := 1; i < FFTBins; i++ {
		chromaOut[chromaBinMap[i]] += magFrame[i]
	}
}
