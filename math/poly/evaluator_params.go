package poly

import (
	"math"
	"math/big"
	"math/rand"

	"github.com/sp301415/carousel/math/num"
	"github.com/sp301415/carousel/math/vec"
)

const (
	// OrderBound is the bound of the order of the baes prime.
	// Given a degree N, we search the possible cyclotomic degree M
	// up to OrderBound*N + 1.
	OrderBound = 128
)

var (
	smallPrimes = []int{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37,
		41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83,
		89, 97, 101, 103, 107, 109, 113, 127,
	}
)

// EvaluatorParameters is the parameters for initializing [Evaluator].
// We find suitable parameters for initializing a degree N subring of
// Z[X] / Phi_M(X), invariant under X -> X^p.
type EvaluatorParameters struct {
	// cyclotomicDegree is the degree of the ambient cyclotomic polynomial.
	// Must be a prime.
	cyclotomicDegree int

	// degree is the degree of the subring used.
	// Must be a power of two.
	degree int

	// basePrime is the base prime used for the subring.
	// For [NewEvaluatorParametersForPacking], basePrime is always 2.
	basePrime uint64

	// order is the multiplicative order of basePrime modulo cyclotomicDegree.
	order int

	// generator is the generator of multplicative group modulo cyclotomicDegree.
	generator uint64

	// resolution are the resolution of unity.
	// Empty if the parameters are for coefficient encoding.
	resolution []uint64
}

// NewEvaluatorParameters creates a new EvaluatorParameters.
// Panics if the parameters are invalid or not found.
func NewEvaluatorParameters(N int) EvaluatorParameters {
	switch {
	case N < MinDegree:
		panic("NewEvaluatorParameters: N smaller than MinDegree")
	case !num.IsPowerOfTwo(N):
		panic("NewEvaluatorParameters: N must be a power of two")
	}

	for k := 2; k <= OrderBound; k++ {
		M := k*N + 1
		if !big.NewInt(int64(M)).ProbablyPrime(0) {
			continue
		}

		kFactor := []int{}
		for _, fac := range smallPrimes {
			if k%fac == 0 {
				kFactor = append(kFactor, fac)
			}
		}

		var g uint64
		for g = 2; g < uint64(M); g++ {
			isFound := true
			for _, fac := range kFactor {
				isFound = isFound && num.ModExp(g, uint64((M-1)/fac), uint64(M)) != 1
			}
			if isFound {
				break
			}
		}

		p := num.ModExp(g, uint64(N), uint64(M))

		return EvaluatorParameters{
			cyclotomicDegree: M,
			degree:           N,
			basePrime:        p,
			order:            k,
			generator:        g,
		}
	}

	panic("NewEvaluatorParameters: failed to find suitable cyclotomicDegree")
}

// NewEvaluatorParametersForPacking creats a new EvaluatorParameters for slot packing,
// with the plaintext modulus 2^bits.
// Panics if the parameters are invalid or not found.
func NewEvaluatorParametersForPacking(N, bits int) EvaluatorParameters {
	switch {
	case N < MinDegree:
		panic("NewEvaluatorParameters: N smaller than MinDegree")
	case !num.IsPowerOfTwo(N):
		panic("NewEvaluatorParameters: N must be a power of two")
	case bits <= 0:
		panic("NewEvaluatorParameters: bits must be positive")
	}

	for k := 2; k <= OrderBound; k++ {
		M := k*N + 1
		if !big.NewInt(int64(M)).ProbablyPrime(0) {
			continue
		}

		kFactor := []int{}
		for _, fac := range smallPrimes {
			if k%fac == 0 {
				kFactor = append(kFactor, fac)
			}
		}

		var g uint64
		for g = 2; g < uint64(M); g++ {
			isFound := true
			for _, fac := range kFactor {
				isFound = isFound && num.ModExp(g, uint64((M-1)/fac), uint64(M)) != 1
			}
			if isFound {
				break
			}
		}

		order := 0
		pPow := uint64(1)
		for {
			pPow = (pPow * 2) % uint64(M)
			order++
			if pPow == 1 {
				break
			}
		}
		if k != order {
			continue
		}

		resolution := genResolution(M, bits, k, g)

		return EvaluatorParameters{
			cyclotomicDegree: M,
			degree:           N,
			basePrime:        2,
			order:            k,
			generator:        g,
			resolution:       resolution,
		}
	}

	panic("NewEvaluatorParameters: failed to find suitable cyclotomicDegree")
}

