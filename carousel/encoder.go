package carousel

import (
	"math"

	"github.com/sp301415/carousel/math/num"
	"github.com/sp301415/carousel/math/poly"
	"github.com/sp301415/carousel/math/vec"
)

// Encoder encodes integer messages to plaintexts.
// Encoder is embedded in Encryptor and Evaluator,
// so usually manual instantiation isn't needed.
//
// Encoder is not safe for concurrent use.
// Use [*Encoder.ShallowCopy] to get a safe copy.
type Encoder struct {
	// Parameters is the parameters for this Encoder.
	Parameters Parameters

	// FourierTransformer is the FourierTransformer for this Encoder.
	FourierTransformer *poly.RealFourierTransformer

	// resolution is the fourier transformed resolution of unities.
	// Empty if parameters does not support packing.
	resolution []float64
	// resolutionInv is the fourier transformed inverse resolution of unities.
	// Empty if parameters does not support packing.
	resolutionInv []float64

	buffer encodeBuffer
}

// encodeBuffer is a buffer for Encoder.
type encodeBuffer struct {
	// fp is the intermediate fourier polynomial.
	fp poly.FourierPoly
}

// NewEncoder creates a new Encoder with given parameters.
func NewEncoder(params Parameters) *Encoder {
	fft := poly.NewRealFourierTransformer(params.polyDegree)

	var resolution, resolutionInv []float64
	if params.encodeType == EncodeTypeSlots {
		pParams := params.polyEvaluatorParameters

		resolutionRef := params.polyEvaluatorParameters.Resolution()
		resolutionInvRef := make([]uint64, params.polyDegree)
		for i := 0; i < params.polyDegree; i++ {
			resolutionInvRef[i] = uint64(pParams.CyclotomicDegree())*resolutionRef[i] + uint64(pParams.Order())
			resolutionInvRef[i] %= params.messageModulus
		}

		resolution = vec.Cast[uint64, float64](resolutionRef)
		resolutionInv = vec.Cast[uint64, float64](resolutionInvRef)

		fft.FourierTransformInPlace(resolution)
		fft.FourierTransformInPlace(resolutionInv)
		for i := 0; i < params.polyDegree; i++ {
			resolution[i] /= float64(params.polyDegree) / 2
			resolutionInv[i] /= float64(params.polyDegree) / 2
		}
	}

	return &Encoder{
		Parameters:         params,
		FourierTransformer: fft,

		resolution:    resolution,
		resolutionInv: resolutionInv,

		buffer: newEncodeBuffer(params),
	}
}

// newEncodeBuffer allocates a new encodeBuffer.
func newEncodeBuffer(params Parameters) encodeBuffer {
	return encodeBuffer{
		fp: poly.NewFourierPoly(params.polyDegree),
	}
}

// ShallowCopy returns a shallow copy of this Encoder.
// Returned Encoder is safe for concurrent use.
func (e *Encoder) ShallowCopy() *Encoder {
	return &Encoder{
		Parameters: e.Parameters,

		FourierTransformer: e.FourierTransformer,

		resolution:    e.resolution,
		resolutionInv: e.resolutionInv,

		buffer: newEncodeBuffer(e.Parameters),
	}
}

// EncodeLWE encodes integer message to LWE plaintext.
func (e *Encoder) EncodeLWE(message int) LWEPlaintext {
	return LWEPlaintext{Value: (uint64(message) % e.Parameters.messageModulus) * e.Parameters.scale}
}

// DecodeLWE decodes LWE plaintext to integer message.
func (e *Encoder) DecodeLWE(pt LWEPlaintext) int {
	return int(num.DivRound(pt.Value, e.Parameters.scale) % e.Parameters.messageModulus)
}

// EncodeRLWE encodes up to PolyDegree integer messages into one RLWE plaintext.
// Encoding type is automatically determined using Parameter's EncodeType.
//
//   - If len(messages) < PolyDegree, the leftovers are padded with zero.
//   - If len(messages) > PolyDegree, the leftovers are discarded.
func (e *Encoder) EncodeRLWE(messages []int) RLWEPlaintext {
	pt := NewRLWEPlaintext(e.Parameters)
	e.EncodeRLWEAssign(messages, pt)
	return pt
}

