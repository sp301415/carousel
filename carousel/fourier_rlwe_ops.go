package carousel

import "github.com/sp301415/carousel/math/poly"

// AddFourierRLWE returns ct0 + ct1.
func (e *Evaluator) AddFourierRLWE(ct0, ct1 FourierRLWECiphertext) FourierRLWECiphertext {
	ctOut := NewFourierRLWECiphertext(e.Parameters)
	e.AddFourierRLWEAssign(ct0, ct1, ctOut)
	return ctOut
}

// AddFourierRLWEAssign computes ctOut = ct0 + ct1.
func (e *Evaluator) AddFourierRLWEAssign(ct0, ct1, ctOut FourierRLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.AddFourierPolyAssign(ct0.Value[i], ct1.Value[i], ctOut.Value[i])
	}
}

// SubFourierRLWE returns ct0 - ct1.
func (e *Evaluator) SubFourierRLWE(ct0, ct1 FourierRLWECiphertext) FourierRLWECiphertext {
	ctOut := NewFourierRLWECiphertext(e.Parameters)
	e.SubFourierRLWEAssign(ct0, ct1, ctOut)
	return ctOut
}

// SubFourierRLWEAssign computes ctOut = ct0 - ct1.
func (e *Evaluator) SubFourierRLWEAssign(ct0, ct1, ctOut FourierRLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.SubFourierPolyAssign(ct0.Value[i], ct1.Value[i], ctOut.Value[i])
	}
}

// NegFourierRLWE returns -ct0.
func (e *Evaluator) NegFourierRLWE(ct0 FourierRLWECiphertext) FourierRLWECiphertext {
	ctOut := NewFourierRLWECiphertext(e.Parameters)
	e.NegFourierRLWEAssign(ct0, ctOut)
	return ctOut
}

// NegFourierRLWEAssign computes ctOut = -ct0.
func (e *Evaluator) NegFourierRLWEAssign(ct0, ctOut FourierRLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.NegFourierPolyAssign(ct0.Value[i], ctOut.Value[i])
	}
}

// ScalarMulFourierRLWE returns c * ct0.
func (e *Evaluator) ScalarMulFourierRLWE(ct0 FourierRLWECiphertext, c float64) FourierRLWECiphertext {
	ctOut := NewFourierRLWECiphertext(e.Parameters)
	e.ScalarMulFourierRLWEAssign(ct0, c, ctOut)
	return ctOut
}

// ScalarMulFourierRLWEAssign computes ctOut = c * ct0.
func (e *Evaluator) ScalarMulFourierRLWEAssign(ct0 FourierRLWECiphertext, c float64, ctOut FourierRLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.ScalarMulFourierPolyAssign(ct0.Value[i], c, ctOut.Value[i])
	}
}

// ScalarMulAddFourierRLWEAssign computes ctOut += c * ct0.
func (e *Evaluator) ScalarMulAddFourierRLWEAssign(ct0 FourierRLWECiphertext, c float64, ctOut FourierRLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.ScalarMulAddFourierPolyAssign(ct0.Value[i], c, ctOut.Value[i])
	}
}

// ScalarMulSubFourierRLWEAssign computes ctOut -= c * ct0.
func (e *Evaluator) ScalarMulSubFourierRLWEAssign(ct0 FourierRLWECiphertext, c float64, ctOut FourierRLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.ScalarMulSubFourierPolyAssign(ct0.Value[i], c, ctOut.Value[i])
	}
}

// PolyMulFourierRLWE returns p * ct0.
func (e *Evaluator) PolyMulFourierRLWE(ct0 FourierRLWECiphertext, p poly.Poly) FourierRLWECiphertext {
	ctOut := NewFourierRLWECiphertext(e.Parameters)
	e.PolyMulFourierRLWEAssign(ct0, p, ctOut)
	return ctOut
}

// PolyMulFourierRLWEAssign computes ctOut = p * ct0.
func (e *Evaluator) PolyMulFourierRLWEAssign(ct0 FourierRLWECiphertext, p poly.Poly, ctOut FourierRLWECiphertext) {
	e.PolyEvaluator.ToFourierPolyAssign(p, e.buffer.fpMul)
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.MulFourierPolyAssign(ct0.Value[i], e.buffer.fpMul, ctOut.Value[i])
	}
}

