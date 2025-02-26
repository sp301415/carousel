//go:build !(amd64 && !purego)

package poly

// rfftInPlace computes the FFT of coeffs in-place.
func rfftInPlace(coeffs []float64, tw []complex128) {
	N := len(coeffs)
	w := 0

	w++
	for j := 0; j < N/4; j += 4 {
		uReal0 := coeffs[j+0]
		uReal1 := coeffs[j+1]
		uReal2 := coeffs[j+2]
		uReal3 := coeffs[j+3]

		uImag0 := coeffs[j+N/4+0]
		uImag1 := coeffs[j+N/4+1]
		uImag2 := coeffs[j+N/4+2]
		uImag3 := coeffs[j+N/4+3]

		vReal0 := coeffs[j+2*N/4+0]
		vReal1 := coeffs[j+2*N/4+1]
		vReal2 := coeffs[j+2*N/4+2]
		vReal3 := coeffs[j+2*N/4+3]

		vImag0 := coeffs[j+3*N/4+0]
		vImag1 := coeffs[j+3*N/4+1]
		vImag2 := coeffs[j+3*N/4+2]
		vImag3 := coeffs[j+3*N/4+3]

		uOutReal0 := uReal0 + vReal0
		uOutReal1 := uReal1 + vReal1
		uOutReal2 := uReal2 + vReal2
		uOutReal3 := uReal3 + vReal3

		uOutImag0 := uImag0 + vImag0
		uOutImag1 := uImag1 + vImag1
		uOutImag2 := uImag2 + vImag2
		uOutImag3 := uImag3 + vImag3

		vOutReal0 := uReal0 - vReal0
		vOutReal1 := uReal1 - vReal1
		vOutReal2 := uReal2 - vReal2
		vOutReal3 := uReal3 - vReal3

		vOutImag0 := -uImag0 + vImag0
		vOutImag1 := -uImag1 + vImag1
		vOutImag2 := -uImag2 + vImag2
		vOutImag3 := -uImag3 + vImag3

		coeffs[j+0] = uOutReal0
		coeffs[j+1] = uOutReal1
		coeffs[j+2] = uOutReal2
		coeffs[j+3] = uOutReal3

		coeffs[j+N/4+0] = uOutImag0
		coeffs[j+N/4+1] = uOutImag1
		coeffs[j+N/4+2] = uOutImag2
		coeffs[j+N/4+3] = uOutImag3

		coeffs[j+2*N/4+0] = vOutReal0
		coeffs[j+2*N/4+1] = vOutReal1
		coeffs[j+2*N/4+2] = vOutReal2
		coeffs[j+2*N/4+3] = vOutReal3

		coeffs[j+3*N/4+0] = vOutImag0
		coeffs[j+3*N/4+1] = vOutImag1
		coeffs[j+3*N/4+2] = vOutImag2
		coeffs[j+3*N/4+3] = vOutImag3
	}

	t := N >> 2
	for m := 2; m <= N/16; m <<= 1 {
		t >>= 1

		w++
		for j := 0; j < t; j += 4 {
			uReal0 := coeffs[j+0]
			uReal1 := coeffs[j+1]
			uReal2 := coeffs[j+2]
			uReal3 := coeffs[j+3]

			uImag0 := coeffs[j+t+0]
			uImag1 := coeffs[j+t+1]
			uImag2 := coeffs[j+t+2]
			uImag3 := coeffs[j+t+3]

			vReal0 := coeffs[j+2*t+0]
			vReal1 := coeffs[j+2*t+1]
			vReal2 := coeffs[j+2*t+2]
			vReal3 := coeffs[j+2*t+3]

			vImag0 := coeffs[j+3*t+0]
			vImag1 := coeffs[j+3*t+1]
			vImag2 := coeffs[j+3*t+2]
			vImag3 := coeffs[j+3*t+3]

			uOutReal0 := uReal0 + vReal0
			uOutReal1 := uReal1 + vReal1
			uOutReal2 := uReal2 + vReal2
			uOutReal3 := uReal3 + vReal3

			uOutImag0 := uImag0 + vImag0
			uOutImag1 := uImag1 + vImag1
			uOutImag2 := uImag2 + vImag2
			uOutImag3 := uImag3 + vImag3

			vOutReal0 := uReal0 - vReal0
			vOutReal1 := uReal1 - vReal1
			vOutReal2 := uReal2 - vReal2
			vOutReal3 := uReal3 - vReal3

			vOutImag0 := -uImag0 + vImag0
			vOutImag1 := -uImag1 + vImag1
			vOutImag2 := -uImag2 + vImag2
			vOutImag3 := -uImag3 + vImag3

			coeffs[j+0] = uOutReal0
			coeffs[j+1] = uOutReal1
			coeffs[j+2] = uOutReal2
			coeffs[j+3] = uOutReal3

			coeffs[j+t+0] = uOutImag0
			coeffs[j+t+1] = uOutImag1
			coeffs[j+t+2] = uOutImag2
			coeffs[j+t+3] = uOutImag3

			coeffs[j+2*t+0] = vOutReal0
			coeffs[j+2*t+1] = vOutReal1
			coeffs[j+2*t+2] = vOutReal2
			coeffs[j+2*t+3] = vOutReal3

			coeffs[j+3*t+0] = vOutImag0
			coeffs[j+3*t+1] = vOutImag1
			coeffs[j+3*t+2] = vOutImag2
			coeffs[j+3*t+3] = vOutImag3
		}

		for i := 1; i < m; i++ {
			wReal := real(tw[w])
			wImag := imag(tw[w])
			w++

			j1 := 4 * i * t
			j2 := j1 + t
			for j := j1; j < j2; j += 4 {
				uReal0 := coeffs[j+0]
				uReal1 := coeffs[j+1]
				uReal2 := coeffs[j+2]
				uReal3 := coeffs[j+3]

				vReal0 := coeffs[j+t]
				vReal1 := coeffs[j+t+1]
				vReal2 := coeffs[j+t+2]
				vReal3 := coeffs[j+t+3]

				uImag0 := coeffs[j+2*t]
				uImag1 := coeffs[j+2*t+1]
				uImag2 := coeffs[j+2*t+2]
				uImag3 := coeffs[j+2*t+3]

				vImag0 := coeffs[j+3*t]
				vImag1 := coeffs[j+3*t+1]
				vImag2 := coeffs[j+3*t+2]
				vImag3 := coeffs[j+3*t+3]

				vTwReal0 := wReal*vReal0 - wImag*vImag0
				vTwReal1 := wReal*vReal1 - wImag*vImag1
				vTwReal2 := wReal*vReal2 - wImag*vImag2
				vTwReal3 := wReal*vReal3 - wImag*vImag3

				vTwImag0 := wReal*vImag0 + wImag*vReal0
				vTwImag1 := wReal*vImag1 + wImag*vReal1
				vTwImag2 := wReal*vImag2 + wImag*vReal2
				vTwImag3 := wReal*vImag3 + wImag*vReal3

				uOutReal0 := uReal0 + vTwReal0
				uOutReal1 := uReal1 + vTwReal1
				uOutReal2 := uReal2 + vTwReal2
				uOutReal3 := uReal3 + vTwReal3

				uOutImag0 := uImag0 + vTwImag0
				uOutImag1 := uImag1 + vTwImag1
				uOutImag2 := uImag2 + vTwImag2
				uOutImag3 := uImag3 + vTwImag3

				vOutReal0 := uReal0 - vTwReal0
				vOutReal1 := uReal1 - vTwReal1
				vOutReal2 := uReal2 - vTwReal2
				vOutReal3 := uReal3 - vTwReal3

				vOutImag0 := -uImag0 + vTwImag0
				vOutImag1 := -uImag1 + vTwImag1
				vOutImag2 := -uImag2 + vTwImag2
				vOutImag3 := -uImag3 + vTwImag3

				coeffs[j+0] = uOutReal0
				coeffs[j+1] = uOutReal1
				coeffs[j+2] = uOutReal2
				coeffs[j+3] = uOutReal3

				coeffs[j+t] = uOutImag0
				coeffs[j+t+1] = uOutImag1
				coeffs[j+t+2] = uOutImag2
				coeffs[j+t+3] = uOutImag3

				coeffs[j+2*t] = vOutReal0
				coeffs[j+2*t+1] = vOutReal1
				coeffs[j+2*t+2] = vOutReal2
				coeffs[j+2*t+3] = vOutReal3

				coeffs[j+3*t] = vOutImag0
				coeffs[j+3*t+1] = vOutImag1
				coeffs[j+3*t+2] = vOutImag2
				coeffs[j+3*t+3] = vOutImag3
			}
		}
	}

	{
		w++

		uReal0 := coeffs[0]
		uReal1 := coeffs[1]

		uImag0 := coeffs[2]
		uImag1 := coeffs[3]

		vReal0 := coeffs[4]
		vReal1 := coeffs[5]

		vImag0 := coeffs[6]
		vImag1 := coeffs[7]

		uOutReal0 := uReal0 + vReal0
		uOutReal1 := uReal1 + vReal1

		uOutImag0 := uImag0 + vImag0
		uOutImag1 := uImag1 + vImag1

		vOutReal0 := uReal0 - vReal0
		vOutReal1 := uReal1 - vReal1

		vOutImag0 := -uImag0 + vImag0
		vOutImag1 := -uImag1 + vImag1

		coeffs[0] = uOutReal0
		coeffs[1] = uOutReal1

		coeffs[2] = uOutImag0
		coeffs[3] = uOutImag1

		coeffs[4] = vOutReal0
		coeffs[5] = vOutReal1

		coeffs[6] = vOutImag0
		coeffs[7] = vOutImag1
	}

	for i := 1; i < N/8; i++ {
		wReal := real(tw[w])
		wImag := imag(tw[w])
		w++

		j := 8 * i

		uReal0 := coeffs[j+0]
		uReal1 := coeffs[j+1]

		vReal0 := coeffs[j+2]
		vReal1 := coeffs[j+3]

		uImag0 := coeffs[j+4]
		uImag1 := coeffs[j+5]

		vImag0 := coeffs[j+6]
		vImag1 := coeffs[j+7]

		vTwReal0 := wReal*vReal0 - wImag*vImag0
		vTwReal1 := wReal*vReal1 - wImag*vImag1

		vTwImag0 := wReal*vImag0 + wImag*vReal0
		vTwImag1 := wReal*vImag1 + wImag*vReal1

		uOutReal0 := uReal0 + vTwReal0
		uOutReal1 := uReal1 + vTwReal1

		uOutImag0 := uImag0 + vTwImag0
		uOutImag1 := uImag1 + vTwImag1

		vOutReal0 := uReal0 - vTwReal0
		vOutReal1 := uReal1 - vTwReal1

		vOutImag0 := -uImag0 + vTwImag0
		vOutImag1 := -uImag1 + vTwImag1

		coeffs[j+0] = uOutReal0
		coeffs[j+1] = uOutReal1

		coeffs[j+2] = uOutImag0
		coeffs[j+3] = uOutImag1

		coeffs[j+4] = vOutReal0
		coeffs[j+5] = vOutReal1

		coeffs[j+6] = vOutImag0
		coeffs[j+7] = vOutImag1
	}

	{
		w++

		uReal := coeffs[0]
		uImag := coeffs[1]

		vReal := coeffs[2]
		vImag := coeffs[3]

		uOutReal := uReal + vReal
		uOutImag := uImag + vImag

		vOutReal := uReal - vReal
		vOutImag := -uImag + vImag

		coeffs[0] = uOutReal
		coeffs[1] = uOutImag

		coeffs[2] = vOutReal
		coeffs[3] = vOutImag
	}

	for i := 1; i < N/4; i++ {
		wReal := real(tw[w])
		wImag := imag(tw[w])
		w++

		j := 4 * i

		uReal := coeffs[j+0]
		vReal := coeffs[j+1]

		uImag := coeffs[j+2]
		vImag := coeffs[j+3]

		vTwReal := wReal*vReal - wImag*vImag
		vTwImag := wReal*vImag + wImag*vReal

		uOutReal := uReal + vTwReal
		uOutImag := uImag + vTwImag

		vOutReal := uReal - vTwReal
		vOutImag := -uImag + vTwImag

		coeffs[j+0] = uOutReal
		coeffs[j+1] = uOutImag

		coeffs[j+2] = vOutReal
		coeffs[j+3] = vOutImag
	}

	coeffs[0], coeffs[1] = coeffs[0]+coeffs[1], coeffs[0]-coeffs[1]
}

