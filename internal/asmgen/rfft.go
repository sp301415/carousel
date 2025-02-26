package main

import (
	"math"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func rfftInPlaceAVX2() {
	TEXT("rfftInPlaceAVX2", NOSPLIT, "func(coeffs []float64, tw []complex128)")
	Pragma("noescape")

	coeffs := Load(Param("coeffs").Base(), GP64())
	tw := Load(Param("tw").Base(), GP64())

	N := Load(Param("coeffs").Len(), GP64())
	w := GP64()
	XORQ(w, w)

	NDiv16 := GP64()
	MOVQ(N, NDiv16)
	SHRQ(Imm(4), NDiv16)

	NDiv8 := GP64()
	MOVQ(N, NDiv8)
	SHRQ(Imm(3), NDiv8)

	NDiv4 := GP64()
	MOVQ(N, NDiv4)
	SHRQ(Imm(2), NDiv4)

	ADDQ(Imm(2), w)

	j := GP64()
	XORQ(j, j)
	jt := GP64()
	MOVQ(j, jt)
	ADDQ(NDiv4, jt)
	j2t := GP64()
	MOVQ(jt, j2t)
	ADDQ(NDiv4, j2t)
	j3t := GP64()
	MOVQ(j2t, j3t)
	ADDQ(NDiv4, j3t)
	JMP(LabelRef("first_loop_end"))
	Label("first_loop_body")

	uReal, uImag, vReal, vImag := YMM(), YMM(), YMM(), YMM()
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8}, uReal)
	VMOVUPD(Mem{Base: coeffs, Index: jt, Scale: 8}, uImag)
	VMOVUPD(Mem{Base: coeffs, Index: j2t, Scale: 8}, vReal)
	VMOVUPD(Mem{Base: coeffs, Index: j3t, Scale: 8}, vImag)

	uOutReal := YMM()
	VADDPD(vReal, uReal, uOutReal)
	uOutImag := YMM()
	VADDPD(vImag, uImag, uOutImag)

	vOutReal := YMM()
	VSUBPD(vReal, uReal, vOutReal)
	vOutImag := YMM()
	VSUBPD(uImag, vImag, vOutImag)

	VMOVUPD(uOutReal, Mem{Base: coeffs, Index: j, Scale: 8})
	VMOVUPD(uOutImag, Mem{Base: coeffs, Index: jt, Scale: 8})
	VMOVUPD(vOutReal, Mem{Base: coeffs, Index: j2t, Scale: 8})
	VMOVUPD(vOutImag, Mem{Base: coeffs, Index: j3t, Scale: 8})

	ADDQ(Imm(4), j)
	ADDQ(Imm(4), jt)
	ADDQ(Imm(4), j2t)
	ADDQ(Imm(4), j3t)

	Label("first_loop_end")
	CMPQ(j, NDiv4)
	JL(LabelRef("first_loop_body"))

	t := GP64()
	MOVQ(N, t)
	SHRQ(Imm(2), t)

	m := GP64()
	MOVQ(U64(2), m)
	JMP(LabelRef("m_loop_end"))
	Label("m_loop_body")

	SHRQ(Imm(1), t)

	ADDQ(Imm(2), w)

	j = GP64()
	XORQ(j, j)
	jt = GP64()
	MOVQ(j, jt)
	ADDQ(t, jt)
	j2t = GP64()
	MOVQ(jt, j2t)
	ADDQ(t, j2t)
	j3t = GP64()
	MOVQ(j2t, j3t)
	ADDQ(t, j3t)
	JMP(LabelRef("j_first_loop_end"))
	Label("j_first_loop_body")

	uReal, uImag, vReal, vImag = YMM(), YMM(), YMM(), YMM()
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8}, uReal)
	VMOVUPD(Mem{Base: coeffs, Index: jt, Scale: 8}, uImag)
	VMOVUPD(Mem{Base: coeffs, Index: j2t, Scale: 8}, vReal)
	VMOVUPD(Mem{Base: coeffs, Index: j3t, Scale: 8}, vImag)

	uOutReal = YMM()
	VADDPD(vReal, uReal, uOutReal)
	uOutImag = YMM()
	VADDPD(vImag, uImag, uOutImag)

	vOutReal = YMM()
	VSUBPD(vReal, uReal, vOutReal)
	vOutImag = YMM()
	VSUBPD(uImag, vImag, vOutImag)

	VMOVUPD(uOutReal, Mem{Base: coeffs, Index: j, Scale: 8})
	VMOVUPD(uOutImag, Mem{Base: coeffs, Index: jt, Scale: 8})
	VMOVUPD(vOutReal, Mem{Base: coeffs, Index: j2t, Scale: 8})
	VMOVUPD(vOutImag, Mem{Base: coeffs, Index: j3t, Scale: 8})

	ADDQ(Imm(4), j)
	ADDQ(Imm(4), jt)
	ADDQ(Imm(4), j2t)
	ADDQ(Imm(4), j3t)

	Label("j_first_loop_end")
	CMPQ(j, t)
	JL(LabelRef("j_first_loop_body"))

	i := GP64()
	MOVQ(U64(1), i)
	JMP(LabelRef("i_loop_end"))
	Label("i_loop_body")

	wReal, wImag := YMM(), YMM()
	VBROADCASTSD(Mem{Base: tw, Index: w, Scale: 8}, wReal)
	VBROADCASTSD(Mem{Base: tw, Index: w, Scale: 8, Disp: 8}, wImag)
	ADDQ(Imm(2), w)

	j1 := GP64()
	MOVQ(t, j1)
	IMULQ(i, j1)
	SHLQ(Imm(2), j1)

	j2 := GP64()
	MOVQ(j1, j2)
	ADDQ(t, j2)

	j = GP64()
	MOVQ(j1, j)
	jt = GP64()
	MOVQ(j, jt)
	ADDQ(t, jt)
	j2t = GP64()
	MOVQ(jt, j2t)
	ADDQ(t, j2t)
	j3t = GP64()
	MOVQ(j2t, j3t)
	ADDQ(t, j3t)
	JMP(LabelRef("j_loop_end"))
	Label("j_loop_body")

	uReal, uImag, vReal, vImag = YMM(), YMM(), YMM(), YMM()
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8}, uReal)
	VMOVUPD(Mem{Base: coeffs, Index: jt, Scale: 8}, vReal)
	VMOVUPD(Mem{Base: coeffs, Index: j2t, Scale: 8}, uImag)
	VMOVUPD(Mem{Base: coeffs, Index: j3t, Scale: 8}, vImag)

	vTwReal := YMM()
	VMULPD(vReal, wReal, vTwReal)
	VFNMADD231PD(vImag, wImag, vTwReal)

	vTwImag := YMM()
	VMULPD(vImag, wReal, vTwImag)
	VFMADD231PD(vReal, wImag, vTwImag)

	uOutReal = YMM()
	VADDPD(vTwReal, uReal, uOutReal)
	uOutImag = YMM()
	VADDPD(vTwImag, uImag, uOutImag)

	vOutReal = YMM()
	VSUBPD(vTwReal, uReal, vOutReal)
	vOutImag = YMM()
	VSUBPD(uImag, vTwImag, vOutImag)

	VMOVUPD(uOutReal, Mem{Base: coeffs, Index: j, Scale: 8})
	VMOVUPD(uOutImag, Mem{Base: coeffs, Index: jt, Scale: 8})
	VMOVUPD(vOutReal, Mem{Base: coeffs, Index: j2t, Scale: 8})
	VMOVUPD(vOutImag, Mem{Base: coeffs, Index: j3t, Scale: 8})

	ADDQ(Imm(4), j)
	ADDQ(Imm(4), jt)
	ADDQ(Imm(4), j2t)
	ADDQ(Imm(4), j3t)

	Label("j_loop_end")
	CMPQ(j, j2)
	JL(LabelRef("j_loop_body"))

	ADDQ(Imm(1), i)

	Label("i_loop_end")
	CMPQ(i, m)
	JL(LabelRef("i_loop_body"))

	SHLQ(Imm(1), m)

	Label("m_loop_end")
	CMPQ(m, NDiv16)
	JLE(LabelRef("m_loop_body"))

	zero := XMM()
	VPXOR(zero, zero, zero)

	ADDQ(Imm(2), w)

	uReal, uImag, vReal, vImag = XMM(), XMM(), XMM(), XMM()
	VMOVUPD(Mem{Base: coeffs, Disp: 0, Scale: 8}, uReal)
	VMOVUPD(Mem{Base: coeffs, Disp: 16, Scale: 8}, uImag)
	VMOVUPD(Mem{Base: coeffs, Disp: 32, Scale: 8}, vReal)
	VMOVUPD(Mem{Base: coeffs, Disp: 48, Scale: 8}, vImag)

	uOutReal = XMM()
	VADDPD(vReal, uReal, uOutReal)
	uOutImag = XMM()
	VADDPD(vImag, uImag, uOutImag)

	vOutReal = XMM()
	VSUBPD(vReal, uReal, vOutReal)
	vOutImag = XMM()
	VSUBPD(uImag, vImag, vOutImag)

	VMOVUPD(uOutReal, Mem{Base: coeffs, Disp: 0, Scale: 8})
	VMOVUPD(uOutImag, Mem{Base: coeffs, Disp: 16, Scale: 8})
	VMOVUPD(vOutReal, Mem{Base: coeffs, Disp: 32, Scale: 8})
	VMOVUPD(vOutImag, Mem{Base: coeffs, Disp: 48, Scale: 8})

	i = GP64()
	MOVQ(U64(1), i)
	JMP(LabelRef("last_loop_2_end"))
	Label("last_loop_2_body")

	wReal, wImag = XMM(), XMM()
	VMOVUPD(Mem{Base: tw, Index: w, Scale: 8}, wReal)
	VSHUFPD(Imm(0b11), wReal, wReal, wImag)
	VSHUFPD(Imm(0b00), wReal, wReal, wReal)
	ADDQ(Imm(2), w)

	j = GP64()
	MOVQ(i, j)
	SHLQ(Imm(3), j)

	uReal, uImag, vReal, vImag = XMM(), XMM(), XMM(), XMM()
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8}, uReal)
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8, Disp: 16}, vReal)
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8, Disp: 32}, uImag)
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8, Disp: 48}, vImag)

	vTwReal = XMM()
	VMULPD(vReal, wReal, vTwReal)
	VFNMADD231PD(vImag, wImag, vTwReal)

	vTwImag = XMM()
	VMULPD(vImag, wReal, vTwImag)
	VFMADD231PD(vReal, wImag, vTwImag)

	uOutReal = XMM()
	VADDPD(vTwReal, uReal, uOutReal)
	uOutImag = XMM()
	VADDPD(vTwImag, uImag, uOutImag)

	vOutReal = XMM()
	VSUBPD(vTwReal, uReal, vOutReal)
	vOutImag = XMM()
	VSUBPD(uImag, vTwImag, vOutImag)

	VMOVUPD(uOutReal, Mem{Base: coeffs, Index: j, Scale: 8})
	VMOVUPD(uOutImag, Mem{Base: coeffs, Index: j, Scale: 8, Disp: 16})
	VMOVUPD(vOutReal, Mem{Base: coeffs, Index: j, Scale: 8, Disp: 32})
	VMOVUPD(vOutImag, Mem{Base: coeffs, Index: j, Scale: 8, Disp: 48})

	ADDQ(Imm(1), i)

	Label("last_loop_2_end")
	CMPQ(i, NDiv8)
	JL(LabelRef("last_loop_2_body"))

	ADDQ(Imm(2), w)

	uRealImag, vRealImag := XMM(), XMM()
	VMOVUPD(Mem{Base: coeffs, Disp: 0, Scale: 8}, uRealImag)
	VMOVUPD(Mem{Base: coeffs, Disp: 16, Scale: 8}, vRealImag)

	uOutRealImag := XMM()
	VADDPD(vRealImag, uRealImag, uOutRealImag)
	vOutRealImag := XMM()
	VSUBPD(vRealImag, uRealImag, vOutRealImag)

	VADDSUBPD(vOutRealImag, zero, vOutRealImag)
	VSUBPD(vOutRealImag, zero, vOutRealImag)

	VMOVUPD(uOutRealImag, Mem{Base: coeffs, Disp: 0, Scale: 8})
	VMOVUPD(vOutRealImag, Mem{Base: coeffs, Disp: 16, Scale: 8})

	i = GP64()
	MOVQ(U64(1), i)
	JMP(LabelRef("last_loop_1_end"))
	Label("last_loop_1_body")

	wReal, wImag = XMM(), XMM()
	VMOVUPD(Mem{Base: tw, Index: w, Scale: 8}, wReal)
	VSHUFPD(Imm(0b11), wReal, wReal, wImag)
	VSHUFPD(Imm(0b00), wReal, wReal, wReal)
	ADDQ(Imm(2), w)

	j = GP64()
	MOVQ(i, j)
	SHLQ(Imm(2), j)

	uRealvReal, uImagvImag := XMM(), XMM()
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8}, uRealvReal)
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8, Disp: 16}, uImagvImag)

	uRealImag = XMM()
	VSHUFPD(Imm(0b00), uImagvImag, uRealvReal, uRealImag)

	vRealImag, vImagReal := XMM(), XMM()
	VSHUFPD(Imm(0b11), uImagvImag, uRealvReal, vRealImag)
	VSHUFPD(Imm(0b11), uRealvReal, uImagvImag, vImagReal)

	vTwRealImag := XMM()
	VMULPD(vImagReal, wImag, vTwRealImag)
	VFMADDSUB231PD(vRealImag, wReal, vTwRealImag)

	uOutRealImag = XMM()
	VADDPD(vTwRealImag, uRealImag, uOutRealImag)
	vOutRealImag = XMM()
	VSUBPD(vTwRealImag, uRealImag, vOutRealImag)

	VADDSUBPD(vOutRealImag, zero, vOutRealImag)
	VSUBPD(vOutRealImag, zero, vOutRealImag)

	VMOVUPD(uOutRealImag, Mem{Base: coeffs, Index: j, Scale: 8})
	VMOVUPD(vOutRealImag, Mem{Base: coeffs, Index: j, Scale: 8, Disp: 16})

	ADDQ(Imm(1), i)

	Label("last_loop_1_end")
	CMPQ(i, NDiv4)
	JL(LabelRef("last_loop_1_body"))

	c0, c1 := XMM(), XMM()
	MOVQ(Mem{Base: coeffs, Disp: 0, Scale: 8}, c0)
	MOVQ(Mem{Base: coeffs, Disp: 8, Scale: 8}, c1)

	cAdd := XMM()
	MOVUPD(c0, cAdd)
	ADDPD(c1, cAdd)

	cSub := XMM()
	MOVUPD(c0, cSub)
	SUBPD(c1, cSub)

	MOVQ(cAdd, Mem{Base: coeffs, Disp: 0, Scale: 8})
	MOVQ(cSub, Mem{Base: coeffs, Disp: 8, Scale: 8})

	RET()
}

