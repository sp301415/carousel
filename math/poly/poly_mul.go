package poly

// MulPoly returns p0 * p1.
func (e *Evaluator) MulPoly(p0, p1 Poly) Poly {
	pOut := e.NewPoly()
	e.MulPolyAssign(p0, p1, pOut)
	return pOut
}

// MulPolyAssign computes pOut = p0 * p1.
func (e *Evaluator) MulPolyAssign(p0, p1, pOut Poly) {
	splitBits, splitCount := splitParameters(e.Parameters.degree)

	fp0Split := make([]FourierPoly, splitCount)
	fp1Split := make([]FourierPoly, splitCount)
	fpOutSplit := make([]FourierPoly, splitCount)
	for i := 0; i < splitCount; i++ {
		fp0Split[i] = e.NewFourierPoly()
		fp1Split[i] = e.NewFourierPoly()
		fpOutSplit[i] = e.NewFourierPoly()
	}

	var splitMask uint64 = 1<<splitBits - 1
	for i := 0; i < splitCount; i++ {
		splitLowBits := i * splitBits
		for j := 0; j < e.Parameters.degree; j++ {
			e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] >> splitLowBits) & splitMask
		}
		e.ToFourierPolyAssign(e.buffer.pSplit, fp0Split[i])

		for j := 0; j < e.Parameters.degree; j++ {
			e.buffer.pSplit.Coeffs[j] = (p1.Coeffs[j] >> splitLowBits) & splitMask
		}
		e.ToFourierPolyAssign(e.buffer.pSplit, fp1Split[i])
	}

	for i := 0; i < splitCount; i++ {
		for j := 0; j < splitCount-i; j++ {
			e.MulAddFourierPolyAssign(fp0Split[i], fp1Split[j], fpOutSplit[i+j])
		}
	}

	e.ToPolyAssignUnsafe(fpOutSplit[0], pOut)
	for i := 1; i < splitCount; i++ {
		e.ToPolyAssignUnsafe(fpOutSplit[i], e.buffer.pSplit)
		splitLowBits := i * splitBits
		for j := 0; j < e.Parameters.degree; j++ {
			pOut.Coeffs[j] += e.buffer.pSplit.Coeffs[j] << splitLowBits
		}
	}
}

// MulAddPolyAssign computes pOut += p0 * p1.
func (e *Evaluator) MulAddPolyAssign(p0, p1, pOut Poly) {
	splitBits, splitCount := splitParameters(e.Parameters.degree)

	fp0Split := make([]FourierPoly, splitCount)
	fp1Split := make([]FourierPoly, splitCount)
	fpOutSplit := make([]FourierPoly, splitCount)
	for i := 0; i < splitCount; i++ {
		fp0Split[i] = e.NewFourierPoly()
		fp1Split[i] = e.NewFourierPoly()
		fpOutSplit[i] = e.NewFourierPoly()
	}

	var splitMask uint64 = 1<<splitBits - 1
	for i := 0; i < splitCount; i++ {
		splitLowBits := i * splitBits
		for j := 0; j < e.Parameters.degree; j++ {
			e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] >> splitLowBits) & splitMask
		}
		e.ToFourierPolyAssign(e.buffer.pSplit, fp0Split[i])

		for j := 0; j < e.Parameters.degree; j++ {
			e.buffer.pSplit.Coeffs[j] = (p1.Coeffs[j] >> splitLowBits) & splitMask
		}
		e.ToFourierPolyAssign(e.buffer.pSplit, fp1Split[i])
	}

	for i := 0; i < splitCount; i++ {
		for j := 0; j < splitCount-i; j++ {
			e.MulAddFourierPolyAssign(fp0Split[i], fp1Split[j], fpOutSplit[i+j])
		}
	}

	e.ToPolyAddAssignUnsafe(fpOutSplit[0], pOut)
	for i := 1; i < splitCount; i++ {
		e.ToPolyAssignUnsafe(fpOutSplit[i], e.buffer.pSplit)
		splitLowBits := i * splitBits
		for j := 0; j < e.Parameters.degree; j++ {
			pOut.Coeffs[j] += e.buffer.pSplit.Coeffs[j] << splitLowBits
		}
	}
}

// MulSubPolyAssign computes pOut -= p0 * p1.
func (e *Evaluator) MulSubPolyAssign(p0, p1, pOut Poly) {
	splitBits, splitCount := splitParameters(e.Parameters.degree)

	fp0Split := make([]FourierPoly, splitCount)
	fp1Split := make([]FourierPoly, splitCount)
	fpOutSplit := make([]FourierPoly, splitCount)
	for i := 0; i < splitCount; i++ {
		fp0Split[i] = e.NewFourierPoly()
		fp1Split[i] = e.NewFourierPoly()
		fpOutSplit[i] = e.NewFourierPoly()
	}

	var splitMask uint64 = 1<<splitBits - 1
	for i := 0; i < splitCount; i++ {
		splitLowBits := i * splitBits
		for j := 0; j < e.Parameters.degree; j++ {
			e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] >> splitLowBits) & splitMask
		}
		e.ToFourierPolyAssign(e.buffer.pSplit, fp0Split[i])

		for j := 0; j < e.Parameters.degree; j++ {
			e.buffer.pSplit.Coeffs[j] = (p1.Coeffs[j] >> splitLowBits) & splitMask
		}
		e.ToFourierPolyAssign(e.buffer.pSplit, fp1Split[i])
	}

	for i := 0; i < splitCount; i++ {
		for j := 0; j < splitCount-i; j++ {
			e.MulAddFourierPolyAssign(fp0Split[i], fp1Split[j], fpOutSplit[i+j])
		}
	}

	e.ToPolySubAssignUnsafe(fpOutSplit[0], pOut)
	for i := 1; i < splitCount; i++ {
		e.ToPolyAssignUnsafe(fpOutSplit[i], e.buffer.pSplit)
		splitLowBits := i * splitBits
		for j := 0; j < e.Parameters.degree; j++ {
			pOut.Coeffs[j] -= e.buffer.pSplit.Coeffs[j] << splitLowBits
		}
	}
}

