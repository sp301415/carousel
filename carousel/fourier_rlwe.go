package carousel

import "github.com/sp301415/carousel/math/poly"

// FourierRLWESecretKey is a RLWE key in Fourier domain.
type FourierRLWESecretKey struct {
	// Value is a FourierPoly.
	Value poly.FourierPoly
}

// NewFourierRLWESecretKey creates a new FourierRLWESecretKey.
func NewFourierRLWESecretKey(params Parameters) FourierRLWESecretKey {
	return FourierRLWESecretKey{Value: poly.NewFourierPoly(params.polyDegree)}
}

// NewFourierRLWESecretKeyCustom creates a new FourierRLWESecretKey with given dimension and polyDegree.
func NewFourierRLWESecretKeyCustom(polyDegree int) FourierRLWESecretKey {
	return FourierRLWESecretKey{Value: poly.NewFourierPoly(polyDegree)}
}

// Copy returns a copy of the key.
func (sk FourierRLWESecretKey) Copy() FourierRLWESecretKey {
	return FourierRLWESecretKey{Value: sk.Value.Copy()}
}

// CopyFrom copies values from the key.
func (sk *FourierRLWESecretKey) CopyFrom(skIn FourierRLWESecretKey) {
	sk.Value.CopyFrom(skIn.Value)
}

// Clear clears the key.
func (sk *FourierRLWESecretKey) Clear() {
	sk.Value.Clear()
}

// FourierRLWECiphertext is a RLWE ciphertext in Fourier domain.
type FourierRLWECiphertext struct {
	// Value is ordered as [body, mask],
	// since Go doesn't provide an easy way to take last element of slice.
	// Therefore, value has length 2.
	Value []poly.FourierPoly
}

// NewFourierRLWECiphertext creates a new FourierRLWECiphertext.
func NewFourierRLWECiphertext(params Parameters) FourierRLWECiphertext {
	ct := make([]poly.FourierPoly, 2)
	ct[0] = poly.NewFourierPoly(params.polyDegree)
	ct[1] = poly.NewFourierPoly(params.polyDegree)
	return FourierRLWECiphertext{Value: ct}
}

// NewFourierRLWECiphertextCustom creates a new FourierRLWECiphertext with given rank and polyDegree.
func NewFourierRLWECiphertextCustom(rank, polyDegree int) FourierRLWECiphertext {
	ct := make([]poly.FourierPoly, rank+1)
	for i := range ct {
		ct[i] = poly.NewFourierPoly(polyDegree)
	}
	return FourierRLWECiphertext{Value: ct}
}

// Copy returns a copy of the ciphertext.
func (ct FourierRLWECiphertext) Copy() FourierRLWECiphertext {
	ctCopy := make([]poly.FourierPoly, len(ct.Value))
	for i := range ctCopy {
		ctCopy[i] = ct.Value[i].Copy()
	}
	return FourierRLWECiphertext{Value: ctCopy}
}

// CopyFrom copies values from the ciphertext.
func (ct *FourierRLWECiphertext) CopyFrom(ctIn FourierRLWECiphertext) {
	for i := range ct.Value {
		ct.Value[i].CopyFrom(ctIn.Value[i])
	}
}

// Clear clears the ciphertext.
func (ct *FourierRLWECiphertext) Clear() {
	for i := range ct.Value {
		ct.Value[i].Clear()
	}
}

// FourierRLevCiphertext is a leveled RLWE ciphertext in Fourier domain.
type FourierRLevCiphertext struct {
	GadgetParameters GadgetParameters

	// Value has length Level.
	Value []FourierRLWECiphertext
}

// NewFourierRLevCiphertext creates a new FourierRLevCiphertext.
func NewFourierRLevCiphertext(params Parameters, gadgetParams GadgetParameters) FourierRLevCiphertext {
	ct := make([]FourierRLWECiphertext, gadgetParams.level)
	for i := 0; i < gadgetParams.level; i++ {
		ct[i] = NewFourierRLWECiphertext(params)
	}
	return FourierRLevCiphertext{Value: ct, GadgetParameters: gadgetParams}
}

