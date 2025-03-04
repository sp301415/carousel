package carousel

import "github.com/sp301415/carousel/math/poly"

// RLWETransformer is a struct for transforming
// RLWE entities to Fourier RLWE entities and vice versa.
// RLWETransformer is embedded in Encryptor and Evaluator,
// so usually manual instantiation isn't needed.
//
// RLWETransformer is not safe for concurrent use.
// Use [*RLWETransformer.ShallowCopy] to get a safe copy.
type RLWETransformer struct {
	// PolyEvaluator is a PolyEvaluator for this RLWETransformer.
	PolyEvaluator *poly.Evaluator
}

// NewRLWETransformer returns a new RLWETransformer with given parameters.
func NewRLWETransformer(params poly.EvaluatorParameters) *RLWETransformer {
	return &RLWETransformer{
		PolyEvaluator: poly.NewEvaluator(params),
	}
}

// ShallowCopy returns a shallow copy of this RLWETransformer.
// Returned RLWETransformer is safe for concurrent use.
func (e *RLWETransformer) ShallowCopy() *RLWETransformer {
	return &RLWETransformer{
		PolyEvaluator: e.PolyEvaluator.ShallowCopy(),
	}
}

// ToFourierRLWESecretKey transforms RLWE secret key to Fourier RLWE secret key.
func (e *RLWETransformer) ToFourierRLWESecretKey(sk RLWESecretKey) FourierRLWESecretKey {
	skOut := NewFourierRLWESecretKeyCustom(e.PolyEvaluator.Parameters.Degree())
	e.ToFourierRLWESecretKeyAssign(sk, skOut)
	return skOut
}

// ToFourierRLWESecretKeyAssign transforms RLWE secret key to Fourier RLWE secret key and writes it to skOut.
func (e *RLWETransformer) ToFourierRLWESecretKeyAssign(sk RLWESecretKey, skOut FourierRLWESecretKey) {
	e.PolyEvaluator.ToFourierPolyAssign(sk.Value, skOut.Value)
}

// ToRLWESecretKey transforms Fourier RLWE secret key to RLWE secret key.
func (e *RLWETransformer) ToRLWESecretKey(sk FourierRLWESecretKey) RLWESecretKey {
	skOut := NewRLWESecretKeyCustom(e.PolyEvaluator.Parameters.Degree())
	e.ToRLWESecretKeyAssign(sk, skOut)
	return skOut
}

// ToRLWESecretKeyAssign transforms Fourier RLWE secret key to RLWE secret key and writes it to skOut.
func (e *RLWETransformer) ToRLWESecretKeyAssign(sk FourierRLWESecretKey, skOut RLWESecretKey) {
	e.PolyEvaluator.ToPolyAssign(sk.Value, skOut.Value)
}

// ToFourierRLWECiphertext transforms RLWE ciphertext to Fourier RLWE ciphertext.
func (e *RLWETransformer) ToFourierRLWECiphertext(ct RLWECiphertext) FourierRLWECiphertext {
	ctOut := NewFourierRLWECiphertextCustom(len(ct.Value)-1, e.PolyEvaluator.Parameters.Degree())
	e.ToFourierRLWECiphertextAssign(ct, ctOut)
	return ctOut
}

// ToFourierRLWECiphertextAssign transforms RLWE ciphertext to Fourier RLWE ciphertext and writes it to ctOut.
func (e *RLWETransformer) ToFourierRLWECiphertextAssign(ct RLWECiphertext, ctOut FourierRLWECiphertext) {
	for i := 0; i < len(ct.Value); i++ {
		e.PolyEvaluator.ToFourierPolyAssign(ct.Value[i], ctOut.Value[i])
	}
}

// ToRLWECiphertext transforms Fourier RLWE ciphertext to RLWE ciphertext.
func (e *RLWETransformer) ToRLWECiphertext(ct FourierRLWECiphertext) RLWECiphertext {
	ctOut := NewRLWECiphertextCustom(len(ct.Value)-1, e.PolyEvaluator.Parameters.Degree())
	e.ToRLWECiphertextAssign(ct, ctOut)
	return ctOut
}

