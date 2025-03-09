package carousel

var (
	// ParamsSlotsUint2 is the parameters for carousel-slots with 2 bits of message space.
	ParamsSlotsUint2 = ParametersLiteral{
		LWEDimension: 608,
		PolyDegree:   2048,

		LWEStdDev:  1597516338809342.5,
		RLWEStdDev: 6148,

		BlockSize: 2,

		MessageModulus: 1 << 2,

		BlindRotateParameters: GadgetParametersLiteral{
			Base:  1 << 13,
			Level: 2,
		},
		KeySwitchParameters: GadgetParametersLiteral{
			Base:  1 << 3,
			Level: 4,
		},

		EncodeType: EncodeTypeSlots,
	}

	// ParamsSlotsUint3 is the parameters for carousel-slots with 3 bits of message space.
	ParamsSlotsUint3 = ParametersLiteral{
		LWEDimension: 687,
		PolyDegree:   2048,

		LWEStdDev:  373833953443844,
		RLWEStdDev: 6148,

		BlockSize: 3,

		MessageModulus: 1 << 3,

		BlindRotateParameters: GadgetParametersLiteral{
			Base:  1 << 10,
			Level: 3,
		},
		KeySwitchParameters: GadgetParametersLiteral{
			Base:  1 << 4,
			Level: 3,
		},

		EncodeType: EncodeTypeSlots,
	}
)
