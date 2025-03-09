package carousel

import (
	"github.com/sp301415/carousel/math/num"
	"github.com/sp301415/carousel/math/vec"
)

// EncryptLWE encodes and encrypts integer message to LWE ciphertext.
func (e *Encryptor) EncryptLWE(message int) LWECiphertext {
	return e.EncryptLWEPlaintext(e.EncodeLWE(message))
}

// EncryptLWEAssign encodes and encrypts integer message to LWE ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptLWEAssign(message int, ctOut LWECiphertext) {
	e.EncryptLWEPlaintextAssign(e.EncodeLWE(message), ctOut)
}

// EncryptLWEPlaintext encrypts LWE plaintext to LWE ciphertext.
func (e *Encryptor) EncryptLWEPlaintext(pt LWEPlaintext) LWECiphertext {
	ctOut := NewLWECiphertext(e.Parameters)
	e.EncryptLWEPlaintextAssign(pt, ctOut)
	return ctOut
}

// EncryptLWEPlaintextAssign encrypts LWE plaintext to LWE ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptLWEPlaintextAssign(pt LWEPlaintext, ctOut LWECiphertext) {
	ctOut.Value[0] = pt.Value
	e.EncryptLWEBody(ctOut)
}

// EncryptLWEBody encrypts the value in the body of LWE ciphertext and overrides it.
// This avoids the need for most buffers.
func (e *Encryptor) EncryptLWEBody(ct LWECiphertext) {
	e.UniformSampler.SampleVecAssign(ct.Value[1:])
	ct.Value[0] += -vec.Dot(ct.Value[1:], e.SecretKey.LWEKey.Value)
	ct.Value[0] += e.GaussianSampler.Sample(e.Parameters.lweStdDev)
}

// DecryptLWE decrypts and decodes LWE ciphertext to integer message.
func (e *Encryptor) DecryptLWE(ct LWECiphertext) int {
	return e.DecodeLWE(e.DecryptLWEPhase(ct))
}

// DecryptLWEPhase decrypts LWE ciphertext to LWE plaintext with errors.
func (e *Encryptor) DecryptLWEPhase(ct LWECiphertext) LWEPlaintext {
	pt := ct.Value[0] + vec.Dot(ct.Value[1:], e.SecretKey.LWEKey.Value)
	return LWEPlaintext{Value: pt}
}

// EncryptLev encrypts integer message to Lev ciphertext.
func (e *Encryptor) EncryptLev(message int, gadgetParams GadgetParameters) LevCiphertext {
	return e.EncryptLevScalar(uint64(message)%e.Parameters.messageModulus, gadgetParams)
}

// EncryptLevAssign encrypts integer message to Lev ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptLevAssign(message int, ctOut LevCiphertext) {
	e.EncryptLevScalarAssign(uint64(message)%e.Parameters.messageModulus, ctOut)
}

// EncryptLevScalar encrypts scalar to Lev ciphertext.
func (e *Encryptor) EncryptLevScalar(c uint64, gadgetParams GadgetParameters) LevCiphertext {
	ct := NewLevCiphertext(e.Parameters, gadgetParams)
	e.EncryptLevScalarAssign(c, ct)
	return ct
}

// EncryptLevScalarAssign encrypts scalar to Lev ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptLevScalarAssign(c uint64, ctOut LevCiphertext) {
	for i := 0; i < ctOut.GadgetParameters.level; i++ {
		ctOut.Value[i].Value[0] = c << ctOut.GadgetParameters.LogBaseQ(i)
		e.EncryptLWEBody(ctOut.Value[i])
	}
}

// DecryptLev decrypts Lev ciphertext to integer message.
func (e *Encryptor) DecryptLev(ct LevCiphertext) int {
	return int(e.DecryptLevScalar(ct) % e.Parameters.messageModulus)
}

// DecryptLevScalar decrypts Lev ciphertext to scalar.
func (e *Encryptor) DecryptLevScalar(ct LevCiphertext) uint64 {
	pt := e.DecryptLWEPhase(ct.Value[0])
	return num.DivRoundBits(pt.Value, ct.GadgetParameters.LogFirstBaseQ()) % ct.GadgetParameters.base
}

// EncryptGSW encrypts integer message to GSW ciphertext.
func (e *Encryptor) EncryptGSW(message int, gadgetParams GadgetParameters) GSWCiphertext {
	return e.EncryptGSWScalar(uint64(message)%e.Parameters.messageModulus, gadgetParams)
}

// EncryptGSWAssign encrypts integer message to GSW ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptGSWAssign(message int, ctOut GSWCiphertext) {
	e.EncryptGSWScalarAssign(uint64(message)%e.Parameters.messageModulus, ctOut)
}

// EncryptGSWScalar encrypts scalar to GSW ciphertext.
func (e *Encryptor) EncryptGSWScalar(c uint64, gadgetParams GadgetParameters) GSWCiphertext {
	ct := NewGSWCiphertext(e.Parameters, gadgetParams)
	e.EncryptGSWScalarAssign(c, ct)
	return ct
}

// EncryptGSWScalarAssign encrypts scalar to GSW ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptGSWScalarAssign(c uint64, ctOut GSWCiphertext) {
	e.EncryptLevScalarAssign(c, ctOut.Value[0])

	for i := 0; i < e.Parameters.lweDimension; i++ {
		for j := 0; j < ctOut.GadgetParameters.level; j++ {
			ctOut.Value[i+1].Value[j].Value[0] = e.SecretKey.LWEKey.Value[i] * c << ctOut.GadgetParameters.LogBaseQ(j)
			e.EncryptLWEBody(ctOut.Value[i+1].Value[j])
		}
	}
}

// DecryptGSW decrypts GSW ciphertext to integer message.
func (e *Encryptor) DecryptGSW(ct GSWCiphertext) int {
	return e.DecryptLev(ct.Value[0])
}

// DecryptGSWScalar decrypts GSW ciphertext to scalar.
func (e *Encryptor) DecryptGSWScalar(ct GSWCiphertext) uint64 {
	return e.DecryptLevScalar(ct.Value[0])
}
