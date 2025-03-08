package carousel

import "github.com/sp301415/carousel/math/poly"

// AddRLWE returns ct0 + ct1.
func (e *Evaluator) AddRLWE(ct0, ct1 RLWECiphertext) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.AddRLWEAssign(ct0, ct1, ctOut)
	return ctOut
}

// AddRLWEAssign computes ctOut = ct0 + ct1.
func (e *Evaluator) AddRLWEAssign(ct0, ct1, ctOut RLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.AddPolyAssign(ct0.Value[i], ct1.Value[i], ctOut.Value[i])
	}
}

// AddPlainRLWE returns ct0 + pt.
func (e *Evaluator) AddPlainRLWE(ct0 RLWECiphertext, pt RLWEPlaintext) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.AddPlainRLWEAssign(ct0, pt, ctOut)
	return ctOut
}

// AddPlainRLWEAssign computes ctOut = ct0 + pt.
func (e *Evaluator) AddPlainRLWEAssign(ct0 RLWECiphertext, pt RLWEPlaintext, ctOut RLWECiphertext) {
	for i := 1; i < len(ctOut.Value); i++ {
		ctOut.Value[i].CopyFrom(ct0.Value[i])
	}
	e.PolyEvaluator.AddPolyAssign(ct0.Value[0], pt.Value, ctOut.Value[0])
}

// SubRLWE returns ct0 - ct1.
func (e *Evaluator) SubRLWE(ct0, ct1 RLWECiphertext) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.SubRLWEAssign(ct0, ct1, ctOut)
	return ctOut
}

// SubRLWEAssign computes ctOut = ct0 - ct1.
func (e *Evaluator) SubRLWEAssign(ct0, ct1, ctOut RLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.SubPolyAssign(ct0.Value[i], ct1.Value[i], ctOut.Value[i])
	}
}

// SubPlainRLWE returns ct0 - pt.
func (e *Evaluator) SubPlainRLWE(ct0 RLWECiphertext, pt RLWEPlaintext) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.SubPlainRLWEAssign(ct0, pt, ctOut)
	return ctOut
}

// SubPlainRLWEAssign computes ctOut = ct0 - pt.
func (e *Evaluator) SubPlainRLWEAssign(ct0 RLWECiphertext, pt RLWEPlaintext, ctOut RLWECiphertext) {
	for i := 1; i < len(ctOut.Value); i++ {
		ctOut.Value[i].CopyFrom(ct0.Value[i])
	}
	e.PolyEvaluator.SubPolyAssign(ct0.Value[0], pt.Value, ctOut.Value[0])
}

// NegRLWE returns -ct0.
func (e *Evaluator) NegRLWE(ct0 RLWECiphertext) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.NegRLWEAssign(ct0, ctOut)
	return ctOut
}

// NegRLWEAssign computes ctOut = -ct0.
func (e *Evaluator) NegRLWEAssign(ct0, ctOut RLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.NegPolyAssign(ct0.Value[i], ctOut.Value[i])
	}
}

// ScalarMulRLWE returns c * ct0.
func (e *Evaluator) ScalarMulRLWE(ct0 RLWECiphertext, c uint64) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.ScalarMulRLWEAssign(ct0, c, ctOut)
	return ctOut
}

// ScalarMulRLWEAssign computes ctOut = c * ct0.
func (e *Evaluator) ScalarMulRLWEAssign(ct0 RLWECiphertext, c uint64, ctOut RLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.ScalarMulPolyAssign(ct0.Value[i], c, ctOut.Value[i])
	}
}

// ScalarMulAddRLWEAssign computes ctOut += c * ct0.
func (e *Evaluator) ScalarMulAddRLWEAssign(ct0 RLWECiphertext, c uint64, ctOut RLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.ScalarMulAddPolyAssign(ct0.Value[i], c, ctOut.Value[i])
	}
}

// ScalarMulSubRLWEAssign computes ctOut -= c * ct0.
func (e *Evaluator) ScalarMulSubRLWEAssign(ct0 RLWECiphertext, c uint64, ctOut RLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.ScalarMulSubPolyAssign(ct0.Value[i], c, ctOut.Value[i])
	}
}

