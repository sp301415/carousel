package carousel

import (
	"math"

	"github.com/sp301415/carousel/math/num"
	"github.com/sp301415/carousel/math/poly"
)

// Evaluator evaluates homomorphic operations on ciphertexts.
// This is meant to be public, usually for servers.
//
// Evaluator is not safe for concurrent use.
// Use [*Evaluator.ShallowCopy] to get a safe copy.
type Evaluator struct {
	// Encoder is an embedded encoder for this Evaluator.
	*Encoder
	// RLWETransformer is an embedded RLWETransformer for this Evaluator.
	*RLWETransformer

	// Parameters is the parameters for this Evaluator.
	Parameters Parameters

	// Decomposer is an Decomposer for this Evaluator.
	Decomposer *Decomposer
	// PolyEvaluator is a PolyEvaluator for this Evaluator.
	PolyEvaluator *poly.Evaluator

	// EvaluationKey is the evaluation key for this Evaluator.
	EvaluationKey EvaluationKey

	// modSwitchConstant is a constant for modulus switching.
	modSwitchConstant float64
	// mask is the masking polynomial for sample extraction.
	fMask poly.FourierPoly
	// extractIdx is an index for sample extraction.
	extractIdx []int

	buffer evaluationBuffer
}

// evaluationBuffer is a buffer for Evaluator.
type evaluationBuffer struct {
	// fpMul is the fourier transformed polynomial for multiplications.
	fpMul poly.FourierPoly

	// ctProdLWE is the LWE ciphertext buffer for ExternalProductLWE and KeySwitchLWE.
	ctProdLWE LWECiphertext
	// ctProdFourierRLWE is the fourier transformed ctRLWEOut in ExternalProductRLWE and KeySwitchRLWE.
	ctProdFourierRLWE FourierRLWECiphertext

	// ctAcc is the accumulator in BlindRotateExtended.
	// This has length PolyExtendFactor.
	ctAcc []RLWECiphertext
	// ctFourierAcc is the fourier transformed accumulator in Blind Rotation.
	// In case of BlindRotateBlock and BlindRotateOriginal, only the first element is used.
	// This has length PolyExtendFactor.
	ctFourierAcc []FourierRLWECiphertext
	// ctBlockFourierAcc is the auxiliary accumulator in BlindRotateBlock and BlindRotateExtended.
	ctBlockFourierAcc []FourierRLWECiphertext
	// ctBlockAcc is the auxiliary accumulator in BlindRotateBlock and BlindRotateExtended.
	ctBlockAcc []RLWECiphertext
	// ctAccFourierDecomposed is the decomposed ctAcc in Blind Rotation.
	// In case of BlindRotateBlock and BlindRotateOriginal, only the first element is used.
	// This has length PolyExpandFactor.
	ctAccFourierDecomposed [][][]poly.FourierPoly
	// pPermute is the permuted ciphertext during Blind Rotation.
	pPermtue poly.Poly
	// pPermuteFourierDecomposed is the fourier transformed pPermute during Blind Rotation.
	pPermtueFourierDecomposed []poly.FourierPoly

	// ctRotate is the blind rotated RLWE ciphertext for bootstrapping.
	ctRotate RLWECiphertext
	// ctMask is the masked RLWE ciphertext for Sample Extraction.
	ctMask RLWECiphertext
	// ctExtract is the extracted LWE ciphertext after Blind Rotation.
	ctExtract LWECiphertext
	// ctKeySwitchForBootstrap is the LWEDimension sized ciphertext from keyswitching for bootstrapping.
	ctKeySwitchForBootstrap LWECiphertext

	// lut is an empty lut, used for BlindRotateFunc.
	lut LookUpTable
	// lutRaw is an full-sized LUT.
	lutRaw [][]int
	// lutReorder is the reordering buffer for LUT.
	lutReorder []int
}

