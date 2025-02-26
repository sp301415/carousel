package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	"github.com/mmcloughlin/avo/reg"
)

func addUint64AssignAVX2() {
	TEXT("addUint64AssignAVX2", NOSPLIT, "func(v0, v1, vOut []uint64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	v1 := Load(Param("v1").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	NN := GP64()
	MOVQ(N, NN)
	SHRQ(Imm(2), NN)
	SHLQ(Imm(2), NN)

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0, x1 := YMM(), YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)
	VMOVDQU(Mem{Base: v1, Index: i, Scale: 8}, x1)

	xOut := YMM()
	VPADDQ(x1, x0, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, NN)
	JL(LabelRef("loop_body"))

	JMP(LabelRef("leftover_loop_end"))
	Label("leftover_loop_body")

	y0, y1 := GP64(), GP64()
	MOVQ(Mem{Base: v0, Index: i, Scale: 8}, y0)
	MOVQ(Mem{Base: v1, Index: i, Scale: 8}, y1)

	ADDQ(y1, y0)

	MOVQ(y0, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(1), i)

	Label("leftover_loop_end")
	CMPQ(i, N)
	JL(LabelRef("leftover_loop_body"))

	RET()
}

func subUint64AssignAVX2() {
	TEXT("subUint64AssignAVX2", NOSPLIT, "func(v0, v1, vOut []uint64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	v1 := Load(Param("v1").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	NN := GP64()
	MOVQ(N, NN)
	SHRQ(Imm(2), NN)
	SHLQ(Imm(2), NN)

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0, x1 := YMM(), YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)
	VMOVDQU(Mem{Base: v1, Index: i, Scale: 8}, x1)

	xOut := YMM()
	VPSUBQ(x1, x0, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, NN)
	JL(LabelRef("loop_body"))

	JMP(LabelRef("leftover_loop_end"))
	Label("leftover_loop_body")

	y0, y1 := GP64(), GP64()
	MOVQ(Mem{Base: v0, Index: i, Scale: 8}, y0)
	MOVQ(Mem{Base: v1, Index: i, Scale: 8}, y1)

	SUBQ(y1, y0)

	MOVQ(y0, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(1), i)

	Label("leftover_loop_end")
	CMPQ(i, N)
	JL(LabelRef("leftover_loop_body"))

	RET()
}

func precomputeMulUint64(x, xSwap reg.VecVirtual) {
	VPSHUFD(Imm(0b10110001), x, xSwap)
}

func mulUint64(x0, x0Swap, x1, xOut reg.VecVirtual) {
	xCross := YMM()
	VPMULLD(x1, x0Swap, xCross)

	xOutHiLo := YMM()
	VPSRLQ(Imm(32), xCross, xOutHiLo)

	xOut0 := YMM()
	VPADDQ(xCross, xOutHiLo, xOut0)
	VPSLLQ(Imm(32), xOut0, xOut0)

	xOut1 := YMM()
	VPMULUDQ(x1, x0, xOut1)

	VPADDQ(xOut1, xOut0, xOut)
}

func scalarMulUint64AssignAVX2() {
	TEXT("scalarMulUint64AssignAVX2", NOSPLIT, "func(v0 []uint64, c uint64, vOut []uint64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	NN := GP64()
	MOVQ(N, NN)
	SHRQ(Imm(2), NN)
	SHLQ(Imm(2), NN)

	cv, cvSwap := YMM(), YMM()
	VPBROADCASTQ(NewParamAddr("c", 24), cv)
	precomputeMulUint64(cv, cvSwap)

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0 := YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)

	xOut := YMM()
	mulUint64(cv, cvSwap, x0, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, NN)
	JL(LabelRef("loop_body"))

	c := GP64()
	Load(Param("c"), c)

	JMP(LabelRef("leftover_loop_end"))
	Label("leftover_loop_body")

	y0 := GP64()
	MOVQ(Mem{Base: v0, Index: i, Scale: 8}, y0)

	IMULQ(c, y0)

	MOVQ(y0, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(1), i)

	Label("leftover_loop_end")
	CMPQ(i, N)
	JL(LabelRef("leftover_loop_body"))

	RET()
}

func scalarMulAddUint64AssignAVX2() {
	TEXT("scalarMulAddUint64AssignAVX2", NOSPLIT, "func(v0 []uint64, c uint64, vOut []uint64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	NN := GP64()
	MOVQ(N, NN)
	SHRQ(Imm(2), NN)
	SHLQ(Imm(2), NN)

	cv, cvSwap := YMM(), YMM()
	VPBROADCASTQ(NewParamAddr("c", 24), cv)
	precomputeMulUint64(cv, cvSwap)

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0 := YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)
	xOut := YMM()
	VMOVDQU(Mem{Base: vOut, Index: i, Scale: 8}, xOut)

	mulUint64(cv, cvSwap, x0, x0)
	VPADDQ(x0, xOut, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, NN)
	JL(LabelRef("loop_body"))

	c := GP64()
	Load(Param("c"), c)

	JMP(LabelRef("leftover_loop_end"))
	Label("leftover_loop_body")

	y0 := GP64()
	MOVQ(Mem{Base: v0, Index: i, Scale: 8}, y0)
	yOut := GP64()
	MOVQ(Mem{Base: vOut, Index: i, Scale: 8}, yOut)

	IMULQ(c, y0)
	ADDQ(y0, yOut)

	MOVQ(yOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(1), i)

	Label("leftover_loop_end")
	CMPQ(i, N)
	JL(LabelRef("leftover_loop_body"))

	RET()
}

