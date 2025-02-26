// Package poly implements polynomial and its operations.
package poly

import (
	"math"

	"github.com/sp301415/carousel/math/num"
	"github.com/sp301415/carousel/math/vec"
)

// Poly is a polynomial over a subring of a cyclotomic ring.
type Poly struct {
	Coeffs []uint64
}

// NewPoly creates a polynomial with degree N with empty coefficients.
//
// Panics when N is not a power of two, or when N is smaller than MinDegree or larger than MaxDegree.
func NewPoly(N int) Poly {
	switch {
	case !num.IsPowerOfTwo(N):
		panic("degree not power of two")
	case N < MinDegree:
		panic("degree smaller than MinDegree")
	}

	return Poly{Coeffs: make([]uint64, N)}
}

// From creates a new polynomial from given coefficient slice.
// The given slice is copied, and extended to degree N.
func From(coeffs []uint64, N int) Poly {
	p := NewPoly(N)
	vec.CopyAssign(coeffs, p.Coeffs)
	return p
}

// Copy returns a copy of the polynomial.
func (p Poly) Copy() Poly {
	return Poly{Coeffs: vec.Copy(p.Coeffs)}
}

// CopyFrom copies p0 to p.
func (p *Poly) CopyFrom(p0 Poly) {
	vec.CopyAssign(p0.Coeffs, p.Coeffs)
}

// Degree returns the degree of the polynomial.
// This is equivalent with length of coefficients.
func (p Poly) Degree() int {
	return len(p.Coeffs)
}

// Clear clears all the coefficients to zero.
func (p Poly) Clear() {
	vec.Fill(p.Coeffs, 0)
}

// Equals checks if p0 is equal with p.
func (p Poly) Equals(p0 Poly) bool {
	return vec.Equals(p.Coeffs, p0.Coeffs)
}

// FourierPoly is a fourier transformed polynomial.
type FourierPoly struct {
	// Coeffs is ordered in reverse for efficient computation.
	//
	// Namely,
	//
	//  [a[0], a[1], a[2], ..., a[N-2], a[N-1]]
	//
	// is represented as
	//
	//  [a[0], a[N-1], a[N-2], ..., a[2], a[1]]
	//
	Coeffs []float64
}

// NewFourierPoly creates a fourier polynomial with degree N/2 with empty coefficients.
//
// Panics when N is not a power of two, or when N is smaller than MinDegree.
func NewFourierPoly(N int) FourierPoly {
	switch {
	case !num.IsPowerOfTwo(N):
		panic("degree not power of two")
	case N < MinDegree:
		panic("degree smaller than MinDegree")
	}

	return FourierPoly{Coeffs: make([]float64, N)}
}

// Degree returns the degree of the polynomial.
func (p FourierPoly) Degree() int {
	return len(p.Coeffs)
}

// Copy returns a copy of the polynomial.
func (p FourierPoly) Copy() FourierPoly {
	return FourierPoly{Coeffs: vec.Copy(p.Coeffs)}
}

// CopyFrom copies p0 to p.
func (p *FourierPoly) CopyFrom(p0 FourierPoly) {
	vec.CopyAssign(p0.Coeffs, p.Coeffs)
}

// Clear clears all the coefficients to zero.
func (p FourierPoly) Clear() {
	vec.Fill(p.Coeffs, 0)
}

// Equals checks if p0 is equal with p.
// Note that due to floating point errors,
// this function may return false even if p0 and p are equal.
func (p FourierPoly) Equals(p0 FourierPoly) bool {
	return vec.Equals(p.Coeffs, p0.Coeffs)
}

// Approx checks if p0 is approximately equal with p,
// with a difference smaller than eps.
func (p FourierPoly) Approx(p0 FourierPoly, eps float64) bool {
	for i := range p.Coeffs {
		if math.Abs(p.Coeffs[i]-p0.Coeffs[i]) > eps {
			return false
		}
	}
	return true
}