// NewEvaluator creates a new Evaluator based on parameters.
// This does not copy evaluation keys, since they may be large.
func NewEvaluator(params Parameters, evk EvaluationKey) *Evaluator {
	decomposer := NewDecomposer(params.polyEvaluatorParameters)
	decomposer.ScalarDecomposedBuffer(params.keySwitchParameters)
	decomposer.PolyDecomposedBuffer(params.blindRotateParameters)
	decomposer.PolyFourierDecomposedBuffer(params.blindRotateParameters)

	ecd := NewEncoder(params)
	polyEvaluator := poly.NewEvaluator(params.polyEvaluatorParameters)
	fMask := poly.NewFourierPoly(params.polyDegree)
	if params.encodeType == EncodeTypeSlots {
		maskConst := num.ModInverse(params.polyEvaluatorParameters.Resolution()[0], params.messageModulus)
		mask := ecd.EncodeRLWESlots([]int{int(maskConst)})
		for i := 0; i < params.polyDegree; i++ {
			mask.Value.Coeffs[i] = num.DivRound(mask.Value.Coeffs[i], params.scale)
		}
		polyEvaluator.ToFourierPolyAssign(mask.Value, fMask)
	}

	extractIdx := make([]int, params.polyEvaluatorParameters.CyclotomicDegree())
	for i := 0; i < params.polyEvaluatorParameters.CyclotomicDegree(); i++ {
		gIdx := num.ModExp(params.polyEvaluatorParameters.Generator(), uint64(i), uint64(params.polyEvaluatorParameters.CyclotomicDegree()))
		extractIdx[gIdx] = i % params.polyEvaluatorParameters.Degree()
	}

	return &Evaluator{
		Encoder:         ecd,
		RLWETransformer: NewRLWETransformer(params.polyEvaluatorParameters),

		Parameters: params,

		Decomposer:    decomposer,
		PolyEvaluator: polyEvaluator,

		EvaluationKey: evk,

		modSwitchConstant: float64(params.polyExtendFactor) / math.Exp2(float64(64-params.logPolyDegree)),
		fMask:             fMask,
		extractIdx:        extractIdx,

		buffer: newEvaluationBuffer(params),
	}
}

// newEvaluationBuffer creates a new evaluationBuffer.
func newEvaluationBuffer(params Parameters) evaluationBuffer {
	ctAcc := make([]RLWECiphertext, params.polyExtendFactor)
	ctFourierAcc := make([]FourierRLWECiphertext, params.polyExtendFactor)
	ctBlockFourierAcc := make([]FourierRLWECiphertext, params.polyExtendFactor)
	ctBlockAcc := make([]RLWECiphertext, params.polyExtendFactor)
	for i := 0; i < params.polyExtendFactor; i++ {
		ctAcc[i] = NewRLWECiphertext(params)
		ctFourierAcc[i] = NewFourierRLWECiphertext(params)
		ctBlockFourierAcc[i] = NewFourierRLWECiphertext(params)
		ctBlockAcc[i] = NewRLWECiphertext(params)
	}

	ctAccFourierDecomposed := make([][][]poly.FourierPoly, params.polyExtendFactor)
	for i := 0; i < params.polyExtendFactor; i++ {
		ctAccFourierDecomposed[i] = make([][]poly.FourierPoly, 2)
		for j := 0; j < 2; j++ {
			ctAccFourierDecomposed[i][j] = make([]poly.FourierPoly, params.blindRotateParameters.level)
			for k := 0; k < params.blindRotateParameters.level; k++ {
				ctAccFourierDecomposed[i][j][k] = poly.NewFourierPoly(params.polyDegree)
			}
		}
	}

	pPermuteFourierDecomposed := make([]poly.FourierPoly, params.blindRotateParameters.level)
	for i := 0; i < params.blindRotateParameters.level; i++ {
		pPermuteFourierDecomposed[i] = poly.NewFourierPoly(params.polyDegree)
	}

	lutRaw := make([][]int, params.polyExtendFactor)
	for i := 0; i < params.polyExtendFactor; i++ {
		lutRaw[i] = make([]int, params.polyDegree)
	}

	return evaluationBuffer{
		fpMul: poly.NewFourierPoly(params.polyDegree),

		ctProdLWE:         NewLWECiphertext(params),
		ctProdFourierRLWE: NewFourierRLWECiphertext(params),

		ctAcc:                     ctAcc,
		ctFourierAcc:              ctFourierAcc,
		ctBlockFourierAcc:         ctBlockFourierAcc,
		ctBlockAcc:                ctBlockAcc,
		ctAccFourierDecomposed:    ctAccFourierDecomposed,
		pPermtue:                  poly.NewPoly(params.polyDegree),
		pPermtueFourierDecomposed: pPermuteFourierDecomposed,

		ctRotate:                NewRLWECiphertext(params),
		ctMask:                  NewRLWECiphertext(params),
		ctExtract:               NewLWECiphertextCustom(params.polyDegree),
		ctKeySwitchForBootstrap: NewLWECiphertextCustom(params.lweDimension),

		lut:        NewLookUpTable(params),
		lutRaw:     lutRaw,
		lutReorder: make([]int, params.lookUpTableSize),
	}
}

// ShallowCopy returns a shallow copy of this Evaluator.
// Returned Evaluator is safe for concurrent use.
func (e *Evaluator) ShallowCopy() *Evaluator {
	return &Evaluator{
		Encoder:         e.Encoder,
		RLWETransformer: e.RLWETransformer.ShallowCopy(),

		Parameters: e.Parameters,

		Decomposer:    e.Decomposer.ShallowCopy(),
		PolyEvaluator: e.PolyEvaluator.ShallowCopy(),

		EvaluationKey: e.EvaluationKey,

		modSwitchConstant: e.modSwitchConstant,

		buffer: newEvaluationBuffer(e.Parameters),
	}
}
