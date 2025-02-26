//go:build !(amd64 && !purego)

package poly

// addFloat64Assign computes vOut = v0 + v1.
func addFloat64Assign(v0, v1, vOut []float64) {
	for i := range vOut {
		vOut[i] = v0[i] + v1[i]
	}
}

// subFloat64Assign computes vOut = v0 - v1.
func subFloat64Assign(v0, v1, vOut []float64) {
	for i := range vOut {
		vOut[i] = v0[i] - v1[i]
	}
}

// negFloat64Assign computes vOut = -v0.
func negFloat64Assign(v0, vOut []float64) {
	for i := range vOut {
		vOut[i] = -v0[i]
	}
}

// scalarMulFloat64Assign computes vOut = c * v0.
func scalarMulFloat64Assign(v0 []float64, c float64, vOut []float64) {
	for i := range vOut {
		vOut[i] = c * v0[i]
	}
}

// scalarMulAddFloat64Assign computes vOut = v0 + c*v1.
func scalarMulAddFloat64Assign(v0 []float64, c float64, vOut []float64) {
	for i := range vOut {
		vOut[i] += c * v0[i]
	}
}

// scalarMulSubFloat64Assign computes vOut = v0 - c*v1.
func scalarMulSubFloat64Assign(v0 []float64, c float64, vOut []float64) {
	for i := range vOut {
		vOut[i] -= c * v0[i]
	}
}

// elementWiseMulFloat64Assign computes vOut = v0 * v1.
func elementWiseMulFloat64Assign(v0, v1, vOut []float64) {
	for i := range vOut {
		vOut[i] = v0[i] * v1[i]
	}
}

// elementWiseMulAddFloat64Assign computes vOut += v0 +*v1.
func elementWiseMulAddFloat64Assign(v0, v1, vOut []float64) {
	for i := range vOut {
		vOut[i] += v0[i] * v1[i]
	}
}

// elementWiseMulSubFloat64Assign computes vOut -= v0 * v1.
func elementWiseMulSubFloat64Assign(v0, v1, vOut []float64) {
	for i := range vOut {
		vOut[i] -= v0[i] * v1[i]
	}
}
