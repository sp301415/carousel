package carousel

import (
	"github.com/sp301415/carousel/math/poly"
	"github.com/sp301415/carousel/math/vec"
)

// GadgetProductLWE returns the gadget product between c and ctLev.
func (e *Evaluator) GadgetProductLWE(ctLev LevCiphertext, c uint64) LWECiphertext {
	ctOut := NewLWECiphertext(e.Parameters)
	e.GadgetProductLWEAssign(ctLev, c, ctOut)
	return ctOut
}

// GadgetProductLWEAssign computes the gadget product between c and ctLev and writes it to ctLWEOut.
func (e *Evaluator) GadgetProductLWEAssign(ctLev LevCiphertext, c uint64, ctLWEOut LWECiphertext) {
	scalarDecomposed := e.Decomposer.ScalarDecomposedBuffer(ctLev.GadgetParameters)

	e.Decomposer.DecomposeScalarAssign(c, ctLev.GadgetParameters, scalarDecomposed)
	vec.ScalarMulAssign(ctLev.Value[0].Value, scalarDecomposed[0], ctLWEOut.Value)
	for i := 1; i < ctLev.GadgetParameters.level; i++ {
		vec.ScalarMulAddAssign(ctLev.Value[i].Value, scalarDecomposed[i], ctLWEOut.Value)
	}
}

// GadgetProductAddLWEAssign computes the gadget product between c and ctLev and adds it to ctLWEOut.
func (e *Evaluator) GadgetProductAddLWEAssign(ctLev LevCiphertext, c uint64, ctLWEOut LWECiphertext) {
	scalarDecomposed := e.Decomposer.ScalarDecomposedBuffer(ctLev.GadgetParameters)

	e.Decomposer.DecomposeScalarAssign(c, ctLev.GadgetParameters, scalarDecomposed)
	for i := 0; i < ctLev.GadgetParameters.level; i++ {
		vec.ScalarMulAddAssign(ctLev.Value[i].Value, scalarDecomposed[i], ctLWEOut.Value)
	}
}

// GadgetProductSubLWEAssign computes the gadget product between c and ctLev and subtracts it from ctLWEOut.
func (e *Evaluator) GadgetProductSubLWEAssign(ctLev LevCiphertext, c uint64, ctLWEOut LWECiphertext) {
	scalarDecomposed := e.Decomposer.ScalarDecomposedBuffer(ctLev.GadgetParameters)

	e.Decomposer.DecomposeScalarAssign(c, ctLev.GadgetParameters, scalarDecomposed)
	for i := 0; i < ctLev.GadgetParameters.level; i++ {
		vec.ScalarMulSubAssign(ctLev.Value[i].Value, scalarDecomposed[i], ctLWEOut.Value)
	}
}

// GadgetProductRLWE returns the gadget product between p and ctFourierRLev.
func (e *Evaluator) GadgetProductRLWE(ctFourierRLev FourierRLevCiphertext, p poly.Poly) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.GadgetProductRLWEAssign(ctFourierRLev, p, ctOut)
	return ctOut
}

// GadgetProductRLWEAssign computes the gadget product between p and ctFourierRLev and writes it to ctRLWEOut.
func (e *Evaluator) GadgetProductRLWEAssign(ctFourierRLev FourierRLevCiphertext, p poly.Poly, ctRLWEOut RLWECiphertext) {
	polyDecomposed := e.Decomposer.PolyDecomposedBuffer(ctFourierRLev.GadgetParameters)
	polyFourierDecomposed := e.Decomposer.PolyFourierDecomposedBuffer(ctFourierRLev.GadgetParameters)

	e.Decomposer.DecomposePolyAssign(p, ctFourierRLev.GadgetParameters, polyDecomposed)
	for i := 0; i < ctFourierRLev.GadgetParameters.level; i++ {
		e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[i], polyFourierDecomposed[i])
	}

	e.FourierPolyMulFourierRLWEAssign(ctFourierRLev.Value[0], polyFourierDecomposed[0], e.buffer.ctProdFourierRLWE)
	for i := 1; i < ctFourierRLev.GadgetParameters.level; i++ {
		e.FourierPolyMulAddFourierRLWEAssign(ctFourierRLev.Value[i], polyFourierDecomposed[i], e.buffer.ctProdFourierRLWE)
	}

	e.PolyEvaluator.ToPolyAssignUnsafe(e.buffer.ctProdFourierRLWE.Value[0], ctRLWEOut.Value[0])
	e.PolyEvaluator.ToPolyAssignUnsafe(e.buffer.ctProdFourierRLWE.Value[1], ctRLWEOut.Value[1])
}