// rifftInPlace computes the inverse FFT of coeffs in-place.
func rifftInPlace(coeffs []float64, twInv []complex128) {
	N := len(coeffs)
	w := 0

	coeffs[0], coeffs[1] = 0.5*(coeffs[0]+coeffs[1]), 0.5*(coeffs[0]-coeffs[1])

	{
		w++

		uReal := coeffs[0]
		uImag := coeffs[1]
		vReal := coeffs[2]
		vImag := -coeffs[3]

		uOutReal := uReal + vReal
		uOutImag := uImag + vImag
		vOutReal := uReal - vReal
		vOutImag := uImag - vImag

		coeffs[0] = uOutReal
		coeffs[1] = uOutImag
		coeffs[2] = vOutReal
		coeffs[3] = vOutImag
	}

	for i := 1; i < N/4; i++ {
		wReal := real(twInv[w])
		wImag := imag(twInv[w])
		w++

		j := 4 * i

		uReal := coeffs[j+0]
		uImag := coeffs[j+1]
		vReal := coeffs[j+2]
		vImag := -coeffs[j+3]

		uOutReal := uReal + vReal
		uOutImag := uImag + vImag
		vOutReal := uReal - vReal
		vOutImag := uImag - vImag

		vOutTwReal := wReal*vOutReal - wImag*vOutImag
		vOutTwImag := wReal*vOutImag + wImag*vOutReal

		coeffs[j+0] = uOutReal
		coeffs[j+1] = vOutTwReal
		coeffs[j+2] = uOutImag
		coeffs[j+3] = vOutTwImag
	}

	{
		w++

		uReal0 := coeffs[0]
		uReal1 := coeffs[1]

		uImag0 := coeffs[2]
		uImag1 := coeffs[3]

		vReal0 := coeffs[4]
		vReal1 := coeffs[5]

		vImag0 := -coeffs[6]
		vImag1 := -coeffs[7]

		uOutReal0 := uReal0 + vReal0
		uOutReal1 := uReal1 + vReal1

		uOutImag0 := uImag0 + vImag0
		uOutImag1 := uImag1 + vImag1

		vOutReal0 := uReal0 - vReal0
		vOutReal1 := uReal1 - vReal1

		vOutImag0 := uImag0 - vImag0
		vOutImag1 := uImag1 - vImag1

		coeffs[0] = uOutReal0
		coeffs[1] = uOutReal1

		coeffs[2] = uOutImag0
		coeffs[3] = uOutImag1

		coeffs[4] = vOutReal0
		coeffs[5] = vOutReal1

		coeffs[6] = vOutImag0
		coeffs[7] = vOutImag1
	}

	for i := 1; i < N/8; i++ {
		wReal := real(twInv[w])
		wImag := imag(twInv[w])
		w++

		j := 8 * i

		uReal0 := coeffs[j+0]
		uReal1 := coeffs[j+1]

		uImag0 := coeffs[j+2]
		uImag1 := coeffs[j+3]

		vReal0 := coeffs[j+4]
		vReal1 := coeffs[j+5]

		vImag0 := -coeffs[j+6]
		vImag1 := -coeffs[j+7]

		uOutReal0 := uReal0 + vReal0
		uOutReal1 := uReal1 + vReal1

		uOutImag0 := uImag0 + vImag0
		uOutImag1 := uImag1 + vImag1

		vOutReal0 := uReal0 - vReal0
		vOutReal1 := uReal1 - vReal1

		vOutImag0 := uImag0 - vImag0
		vOutImag1 := uImag1 - vImag1

		vOutTwReal0 := wReal*vOutReal0 - wImag*vOutImag0
		vOutTwReal1 := wReal*vOutReal1 - wImag*vOutImag1

		vOutTwImag0 := wReal*vOutImag0 + wImag*vOutReal0
		vOutTwImag1 := wReal*vOutImag1 + wImag*vOutReal1

		coeffs[j+0] = uOutReal0
		coeffs[j+1] = uOutReal1

		coeffs[j+2] = vOutTwReal0
		coeffs[j+3] = vOutTwReal1

		coeffs[j+4] = uOutImag0
		coeffs[j+5] = uOutImag1

		coeffs[j+6] = vOutTwImag0
		coeffs[j+7] = vOutTwImag1
	}

	t := 4
	for m := N / 16; m >= 2; m >>= 1 {
		w++
		for j := 0; j < t; j += 4 {
			uReal0 := coeffs[j+0]
			uReal1 := coeffs[j+1]
			uReal2 := coeffs[j+2]
			uReal3 := coeffs[j+3]

			uImag0 := coeffs[j+t]
			uImag1 := coeffs[j+t+1]
			uImag2 := coeffs[j+t+2]
			uImag3 := coeffs[j+t+3]

			vReal0 := coeffs[j+2*t]
			vReal1 := coeffs[j+2*t+1]
			vReal2 := coeffs[j+2*t+2]
			vReal3 := coeffs[j+2*t+3]

			vImag0 := -coeffs[j+3*t]
			vImag1 := -coeffs[j+3*t+1]
			vImag2 := -coeffs[j+3*t+2]
			vImag3 := -coeffs[j+3*t+3]

			uOutReal0 := uReal0 + vReal0
			uOutReal1 := uReal1 + vReal1
			uOutReal2 := uReal2 + vReal2
			uOutReal3 := uReal3 + vReal3

			uOutImag0 := uImag0 + vImag0
			uOutImag1 := uImag1 + vImag1
			uOutImag2 := uImag2 + vImag2
			uOutImag3 := uImag3 + vImag3

			vOutReal0 := uReal0 - vReal0
			vOutReal1 := uReal1 - vReal1
			vOutReal2 := uReal2 - vReal2
			vOutReal3 := uReal3 - vReal3

			vOutImag0 := uImag0 - vImag0
			vOutImag1 := uImag1 - vImag1
			vOutImag2 := uImag2 - vImag2
			vOutImag3 := uImag3 - vImag3

			coeffs[j+0] = uOutReal0
			coeffs[j+1] = uOutReal1
			coeffs[j+2] = uOutReal2
			coeffs[j+3] = uOutReal3

			coeffs[j+t] = uOutImag0
			coeffs[j+t+1] = uOutImag1
			coeffs[j+t+2] = uOutImag2
			coeffs[j+t+3] = uOutImag3

			coeffs[j+2*t] = vOutReal0
			coeffs[j+2*t+1] = vOutReal1
			coeffs[j+2*t+2] = vOutReal2
			coeffs[j+2*t+3] = vOutReal3

			coeffs[j+3*t] = vOutImag0
			coeffs[j+3*t+1] = vOutImag1
			coeffs[j+3*t+2] = vOutImag2
			coeffs[j+3*t+3] = vOutImag3
		}

		for i := 1; i < m; i++ {
			wReal := real(twInv[w])
			wImag := imag(twInv[w])
			w++

			j1 := 4 * i * t
			j2 := j1 + t
			for j := j1; j < j2; j += 4 {
				uReal0 := coeffs[j+0]
				uReal1 := coeffs[j+1]
				uReal2 := coeffs[j+2]
				uReal3 := coeffs[j+3]

				uImag0 := coeffs[j+t]
				uImag1 := coeffs[j+t+1]
				uImag2 := coeffs[j+t+2]
				uImag3 := coeffs[j+t+3]

				vReal0 := coeffs[j+2*t]
				vReal1 := coeffs[j+2*t+1]
				vReal2 := coeffs[j+2*t+2]
				vReal3 := coeffs[j+2*t+3]

				vImag0 := -coeffs[j+3*t]
				vImag1 := -coeffs[j+3*t+1]
				vImag2 := -coeffs[j+3*t+2]
				vImag3 := -coeffs[j+3*t+3]

				uOutReal0 := uReal0 + vReal0
				uOutReal1 := uReal1 + vReal1
				uOutReal2 := uReal2 + vReal2
				uOutReal3 := uReal3 + vReal3

				uOutImag0 := uImag0 + vImag0
				uOutImag1 := uImag1 + vImag1
				uOutImag2 := uImag2 + vImag2
				uOutImag3 := uImag3 + vImag3

				vOutReal0 := uReal0 - vReal0
				vOutReal1 := uReal1 - vReal1
				vOutReal2 := uReal2 - vReal2
				vOutReal3 := uReal3 - vReal3

				vOutImag0 := uImag0 - vImag0
				vOutImag1 := uImag1 - vImag1
				vOutImag2 := uImag2 - vImag2
				vOutImag3 := uImag3 - vImag3

				vOutTwReal0 := wReal*vOutReal0 - wImag*vOutImag0
				vOutTwReal1 := wReal*vOutReal1 - wImag*vOutImag1
				vOutTwReal2 := wReal*vOutReal2 - wImag*vOutImag2
				vOutTwReal3 := wReal*vOutReal3 - wImag*vOutImag3

				vOutTwImag0 := wReal*vOutImag0 + wImag*vOutReal0
				vOutTwImag1 := wReal*vOutImag1 + wImag*vOutReal1
				vOutTwImag2 := wReal*vOutImag2 + wImag*vOutReal2
				vOutTwImag3 := wReal*vOutImag3 + wImag*vOutReal3

				coeffs[j+0] = uOutReal0
				coeffs[j+1] = uOutReal1
				coeffs[j+2] = uOutReal2
				coeffs[j+3] = uOutReal3

				coeffs[j+t] = vOutTwReal0
				coeffs[j+t+1] = vOutTwReal1
				coeffs[j+t+2] = vOutTwReal2
				coeffs[j+t+3] = vOutTwReal3

				coeffs[j+2*t] = uOutImag0
				coeffs[j+2*t+1] = uOutImag1
				coeffs[j+2*t+2] = uOutImag2
				coeffs[j+2*t+3] = uOutImag3

				coeffs[j+3*t] = vOutTwImag0
				coeffs[j+3*t+1] = vOutTwImag1
				coeffs[j+3*t+2] = vOutTwImag2
				coeffs[j+3*t+3] = vOutTwImag3
			}
		}
		t <<= 1
	}

	w++
	for j := 0; j < N/4; j += 4 {
		uReal0 := coeffs[j+0]
		uReal1 := coeffs[j+1]
		uReal2 := coeffs[j+2]
		uReal3 := coeffs[j+3]

		uImag0 := coeffs[j+N/4+0]
		uImag1 := coeffs[j+N/4+1]
		uImag2 := coeffs[j+N/4+2]
		uImag3 := coeffs[j+N/4+3]

		vReal0 := coeffs[j+2*N/4+0]
		vReal1 := coeffs[j+2*N/4+1]
		vReal2 := coeffs[j+2*N/4+2]
		vReal3 := coeffs[j+2*N/4+3]

		vImag0 := -coeffs[j+3*N/4+0]
		vImag1 := -coeffs[j+3*N/4+1]
		vImag2 := -coeffs[j+3*N/4+2]
		vImag3 := -coeffs[j+3*N/4+3]

		uOutReal0 := uReal0 + vReal0
		uOutReal1 := uReal1 + vReal1
		uOutReal2 := uReal2 + vReal2
		uOutReal3 := uReal3 + vReal3

		uOutImag0 := uImag0 + vImag0
		uOutImag1 := uImag1 + vImag1
		uOutImag2 := uImag2 + vImag2
		uOutImag3 := uImag3 + vImag3

		vOutReal0 := uReal0 - vReal0
		vOutReal1 := uReal1 - vReal1
		vOutReal2 := uReal2 - vReal2
		vOutReal3 := uReal3 - vReal3

		vOutImag0 := uImag0 - vImag0
		vOutImag1 := uImag1 - vImag1
		vOutImag2 := uImag2 - vImag2
		vOutImag3 := uImag3 - vImag3

		coeffs[j+0] = uOutReal0
		coeffs[j+1] = uOutReal1
		coeffs[j+2] = uOutReal2
		coeffs[j+3] = uOutReal3

		coeffs[j+N/4+0] = uOutImag0
		coeffs[j+N/4+1] = uOutImag1
		coeffs[j+N/4+2] = uOutImag2
		coeffs[j+N/4+3] = uOutImag3

		coeffs[j+2*N/4+0] = vOutReal0
		coeffs[j+2*N/4+1] = vOutReal1
		coeffs[j+2*N/4+2] = vOutReal2
		coeffs[j+2*N/4+3] = vOutReal3

		coeffs[j+3*N/4+0] = vOutImag0
		coeffs[j+3*N/4+1] = vOutImag1
		coeffs[j+3*N/4+2] = vOutImag2
		coeffs[j+3*N/4+3] = vOutImag3
	}
}

