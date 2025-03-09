package carousel

import (
	"github.com/sp301415/carousel/math/num"
	"github.com/sp301415/carousel/math/poly"
)

// Encrypt encodes and encrypts integer message to RLWE ciphertext.
func (e *Encryptor) Encrypt(message int) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.EncryptAssign(message, ctOut)
	return ctOut
}

// EncryptRLWEAssign encodes and encrypts integer messages to RLWE ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptAssign(message int, ctOut RLWECiphertext) {
	e.EncodeRLWEAssign([]int{message}, e.buffer.ptRLWE)
	e.EncryptRLWEPlaintextAssign(e.buffer.ptRLWE, ctOut)
}

// EncryptRLWE encodes and encrypts integer messages to RLWE ciphertext.
func (e *Encryptor) EncryptRLWE(messages []int) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.EncryptRLWEAssign(messages, ctOut)
	return ctOut
}

// EncryptRLWEAssign encodes and encrypts integer messages to RLWE ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptRLWEAssign(messages []int, ctOut RLWECiphertext) {
	e.EncodeRLWEAssign(messages, e.buffer.ptRLWE)
	e.EncryptRLWEPlaintextAssign(e.buffer.ptRLWE, ctOut)
}

// EncryptRLWECoeffs encodes encrypts integer messages to RLWE ciphertext using coefficient encoding.
func (e *Encryptor) EncryptRLWECoeffs(messages []int) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.EncryptRLWECoeffsAssign(messages, ctOut)
	return ctOut
}

// EncryptRLWECoeffsAssign encodes and encrypts integer messages to RLWE ciphertext using coefficient encoding
// and writes it to ctOut.
func (e *Encryptor) EncryptRLWECoeffsAssign(messages []int, ctOut RLWECiphertext) {
	e.EncodeRLWECoeffsAssign(messages, e.buffer.ptRLWE)
	e.EncryptRLWEPlaintextAssign(e.buffer.ptRLWE, ctOut)
}

// EncryptRLWESlots encodes and encrypts integer messages to RLWE ciphertext using slot encoding.
func (e *Encryptor) EncryptRLWESlots(messages []int) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.EncryptRLWESlotsAssign(messages, ctOut)
	return ctOut
}

// EncryptRLWESlotsAssign encodes and encrypts integer messages to RLWE ciphertext using slot encoding
// and writes it to ctOut.
func (e *Encryptor) EncryptRLWESlotsAssign(messages []int, ctOut RLWECiphertext) {
	e.EncodeRLWESlotsAssign(messages, e.buffer.ptRLWE)
	e.EncryptRLWEPlaintextAssign(e.buffer.ptRLWE, ctOut)
}

// EncryptRLWEPlaintext encrypts RLWE plaintext to RLWE ciphertext.
func (e *Encryptor) EncryptRLWEPlaintext(pt RLWEPlaintext) RLWECiphertext {
	ctOut := NewRLWECiphertext(e.Parameters)
	e.EncryptRLWEPlaintextAssign(pt, ctOut)
	return ctOut
}

// EncryptRLWEPlaintextAssign encrypts RLWE plaintext to RLWE ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptRLWEPlaintextAssign(pt RLWEPlaintext, ctOut RLWECiphertext) {
	ctOut.Value[0].CopyFrom(pt.Value)
	e.EncryptRLWEBody(ctOut)
}

// EncryptRLWEBody encrypts the value in the body of RLWE ciphertext and overrides it.
// This avoids the need for most buffers.
func (e *Encryptor) EncryptRLWEBody(ct RLWECiphertext) {
	e.UniformSampler.SamplePolyAssign(ct.Value[1])
	e.PolyEvaluator.ShortFourierPolyMulSubPolyAssign(ct.Value[1], e.SecretKey.FourierRLWEKey.Value, ct.Value[0])
	e.GaussianSampler.SamplePolyAddAssign(e.Parameters.RLWEStdDev(), ct.Value[0])
}

// Decrypt decodes and decrypts RLWE ciphertext to integer message.
func (e *Encryptor) Decrypt(ct RLWECiphertext) int {
	e.DecryptRLWEPhaseAssign(ct, e.buffer.ptRLWE)
	return e.DecodeRLWE(e.buffer.ptRLWE)[0]
}

// DecryptRLWE decrypts and decodes RLWE ciphertext to integer message.
func (e *Encryptor) DecryptRLWE(ct RLWECiphertext) []int {
	e.DecryptRLWEPhaseAssign(ct, e.buffer.ptRLWE)
	return e.DecodeRLWE(e.buffer.ptRLWE)
}

// DecryptRLWEAssign decrypts and decodes RLWE ciphertext to integer message and writes it to messagesOut.
func (e *Encryptor) DecryptRLWEAssign(ct RLWECiphertext, messagesOut []int) {
	e.DecryptRLWEPhaseAssign(ct, e.buffer.ptRLWE)
	e.DecodeRLWEAssign(e.buffer.ptRLWE, messagesOut)
}

// DecryptRLWECoeffs decrypts and decodes RLWE ciphertext to integer messages using coefficient decoding.
func (e *Encryptor) DecryptRLWECoeffs(ct RLWECiphertext) []int {
	e.DecryptRLWEPhaseAssign(ct, e.buffer.ptRLWE)
	return e.DecodeRLWECoeffs(e.buffer.ptRLWE)
}

