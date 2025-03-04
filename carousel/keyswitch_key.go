package carousel

// LWEKeySwitchKey is a LWE keyswitch key from one LWEKey to another LWEKey.
type LWEKeySwitchKey struct {
	GadgetParameters GadgetParameters

	// Value has length InputLWEDimension.
	Value []LevCiphertext
}

// NewLWEKeySwitchKey creates a new LWEKeySwitchingKey.
func NewLWEKeySwitchKey(params Parameters, inputDimension int, gadgetParams GadgetParameters) LWEKeySwitchKey {
	ksk := make([]LevCiphertext, inputDimension)
	for i := 0; i < inputDimension; i++ {
		ksk[i] = NewLevCiphertext(params, gadgetParams)
	}
	return LWEKeySwitchKey{Value: ksk, GadgetParameters: gadgetParams}
}

// NewLWEKeySwitchKeyCustom creates a new LWEKeySwitchingKey with custom parameters.
func NewLWEKeySwitchKeyCustom(inputDimension, outputDimension int, gadgetParams GadgetParameters) LWEKeySwitchKey {
	ksk := make([]LevCiphertext, inputDimension)
	for i := 0; i < inputDimension; i++ {
		ksk[i] = NewLevCiphertextCustom(outputDimension, gadgetParams)
	}
	return LWEKeySwitchKey{Value: ksk, GadgetParameters: gadgetParams}
}

// NewKeySwitchKeyForBootstrap creates a new LWEKeySwitchingKey for bootstrapping.
func NewKeySwitchKeyForBootstrap(params Parameters) LWEKeySwitchKey {
	return NewLWEKeySwitchKeyCustom(params.polyDegree-params.lweDimension, params.lweDimension, params.keySwitchParameters)
}

// NewKeySwitchKeyForBootstrapCustom creates a new LWEKeySwitchingKey with custom parameters.
func NewKeySwitchKeyForBootstrapCustom(lweDimension, polyDegree int, gadgetParams GadgetParameters) LWEKeySwitchKey {
	return NewLWEKeySwitchKeyCustom(polyDegree-lweDimension, lweDimension, gadgetParams)
}

// InputLWEDimension returns the input LWEDimension of this key.
func (ksk LWEKeySwitchKey) InputLWEDimension() int {
	return len(ksk.Value)
}

// Copy returns a copy of the key.
func (ksk LWEKeySwitchKey) Copy() LWEKeySwitchKey {
	kskCopy := make([]LevCiphertext, len(ksk.Value))
	for i := range ksk.Value {
		kskCopy[i] = ksk.Value[i].Copy()
	}
	return LWEKeySwitchKey{Value: kskCopy, GadgetParameters: ksk.GadgetParameters}
}

// CopyFrom copies values from key.
func (ksk *LWEKeySwitchKey) CopyFrom(kskIn LWEKeySwitchKey) {
	for i := range ksk.Value {
		ksk.Value[i].CopyFrom(kskIn.Value[i])
	}
	ksk.GadgetParameters = kskIn.GadgetParameters
}

// Clear clears the key.
func (ksk *LWEKeySwitchKey) Clear() {
	for i := range ksk.Value {
		ksk.Value[i].Clear()
	}
}

// RLWEKeySwitchKey is a RLWE keyswitch key from one RLWEKey to another RLWEKey.
type RLWEKeySwitchKey struct {
	GadgetParameters GadgetParameters

	// Value is a FourierRLevCiphertext.
	Value FourierRLevCiphertext
}

// NewRLWEKeySwitchKey creates a new RLWEKeySwitchingKey.
func NewRLWEKeySwitchKey(params Parameters, gadgetParams GadgetParameters) RLWEKeySwitchKey {
	return RLWEKeySwitchKey{
		Value:            NewFourierRLevCiphertext(params, gadgetParams),
		GadgetParameters: gadgetParams,
	}
}

// NewRLWEKeySwitchKeyCustom creates a new RLWEKeySwitchingKey with custom parameters.
func NewRLWEKeySwitchKeyCustom(polyDegree int, gadgetParams GadgetParameters) RLWEKeySwitchKey {
	return RLWEKeySwitchKey{
		Value:            NewFourierRLevCiphertextCustom(polyDegree, gadgetParams),
		GadgetParameters: gadgetParams,
	}
}

// Copy returns a copy of the key.
func (ksk RLWEKeySwitchKey) Copy() RLWEKeySwitchKey {
	return RLWEKeySwitchKey{Value: ksk.Value.Copy(), GadgetParameters: ksk.GadgetParameters}
}

// CopyFrom copies values from key.
func (ksk *RLWEKeySwitchKey) CopyFrom(kskIn RLWEKeySwitchKey) {
	ksk.Value.CopyFrom(kskIn.Value)
	ksk.GadgetParameters = kskIn.GadgetParameters
}

// Clear clears the key.
func (ksk *RLWEKeySwitchKey) Clear() {
	ksk.Value.Clear()
}
