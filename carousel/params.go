package carousel

import (
	"math"

	"github.com/sp301415/carousel/math/num"
	"github.com/sp301415/carousel/math/poly"
)

// GadgetParametersLiteral is a structure for Gadget Decomposition,
// which is used in Lev, GSW, GLev and RGSW encryptions.
type GadgetParametersLiteral struct {
	// Base is a base of gadget. It must be power of two.
	Base uint64
	// Level is a length of gadget.
	Level int
}

// WithBase sets the base and returns the new GadgetParametersLiteral.
func (p GadgetParametersLiteral) WithBase(base uint64) GadgetParametersLiteral {
	p.Base = base
	return p
}

// WithLevel sets the level and returns the new GadgetParametersLiteral.
func (p GadgetParametersLiteral) WithLevel(level int) GadgetParametersLiteral {
	p.Level = level
	return p
}

// Compile transforms GadgetParametersLiteral to read-only GadgetParameters.
// If there is any invalid parameter in the literal, it panics.
func (p GadgetParametersLiteral) Compile() GadgetParameters {
	switch {
	case p.Base < 2:
		panic("Base smaller than two")
	case !num.IsPowerOfTwo(p.Base):
		panic("Base not power of two")
	case p.Level <= 0:
		panic("Level smaller than zero")
	case 64 < num.Log2(p.Base)*p.Level:
		panic("Base * Level larger than Q")
	}

	return GadgetParameters{
		base:    p.Base,
		logBase: num.Log2(p.Base),
		level:   p.Level,
	}
}

// GadgetParameters is a read-only, compiled parameters based on GadgetParametersLiteral.
type GadgetParameters struct {
	// Base is a base of gadget. It must be power of two.
	base uint64
	// LogBase equals log(Base).
	logBase int
	// Level is a length of gadget.
	level int
}

// Base is a base of gadget. It must be power of two.
func (p GadgetParameters) Base() uint64 {
	return p.base
}

// LogBase equals log(Base).
func (p GadgetParameters) LogBase() int {
	return p.logBase
}

// Level is a length of gadget.
func (p GadgetParameters) Level() int {
	return p.level
}

// BaseQ returns Q / Base^(i+1) for 0 <= i < Level.
// For the most common usages i = 0 and i = Level-1, use [GadgetParameters.FirstBaseQ] and [GadgetParameters.LastBaseQ].
func (p GadgetParameters) BaseQ(i int) uint64 {
	return 1 << (64 - (i+1)*p.logBase)
}

// FirstBaseQ returns Q / Base.
func (p GadgetParameters) FirstBaseQ() uint64 {
	return 1 << (64 - p.logBase)
}

// LastBaseQ returns Q / Base^Level.
func (p GadgetParameters) LastBaseQ() uint64 {
	return 1 << (64 - p.level*p.logBase)
}

// LogBaseQ returns log(Q / Base^(i+1)) for 0 <= i < Level.
// For the most common usages i = 0 and i = Level-1, use [GadgetParameters.LogFirstBaseQ] and [GadgetParameters.LogLastBaseQ].
func (p GadgetParameters) LogBaseQ(i int) int {
	return 64 - (i+1)*p.logBase
}

// LogFirstBaseQ returns log(Q / Base).
func (p GadgetParameters) LogFirstBaseQ() int {
	return 64 - p.logBase
}

// LogLastBaseQ returns log(Q / Base^Level).
func (p GadgetParameters) LogLastBaseQ() int {
	return 64 - p.level*p.logBase
}

// Literal returns a GadgetParametersLiteral from this GadgetParameters.
func (p GadgetParameters) Literal() GadgetParametersLiteral {
	return GadgetParametersLiteral{
		Base:  p.base,
		Level: p.level,
	}
}

// EncodeType is an enum type for the type of encoding.
type EncodeType int

const (
	// EncodeTypeCoeffs uses the coefficient encoding to encode the message.
	EncodeTypeCoeffs EncodeType = iota
	// EncodeTypeSlots uses the slot encoding to encode the message.
	EncodeTypeSlots
)