func rifftInPlaceAVX2() {
	TEXT("rifftInPlaceAVX2", NOSPLIT, "func(coeffs []float64, tw []complex128)")
	Pragma("noescape")

	coeffs := Load(Param("coeffs").Base(), GP64())
	tw := Load(Param("tw").Base(), GP64())

	N := Load(Param("coeffs").Len(), GP64())
	w := GP64()
	XORQ(w, w)

	NDiv16 := GP64()
	MOVQ(N, NDiv16)
	SHRQ(Imm(4), NDiv16)

	NDiv8 := GP64()
	MOVQ(N, NDiv8)
	SHRQ(Imm(3), NDiv8)

	NDiv4 := GP64()
	MOVQ(N, NDiv4)
	SHRQ(Imm(2), NDiv4)

	half64, half := GP64(), XMM()
	MOVQ(U64(math.Float64bits(0.5)), half64)
	MOVQ(half64, half)

	c0, c1 := XMM(), XMM()
	MOVQ(Mem{Base: coeffs, Disp: 0, Scale: 8}, c0)
	MOVQ(Mem{Base: coeffs, Disp: 8, Scale: 8}, c1)

	cAdd := XMM()
	MOVUPD(c0, cAdd)
	ADDPD(c1, cAdd)
	MULPD(half, cAdd)

	cSub := XMM()
	MOVUPD(c0, cSub)
	SUBPD(c1, cSub)
	MULPD(half, cSub)

	MOVQ(cAdd, Mem{Base: coeffs, Disp: 0, Scale: 8})
	MOVQ(cSub, Mem{Base: coeffs, Disp: 8, Scale: 8})

	zero := XMM()
	VPXOR(zero, zero, zero)

	ADDQ(Imm(2), w)

	uRealImag, vRealImag := XMM(), XMM()
	VMOVUPD(Mem{Base: coeffs, Disp: 0, Scale: 8}, uRealImag)
	VMOVUPD(Mem{Base: coeffs, Disp: 16, Scale: 8}, vRealImag)

	VADDSUBPD(vRealImag, zero, vRealImag)
	VSUBPD(vRealImag, zero, vRealImag)

	uOutRealImag := XMM()
	VADDPD(vRealImag, uRealImag, uOutRealImag)
	vOutRealImag := XMM()
	VSUBPD(vRealImag, uRealImag, vOutRealImag)

	VMOVUPD(uOutRealImag, Mem{Base: coeffs, Disp: 0, Scale: 8})
	VMOVUPD(vOutRealImag, Mem{Base: coeffs, Disp: 16, Scale: 8})

	i := GP64()
	MOVQ(U64(1), i)
	JMP(LabelRef("first_loop_end"))
	Label("first_loop_body")

	wReal, wImag := XMM(), XMM()
	VMOVUPD(Mem{Base: tw, Index: w, Scale: 8}, wReal)
	VSHUFPD(Imm(0b11), wReal, wReal, wImag)
	VSHUFPD(Imm(0b00), wReal, wReal, wReal)
	ADDQ(Imm(2), w)

	j := GP64()
	MOVQ(i, j)
	SHLQ(Imm(2), j)

	uRealImag, vRealImag = XMM(), XMM()
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8}, uRealImag)
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8, Disp: 16}, vRealImag)

	VADDSUBPD(vRealImag, zero, vRealImag)
	VSUBPD(vRealImag, zero, vRealImag)

	uOutRealImag = XMM()
	VADDPD(vRealImag, uRealImag, uOutRealImag)
	vOutRealImag = XMM()
	VSUBPD(vRealImag, uRealImag, vOutRealImag)

	vOutImagReal := XMM()
	VSHUFPD(Imm(0b01), vOutRealImag, vOutRealImag, vOutImagReal)

	vOutTwRealImag := XMM()
	VMULPD(vOutImagReal, wImag, vOutTwRealImag)
	VFMADDSUB231PD(vOutRealImag, wReal, vOutTwRealImag)

	uOutvOutTwReal, uOutvOutTwImag := XMM(), XMM()
	VSHUFPD(Imm(0b00), vOutTwRealImag, uOutRealImag, uOutvOutTwReal)
	VSHUFPD(Imm(0b11), vOutTwRealImag, uOutRealImag, uOutvOutTwImag)

	VMOVUPD(uOutvOutTwReal, Mem{Base: coeffs, Index: j, Scale: 8})
	VMOVUPD(uOutvOutTwImag, Mem{Base: coeffs, Index: j, Scale: 8, Disp: 16})

	ADDQ(Imm(1), i)

	Label("first_loop_end")
	CMPQ(i, NDiv4)
	JL(LabelRef("first_loop_body"))

	ADDQ(Imm(2), w)

	uReal, uImag, vReal, vImag := XMM(), XMM(), XMM(), XMM()
	VMOVUPD(Mem{Base: coeffs, Disp: 0, Scale: 8}, uReal)
	VMOVUPD(Mem{Base: coeffs, Disp: 16, Scale: 8}, uImag)
	VMOVUPD(Mem{Base: coeffs, Disp: 32, Scale: 8}, vReal)
	VMOVUPD(Mem{Base: coeffs, Disp: 48, Scale: 8}, vImag)
	VSUBPD(vImag, zero, vImag)

	uOutReal := XMM()
	VADDPD(vReal, uReal, uOutReal)
	uOutImag := XMM()
	VADDPD(vImag, uImag, uOutImag)

	vOutReal := XMM()
	VSUBPD(vReal, uReal, vOutReal)
	vOutImag := XMM()
	VSUBPD(vImag, uImag, vOutImag)

	VMOVUPD(uOutReal, Mem{Base: coeffs, Disp: 0, Scale: 8})
	VMOVUPD(uOutImag, Mem{Base: coeffs, Disp: 16, Scale: 8})
	VMOVUPD(vOutReal, Mem{Base: coeffs, Disp: 32, Scale: 8})
	VMOVUPD(vOutImag, Mem{Base: coeffs, Disp: 48, Scale: 8})

	i = GP64()
	MOVQ(U64(1), i)
	JMP(LabelRef("first_loop_2_end"))
	Label("first_loop_2_body")

	wReal, wImag = XMM(), XMM()
	VMOVUPD(Mem{Base: tw, Index: w, Scale: 8}, wReal)
	VSHUFPD(Imm(0b11), wReal, wReal, wImag)
	VSHUFPD(Imm(0b00), wReal, wReal, wReal)
	ADDQ(Imm(2), w)

	j = GP64()
	MOVQ(i, j)
	SHLQ(Imm(3), j)

	uReal, uImag, vReal, vImag = XMM(), XMM(), XMM(), XMM()
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8}, uReal)
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8, Disp: 16}, uImag)
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8, Disp: 32}, vReal)
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8, Disp: 48}, vImag)
	VSUBPD(vImag, zero, vImag)

	uOutReal = XMM()
	VADDPD(vReal, uReal, uOutReal)
	uOutImag = XMM()
	VADDPD(vImag, uImag, uOutImag)

	vOutReal = XMM()
	VSUBPD(vReal, uReal, vOutReal)
	vOutImag = XMM()
	VSUBPD(vImag, uImag, vOutImag)

	vOutTwReal := XMM()
	VMULPD(vOutReal, wReal, vOutTwReal)
	VFNMADD231PD(vOutImag, wImag, vOutTwReal)
	vOutTwImag := XMM()
	VMULPD(vOutImag, wReal, vOutTwImag)
	VFMADD231PD(vOutReal, wImag, vOutTwImag)

	VMOVUPD(uOutReal, Mem{Base: coeffs, Index: j, Scale: 8})
	VMOVUPD(vOutTwReal, Mem{Base: coeffs, Index: j, Scale: 8, Disp: 16})
	VMOVUPD(uOutImag, Mem{Base: coeffs, Index: j, Scale: 8, Disp: 32})
	VMOVUPD(vOutTwImag, Mem{Base: coeffs, Index: j, Scale: 8, Disp: 48})

	ADDQ(Imm(1), i)

	Label("first_loop_2_end")
	CMPQ(i, NDiv8)
	JL(LabelRef("first_loop_2_body"))

	zero = YMM()
	VPXOR(zero, zero, zero)

	t := GP64()
	MOVQ(U64(4), t)

	m := GP64()
	MOVQ(NDiv16, m)
	JMP(LabelRef("m_loop_end"))
	Label("m_loop_body")

	ADDQ(Imm(2), w)

	j = GP64()
	XORQ(j, j)
	jt := GP64()
	MOVQ(j, jt)
	ADDQ(t, jt)
	j2t := GP64()
	MOVQ(jt, j2t)
	ADDQ(t, j2t)
	j3t := GP64()
	MOVQ(j2t, j3t)
	ADDQ(t, j3t)
	JMP(LabelRef("j_first_loop_end"))
	Label("j_first_loop_body")

	uReal, uImag, vReal, vImag = YMM(), YMM(), YMM(), YMM()
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8}, uReal)
	VMOVUPD(Mem{Base: coeffs, Index: jt, Scale: 8}, uImag)
	VMOVUPD(Mem{Base: coeffs, Index: j2t, Scale: 8}, vReal)
	VMOVUPD(Mem{Base: coeffs, Index: j3t, Scale: 8}, vImag)
	VSUBPD(vImag, zero, vImag)

	uOutReal = YMM()
	VADDPD(vReal, uReal, uOutReal)
	uOutImag = YMM()
	VADDPD(vImag, uImag, uOutImag)

	vOutReal = YMM()
	VSUBPD(vReal, uReal, vOutReal)
	vOutImag = YMM()
	VSUBPD(vImag, uImag, vOutImag)

	VMOVUPD(uOutReal, Mem{Base: coeffs, Index: j, Scale: 8})
	VMOVUPD(uOutImag, Mem{Base: coeffs, Index: jt, Scale: 8})
	VMOVUPD(vOutReal, Mem{Base: coeffs, Index: j2t, Scale: 8})
	VMOVUPD(vOutImag, Mem{Base: coeffs, Index: j3t, Scale: 8})

	ADDQ(Imm(4), j)
	ADDQ(Imm(4), jt)
	ADDQ(Imm(4), j2t)
	ADDQ(Imm(4), j3t)

	Label("j_first_loop_end")
	CMPQ(j, t)
	JL(LabelRef("j_first_loop_body"))

	i = GP64()
	MOVQ(U64(1), i)
	JMP(LabelRef("i_loop_end"))
	Label("i_loop_body")

	wReal, wImag = YMM(), YMM()
	VBROADCASTSD(Mem{Base: tw, Index: w, Scale: 8}, wReal)
	VBROADCASTSD(Mem{Base: tw, Index: w, Scale: 8, Disp: 8}, wImag)
	ADDQ(Imm(2), w)

	j1 := GP64()
	MOVQ(t, j1)
	IMULQ(i, j1)
	SHLQ(Imm(2), j1)

	j2 := GP64()
	MOVQ(j1, j2)
	ADDQ(t, j2)

	j = GP64()
	MOVQ(j1, j)
	jt = GP64()
	MOVQ(j, jt)
	ADDQ(t, jt)
	j2t = GP64()
	MOVQ(jt, j2t)
	ADDQ(t, j2t)
	j3t = GP64()
	MOVQ(j2t, j3t)
	ADDQ(t, j3t)
	JMP(LabelRef("j_loop_end"))
	Label("j_loop_body")

	uReal, uImag, vReal, vImag = YMM(), YMM(), YMM(), YMM()
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8}, uReal)
	VMOVUPD(Mem{Base: coeffs, Index: jt, Scale: 8}, uImag)
	VMOVUPD(Mem{Base: coeffs, Index: j2t, Scale: 8}, vReal)
	VMOVUPD(Mem{Base: coeffs, Index: j3t, Scale: 8}, vImag)
	VSUBPD(vImag, zero, vImag)

	uOutReal = YMM()
	VADDPD(vReal, uReal, uOutReal)
	uOutImag = YMM()
	VADDPD(vImag, uImag, uOutImag)

	vOutReal = YMM()
	VSUBPD(vReal, uReal, vOutReal)
	vOutImag = YMM()
	VSUBPD(vImag, uImag, vOutImag)

	vOutTwReal = YMM()
	VMULPD(vOutReal, wReal, vOutTwReal)
	VFNMADD231PD(vOutImag, wImag, vOutTwReal)

	vOutTwImag = YMM()
	VMULPD(vOutImag, wReal, vOutTwImag)
	VFMADD231PD(vOutReal, wImag, vOutTwImag)

	VMOVUPD(uOutReal, Mem{Base: coeffs, Index: j, Scale: 8})
	VMOVUPD(vOutTwReal, Mem{Base: coeffs, Index: jt, Scale: 8})
	VMOVUPD(uOutImag, Mem{Base: coeffs, Index: j2t, Scale: 8})
	VMOVUPD(vOutTwImag, Mem{Base: coeffs, Index: j3t, Scale: 8})

	ADDQ(Imm(4), j)
	ADDQ(Imm(4), jt)
	ADDQ(Imm(4), j2t)
	ADDQ(Imm(4), j3t)

	Label("j_loop_end")
	CMPQ(j, j2)
	JL(LabelRef("j_loop_body"))

	ADDQ(Imm(1), i)

	Label("i_loop_end")
	CMPQ(i, m)
	JL(LabelRef("i_loop_body"))

	SHRQ(Imm(1), m)
	SHLQ(Imm(1), t)

	Label("m_loop_end")
	CMPQ(m, Imm(2))
	JGE(LabelRef("m_loop_body"))

	ADDQ(Imm(2), w)

	j = GP64()
	XORQ(j, j)
	jt = GP64()
	MOVQ(j, jt)
	ADDQ(NDiv4, jt)
	j2t = GP64()
	MOVQ(jt, j2t)
	ADDQ(NDiv4, j2t)
	j3t = GP64()
	MOVQ(j2t, j3t)
	ADDQ(NDiv4, j3t)
	JMP(LabelRef("last_loop_end"))
	Label("last_loop_body")

	uReal, uImag, vReal, vImag = YMM(), YMM(), YMM(), YMM()
	VMOVUPD(Mem{Base: coeffs, Index: j, Scale: 8}, uReal)
	VMOVUPD(Mem{Base: coeffs, Index: jt, Scale: 8}, uImag)
	VMOVUPD(Mem{Base: coeffs, Index: j2t, Scale: 8}, vReal)
	VMOVUPD(Mem{Base: coeffs, Index: j3t, Scale: 8}, vImag)
	VSUBPD(vImag, zero, vImag)

	uOutReal = YMM()
	VADDPD(vReal, uReal, uOutReal)
	uOutImag = YMM()
	VADDPD(vImag, uImag, uOutImag)

	vOutReal = YMM()
	VSUBPD(vReal, uReal, vOutReal)
	vOutImag = YMM()
	VSUBPD(vImag, uImag, vOutImag)

	VMOVUPD(uOutReal, Mem{Base: coeffs, Index: j, Scale: 8})
	VMOVUPD(uOutImag, Mem{Base: coeffs, Index: jt, Scale: 8})
	VMOVUPD(vOutReal, Mem{Base: coeffs, Index: j2t, Scale: 8})
	VMOVUPD(vOutImag, Mem{Base: coeffs, Index: j3t, Scale: 8})

	ADDQ(Imm(4), j)
	ADDQ(Imm(4), jt)
	ADDQ(Imm(4), j2t)
	ADDQ(Imm(4), j3t)

	Label("last_loop_end")
	CMPQ(j, NDiv4)
	JL(LabelRef("last_loop_body"))

	RET()
}