// NewFourierRLevCiphertextCustom creates a new FourierRLevCiphertext with given polyDegree.
func NewFourierRLevCiphertextCustom(polyDegree int, gadgetParams GadgetParameters) FourierRLevCiphertext {
	ct := make([]FourierRLWECiphertext, gadgetParams.level)
	for i := 0; i < gadgetParams.level; i++ {
		ct[i] = NewFourierRLWECiphertextCustom(1, polyDegree)
	}
	return FourierRLevCiphertext{Value: ct, GadgetParameters: gadgetParams}
}

// Copy returns a copy of the ciphertext.
func (ct FourierRLevCiphertext) Copy() FourierRLevCiphertext {
	ctCopy := make([]FourierRLWECiphertext, len(ct.Value))
	for i := range ct.Value {
		ctCopy[i] = ct.Value[i].Copy()
	}
	return FourierRLevCiphertext{Value: ctCopy, GadgetParameters: ct.GadgetParameters}
}

// CopyFrom copies values from the ciphertext.
func (ct *FourierRLevCiphertext) CopyFrom(ctIn FourierRLevCiphertext) {
	for i := range ct.Value {
		ct.Value[i].CopyFrom(ctIn.Value[i])
	}
	ct.GadgetParameters = ctIn.GadgetParameters
}

// Clear clears the ciphertext.
func (ct *FourierRLevCiphertext) Clear() {
	for i := range ct.Value {
		ct.Value[i].Clear()
	}
}

// FourierRGSWCiphertext represents an encrypted RGSW ciphertext in Fourier domain.
type FourierRGSWCiphertext struct {
	GadgetParameters GadgetParameters

	// Value has length 2.
	Value []FourierRLevCiphertext
}

// NewFourierRGSWCiphertext creates a new RGSW ciphertext.
func NewFourierRGSWCiphertext(params Parameters, gadgetParams GadgetParameters) FourierRGSWCiphertext {
	ct := make([]FourierRLevCiphertext, 2)
	ct[0] = NewFourierRLevCiphertext(params, gadgetParams)
	ct[1] = NewFourierRLevCiphertext(params, gadgetParams)
	return FourierRGSWCiphertext{Value: ct, GadgetParameters: gadgetParams}
}

// NewFourierRGSWCiphertextCustom creates a new RGSW ciphertext with given  polyDegree.
func NewFourierRGSWCiphertextCustom(polyDegree int, gadgetParams GadgetParameters) FourierRGSWCiphertext {
	ct := make([]FourierRLevCiphertext, 2)
	ct[0] = NewFourierRLevCiphertextCustom(polyDegree, gadgetParams)
	ct[1] = NewFourierRLevCiphertextCustom(polyDegree, gadgetParams)
	return FourierRGSWCiphertext{Value: ct, GadgetParameters: gadgetParams}
}

// Copy returns a copy of the ciphertext.
func (ct FourierRGSWCiphertext) Copy() FourierRGSWCiphertext {
	ctCopy := make([]FourierRLevCiphertext, len(ct.Value))
	for i := range ct.Value {
		ctCopy[i] = ct.Value[i].Copy()
	}
	return FourierRGSWCiphertext{Value: ctCopy, GadgetParameters: ct.GadgetParameters}
}

// CopyFrom copies values from the ciphertext.
func (ct *FourierRGSWCiphertext) CopyFrom(ctIn FourierRGSWCiphertext) {
	for i := range ct.Value {
		ct.Value[i].CopyFrom(ctIn.Value[i])
	}
	ct.GadgetParameters = ctIn.GadgetParameters
}

// Clear clears the ciphertext.
func (ct *FourierRGSWCiphertext) Clear() {
	for i := range ct.Value {
		ct.Value[i].Clear()
	}
}
