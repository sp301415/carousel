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
	ct := NewLWECiphertext(e.Parameters)
	e.EncryptLWEPlaintextAssign(pt, ct)
	return ct
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
	return e.DecodeLWE(e.DecryptLWEPlaintext(ct))
}

// DecryptLWEPlaintext decrypts LWE ciphertext to LWE plaintext.
func (e *Encryptor) DecryptLWEPlaintext(ct LWECiphertext) LWEPlaintext {
	pt := ct.Value[0] + vec.Dot(ct.Value[1:], e.SecretKey.LWEKey.Value)
	return LWEPlaintext{Value: pt}
}

// EncryptLev encrypts integer message to Lev ciphertext.
func (e *Encryptor) EncryptLev(message int, gadgetParams GadgetParameters) LevCiphertext {
	pt := LWEPlaintext{Value: uint64(message) % e.Parameters.messageModulus}
	return e.EncryptLevPlaintext(pt, gadgetParams)
}

// EncryptLevAssign encrypts integer message to Lev ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptLevAssign(message int, ctOut LevCiphertext) {
	pt := LWEPlaintext{Value: uint64(message) % e.Parameters.messageModulus}
	e.EncryptLevPlaintextAssign(pt, ctOut)
}

// EncryptLevPlaintext encrypts LWE plaintext to Lev ciphertext.
func (e *Encryptor) EncryptLevPlaintext(pt LWEPlaintext, gadgetParams GadgetParameters) LevCiphertext {
	ct := NewLevCiphertext(e.Parameters, gadgetParams)
	e.EncryptLevPlaintextAssign(pt, ct)
	return ct
}

// EncryptLevPlaintextAssign encrypts LWE plaintext to Lev ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptLevPlaintextAssign(pt LWEPlaintext, ctOut LevCiphertext) {
	for i := 0; i < ctOut.GadgetParameters.level; i++ {
		ctOut.Value[i].Value[0] = pt.Value << ctOut.GadgetParameters.LogBaseQ(i)
		e.EncryptLWEBody(ctOut.Value[i])
	}
}

// DecryptLev decrypts Lev ciphertext to integer message.
func (e *Encryptor) DecryptLev(ct LevCiphertext) int {
	pt := e.DecryptLevPlaintext(ct)
	return int(pt.Value % e.Parameters.messageModulus)
}

// DecryptLevPlaintext decrypts Lev ciphertext to LWE plaintext.
func (e *Encryptor) DecryptLevPlaintext(ct LevCiphertext) LWEPlaintext {
	pt := e.DecryptLWEPlaintext(ct.Value[0])
	return LWEPlaintext{Value: num.DivRoundBits(pt.Value, ct.GadgetParameters.LogFirstBaseQ())}
}

// EncryptGSW encrypts integer message to GSW ciphertext.
func (e *Encryptor) EncryptGSW(message int, gadgetParams GadgetParameters) GSWCiphertext {
	pt := LWEPlaintext{Value: uint64(message) % e.Parameters.messageModulus}
	return e.EncryptGSWPlaintext(pt, gadgetParams)
}

// EncryptGSWAssign encrypts integer message to GSW ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptGSWAssign(message int, ctOut GSWCiphertext) {
	pt := LWEPlaintext{Value: uint64(message) % e.Parameters.messageModulus}
	e.EncryptGSWPlaintextAssign(pt, ctOut)
}

// EncryptGSWPlaintext encrypts LWE plaintext to GSW ciphertext.
func (e *Encryptor) EncryptGSWPlaintext(pt LWEPlaintext, gadgetParams GadgetParameters) GSWCiphertext {
	ct := NewGSWCiphertext(e.Parameters, gadgetParams)
	e.EncryptGSWPlaintextAssign(pt, ct)
	return ct
}

// EncryptGSWPlaintextAssign encrypts LWE plaintext to GSW ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptGSWPlaintextAssign(pt LWEPlaintext, ctOut GSWCiphertext) {
	e.EncryptLevPlaintextAssign(pt, ctOut.Value[0])

	for i := 0; i < e.Parameters.lweDimension; i++ {
		for j := 0; j < ctOut.GadgetParameters.level; j++ {
			ctOut.Value[i+1].Value[j].Value[0] = e.SecretKey.LWEKey.Value[i] * pt.Value << ctOut.GadgetParameters.LogBaseQ(j)
			e.EncryptLWEBody(ctOut.Value[i+1].Value[j])
		}
	}
}

// DecryptGSW decrypts GSW ciphertext to integer message.
func (e *Encryptor) DecryptGSW(ct GSWCiphertext) int {
	return e.DecryptLev(ct.Value[0])
}

// DecryptGSWPlaintext decrypts GSW ciphertext to LWE plaintext.
func (e *Encryptor) DecryptGSWPlaintext(ct GSWCiphertext) LWEPlaintext {
	return e.DecryptLevPlaintext(ct.Value[0])
}
