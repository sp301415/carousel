package poly

import (
	"math/big"

	"github.com/sp301415/carousel/math/vec"
)

// z2xPoly is a polynomial in Z_2^r[X] (or its quotient ring).
type z2xPoly []uint64

// degree returns the degree of the polynomial.
func (p z2xPoly) degree() int {
	return len(p) - 1
}

// copy returns a copy of the polynomial.
func (p z2xPoly) copy() z2xPoly {
	return vec.Copy(p)
}

// isZero returns true if the polynomial is zero.
func (p z2xPoly) isZero() bool {
	for i := 0; i < len(p); i++ {
		if p[i] != 0 {
			return false
		}
	}

	return true
}

// isOne returns true if the polynomial is one.
func (p z2xPoly) isOne() bool {
	if p[0] != 1 {
		return false
	}

	for i := 1; i < len(p); i++ {
		if p[i] != 0 {
			return false
		}
	}

	return true
}

// z2xQuoRing is the polynomial ring Z_2^r[X] / (f).
type z2xQuoRing struct {
	r      int
	mask   uint64
	quo    z2xPoly
	bufMul z2xPoly
}

// newZ2XQuoRing returns Z_2^r[X] / (f).
func newZ2XQuoRing(r int, f z2xPoly) z2xQuoRing {
	return z2xQuoRing{
		r:      r,
		mask:   (1 << r) - 1,
		quo:    f.copy(),
		bufMul: make(z2xPoly, 2*len(f)),
	}
}

// newPoly returns a new polynomial in ring.
func (r z2xQuoRing) newPoly() z2xPoly {
	return make(z2xPoly, r.quo.degree())
}

// neg returns pOut = -p0.
func (r z2xQuoRing) neg(p0 z2xPoly) z2xPoly {
	pOut := r.newPoly()
	for i := 0; i < r.quo.degree(); i++ {
		pOut[i] = -p0[i] & r.mask
	}
	return pOut
}

// sub returns pOut = p0 - p1.
func (r z2xQuoRing) sub(p0, p1 z2xPoly) z2xPoly {
	pOut := r.newPoly()
	for i := 0; i < r.quo.degree(); i++ {
		pOut[i] = (p0[i] - p1[i]) & r.mask
	}
	return pOut
}

// mulAssign computes pOut = p0 * p1.
func (r z2xQuoRing) mul(p0, p1 z2xPoly) z2xPoly {
	vec.Fill(r.bufMul, 0)
	for i := 0; i < r.quo.degree(); i++ {
		for j := 0; j < r.quo.degree(); j++ {
			r.bufMul[i+j] += p0[i] * p1[j]
		}
	}

	return r.reduce(r.bufMul)
}

// reduce returns pOut = p0 mod f.
func (r z2xQuoRing) reduce(p0 z2xPoly) z2xPoly {
	pOut := z2xPoly(vec.Copy(p0))
	for i := len(p0) - 1; i >= r.quo.degree(); i-- {
		if pOut[i] == 0 {
			continue
		}

		c := pOut[i]
		for j := 0; j <= r.quo.degree(); j++ {
			pOut[i-j] -= c * r.quo[r.quo.degree()-j]
			pOut[i-j] &= r.mask
		}
	}

	return pOut[:r.quo.degree()]
}

// expBig returns pOut = p0^e.
func (r z2xQuoRing) expBig(p0 z2xPoly, e *big.Int) z2xPoly {
	switch e.Cmp(big.NewInt(0)) {
	case -1:
		panic("expBigAssign: negative exponent")
	case 0:
		pOut := r.newPoly()
		pOut[0] = 1
		return pOut
	}

	pOut := r.newPoly()
	pOut[0] = 1

	p0Copy := p0.copy()
	eCopy := big.NewInt(0).Set(e)
	for eCopy.Cmp(big.NewInt(0)) > 0 {
		if eCopy.Bit(0) == 1 {
			pOut = r.mul(pOut, p0Copy)
		}
		p0Copy = r.mul(p0Copy, p0Copy)
		eCopy.Rsh(eCopy, 1)
	}
	return pOut
}