// ShortFourierPolyMulPoly returns p0 * fpShort, under the assumption that fpShort is a short polynomial.
// (i.e., all coefficients are bounded by [ShortLogBound] bits.)
// This is faster than [*Evaluator.MulPoly], and the result is exact unlike [*Evaluator.FourierPolyMulPoly].
func (e *Evaluator) ShortFourierPolyMulPoly(p0 Poly, fpShort FourierPoly) Poly {
	pOut := e.NewPoly()
	e.ShortFourierPolyMulPolyAssign(p0, fpShort, pOut)
	return pOut
}

// ShortFourierPolyMulPolyAssign computes pOut = p0 * fpShort, under the assumption that fpShort is a short polynomial.
// (i.e., all coefficients are bounded by [ShortLogBound] bits.)
// This is faster than [*Evaluator.MulPolyAssign], and the result is exact unlike [*Evaluator.FourierPolyMulPolyAssign].
func (e *Evaluator) ShortFourierPolyMulPolyAssign(p0 Poly, fpShort FourierPoly, pOut Poly) {
	splitBits, splitCount := splitParametersShort(e.Parameters.degree)

	var splitMask uint64 = 1<<splitBits - 1
	for i := 0; i < splitCount; i++ {
		splitLowBits := i * splitBits
		for j := 0; j < e.Parameters.degree; j++ {
			e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] >> splitLowBits) & splitMask
		}
		e.ToFourierPolyAssign(e.buffer.pSplit, e.buffer.fpShortSplit[i])
		e.MulFourierPolyAssign(e.buffer.fpShortSplit[i], fpShort, e.buffer.fpShortSplit[i])
	}

	e.ToPolyAssignUnsafe(e.buffer.fpShortSplit[0], pOut)
	for i := 1; i < splitCount; i++ {
		e.ToPolyAssignUnsafe(e.buffer.fpShortSplit[i], e.buffer.pSplit)
		splitLowBits := i * int(splitBits)
		for j := 0; j < e.Parameters.degree; j++ {
			pOut.Coeffs[j] += e.buffer.pSplit.Coeffs[j] << splitLowBits
		}
	}
}

// ShortFourierPolyMulAddPolyAssign computes pOut += p0 * fpShort, under the assumption that fpShort is a short polynomial.
// (i.e., all coefficients are bounded by [ShortLogBound] bits.)
// This is faster than [*Evaluator.MulAddPolyAssign], and the result is exact unlike [*Evaluator.FourierPolyMulAddPolyAssign].
func (e *Evaluator) ShortFourierPolyMulAddPolyAssign(p0 Poly, fpShort FourierPoly, pOut Poly) {
	splitBits, splitCount := splitParametersShort(e.Parameters.degree)

	var splitMask uint64 = 1<<splitBits - 1
	for i := 0; i < splitCount; i++ {
		splitLowBits := i * splitBits
		for j := 0; j < e.Parameters.degree; j++ {
			e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] >> splitLowBits) & splitMask
		}
		e.ToFourierPolyAssign(e.buffer.pSplit, e.buffer.fpShortSplit[i])
		e.MulFourierPolyAssign(e.buffer.fpShortSplit[i], fpShort, e.buffer.fpShortSplit[i])
	}

	e.ToPolyAddAssignUnsafe(e.buffer.fpShortSplit[0], pOut)
	for i := 1; i < splitCount; i++ {
		e.ToPolyAssignUnsafe(e.buffer.fpShortSplit[i], e.buffer.pSplit)
		splitLowBits := i * splitBits
		for j := 0; j < e.Parameters.degree; j++ {
			pOut.Coeffs[j] += e.buffer.pSplit.Coeffs[j] << splitLowBits
		}
	}
}

// ShortFourierPolyMulSubPolyAssign computes pOut -= p0 * fpShort, under the assumption that fpShort is a short polynomial.
// (i.e., all coefficients are bounded by [ShortLogBound] bits.)
// This is faster than [*Evaluator.MulSubPolyAssign], and the result is exact unlike [*Evaluator.FourierPolyMulSubPolyAssign].
func (e *Evaluator) ShortFourierPolyMulSubPolyAssign(p0 Poly, fpShort FourierPoly, pOut Poly) {
	splitBits, splitCount := splitParametersShort(e.Parameters.degree)

	var splitMask uint64 = 1<<splitBits - 1
	for i := 0; i < splitCount; i++ {
		splitLowBits := i * splitBits
		for j := 0; j < e.Parameters.degree; j++ {
			e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] >> splitLowBits) & splitMask
		}
		e.ToFourierPolyAssign(e.buffer.pSplit, e.buffer.fpShortSplit[i])
		e.MulFourierPolyAssign(e.buffer.fpShortSplit[i], fpShort, e.buffer.fpShortSplit[i])
	}

	e.ToPolySubAssignUnsafe(e.buffer.fpShortSplit[0], pOut)
	for i := 1; i < splitCount; i++ {
		e.ToPolyAssignUnsafe(e.buffer.fpShortSplit[i], e.buffer.pSplit)
		splitLowBits := i * splitBits
		for j := 0; j < e.Parameters.degree; j++ {
			pOut.Coeffs[j] -= e.buffer.pSplit.Coeffs[j] << splitLowBits
		}
	}
}
