//go:build amd64 && !purego

package carousel

import (
	"unsafe"

	"github.com/sp301415/carousel/math/num"
	"github.com/sp301415/carousel/math/poly"
	"golang.org/x/sys/cpu"
)

// decomposePolyAssign decomposes p with respect to gadgetParams and writes it to decomposedOut.
func decomposePolyAssign(p poly.Poly, gadgetParams GadgetParameters, decomposedOut []poly.Poly) {
	if cpu.X86.HasAVX2 && cpu.X86.HasFMA {
		decomposePolyAssignUint64AVX2(
			p.Coeffs,
			gadgetParams.base,
			uint64(gadgetParams.logBase),
			uint64(gadgetParams.LogLastBaseQ()),
			*(*[][]uint64)(unsafe.Pointer(&decomposedOut)),
		)
		return
	}

	logLastBaseQ := gadgetParams.LogLastBaseQ()
	for i := 0; i < p.Degree(); i++ {
		c := num.DivRoundBits(p.Coeffs[i], logLastBaseQ)
		for j := gadgetParams.level - 1; j >= 1; j-- {
			decomposedOut[j].Coeffs[i] = c & (gadgetParams.base - 1)
			c >>= gadgetParams.logBase
			c += decomposedOut[j].Coeffs[i] >> (gadgetParams.logBase - 1)
			decomposedOut[j].Coeffs[i] -= (decomposedOut[j].Coeffs[i] & (gadgetParams.base >> 1)) << 1
		}
		decomposedOut[0].Coeffs[i] = c & (gadgetParams.base - 1)
		decomposedOut[0].Coeffs[i] -= (decomposedOut[0].Coeffs[i] & (gadgetParams.base >> 1)) << 1
	}
}
