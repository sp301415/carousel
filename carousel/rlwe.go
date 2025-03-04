package carousel

import (
	"github.com/sp301415/carousel/math/poly"
)

// RLWESecretKey is a RLWE secret key, sampled from uniform binary distribution.
type RLWESecretKey struct {
	// Value is a single polynomial.
	Value poly.Poly
}

// NewRLWESecretKey creates a new RLWESecretKey.
func NewRLWESecretKey(params Parameters) RLWESecretKey {
	return RLWESecretKey{Value: poly.NewPoly(params.polyDegree)}
}

// NewRLWESecretKeyCustom creates a new RLWESecretKey with given polyDegree.
func NewRLWESecretKeyCustom(polyDegree int) RLWESecretKey {
	return RLWESecretKey{Value: poly.NewPoly(polyDegree)}
}

// Copy returns a copy of the key.
func (sk RLWESecretKey) Copy() RLWESecretKey {
	return RLWESecretKey{Value: sk.Value.Copy()}
}

// CopyFrom copies values from the key.
func (sk *RLWESecretKey) CopyFrom(skIn RLWESecretKey) {
	sk.Value.CopyFrom(skIn.Value)
}

// Clear clears the key.
func (sk *RLWESecretKey) Clear() {
	sk.Value.Clear()
}

// RLWEPlaintext represents an encoded RLWE plaintext.
type RLWEPlaintext struct {
	// Value is a single polynomial.
	Value poly.Poly
}

// NewRLWEPlaintext creates a new RLWEPlaintext.
func NewRLWEPlaintext(params Parameters) RLWEPlaintext {
	return RLWEPlaintext{Value: poly.NewPoly(params.polyDegree)}
}

// NewRLWEPlaintextCustom creates a new RLWEPlaintext with given polyDegree.
func NewRLWEPlaintextCustom(polyDegree int) RLWEPlaintext {
	return RLWEPlaintext{Value: poly.NewPoly(polyDegree)}
}

// Copy returns a copy of the plaintext.
func (pt RLWEPlaintext) Copy() RLWEPlaintext {
	return RLWEPlaintext{Value: pt.Value.Copy()}
}

// CopyFrom copies values from the plaintext.
func (pt *RLWEPlaintext) CopyFrom(ptIn RLWEPlaintext) {
	pt.Value.CopyFrom(ptIn.Value)
}

// Clear clears the plaintext.
func (pt *RLWEPlaintext) Clear() {
	pt.Value.Clear()
}

// RLWECiphertext represents an encrypted RLWE ciphertext.
type RLWECiphertext struct {
	// Value is ordered as [body, mask],
	// since Go doesn't provide an easy way to take last element of slice.
	// Therefore, value has length 2.
	Value []poly.Poly
}

// NewRLWECiphertext creates a new RLWECiphertext.
func NewRLWECiphertext(params Parameters) RLWECiphertext {
	ct := make([]poly.Poly, 2)
	ct[0] = poly.NewPoly(params.polyDegree)
	ct[1] = poly.NewPoly(params.polyDegree)
	return RLWECiphertext{Value: ct}
}

// NewRLWECiphertextCustom creates a new RLWECiphertext with given rank and polyDegree.
func NewRLWECiphertextCustom(rank, polyDegree int) RLWECiphertext {
	ct := make([]poly.Poly, rank+1)
	for i := 0; i < rank+1; i++ {
		ct[i] = poly.NewPoly(polyDegree)
	}
	return RLWECiphertext{Value: ct}
}

// Copy returns a copy of the ciphertext.
func (ct RLWECiphertext) Copy() RLWECiphertext {
	ctCopy := make([]poly.Poly, len(ct.Value))
	for i := range ct.Value {
		ctCopy[i] = ct.Value[i].Copy()
	}
	return RLWECiphertext{Value: ctCopy}
}

// CopyFrom copies values from the ciphertext.
func (ct *RLWECiphertext) CopyFrom(ctIn RLWECiphertext) {
	for i := range ct.Value {
		ct.Value[i].CopyFrom(ctIn.Value[i])
	}
}

// Clear clears the ciphertext.
func (ct *RLWECiphertext) Clear() {
	for i := range ct.Value {
		ct.Value[i].Clear()
	}
}

