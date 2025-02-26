package poly

import (
	"math"
	"math/cmplx"

	"github.com/sp301415/carousel/math/num"
	"github.com/sp301415/carousel/math/vec"
)

const (
	// MinDegree is the minimum degree of polynomial that Evaluator can handle.
	// Currently, this is set to 2^4, because AVX2 implementation of FFT and inverse FFT
	// handles first/last two loops separately.
	MinDegree = 1 << 4
)

// Evaluator computes polynomial operation over a subring.
type Evaluator struct {
	Parameters EvaluatorParameters

	// twRealFFT is the twiddle factor for real FFT.
	twRealFFT []complex128
	// twRealInvFFT is the twiddle factor for inverse real FFT.
	twRealInvFFT []complex128

	// twFFT is the factor for FFT over the subring.
	twFFT []float64
	// twInvFFT is the factor for inverse FFT over the subring.
	twInvFFT []float64

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
}

// NewEvaluatorWithParameters creates a new Evaluator with the given parameters.
func NewEvaluator(params EvaluatorParameters) *Evaluator {
	N := params.degree

	switch {
	case !num.IsPowerOfTwo(N):
		panic("NewEvaluator: degree must be a power of two")
	case N < MinDegree:
		panic("NewEvaluator: degree smaller than MinDegree")
	}

	twRealFFT, twRealInvFFT := genTwiddleFactorsRealFFT(N)

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
	rfftInPlace(twFFT, twRealFFT)
	rfftInPlace(twInvFFT, twRealFFT)

	for i := 0; i < N; i++ {
		twFFT[i] /= float64(N / 2)
		twInvFFT[i] /= float64(N / 2)
	}

	return &Evaluator{
		Parameters: params,

		twRealFFT:    twRealFFT,
		twRealInvFFT: twRealInvFFT,

		twFFT:    twFFT,
		twInvFFT: twInvFFT,

		buffer: newEvaluationBuffer(N),
	}
}

func genTwiddleFactorsRealFFT(N int) (twRealFFT, twRealInvFFT []complex128) {
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
	twRealFFT = make([]complex128, N/2)
	w = 0
	for m := 1; m <= N/4; m <<= 1 {
		for i := 0; i < m; i++ {
			twRealFFT[w] = twRealFFTRef[i]
			w++
		}
	}

	twRealInvFFT = make([]complex128, N/2)
	w = 0
	for m := N / 4; m >= 1; m >>= 1 {
		for i := 0; i < m; i++ {
			twRealInvFFT[w] = twRealInvFFTRef[i]
			w++
		}
	}

	return twRealFFT, twRealInvFFT
}

// newEvaluationBuffer creates a new evaluationBuffer.
func newEvaluationBuffer(N int) evaluationBuffer {
	return evaluationBuffer{
		pOut:  NewPoly(N),
		fpOut: NewFourierPoly(N),

		fp:    NewFourierPoly(N),
		fpInv: NewFourierPoly(N),
		pInv:  NewPoly(N),
	}
}

// ShallowCopy returns a shallow copy of this Evaluator.
// Returned Evaluator is safe for concurrent use.
func (e *Evaluator) ShallowCopy() *Evaluator {
	return &Evaluator{
		Parameters: e.Parameters,

		twRealFFT:    e.twRealFFT,
		twRealInvFFT: e.twRealInvFFT,

		twFFT:    e.twFFT,
		twInvFFT: e.twInvFFT,

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
