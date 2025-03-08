package carousel

import "github.com/sp301415/carousel/math/vec"

// KeySwitchLWE switches key of ct.
// Input ciphertext should be of length ksk.InputLWEDimension + 1.
func (e *Evaluator) KeySwitchLWE(ct LWECiphertext, ksk LWEKeySwitchKey) LWECiphertext {
	ctOut := NewLWECiphertext(e.Parameters)
	e.KeySwitchLWEAssign(ct, ksk, ctOut)
	return ctOut
}

// KeySwitchLWEAssign switches key of ct and writes it to ctOut.
// Input ciphertext should be of length ksk.InputLWEDimension + 1.
func (e *Evaluator) KeySwitchLWEAssign(ct LWECiphertext, ksk LWEKeySwitchKey, ctOut LWECiphertext) {
	scalarDecomposed := e.Decomposer.ScalarDecomposedBuffer(ksk.GadgetParameters)

	e.Decomposer.DecomposeScalarAssign(ct.Value[1], ksk.GadgetParameters, scalarDecomposed)
	e.ScalarMulLWEAssign(ksk.Value[0].Value[0], scalarDecomposed[0], e.buffer.ctProdLWE)
	for j := 1; j < ksk.GadgetParameters.level; j++ {
		e.ScalarMulAddLWEAssign(ksk.Value[0].Value[j], scalarDecomposed[j], e.buffer.ctProdLWE)
	}

	for i := 1; i < ksk.InputLWEDimension(); i++ {
		e.Decomposer.DecomposeScalarAssign(ct.Value[i+1], ksk.GadgetParameters, scalarDecomposed)
		for j := 0; j < ksk.GadgetParameters.level; j++ {
			e.ScalarMulAddLWEAssign(ksk.Value[i].Value[j], scalarDecomposed[j], e.buffer.ctProdLWE)
		}
	}

	ctOut.Value[0] = ct.Value[0] + e.buffer.ctProdLWE.Value[0]
	vec.CopyAssign(e.buffer.ctProdLWE.Value[1:], ctOut.Value[1:])
}

// KeySwitchRLWE switches key of ct.
func (e *Evaluator) KeySwitchRLWE(ct RLWECiphertext, ksk RLWEKeySwitchKey) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.KeySwitchRLWEAssign(ct, ksk, ctOut)
	return ctOut
}

// KeySwitchRLWEAssign switches key of ct and writes it to ctOut.
func (e *Evaluator) KeySwitchRLWEAssign(ct RLWECiphertext, ksk RLWEKeySwitchKey, ctOut RLWECiphertext) {
	fourierDecomposed := e.Decomposer.PolyFourierDecomposedBuffer(ksk.GadgetParameters)

	e.Decomposer.FourierDecomposePolyAssign(ct.Value[1], ksk.GadgetParameters, fourierDecomposed)
	e.FourierPolyMulFourierRLWEAssign(ksk.Value.Value[0], fourierDecomposed[0], e.buffer.ctProdFourierRLWE)
	for j := 1; j < ksk.GadgetParameters.level; j++ {
		e.FourierPolyMulAddFourierRLWEAssign(ksk.Value.Value[j], fourierDecomposed[j], e.buffer.ctProdFourierRLWE)
	}

	ctOut.Value[0].CopyFrom(ct.Value[0])
	e.PolyEvaluator.ToPolyAddAssignUnsafe(e.buffer.ctProdFourierRLWE.Value[0], ctOut.Value[0])
	e.PolyEvaluator.ToPolyAddAssignUnsafe(e.buffer.ctProdFourierRLWE.Value[1], ctOut.Value[1])
}