// ToRLWECiphertextAssign transforms Fourier RLWE ciphertext to RLWE ciphertext and writes it to ctOut.
func (e *RLWETransformer) ToRLWECiphertextAssign(ct FourierRLWECiphertext, ctOut RLWECiphertext) {
	for i := 0; i < len(ct.Value); i++ {
		e.PolyEvaluator.ToPolyAssign(ct.Value[i], ctOut.Value[i])
	}
}

// ToFourierRLevCiphertext transforms RLev ciphertext to Fourier RLev ciphertext.
func (e *RLWETransformer) ToFourierRLevCiphertext(ct RLevCiphertext) FourierRLevCiphertext {
	ctOut := NewFourierRLevCiphertextCustom(e.PolyEvaluator.Parameters.Degree(), ct.GadgetParameters)
	e.ToFourierRLevCiphertextAssign(ct, ctOut)
	return ctOut
}

// ToFourierRLevCiphertextAssign transforms RLev ciphertext to Fourier RLev ciphertext and writes it to ctOut.
func (e *RLWETransformer) ToFourierRLevCiphertextAssign(ct RLevCiphertext, ctOut FourierRLevCiphertext) {
	for i := 0; i < ct.GadgetParameters.level; i++ {
		e.ToFourierRLWECiphertextAssign(ct.Value[i], ctOut.Value[i])
	}
}

// ToRLevCiphertext transforms Fourier RLev ciphertext to RLev ciphertext.
func (e *RLWETransformer) ToRLevCiphertext(ct FourierRLevCiphertext) RLevCiphertext {
	ctOut := NewRLevCiphertextCustom(e.PolyEvaluator.Parameters.Degree(), ct.GadgetParameters)
	e.ToRLevCiphertextAssign(ct, ctOut)
	return ctOut
}

// ToRLevCiphertextAssign transforms Fourier RLev ciphertext to RLev ciphertext and writes it to ctOut.
func (e *RLWETransformer) ToRLevCiphertextAssign(ct FourierRLevCiphertext, ctOut RLevCiphertext) {
	for i := 0; i < ct.GadgetParameters.level; i++ {
		e.ToRLWECiphertextAssign(ct.Value[i], ctOut.Value[i])
	}
}

// ToFourierRGSWCiphertext transforms RGSW ciphertext to Fourier RGSW ciphertext.
func (e *RLWETransformer) ToFourierRGSWCiphertext(ct RGSWCiphertext) FourierRGSWCiphertext {
	ctOut := NewFourierRGSWCiphertextCustom(e.PolyEvaluator.Parameters.Degree(), ct.GadgetParameters)
	e.ToFourierRGSWCiphertextAssign(ct, ctOut)
	return ctOut
}

// ToFourierRGSWCiphertextAssign transforms RGSW ciphertext to Fourier RGSW ciphertext and writes it to ctOut.
func (e *RLWETransformer) ToFourierRGSWCiphertextAssign(ct RGSWCiphertext, ctOut FourierRGSWCiphertext) {
	for i := 0; i < len(ct.Value); i++ {
		e.ToFourierRLevCiphertextAssign(ct.Value[i], ctOut.Value[i])
	}
}

// ToRGSWCiphertext transforms Fourier RGSW ciphertext to RGSW ciphertext.
func (e *RLWETransformer) ToRGSWCiphertext(ct FourierRGSWCiphertext) RGSWCiphertext {
	ctOut := NewRGSWCiphertextCustom(e.PolyEvaluator.Parameters.Degree(), ct.GadgetParameters)
	e.ToRGSWCiphertextAssign(ct, ctOut)
	return ctOut
}

// ToRGSWCiphertextAssign transforms Fourier RGSW ciphertext to RGSW ciphertext and writes it to ctOut.
func (e *RLWETransformer) ToRGSWCiphertextAssign(ct FourierRGSWCiphertext, ctOut RGSWCiphertext) {
	for i := 0; i < len(ct.Value); i++ {
		e.ToRLevCiphertextAssign(ct.Value[i], ctOut.Value[i])
	}
}
