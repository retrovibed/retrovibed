package acoustics

import (
	"math"
	"math/cmplx"
)

const (
	SampleRate       = 11025
	WindowSize       = 1024
	HopSize          = 512
	FFTBins          = WindowSize/2 + 1 // 513
	NumMelFilters    = 26
	NumMFCC          = 13 // C1 through C13, excluding C0
	NumChroma        = 12
	NumContrastBands = 6
	FeatureDim       = 64  // per-window feature count
	VectorDim        = 128 // final vector: mean + inter-window σ
)

var (
	twiddleFactors [WindowSize / 2]complex128
	bitReversal    [WindowSize]uint16
	hannWindow     [WindowSize]float32
	chromaBinMap   [FFTBins]uint8
	contrastEdges  [NumContrastBands + 1]int
	melFilterbank  []melFilter
)

type melFilter struct {
	start   int
	weights []float32
}

func init() {
	initTwiddle()
	initBitReversal()
	initHann()
	initMelFilterbank()
	initChromaBinMap()
	initContrastEdges()
}

func initTwiddle() {
	for k := range WindowSize / 2 {
		twiddleFactors[k] = cmplx.Rect(1, -2.0*math.Pi*float64(k)/float64(WindowSize))
	}
}

func initBitReversal() {
	bits := 0
	for v := WindowSize; v > 1; v >>= 1 {
		bits++
	}
	for i := range WindowSize {
		rev, val := 0, i
		for range bits {
			rev = (rev << 1) | (val & 1)
			val >>= 1
		}
		bitReversal[i] = uint16(rev)
	}
}

func initHann() {
	for i := range WindowSize {
		hannWindow[i] = float32(0.5 * (1.0 - math.Cos(2.0*math.Pi*float64(i)/float64(WindowSize-1))))
	}
}

func hzToMel(hz float64) float64  { return 2595.0 * math.Log10(1.0+hz/700.0) }
func melToHz(mel float64) float64 { return 700.0 * (math.Pow(10.0, mel/2595.0) - 1.0) }

func initMelFilterbank() {
	nyquist := float64(SampleRate) / 2.0
	melHigh := hzToMel(nyquist)

	nPoints := NumMelFilters + 2
	binPoints := make([]int, nPoints)
	for i := range nPoints {
		hz := melToHz(float64(i) * melHigh / float64(nPoints-1))
		binPoints[i] = int(math.Floor(hz / nyquist * float64(FFTBins-1)))
	}

	melFilterbank = make([]melFilter, NumMelFilters)
	for f := range NumMelFilters {
		start, mid, end := binPoints[f], binPoints[f+1], binPoints[f+2]
		if start == mid {
			mid = start + 1
		}
		if mid == end {
			end = mid + 1
		}
		end = min(end, FFTBins)

		weights := make([]float32, end-start)
		for b := start; b < mid; b++ {
			weights[b-start] = float32(b-start) / float32(mid-start)
		}
		for b := mid; b < end; b++ {
			weights[b-start] = float32(end-b) / float32(end-mid)
		}
		melFilterbank[f] = melFilter{start: start, weights: weights}
	}
}

func initChromaBinMap() {
	for i := 1; i < FFTBins; i++ {
		freq := float64(i) * float64(SampleRate) / float64(WindowSize)
		class := int(math.Round(12.0*math.Log2(freq/440.0))) % NumChroma
		if class < 0 {
			class += NumChroma
		}
		chromaBinMap[i] = uint8(class)
	}
}

func initContrastEdges() {
	lowBin := 4 // ~43 Hz
	contrastEdges[0] = lowBin
	for i := 1; i <= NumContrastBands; i++ {
		contrastEdges[i] = min(lowBin*(1<<i), FFTBins)
	}
}