func convolveAssignAVX2() {
	TEXT("convolveAssignAVX2", NOSPLIT, "func(fp0, fp1, fpOut []float64)")
	Pragma("noescape")

	fp0 := Load(Param("fp0").Base(), GP64())
	fp1 := Load(Param("fp1").Base(), GP64())
	fpOut := Load(Param("fpOut").Base(), GP64())
	N := Load(Param("fpOut").Len(), GP64())

	fp0Real, fp1Real, fpOutReal := XMM(), XMM(), XMM()
	VMOVUPD(Mem{Base: fp0, Disp: 0, Scale: 8}, fp0Real)
	VMOVUPD(Mem{Base: fp1, Disp: 0, Scale: 8}, fp1Real)
	VMULPD(fp0Real, fp1Real, fpOutReal)
	VMOVUPD(fpOutReal, Mem{Base: fpOut, Disp: 0, Scale: 8})

	fp0RealImag, fp1RealImag := XMM(), XMM()
	VMOVUPD(Mem{Base: fp0, Disp: 16, Scale: 8}, fp0RealImag)
	VMOVUPD(Mem{Base: fp1, Disp: 16, Scale: 8}, fp1RealImag)

	fp0Real, fp0Imag := XMM(), XMM()
	VSHUFPD(Imm(0b00), fp0RealImag, fp0RealImag, fp0Real)
	VSHUFPD(Imm(0b11), fp0RealImag, fp0RealImag, fp0Imag)

	fp1ImagReal := XMM()
	VSHUFPD(Imm(0b01), fp1RealImag, fp1RealImag, fp1ImagReal)

	fpOutRealImag := XMM()
	VMULPD(fp1ImagReal, fp0Imag, fpOutRealImag)
	VFMADDSUB231PD(fp1RealImag, fp0Real, fpOutRealImag)

	VMOVUPD(fpOutRealImag, Mem{Base: fpOut, Disp: 16, Scale: 8})

	i := GP64()
	MOVQ(U64(4), i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	fp0RealImag, fp1RealImag = YMM(), YMM()
	VMOVUPD(Mem{Base: fp0, Index: i, Scale: 8}, fp0RealImag)
	VMOVUPD(Mem{Base: fp1, Index: i, Scale: 8}, fp1RealImag)

	fp0Real, fp0Imag = YMM(), YMM()
	VSHUFPD(Imm(0b0000), fp0RealImag, fp0RealImag, fp0Real)
	VSHUFPD(Imm(0b1111), fp0RealImag, fp0RealImag, fp0Imag)

	fp1ImagReal = YMM()
	VSHUFPD(Imm(0b0101), fp1RealImag, fp1RealImag, fp1ImagReal)

	fpOutRealImag = YMM()
	VMULPD(fp1ImagReal, fp0Imag, fpOutRealImag)
	VFMADDSUB231PD(fp1RealImag, fp0Real, fpOutRealImag)

	VMOVUPD(fpOutRealImag, Mem{Base: fpOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, N)
	JL(LabelRef("loop_body"))

	RET()
}
