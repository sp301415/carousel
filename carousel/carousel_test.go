package carousel_test

import (
	"fmt"
	"testing"

	"github.com/sp301415/carousel/carousel"
	"github.com/sp301415/carousel/math/num"
	"github.com/sp301415/carousel/math/poly"
	"github.com/stretchr/testify/assert"
)

var (
	params = carousel.ParamsSlotsUint3.Compile()
	enc    = carousel.NewEncryptor(params)
	eval   = carousel.NewEvaluator(params, enc.GenEvaluationKeyParallel())

	paramsList = []carousel.ParametersLiteral{
		carousel.ParamsSlotsUint2,
		carousel.ParamsSlotsUint3,
	}
)

func paramsString(p carousel.ParametersLiteral) string {
	switch p.EncodeType {
	case carousel.EncodeTypeCoeffs:
		return fmt.Sprintf("ParamsCoeffsUint%v", num.Log2(p.MessageModulus))
	case carousel.EncodeTypeSlots:
		return fmt.Sprintf("ParamsSlotsUint%v", num.Log2(p.MessageModulus))
	}
	return ""
}

func TestParams(t *testing.T) {
	for _, params := range paramsList {
		t.Run(fmt.Sprintf("Compile/%v", paramsString(params)), func(t *testing.T) {
			assert.NotPanics(t, func() { params.Compile() })
		})
	}
}

func TestEncoder(t *testing.T) {
	messages := []int{1, 2, 3}

	t.Run("Coeffs", func(t *testing.T) {
		pt := enc.EncodeRLWECoeffs(messages)
		assert.Equal(t, messages, enc.DecodeRLWECoeffs(pt)[:len(messages)])
	})

	t.Run("Slots", func(t *testing.T) {
		pt := enc.EncodeRLWESlots(messages)
		assert.Equal(t, messages, enc.DecodeRLWESlots(pt)[:len(messages)])
	})
}

func TestEncryptor(t *testing.T) {
	messages := []int{1, 2, 3}
	gadgetParams := params.KeySwitchParameters()

	t.Run("LWE", func(t *testing.T) {
		for _, m := range messages {
			ct := enc.EncryptLWE(m)
			assert.Equal(t, m, enc.DecryptLWE(ct))
		}
	})

	t.Run("Lev", func(t *testing.T) {
		for _, m := range messages {
			ct := enc.EncryptLev(m, gadgetParams)
			assert.Equal(t, m, enc.DecryptLev(ct))
		}
	})

	t.Run("GSW", func(t *testing.T) {
		for _, m := range messages {
			ct := enc.EncryptGSW(m, gadgetParams)
			assert.Equal(t, m, enc.DecryptGSW(ct))
		}
	})

	t.Run("RLWE", func(t *testing.T) {
		ct := enc.EncryptRLWE(messages)
		assert.Equal(t, messages, enc.DecryptRLWE(ct)[:len(messages)])
	})

	pt := poly.NewPoly(params.PolyDegree())
	for i := 0; i < len(messages); i++ {
		pt.Coeffs[i] = uint64(messages[i])
	}

	t.Run("RLev", func(t *testing.T) {
		ct := enc.EncryptRLevPoly(pt, gadgetParams)
		assert.Equal(t, pt, enc.DecryptRLevPoly(ct))
	})

	t.Run("RGSW", func(t *testing.T) {
		ct := enc.EncryptRGSWPoly(pt, gadgetParams)
		assert.Equal(t, pt, enc.DecryptRGSWPoly(ct))
	})

	t.Run("FourierRLWE", func(t *testing.T) {
		ct := enc.EncryptFourierRLWE(messages)
		assert.Equal(t, messages, enc.DecryptFourierRLWE(ct)[:len(messages)])
	})

	t.Run("FourierRLev", func(t *testing.T) {
		ct := enc.EncryptFourierRLevPoly(pt, gadgetParams)
		assert.Equal(t, pt, enc.DecryptFourierRLevPoly(ct))
	})

	t.Run("FourierRGSW", func(t *testing.T) {
		ct := enc.EncryptFourierRGSWPoly(pt, gadgetParams)
		assert.Equal(t, pt, enc.DecryptFourierRGSWPoly(ct))
	})
}

func TestEvaluator(t *testing.T) {
	messages := []int{1, 2, 3}

	t.Run("KeySwitch", func(t *testing.T) {
		kskGLWEParams := carousel.GadgetParametersLiteral{
			Base:  1 << 12,
			Level: 3,
		}.Compile()

		encOut := carousel.NewEncryptor(params)

		ctLWEIn := enc.EncryptLWE(messages[0])
		ctGLWEIn := enc.EncryptRLWE(messages)

		kskLWE := encOut.GenLWEKeySwitchKey(enc.SecretKey.LWEKey, params.KeySwitchParameters())
		kskRLWE := encOut.GenGLWEKeySwitchKey(enc.SecretKey.RLWEKey, kskGLWEParams)

		ctLWEOut := eval.KeySwitchLWE(ctLWEIn, kskLWE)
		ctGLWEOut := eval.KeySwitchRLWE(ctGLWEIn, kskRLWE)

		assert.Equal(t, messages[0], encOut.DecryptLWE(ctLWEOut))
		assert.Equal(t, messages, encOut.DecryptRLWE(ctGLWEOut)[:len(messages)])
	})

	t.Run("BootstrapOriginalFunc", func(t *testing.T) {
		f := func(x int) int { return 2 * x }

		paramsOriginal := params.Literal().WithBlockSize(1).Compile()
		evalOriginal := carousel.NewEvaluator(paramsOriginal, eval.EvaluationKey)

		for _, m := range messages {
			ct := enc.Encrypt(m)
			ctOut := evalOriginal.BootstrapFunc(ct, f)
			assert.Equal(t, f(m), enc.Decrypt(ctOut))
		}
	})

	t.Run("BootstrapBlockFunc", func(t *testing.T) {
		f := func(x int) int { return 2 * x }

		for _, m := range messages {
			ct := enc.Encrypt(m)
			ctOut := eval.BootstrapFunc(ct, f)
			assert.Equal(t, f(m), enc.Decrypt(ctOut))
		}
	})
}

func BenchmarkProgrammableBootstrap(b *testing.B) {
	for _, paramsLiteral := range paramsList {
		params := paramsLiteral.Compile()
		enc := carousel.NewEncryptor(params)
		eval := carousel.NewEvaluator(params, enc.GenEvaluationKeyParallel())

		ct := enc.Encrypt(0)
		ctOut := ct.Copy()
		lut := eval.GenLookUpTable(func(x int) int { return 2*x + 1 })

		b.Run(paramsString(paramsLiteral), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				eval.BootstrapLUTAssign(ct, lut, ctOut)
			}
		})
	}
}