// ParametersLiteral is a structure for Carousel parameters.
//
// # Warning
//
// Unless you are a cryptographic expert, DO NOT set these by yourself;
// always use the default parameters provided.
type ParametersLiteral struct {
	// LWEDimension is the dimension of LWE lattice used. Usually this is denoted by n.
	LWEDimension int
	// PolyDegree is the degree of polynomials in RLWE entities. Usually this is denoted by N.
	PolyDegree int

	// LWEStdDev is the normalized standard deviation used for gaussian error sampling in LWE encryption.
	LWEStdDev float64
	// RLWEStdDev is the normalized standard deviation used for gaussian error sampling in RLWE encryption.
	RLWEStdDev float64

	// BlockSize is the size of block to be used for LWE key sampling.
	//
	// This is used in Block Binary Key distribution, as explained in https://eprint.iacr.org/2023/958.
	// To use the original TFHE bootstrapping, set this to 1.
	//
	// If zero, then it is set to 1.
	BlockSize int

	// MessageModulus is the modulus of the encoded message.
	MessageModulus uint64

	// BlindRotateParameters is the gadget parameters for Blind Rotation.
	BlindRotateParameters GadgetParametersLiteral
	// KeySwitchParameters is the gadget parameters for KeySwitching.
	KeySwitchParameters GadgetParametersLiteral

	// EncodeType is the type of encoding.
	EncodeType EncodeType
}

// WithLWEDimension sets the LWEDimension and returns the new ParametersLiteral.
func (p ParametersLiteral) WithLWEDimension(lweDimension int) ParametersLiteral {
	p.LWEDimension = lweDimension
	return p
}

// WithPolyDegree sets the PolyDegree and returns the new ParametersLiteral.
func (p ParametersLiteral) WithPolyDegree(polyDegree int) ParametersLiteral {
	p.PolyDegree = polyDegree
	return p
}

// WithLWEStdDev sets the LWEStdDev and returns the new ParametersLiteral.
func (p ParametersLiteral) WithLWEStdDev(lweStdDev float64) ParametersLiteral {
	p.LWEStdDev = lweStdDev
	return p
}

// WithRLWEStdDev sets the RLWEStdDev and returns the new ParametersLiteral.
func (p ParametersLiteral) WithRLWEStdDev(rlweStdDev float64) ParametersLiteral {
	p.RLWEStdDev = rlweStdDev
	return p
}

// WithBlockSize sets the BlockSize and returns the new ParametersLiteral.
func (p ParametersLiteral) WithBlockSize(blockSize int) ParametersLiteral {
	p.BlockSize = blockSize
	return p
}

// WithMessageModulus sets the MessageModulus and returns the new ParametersLiteral.
func (p ParametersLiteral) WithMessageModulus(messageModulus uint64) ParametersLiteral {
	p.MessageModulus = messageModulus
	return p
}

// WithBlindRotateParameters sets the BlindRotateParameters and returns the new ParametersLiteral.
func (p ParametersLiteral) WithBlindRotateParameters(blindRotateParameters GadgetParametersLiteral) ParametersLiteral {
	p.BlindRotateParameters = blindRotateParameters
	return p
}

// WithKeySwitchParameters sets the KeySwitchParameters and returns the new ParametersLiteral.
func (p ParametersLiteral) WithKeySwitchParameters(keySwitchParameters GadgetParametersLiteral) ParametersLiteral {
	p.KeySwitchParameters = keySwitchParameters
	return p
}

// WithEncodeType sets the EncodeType and returns the new ParametersLiteral.
func (p ParametersLiteral) WithEncodeType(encodeType EncodeType) ParametersLiteral {
	p.EncodeType = encodeType
	return p
}