// PolyMulRLWE returns p * ct0.
func (e *Evaluator) PolyMulRLWE(ct0 RLWECiphertext, p poly.Poly) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.PolyMulRLWEAssign(ct0, p, ctOut)
	return ctOut
}

// PolyMulRLWEAssign computes ctOut = p * ct0.
func (e *Evaluator) PolyMulRLWEAssign(ct0 RLWECiphertext, p poly.Poly, ctOut RLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.MulPolyAssign(ct0.Value[i], p, ctOut.Value[i])
	}
}

// PolyMulAddRLWEAssign computes ctOut += p * ct0.
func (e *Evaluator) PolyMulAddRLWEAssign(ct0 RLWECiphertext, p poly.Poly, ctOut RLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.MulAddPolyAssign(ct0.Value[i], p, ctOut.Value[i])
	}
}

// PolyMulSubRLWEAssign computes ctOut -= p * ct0.
func (e *Evaluator) PolyMulSubRLWEAssign(ct0 RLWECiphertext, p poly.Poly, ctOut RLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.MulSubPolyAssign(ct0.Value[i], p, ctOut.Value[i])
	}
}

// FourierPolyMulRLWE returns fp * ct0.
func (e *Evaluator) FourierPolyMulRLWE(ct0 RLWECiphertext, fp poly.FourierPoly) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.FourierPolyMulRLWEAssign(ct0, fp, ctOut)
	return ctOut
}

// FourierPolyMulRLWEAssign computes ctOut = fp * ct0.
func (e *Evaluator) FourierPolyMulRLWEAssign(ct0 RLWECiphertext, fp poly.FourierPoly, ctOut RLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.FourierPolyMulPolyAssign(ct0.Value[i], fp, ctOut.Value[i])
	}
}

// FourierPolyMulAddRLWEAssign computes ctOut += fp * ct0.
func (e *Evaluator) FourierPolyMulAddRLWEAssign(ct0 RLWECiphertext, fp poly.FourierPoly, ctOut RLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.FourierPolyMulAddPolyAssign(ct0.Value[i], fp, ctOut.Value[i])
	}
}

// FourierPolyMulSubRLWEAssign computes ctOut -= fp * ct0.
func (e *Evaluator) FourierPolyMulSubRLWEAssign(ct0 RLWECiphertext, fp poly.FourierPoly, ctOut RLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.FourierPolyMulSubPolyAssign(ct0.Value[i], fp, ctOut.Value[i])
	}
}

// PermuteRLWE returns ctOut = ct0(X^d).
func (e *Evaluator) PermuteRLWE(ct0 RLWECiphertext, d int) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.PermuteRLWEAssign(ct0, d, ctOut)
	return ctOut
}

// PermuteRLWEAssign computes ctOut = ct0(X^d).
//
// ct0 and ctOut should not overlap. For inplace permutation,
// use [*Evaluator.PermuteRLWEInPlace].
func (e *Evaluator) PermuteRLWEAssign(ct0 RLWECiphertext, d int, ctOut RLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.PermutePolyAssign(ct0.Value[i], d, ctOut.Value[i])
	}
}

// PermuteRLWEInPlace computes ct0 = ct0(X^d).
func (e *Evaluator) PermuteRLWEInPlace(ct0 RLWECiphertext, d int) {
	for i := 0; i < len(ct0.Value); i++ {
		e.PolyEvaluator.PermutePolyInPlace(ct0.Value[i], d)
	}
}

// PermuteAddRLWEAssign computes ctOut += ct0(X^d).
//
// ct0 and ctOut should not overlap.
func (e *Evaluator) PermuteAddRLWEAssign(ct0 RLWECiphertext, d int, ctOut RLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.PermuteAddPolyAssign(ct0.Value[i], d, ctOut.Value[i])
	}
}

// PermuteSubRLWEAssign computes ctOut -= ct0(X^d).
//
// ct0 and ctOut should not overlap.
func (e *Evaluator) PermuteSubRLWEAssign(ct0 RLWECiphertext, d int, ctOut RLWECiphertext) {
	for i := 0; i < len(ctOut.Value); i++ {
		e.PolyEvaluator.PermuteSubPolyAssign(ct0.Value[i], d, ctOut.Value[i])
	}
}
