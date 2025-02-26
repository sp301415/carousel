package poly

import "github.com/sp301415/carousel/math/vec"

// ToFourierPoly transforms Poly to FourierPoly.
func (e *Evaluator) ToFourierPoly(p Poly) FourierPoly {
	fpOut := NewFourierPoly(e.Parameters.degree)
	e.ToFourierPolyAssign(p, fpOut)
	return fpOut
}

// ToFourierPolyAssign transforms Poly to FourierPoly and writes it to fpOut.
func (e *Evaluator) ToFourierPolyAssign(p Poly, fpOut FourierPoly) {
	convertUint64ToFloat64Assign(p.Coeffs, fpOut.Coeffs)
	rfftInPlace(fpOut.Coeffs, e.twRealFFT)
	convolveAssign(fpOut.Coeffs, e.twFFT, fpOut.Coeffs)
	rifftInPlace(fpOut.Coeffs, e.twRealInvFFT)
}

// ToFourierPolyAddAssign transforms Poly to FourierPoly and adds it to fpOut.
func (e *Evaluator) ToFourierPolyAddAssign(p Poly, fpOut FourierPoly) {
	convertUint64ToFloat64Assign(p.Coeffs, e.buffer.fp.Coeffs)
	rfftInPlace(e.buffer.fp.Coeffs, e.twRealFFT)
	convolveAssign(e.buffer.fp.Coeffs, e.twFFT, e.buffer.fp.Coeffs)
	rifftInPlace(e.buffer.fp.Coeffs, e.twRealInvFFT)
	addFloat64Assign(fpOut.Coeffs, e.buffer.fp.Coeffs, fpOut.Coeffs)
}

// ToFourierPolySubAssign transforms Poly to FourierPoly and subtracts it from fpOut.
func (e *Evaluator) ToFourierPolySubAssign(p Poly, fpOut FourierPoly) {
	convertUint64ToFloat64Assign(p.Coeffs, e.buffer.fp.Coeffs)
	rfftInPlace(e.buffer.fp.Coeffs, e.twRealFFT)
	convolveAssign(e.buffer.fp.Coeffs, e.twFFT, e.buffer.fp.Coeffs)
	rifftInPlace(e.buffer.fp.Coeffs, e.twRealInvFFT)
	subFloat64Assign(fpOut.Coeffs, e.buffer.fp.Coeffs, fpOut.Coeffs)
}

// ToPoly transforms FourierPoly to Poly.
func (e *Evaluator) ToPoly(fp FourierPoly) Poly {
	pOut := NewPoly(e.Parameters.degree)
	e.ToPolyAssign(fp, pOut)
	return pOut
}

// ToPolyAssign transforms FourierPoly to Poly and writes it to pOut.
func (e *Evaluator) ToPolyAssign(fp FourierPoly, pOut Poly) {
	e.buffer.fpInv.CopyFrom(fp)
	rfftInPlace(e.buffer.fpInv.Coeffs, e.twRealFFT)
	convolveAssign(e.buffer.fpInv.Coeffs, e.twInvFFT, e.buffer.fpInv.Coeffs)
	rifftInPlace(e.buffer.fpInv.Coeffs, e.twRealInvFFT)
	floatModQInPlace(e.buffer.fpInv.Coeffs)
	convertFloat64ToUint64Assign(e.buffer.fpInv.Coeffs, pOut.Coeffs)
}

// ToPolyAddAssign transforms FourierPoly to Poly and adds it to pOut.
func (e *Evaluator) ToPolyAddAssign(fp FourierPoly, pOut Poly) {
	e.buffer.fpInv.CopyFrom(fp)
	rfftInPlace(e.buffer.fpInv.Coeffs, e.twRealFFT)
	convolveAssign(e.buffer.fpInv.Coeffs, e.twInvFFT, e.buffer.fpInv.Coeffs)
	rifftInPlace(e.buffer.fpInv.Coeffs, e.twRealInvFFT)
	floatModQInPlace(e.buffer.fpInv.Coeffs)
	convertFloat64ToUint64Assign(e.buffer.fpInv.Coeffs, e.buffer.pInv.Coeffs)
	vec.AddAssign(pOut.Coeffs, e.buffer.pInv.Coeffs, e.buffer.pInv.Coeffs)
}

