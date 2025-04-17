package carousel

import (
	"math"

	"github.com/sp301415/carousel/math/vec"
)

// BootstrapFunc returns a bootstrapped LWE ciphertext with respect to given function.
func (e *Evaluator) BootstrapFunc(ct RLWECiphertext, f func(int) int) RLWECiphertext {
	e.GenLookUpTableAssign(f, e.buffer.lut)
	return e.BootstrapLUT(ct, e.buffer.lut)
}

// BootstrapFuncAssign bootstraps LWE ciphertext with respect to given function and writes it to ctOut.
func (e *Evaluator) BootstrapFuncAssign(ct RLWECiphertext, f func(int) int, ctOut RLWECiphertext) {
	e.GenLookUpTableAssign(f, e.buffer.lut)
	e.BootstrapLUTAssign(ct, e.buffer.lut, ctOut)
}

// BootstrapLUT returns a bootstrapped LWE ciphertext with respect to given LUT.
func (e *Evaluator) BootstrapLUT(ct RLWECiphertext, lut LookUpTable) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.BootstrapLUTAssign(ct, lut, ctOut)
	return ctOut
}

// BootstrapLUTAssign bootstraps LWE ciphertext with respect to given LUT and writes it to ctOut.
func (e *Evaluator) BootstrapLUTAssign(ct RLWECiphertext, lut LookUpTable, ctOut RLWECiphertext) {
	e.SampleExtractAssign(ct, e.buffer.ctExtract)
	e.KeySwitchForBootstrapAssign(e.buffer.ctExtract, e.buffer.ctKeySwitchForBootstrap)
	e.BlindRotateAssign(e.buffer.ctKeySwitchForBootstrap, lut, ctOut)
}

// ModSwitch switches the modulus of x from Q to LookUpTableSize.
func (e *Evaluator) ModSwitch(x uint64) int {
	return int(math.Round(e.modSwitchConstant*float64(x))) % e.Parameters.lookUpTableSize
}

// BlindRotate returns the blind rotation of LWE ciphertext with respect to LUT.
func (e *Evaluator) BlindRotate(ct LWECiphertext, lut LookUpTable) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.BlindRotateAssign(ct, lut, ctOut)
	return ctOut
}

// BlindRotateAssign computes the blind rotation of LWE ciphertext with respect to LUT, and writes it to ctOut.
func (e *Evaluator) BlindRotateAssign(ct LWECiphertext, lut LookUpTable, ctOut RLWECiphertext) {
	switch {
	case e.Parameters.blockSize > 1:
		e.blindRotateBlockAssign(ct, lut, ctOut)
	default:
		e.blindRotateOriginalAssign(ct, lut, ctOut)
	}
}