// GadgetProductAddRLWEAssign computes the gadget product between p and ctFourierRLev and adds it to ctRLWEOut.
func (e *Evaluator) GadgetProductAddRLWEAssign(ctFourierRLev FourierRLevCiphertext, p poly.Poly, ctRLWEOut RLWECiphertext) {
	polyDecomposed := e.Decomposer.PolyDecomposedBuffer(ctFourierRLev.GadgetParameters)
	polyFourierDecomposed := e.Decomposer.PolyFourierDecomposedBuffer(ctFourierRLev.GadgetParameters)

	e.Decomposer.DecomposePolyAssign(p, ctFourierRLev.GadgetParameters, polyDecomposed)
	for i := 0; i < ctFourierRLev.GadgetParameters.level; i++ {
		e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[i], polyFourierDecomposed[i])
	}

	e.FourierPolyMulFourierRLWEAssign(ctFourierRLev.Value[0], polyFourierDecomposed[0], e.buffer.ctProdFourierRLWE)
	for i := 1; i < ctFourierRLev.GadgetParameters.level; i++ {
		e.FourierPolyMulAddFourierRLWEAssign(ctFourierRLev.Value[i], polyFourierDecomposed[i], e.buffer.ctProdFourierRLWE)
	}

	e.PolyEvaluator.ToPolyAddAssignUnsafe(e.buffer.ctProdFourierRLWE.Value[0], ctRLWEOut.Value[0])
	e.PolyEvaluator.ToPolyAddAssignUnsafe(e.buffer.ctProdFourierRLWE.Value[1], ctRLWEOut.Value[1])
}

// GadgetProductSubRLWEAssign computes the gadget product between p and ctFourierRLev and subtracts it from ctRLWEOut.
func (e *Evaluator) GadgetProductSubRLWEAssign(ctFourierRLev FourierRLevCiphertext, p poly.Poly, ctRLWEOut RLWECiphertext) {
	polyDecomposed := e.Decomposer.PolyDecomposedBuffer(ctFourierRLev.GadgetParameters)
	polyFourierDecomposed := e.Decomposer.PolyFourierDecomposedBuffer(ctFourierRLev.GadgetParameters)

	e.Decomposer.DecomposePolyAssign(p, ctFourierRLev.GadgetParameters, polyDecomposed)
	for i := 0; i < ctFourierRLev.GadgetParameters.level; i++ {
		e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[i], polyFourierDecomposed[i])
	}

	e.FourierPolyMulFourierRLWEAssign(ctFourierRLev.Value[0], polyFourierDecomposed[0], e.buffer.ctProdFourierRLWE)
	for i := 1; i < ctFourierRLev.GadgetParameters.level; i++ {
		e.FourierPolyMulAddFourierRLWEAssign(ctFourierRLev.Value[i], polyFourierDecomposed[i], e.buffer.ctProdFourierRLWE)
	}

	e.PolyEvaluator.ToPolySubAssignUnsafe(e.buffer.ctProdFourierRLWE.Value[0], ctRLWEOut.Value[0])
	e.PolyEvaluator.ToPolySubAssignUnsafe(e.buffer.ctProdFourierRLWE.Value[1], ctRLWEOut.Value[1])
}

