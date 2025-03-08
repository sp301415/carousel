package carousel

import (
	"github.com/sp301415/carousel/math/num"
	"github.com/sp301415/carousel/math/poly"
)

// Decomposer decomposes a scalar or a polynomial
// according to the gadget parameters.
//
// Decomposer is not safe for concurrent use.
// Use [*Decomposer.ShallowCopy] to get a safe copy.
type Decomposer struct {
	// PolyEvaluator is a PolyEvaluator for this Decomposer.
	PolyEvaluator *poly.Evaluator

	buffer decompositionBuffer
}

// decompositionBuffer is a buffer for Decomposer.
type decompositionBuffer struct {
	// scalarDecomposed is the scalarDecomposed scalar.
	scalarDecomposed []uint64
	// polyDecomposed is the decomposed polynomial.
	polyDecomposed []poly.Poly
	// polyFourierDecomposed is the decomposed polynomial in Fourier domain.
	polyFourierDecomposed []poly.FourierPoly
}

// NewDecomposer creates a new Decomposer.
func NewDecomposer(params poly.EvaluatorParameters) *Decomposer {
	return &Decomposer{
		PolyEvaluator: poly.NewEvaluator(params),

		buffer: decompositionBuffer{},
	}
}

// ShallowCopy returns a shallow copy of this Decomposer.
// Returned Decomposer is safe for concurrent use.
func (d *Decomposer) ShallowCopy() *Decomposer {
	decompositionBufferCopy := decompositionBuffer{
		scalarDecomposed:      make([]uint64, len(d.buffer.scalarDecomposed)),
		polyDecomposed:        make([]poly.Poly, len(d.buffer.polyDecomposed)),
		polyFourierDecomposed: make([]poly.FourierPoly, len(d.buffer.polyFourierDecomposed)),
	}

	for i := range decompositionBufferCopy.polyDecomposed {
		decompositionBufferCopy.polyDecomposed[i] = d.PolyEvaluator.NewPoly()
	}
	for i := range decompositionBufferCopy.polyFourierDecomposed {
		decompositionBufferCopy.polyFourierDecomposed[i] = d.PolyEvaluator.NewFourierPoly()
	}

	return &Decomposer{
		PolyEvaluator: d.PolyEvaluator.ShallowCopy(),

		buffer: decompositionBufferCopy,
	}
}

// ScalarDecomposedBuffer returns a internal buffer for scalar decomposition with respect to gadgetParams.
//
// You can also set the length of the internal buffer by calling this function and ignoring the return value.
func (d *Decomposer) ScalarDecomposedBuffer(gadgetParams GadgetParameters) []uint64 {
	if len(d.buffer.scalarDecomposed) >= gadgetParams.level {
		return d.buffer.scalarDecomposed[:gadgetParams.level]
	}

	oldLen := len(d.buffer.scalarDecomposed)
	d.buffer.scalarDecomposed = append(d.buffer.scalarDecomposed, make([]uint64, gadgetParams.level-oldLen)...)
	return d.buffer.scalarDecomposed
}

// PolyDecomposedBuffer returns a internal buffer for polynomial decomposition with respect to gadgetParams.
//
// You can also set the length of the internal buffer by calling this function and ignoring the return value.
func (d *Decomposer) PolyDecomposedBuffer(gadgetParams GadgetParameters) []poly.Poly {
	if len(d.buffer.polyDecomposed) >= gadgetParams.level {
		return d.buffer.polyDecomposed[:gadgetParams.level]
	}

	oldLen := len(d.buffer.polyDecomposed)
	d.buffer.polyDecomposed = append(d.buffer.polyDecomposed, make([]poly.Poly, gadgetParams.level-oldLen)...)
	for i := oldLen; i < gadgetParams.level; i++ {
		d.buffer.polyDecomposed[i] = d.PolyEvaluator.NewPoly()
	}
	return d.buffer.polyDecomposed
}

// PolyFourierDecomposedBuffer returns a internal buffer for polynomial decomposition in Fourier domain with respect to gadgetParams.
//
// You can also set the length of the internal buffer by calling this function and ignoring the return value.
func (d *Decomposer) PolyFourierDecomposedBuffer(gadgetParams GadgetParameters) []poly.FourierPoly {
	if len(d.buffer.polyFourierDecomposed) >= gadgetParams.level {
		return d.buffer.polyFourierDecomposed[:gadgetParams.level]
	}

	oldLen := len(d.buffer.polyFourierDecomposed)
	d.buffer.polyFourierDecomposed = append(d.buffer.polyFourierDecomposed, make([]poly.FourierPoly, gadgetParams.level-oldLen)...)
	for i := oldLen; i < gadgetParams.level; i++ {
		d.buffer.polyFourierDecomposed[i] = d.PolyEvaluator.NewFourierPoly()
	}
	return d.buffer.polyFourierDecomposed
}

