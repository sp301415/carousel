//go:build amd64 && !purego

package poly

import (
	"math"

	"golang.org/x/sys/cpu"
)

// convertUint64ToFloat64Assign casts p to fpOut.
func convertUint64ToFloat64Assign(p []uint64, fpOut []float64) {
	if cpu.X86.HasAVX2 && cpu.X86.HasFMA {
		convertUint64ToFloat64AssignAVX2(p, fpOut)
		return
	}

	for i := 0; i < len(fpOut); i++ {
		fpOut[i] = float64(int64(p[i]))
	}
}

// convertFloat64ToUint64Assign casts fp to pOut.
func convertFloat64ToUint64Assign(fp []float64, pOut []uint64) {
	if cpu.X86.HasAVX2 && cpu.X86.HasFMA {
		convertFloat64ToUint64AssignAVX2(fp, pOut)
		return
	}

	for i := 0; i < len(pOut); i++ {
		pOut[i] = uint64(int64(fp[i]))
	}
}

// floatModQInPlace computes fp mod Q in-place.
func floatModQInPlace(coeffs []float64) {
	if cpu.X86.HasAVX2 {
		floatModQInPlaceAVX2(coeffs)
		return
	}

	q := math.Exp2(64)
	for i := 0; i < len(coeffs); i++ {
		coeffs[i] = math.Round(coeffs[i] - q*math.Round(coeffs[i]/q))
	}
}
