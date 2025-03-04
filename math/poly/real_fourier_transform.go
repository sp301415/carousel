package poly

import (
	"math"
	"math/cmplx"

	"github.com/sp301415/carousel/math/num"
)

// RealFourierTransformer computes the Fast Fourier Transform over R[X]/(X^N - 1).
//
// RealFourierTransformer is safe for concurrent use.
type RealFourierTransformer struct {
	degree int

	// tw is the twiddle factor for FFT.
	tw []complex128
	// twInv is the twiddle factor for inverse FFT.
	twInv []complex128
}

// NewRealFourierTransformer creates a new RealFourierTransformer.
func NewRealFourierTransformer(N int) *RealFourierTransformer {
	switch {
	case !num.IsPowerOfTwo(N):
		panic("NewFourierTransformer: degree must be a power of two")
	case N < MinDegree:
		panic("NewFourierTransformer: degree smaller than MinDegree")
	}

	tw, twInv := genTwiddleFactorsRealFFT(N)
	return &RealFourierTransformer{
		degree: N,
		tw:     tw,
		twInv:  twInv,
	}
}

func genTwiddleFactorsRealFFT(N int) (tw, twInv []complex128) {
	twRealFFTIdx := make([]int, N/4)
	t := N / 8
	twRealFFTIdx[t] = 1
	for m := 4; m <= N/4; m <<= 1 {
		t >>= 1
		twRealFFTIdx[t] = twRealFFTIdx[t<<1] << 1
		for j := 3; j < m; j += 2 {
			twRealFFTIdx[j*t] = 2*twRealFFTIdx[t] - twRealFFTIdx[(j-1)*t]
		}
	}

	twRealFFTRef := make([]complex128, N/4)
	twRealInvFFTRef := make([]complex128, N/4)
	for i, e := range twRealFFTIdx {
		x := -2 * math.Pi * float64(e) / float64(N)
		twRealFFTRef[i] = cmplx.Exp(complex(0, x))
		twRealInvFFTRef[i] = cmplx.Exp(-complex(0, x))
	}

	var w int
	tw = make([]complex128, N/2)
	w = 0
	for m := 1; m <= N/4; m <<= 1 {
		for i := 0; i < m; i++ {
			tw[w] = twRealFFTRef[i]
			w++
		}
	}

	twInv = make([]complex128, N/2)
	w = 0
	for m := N / 4; m >= 1; m >>= 1 {
		for i := 0; i < m; i++ {
			twInv[w] = twRealInvFFTRef[i]
			w++
		}
	}

	return tw, twInv
}

// FourierTransformInPlace computes the FFT of p in place.
func (fft *RealFourierTransformer) FourierTransformInPlace(p []float64) {
	rfftInPlace(p, fft.tw)
}

// InvFourierTransformInPlace computes the inverse FFT of p in place.
func (fft *RealFourierTransformer) InvFourierTransformInPlace(p []float64) {
	rifftInPlace(p, fft.twInv)
}

// ConvolveAssign computes the convolution of fp0, fp1 and assigns the result to fpOut.
func (fft *RealFourierTransformer) ConvolveAssign(fp0, fp1, fpOut []float64) {
	convolveAssign(fp0, fp1, fpOut)
}