// DecomposeScalar decomposes x with respect to gadgetParams.
func (d *Decomposer) DecomposeScalar(x uint64, gadgetParams GadgetParameters) []uint64 {
	decomposedOut := make([]uint64, gadgetParams.level)
	d.DecomposeScalarAssign(x, gadgetParams, decomposedOut)
	return decomposedOut
}

// DecomposeScalarAssign decomposes x with respect to gadgetParams and writes it to decomposedOut.
func (d *Decomposer) DecomposeScalarAssign(x uint64, gadgetParams GadgetParameters, decomposedOut []uint64) {
	u := num.DivRoundBits(x, gadgetParams.LogLastBaseQ())
	for i := gadgetParams.level - 1; i >= 1; i-- {
		decomposedOut[i] = u & (gadgetParams.base - 1)
		u >>= gadgetParams.logBase
		u += decomposedOut[i] >> (gadgetParams.logBase - 1)
		decomposedOut[i] -= (decomposedOut[i] & (gadgetParams.base >> 1)) << 1
	}
	decomposedOut[0] = u & (gadgetParams.base - 1)
	decomposedOut[0] -= (decomposedOut[0] & (gadgetParams.base >> 1)) << 1
}

// DecomposePoly decomposes p with respect to gadgetParams.
func (d *Decomposer) DecomposePoly(p poly.Poly, gadgetParams GadgetParameters) []poly.Poly {
	decomposedOut := make([]poly.Poly, gadgetParams.level)
	for i := 0; i < gadgetParams.level; i++ {
		decomposedOut[i] = poly.NewPoly(p.Degree())
	}
	d.DecomposePolyAssign(p, gadgetParams, decomposedOut)
	return decomposedOut
}

// DecomposePolyAssign decomposes p with respect to gadgetParams and writes it to decomposedOut.
func (d *Decomposer) DecomposePolyAssign(p poly.Poly, gadgetParams GadgetParameters, decomposedOut []poly.Poly) {
	decomposePolyAssign(p, gadgetParams, decomposedOut)
}

// FourierDecomposePoly decomposes p with respect to gadgetParams in Fourier domain.
func (d *Decomposer) FourierDecomposePoly(p poly.Poly, gadgetParams GadgetParameters) []poly.FourierPoly {
	decomposedOut := make([]poly.FourierPoly, gadgetParams.level)
	for i := 0; i < gadgetParams.level; i++ {
		decomposedOut[i] = d.PolyEvaluator.NewFourierPoly()
	}
	d.FourierDecomposePolyAssign(p, gadgetParams, decomposedOut)
	return decomposedOut
}

// FourierDecomposePolyAssign decomposes p with respect to gadgetParams in Fourier domain and writes it to decomposedOut.
func (d *Decomposer) FourierDecomposePolyAssign(p poly.Poly, gadgetParams GadgetParameters, decomposedOut []poly.FourierPoly) {
	polyDecomposed := d.PolyDecomposedBuffer(gadgetParams)
	decomposePolyAssign(p, gadgetParams, polyDecomposed)
	for i := 0; i < gadgetParams.level; i++ {
		d.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[i], decomposedOut[i])
	}
}

// RecomposeScalar recomposes decomposed with respect to gadgetParams.
func (d *Decomposer) RecomposeScalar(decomposed []uint64, gadgetParams GadgetParameters) uint64 {
	var x uint64
	for i := 0; i < gadgetParams.level; i++ {
		x += decomposed[i] << gadgetParams.LogBaseQ(i)
	}
	return x
}

// RecomposePoly recomposes decomposed with respect to gadgetParams.
func (d *Decomposer) RecomposePoly(decomposed []poly.Poly, gadgetParams GadgetParameters) poly.Poly {
	pOut := poly.NewPoly(decomposed[0].Degree())
	d.RecomposePolyAssign(decomposed, gadgetParams, pOut)
	return pOut
}

// RecomposePolyAssign recomposes decomposed with respect to gadgetParams and writes it to pOut.
func (d *Decomposer) RecomposePolyAssign(decomposed []poly.Poly, gadgetParams GadgetParameters, pOut poly.Poly) {
	for i := 0; i < pOut.Degree(); i++ {
		pOut.Coeffs[i] = 0
		for j := 0; j < gadgetParams.level; j++ {
			pOut.Coeffs[i] += decomposed[j].Coeffs[i] << gadgetParams.LogBaseQ(j)
		}
	}
}