// convolveAssign computes the convolution of fourier transformed coefficients fp0, fp1
// and writes it to fpOut.
func convolveAssign(fp0, fp1, fpOut []float64) {
	N := len(fpOut)

	fpOut[0] = fp0[0] * fp1[0]
	fpOut[1] = fp0[1] * fp1[1]

	fp0Real := fp0[2]
	fp0Imag := fp0[3]

	fp1Real := fp1[2]
	fp1Imag := fp1[3]

	fpOutReal := fp0Real*fp1Real - fp0Imag*fp1Imag
	fpOutImag := fp0Real*fp1Imag + fp0Imag*fp1Real

	fpOut[2] = fpOutReal
	fpOut[3] = fpOutImag

	for i := 4; i < N; i += 4 {
		fp0Real0 := fp0[i+0]
		fp0Imag0 := fp0[i+1]
		fp0Real1 := fp0[i+2]
		fp0Imag1 := fp0[i+3]

		fp1Real0 := fp1[i+0]
		fp1Imag0 := fp1[i+1]
		fp1Real1 := fp1[i+2]
		fp1Imag1 := fp1[i+3]

		fpOutReal0 := fp0Real0*fp1Real0 - fp0Imag0*fp1Imag0
		fpOutImag0 := fp0Real0*fp1Imag0 + fp0Imag0*fp1Real0
		fpOutReal1 := fp0Real1*fp1Real1 - fp0Imag1*fp1Imag1
		fpOutImag1 := fp0Real1*fp1Imag1 + fp0Imag1*fp1Real1

		fpOut[i+0] = fpOutReal0
		fpOut[i+1] = fpOutImag0
		fpOut[i+2] = fpOutReal1
		fpOut[i+3] = fpOutImag1
	}
}
