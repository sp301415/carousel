package csprng_test

import (
	"math"
	"testing"

	"github.com/sp301415/carousel/math/csprng"
	"github.com/stretchr/testify/assert"
)

func meanStdDev(v []float64) (mean, stdDev float64) {
	sum := 0.0
	for _, x := range v {
		sum += x
	}

	mean = sum / float64(len(v))

	variance := 0.0
	for _, x := range v {
		variance += (x - mean) * (x - mean)
	}
	stdDev = math.Sqrt(variance / float64(len(v)))

	return
}

func TestGaussianSampler(t *testing.T) {
	mean := 0.0
	sigma := math.Exp2(16)

	gs := csprng.NewGaussianSampler()
	samples := make([]uint64, 1024)
	gs.SampleVecAssign(sigma, samples)
	samplesFloat := make([]float64, len(samples))
	for i, s := range samples {
		samplesFloat[i] = float64(int64(s))
	}
	meanSample, stdDevSample := meanStdDev(samplesFloat)

	k := 3.29 // From the GLITCH test suite
	N := float64(len(samples))
	meanBound := meanSample + k*stdDevSample/math.Sqrt(N)
	stdDevBound := stdDevSample + k*stdDevSample/math.Sqrt(2*(N-1))

	assert.GreaterOrEqual(t, meanBound, mean)
	assert.GreaterOrEqual(t, stdDevBound, sigma)
}