// blindRotateBlockAssign computes the blind rotation when PolyDegree = LookUpTableSize and BlockSize > 1.
// This is equivalent to the blind rotation algorithm using block binary keys, as explained in https://eprint.iacr.org/2023/958.
func (e *Evaluator) blindRotateBlockAssign(ct LWECiphertext, lut LookUpTable, ctOut RLWECiphertext) {
	polyDecomposed := e.Decomposer.buffer.polyDecomposed[:e.Parameters.blindRotateParameters.level]

	e.PolyEvaluator.PermutePolyAssign(lut.Value[0], e.ModSwitch(ct.Value[0]), ctOut.Value[0])
	ctOut.Value[1].Clear()

	e.Decomposer.DecomposePolyAssign(ctOut.Value[0], e.Parameters.blindRotateParameters, polyDecomposed)
	for j := 0; j < e.Parameters.blindRotateParameters.level; j++ {
		e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[j], e.buffer.ctAccFourierDecomposed[0][0][j])
	}

	e.GadgetProductFourierDecomposedFourierRLWEAssign(e.EvaluationKey.BlindRotateKey.Value[0].Value[0], e.buffer.ctAccFourierDecomposed[0][0], e.buffer.ctBlockFourierAcc[0])
	e.PolyEvaluator.NegFourierPolyAssign(e.buffer.ctBlockFourierAcc[0].Value[0], e.buffer.ctFourierAcc[0].Value[0])
	e.PolyEvaluator.NegFourierPolyAssign(e.buffer.ctBlockFourierAcc[0].Value[1], e.buffer.ctFourierAcc[0].Value[1])

	aN := e.ModSwitch(ct.Value[1])
	e.PolyEvaluator.PermuteAddFourierPolyAssign(e.buffer.ctBlockFourierAcc[0].Value[0], aN, e.buffer.ctFourierAcc[0].Value[0])
	e.PolyEvaluator.ToPolyAssignUnsafe(e.buffer.ctBlockFourierAcc[0].Value[1], e.buffer.ctBlockAcc[0].Value[1])
	e.PolyEvaluator.PermutePolyAssign(e.buffer.ctBlockAcc[0].Value[1], aN, e.buffer.pPermute)
	e.Decomposer.DecomposePolyAssign(e.buffer.pPermute, e.Parameters.blindRotateParameters, polyDecomposed)
	for k := 0; k < e.Parameters.blindRotateParameters.level; k++ {
		e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[k], e.buffer.pPermuteFourierDecomposed[k])
	}
	e.GadgetProductFourierDecomposedAddFourierRLWEAssign(e.EvaluationKey.AutomorphismKey[aN].Value, e.buffer.pPermuteFourierDecomposed, e.buffer.ctFourierAcc[0])

	for j := 1; j < e.Parameters.blockSize; j++ {
		e.GadgetProductFourierDecomposedFourierRLWEAssign(e.EvaluationKey.BlindRotateKey.Value[j].Value[0], e.buffer.ctAccFourierDecomposed[0][0], e.buffer.ctBlockFourierAcc[0])
		e.PolyEvaluator.SubFourierPolyAssign(e.buffer.ctFourierAcc[0].Value[0], e.buffer.ctBlockFourierAcc[0].Value[0], e.buffer.ctFourierAcc[0].Value[0])
		e.PolyEvaluator.SubFourierPolyAssign(e.buffer.ctFourierAcc[0].Value[1], e.buffer.ctBlockFourierAcc[0].Value[1], e.buffer.ctFourierAcc[0].Value[1])

		aN := e.ModSwitch(ct.Value[j+1])
		e.PolyEvaluator.PermuteAddFourierPolyAssign(e.buffer.ctBlockFourierAcc[0].Value[0], aN, e.buffer.ctFourierAcc[0].Value[0])
		e.PolyEvaluator.ToPolyAssignUnsafe(e.buffer.ctBlockFourierAcc[0].Value[1], e.buffer.ctBlockAcc[0].Value[1])
		e.PolyEvaluator.PermutePolyAssign(e.buffer.ctBlockAcc[0].Value[1], aN, e.buffer.pPermute)
		e.Decomposer.DecomposePolyAssign(e.buffer.pPermute, e.Parameters.blindRotateParameters, polyDecomposed)
		for k := 0; k < e.Parameters.blindRotateParameters.level; k++ {
			e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[k], e.buffer.pPermuteFourierDecomposed[k])
		}
		e.GadgetProductFourierDecomposedAddFourierRLWEAssign(e.EvaluationKey.AutomorphismKey[aN].Value, e.buffer.pPermuteFourierDecomposed, e.buffer.ctFourierAcc[0])
	}

	e.PolyEvaluator.ToPolyAddAssignUnsafe(e.buffer.ctFourierAcc[0].Value[0], ctOut.Value[0])
	e.PolyEvaluator.ToPolyAddAssignUnsafe(e.buffer.ctFourierAcc[0].Value[1], ctOut.Value[1])

	for i := 1; i < e.Parameters.blockCount; i++ {
		e.Decomposer.DecomposePolyAssign(ctOut.Value[0], e.Parameters.blindRotateParameters, polyDecomposed)
		for j := 0; j < e.Parameters.blindRotateParameters.level; j++ {
			e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[j], e.buffer.ctAccFourierDecomposed[0][0][j])
		}

		e.Decomposer.DecomposePolyAssign(ctOut.Value[1], e.Parameters.blindRotateParameters, polyDecomposed)
		for j := 0; j < e.Parameters.blindRotateParameters.level; j++ {
			e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[j], e.buffer.ctAccFourierDecomposed[0][1][j])
		}

		e.ExternalProductFourierDecomposedFourierRLWEAssign(e.EvaluationKey.BlindRotateKey.Value[i*e.Parameters.blockSize], e.buffer.ctAccFourierDecomposed[0], e.buffer.ctBlockFourierAcc[0])
		e.PolyEvaluator.NegFourierPolyAssign(e.buffer.ctBlockFourierAcc[0].Value[0], e.buffer.ctFourierAcc[0].Value[0])
		e.PolyEvaluator.NegFourierPolyAssign(e.buffer.ctBlockFourierAcc[0].Value[1], e.buffer.ctFourierAcc[0].Value[1])

		aN := e.ModSwitch(ct.Value[i*e.Parameters.blockSize+1])
		e.PolyEvaluator.PermuteAddFourierPolyAssign(e.buffer.ctBlockFourierAcc[0].Value[0], aN, e.buffer.ctFourierAcc[0].Value[0])
		e.PolyEvaluator.ToPolyAssignUnsafe(e.buffer.ctBlockFourierAcc[0].Value[1], e.buffer.ctBlockAcc[0].Value[1])
		e.PolyEvaluator.PermutePolyAssign(e.buffer.ctBlockAcc[0].Value[1], aN, e.buffer.pPermute)
		e.Decomposer.DecomposePolyAssign(e.buffer.pPermute, e.Parameters.blindRotateParameters, polyDecomposed)
		for k := 0; k < e.Parameters.blindRotateParameters.level; k++ {
			e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[k], e.buffer.pPermuteFourierDecomposed[k])
		}
		e.GadgetProductFourierDecomposedAddFourierRLWEAssign(e.EvaluationKey.AutomorphismKey[aN].Value, e.buffer.pPermuteFourierDecomposed, e.buffer.ctFourierAcc[0])

		for j := i*e.Parameters.blockSize + 1; j < (i+1)*e.Parameters.blockSize; j++ {
			e.ExternalProductFourierDecomposedFourierRLWEAssign(e.EvaluationKey.BlindRotateKey.Value[j], e.buffer.ctAccFourierDecomposed[0], e.buffer.ctBlockFourierAcc[0])
			e.PolyEvaluator.SubFourierPolyAssign(e.buffer.ctFourierAcc[0].Value[0], e.buffer.ctBlockFourierAcc[0].Value[0], e.buffer.ctFourierAcc[0].Value[0])
			e.PolyEvaluator.SubFourierPolyAssign(e.buffer.ctFourierAcc[0].Value[1], e.buffer.ctBlockFourierAcc[0].Value[1], e.buffer.ctFourierAcc[0].Value[1])

			aN := e.ModSwitch(ct.Value[j+1])
			e.PolyEvaluator.PermuteAddFourierPolyAssign(e.buffer.ctBlockFourierAcc[0].Value[0], aN, e.buffer.ctFourierAcc[0].Value[0])
			e.PolyEvaluator.ToPolyAssignUnsafe(e.buffer.ctBlockFourierAcc[0].Value[1], e.buffer.ctBlockAcc[0].Value[1])
			e.PolyEvaluator.PermutePolyAssign(e.buffer.ctBlockAcc[0].Value[1], aN, e.buffer.pPermute)
			e.Decomposer.DecomposePolyAssign(e.buffer.pPermute, e.Parameters.blindRotateParameters, polyDecomposed)
			for k := 0; k < e.Parameters.blindRotateParameters.level; k++ {
				e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[k], e.buffer.pPermuteFourierDecomposed[k])
			}
			e.GadgetProductFourierDecomposedAddFourierRLWEAssign(e.EvaluationKey.AutomorphismKey[aN].Value, e.buffer.pPermuteFourierDecomposed, e.buffer.ctFourierAcc[0])
		}

		e.PolyEvaluator.ToPolyAddAssignUnsafe(e.buffer.ctFourierAcc[0].Value[0], ctOut.Value[0])
		e.PolyEvaluator.ToPolyAddAssignUnsafe(e.buffer.ctFourierAcc[0].Value[1], ctOut.Value[1])
	}
}