// Compile transforms ParametersLiteral to read-only Parameters.
// If there is any invalid parameter in the literal, it panics.
// Default parameters are guaranteed to be compiled without panics.
//
// # Warning
//
// This method performs only basic sanity checks.
// Just because a parameter compiles does not necessarily mean it is safe or correct.
// Unless you are a cryptographic expert, DO NOT set parameters by yourself;
// always use the default parameters provided.
func (p ParametersLiteral) Compile() Parameters {
	if p.BlockSize == 0 {
		p.BlockSize = 1
	}

	switch {
	case p.LWEDimension <= 0:
		panic("LWEDimension smaller than zero")
	case p.LWEDimension > p.PolyDegree:
		panic("LWEDimension larger than PolyDegree")
	case p.LWEStdDev <= 0:
		panic("LWEStdDev smaller than zero")
	case p.RLWEStdDev <= 0:
		panic("RLWEStdDev smaller than zero")
	case p.BlockSize <= 0:
		panic("BlockSize smaller than zero")
	case p.LWEDimension%p.BlockSize != 0:
		panic("LWEDimension not multiple of BlockSize")
	case !num.IsPowerOfTwo(p.PolyDegree):
		panic("PolyDegree not power of two")
	case !num.IsPowerOfTwo(p.MessageModulus):
		panic("MessageModulus not power of two")
	case !(p.EncodeType == EncodeTypeCoeffs || p.EncodeType == EncodeTypeSlots):
		panic("BootstrapOrder not valid")
	}

	var polyEvaluatorParameters poly.EvaluatorParameters
	switch p.EncodeType {
	case EncodeTypeCoeffs:
		polyEvaluatorParameters = poly.NewEvaluatorParameters(p.PolyDegree)
	case EncodeTypeSlots:
		polyEvaluatorParameters = poly.NewEvaluatorParametersForPacking(p.PolyDegree, num.Log2(p.MessageModulus))
	}

	return Parameters{
		lweDimension:            p.LWEDimension,
		polyDegree:              p.PolyDegree,
		logPolyDegree:           num.Log2(p.PolyDegree),
		polyEvaluatorParameters: polyEvaluatorParameters,
		lookUpTableSize:         p.PolyDegree,
		polyExtendFactor:        1,

		lweStdDev:  p.LWEStdDev,
		rlweStdDev: p.RLWEStdDev,

		blockSize:  p.BlockSize,
		blockCount: p.LWEDimension / p.BlockSize,

		messageModulus:    p.MessageModulus,
		logMessageModulus: num.Log2(p.MessageModulus),
		scale:             1 << (64 - num.Log2(p.MessageModulus)),
		logScale:          64 - num.Log2(p.MessageModulus),

		blindRotateParameters: p.BlindRotateParameters.Compile(),
		keySwitchParameters:   p.KeySwitchParameters.Compile(),

		encodeType: p.EncodeType,
	}
}

// Parameters are read-only, compiled parameters based on ParametersLiteral.
type Parameters struct {
	// LWEDimension is the dimension of LWE lattice used. Usually this is denoted by n.
	lweDimension int
	// PolyDegree is the degree of polynomials in RLWE entities. Usually this is denoted by N.
	polyDegree int
	// LogPolyDegree equals log(PolyDegree).
	logPolyDegree int
	// polyEvaluatorParameters is the parameters for PolyEvaluator.
	polyEvaluatorParameters poly.EvaluatorParameters
	// LookUpTableSize is the size of Lookup Table used in Blind Rotation.
	lookUpTableSize int
	// PolyExtendFactor equals LookUpTableSize / PolyDegree.
	polyExtendFactor int

	// LWEStdDev is the normalized standard deviation used for gaussian error sampling in LWE encryption.
	lweStdDev float64
	// RLWEStdDev is the normalized standard deviation used for gaussian error sampling in RLWE encryption.
	rlweStdDev float64

	// BlockSize is the size of block to be used for LWE key sampling.
	blockSize int
	// BlockCount is a number of blocks in LWESecretkey. Equal to LWEDimension / BlockSize.
	blockCount int

	// MessageModulus is the modulus of the encoded message.
	messageModulus uint64
	// logMessageModulus is the value of log(MessageModulus).
	logMessageModulus int
	// Scale is the scaling factor used for message encoding.
	// The lower log(Scale) bits are reserved for errors.
	scale uint64
	// LogScale is the value of log(Scale).
	logScale int

	// blindRotateParameters is the gadget parameters for Blind Rotation.
	blindRotateParameters GadgetParameters
	// keySwitchParameters is the gadget parameters for KeySwitching.
	keySwitchParameters GadgetParameters

	// encodeType is the type of encoding.
	encodeType EncodeType
}