// RLevCiphertext is a leveled RLWE ciphertext, decomposed according to GadgetParameters.
type RLevCiphertext struct {
	GadgetParameters GadgetParameters

	// Value has length Level.
	Value []RLWECiphertext
}

// NewRLevCiphertext creates a new RLevCiphertext.
func NewRLevCiphertext(params Parameters, gadgetParams GadgetParameters) RLevCiphertext {
	ct := make([]RLWECiphertext, gadgetParams.level)
	for i := 0; i < gadgetParams.level; i++ {
		ct[i] = NewRLWECiphertext(params)
	}
	return RLevCiphertext{Value: ct, GadgetParameters: gadgetParams}
}

// NewRLevCiphertextCustom creates a new RLevCiphertext with given polyDegree.
func NewRLevCiphertextCustom(polyDegree int, gadgetParams GadgetParameters) RLevCiphertext {
	ct := make([]RLWECiphertext, gadgetParams.level)
	for i := 0; i < gadgetParams.level; i++ {
		ct[i] = NewRLWECiphertextCustom(1, polyDegree)
	}
	return RLevCiphertext{Value: ct, GadgetParameters: gadgetParams}
}

// Copy returns a copy of the ciphertext.
func (ct RLevCiphertext) Copy() RLevCiphertext {
	ctCopy := make([]RLWECiphertext, len(ct.Value))
	for i := range ct.Value {
		ctCopy[i] = ct.Value[i].Copy()
	}
	return RLevCiphertext{Value: ctCopy, GadgetParameters: ct.GadgetParameters}
}

// CopyFrom copies values from ciphertext.
func (ct *RLevCiphertext) CopyFrom(ctIn RLevCiphertext) {
	for i := range ct.Value {
		ct.Value[i].CopyFrom(ctIn.Value[i])
	}
	ct.GadgetParameters = ctIn.GadgetParameters
}

// Clear clears the ciphertext.
func (ct *RLevCiphertext) Clear() {
	for i := range ct.Value {
		ct.Value[i].Clear()
	}
}

// RGSWCiphertext represents an encrypted RGSW ciphertext,
// which is a collection of RLev ciphertexts.
type RGSWCiphertext struct {
	GadgetParameters GadgetParameters

	// Value has length 2.
	Value []RLevCiphertext
}

// NewRGSWCiphertext creates a new RGSW ciphertext.
func NewRGSWCiphertext(params Parameters, gadgetParams GadgetParameters) RGSWCiphertext {
	ct := make([]RLevCiphertext, 2)
	ct[0] = NewRLevCiphertext(params, gadgetParams)
	ct[1] = NewRLevCiphertext(params, gadgetParams)
	return RGSWCiphertext{Value: ct, GadgetParameters: gadgetParams}
}

// NewRGSWCiphertextCustom creates a new RGSW ciphertext with given polyDegree.
func NewRGSWCiphertextCustom(polyDegree int, gadgetParams GadgetParameters) RGSWCiphertext {
	ct := make([]RLevCiphertext, 2)
	ct[0] = NewRLevCiphertextCustom(polyDegree, gadgetParams)
	ct[1] = NewRLevCiphertextCustom(polyDegree, gadgetParams)
	return RGSWCiphertext{Value: ct, GadgetParameters: gadgetParams}
}

// Copy returns a copy of the ciphertext.
func (ct RGSWCiphertext) Copy() RGSWCiphertext {
	ctCopy := make([]RLevCiphertext, len(ct.Value))
	for i := range ct.Value {
		ctCopy[i] = ct.Value[i].Copy()
	}
	return RGSWCiphertext{Value: ctCopy, GadgetParameters: ct.GadgetParameters}
}

// CopyFrom copies values from the ciphertext.
func (ct *RGSWCiphertext) CopyFrom(ctIn RGSWCiphertext) {
	for i := range ct.Value {
		ct.Value[i].CopyFrom(ctIn.Value[i])
	}
	ct.GadgetParameters = ctIn.GadgetParameters
}

// Clear clears the ciphertext.
func (ct *RGSWCiphertext) Clear() {
	for i := range ct.Value {
		ct.Value[i].Clear()
	}
}