// blindRotateOriginalAssign computes the blind rotation when PolyDegree = LookUpTableSize and BlockSize = 1.
// This is equivalent to the original blind rotation algorithm.
func (e *Evaluator) blindRotateOriginalAssign(ct LWECiphertext, lut LookUpTable, ctOut RLWECiphertext) {
	polyDecomposed := e.Decomposer.buffer.polyDecomposed[:e.Parameters.blindRotateParameters.level]

	e.PolyEvaluator.PermutePolyAssign(lut.Value[0], e.ModSwitch(ct.Value[0]), ctOut.Value[0])
	ctOut.Value[1].Clear()

	e.Decomposer.DecomposePolyAssign(ctOut.Value[0], e.Parameters.blindRotateParameters, polyDecomposed)
	for j := 0; j < e.Parameters.blindRotateParameters.level; j++ {
		e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[j], e.buffer.ctAccFourierDecomposed[0][0][j])
	}

	e.GadgetProductFourierDecomposedFourierRLWEAssign(e.EvaluationKey.BlindRotateKey.Value[0].Value[0], e.buffer.ctAccFourierDecomposed[0][0], e.buffer.ctBlockFourierAcc[0])
	e.PolyEvaluator.ToPolyAssignUnsafe(e.buffer.ctBlockFourierAcc[0].Value[0], e.buffer.ctBlockAcc[0].Value[0])
	e.PolyEvaluator.ToPolyAssignUnsafe(e.buffer.ctBlockFourierAcc[0].Value[1], e.buffer.ctBlockAcc[0].Value[1])

	e.PolyEvaluator.SubPolyAssign(ctOut.Value[0], e.buffer.ctBlockAcc[0].Value[0], ctOut.Value[0])
	e.PolyEvaluator.SubPolyAssign(ctOut.Value[1], e.buffer.ctBlockAcc[0].Value[1], ctOut.Value[1])

	aN := e.ModSwitch(ct.Value[1])
	e.PolyEvaluator.PermuteAddPolyAssign(e.buffer.ctBlockAcc[0].Value[0], aN, ctOut.Value[0])
	e.PolyEvaluator.PermutePolyAssign(e.buffer.ctBlockAcc[0].Value[1], aN, e.buffer.pPermute)
	e.Decomposer.DecomposePolyAssign(e.buffer.pPermute, e.Parameters.blindRotateParameters, polyDecomposed)
	for j := 0; j < e.Parameters.blindRotateParameters.level; j++ {
		e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[j], e.buffer.pPermuteFourierDecomposed[j])
	}
	e.GadgetProductFourierDecomposedFourierRLWEAssign(e.EvaluationKey.AutomorphismKey[aN].Value, e.buffer.pPermuteFourierDecomposed, e.buffer.ctFourierAcc[0])
	e.PolyEvaluator.ToPolyAddAssignUnsafe(e.buffer.ctFourierAcc[0].Value[0], ctOut.Value[0])
	e.PolyEvaluator.ToPolyAddAssignUnsafe(e.buffer.ctFourierAcc[0].Value[1], ctOut.Value[1])

	for i := 1; i < e.Parameters.lweDimension; i++ {
		e.Decomposer.DecomposePolyAssign(ctOut.Value[0], e.Parameters.blindRotateParameters, polyDecomposed)
		for j := 0; j < e.Parameters.blindRotateParameters.level; j++ {
			e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[j], e.buffer.ctAccFourierDecomposed[0][0][j])
		}

		e.Decomposer.DecomposePolyAssign(ctOut.Value[1], e.Parameters.blindRotateParameters, polyDecomposed)
		for j := 0; j < e.Parameters.blindRotateParameters.level; j++ {
			e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[j], e.buffer.ctAccFourierDecomposed[0][1][j])
		}

		e.ExternalProductFourierDecomposedFourierRLWEAssign(e.EvaluationKey.BlindRotateKey.Value[i], e.buffer.ctAccFourierDecomposed[0], e.buffer.ctBlockFourierAcc[0])
		e.PolyEvaluator.ToPolyAssignUnsafe(e.buffer.ctBlockFourierAcc[0].Value[0], e.buffer.ctBlockAcc[0].Value[0])
		e.PolyEvaluator.ToPolyAssignUnsafe(e.buffer.ctBlockFourierAcc[0].Value[1], e.buffer.ctBlockAcc[0].Value[1])

		e.PolyEvaluator.SubPolyAssign(ctOut.Value[0], e.buffer.ctBlockAcc[0].Value[0], ctOut.Value[0])
		e.PolyEvaluator.SubPolyAssign(ctOut.Value[1], e.buffer.ctBlockAcc[0].Value[1], ctOut.Value[1])

		aN := e.ModSwitch(ct.Value[i+1])
		e.PolyEvaluator.PermuteAddPolyAssign(e.buffer.ctBlockAcc[0].Value[0], aN, ctOut.Value[0])
		e.PolyEvaluator.PermutePolyAssign(e.buffer.ctBlockAcc[0].Value[1], aN, e.buffer.pPermute)
		e.Decomposer.DecomposePolyAssign(e.buffer.pPermute, e.Parameters.blindRotateParameters, polyDecomposed)
		for j := 0; j < e.Parameters.blindRotateParameters.level; j++ {
			e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[j], e.buffer.pPermuteFourierDecomposed[j])
		}
		e.GadgetProductFourierDecomposedFourierRLWEAssign(e.EvaluationKey.AutomorphismKey[aN].Value, e.buffer.pPermuteFourierDecomposed, e.buffer.ctFourierAcc[0])
		e.PolyEvaluator.ToPolyAddAssignUnsafe(e.buffer.ctFourierAcc[0].Value[0], ctOut.Value[0])
		e.PolyEvaluator.ToPolyAddAssignUnsafe(e.buffer.ctFourierAcc[0].Value[1], ctOut.Value[1])
	}
}

