package carousel

import (
	"github.com/sp301415/carousel/math/num"
	"github.com/sp301415/carousel/math/poly"
	"github.com/sp301415/carousel/math/vec"
)

// LookUpTable is a polynomial that is the lookup table
// for function evaluations during programmable bootstrapping.
type LookUpTable struct {
	// Value has length polyExtendFactor.
	Value []poly.Poly
}

// NewLookUpTable creates a new lookup table.
func NewLookUpTable(params Parameters) LookUpTable {
	lut := make([]poly.Poly, params.polyExtendFactor)
	for i := 0; i < params.polyExtendFactor; i++ {
		lut[i] = poly.NewPoly(params.polyDegree)
	}

	return LookUpTable{Value: lut}
}

// NewLookUpTableCustom creates a new lookup table with custom size.
func NewLookUpTableCustom(extendFactor, polyDegree int) LookUpTable {
	lut := make([]poly.Poly, extendFactor)
	for i := 0; i < extendFactor; i++ {
		lut[i] = poly.NewPoly(polyDegree)
	}

	return LookUpTable{Value: lut}
}

// Copy returns a copy of the LUT.
func (lut LookUpTable) Copy() LookUpTable {
	lutCopy := make([]poly.Poly, len(lut.Value))
	for i := 0; i < len(lut.Value); i++ {
		lutCopy[i] = lut.Value[i].Copy()
	}
	return LookUpTable{Value: lutCopy}
}

// CopyFrom copies values from the LUT.
func (lut *LookUpTable) CopyFrom(lutIn LookUpTable) {
	for i := 0; i < len(lut.Value); i++ {
		lut.Value[i].CopyFrom(lutIn.Value[i])
	}
}

// Clear clears the LUT.
func (lut *LookUpTable) Clear() {
	for i := 0; i < len(lut.Value); i++ {
		lut.Value[i].Clear()
	}
}

// GenLookUpTableCoeffs generates a lookup table based on function f.
// Input and output of f is cut by MessageModulus.
func (e *Evaluator) GenLookUpTable(f func(int) int) LookUpTable {
	lutOut := NewLookUpTable(e.Parameters)
	e.GenLookUpTableAssign(f, lutOut)
	return lutOut
}

// GenLookUpTableCoeffsAssign generates a lookup table based on function f and writes it to lutOut.
// Input and output of f is cut by MessageModulus.
func (e *Evaluator) GenLookUpTableAssign(f func(int) int, lutOut LookUpTable) {
	switch e.Parameters.encodeType {
	case EncodeTypeCoeffs:
		e.GenLookUpTableCoeffsAssign(f, lutOut)
	case EncodeTypeSlots:
		e.GenLookUpTableSlotsAssign(f, lutOut)
	}
}

// GenLookUpTableCoeffs generates a coefficient encoded lookup table based on function f.
// Input and output of f is cut by MessageModulus.
func (e *Evaluator) GenLookUpTableCoeffs(f func(int) int) LookUpTable {
	lutOut := NewLookUpTable(e.Parameters)
	e.GenLookUpTableCoeffsAssign(f, lutOut)
	return lutOut
}

// GenLookUpTableCoeffsAssign generates a coefficient encoded lookup table based on function f and writes it to lutOut.
// Input and output of f is cut by MessageModulus.
func (e *Evaluator) GenLookUpTableCoeffsAssign(f func(int) int, lutOut LookUpTable) {
	e.genLookUpTableRawAssign(f, e.buffer.lutRaw)
	for i := 0; i < e.Parameters.polyExtendFactor; i++ {
		vec.ReverseInPlace(e.buffer.lutRaw[i])
		e.Encoder.EncodeRLWECoeffsAssign(e.buffer.lutRaw[i], RLWEPlaintext{Value: lutOut.Value[i]})
	}
}

// GenLookUpTableCoeffs generates a slot encoded lookup table based on function f.
// Input and output of f is cut by MessageModulus.
func (e *Evaluator) GenLookUpTableSlots(f func(int) int) LookUpTable {
	lutOut := NewLookUpTable(e.Parameters)
	e.GenLookUpTableSlotsAssign(f, lutOut)
	return lutOut
}

// GenLookUpTableCoeffsAssign generates a slot encoded lookup table based on function f and writes it to lutOut.
// Input and output of f is cut by MessageModulus.
func (e *Evaluator) GenLookUpTableSlotsAssign(f func(int) int, lutOut LookUpTable) {
	e.genLookUpTableRawAssign(f, e.buffer.lutRaw)
	for i := 0; i < e.Parameters.polyExtendFactor; i++ {
		e.Encoder.EncodeRLWESlotsAssign(e.buffer.lutRaw[i], RLWEPlaintext{Value: lutOut.Value[i]})
	}
}

// genLookUpTableRawAssign generates a lookup table based on function f and writes it to lutOut.
// Input and output of f and lutOut is cut by MessageModulus.
func (e *Evaluator) genLookUpTableRawAssign(f func(int) int, lutRawOut [][]int) {
	for x := 0; x < int(e.Parameters.messageModulus); x++ {
		start := num.DivRound(x*e.Parameters.lookUpTableSize, int(e.Parameters.messageModulus))
		end := num.DivRound((x+1)*e.Parameters.lookUpTableSize, int(e.Parameters.messageModulus))
		y := f(x) % int(e.Parameters.messageModulus)
		for xx := start; xx < end; xx++ {
			e.buffer.lutReorder[xx] = y
		}
	}

	offset := num.DivRound(e.Parameters.lookUpTableSize, int(2*e.Parameters.messageModulus))
	vec.RotateInPlace(e.buffer.lutReorder, -offset)
	for i := e.Parameters.lookUpTableSize - offset; i < e.Parameters.lookUpTableSize; i++ {
		e.buffer.lutReorder[i] = -e.buffer.lutReorder[i]
	}

	for i := 0; i < e.Parameters.polyExtendFactor; i++ {
		for j := 0; j < e.Parameters.polyDegree; j++ {
			lutRawOut[i][j] = e.buffer.lutReorder[j*e.Parameters.polyExtendFactor+i]
		}
	}
}