// EncodeRLWEAssign encodes up to PolyDegree integer messages into one RLWE plaintext.
// Encoding type is automatically determined using Parameter's EncodeType.
//
//   - If len(messages) < PolyDegree, the leftovers are padded with zero.
//   - If len(messages) > PolyDegree, the leftovers are discarded.
func (e *Encoder) EncodeRLWEAssign(messages []int, pt RLWEPlaintext) {
	switch e.Parameters.encodeType {
	case EncodeTypeCoeffs:
		e.EncodeRLWECoeffsAssign(messages, pt)
	case EncodeTypeSlots:
		e.EncodeRLWESlotsAssign(messages, pt)
	}
}

// DecodeRLWE decodes RLWE plaintext to integer messages.
// The returned messages are always of length PolyDegree.
func (e *Encoder) DecodeRLWE(pt RLWEPlaintext) []int {
	messages := make([]int, e.Parameters.polyDegree)
	e.DecodeRLWEAssign(pt, messages)
	return messages
}

// DecodeRLWEAssign decodes RLWE plaintext to integer messages.
//
//   - If len(messagesOut) < PolyDegree, the leftovers are discarded.
//   - If len(messagesOut) > PolyDegree, only the first PolyDegree elements are written.
func (e *Encoder) DecodeRLWEAssign(pt RLWEPlaintext, messagesOut []int) {
	switch e.Parameters.encodeType {
	case EncodeTypeCoeffs:
		e.DecodeRLWECoeffsAssign(pt, messagesOut)
	case EncodeTypeSlots:
		e.DecodeRLWESlotsAssign(pt, messagesOut)
	}
}

// EncodeRLWECoeffs encodes up to PolyDegree integer messages into one RLWE plaintext
// using coefficient packing.
//
//   - If len(messages) < PolyDegree, the leftovers are padded with zero.
//   - If len(messages) > PolyDegree, the leftovers are discarded.
func (e *Encoder) EncodeRLWECoeffs(messages []int) RLWEPlaintext {
	pt := NewRLWEPlaintext(e.Parameters)
	e.EncodeRLWECoeffsAssign(messages, pt)
	return pt
}

// EncodeRLWECoeffsAssign encodes up to PolyDegree integer messages into one RLWE plaintext
// using coefficient packing.
//
//   - If len(messages) < PolyDegree, the leftovers are padded with zero.
//   - If len(messages) > PolyDegree, the leftovers are discarded.
func (e *Encoder) EncodeRLWECoeffsAssign(messages []int, pt RLWEPlaintext) {
	length := num.Min(len(messages), e.Parameters.polyDegree)
	for i := 0; i < length; i++ {
		pt.Value.Coeffs[i] = (uint64(messages[i]) % e.Parameters.messageModulus) * e.Parameters.scale
	}
	vec.Fill(pt.Value.Coeffs[length:], 0)
}

// DecodeRLWECoeffs decodes coefficient-packed RLWE plaintext to integer messages.
// The returned messages are always of length PolyDegree.
func (e *Encoder) DecodeRLWECoeffs(pt RLWEPlaintext) []int {
	messages := make([]int, e.Parameters.polyDegree)
	e.DecodeRLWECoeffsAssign(pt, messages)
	return messages
}

// DecodeRLWECoeffsAssign decodes coefficient-packed RLWE plaintext to integer messages.
//
//   - If len(messagesOut) < PolyDegree, the leftovers are discarded.
//   - If len(messagesOut) > PolyDegree, only the first PolyDegree elements are written.
func (e *Encoder) DecodeRLWECoeffsAssign(pt RLWEPlaintext, messagesOut []int) {
	length := num.Min(len(messagesOut), e.Parameters.polyDegree)
	for i := 0; i < length; i++ {
		messagesOut[i] = int(num.DivRound(pt.Value.Coeffs[i], e.Parameters.scale) % e.Parameters.messageModulus)
	}
}

