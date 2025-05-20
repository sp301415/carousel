package carousel

import "math"

var (
	// ParamsCoeffsUint2 is the parameters for carousel-coeffs with 2 bits of message space.
	ParamsCoeffsUint2 = ParametersLiteral{
		LWEDimension: 687,
		PolyDegree:   2048,

		LWEStdDev:  math.Exp2(49.07),
		RLWEStdDev: math.Exp2(13.05),

		BlockSize: 3,

		MessageModulus: 1 << 2,

		BlindRotateParameters: GadgetParametersLiteral{
			Base:  1 << 20,
			Level: 1,
		},
		KeySwitchParameters: GadgetParametersLiteral{
			Base:  1 << 2,
			Level: 6,
		},

		EncodeType: EncodeTypeCoeffs,
	}

	// ParamsCoeffsUint3 is the parameters for carousel-coeffs with 3 bits of message space.
	ParamsCoeffsUint3 = ParametersLiteral{
		LWEDimension: 788,
		PolyDegree:   2048,

		LWEStdDev:  math.Exp2(46.60),
		RLWEStdDev: math.Exp2(13.09),

		BlockSize: 4,

		MessageModulus: 1 << 3,

		BlindRotateParameters: GadgetParametersLiteral{
			Base:  1 << 22,
			Level: 1,
		},
		KeySwitchParameters: GadgetParametersLiteral{
			Base:  1 << 4,
			Level: 3,
		},

		EncodeType: EncodeTypeCoeffs,
	}

	// ParamsCoeffsUint4 is the parameters for carousel-coeffs with 4 bits of message space.
	ParamsCoeffsUint4 = ParametersLiteral{
		LWEDimension: 828,
		PolyDegree:   2048,

		LWEStdDev:  math.Exp2(45.57),
		RLWEStdDev: math.Exp2(13.09),

		BlockSize: 4,

		MessageModulus: 1 << 4,

		BlindRotateParameters: GadgetParametersLiteral{
			Base:  1 << 23,
			Level: 1,
		},
		KeySwitchParameters: GadgetParametersLiteral{
			Base:  1 << 3,
			Level: 5,
		},

		EncodeType: EncodeTypeCoeffs,
	}

	// ParamsSlotsUint2 is the parameters for carousel-slots with 2 bits of message space.
	ParamsSlotsUint2 = ParametersLiteral{
		LWEDimension: 687,
		PolyDegree:   2048,

		LWEStdDev:  math.Exp2(49.07),
		RLWEStdDev: math.Exp2(13.05),

		BlockSize: 3,

		MessageModulus: 1 << 2,

		BlindRotateParameters: GadgetParametersLiteral{
			Base:  1 << 14,
			Level: 2,
		},
		KeySwitchParameters: GadgetParametersLiteral{
			Base:  1 << 2,
			Level: 6,
		},

		EncodeType: EncodeTypeSlots,
	}

	// ParamsSlotsUint3 is the parameters for carousel-slots with 3 bits of message space.
	ParamsSlotsUint3 = ParametersLiteral{
		LWEDimension: 788,
		PolyDegree:   2048,

		LWEStdDev:  math.Exp2(46.60),
		RLWEStdDev: math.Exp2(13.09),

		BlockSize: 4,

		MessageModulus: 1 << 3,

		BlindRotateParameters: GadgetParametersLiteral{
			Base:  1 << 16,
			Level: 2,
		},
		KeySwitchParameters: GadgetParametersLiteral{
			Base:  1 << 4,
			Level: 3,
		},

		EncodeType: EncodeTypeSlots,
	}

	// ParamsSlotsUint4 is the parameters for carousel-slots with 4 bits of message space.
	ParamsSlotsUint4 = ParametersLiteral{
		LWEDimension: 828,
		PolyDegree:   2048,

		LWEStdDev:  math.Exp2(45.57),
		RLWEStdDev: math.Exp2(13.09),

		BlockSize: 4,

		MessageModulus: 1 << 4,

		BlindRotateParameters: GadgetParametersLiteral{
			Base:  1 << 17,
			Level: 2,
		},
		KeySwitchParameters: GadgetParametersLiteral{
			Base:  1 << 3,
			Level: 5,
		},

		EncodeType: EncodeTypeSlots,
	}
)
