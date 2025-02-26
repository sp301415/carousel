package poly

import (
	"math"
	"math/cmplx"
	"math/rand"
	"testing"

	"github.com/sp301415/carousel/math/vec"
	"github.com/stretchr/testify/assert"
)

func fftInPlaceRef(coeffs, tw []complex128) {
	N := len(coeffs)

	t := N
	for m := 1; m <= N/2; m <<= 1 {
		t >>= 1
		for i := 0; i < m; i++ {
			j1 := i * t << 1
			j2 := j1 + t
			for j := j1; j < j2; j++ {
				U, V := coeffs[j], coeffs[j+t]*tw[i]
				coeffs[j], coeffs[j+t] = U+V, U-V
			}
		}
	}
}

func invFFTInPlaceRef(coeffs, twInv []complex128) {
	N := len(coeffs)

	t := 1
	for m := N / 2; m >= 1; m >>= 1 {
		for i := 0; i < m; i++ {
			j1 := i * t << 1
			j2 := j1 + t
			for j := j1; j < j2; j++ {
				U, V := coeffs[j], coeffs[j+t]
				coeffs[j], coeffs[j+t] = U+V, (U-V)*twInv[i]
			}
		}
		t <<= 1
	}
}

func TestRFFTAssembly(t *testing.T) {
	N := 32
	eps := 1e-10
	evParams := NewEvaluatorParameters(N)
	ev := NewEvaluator(evParams)

	r := rand.New(rand.NewSource(0))

	fp0 := make([]float64, N)
	fp1 := make([]float64, N)
	fpOut := make([]float64, N)

	fpCmplx0 := make([]complex128, N)
	fpCmplx1 := make([]complex128, N)
	fpCmplxOut := make([]complex128, N)
	fpCmplxOutReal := make([]float64, N)

	for i := 0; i < N; i++ {
		fp0[i] = r.Float64()
		fp1[i] = r.Float64()

		fpCmplx0[i] = complex(fp0[i], 0)
		fpCmplx1[i] = complex(fp1[i], 0)
	}

	twCmplxFFT := make([]complex128, N/2)
	twCmplxInvFFT := make([]complex128, N/2)
	for i := 0; i < N/2; i++ {
		e := -2 * math.Pi * float64(i) / float64(N)
		twCmplxFFT[i] = cmplx.Exp(complex(0, e))
		twCmplxInvFFT[i] = cmplx.Exp(-complex(0, e))
	}
	vec.BitReverseInPlace(twCmplxFFT)
	vec.BitReverseInPlace(twCmplxInvFFT)

	rfftInPlace(fp0, ev.twRealFFT)
	rfftInPlace(fp1, ev.twRealFFT)
	convolveAssign(fp0, fp1, fpOut)
	rifftInPlace(fpOut, ev.twRealInvFFT)
	for i := 0; i < N; i++ {
		fpOut[i] /= float64(N / 2)
	}

	fftInPlaceRef(fpCmplx0, twCmplxFFT)
	fftInPlaceRef(fpCmplx1, twCmplxFFT)
	for i := 0; i < N; i++ {
		fpCmplxOut[i] = fpCmplx0[i] * fpCmplx1[i]
	}
	invFFTInPlaceRef(fpCmplxOut, twCmplxInvFFT)
	for i := 0; i < N; i++ {
		fpCmplxOutReal[i] = real(fpCmplxOut[i]) / float64(N)
	}

	assert.InEpsilonSlice(t, fpOut, fpCmplxOutReal, eps)
}