// genResolution computes the resolution of unities.
func genResolution(M, bits, order int, gen uint64) []uint64 {
	Gptd := newZ2XQuoRing(bits, loadConaway(order))

	bigOrder := big.NewInt(0).Lsh(big.NewInt(1), uint(order))
	bigOrder.Sub(bigOrder, big.NewInt(1))
	bigOrder.Lsh(bigOrder, uint((bits-1)*order))

	bigOrderDivM := big.NewInt(0).Div(bigOrder, big.NewInt(int64(M)))

	rSrc := rand.New(rand.NewSource(0))
	var g z2xPoly
	for {
		gg := Gptd.newPoly()
		for i := 0; i < order; i++ {
			gg[i] = uint64(rSrc.Intn(1 << bits))
		}

		ggPow := Gptd.expBig(gg, bigOrder)
		if ggPow.isOne() {
			g = Gptd.expBig(gg, bigOrderDivM)
			if !g.isOne() {
				break
			}
		}
	}

	factorLift := make([]z2xPoly, order+1)
	for i := 0; i < order+1; i++ {
		factorLift[i] = Gptd.newPoly()
	}
	factorLift[0][0] = 1

	for i := 0; i < order; i++ {
		e := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(int64(i)), big.NewInt(int64(M)))
		gPow := Gptd.expBig(g, e)

		vec.CopyAssign(factorLift[i], factorLift[i+1])
		for j := i; j > 0; j-- {
			factorLift[j] = Gptd.sub(factorLift[j-1], Gptd.mul(factorLift[j], gPow))
		}
		factorLift[0] = Gptd.neg(Gptd.mul(factorLift[0], gPow))
	}

	factor := make(z2xPoly, order+1)
	for i := 0; i < order+1; i++ {
		factor[i] = factorLift[i][0]
	}

	comp, rem := make(z2xPoly, M-order), make(z2xPoly, M)
	vec.Fill(rem, 1)
	for i := M - 1; i >= order; i-- {
		if rem[i] == 0 {
			continue
		}

		c := rem[i]
		for j := 0; j <= order; j++ {
			rem[i-j] -= c * factor[order-j]
			rem[i-j] &= Gptd.mask
		}
		comp[i-order] += c & Gptd.mask
	}
	if !rem.isZero() {
		panic("GenResolution: failed to compute resolution, this is a bug")
	}

	factorRing := newZ2XQuoRing(bits, factor)
	compFactorRing := factorRing.reduce(comp)
	compInv := factorRing.expBig(compFactorRing, big.NewInt(0).Sub(bigOrder, big.NewInt(1)))

	fftDegree := 1 << (num.Log2(M) + 1)
	compFFTUint64 := make([]uint64, fftDegree)
	compInvFFTUint64 := make([]uint64, fftDegree)
	vec.CopyAssign(comp, compFFTUint64)
	vec.CopyAssign(compInv, compInvFFTUint64)

	compFFT := make([]float64, fftDegree)
	compInvFFT := make([]float64, fftDegree)
	convertUint64ToFloat64Assign(compFFTUint64, compFFT)
	convertUint64ToFloat64Assign(compInvFFTUint64, compInvFFT)

	twRealFFT, twRealInvFFT := genTwiddleFactorsRealFFT(fftDegree)
	rfftInPlace(compFFT, twRealFFT)
	rfftInPlace(compInvFFT, twRealFFT)

	crtBasisFFT := make([]float64, fftDegree)
	convolveAssign(compFFT, compInvFFT, crtBasisFFT)
	rifftInPlace(crtBasisFFT, twRealInvFFT)
	for i := 0; i < fftDegree; i++ {
		crtBasisFFT[i] /= float64(fftDegree) / 2
	}

	crtBasis := make([]uint64, fftDegree)
	for i := 0; i < fftDegree; i++ {
		crtBasis[i] = uint64(int64(math.Round(crtBasisFFT[i])))
	}

	N := (M - 1) / order
	resolution := make([]uint64, N)
	idx := uint64(1)
	for i := 0; i < N; i++ {
		resolution[i] = (crtBasis[idx] - crtBasis[0]) & ((1 << bits) - 1)
		idx = (idx * gen) % uint64(M)
	}

	rotIdx := 0
	for i := 0; i < N; i++ {
		if resolution[i]&1 == 1 {
			rotIdx = i
			break
		}
	}
	vec.RotateInPlace(resolution, -rotIdx)

	return resolution
}

// CyclotomicDegree returns the degree of the ambient cyclotomic polynomial.
func (p EvaluatorParameters) CyclotomicDegree() int {
	return p.cyclotomicDegree
}

// Degree returns the degree of the subring used.
func (p EvaluatorParameters) Degree() int {
	return p.degree
}

// Order returns the multiplicative order of basePrime modulo cyclotomicDegree.
func (p EvaluatorParameters) Order() int {
	return p.order
}

// Generator returns the generator of multplicative group modulo cyclotomicDegree.
func (p EvaluatorParameters) Generator() uint64 {
	return p.generator
}

// Resolution returns the resolution of unities.
func (p EvaluatorParameters) Resolution() []uint64 {
	return p.resolution
}

// IsPackable returns true if the parameters are for slot packing.
func (p EvaluatorParameters) IsPackable() bool {
	return len(p.resolution) != 0
}
