package carousel

// GenLWEKeySwitchKey samples a new keyswitch key skIn -> LWEKey.
func (e *Encryptor) GenLWEKeySwitchKey(skIn LWESecretKey, gadgetParams GadgetParameters) LWEKeySwitchKey {
	ksk := NewLWEKeySwitchKey(e.Parameters, len(skIn.Value), gadgetParams)

	for i := 0; i < ksk.InputLWEDimension(); i++ {
		e.EncryptLevScalarAssign(skIn.Value[i], ksk.Value[i])
	}

	return ksk
}

// GenGLWEKeySwitchKey samples a new keyswitch key skIn -> RLWEKey.
func (e *Encryptor) GenGLWEKeySwitchKey(skIn RLWESecretKey, gadgetParams GadgetParameters) RLWEKeySwitchKey {
	ksk := NewRLWEKeySwitchKey(e.Parameters, gadgetParams)

	e.EncryptFourierRLevPolyAssign(skIn.Value, ksk.Value)

	return ksk
}
