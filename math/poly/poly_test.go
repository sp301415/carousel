package poly_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/sp301415/carousel/math/poly"
)

var (
	LogN = []int{9, 10, 11, 12, 13, 14, 15}
)

func BenchmarkOps(b *testing.B) {
	for _, logN := range LogN {
		N := 1 << logN

		pParams := poly.NewEvaluatorParameters(N)
		ev := poly.NewEvaluator(pParams)

		p0 := ev.NewPoly()
		p1 := ev.NewPoly()
		pOut := ev.NewPoly()

		for i := 0; i < ev.Parameters.Degree(); i++ {
			p0.Coeffs[i] = rand.Uint64()
			p1.Coeffs[i] = rand.Uint64()
		}

		fp := ev.ToFourierPoly(p0)

		b.Run(fmt.Sprintf("LogN=%v/op=Add", logN), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ev.AddPolyAssign(p0, p1, pOut)
			}
		})

		b.Run(fmt.Sprintf("LogN=%v/op=Sub", logN), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ev.SubPolyAssign(p0, p1, pOut)
			}
		})

		b.Run(fmt.Sprintf("LogN=%v/op=BinaryMul", logN), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ev.ShortFourierPolyMulPolyAssign(p0, fp, pOut)
			}
		})

		b.Run(fmt.Sprintf("LogN=%v/op=Mul", logN), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ev.MulPolyAssign(p0, p1, pOut)
			}
		})
	}
}

func BenchmarkFourierOps(b *testing.B) {
	for _, logN := range LogN {
		N := 1 << logN

		pParams := poly.NewEvaluatorParameters(N)
		ev := poly.NewEvaluator(pParams)

		fp0 := ev.NewFourierPoly()
		fp1 := ev.NewFourierPoly()
		fpOut := ev.NewFourierPoly()

		for i := 0; i < ev.Parameters.Degree(); i++ {
			fp0.Coeffs[i] = (2*rand.Float64() - 1.0) * math.Exp(63)
			fp1.Coeffs[i] = (2*rand.Float64() - 1.0) * math.Exp(63)
		}

		b.Run(fmt.Sprintf("LogN=%v/op=Add", logN), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ev.AddFourierPolyAssign(fp0, fp1, fpOut)
			}
		})

		b.Run(fmt.Sprintf("LogN=%v/op=Sub", logN), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ev.SubFourierPolyAssign(fp0, fp1, fpOut)
			}
		})

		b.Run(fmt.Sprintf("LogN=%v/op=Mul", logN), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ev.MulFourierPolyAssign(fp0, fp1, fpOut)
			}
		})
	}
}

func BenchmarkFourierTransform(b *testing.B) {
	for _, logN := range LogN {
		N := 1 << logN

		pParams := poly.NewEvaluatorParameters(N)
		ev := poly.NewEvaluator(pParams)

		p := ev.NewPoly()
		fp := ev.NewFourierPoly()

		for i := 0; i < ev.Parameters.Degree(); i++ {
			p.Coeffs[i] = rand.Uint64()
		}

		b.Run(fmt.Sprintf("LogN=%v/op=ToFourierPoly", logN), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ev.ToFourierPolyAssign(p, fp)
			}
		})

		b.Run(fmt.Sprintf("LogN=%v/op=ToPoly", logN), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ev.ToPolyAssignUnsafe(fp, p)
			}
		})
	}
}
