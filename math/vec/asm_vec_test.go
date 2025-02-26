package vec_test

import (
	"math/rand"
	"testing"

	"github.com/sp301415/carousel/math/vec"
	"github.com/stretchr/testify/assert"
)

func TestVec(t *testing.T) {
	N := 687

	v0 := make([]uint64, N)
	v1 := make([]uint64, N)
	vOut := make([]uint64, N)
	vOutAVX := make([]uint64, N)

	for i := 0; i < N; i++ {
		v0[i] = rand.Uint64()
		v1[i] = rand.Uint64()
	}

	t.Run("AddAssign", func(t *testing.T) {
		for i := 0; i < N; i++ {
			vOut[i] = v0[i] + v1[i]
		}
		vec.AddAssign(v0, v1, vOutAVX)

		assert.Equal(t, vOut, vOutAVX)
	})

	t.Run("SubAssign", func(t *testing.T) {
		for i := 0; i < N; i++ {
			vOut[i] = v0[i] - v1[i]
		}
		vec.SubAssign(v0, v1, vOutAVX)

		assert.Equal(t, vOut, vOutAVX)
	})

	t.Run("ScalarMulAssign", func(t *testing.T) {
		cv := vOut[0]
		for i := 0; i < N; i++ {
			vOut[i] = cv * v0[i]
		}
		vec.ScalarMulAssign(v0, cv, vOutAVX)

		assert.Equal(t, vOut, vOutAVX)
	})

	t.Run("ScalarMulAddAssign", func(t *testing.T) {
		cv := vOut[0]
		for i := 0; i < N; i++ {
			vOut[i] += cv * v0[i]
		}
		vec.ScalarMulAddAssign(v0, cv, vOutAVX)

		assert.Equal(t, vOut, vOutAVX)
	})

	t.Run("ScalarMulSubAssign", func(t *testing.T) {
		cv := vOut[0]
		for i := 0; i < N; i++ {
			vOut[i] -= cv * v0[i]
		}
		vec.ScalarMulSubAssign(v0, cv, vOutAVX)

		assert.Equal(t, vOut, vOutAVX)
	})

	t.Run("ElementWiseMulAssign", func(t *testing.T) {
		for i := 0; i < N; i++ {
			vOut[i] = v0[i] * v1[i]
		}
		vec.ElementWiseMulAssign(v0, v1, vOutAVX)

		assert.Equal(t, vOut, vOutAVX)
	})

	t.Run("ElementWiseMulAddAssign", func(t *testing.T) {
		for i := 0; i < N; i++ {
			vOut[i] += v0[i] * v1[i]
		}
		vec.ElementWiseMulAddAssign(v0, v1, vOutAVX)

		assert.Equal(t, vOut, vOutAVX)
	})

	t.Run("ElementWiseMulSubAssign", func(t *testing.T) {
		for i := 0; i < N; i++ {
			vOut[i] -= v0[i] * v1[i]
		}
		vec.ElementWiseMulSubAssign(v0, v1, vOutAVX)

		assert.Equal(t, vOut, vOutAVX)
	})
}
