package poly

import (
	"math/rand"
	"testing"

	"github.com/sp301415/carousel/math/vec"
	"github.com/stretchr/testify/assert"
)

func TestVecCmplxAssembly(t *testing.T) {
	N := 1 << 10
	eps := 1e-10

	r := rand.New(rand.NewSource(0))

	v0 := make([]float64, N)
	v1 := make([]float64, N)
	for i := 0; i < N; i++ {
		v0[i] = r.Float64()
		v1[i] = r.Float64()
	}

	vOut := make([]float64, N)
	vOutAVX2 := make([]float64, N)

	t.Run("Add", func(t *testing.T) {
		for i := 0; i < N; i++ {
			vOut[i] = v0[i] + v1[i]
		}
		addFloat64Assign(v0, v1, vOutAVX2)
		assert.InEpsilonSlice(t, vOut, vOutAVX2, eps)
	})

	t.Run("Sub", func(t *testing.T) {
		for i := 0; i < N; i++ {
			vOut[i] = v0[i] - v1[i]
		}
		subFloat64Assign(v0, v1, vOutAVX2)
		assert.InEpsilonSlice(t, vOut, vOutAVX2, eps)
	})

	t.Run("Neg", func(t *testing.T) {
		for i := 0; i < N; i++ {
			vOut[i] = -v0[i]
		}
		negFloat64Assign(v0, vOutAVX2)
		assert.InEpsilonSlice(t, vOut, vOutAVX2, eps)
	})

	t.Run("ScalarMul", func(t *testing.T) {
		c := r.Float64()
		for i := 0; i < N; i++ {
			vOut[i] = c * v0[i]
		}
		scalarMulFloat64Assign(v0, c, vOutAVX2)
		assert.InEpsilonSlice(t, vOut, vOutAVX2, eps)
	})

	t.Run("ScalarMulAdd", func(t *testing.T) {
		vec.Fill(vOut, 0)
		vec.Fill(vOutAVX2, 0)

		c := r.Float64()
		for i := 0; i < N; i++ {
			vOut[i] += c * v0[i]
		}
		scalarMulAddFloat64Assign(v0, c, vOutAVX2)
		assert.InEpsilonSlice(t, vOut, vOutAVX2, eps)
	})

	t.Run("ScalarMulSub", func(t *testing.T) {
		vec.Fill(vOut, 0)
		vec.Fill(vOutAVX2, 0)

		c := r.Float64()
		for i := 0; i < N; i++ {
			vOut[i] -= c * v0[i]
		}
		scalarMulSubFloat64Assign(v0, c, vOutAVX2)
		assert.InEpsilonSlice(t, vOut, vOutAVX2, eps)
	})

	t.Run("Mul", func(t *testing.T) {
		for i := 0; i < N; i++ {
			vOut[i] = v0[i] * v1[i]
		}
		elementWiseMulFloat64Assign(v0, v1, vOutAVX2)
		assert.InEpsilonSlice(t, vOut, vOutAVX2, eps)
	})

	t.Run("MulAdd", func(t *testing.T) {
		vec.Fill(vOut, 0)
		vec.Fill(vOutAVX2, 0)

		for i := 0; i < N; i++ {
			vOut[i] += v0[i] * v1[i]
		}
		elementWiseMulAddFloat64Assign(v0, v1, vOutAVX2)
		assert.InEpsilonSlice(t, vOut, vOutAVX2, eps)
	})

	t.Run("MulSub", func(t *testing.T) {
		vec.Fill(vOut, 0)
		vec.Fill(vOutAVX2, 0)

		for i := 0; i < N; i++ {
			vOut[i] -= v0[i] * v1[i]
		}
		elementWiseMulSubFloat64Assign(v0, v1, vOutAVX2)
		assert.InEpsilonSlice(t, vOut, vOutAVX2, eps)
	})
}
