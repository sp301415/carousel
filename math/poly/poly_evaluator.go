package poly

import (
	"math"

	"github.com/sp301415/carousel/math/num"
	"github.com/sp301415/carousel/math/vec"
)

const (
	// MinDegree is the minimum degree of polynomial that Evaluator can handle.
	// Currently, this is set to 2^4, because AVX2 implementation of FFT and inverse FFT
	// handles first/last two loops separately.
	MinDegree = 1 << 4

	// ShortLogBound is a maximum bound for the coefficients of "short" polynomials
	// used in [*Evaluator.ShortFourierPolyMulPoly] functions.
	// Currently, this is set to 16 bits.
	ShortLogBound = 16

	// splitLogBound is denotes the maximum bits of N*B1^2*B2^2, where B1, B2 is the splitting bound of polynomial multiplication.
	// Currently, this is set to 48, which gives failure rate less than 2^-284.
	splitLogBound = 48
)

// Evaluator computes polynomial operation over a subring.
type Evaluator struct {
	Parameters EvaluatorParameters

	FourierTransformer *RealFourierTransformer

	// tw is the twiddle factor for FFT over the subring.
	tw []float64
	// twInv is the twiddle factor for inverse FFT over the subring.
	twInv []float64

	buffer evaluationBuffer
}

// evaluationBuffer is a buffer for Evaluator.
type evaluationBuffer struct {
	// pOut is the intermediate output polynomial for InPlace operations.
	pOut Poly
	// fpOut is the intermediate output fourier polynomial for InPlace operations.
	fpOut FourierPoly

	// fp is the FFT value of p.
	fp FourierPoly
	// fpInv is the InvFFT value of fp.
	fpInv FourierPoly
	// pInv is the InvFFT value of fp.
	pInv Poly

	// pSplit is the split value of p0 in [*Evaluator.ShortFourierPolyMulPoly].
	pSplit Poly
	// fpShortSplit is the fourier transformed pSplit in [*Evaluator.ShortFourierPolyMulPoly].
	fpShortSplit []FourierPoly
}

// NewEvaluatorWithParameters creates a new Evaluator with the given parameters.
func NewEvaluator(params EvaluatorParameters) *Evaluator {
	N := params.Degree()
	fft := NewRealFourierTransformer(N)

	twFFT := make([]float64, N)
	t := params.generator
	for i := 0; i < params.order; i++ {
		for j := 0; j < N; j++ {
			twFFT[j] += math.Cos(-2 * math.Pi * float64(t) / float64(params.cyclotomicDegree))
			t = (t * params.generator) % uint64(params.cyclotomicDegree)
		}
	}

	twInvFFT := make([]float64, N)
	for i := 0; i < N; i++ {
		twInvFFT[i] = (twFFT[i] - float64(params.order)) / float64(params.cyclotomicDegree)
	}

	vec.ReverseInPlace(twFFT[1:])
	fft.FourierTransformInPlace(twFFT)
	fft.FourierTransformInPlace(twInvFFT)

	for i := 0; i < N; i++ {
		twFFT[i] /= float64(N / 2)
		twInvFFT[i] /= float64(N / 2)
	}

	return &Evaluator{
		Parameters: params,

		FourierTransformer: fft,

		tw:    twFFT,
		twInv: twInvFFT,

		buffer: newEvaluationBuffer(N),
	}
}

// splitParameters generates splitBits and splitCount for [*Evaluator.MulPoly].
func splitParameters(N int) (splitBits int, splitCount int) {
	splitBits = (splitLogBound - num.Log2(N)) / 4
	splitCount = int(math.Ceil(64 / float64(splitBits)))
	return
}

// splitParametersShort generates splitBits and splitCount for [*Evaluator.ShortFourierPolyMulPoly].
func splitParametersShort(N int) (splitBits int, splitCount int) {
	splitBits = (splitLogBound - 2*ShortLogBound - num.Log2(N)) / 2
	splitCount = int(math.Ceil(64 / float64(splitBits)))
	return
}

// newEvaluationBuffer creates a new evaluationBuffer.
func newEvaluationBuffer(N int) evaluationBuffer {
	_, splitCount := splitParametersShort(N)

	fpShortSplit := make([]FourierPoly, splitCount)
	for i := 0; i < splitCount; i++ {
		fpShortSplit[i] = NewFourierPoly(N)
	}

	return evaluationBuffer{
		pOut:  NewPoly(N),
		fpOut: NewFourierPoly(N),

		fp:    NewFourierPoly(N),
		fpInv: NewFourierPoly(N),
		pInv:  NewPoly(N),

		pSplit:       NewPoly(N),
		fpShortSplit: fpShortSplit,
	}
}

// ShallowCopy returns a shallow copy of this Evaluator.
// Returned Evaluator is safe for concurrent use.
func (e *Evaluator) ShallowCopy() *Evaluator {
	return &Evaluator{
		Parameters: e.Parameters,

		FourierTransformer: e.FourierTransformer,

		tw:    e.tw,
		twInv: e.twInv,

		buffer: newEvaluationBuffer(e.Parameters.Degree()),
	}
}

// NewPoly creates a new polynomial with the same degree as the evaluator.
func (e *Evaluator) NewPoly() Poly {
	return Poly{Coeffs: make([]uint64, e.Parameters.degree)}
}

// NewFourierPoly creates a new fourier polynomial with the same degree as the evaluator.
func (e *Evaluator) NewFourierPoly() FourierPoly {
	return FourierPoly{Coeffs: make([]float64, e.Parameters.degree)}
}
