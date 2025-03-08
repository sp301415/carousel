package carousel

import "github.com/sp301415/carousel/math/vec"

// AddLWE returns ct0 + ct1.
func (e *Evaluator) AddLWE(ct0, ct1 LWECiphertext) LWECiphertext {
	ctOut := NewLWECiphertext(e.Parameters)
	e.AddLWEAssign(ct0, ct1, ctOut)
	return ctOut
}

// AddLWEAssign computes ctOut = ct0 + ct1.
func (e *Evaluator) AddLWEAssign(ct0, ct1, ctOut LWECiphertext) {
	vec.AddAssign(ct0.Value, ct1.Value, ctOut.Value)
}

// AddPlainLWE returns ct0 + pt.
func (e *Evaluator) AddPlainLWE(ct0 LWECiphertext, pt LWEPlaintext) LWECiphertext {
	ctOut := NewLWECiphertext(e.Parameters)
	e.AddPlainLWEAssign(ct0, pt, ctOut)
	return ctOut
}

// AddPlainLWEAssign computes ctOut = ct0 + pt.
func (e *Evaluator) AddPlainLWEAssign(ct0 LWECiphertext, pt LWEPlaintext, ctOut LWECiphertext) {
	ctOut.CopyFrom(ct0)
	ctOut.Value[0] += pt.Value
}

// SubLWE returns ct0 - ct1.
func (e *Evaluator) SubLWE(ct0, ct1 LWECiphertext) LWECiphertext {
	ctOut := NewLWECiphertext(e.Parameters)
	e.SubLWEAssign(ct0, ct1, ctOut)
	return ctOut
}

// SubLWEAssign computes ctOut = ct0 - ct1.
func (e *Evaluator) SubLWEAssign(ct0, ct1, ctOut LWECiphertext) {
	vec.SubAssign(ct0.Value, ct1.Value, ctOut.Value)
}

// SubPlainLWE returns ct0 - pt.
func (e *Evaluator) SubPlainLWE(ct0 LWECiphertext, pt LWEPlaintext) LWECiphertext {
	ctOut := NewLWECiphertext(e.Parameters)
	e.SubPlainLWEAssign(ct0, pt, ctOut)
	return ctOut
}

// SubPlainLWEAssign computes ctOut = ct0 - pt.
func (e *Evaluator) SubPlainLWEAssign(ct0 LWECiphertext, pt LWEPlaintext, ctOut LWECiphertext) {
	ctOut.CopyFrom(ct0)
	ctOut.Value[0] -= pt.Value
}

// NegLWE returns -ct0.
func (e *Evaluator) NegLWE(ct0 LWECiphertext) LWECiphertext {
	ctOut := NewLWECiphertext(e.Parameters)
	e.NegLWEAssign(ct0, ctOut)
	return ctOut
}

// NegLWEAssign computes ctOut = -ct0.
func (e *Evaluator) NegLWEAssign(ct0, ctOut LWECiphertext) {
	vec.NegAssign(ct0.Value, ctOut.Value)
}

// ScalarMulLWE returns c * ct0.
func (e *Evaluator) ScalarMulLWE(ct0 LWECiphertext, c uint64) LWECiphertext {
	ctOut := NewLWECiphertext(e.Parameters)
	e.ScalarMulLWEAssign(ct0, c, ctOut)
	return ctOut
}

// ScalarMulLWEAssign computes ctOut = c * ct0.
func (e *Evaluator) ScalarMulLWEAssign(ct0 LWECiphertext, c uint64, ctOut LWECiphertext) {
	vec.ScalarMulAssign(ct0.Value, c, ctOut.Value)
}

// ScalarMulAddLWEAssign computes ctOut += c * ct0.
func (e *Evaluator) ScalarMulAddLWEAssign(ct0 LWECiphertext, c uint64, ctOut LWECiphertext) {
	vec.ScalarMulAddAssign(ct0.Value, c, ctOut.Value)
}

// ScalarMulSubLWEAssign computes ctOut -= c * ct0.
func (e *Evaluator) ScalarMulSubLWEAssign(ct0 LWECiphertext, c uint64, ctOut LWECiphertext) {
	vec.ScalarMulSubAssign(ct0.Value, c, ctOut.Value)
}