// DecryptRLWECoeffsAssign decrypts and decodes RLWE ciphertext to integer messages using coefficient decoding
// and writes it to messagesOut.
func (e *Encryptor) DecryptRLWECoeffsAssign(ct RLWECiphertext, messagesOut []int) {
	e.DecryptRLWEPhaseAssign(ct, e.buffer.ptRLWE)
	e.DecodeRLWECoeffsAssign(e.buffer.ptRLWE, messagesOut)
}

// DecryptRLWESlots decrypts and decodes RLWE ciphertext to integer messages using slot decoding.
func (e *Encryptor) DecryptRLWESlots(ct RLWECiphertext) []int {
	e.DecryptRLWEPhaseAssign(ct, e.buffer.ptRLWE)
	return e.DecodeRLWESlots(e.buffer.ptRLWE)
}

// DecryptRLWESlotsAssign decrypts and decodes RLWE ciphertext to integer messages using slot decoding
// and writes it to messagesOut.
func (e *Encryptor) DecryptRLWESlotsAssign(ct RLWECiphertext, messagesOut []int) {
	e.DecryptRLWEPhaseAssign(ct, e.buffer.ptRLWE)
	e.DecodeRLWESlotsAssign(e.buffer.ptRLWE, messagesOut)
}

// DecryptRLWEPhase decrypts RLWE ciphertext to RLWE plaintext with errors.
func (e *Encryptor) DecryptRLWEPhase(ct RLWECiphertext) RLWEPlaintext {
	pt := NewRLWEPlaintext(e.Parameters)
	e.DecryptRLWEPhaseAssign(ct, pt)
	return pt
}

// DecryptRLWEPhaseAssign decrypts RLWE ciphertext to RLWE plaintext with errors and writes it to ptOut.
func (e *Encryptor) DecryptRLWEPhaseAssign(ct RLWECiphertext, ptOut RLWEPlaintext) {
	ptOut.Value.CopyFrom(ct.Value[0])
	e.PolyEvaluator.ShortFourierPolyMulAddPolyAssign(ct.Value[1], e.SecretKey.FourierRLWEKey.Value, ptOut.Value)
}

// EncryptRLevPoly encrypts polynomial to RLev ciphertext.
func (e *Encryptor) EncryptRLevPoly(p poly.Poly, gadgetParams GadgetParameters) RLevCiphertext {
	ctOut := NewRLevCiphertext(e.Parameters, gadgetParams)
	e.EncryptRLevPolyAssign(p, ctOut)
	return ctOut
}

// EncryptRLevPolyAssign encrypts polynomial to RLev ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptRLevPolyAssign(p poly.Poly, ctOut RLevCiphertext) {
	for i := 0; i < ctOut.GadgetParameters.level; i++ {
		e.PolyEvaluator.ScalarMulPolyAssign(p, ctOut.GadgetParameters.BaseQ(i), ctOut.Value[i].Value[0])
		e.EncryptRLWEBody(ctOut.Value[i])
	}
}

// DecryptRLevPoly decrypts RLev ciphertext to polynomial.
func (e *Encryptor) DecryptRLevPoly(ct RLevCiphertext) poly.Poly {
	pOut := poly.NewPoly(e.Parameters.polyDegree)
	e.DecryptRLevPolyAssign(ct, pOut)
	return pOut
}

// DecryptRLevPolyAssign decrypts RLev ciphertext to polynomial and writes it to pOut.
func (e *Encryptor) DecryptRLevPolyAssign(ct RLevCiphertext, pOut poly.Poly) {
	e.DecryptRLWEPhaseAssign(ct.Value[0], RLWEPlaintext{Value: pOut})
	for i := 0; i < e.Parameters.polyDegree; i++ {
		pOut.Coeffs[i] = num.DivRoundBits(pOut.Coeffs[i], ct.GadgetParameters.LogFirstBaseQ()) % ct.GadgetParameters.base
	}
}

// EncryptRGSWPoly encrypts polynomial to RGSW ciphertext.
func (e *Encryptor) EncryptRGSWPoly(p poly.Poly, gadgetParams GadgetParameters) RGSWCiphertext {
	ctOut := NewRGSWCiphertext(e.Parameters, gadgetParams)
	e.EncryptRGSWPolyAssign(p, ctOut)
	return ctOut
}

// EncryptRGSWPolyAssign encrypts polynomial to RGSW ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptRGSWPolyAssign(p poly.Poly, ctOut RGSWCiphertext) {
	e.EncryptRLevPolyAssign(p, ctOut.Value[0])
	e.PolyEvaluator.ShortFourierPolyMulPolyAssign(p, e.SecretKey.FourierRLWEKey.Value, e.buffer.ptRGSW)
	for i := 0; i < ctOut.GadgetParameters.level; i++ {
		e.PolyEvaluator.ScalarMulPolyAssign(e.buffer.ptRGSW, ctOut.GadgetParameters.BaseQ(i), ctOut.Value[1].Value[i].Value[0])
		e.EncryptRLWEBody(ctOut.Value[1].Value[i])
	}
}

// DecryptRGSWPoly decrypts RGSW ciphertext to polynomial.
func (e *Encryptor) DecryptRGSWPoly(ct RGSWCiphertext) poly.Poly {
	return e.DecryptRLevPoly(ct.Value[0])
}

// DecryptRGSWPolyAssign decrypts RGSW ciphertext to polynomial and writes it to ptOut.
func (e *Encryptor) DecryptRGSWPolyAssign(ct RGSWCiphertext, pOut poly.Poly) {
	e.DecryptRLevPolyAssign(ct.Value[0], pOut)
}