// SampleExtract extracts a LWE ciphertext from RLWE ciphertext.
// Output ciphertext will be of length PolyDegree + 1.
func (e *Evaluator) SampleExtract(ct RLWECiphertext) LWECiphertext {
	ctOut := NewLWECiphertextCustom(e.Parameters.polyDegree)
	e.SampleExtractAssign(ct, ctOut)
	return ctOut
}

// SampleExtractAssign extracts a LWE ciphertext from RLWE ciphertext.
// Output ciphertext should be of length PolyDegree + 1.
func (e *Evaluator) SampleExtractAssign(ct RLWECiphertext, ctOut LWECiphertext) {
	switch e.Parameters.encodeType {
	case EncodeTypeCoeffs:
		e.sampleExtractAssign(ct, ctOut)
	case EncodeTypeSlots:
		e.PolyEvaluator.ShortFourierPolyMulPolyAssign(ct.Value[0], e.fMask, e.buffer.ctMask.Value[0])
		e.PolyEvaluator.ShortFourierPolyMulPolyAssign(ct.Value[1], e.fMask, e.buffer.ctMask.Value[1])
		e.sampleExtractAssign(e.buffer.ctMask, ctOut)
	}
}

// sampleExtractAssign extracts a LWE ciphertext from RLWE ciphertext.
// Output ciphertext should be of length PolyDegree + 1.
func (e *Evaluator) sampleExtractAssign(ct RLWECiphertext, ctOut LWECiphertext) {
	ctOut.Value[0] = ct.Value[0].Coeffs[0]

	for i := 0; i < e.Parameters.polyDegree; i++ {
		ctOut.Value[1+i] = -uint64(e.PolyEvaluator.Parameters.Order()) * ct.Value[1].Coeffs[i]
	}

	for i := 1; i < e.PolyEvaluator.Parameters.CyclotomicDegree()-1; i++ {
		ctOut.Value[1+e.extractIdx[i]] += ct.Value[1].Coeffs[e.extractIdx[e.PolyEvaluator.Parameters.CyclotomicDegree()-i-1]]
	}
}

