//go:build amd64 && !purego

package vec

import (
	"golang.org/x/sys/cpu"
)

// AddAssign computes vOut = v0 + v1.
func AddAssign(v0, v1, vOut []uint64) {
	if cpu.X86.HasAVX2 {
		addUint64AssignAVX2(v0, v1, vOut)
		return
	}

	for i := range vOut {
		vOut[i] = v0[i] + v1[i]
	}
}

// SubAssign computes vOut = v0 - v1.
func SubAssign(v0, v1, vOut []uint64) {
	if cpu.X86.HasAVX2 {
		subUint64AssignAVX2(v0, v1, vOut)
		return
	}

	for i := range vOut {
		vOut[i] = v0[i] - v1[i]
	}
}

// ScalarMulAssign computes vOut = c * v0.
func ScalarMulAssign(v0 []uint64, c uint64, vOut []uint64) {
	if cpu.X86.HasAVX2 {
		scalarMulUint64AssignAVX2(v0, c, vOut)
		return
	}

	for i := range vOut {
		vOut[i] = c * v0[i]
	}
}

// ScalarMulAddAssign computes vOut += c * v0.
func ScalarMulAddAssign(v0 []uint64, c uint64, vOut []uint64) {
	if cpu.X86.HasAVX2 {
		scalarMulAddUint64AssignAVX2(v0, c, vOut)
		return
	}

	for i := range vOut {
		vOut[i] += c * v0[i]
	}
}

// ScalarMulSubAssign computes vOut -= c * v0.
func ScalarMulSubAssign(v0 []uint64, c uint64, vOut []uint64) {
	if cpu.X86.HasAVX2 {
		scalarMulSubUint64AssignAVX2(v0, c, vOut)
		return
	}

	for i := range vOut {
		vOut[i] -= c * v0[i]
	}
}

// ElementWiseMulAssign computes vOut = v0 * v1, where * is an elementwise multiplication.
func ElementWiseMulAssign(v0, v1, vOut []uint64) {
	if cpu.X86.HasAVX2 {
		elementWiseMulUint64AssignAVX2(v0, v1, vOut)
		return
	}

	for i := range vOut {
		vOut[i] = v0[i] * v1[i]
	}
}

// ElementWiseMulAddAssign computes vOut += v0 * v1, where * is an elementwise multiplication.
func ElementWiseMulAddAssign(v0, v1, vOut []uint64) {
	if cpu.X86.HasAVX2 {
		elementWiseMulAddUint64AssignAVX2(v0, v1, vOut)
		return
	}

	for i := range vOut {
		vOut[i] += v0[i] * v1[i]
	}
}

// ElementWiseMulSubAssign computes vOut -= v0 * v1, where * is an elementwise multiplication.
func ElementWiseMulSubAssign(v0, v1, vOut []uint64) {
	if cpu.X86.HasAVX2 {
		elementWiseMulSubUint64AssignAVX2(v0, v1, vOut)
		return
	}

	for i := range vOut {
		vOut[i] -= v0[i] * v1[i]
	}
}
