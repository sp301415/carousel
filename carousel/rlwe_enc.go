package carousel

import "github.com/sp301415/carousel/math/num"

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
	ct := NewRLWECiphertext(e.Parameters)
	e.EncryptRLWEPlaintextAssign(pt, ct)
	return ct
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
	e.DecryptRLWEPlaintextAssign(ct, e.buffer.ptRLWE)
	return e.DecodeRLWE(e.buffer.ptRLWE)[0]
}

// DecryptRLWE decrypts and decodes RLWE ciphertext to integer message.
func (e *Encryptor) DecryptRLWE(ct RLWECiphertext) []int {
	e.DecryptRLWEPlaintextAssign(ct, e.buffer.ptRLWE)
	return e.DecodeRLWE(e.buffer.ptRLWE)
}

// DecryptRLWEAssign decrypts and decodes RLWE ciphertext to integer message and writes it to messagesOut.
func (e *Encryptor) DecryptRLWEAssign(ct RLWECiphertext, messagesOut []int) {
	e.DecryptRLWEPlaintextAssign(ct, e.buffer.ptRLWE)
	e.DecodeRLWEAssign(e.buffer.ptRLWE, messagesOut)
}

// DecryptRLWECoeffs decrypts and decodes RLWE ciphertext to integer messages using coefficient decoding.
func (e *Encryptor) DecryptRLWECoeffs(ct RLWECiphertext) []int {
	e.DecryptRLWEPlaintextAssign(ct, e.buffer.ptRLWE)
	return e.DecodeRLWECoeffs(e.buffer.ptRLWE)
}

// DecryptRLWECoeffsAssign decrypts and decodes RLWE ciphertext to integer messages using coefficient decoding
// and writes it to messagesOut.
func (e *Encryptor) DecryptRLWECoeffsAssign(ct RLWECiphertext, messagesOut []int) {
	e.DecryptRLWEPlaintextAssign(ct, e.buffer.ptRLWE)
	e.DecodeRLWECoeffsAssign(e.buffer.ptRLWE, messagesOut)
}

// DecryptRLWESlots decrypts and decodes RLWE ciphertext to integer messages using slot decoding.
func (e *Encryptor) DecryptRLWESlots(ct RLWECiphertext) []int {
	e.DecryptRLWEPlaintextAssign(ct, e.buffer.ptRLWE)
	return e.DecodeRLWESlots(e.buffer.ptRLWE)
}

// DecryptRLWESlotsAssign decrypts and decodes RLWE ciphertext to integer messages using slot decoding
// and writes it to messagesOut.
func (e *Encryptor) DecryptRLWESlotsAssign(ct RLWECiphertext, messagesOut []int) {
	e.DecryptRLWEPlaintextAssign(ct, e.buffer.ptRLWE)
	e.DecodeRLWESlotsAssign(e.buffer.ptRLWE, messagesOut)
}

// DecryptRLWEPlaintext decrypts RLWE ciphertext to RLWE plaintext.
func (e *Encryptor) DecryptRLWEPlaintext(ct RLWECiphertext) RLWEPlaintext {
	pt := NewRLWEPlaintext(e.Parameters)
	e.DecryptRLWEPlaintextAssign(ct, pt)
	return pt
}

// DecryptRLWEPlaintextAssign decrypts RLWE ciphertext to RLWE plaintext and writes it to ptOut.
func (e *Encryptor) DecryptRLWEPlaintextAssign(ct RLWECiphertext, ptOut RLWEPlaintext) {
	ptOut.Value.CopyFrom(ct.Value[0])
	e.PolyEvaluator.ShortFourierPolyMulAddPolyAssign(ct.Value[1], e.SecretKey.FourierRLWEKey.Value, ptOut.Value)
}

// EncryptRLevPlaintext encrypts RLWE plaintext to RLev ciphertext.
func (e *Encryptor) EncryptRLevPlaintext(pt RLWEPlaintext, gadgetParams GadgetParameters) RLevCiphertext {
	ct := NewRLevCiphertext(e.Parameters, gadgetParams)
	e.EncryptRLevPlaintextAssign(pt, ct)
	return ct
}

// EncryptRLevPlaintextAssign encrypts RLWE plaintext to RLev ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptRLevPlaintextAssign(pt RLWEPlaintext, ctOut RLevCiphertext) {
	for i := 0; i < ctOut.GadgetParameters.level; i++ {
		e.PolyEvaluator.ScalarMulPolyAssign(pt.Value, ctOut.GadgetParameters.BaseQ(i), ctOut.Value[i].Value[0])
		e.EncryptRLWEBody(ctOut.Value[i])
	}
}

// DecryptRLevPlaintext decrypts RLev ciphertext to RLWE plaintext.
func (e *Encryptor) DecryptRLevPlaintext(ct RLevCiphertext) RLWEPlaintext {
	pt := NewRLWEPlaintext(e.Parameters)
	e.DecryptRLevPlaintextAssign(ct, pt)
	return pt
}

// DecryptRLevPlaintextAssign decrypts RLev ciphertext to RLWE plaintext and writes it to ptOut.
func (e *Encryptor) DecryptRLevPlaintextAssign(ct RLevCiphertext, ptOut RLWEPlaintext) {
	e.DecryptRLWEPlaintextAssign(ct.Value[0], ptOut)
	for i := 0; i < e.Parameters.polyDegree; i++ {
		ptOut.Value.Coeffs[i] = num.DivRoundBits(ptOut.Value.Coeffs[i], ct.GadgetParameters.LogFirstBaseQ())
	}
}

// EncryptRGSWPlaintext encrypts RLWE plaintext to RGSW ciphertext.
func (e *Encryptor) EncryptRGSWPlaintext(pt RLWEPlaintext, gadgetParams GadgetParameters) RGSWCiphertext {
	ct := NewRGSWCiphertext(e.Parameters, gadgetParams)
	e.EncryptRGSWPlaintextAssign(pt, ct)
	return ct
}

// EncryptRGSWPlaintextAssign encrypts RLWE plaintext to RGSW ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptRGSWPlaintextAssign(pt RLWEPlaintext, ctOut RGSWCiphertext) {
	e.EncryptRLevPlaintextAssign(pt, ctOut.Value[0])
	e.PolyEvaluator.ShortFourierPolyMulPolyAssign(pt.Value, e.SecretKey.FourierRLWEKey.Value, e.buffer.ptRGSW)
	for i := 0; i < ctOut.GadgetParameters.level; i++ {
		e.PolyEvaluator.ScalarMulPolyAssign(e.buffer.ptRGSW, ctOut.GadgetParameters.BaseQ(i), ctOut.Value[1].Value[i].Value[0])
		e.EncryptRLWEBody(ctOut.Value[1].Value[i])
	}
}

// DecryptRGSWPlaintext decrypts RGSW ciphertext to RLWE plaintext.
func (e *Encryptor) DecryptRGSWPlaintext(ct RGSWCiphertext) RLWEPlaintext {
	return e.DecryptRLevPlaintext(ct.Value[0])
}

// DecryptRGSWPlaintextAssign decrypts RGSW ciphertext to RLWE plaintext and writes it to ptOut.
func (e *Encryptor) DecryptRGSWPlaintextAssign(ct RGSWCiphertext, ptOut RLWEPlaintext) {
	e.DecryptRLevPlaintextAssign(ct.Value[0], ptOut)
}
