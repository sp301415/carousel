package carousel

import "github.com/sp301415/carousel/math/vec"

// LWESecretKey is a LWE secret key, sampled from uniform or block binary distribution.
type LWESecretKey struct {
	// Value has length DefaultLWEDimension.
	Value []uint64
}

// NewLWESecretKey creates a new LWESecretKey.
func NewLWESecretKey(params Parameters) LWESecretKey {
	return LWESecretKey{Value: make([]uint64, params.lweDimension)}
}

// NewLWESecretKeyCustom creates a new LWESecretKey with given dimension.
func NewLWESecretKeyCustom(lweDimension int) LWESecretKey {
	return LWESecretKey{Value: make([]uint64, lweDimension)}
}

// Copy returns a copy of the key.
func (sk LWESecretKey) Copy() LWESecretKey {
	return LWESecretKey{Value: vec.Copy(sk.Value)}
}

// CopyFrom copies values from the key.
func (sk *LWESecretKey) CopyFrom(skIn LWESecretKey) {
	vec.CopyAssign(skIn.Value, sk.Value)
}

// Clear clears the key.
func (sk *LWESecretKey) Clear() {
	vec.Fill(sk.Value, 0)
}

// LWEPlaintext represents an encoded LWE plaintext.
type LWEPlaintext struct {
	// Value is a scalar.
	Value uint64
}

// NewLWEPlaintext creates a new LWEPlaintext.
func NewLWEPlaintext() LWEPlaintext {
	return LWEPlaintext{}
}

// Copy returns a copy of the plaintext.
func (pt LWEPlaintext) Copy() LWEPlaintext {
	return LWEPlaintext{Value: pt.Value}
}

// CopyFrom copies values from the plaintext.
func (pt *LWEPlaintext) CopyFrom(ptIn LWEPlaintext) {
	pt.Value = ptIn.Value
}

// Clear clears the plaintext.
func (pt *LWEPlaintext) Clear() {
	pt.Value = 0
}

// LWECiphertext represents an encrypted LWE ciphertext.
//
// LWE ciphertexts are the default encrypted form of the ciphertext.
type LWECiphertext struct {
	// Value is ordered as [body, mask],
	// since Go doesn't provide an easy way to take last element of slice.
	// Therefore, value has length DefaultLWEDimension + 1.
	Value []uint64
}

// NewLWECiphertext creates a new LWECiphertext.
func NewLWECiphertext(params Parameters) LWECiphertext {
	return LWECiphertext{Value: make([]uint64, params.lweDimension+1)}
}

// NewLWECiphertextCustom creates a new LWECiphertext with given dimension.
func NewLWECiphertextCustom(lweDimension int) LWECiphertext {
	return LWECiphertext{Value: make([]uint64, lweDimension+1)}
}

// Copy returns a copy of the ciphertext.
func (ct LWECiphertext) Copy() LWECiphertext {
	return LWECiphertext{Value: vec.Copy(ct.Value)}
}

// CopyFrom copies values from the ciphertext.
func (ct *LWECiphertext) CopyFrom(ctIn LWECiphertext) {
	vec.CopyAssign(ctIn.Value, ct.Value)
}

// Clear clears the ciphertext.
func (ct *LWECiphertext) Clear() {
	vec.Fill(ct.Value, 0)
}

// LevCiphertext is a leveled LWE ciphertext, decomposed according to GadgetParameters.
type LevCiphertext struct {
	GadgetParameters GadgetParameters

	// Value has length Level.
	Value []LWECiphertext
}

// NewLevCiphertext creates a new LevCiphertext.
func NewLevCiphertext(params Parameters, gadgetParams GadgetParameters) LevCiphertext {
	ct := make([]LWECiphertext, gadgetParams.level)
	for i := 0; i < gadgetParams.level; i++ {
		ct[i] = NewLWECiphertext(params)
	}
	return LevCiphertext{Value: ct, GadgetParameters: gadgetParams}
}

// NewLevCiphertextCustom creates a new LevCiphertext with given dimension.
func NewLevCiphertextCustom(lweDimension int, gadgetParams GadgetParameters) LevCiphertext {
	ct := make([]LWECiphertext, gadgetParams.level)
	for i := 0; i < gadgetParams.level; i++ {
		ct[i] = NewLWECiphertextCustom(lweDimension)
	}
	return LevCiphertext{Value: ct, GadgetParameters: gadgetParams}
}

// Copy returns a copy of the ciphertext.
func (ct LevCiphertext) Copy() LevCiphertext {
	ctCopy := make([]LWECiphertext, len(ct.Value))
	for i := range ct.Value {
		ctCopy[i] = ct.Value[i].Copy()
	}
	return LevCiphertext{Value: ctCopy, GadgetParameters: ct.GadgetParameters}
}

// CopyFrom copies values from the ciphertext.
func (ct *LevCiphertext) CopyFrom(ctIn LevCiphertext) {
	for i := range ct.Value {
		ct.Value[i].CopyFrom(ctIn.Value[i])
	}
	ct.GadgetParameters = ctIn.GadgetParameters
}

// Clear clears the ciphertext.
func (ct *LevCiphertext) Clear() {
	for i := range ct.Value {
		ct.Value[i].Clear()
	}
}

// GSWCiphertext represents an encrypted GSW ciphertext,
// which is a DefaultLWEDimension+1 collection of Lev ciphertexts.
type GSWCiphertext struct {
	GadgetParameters GadgetParameters

	// Value has length DefaultLWEDimension + 1.
	Value []LevCiphertext
}

// NewGSWCiphertext creates a new GSW ciphertext.
func NewGSWCiphertext(params Parameters, gadgetParams GadgetParameters) GSWCiphertext {
	lweDimension := params.lweDimension
	ct := make([]LevCiphertext, lweDimension+1)
	for i := 0; i < lweDimension+1; i++ {
		ct[i] = NewLevCiphertext(params, gadgetParams)
	}
	return GSWCiphertext{Value: ct, GadgetParameters: gadgetParams}
}

// NewGSWCiphertextCustom creates a new GSW ciphertext with given dimension.
func NewGSWCiphertextCustom(lweDimension int, gadgetParams GadgetParameters) GSWCiphertext {
	ct := make([]LevCiphertext, lweDimension+1)
	for i := 0; i < lweDimension+1; i++ {
		ct[i] = NewLevCiphertextCustom(lweDimension, gadgetParams)
	}
	return GSWCiphertext{Value: ct, GadgetParameters: gadgetParams}
}

// Copy returns a copy of the ciphertext.
func (ct GSWCiphertext) Copy() GSWCiphertext {
	ctCopy := make([]LevCiphertext, len(ct.Value))
	for i := range ct.Value {
		ctCopy[i] = ct.Value[i].Copy()
	}
	return GSWCiphertext{Value: ctCopy, GadgetParameters: ct.GadgetParameters}
}

// CopyFrom copies values from the ciphertext.
func (ct *GSWCiphertext) CopyFrom(ctIn GSWCiphertext) {
	for i := range ct.Value {
		ct.Value[i].CopyFrom(ctIn.Value[i])
	}
	ct.GadgetParameters = ctIn.GadgetParameters
}

// Clear clears the ciphertext.
func (ct *GSWCiphertext) Clear() {
	for i := range ct.Value {
		ct.Value[i].Clear()
	}
}