// GadgetProductFourierDecomposedFourierRLWEAssign computes the gadget product between fourier decomposed p and ctFourierRLev and writes it to ctFourierRLWEOut.
func (e *Evaluator) GadgetProductFourierDecomposedFourierRLWEAssign(ctFourierRLev FourierRLevCiphertext, pDecomposed []poly.FourierPoly, ctFourierRLWEOut FourierRLWECiphertext) {
	e.FourierPolyMulFourierRLWEAssign(ctFourierRLev.Value[0], pDecomposed[0], ctFourierRLWEOut)
	for i := 1; i < ctFourierRLev.GadgetParameters.level; i++ {
		e.FourierPolyMulAddFourierRLWEAssign(ctFourierRLev.Value[i], pDecomposed[i], ctFourierRLWEOut)
	}
}

// GadgetProductFourierDecomposedAddFourierRLWEAssign computes the gadget product between fourier decomposed p and ctFourierRLev and adds it to ctFourierRLWEOut.
func (e *Evaluator) GadgetProductFourierDecomposedAddFourierRLWEAssign(ctFourierRLev FourierRLevCiphertext, pDecomposed []poly.FourierPoly, ctFourierRLWEOut FourierRLWECiphertext) {
	for i := 0; i < ctFourierRLev.GadgetParameters.level; i++ {
		e.FourierPolyMulAddFourierRLWEAssign(ctFourierRLev.Value[i], pDecomposed[i], ctFourierRLWEOut)
	}
}

// GadgetProductFourierDecomposedSubFourierRLWEAssign computes the gadget product between fourier decomposed p and ctFourierRLev and subtracts it from ctFourierRLWEOut.
func (e *Evaluator) GadgetProductFourierDecomposedSubFourierRLWEAssign(ctFourierRLev FourierRLevCiphertext, pDecomposed []poly.FourierPoly, ctFourierRLWEOut FourierRLWECiphertext) {
	for i := 0; i < ctFourierRLev.GadgetParameters.level; i++ {
		e.FourierPolyMulSubFourierRLWEAssign(ctFourierRLev.Value[i], pDecomposed[i], ctFourierRLWEOut)
	}
}

// ExternalProductLWE returns the external product between ctGSW and ctLWE.
func (e *Evaluator) ExternalProductLWE(ctGSW GSWCiphertext, ctLWE LWECiphertext) LWECiphertext {
	ctOut := NewLWECiphertext(e.Parameters)
	e.ExternalProductLWEAssign(ctGSW, ctLWE, ctOut)
	return ctOut
}

// ExternalProductLWEAssign computes the external product between ctGSW and ctLWE and writes it to ctOut.
func (e *Evaluator) ExternalProductLWEAssign(ctGSW GSWCiphertext, ctLWE LWECiphertext, ctLWEOut LWECiphertext) {
	decomposed := e.Decomposer.ScalarDecomposedBuffer(ctGSW.GadgetParameters)

	e.Decomposer.DecomposeScalarAssign(ctLWE.Value[0], ctGSW.GadgetParameters, decomposed)
	vec.ScalarMulAssign(ctGSW.Value[0].Value[0].Value, decomposed[0], e.buffer.ctProdLWE.Value)
	for j := 1; j < ctGSW.GadgetParameters.level; j++ {
		vec.ScalarMulAddAssign(ctGSW.Value[0].Value[j].Value, decomposed[j], e.buffer.ctProdLWE.Value)
	}

	for i := 1; i < e.Parameters.lweDimension+1; i++ {
		e.Decomposer.DecomposeScalarAssign(ctLWE.Value[i], ctGSW.GadgetParameters, decomposed)
		for j := 0; j < ctGSW.GadgetParameters.level; j++ {
			vec.ScalarMulAddAssign(ctGSW.Value[i].Value[j].Value, decomposed[j], e.buffer.ctProdLWE.Value)
		}
	}

	ctLWEOut.CopyFrom(e.buffer.ctProdLWE)
}