// EncodeRLWESlots encodes up to PolyDegree integer messages into one RLWE plaintext
// using slot packing.
//
//   - If len(messages) < PolyDegree, the leftovers are padded with zero.
//   - If len(messages) > PolyDegree, the leftovers are discarded.
//
// Panics if EncodeType is not EncodeTypeSlots.
func (e *Encoder) EncodeRLWESlots(messages []int) RLWEPlaintext {
	pt := NewRLWEPlaintext(e.Parameters)
	e.EncodeRLWESlotsAssign(messages, pt)
	return pt
}

// EncodeRLWESlotsAssign encodes up to PolyDegree integer messages into one RLWE plaintext
// using slot packing.
//
//   - If len(messages) < PolyDegree, the leftovers are padded with zero.
//   - If len(messages) > PolyDegree, the leftovers are discarded.
//
// Panics if EncodeType is not EncodeTypeSlots.
func (e *Encoder) EncodeRLWESlotsAssign(messages []int, pt RLWEPlaintext) {
	if e.Parameters.encodeType != EncodeTypeSlots {
		panic("EncodeRLWESlotsAssign: invalid encode type")
	}

	length := num.Min(len(messages), e.Parameters.polyDegree)
	for i := 0; i < length; i++ {
		e.buffer.fp.Coeffs[i] = float64(uint64(messages[i]) % e.Parameters.messageModulus)
	}
	vec.Fill(e.buffer.fp.Coeffs[length:], 0)
	vec.ReverseInPlace(e.buffer.fp.Coeffs[1:])

	e.FourierTransformer.FourierTransformInPlace(e.buffer.fp.Coeffs)
	e.FourierTransformer.ConvolveAssign(e.buffer.fp.Coeffs, e.resolution, e.buffer.fp.Coeffs)
	e.FourierTransformer.InvFourierTransformInPlace(e.buffer.fp.Coeffs)

	for i := 0; i < e.Parameters.polyDegree; i++ {
		pt.Value.Coeffs[i] = uint64(int64(math.Round(e.buffer.fp.Coeffs[i]))) % e.Parameters.messageModulus
		pt.Value.Coeffs[i] *= e.Parameters.scale
	}
}

// DecodeRLWESlots decodes slot-packed RLWE plaintext to integer messages.
// The returned messages are always of length PolyDegree.
func (e *Encoder) DecodeRLWESlots(pt RLWEPlaintext) []int {
	messages := make([]int, e.Parameters.polyDegree)
	e.DecodeRLWESlotsAssign(pt, messages)
	return messages
}

// DecodeRLWESlotsAssign decodes slot-packed RLWE plaintext to integer messages.
//
//   - If len(messagesOut) < PolyDegree, the leftovers are discarded.
//   - If len(messagesOut) > PolyDegree, only the first PolyDegree elements are written.
func (e *Encoder) DecodeRLWESlotsAssign(pt RLWEPlaintext, messagesOut []int) {
	if e.Parameters.encodeType != EncodeTypeSlots {
		panic("DecodeRLWESlotsAssign: invalid encode type")
	}

	for i := 0; i < e.Parameters.polyDegree; i++ {
		e.buffer.fp.Coeffs[i] = float64(num.DivRound(pt.Value.Coeffs[i], e.Parameters.scale) % e.Parameters.messageModulus)
	}
	vec.ReverseInPlace(e.buffer.fp.Coeffs[1:])

	e.FourierTransformer.FourierTransformInPlace(e.buffer.fp.Coeffs)
	e.FourierTransformer.ConvolveAssign(e.buffer.fp.Coeffs, e.resolutionInv, e.buffer.fp.Coeffs)
	e.FourierTransformer.InvFourierTransformInPlace(e.buffer.fp.Coeffs)

	length := num.Min(len(messagesOut), e.Parameters.polyDegree)
	for i := 0; i < length; i++ {
		messagesOut[i] = int(uint64(int64(math.Round(e.buffer.fp.Coeffs[i]))) % e.Parameters.messageModulus)
	}
}
