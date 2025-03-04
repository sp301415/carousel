package carousel

var (
	// ParamsSlotsUint3 is the parameters for carousel-slots with 3 bits of message space.
	ParamsSlotsUint3 = ParametersLiteral{
		LWEDimension: 680,
		PolyDegree:   2048,

		LWEStdDev:  430131865291752.9,
		RLWEStdDev: 6405.772103413834,

		BlockSize: 1,

		MessageModulus: 1 << 3,

		BlindRotateParameters: GadgetParametersLiteral{
			Base:  1 << 8,
			Level: 4,
		},
		KeySwitchParameters: GadgetParametersLiteral{
			Base:  1 << 2,
			Level: 7,
		},

		EncodeType: EncodeTypeSlots,
	}
)