// ExternalProductAddLWEAssign computes the external product between ctGSW and ctLWE and adds it to ctLWEOut.
func (e *Evaluator) ExternalProductAddLWEAssign(ctGSW GSWCiphertext, ctLWE LWECiphertext, ctLWEOut LWECiphertext) {
	decomposed := e.Decomposer.ScalarDecomposedBuffer(ctGSW.GadgetParameters)

	e.Decomposer.DecomposeScalarAssign(ctLWE.Value[0], ctGSW.GadgetParameters, decomposed)
	vec.ScalarMulAssign(ctGSW.Value[0].Value[0].Value, decomposed[0], e.buffer.ctProdLWE.Value)
	for j := 1; j < ctGSW.GadgetParameters.level; j++ {
		vec.ScalarMulAddAssign(ctGSW.Value[0].Value[j].Value, decomposed[j], e.buffer.ctProdLWE.Value)
	}

	for i := 1; i < e.Parameters.lweDimension+1; i++ {
		e.Decomposer.DecomposeScalarAssign(ctLWE.Value[i], ctGSW.GadgetParameters, decomposed)
		for j := 0; j < ctGSW.GadgetParameters.level; j++ {
			vec.ScalarMulAddAssign(ctGSW.Value[i].Value[j].Value, decomposed[j], e.buffer.ctProdLWE.Value)
		}
	}

	e.AddLWEAssign(ctLWEOut, e.buffer.ctProdLWE, ctLWEOut)
}

// ExternalProductSubLWEAssign computes the external product between ctGSW and ctLWE and subtracts it from ctLWEOut.
func (e *Evaluator) ExternalProductSubLWEAssign(ctGSW GSWCiphertext, ctLWE LWECiphertext, ctLWEOut LWECiphertext) {
	decomposed := e.Decomposer.ScalarDecomposedBuffer(ctGSW.GadgetParameters)

	e.Decomposer.DecomposeScalarAssign(ctLWE.Value[0], ctGSW.GadgetParameters, decomposed)
	vec.ScalarMulAssign(ctGSW.Value[0].Value[0].Value, decomposed[0], e.buffer.ctProdLWE.Value)
	for j := 1; j < ctGSW.GadgetParameters.level; j++ {
		vec.ScalarMulAddAssign(ctGSW.Value[0].Value[j].Value, decomposed[j], e.buffer.ctProdLWE.Value)
	}

	for i := 1; i < e.Parameters.lweDimension+1; i++ {
		e.Decomposer.DecomposeScalarAssign(ctLWE.Value[i], ctGSW.GadgetParameters, decomposed)
		for j := 0; j < ctGSW.GadgetParameters.level; j++ {
			vec.ScalarMulAddAssign(ctGSW.Value[i].Value[j].Value, decomposed[j], e.buffer.ctProdLWE.Value)
		}
	}

	e.SubLWEAssign(ctLWEOut, e.buffer.ctProdLWE, ctLWEOut)
}

// ExternalProductRLWE returns the external product between ctFourierRGSW and ctRLWE.
func (e *Evaluator) ExternalProductRLWE(ctFourierRGSW FourierRGSWCiphertext, ctRLWE RLWECiphertext) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.ExternalProductRLWEAssign(ctFourierRGSW, ctRLWE, ctOut)
	return ctOut
}

