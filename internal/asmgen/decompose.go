package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func decomposeConstants() {
	ConstData("ONE", U64(1))
}

func decomposePolyAssignUint64AVX2() {
	TEXT("decomposePolyAssignUint64AVX2", NOSPLIT, "func(p []uint64, base uint64, logBase uint64, logLastBaseQ uint64, decomposedOut [][]uint64)")
	Pragma("noescape")

	p := Load(Param("p").Base(), GP64())
	decomposedOut := Load(Param("decomposedOut").Base(), GP64())

	N := Load(Param("p").Len(), GP64())
	level := Load(Param("decomposedOut").Len(), GP64())

	base, logBase, logLastBaseQ := YMM(), YMM(), YMM()
	VPBROADCASTQ(NewParamAddr("base", 24), base)
	VPBROADCASTQ(NewParamAddr("logBase", 32), logBase)
	VPBROADCASTQ(NewParamAddr("logLastBaseQ", 40), logLastBaseQ)

	one := YMM()
	VPBROADCASTQ(NewDataAddr(NewStaticSymbol("ONE"), 0), one)

	baseMask, baseHalf, logBaseSubOne := YMM(), YMM(), YMM()
	VPSUBQ(one, base, baseMask)
	VPSRLQ(Imm(1), base, baseHalf)
	VPSUBQ(one, logBase, logBaseSubOne)

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("N_loop_end"))
	Label("N_loop_body")

	c := YMM()
	VMOVDQU(Mem{Base: p, Index: i, Scale: 8}, c)

	cDiv, cCarry := YMM(), YMM()
	VPSRLVQ(logLastBaseQ, c, cDiv)
	VPSLLQ(Imm(1), c, cCarry)
	VPSRLVQ(logLastBaseQ, cCarry, cCarry)
	VANDPD(one, cCarry, cCarry)
	VPADDQ(cCarry, cDiv, c)

	j := GP64()
	MOVQ(level, j)
	SUBQ(Imm(1), j)

	jj := GP64()
	MOVQ(j, jj)
	ADDQ(j, jj)
	ADDQ(j, jj)
	JMP(LabelRef("level_loop_end"))
	Label("level_loop_body")

	u := YMM()
	VANDPD(baseMask, c, u)
	VPSRLVQ(logBase, c, c)

	uDiv := YMM()
	VPSRLVQ(logBaseSubOne, u, uDiv)
	VPADDQ(uDiv, c, c)

	uCarry := YMM()
	VANDPD(baseHalf, u, uCarry)
	VPSLLQ(Imm(1), uCarry, uCarry)
	VPSUBQ(uCarry, u, u)

	decomposedOutj := GP64()
	MOVQ(Mem{Base: decomposedOut, Index: jj, Scale: 8}, decomposedOutj)
	VMOVDQU(u, Mem{Base: decomposedOutj, Index: i, Scale: 8})

	SUBQ(Imm(1), j)
	SUBQ(Imm(3), jj)

	Label("level_loop_end")
	CMPQ(j, Imm(1))
	JGE(LabelRef("level_loop_body"))

	u = YMM()
	VANDPD(baseMask, c, u)

	uCarry = YMM()
	VANDPD(baseHalf, u, uCarry)
	VPSLLQ(Imm(1), uCarry, uCarry)
	VPSUBQ(uCarry, u, u)

	decomposedOut0 := GP64()
	MOVQ(Mem{Base: decomposedOut}, decomposedOut0)
	VMOVDQU(u, Mem{Base: decomposedOut0, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("N_loop_end")
	CMPQ(i, N)
	JL(LabelRef("N_loop_body"))

	RET()
}
