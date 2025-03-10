package carousel

var (
	// ParamsSlotsUint2 is the parameters for carousel-slots with 2 bits of message space.
	ParamsSlotsUint2 = ParametersLiteral{
		LWEDimension: 628,
		PolyDegree:   2048,

		LWEStdDev:  1090715534753795.8,
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

	// ParamsSlotsUint4 is the parameters for carousel-slots with 4 bits of message space.
	ParamsSlotsUint4 = ParametersLiteral{
		LWEDimension: 717,
		PolyDegree:   2048,

		LWEStdDev:  215504279044099.94,
		RLWEStdDev: 6148,

		BlockSize: 3,

		MessageModulus: 1 << 4,

		BlindRotateParameters: GadgetParametersLiteral{
			Base:  1 << 8,
			Level: 4,
		},
		KeySwitchParameters: GadgetParametersLiteral{
			Base:  1 << 3,
			Level: 4,
		},

		EncodeType: EncodeTypeSlots,
	}
)
