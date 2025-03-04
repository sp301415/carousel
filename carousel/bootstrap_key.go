package carousel

// EvaluationKey is a public key for Evaluator,
// which consists of BlindRotation Key and KeySwitching Key.
// All keys should be treated as read-only.
// Changing them mid-operation will usually result in wrong results.
type EvaluationKey struct {
	// BlindRotateKey is a blindrotate key.
	BlindRotateKey BlindRotateKey
	// KeySwitchKey is a keyswitch key switching LWELargeKey -> LWEKey.
	KeySwitchKey LWEKeySwitchKey
	// AutomorphismKey is a key for automorphism.
	AutomorphismKey []RLWEKeySwitchKey
}

// NewEvaluationKey creates a new EvaluationKey.
func NewEvaluationKey(params Parameters) EvaluationKey {
	atk := make([]RLWEKeySwitchKey, params.polyDegree)
	for i := 0; i < params.polyDegree; i++ {
		atk[i] = NewRLWEKeySwitchKey(params, params.blindRotateParameters)
	}

	return EvaluationKey{
		BlindRotateKey:  NewBlindRotateKey(params),
		KeySwitchKey:    NewKeySwitchKeyForBootstrap(params),
		AutomorphismKey: atk,
	}
}

// NewEvaluationKeyCustom creates a new EvaluationKey with custom parameters.
func NewEvaluationKeyCustom(lweDimension, polyDegree int, blindRotateParams, keySwitchParams GadgetParameters) EvaluationKey {
	atk := make([]RLWEKeySwitchKey, polyDegree)
	for i := 0; i < polyDegree; i++ {
		atk[i] = NewRLWEKeySwitchKeyCustom(polyDegree, blindRotateParams)
	}

	return EvaluationKey{
		BlindRotateKey:  NewBlindRotateKeyCustom(lweDimension, polyDegree, blindRotateParams),
		KeySwitchKey:    NewKeySwitchKeyForBootstrapCustom(lweDimension, polyDegree, keySwitchParams),
		AutomorphismKey: atk,
	}
}

// Copy returns a copy of the key.
func (evk EvaluationKey) Copy() EvaluationKey {
	atkCopy := make([]RLWEKeySwitchKey, len(evk.AutomorphismKey))
	for i := range evk.AutomorphismKey {
		atkCopy[i] = evk.AutomorphismKey[i].Copy()
	}

	return EvaluationKey{
		BlindRotateKey:  evk.BlindRotateKey.Copy(),
		KeySwitchKey:    evk.KeySwitchKey.Copy(),
		AutomorphismKey: atkCopy,
	}
}

// CopyFrom copies values from key.
func (evk *EvaluationKey) CopyFrom(evkIn EvaluationKey) {
	for i := range evk.AutomorphismKey {
		evk.AutomorphismKey[i].CopyFrom(evkIn.AutomorphismKey[i])
	}
	evk.BlindRotateKey.CopyFrom(evkIn.BlindRotateKey)
	evk.KeySwitchKey.CopyFrom(evkIn.KeySwitchKey)
}

// Clear clears the key.
func (evk *EvaluationKey) Clear() {
	evk.BlindRotateKey.Clear()
	evk.KeySwitchKey.Clear()
	for i := range evk.AutomorphismKey {
		evk.AutomorphismKey[i].Clear()
	}
}

// BlindRotateKey is a key for blind rotation.
// Essentially, this is a RGSW encryption of LWEKey with RLWEKey.
// However, FFT is already applied for fast external product.
type BlindRotateKey struct {
	GadgetParameters GadgetParameters

	// Value has length LWEDimension.
	Value []FourierRGSWCiphertext
}

// NewBlindRotateKey creates a new BlindRotateKey.
func NewBlindRotateKey(params Parameters) BlindRotateKey {
	brk := make([]FourierRGSWCiphertext, params.lweDimension)
	for i := 0; i < params.lweDimension; i++ {
		brk[i] = NewFourierRGSWCiphertext(params, params.blindRotateParameters)
	}
	return BlindRotateKey{Value: brk, GadgetParameters: params.blindRotateParameters}
}

// NewBlindRotateKeyCustom creates a new BlindRotateKey with custom parameters.
func NewBlindRotateKeyCustom(lweDimension, polyDegree int, gadgetParams GadgetParameters) BlindRotateKey {
	brk := make([]FourierRGSWCiphertext, lweDimension)
	for i := 0; i < lweDimension; i++ {
		brk[i] = NewFourierRGSWCiphertextCustom(polyDegree, gadgetParams)
	}
	return BlindRotateKey{Value: brk, GadgetParameters: gadgetParams}
}

// Copy returns a copy of the key.
func (brk BlindRotateKey) Copy() BlindRotateKey {
	brkCopy := make([]FourierRGSWCiphertext, len(brk.Value))
	for i := range brk.Value {
		brkCopy[i] = brk.Value[i].Copy()
	}
	return BlindRotateKey{Value: brkCopy, GadgetParameters: brk.GadgetParameters}
}

// CopyFrom copies values from key.
func (brk *BlindRotateKey) CopyFrom(brkIn BlindRotateKey) {
	for i := range brk.Value {
		brk.Value[i].CopyFrom(brkIn.Value[i])
	}
	brk.GadgetParameters = brkIn.GadgetParameters
}

// Clear clears the key.
func (brk *BlindRotateKey) Clear() {
	for i := range brk.Value {
		brk.Value[i].Clear()
	}
}
