package carousel

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
	ct := NewFourierRLWECiphertext(e.Parameters)
	e.EncryptFourierRLWEPlaintextAssign(pt, ct)
	return ct
}

// EncryptFourierRLWEPlaintextAssign encrypts RLWE plaintext to FourierRLWE ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptFourierRLWEPlaintextAssign(pt RLWEPlaintext, ctOut FourierRLWECiphertext) {
	e.EncryptRLWEPlaintextAssign(pt, e.buffer.ctRLWE)
	e.ToFourierRLWECiphertextAssign(e.buffer.ctRLWE, ctOut)
}

// DecryptFourier decodes and decrypts FourierRLWE ciphertext to integer message.
func (e *Encryptor) DecryptFourier(ct FourierRLWECiphertext) int {
	e.DecryptFourierRLWEPlaintextAssign(ct, e.buffer.ptRLWE)
	return e.DecodeRLWE(e.buffer.ptRLWE)[0]
}

// DecryptFourierRLWE decrypts and decodes FourierRLWE ciphertext to integer message.
func (e *Encryptor) DecryptFourierRLWE(ct FourierRLWECiphertext) []int {
	e.DecryptFourierRLWEPlaintextAssign(ct, e.buffer.ptRLWE)
	return e.DecodeRLWE(e.buffer.ptRLWE)
}

// DecryptFourierRLWEAssign decrypts and decodes FourierRLWE ciphertext to integer message and writes it to messagesOut.
func (e *Encryptor) DecryptFourierRLWEAssign(ct FourierRLWECiphertext, messagesOut []int) {
	e.DecryptFourierRLWEPlaintextAssign(ct, e.buffer.ptRLWE)
	e.DecodeRLWEAssign(e.buffer.ptRLWE, messagesOut)
}

// DecryptFourierRLWEPlaintext decrypts RLWE ciphertext to FourierRLWE plaintext.
func (e *Encryptor) DecryptFourierRLWEPlaintext(ct FourierRLWECiphertext) RLWEPlaintext {
	pt := NewRLWEPlaintext(e.Parameters)
	e.DecryptFourierRLWEPlaintextAssign(ct, pt)
	return pt
}

// DecryptFourierRLWEPlaintextAssign decrypts RLWE ciphertext to FourierRLWE plaintext and writes it to ptOut.
func (e *Encryptor) DecryptFourierRLWEPlaintextAssign(ct FourierRLWECiphertext, ptOut RLWEPlaintext) {
	e.ToRLWECiphertextAssign(ct, e.buffer.ctRLWE)
	e.DecryptRLWEPlaintextAssign(e.buffer.ctRLWE, ptOut)
}

// EncryptFourierRLevPlaintextAssign encrypts RLWE plaintext to FourierRLev ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptFourierRLevPlaintextAssign(pt RLWEPlaintext, ctOut FourierRLevCiphertext) {
	for i := 0; i < ctOut.GadgetParameters.level; i++ {
		e.PolyEvaluator.ScalarMulPolyAssign(pt.Value, ctOut.GadgetParameters.BaseQ(i), e.buffer.ctRLWE.Value[0])
		e.EncryptRLWEBody(e.buffer.ctRLWE)
		e.ToFourierRLWECiphertextAssign(e.buffer.ctRLWE, ctOut.Value[i])
	}
}

// DecryptFourierRLevPlaintext decrypts FourierRLev ciphertext to RLWE plaintext.
func (e *Encryptor) DecryptFourierRLevPlaintext(ct FourierRLevCiphertext) RLWEPlaintext {
	pt := NewRLWEPlaintext(e.Parameters)
	e.DecryptFourierRLevPlaintextAssign(ct, pt)
	return pt
}

// DecryptFourierRLevPlaintextAssign decrypts FourierRLev ciphertext to RLWE plaintext and writes it to ptOut.
func (e *Encryptor) DecryptFourierRLevPlaintextAssign(ct FourierRLevCiphertext, ptOut RLWEPlaintext) {
	ctLastLevel := ct.Value[ct.GadgetParameters.level-1]
	e.DecryptFourierRLWEPlaintextAssign(ctLastLevel, ptOut)
}

// EncryptFourierRGSWPlaintext encrypts RLWE plaintext to FourierRGSW ciphertext.
func (e *Encryptor) EncryptFourierRGSWPlaintext(pt RLWEPlaintext, gadgetParams GadgetParameters) FourierRGSWCiphertext {
	ct := NewFourierRGSWCiphertext(e.Parameters, gadgetParams)
	e.EncryptFourierRGSWPlaintextAssign(pt, ct)
	return ct
}

// EncryptFourierRGSWPlaintextAssign encrypts RLWE plaintext to FourierRGSW ciphertext and writes it to ctOut.
func (e *Encryptor) EncryptFourierRGSWPlaintextAssign(pt RLWEPlaintext, ctOut FourierRGSWCiphertext) {
	e.EncryptFourierRLevPlaintextAssign(pt, ctOut.Value[0])
	e.PolyEvaluator.ShortFourierPolyMulPolyAssign(pt.Value, e.SecretKey.FourierRLWEKey.Value, e.buffer.ptRGSW)
	for i := 0; i < ctOut.GadgetParameters.level; i++ {
		e.PolyEvaluator.ScalarMulPolyAssign(e.buffer.ptRGSW, ctOut.GadgetParameters.BaseQ(i), e.buffer.ctRLWE.Value[0])
		e.EncryptRLWEBody(e.buffer.ctRLWE)
		e.ToFourierRLWECiphertextAssign(e.buffer.ctRLWE, ctOut.Value[1].Value[i])
	}
}
