package csprng

import (
	"github.com/sp301415/carousel/math/poly"
	"github.com/sp301415/carousel/math/vec"
)

// BinarySampler samples values from uniform and block binary distribution.
type BinarySampler struct {
	baseSampler *UniformSampler
}

// NewBinarySampler creates a new BinarySampler.
//
// Panics when read from crypto/rand or blake2b initialization fails.
func NewBinarySampler() *BinarySampler {
	return &BinarySampler{
		baseSampler: NewUniformSampler(),
	}
}

// NewBinarySamplerWithSeed creates a new BinarySampler, with user supplied seed.
//
// Panics when blake2b initialization fails.
func NewBinarySamplerWithSeed(seed []byte) *BinarySampler {
	return &BinarySampler{
		baseSampler: NewUniformSamplerWithSeed(seed),
	}
}

// Sample uniformly samples a random binary integer.
func (s *BinarySampler) Sample() uint64 {
	return s.baseSampler.Sample() & 1
}

// SampleVecAssign samples uniform binary values to vOut.
func (s *BinarySampler) SampleVecAssign(vOut []uint64) {
	var buf uint64
	for i := 0; i < len(vOut); i++ {
		if i&63 == 0 {
			buf = s.baseSampler.Sample()
		}
		vOut[i] = buf & 1
		buf >>= 1
	}
}

// SamplePolyAssign samples uniform binary values to pOut.
func (s *BinarySampler) SamplePolyAssign(pOut poly.Poly) {
	s.SampleVecAssign(pOut.Coeffs)
}

// SampleBlockVecAssign samples block binary values to vOut.
func (s *BinarySampler) SampleBlockVecAssign(blockSize int, vOut []uint64) {
	if len(vOut)%blockSize != 0 {
		panic("length not multiple of blocksize")
	}

	for i := 0; i < len(vOut); i += blockSize {
		vec.Fill(vOut[i:i+blockSize], 0)
		offset := int(s.baseSampler.SampleN(uint64(blockSize) + 1))
		if offset == blockSize {
			continue
		}
		vOut[i+offset] = 1
	}
}

// SampleBlockPolyAssign samples block binary values to pOut.
func (s *BinarySampler) SampleBlockPolyAssign(blockSize int, pOut poly.Poly) {
	s.SampleBlockVecAssign(blockSize, pOut.Coeffs)
}