// LWEDimension is the dimension of LWE lattice used. Usually this is denoted by n.
func (p Parameters) LWEDimension() int {
	return p.lweDimension
}

// PolyDegree is the degree of polynomials in RLWE entities. Usually this is denoted by N.
func (p Parameters) PolyDegree() int {
	return p.polyDegree
}

// LogPolyDegree equals log(PolyDegree).
func (p Parameters) LogPolyDegree() int {
	return p.logPolyDegree
}

// PolyEvaluatorParameters is the parameters for PolyEvaluator.
func (p Parameters) PolyEvaluatorParameters() poly.EvaluatorParameters {
	return p.polyEvaluatorParameters
}

// LookUpTableSize is the size of LookUpTable used in Blind Rotation.
func (p Parameters) LookUpTableSize() int {
	return p.lookUpTableSize
}

// PolyExtendFactor returns LookUpTableSize / PolyDegree.
func (p Parameters) PolyExtendFactor() int {
	return p.polyExtendFactor
}

// LWEStdDev is the standard deviation used for gaussian error sampling in LWE encryption.
func (p Parameters) LWEStdDev() float64 {
	return p.lweStdDev
}

// RLWEStdDev is the standard deviation used for gaussian error sampling in RLWE encryption.
func (p Parameters) RLWEStdDev() float64 {
	return p.rlweStdDev
}

// BlockSize is the size of block to be used for LWE key sampling.
func (p Parameters) BlockSize() int {
	return p.blockSize
}

// BlockCount is a number of blocks in LWESecretkey. Equal to LWEDimension / BlockSize.
func (p Parameters) BlockCount() int {
	return p.blockCount
}

// Scale is the scaling factor used for message encoding.
// The lower log(Scale) bits are reserved for errors.
func (p Parameters) Scale() uint64 {
	return p.scale
}

// MessageModulus is the modulus of the encoded message.
func (p Parameters) MessageModulus() uint64 {
	return p.messageModulus
}

// LogMessageModulus is the value of log(MessageModulus).
func (p Parameters) LogMessageModulus() int {
	return p.logMessageModulus
}

// BlindRotateParameters is the gadget parameters for Programmable Bootstrapping.
func (p Parameters) BlindRotateParameters() GadgetParameters {
	return p.blindRotateParameters
}

// KeySwitchParameters is the gadget parameters for KeySwitching.
func (p Parameters) KeySwitchParameters() GadgetParameters {
	return p.keySwitchParameters
}

// EncodeType is the type of encoding.
func (p Parameters) EncodeType() EncodeType {
	return p.encodeType
}

// Literal returns a ParametersLiteral from this Parameters.
func (p Parameters) Literal() ParametersLiteral {
	return ParametersLiteral{
		LWEDimension: p.lweDimension,
		PolyDegree:   p.polyDegree,

		LWEStdDev:  p.lweStdDev,
		RLWEStdDev: p.rlweStdDev,

		BlockSize: p.blockSize,

		MessageModulus: p.messageModulus,

		BlindRotateParameters: p.blindRotateParameters.Literal(),
		KeySwitchParameters:   p.keySwitchParameters.Literal(),

		EncodeType: p.encodeType,
	}
}

// EstimateModSwitchStdDev returns an estimated standard deviation of error from modulus switching.
func (p Parameters) EstimateModSwitchStdDev() float64 {
	L := float64(p.lookUpTableSize)
	q := math.Exp2(64)

	h := float64(p.blockCount) * (float64(p.blockSize)) / (float64(p.blockSize + 1))

	modSwitchVar := ((h + 1) * q * q) / (12 * L * L)

	return math.Sqrt(modSwitchVar)
}