// PolyMulAddFourierRLWEAssign computes ctOut += p * ct0.
func (e *Evaluator) PolyMulAddFourierRLWEAssign(ct0 FourierRLWECiphertext, p poly.Poly, ctOut FourierRLWECiphertext) {
	e.PolyEvaluator.ToFourierPolyAssign(p, e.buffer.fpMul)
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.MulAddFourierPolyAssign(ct0.Value[i], e.buffer.fpMul, ctOut.Value[i])
	}
}

// PolyMulSubFourierRLWEAssign computes ctOut -= p * ct0.
func (e *Evaluator) PolyMulSubFourierRLWEAssign(ct0 FourierRLWECiphertext, p poly.Poly, ctOut FourierRLWECiphertext) {
	e.PolyEvaluator.ToFourierPolyAssign(p, e.buffer.fpMul)
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.MulSubFourierPolyAssign(ct0.Value[i], e.buffer.fpMul, ctOut.Value[i])
	}
}

// FourierPolyMulFourierRLWE returns fp * ct0.
func (e *Evaluator) FourierPolyMulFourierRLWE(ct0 FourierRLWECiphertext, fp poly.FourierPoly) FourierRLWECiphertext {
	ctOut := NewFourierRLWECiphertext(e.Parameters)
	e.FourierPolyMulFourierRLWEAssign(ct0, fp, ctOut)
	return ctOut
}

// FourierPolyMulFourierRLWEAssign computes ctOut = fp * ct0.
func (e *Evaluator) FourierPolyMulFourierRLWEAssign(ct0 FourierRLWECiphertext, fp poly.FourierPoly, ctOut FourierRLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.MulFourierPolyAssign(ct0.Value[i], fp, ctOut.Value[i])
	}
}

// FourierPolyMulAddFourierRLWEAssign computes ctOut += fp * ct0.
func (e *Evaluator) FourierPolyMulAddFourierRLWEAssign(ct0 FourierRLWECiphertext, fp poly.FourierPoly, ctOut FourierRLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.MulAddFourierPolyAssign(ct0.Value[i], fp, ctOut.Value[i])
	}
}

// FourierPolyMulSubFourierRLWEAssign computes ctOut -= fp * ct0.
func (e *Evaluator) FourierPolyMulSubFourierRLWEAssign(ct0 FourierRLWECiphertext, fp poly.FourierPoly, ctOut FourierRLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.MulSubFourierPolyAssign(ct0.Value[i], fp, ctOut.Value[i])
	}
}

// PermuteRLWEAssign computes ctOut = ct0(X^d).
//
// ct0 and ctOut should not overlap. For inplace permutation,
// use [*Evaluator.PermuteFourierRLWEInPlace].
func (e *Evaluator) PermuteFourierRLWEAssign(ct0 FourierRLWECiphertext, d int, ctOut FourierRLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.PermuteFourierPolyAssign(ct0.Value[i], d, ctOut.Value[i])
	}
}

// PermuteRLWEInPlace computes ct0 = ct0(X^d).
func (e *Evaluator) PermuteFourierRLWEInPlace(ct0 FourierRLWECiphertext, d int) {
	for i := 0; i < len(ct0.Value); i++ {
		e.PolyEvaluator.PermuteFourierPolyInPlace(ct0.Value[i], d)
	}
}

// PermuteAddRLWEAssign computes ctOut += ct0(X^d).
//
// ct0 and ctOut should not overlap.
func (e *Evaluator) PermuteAddFourierRLWEAssign(ct0 FourierRLWECiphertext, d int, ctOut FourierRLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.PermuteAddFourierPolyAssign(ct0.Value[i], d, ctOut.Value[i])
	}
}

// PermuteSubRLWEAssign computes ctOut -= ct0(X^d).
//
// ct0 and ctOut should not overlap.
func (e *Evaluator) PermuteSubFourierRLWEAssign(ct0 FourierRLWECiphertext, d int, ctOut FourierRLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.PermuteSubFourierPolyAssign(ct0.Value[i], d, ctOut.Value[i])
	}
}
