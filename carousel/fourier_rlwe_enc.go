package carousel

import (
	"github.com/sp301415/carousel/math/num"
	"github.com/sp301415/carousel/math/poly"
)

// EncryptFourier encodes and encrypts integer message to FourierRLWE ciphertext.
func (e *Encryptor) EncryptFourier(message int) FourierRLWECiphertext {
	ctOut := NewFourierRLWECiphertext(e.Parameters)
	e.EncryptFourierAssign(message, ctOut)
	return ctOut
}

// EncryptRLWEFourierAssign encodes and encrypts integer messages to FourierRLWE ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptFourierAssign(message int, ctOut FourierRLWECiphertext) {
	e.EncodeRLWEAssign([]int{message}, e.buffer.ptRLWE)
	e.EncryptFourierRLWEPlaintextAssign(e.buffer.ptRLWE, ctOut)
}

// EncryptFourierRLWE encodes and encrypts integer messages to FourierRLWE ciphertext.
func (e *Encryptor) EncryptFourierRLWE(messages []int) FourierRLWECiphertext {
	ctOut := NewFourierRLWECiphertext(e.Parameters)
	e.EncryptFourierRLWEAssign(messages, ctOut)
	return ctOut
}

// EncryptFourierRLWEAssign encodes and encrypts integer messages to FourierRLWE ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptFourierRLWEAssign(messages []int, ctOut FourierRLWECiphertext) {
	e.EncodeRLWEAssign(messages, e.buffer.ptRLWE)
	e.EncryptFourierRLWEPlaintextAssign(e.buffer.ptRLWE, ctOut)
}

// EncryptFourierRLWECoeffs encodes encrypts integer messages to FourierRLWE ciphertext using coefficient encoding.
func (e *Encryptor) EncryptFourierRLWECoeffs(messages []int) FourierRLWECiphertext {
	ctOut := NewFourierRLWECiphertext(e.Parameters)
	e.EncryptFourierRLWECoeffsAssign(messages, ctOut)
	return ctOut
}

// EncryptFourierRLWECoeffsAssign encodes and encrypts integer messages to FourierRLWE ciphertext using coefficient encoding
// and writes it to ctOut.
func (e *Encryptor) EncryptFourierRLWECoeffsAssign(messages []int, ctOut FourierRLWECiphertext) {
	e.EncodeRLWECoeffsAssign(messages, e.buffer.ptRLWE)
	e.EncryptFourierRLWEPlaintextAssign(e.buffer.ptRLWE, ctOut)
}

// EncryptFourierRLWESlots encodes and encrypts integer messages to FourierRLWE ciphertext using slot encoding.
func (e *Encryptor) EncryptFourierRLWESlots(messages []int) FourierRLWECiphertext {
	ctOut := NewFourierRLWECiphertext(e.Parameters)
	e.EncryptFourierRLWESlotsAssign(messages, ctOut)
	return ctOut
}

// EncryptFourierRLWESlotsAssign encodes and encrypts integer messages to FourierRLWE ciphertext using slot encoding
// and writes it to ctOut.
func (e *Encryptor) EncryptFourierRLWESlotsAssign(messages []int, ctOut FourierRLWECiphertext) {
	e.EncodeRLWESlotsAssign(messages, e.buffer.ptRLWE)
	e.EncryptFourierRLWEPlaintextAssign(e.buffer.ptRLWE, ctOut)
}

// EncryptFourierRLWEPlaintext encrypts RLWE plaintext to FourierRLWE ciphertext.
func (e *Encryptor) EncryptFourierRLWEPlaintext(pt RLWEPlaintext) FourierRLWECiphertext {
	ctOut := NewFourierRLWECiphertext(e.Parameters)
	e.EncryptFourierRLWEPlaintextAssign(pt, ctOut)
	return ctOut
}

// EncryptFourierRLWEPlaintextAssign encrypts RLWE plaintext to FourierRLWE ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptFourierRLWEPlaintextAssign(pt RLWEPlaintext, ctOut FourierRLWECiphertext) {
	e.EncryptRLWEPlaintextAssign(pt, e.buffer.ctRLWE)
	e.ToFourierRLWECiphertextAssign(e.buffer.ctRLWE, ctOut)
}

// DecryptFourier decodes and decrypts FourierRLWE ciphertext to integer message.
func (e *Encryptor) DecryptFourier(ct FourierRLWECiphertext) int {
	e.DecryptFourierRLWEPhaseAssign(ct, e.buffer.ptRLWE)
	return e.DecodeRLWE(e.buffer.ptRLWE)[0]
}

// DecryptFourierRLWE decrypts and decodes FourierRLWE ciphertext to integer message.
func (e *Encryptor) DecryptFourierRLWE(ct FourierRLWECiphertext) []int {
	e.DecryptFourierRLWEPhaseAssign(ct, e.buffer.ptRLWE)
	return e.DecodeRLWE(e.buffer.ptRLWE)
}

// DecryptFourierRLWEAssign decrypts and decodes FourierRLWE ciphertext to integer message and writes it to messagesOut.
func (e *Encryptor) DecryptFourierRLWEAssign(ct FourierRLWECiphertext, messagesOut []int) {
	e.DecryptFourierRLWEPhaseAssign(ct, e.buffer.ptRLWE)
	e.DecodeRLWEAssign(e.buffer.ptRLWE, messagesOut)
}

