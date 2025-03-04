package carousel

import (
	"github.com/sp301415/carousel/math/vec"
)

// SecretKey is a structure containing LWE and RLWE key.
// All keys should be treated as read-only.
// Changing them mid-operation will usually result in wrong results.
//
// LWEKey and RLWEKey is sampled together, as explained in https://eprint.iacr.org/2023/958.
// As a result, LWEKey and RLWEKey share the same backing slice, so modifying one will affect the other.
type SecretKey struct {
	// LWELargeKey is a LWE key with length PolyDegree.
	// Essentially, this is same as RLWEKey but parsed differently.
	LWELargeKey LWESecretKey
	// RLWEKey is a key used for RLWE encryption and decryption.
	// Essentially, this is same as LWEKey but parsed differently.
	RLWEKey RLWESecretKey
	// FourierRLWEKey is a fourier transformed RLWEKey.
	// Used for RLWE encryption.
	FourierRLWEKey FourierRLWESecretKey
	// LWEKey is a LWE key with length LWEDimension.
	// Essentially, this is the first LWEDimension elements of LWEKey.
	LWEKey LWESecretKey
}

// NewSecretKey creates a new SecretKey.
// Each key shares the same backing slice, held by LWEKey.
func NewSecretKey(params Parameters) SecretKey {
	rlweKey := NewRLWESecretKey(params)
	fourierRLWEKey := NewFourierRLWESecretKey(params)
	lweLargeKey := LWESecretKey{Value: rlweKey.Value.Coeffs}

	lweKey := LWESecretKey{Value: lweLargeKey.Value[:params.lweDimension]}

	return SecretKey{
		LWELargeKey:    lweLargeKey,
		RLWEKey:        rlweKey,
		FourierRLWEKey: fourierRLWEKey,
		LWEKey:         lweKey,
	}
}

// NewSecretKeyCustom creates a new SecretKey with given dimension and polyDegree.
// Each key shares the same backing slice, held by LWEKey.
func NewSecretKeyCustom(lweDimension, polyDegree int) SecretKey {
	rlweKey := NewRLWESecretKeyCustom(polyDegree)
	fourierRLWEKey := NewFourierRLWESecretKeyCustom(polyDegree)
	lweLargeKey := LWESecretKey{Value: rlweKey.Value.Coeffs}

	lweKey := LWESecretKey{Value: lweLargeKey.Value[:lweDimension]}

	return SecretKey{
		LWELargeKey:    lweLargeKey,
		RLWEKey:        rlweKey,
		FourierRLWEKey: fourierRLWEKey,
		LWEKey:         lweKey,
	}
}

// Copy returns a copy of the key.
func (sk SecretKey) Copy() SecretKey {
	rlweKey := sk.RLWEKey.Copy()
	fourierRLWEKey := sk.FourierRLWEKey.Copy()
	lweLargeKey := LWESecretKey{Value: rlweKey.Value.Coeffs}

	lweKey := LWESecretKey{Value: lweLargeKey.Value[:len(sk.LWEKey.Value)]}

	return SecretKey{
		LWELargeKey:    lweLargeKey,
		RLWEKey:        rlweKey,
		FourierRLWEKey: fourierRLWEKey,
		LWEKey:         lweKey,
	}
}

// CopyFrom copies values from the key.
func (sk *SecretKey) CopyFrom(skIn SecretKey) {
	vec.CopyAssign(skIn.LWELargeKey.Value, sk.LWELargeKey.Value)
	sk.FourierRLWEKey.CopyFrom(skIn.FourierRLWEKey)
}

// Clear clears the key.
func (sk *SecretKey) Clear() {
	vec.Fill(sk.LWELargeKey.Value, 0)
	sk.FourierRLWEKey.Clear()
}