// EstimateBlindRotateStdDev returns an estimated standard deviation of error from Blind Rotation.
func (p Parameters) EstimateBlindRotateStdDev() float64 {
	n := float64(p.lweDimension)
	N := float64(p.polyDegree)
	M := float64(p.polyEvaluatorParameters.CyclotomicDegree())
	o := float64(p.polyEvaluatorParameters.Order())
	beta := p.rlweStdDev
	q := math.Exp2(64)

	h := float64(p.blockCount) * (float64(p.blockSize)) / (float64(p.blockSize + 1))

	Bbr := float64(p.blindRotateParameters.Base())
	Lbr := float64(p.blindRotateParameters.Level())

	rotVar1 := n * (o * (h + (N-n)/2) * q * q) / (12 * math.Pow(Bbr, 2*Lbr))
	rotVar2 := n * (M * Lbr * Bbr * Bbr * beta * beta) / 12
	rotVarFFT := n * math.Exp2(-100) * (o * (h + (N-n)/2)) * M * (q * q) * Lbr * (Bbr * Bbr)
	rotVar := rotVar1 + rotVar2 + rotVarFFT

	muxVar1 := h * ((1 + o*(h+(N-n)/2)) * q * q) / (12 * math.Pow(Bbr, 2*Lbr))
	muxVar2 := n * (M * 2 * Lbr * Bbr * Bbr * beta * beta) / 12
	muxVarFFT := n * math.Exp2(-100) * (1 + o*(h+(N-n)/2)) * M * (q * q) * 2 * Lbr * (Bbr * Bbr)
	muxVar := muxVar1 + muxVar2 + muxVarFFT

	return math.Sqrt(rotVar + muxVar)
}

// EstimateBlindRotateExtractStdDev returns an estimated standard deviation of error from Blind Rotation and Sample Extraction.
func (p Parameters) EstimateBlindRotateExtractStdDev() float64 {
	switch p.encodeType {
	case EncodeTypeCoeffs:
		return p.EstimateBlindRotateStdDev()
	case EncodeTypeSlots:
		t := float64(p.messageModulus)
		M := float64(p.polyEvaluatorParameters.CyclotomicDegree())

		return math.Sqrt(M*t*t/12) * p.EstimateBlindRotateStdDev()
	}

	return 0
}

// EstimateKeySwitchForBootstrapStdDev returns an estimated standard deviation of error from Key Switching for bootstrapping.
func (p Parameters) EstimateKeySwitchForBootstrapStdDev() float64 {
	n := float64(p.lweDimension)
	N := float64(p.polyDegree)
	alpha := p.lweStdDev
	q := math.Exp2(64)

	Bks := float64(p.keySwitchParameters.Base())
	Lks := float64(p.keySwitchParameters.Level())

	keySwitchVar1 := ((N - n) / 2) * (q * q) / (12 * math.Pow(Bks, 2*Lks))
	keySwitchVar2 := (N - n) * (alpha * alpha * Lks * Bks * Bks) / 12
	keySwitchVar := keySwitchVar1 + keySwitchVar2

	return math.Sqrt(keySwitchVar)
}

// EstimateMaxErrorStdDev returns an estimated standard deviation of maximum possible error.
func (p Parameters) EstimateMaxErrorStdDev() float64 {
	modSwitchStdDev := p.EstimateModSwitchStdDev()
	blindRotateStdDev := p.EstimateBlindRotateExtractStdDev()
	keySwitchStdDev := p.EstimateKeySwitchForBootstrapStdDev()

	return math.Sqrt(modSwitchStdDev*modSwitchStdDev + blindRotateStdDev*blindRotateStdDev + keySwitchStdDev*keySwitchStdDev)
}

// EstimateFailureProbability returns the failure probability of bootstrapping.
func (p Parameters) EstimateFailureProbability() float64 {
	bound := math.Exp2(64) / (2 * float64(p.messageModulus))
	return math.Erfc(bound / (math.Sqrt2 * p.EstimateMaxErrorStdDev()))
}
