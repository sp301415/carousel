package carousel

var (
	// ParamsSlotsUint2 is the parameters for carousel-slots with 2 bits of message space.
	ParamsSlotsUint2 = ParametersLiteral{
		LWEDimension: 630,
		PolyDegree:   2048,

		LWEStdDev:  1046735069642755.8,
		RLWEStdDev: 6148,

		BlockSize: 2,

		MessageModulus: 1 << 2,

		BlindRotateParameters: GadgetParametersLiteral{
			Base:  1 << 10,
			Level: 3,
		},
		KeySwitchParameters: GadgetParametersLiteral{
			Base:  1 << 2,
			Level: 6,
		},

		EncodeType: EncodeTypeSlots,
	}

	// ParamsSlotsUint3 is the parameters for carousel-slots with 3 bits of message space.
	ParamsSlotsUint3 = ParametersLiteral{
		LWEDimension: 696,
		PolyDegree:   2048,

		LWEStdDev:  424411488321539.9,
		RLWEStdDev: 6148,

		BlockSize: 3,

		MessageModulus: 1 << 3,

		BlindRotateParameters: GadgetParametersLiteral{
			Base:  1 << 22,
			Level: 1,
		},
		KeySwitchParameters: GadgetParametersLiteral{
			Base:  1 << 5,
			Level: 3,
		},

		EncodeType: EncodeTypeCoeffs,
	}
)