func scalarMulSubUint64AssignAVX2() {
	TEXT("scalarMulSubUint64AssignAVX2", NOSPLIT, "func(v0 []uint64, c uint64, vOut []uint64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	NN := GP64()
	MOVQ(N, NN)
	SHRQ(Imm(2), NN)
	SHLQ(Imm(2), NN)

	cv, cvSwap := YMM(), YMM()
	VPBROADCASTQ(NewParamAddr("c", 24), cv)
	precomputeMulUint64(cv, cvSwap)

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0 := YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)
	xOut := YMM()
	VMOVDQU(Mem{Base: vOut, Index: i, Scale: 8}, xOut)

	mulUint64(cv, cvSwap, x0, x0)
	VPSUBQ(x0, xOut, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, NN)
	JL(LabelRef("loop_body"))

	c := GP64()
	Load(Param("c"), c)

	JMP(LabelRef("leftover_loop_end"))
	Label("leftover_loop_body")

	y0 := GP64()
	MOVQ(Mem{Base: v0, Index: i, Scale: 8}, y0)
	yOut := GP64()
	MOVQ(Mem{Base: vOut, Index: i, Scale: 8}, yOut)

	IMULQ(c, y0)
	SUBQ(y0, yOut)

	MOVQ(yOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(1), i)

	Label("leftover_loop_end")
	CMPQ(i, N)
	JL(LabelRef("leftover_loop_body"))

	RET()
}

func elementWiseMulUint64AssignAVX2() {
	TEXT("elementWiseMulUint64AssignAVX2", NOSPLIT, "func(v0, v1, vOut []uint64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	v1 := Load(Param("v1").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	NN := GP64()
	MOVQ(N, NN)
	SHRQ(Imm(2), NN)
	SHLQ(Imm(2), NN)

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0 := YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)
	x1 := YMM()
	VMOVDQU(Mem{Base: v1, Index: i, Scale: 8}, x1)

	x0Swap := YMM()
	precomputeMulUint64(x0, x0Swap)

	xOut := YMM()
	mulUint64(x0, x0Swap, x1, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, NN)
	JL(LabelRef("loop_body"))

	JMP(LabelRef("leftover_loop_end"))
	Label("leftover_loop_body")

	y0 := GP64()
	MOVQ(Mem{Base: v0, Index: i, Scale: 8}, y0)
	y1 := GP64()
	MOVQ(Mem{Base: v1, Index: i, Scale: 8}, y1)

	IMULQ(y1, y0)

	MOVQ(y0, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(1), i)

	Label("leftover_loop_end")
	CMPQ(i, N)
	JL(LabelRef("leftover_loop_body"))

	RET()
}

func elementWiseMulAddUint64AssignAVX2() {
	TEXT("elementWiseMulAddUint64AssignAVX2", NOSPLIT, "func(v0, v1, vOut []uint64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	v1 := Load(Param("v1").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	NN := GP64()
	MOVQ(N, NN)
	SHRQ(Imm(2), NN)
	SHLQ(Imm(2), NN)

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0 := YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)
	x1 := YMM()
	VMOVDQU(Mem{Base: v1, Index: i, Scale: 8}, x1)
	xOut := YMM()
	VMOVDQU(Mem{Base: vOut, Index: i, Scale: 8}, xOut)

	x0Swap := YMM()
	precomputeMulUint64(x0, x0Swap)

	mulUint64(x0, x0Swap, x1, x1)
	VPADDQ(x1, xOut, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, NN)
	JL(LabelRef("loop_body"))

	JMP(LabelRef("leftover_loop_end"))
	Label("leftover_loop_body")

	y0 := GP64()
	MOVQ(Mem{Base: v0, Index: i, Scale: 8}, y0)
	y1 := GP64()
	MOVQ(Mem{Base: v1, Index: i, Scale: 8}, y1)
	yOut := GP64()
	MOVQ(Mem{Base: vOut, Index: i, Scale: 8}, yOut)

	IMULQ(y1, y0)
	ADDQ(y0, yOut)

	MOVQ(yOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(1), i)

	Label("leftover_loop_end")
	CMPQ(i, N)
	JL(LabelRef("leftover_loop_body"))

	RET()
}

func elementWiseMulSubUint64AssignAVX2() {
	TEXT("elementWiseMulSubUint64AssignAVX2", NOSPLIT, "func(v0, v1, vOut []uint64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	v1 := Load(Param("v1").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	NN := GP64()
	MOVQ(N, NN)
	SHRQ(Imm(2), NN)
	SHLQ(Imm(2), NN)

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0 := YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)
	x1 := YMM()
	VMOVDQU(Mem{Base: v1, Index: i, Scale: 8}, x1)
	xOut := YMM()
	VMOVDQU(Mem{Base: vOut, Index: i, Scale: 8}, xOut)

	x0Swap := YMM()
	precomputeMulUint64(x0, x0Swap)

	mulUint64(x0, x0Swap, x1, x1)
	VPSUBQ(x1, xOut, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, NN)
	JL(LabelRef("loop_body"))

	JMP(LabelRef("leftover_loop_end"))
	Label("leftover_loop_body")

	y0 := GP64()
	MOVQ(Mem{Base: v0, Index: i, Scale: 8}, y0)
	y1 := GP64()
	MOVQ(Mem{Base: v1, Index: i, Scale: 8}, y1)
	yOut := GP64()
	MOVQ(Mem{Base: vOut, Index: i, Scale: 8}, yOut)

	IMULQ(y1, y0)
	SUBQ(y0, yOut)

	MOVQ(yOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(1), i)

	Label("leftover_loop_end")
	CMPQ(i, N)
	JL(LabelRef("leftover_loop_body"))

	RET()
}
