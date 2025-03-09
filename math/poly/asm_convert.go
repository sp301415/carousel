//go:build !(amd64 && !purego)

package poly

import "math"

// convertUint64ToFloat64Assign casts p to fpOut.
func convertUint64ToFloat64Assign(p []uint64, fpOut []float64) {
	for i := 0; i < len(fpOut); i++ {
		fpOut[i] = float64(int64(p[i]))
	}
}

// convertFloat64ToUint64Assign casts fp to pOut.
func convertFloat64ToUint64Assign(fp []float64, pOut []uint64) {
	for i := 0; i < len(pOut); i++ {
		pOut[i] = uint64(int64(fp[i]))
	}
}

// floatModQInPlace computes fp mod Q in-place.
func floatModQInPlace(coeffs []float64) {
	q := math.Exp2(64)
	for i := 0; i < len(coeffs); i++ {
		coeffs[i] = math.Round(coeffs[i] - q*math.Round(coeffs[i]/q))
	}
}
