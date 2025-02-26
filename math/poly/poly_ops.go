package poly

import "github.com/sp301415/carousel/math/vec"

// AddPoly returns p0 + p1.
func (e *Evaluator) AddPoly(p0, p1 Poly) Poly {
	pOut := e.NewPoly()
	e.AddPolyAssign(p0, p1, pOut)
	return pOut
}

// AddPolyAssign computes pOut = p0 + p1.
func (e *Evaluator) AddPolyAssign(p0, p1, pOut Poly) {
	vec.AddAssign(p0.Coeffs, p1.Coeffs, pOut.Coeffs)
}

// SubPoly returns p0 - p1.
func (e *Evaluator) SubPoly(p0, p1 Poly) Poly {
	pOut := e.NewPoly()
	e.SubPolyAssign(p0, p1, pOut)
	return pOut
}

// SubPolyAssign computes pOut = p0 - p1.
func (e *Evaluator) SubPolyAssign(p0, p1, pOut Poly) {
	vec.SubAssign(p0.Coeffs, p1.Coeffs, pOut.Coeffs)
}

// NegPoly returns pOut = -p0.
func (e *Evaluator) NegPoly(p0 Poly) Poly {
	pOut := e.NewPoly()
	e.NegPolyAssign(p0, pOut)
	return pOut
}

// NegPolyAssign computes pOut = -p0.
func (e *Evaluator) NegPolyAssign(p0, pOut Poly) {
	vec.NegAssign(p0.Coeffs, pOut.Coeffs)
}

// ScalarMulPoly returns c * p0.
func (e *Evaluator) ScalarMulPoly(p0 Poly, c uint64) Poly {
	pOut := e.NewPoly()
	e.ScalarMulPolyAssign(p0, c, pOut)
	return pOut
}

// ScalarMulPolyAssign computes pOut = c * p0.
func (e *Evaluator) ScalarMulPolyAssign(p0 Poly, c uint64, pOut Poly) {
	vec.ScalarMulAssign(p0.Coeffs, c, pOut.Coeffs)
}

// ScalarMulAddPolyAssign computes pOut += c * p0.
func (e *Evaluator) ScalarMulAddPolyAssign(p0 Poly, c uint64, pOut Poly) {
	vec.ScalarMulAddAssign(p0.Coeffs, c, pOut.Coeffs)
}

// ScalarMulSubPolyAssign computes pOut -= c * p0.
func (e *Evaluator) ScalarMulSubPolyAssign(p0 Poly, c uint64, pOut Poly) {
	vec.ScalarMulSubAssign(p0.Coeffs, c, pOut.Coeffs)
}

// FourierPolyMulPoly returns p0 * fp.
func (e *Evaluator) FourierPolyMulPoly(p0 Poly, fp FourierPoly) Poly {
	pOut := e.NewPoly()
	e.FourierPolyMulPolyAssign(p0, fp, pOut)
	return pOut
}

// FourierPolyMulPolyAssign computes pOut = p0 * fp.
func (e *Evaluator) FourierPolyMulPolyAssign(p0 Poly, fp FourierPoly, pOut Poly) {
	e.ToFourierPolyAssign(p0, e.buffer.fp)
	e.MulFourierPolyAssign(e.buffer.fp, fp, e.buffer.fp)
	e.ToPolyAssignUnsafe(e.buffer.fp, pOut)
}

// FourierPolyMulAddPolyAssign computes pOut += p0 * fp.
func (e *Evaluator) FourierPolyMulAddPolyAssign(p0 Poly, fp FourierPoly, pOut Poly) {
	e.ToFourierPolyAssign(p0, e.buffer.fp)
	e.MulFourierPolyAssign(e.buffer.fp, fp, e.buffer.fp)
	e.ToPolyAddAssignUnsafe(e.buffer.fp, pOut)
}

// FourierPolyMulSubPolyAssign computes pOut -= p0 * fp.
func (e *Evaluator) FourierPolyMulSubPolyAssign(p0 Poly, fp FourierPoly, pOut Poly) {
	e.ToFourierPolyAssign(p0, e.buffer.fp)
	e.MulFourierPolyAssign(e.buffer.fp, fp, e.buffer.fp)
	e.ToPolySubAssignUnsafe(e.buffer.fp, pOut)
}

// PermutePoly returns p0(X^d).
func (e *Evaluator) PermutePoly(p0 Poly, d int) Poly {
	pOut := e.NewPoly()
	e.PermutePolyAssign(p0, d, pOut)
	return pOut
}

// PermutePolyAssign computes pOut = p0(X^d).
//
// p0 and pOut should not overlap. For inplace permutation,
// use [*Evaluator.PermutePolyInPlace].
func (e *Evaluator) PermutePolyAssign(p0 Poly, d int, pOut Poly) {
	vec.RotateAssign(p0.Coeffs, d, pOut.Coeffs)
}

// PermutePolyInPlace computes p0 = p0(X^d).
//
// Panics when d is not odd.
// This is because the permutation is not bijective when d is even.
func (e *Evaluator) PermutePolyInPlace(p0 Poly, d int) {
	vec.RotateInPlace(p0.Coeffs, d)
}