// DecryptFourierRLWEPhase decrypts RLWE ciphertext to FourierRLWE plaintext with errors.
func (e *Encryptor) DecryptFourierRLWEPhase(ct FourierRLWECiphertext) RLWEPlaintext {
	ptOut := NewRLWEPlaintext(e.Parameters)
	e.DecryptFourierRLWEPhaseAssign(ct, ptOut)
	return ptOut
}

// DecryptFourierRLWEPhaseAssign decrypts RLWE ciphertext to FourierRLWE plaintext with errors and writes it to ptOut.
func (e *Encryptor) DecryptFourierRLWEPhaseAssign(ct FourierRLWECiphertext, ptOut RLWEPlaintext) {
	e.ToRLWECiphertextAssign(ct, e.buffer.ctRLWE)
	e.DecryptRLWEPhaseAssign(e.buffer.ctRLWE, ptOut)
}

// EncryptFourierRLevPoly encrypts polynomial to RLev ciphertext.
func (e *Encryptor) EncryptFourierRLevPoly(p poly.Poly, gadgetParams GadgetParameters) FourierRLevCiphertext {
	ctOut := NewFourierRLevCiphertext(e.Parameters, gadgetParams)
	e.EncryptFourierRLevPolyAssign(p, ctOut)
	return ctOut
}

// EncryptFourierRLevPolyAssign encrypts polynomial to FourierRLev ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptFourierRLevPolyAssign(p poly.Poly, ctOut FourierRLevCiphertext) {
	for i := 0; i < ctOut.GadgetParameters.level; i++ {
		e.PolyEvaluator.ScalarMulPolyAssign(p, ctOut.GadgetParameters.BaseQ(i), e.buffer.ctRLWE.Value[0])
		e.EncryptRLWEBody(e.buffer.ctRLWE)
		e.ToFourierRLWECiphertextAssign(e.buffer.ctRLWE, ctOut.Value[i])
	}
}

// DecryptFourierRLevPoly decrypts FourierRLev ciphertext to polynomial.
func (e *Encryptor) DecryptFourierRLevPoly(ct FourierRLevCiphertext) poly.Poly {
	pOut := poly.NewPoly(e.Parameters.polyDegree)
	e.DecryptFourierRLevPolyAssign(ct, pOut)
	return pOut
}

// DecryptFourierRLevPolyAssign decrypts FourierRLev ciphertext to polynomial and writes it to pOut.
func (e *Encryptor) DecryptFourierRLevPolyAssign(ct FourierRLevCiphertext, pOut poly.Poly) {
	e.DecryptFourierRLWEPhaseAssign(ct.Value[0], RLWEPlaintext{Value: pOut})
	for i := 0; i < e.Parameters.polyDegree; i++ {
		pOut.Coeffs[i] = num.DivRoundBits(pOut.Coeffs[i], ct.GadgetParameters.LogFirstBaseQ()) % ct.GadgetParameters.base
	}
}

// EncryptFourierRGSWPoly encrypts polynomial to FourierRGSW ciphertext.
func (e *Encryptor) EncryptFourierRGSWPoly(p poly.Poly, gadgetParams GadgetParameters) FourierRGSWCiphertext {
	ctOut := NewFourierRGSWCiphertext(e.Parameters, gadgetParams)
	e.EncryptFourierRGSWPolyAssign(p, ctOut)
	return ctOut
}

// EncryptFourierRGSWPolyAssign encrypts polynomial to FourierRGSW ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptFourierRGSWPolyAssign(p poly.Poly, ctOut FourierRGSWCiphertext) {
	e.EncryptFourierRLevPolyAssign(p, ctOut.Value[0])
	e.PolyEvaluator.ShortFourierPolyMulPolyAssign(p, e.SecretKey.FourierRLWEKey.Value, e.buffer.ptRGSW)
	for i := 0; i < ctOut.GadgetParameters.level; i++ {
		e.PolyEvaluator.ScalarMulPolyAssign(e.buffer.ptRGSW, ctOut.GadgetParameters.BaseQ(i), e.buffer.ctRLWE.Value[0])
		e.EncryptRLWEBody(e.buffer.ctRLWE)
		e.ToFourierRLWECiphertextAssign(e.buffer.ctRLWE, ctOut.Value[1].Value[i])
	}
}

// DecryptFourierRGSWPoly decrypts FourierRGSW ciphertext to polynomial.
func (e *Encryptor) DecryptFourierRGSWPoly(ct FourierRGSWCiphertext) poly.Poly {
	pOut := poly.NewPoly(e.Parameters.polyDegree)
	e.DecryptFourierRGSWPolyAssign(ct, pOut)
	return pOut
}

// DecryptFourierRGSWPolyAssign decrypts FourierRGSW ciphertext to polynomial and writes it to pOut.
func (e *Encryptor) DecryptFourierRGSWPolyAssign(ct FourierRGSWCiphertext, pOut poly.Poly) {
	e.DecryptFourierRLevPolyAssign(ct.Value[0], pOut)
}