// ExternalProductRLWEAssign computes the external product between ctFourierRGSW and ctRLWE and writes it to ctOut.
func (e *Evaluator) ExternalProductRLWEAssign(ctFourierRGSW FourierRGSWCiphertext, ctRLWE, ctRLWEOut RLWECiphertext) {
	polyDecomposed := e.Decomposer.PolyDecomposedBuffer(ctFourierRGSW.GadgetParameters)

	e.Decomposer.DecomposePolyAssign(ctRLWE.Value[0], ctFourierRGSW.GadgetParameters, polyDecomposed)
	e.PolyMulFourierRLWEAssign(ctFourierRGSW.Value[0].Value[0], polyDecomposed[0], e.buffer.ctProdFourierRLWE)
	for j := 1; j < ctFourierRGSW.GadgetParameters.level; j++ {
		e.PolyMulAddFourierRLWEAssign(ctFourierRGSW.Value[0].Value[j], polyDecomposed[j], e.buffer.ctProdFourierRLWE)
	}

	e.Decomposer.DecomposePolyAssign(ctRLWE.Value[1], ctFourierRGSW.GadgetParameters, polyDecomposed)
	for j := 0; j < ctFourierRGSW.GadgetParameters.level; j++ {
		e.PolyMulAddFourierRLWEAssign(ctFourierRGSW.Value[1].Value[j], polyDecomposed[j], e.buffer.ctProdFourierRLWE)
	}

	e.PolyEvaluator.ToPolyAssignUnsafe(e.buffer.ctProdFourierRLWE.Value[0], ctRLWEOut.Value[0])
	e.PolyEvaluator.ToPolyAssignUnsafe(e.buffer.ctProdFourierRLWE.Value[1], ctRLWEOut.Value[1])
}

// ExternalProductAddRLWEAssign computes the external product between ctFourierRGSW and ctRLWE and adds it to ctRLWEOut.
func (e *Evaluator) ExternalProductAddRLWEAssign(ctFourierRGSW FourierRGSWCiphertext, ctRLWE, ctRLWEOut RLWECiphertext) {
	polyDecomposed := e.Decomposer.PolyDecomposedBuffer(ctFourierRGSW.GadgetParameters)

	e.Decomposer.DecomposePolyAssign(ctRLWE.Value[0], ctFourierRGSW.GadgetParameters, polyDecomposed)
	e.PolyMulFourierRLWEAssign(ctFourierRGSW.Value[0].Value[0], polyDecomposed[0], e.buffer.ctProdFourierRLWE)
	for j := 1; j < ctFourierRGSW.GadgetParameters.level; j++ {
		e.PolyMulAddFourierRLWEAssign(ctFourierRGSW.Value[0].Value[j], polyDecomposed[j], e.buffer.ctProdFourierRLWE)
	}

	e.Decomposer.DecomposePolyAssign(ctRLWE.Value[1], ctFourierRGSW.GadgetParameters, polyDecomposed)
	for j := 0; j < ctFourierRGSW.GadgetParameters.level; j++ {
		e.PolyMulAddFourierRLWEAssign(ctFourierRGSW.Value[1].Value[j], polyDecomposed[j], e.buffer.ctProdFourierRLWE)
	}

	e.PolyEvaluator.ToPolyAddAssignUnsafe(e.buffer.ctProdFourierRLWE.Value[0], ctRLWEOut.Value[0])
	e.PolyEvaluator.ToPolyAddAssignUnsafe(e.buffer.ctProdFourierRLWE.Value[1], ctRLWEOut.Value[1])
}

// ExternalProductSubRLWEAssign computes the external product between ctFourierRGSW and ctRLWE and subtracts it from ctRLWEOut.
func (e *Evaluator) ExternalProductSubRLWEAssign(ctFourierRGSW FourierRGSWCiphertext, ctRLWE, ctRLWEOut RLWECiphertext) {
	polyDecomposed := e.Decomposer.PolyDecomposedBuffer(ctFourierRGSW.GadgetParameters)

	e.Decomposer.DecomposePolyAssign(ctRLWE.Value[0], ctFourierRGSW.GadgetParameters, polyDecomposed)
	e.PolyMulFourierRLWEAssign(ctFourierRGSW.Value[0].Value[0], polyDecomposed[0], e.buffer.ctProdFourierRLWE)
	for j := 1; j < ctFourierRGSW.GadgetParameters.level; j++ {
		e.PolyMulAddFourierRLWEAssign(ctFourierRGSW.Value[0].Value[j], polyDecomposed[j], e.buffer.ctProdFourierRLWE)
	}

	e.Decomposer.DecomposePolyAssign(ctRLWE.Value[1], ctFourierRGSW.GadgetParameters, polyDecomposed)
	for j := 0; j < ctFourierRGSW.GadgetParameters.level; j++ {
		e.PolyMulAddFourierRLWEAssign(ctFourierRGSW.Value[1].Value[j], polyDecomposed[j], e.buffer.ctProdFourierRLWE)
	}

	e.PolyEvaluator.ToPolySubAssignUnsafe(e.buffer.ctProdFourierRLWE.Value[0], ctRLWEOut.Value[0])
	e.PolyEvaluator.ToPolySubAssignUnsafe(e.buffer.ctProdFourierRLWE.Value[1], ctRLWEOut.Value[1])
}