// ToPolySubAssign transforms FourierPoly to Poly and subtracts it from pOut.
func (e *Evaluator) ToPolySubAssign(fp FourierPoly, pOut Poly) {
	e.buffer.fpInv.CopyFrom(fp)
	rfftInPlace(e.buffer.fpInv.Coeffs, e.twRealFFT)
	convolveAssign(e.buffer.fpInv.Coeffs, e.twInvFFT, e.buffer.fpInv.Coeffs)
	rifftInPlace(e.buffer.fpInv.Coeffs, e.twRealInvFFT)
	floatModQInPlace(e.buffer.fpInv.Coeffs)
	convertFloat64ToUint64Assign(e.buffer.fpInv.Coeffs, e.buffer.pInv.Coeffs)
	vec.SubAssign(pOut.Coeffs, e.buffer.pInv.Coeffs, e.buffer.pInv.Coeffs)
}

// ToPolyAssignUnsafe transforms FourierPoly to Poly and writes it to pOut.
//
// This method is slightly faster than [*Evaluator.ToPolyAssign], but it modifies fp directly.
// Use it only if you don't need fp after this method (e.g. fp is a buffer).
func (e *Evaluator) ToPolyAssignUnsafe(fp FourierPoly, pOut Poly) {
	rfftInPlace(fp.Coeffs, e.twRealFFT)
	convolveAssign(fp.Coeffs, e.twInvFFT, fp.Coeffs)
	rifftInPlace(fp.Coeffs, e.twRealInvFFT)
	floatModQInPlace(fp.Coeffs)
	convertFloat64ToUint64Assign(fp.Coeffs, pOut.Coeffs)
}

// ToPolyAddAssignUnsafe transforms FourierPoly to Poly and adds it to pOut.
//
// This method is slightly faster than [*Evaluator.ToPolyAddAssign], but it modifies fp directly.
// Use it only if you don't need fp after this method (e.g. fp is a buffer).
func (e *Evaluator) ToPolyAddAssignUnsafe(fp FourierPoly, pOut Poly) {
	rfftInPlace(fp.Coeffs, e.twRealFFT)
	convolveAssign(fp.Coeffs, e.twInvFFT, fp.Coeffs)
	rifftInPlace(fp.Coeffs, e.twRealInvFFT)
	floatModQInPlace(fp.Coeffs)
	convertFloat64ToUint64Assign(fp.Coeffs, e.buffer.pInv.Coeffs)
	vec.AddAssign(pOut.Coeffs, e.buffer.pInv.Coeffs, e.buffer.pInv.Coeffs)
}

// ToPolySubAssignUnsafe transforms FourierPoly to Poly and subtracts it from pOut.
//
// This method is slightly faster than [*Evaluator.ToPolySubAssign], but it modifies fp directly.
// Use it only if you don't need fp after this method (e.g. fp is a buffer).
func (e *Evaluator) ToPolySubAssignUnsafe(fp FourierPoly, pOut Poly) {
	rfftInPlace(fp.Coeffs, e.twRealFFT)
	convolveAssign(fp.Coeffs, e.twInvFFT, fp.Coeffs)
	rifftInPlace(fp.Coeffs, e.twRealInvFFT)
	floatModQInPlace(fp.Coeffs)
	convertFloat64ToUint64Assign(fp.Coeffs, e.buffer.pInv.Coeffs)
	vec.SubAssign(pOut.Coeffs, e.buffer.pInv.Coeffs, e.buffer.pInv.Coeffs)
}
