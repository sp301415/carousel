package poly

import "github.com/sp301415/carousel/math/vec"

// AddFourierPoly returns fp0 + fp1.
func (e *Evaluator) AddFourierPoly(fp0, fp1 FourierPoly) FourierPoly {
	fpOut := e.NewFourierPoly()
	e.AddFourierPolyAssign(fp0, fp1, fpOut)
	return fpOut
}

// AddFourierPolyAssign computes fpOut = fp0 + fp1.
func (e *Evaluator) AddFourierPolyAssign(fp0, fp1, fpOut FourierPoly) {
	addFloat64Assign(fp0.Coeffs, fp1.Coeffs, fpOut.Coeffs)
}

// SubFourierPoly returns fp0 - fp1.
func (e *Evaluator) SubFourierPoly(fp0, fp1 FourierPoly) FourierPoly {
	fpOut := e.NewFourierPoly()
	e.SubFourierPolyAssign(fp0, fp1, fpOut)
	return fpOut
}

// SubFourierPolyAssign computes fpOut = fp0 - fp1.
func (e *Evaluator) SubFourierPolyAssign(fp0, fp1, fpOut FourierPoly) {
	subFloat64Assign(fp0.Coeffs, fp1.Coeffs, fpOut.Coeffs)
}

// NegFourierPoly returns -fp0.
func (e *Evaluator) NegFourierPoly(fp0 FourierPoly) FourierPoly {
	fpOut := e.NewFourierPoly()
	e.NegFourierPolyAssign(fp0, fpOut)
	return fpOut
}

// NegFourierPolyAssign computes fpOut = -fp0.
func (e *Evaluator) NegFourierPolyAssign(fp0, fpOut FourierPoly) {
	negFloat64Assign(fp0.Coeffs, fpOut.Coeffs)
}

// ScalarMulFourierPoly returns c * fp0.
func (e *Evaluator) ScalarMulFourierPoly(fp0 FourierPoly, c float64) FourierPoly {
	fpOut := e.NewFourierPoly()
	e.ScalarMulFourierPolyAssign(fp0, c, fpOut)
	return fpOut
}

// ScalarMulFourierPolyAssign computes fpOut = c * fp0.
func (e *Evaluator) ScalarMulFourierPolyAssign(fp0 FourierPoly, c float64, fpOut FourierPoly) {
	scalarMulFloat64Assign(fp0.Coeffs, c, fpOut.Coeffs)
}

// ScalarMulAddFourierPolyAssign computes fpOut += c * fp0.
func (e *Evaluator) ScalarMulAddFourierPolyAssign(fp0 FourierPoly, c float64, fpOut FourierPoly) {
	scalarMulAddFloat64Assign(fp0.Coeffs, c, fpOut.Coeffs)
}

// ScalarMulSubFourierPolyAssign computes fpOut -= c * fp0.
func (e *Evaluator) ScalarMulSubFourierPolyAssign(fp0 FourierPoly, c float64, fpOut FourierPoly) {
	scalarMulSubFloat64Assign(fp0.Coeffs, c, fpOut.Coeffs)
}

// MulFourierPolyAssign computes fpOut = fp0 * fp1.
func (e *Evaluator) MulFourierPolyAssign(fp0, fp1, fpOut FourierPoly) {
	elementWiseMulFloat64Assign(fp0.Coeffs, fp1.Coeffs, fpOut.Coeffs)
}

// MulAddFourierPolyAssign computes fpOut += fp0 * fp1.
func (e *Evaluator) MulAddFourierPolyAssign(fp0, fp1, fpOut FourierPoly) {
	elementWiseMulAddFloat64Assign(fp0.Coeffs, fp1.Coeffs, fpOut.Coeffs)
}

// MulSubFourierPolyAssign computes fpOut -= fp0 * fp1.
func (e *Evaluator) MulSubFourierPolyAssign(fp0, fp1, fpOut FourierPoly) {
	elementWiseMulSubFloat64Assign(fp0.Coeffs, fp1.Coeffs, fpOut.Coeffs)
}

// PermuteFourierPoly returns f0(X^d).
func (e *Evaluator) PermuteFourierPoly(fp0 FourierPoly, d int) FourierPoly {
	fpOut := e.NewFourierPoly()
	e.PermuteFourierPolyAssign(fp0, d, fpOut)
	return fpOut
}

// PermuteFourierPolyAssign computes fpOut = fp0(X^d).
//
// fp0 and fpOut should not overlap. For inplace permutation,
// use [*Evaluator.PermuteFourierPolyInPlace].
func (e *Evaluator) PermuteFourierPolyAssign(fp0 FourierPoly, d int, fpOut FourierPoly) {
	vec.RotateAssign(fp0.Coeffs, -d, fpOut.Coeffs)
}

// PermutePolyInPlace computes fp0 = fp0(X^d).
func (e *Evaluator) PermuteFourierPolyInPlace(fp0 FourierPoly, d int) {
	vec.RotateInPlace(fp0.Coeffs, -d)
}