// ExternalProductFourierDecomposedFourierRLWEAssign computes the external product between ctFourierRGSW and fourier decomposed ctRLWE and writes it to ctFourierRLWEOut.
func (e *Evaluator) ExternalProductFourierDecomposedFourierRLWEAssign(ctFourierRGSW FourierRGSWCiphertext, ctRLWEDecomposed [][]poly.FourierPoly, ctFourierRLWEOut FourierRLWECiphertext) {
	e.FourierPolyMulFourierRLWEAssign(ctFourierRGSW.Value[0].Value[0], ctRLWEDecomposed[0][0], ctFourierRLWEOut)
	for j := 1; j < ctFourierRGSW.GadgetParameters.level; j++ {
		e.FourierPolyMulAddFourierRLWEAssign(ctFourierRGSW.Value[0].Value[j], ctRLWEDecomposed[0][j], ctFourierRLWEOut)
	}

	for j := 0; j < ctFourierRGSW.GadgetParameters.level; j++ {
		e.FourierPolyMulAddFourierRLWEAssign(ctFourierRGSW.Value[1].Value[j], ctRLWEDecomposed[1][j], ctFourierRLWEOut)
	}
}

// ExternalProductFourierDecomposedAddFourierRLWEAssign computes the external product between ctFourierRGSW and fourier decomposed ctRLWE and adds it to ctFourierRLWEOut.
func (e *Evaluator) ExternalProductFourierDecomposedAddFourierRLWEAssign(ctFourierRGSW FourierRGSWCiphertext, ctRLWEDecomposed [][]poly.FourierPoly, ctFourierRLWEOut FourierRLWECiphertext) {
	for j := 0; j < ctFourierRGSW.GadgetParameters.level; j++ {
		e.FourierPolyMulAddFourierRLWEAssign(ctFourierRGSW.Value[0].Value[j], ctRLWEDecomposed[0][j], ctFourierRLWEOut)
	}

	for j := 0; j < ctFourierRGSW.GadgetParameters.level; j++ {
		e.FourierPolyMulAddFourierRLWEAssign(ctFourierRGSW.Value[1].Value[j], ctRLWEDecomposed[1][j], ctFourierRLWEOut)
	}
}

// ExternalProductFourierDecomposedSubFourierRLWEAssign computes the external product between ctFourierRGSW and fourier decomposed ctRLWE and subtracts it from ctFourierRLWEOut.
func (e *Evaluator) ExternalProductFourierDecomposedSubFourierRLWEAssign(ctFourierRGSW FourierRGSWCiphertext, ctRLWEDecomposed [][]poly.FourierPoly, ctFourierRLWEOut FourierRLWECiphertext) {
	for j := 0; j < ctFourierRGSW.GadgetParameters.level; j++ {
		e.FourierPolyMulSubFourierRLWEAssign(ctFourierRGSW.Value[0].Value[j], ctRLWEDecomposed[0][j], ctFourierRLWEOut)
	}

	for j := 0; j < ctFourierRGSW.GadgetParameters.level; j++ {
		e.FourierPolyMulSubFourierRLWEAssign(ctFourierRGSW.Value[1].Value[j], ctRLWEDecomposed[1][j], ctFourierRLWEOut)
	}
}