// KeySwitchForBootstrap performs the keyswitching using evaulater's evaluation key.
// Input ciphertext should be of length PolyDegree + 1.
// Output ciphertext will be of length LWEDimension + 1.
func (e *Evaluator) KeySwitchForBootstrap(ct LWECiphertext) LWECiphertext {
	ctOut := NewLWECiphertext(e.Parameters)
	e.KeySwitchForBootstrapAssign(ct, ctOut)
	return ctOut
}

// KeySwitchForBootstrapAssign performs the keyswitching using evaulater's evaluation key.
// Input ciphertext should be of length PolyDegree + 1.
// Output ciphertext should be of length LWEDimension + 1.
func (e *Evaluator) KeySwitchForBootstrapAssign(ct, ctOut LWECiphertext) {
	scalarDecomposed := e.Decomposer.buffer.scalarDecomposed[:e.Parameters.keySwitchParameters.level]

	vec.CopyAssign(ct.Value[:e.Parameters.lweDimension+1], ctOut.Value)
	for i, ii := e.Parameters.lweDimension, 0; i < e.Parameters.polyDegree; i, ii = i+1, ii+1 {
		e.Decomposer.DecomposeScalarAssign(ct.Value[i+1], e.Parameters.keySwitchParameters, scalarDecomposed)
		for j := 0; j < e.Parameters.keySwitchParameters.level; j++ {
			e.ScalarMulAddLWEAssign(e.EvaluationKey.KeySwitchKey.Value[ii].Value[j], scalarDecomposed[j], ctOut)
		}
	}
}
