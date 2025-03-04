package carousel

import (
	"github.com/sp301415/carousel/math/csprng"
	"github.com/sp301415/carousel/math/poly"
)

// Encryptor encrypts and decrypts TFHE plaintexts and ciphertexts.
// This is meant to be private, only for clients.
//
// Encryptor is not safe for concurrent use.
// Use [*Encryptor.ShallowCopy] to get a safe copy.
type Encryptor struct {
	// Encoder is an embedded encoder for this Encryptor.
	*Encoder
	// RLWETransformer is an embedded RLWETransformer for this Encryptor.
	*RLWETransformer

	// Parameters is the parameters for this Encryptor.
	Parameters Parameters

	// UniformSampler is used for sampling the mask of encryptions.
	UniformSampler *csprng.UniformSampler
	// BinarySampler is used for sampling LWE and RLWE key.
	BinarySampler *csprng.BinarySampler
	// GaussainSampler is used for sampling noise in LWE and RLWE encryption.
	GaussianSampler *csprng.GaussianSampler

	// PolyEvaluator is a PolyEvaluator for this Encryptor.
	PolyEvaluator *poly.Evaluator

	// SecretKey is the LWE and RLWE key for this Encryptor.
	SecretKey SecretKey

	buffer encryptionBuffer
}

// encryptionBuffer is a buffer for Encryptor.
type encryptionBuffer struct {
	// ptRLWE is the RLWE plaintext for RLWE encryption / decryptions.
	ptRLWE RLWEPlaintext
	// ctRLWE is the standard RLWE Ciphertext for Fourier encryption / decryptions.
	ctRLWE RLWECiphertext
	// ptRGSW is RLWEKey * Pt in RGSW encryption.
	ptRGSW poly.Poly
}

// NewEncryptor returns a initialized Encryptor with given parameters.
// It also automatically samples LWE and RLWE key.
func NewEncryptor(params Parameters) *Encryptor {
	// Fill samplers to call encryptor.GenSecretKey()
	encryptor := Encryptor{
		Encoder:         NewEncoder(params),
		RLWETransformer: NewRLWETransformer(params.polyEvaluatorParameters),

		Parameters: params,

		UniformSampler:  csprng.NewUniformSampler(),
		BinarySampler:   csprng.NewBinarySampler(),
		GaussianSampler: csprng.NewGaussianSampler(),

		PolyEvaluator: poly.NewEvaluator(params.polyEvaluatorParameters),

		buffer: newEncryptionBuffer(params),
	}

	encryptor.SecretKey = encryptor.GenSecretKey()

	return &encryptor
}

// NewEncryptorWithKey returns a initialized Encryptor with given parameters and key.
// This does not copy secret keys.
func NewEncryptorWithKey(params Parameters, sk SecretKey) *Encryptor {
	return &Encryptor{
		Encoder:         NewEncoder(params),
		RLWETransformer: NewRLWETransformer(params.polyEvaluatorParameters),

		Parameters: params,

		UniformSampler:  csprng.NewUniformSampler(),
		BinarySampler:   csprng.NewBinarySampler(),
		GaussianSampler: csprng.NewGaussianSampler(),

		PolyEvaluator: poly.NewEvaluator(params.polyEvaluatorParameters),

		SecretKey: sk,

		buffer: newEncryptionBuffer(params),
	}
}

// newEncryptionBuffer creates a new encryptionBuffer.
func newEncryptionBuffer(params Parameters) encryptionBuffer {
	return encryptionBuffer{
		ptRLWE: NewRLWEPlaintext(params),
		ctRLWE: NewRLWECiphertext(params),
		ptRGSW: poly.NewPoly(params.polyDegree),
	}
}

// ShallowCopy returns a shallow copy of this Encryptor.
// Returned Encryptor is safe for concurrent use.
func (e *Encryptor) ShallowCopy() *Encryptor {
	return &Encryptor{
		Encoder:         e.Encoder,
		RLWETransformer: e.RLWETransformer.ShallowCopy(),

		Parameters: e.Parameters,

		UniformSampler:  csprng.NewUniformSampler(),
		BinarySampler:   csprng.NewBinarySampler(),
		GaussianSampler: csprng.NewGaussianSampler(),

		SecretKey: e.SecretKey,

		PolyEvaluator: e.PolyEvaluator.ShallowCopy(),

		buffer: newEncryptionBuffer(e.Parameters),
	}
}

// GenSecretKey samples a new SecretKey.
// The SecretKey of the Encryptor is not changed.
func (e *Encryptor) GenSecretKey() SecretKey {
	sk := NewSecretKey(e.Parameters)

	if e.Parameters.blockSize == 1 {
		e.BinarySampler.SampleVecAssign(sk.LWELargeKey.Value)
	} else {
		e.BinarySampler.SampleBlockVecAssign(e.Parameters.blockSize, sk.LWELargeKey.Value[:e.Parameters.lweDimension])
		e.BinarySampler.SampleVecAssign(sk.LWELargeKey.Value[e.Parameters.lweDimension:])
	}

	e.ToFourierRLWESecretKeyAssign(sk.RLWEKey, sk.FourierRLWEKey)

	return sk
}
