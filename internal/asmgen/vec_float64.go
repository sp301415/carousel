package main

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func addFloat64AssignAVX2() {
	TEXT("addFloat64AssignAVX2", NOSPLIT, "func(v0, v1, vOut []float64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	v1 := Load(Param("v1").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0, x1 := YMM(), YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)
	VMOVDQU(Mem{Base: v1, Index: i, Scale: 8}, x1)

	xOut := YMM()
	VADDPD(x1, x0, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, N)
	JL(LabelRef("loop_body"))

	RET()
}

func subFloat64AssignAVX2() {
	TEXT("subFloat64AssignAVX2", NOSPLIT, "func(v0, v1, vOut []float64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	v1 := Load(Param("v1").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0, x1 := YMM(), YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)
	VMOVDQU(Mem{Base: v1, Index: i, Scale: 8}, x1)

	xOut := YMM()
	VSUBPD(x1, x0, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, N)
	JL(LabelRef("loop_body"))

	RET()
}

func negFloat64AssignAVX2() {
	TEXT("negFloat64AssignAVX2", NOSPLIT, "func(v0, vOut []float64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	zero := YMM()
	VPXOR(zero, zero, zero)

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0 := YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)

	xOut := YMM()
	VSUBPD(x0, zero, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, N)
	JL(LabelRef("loop_body"))

	RET()
}

func scalarMulFloat64AssignAVX2() {
	TEXT("scalarMulFloat64AssignAVX2", NOSPLIT, "func(v0 []float64, c float64, vOut []float64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	c := YMM()
	VBROADCASTSD(NewParamAddr("c", 24), c)

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0 := YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)

	xOut := YMM()
	VMULPD(x0, c, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, N)
	JL(LabelRef("loop_body"))

	RET()
}

func scalarMulAddFloat64AssignAVX2() {
	TEXT("scalarMulAddFloat64AssignAVX2", NOSPLIT, "func(v0 []float64, c float64, vOut []float64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	c := YMM()
	VBROADCASTSD(NewParamAddr("c", 24), c)

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0, xOut := YMM(), YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)
	VMOVDQU(Mem{Base: vOut, Index: i, Scale: 8}, xOut)

	VFMADD231PD(x0, c, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, N)
	JL(LabelRef("loop_body"))

	RET()
}

func scalarMulSubFloat64AssignAVX2() {
	TEXT("scalarMulSubFloat64AssignAVX2", NOSPLIT, "func(v0 []float64, c float64, vOut []float64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	c := YMM()
	VBROADCASTSD(NewParamAddr("c", 24), c)

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0, xOut := YMM(), YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)
	VMOVDQU(Mem{Base: vOut, Index: i, Scale: 8}, xOut)

	VFNMADD231PD(x0, c, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, N)
	JL(LabelRef("loop_body"))

	RET()
}

func elementWiseMulFloat64AssignAVX2() {
	TEXT("elementWiseMulFloat64AssignAVX2", NOSPLIT, "func(v0, v1, vOut []float64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	v1 := Load(Param("v1").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0, x1 := YMM(), YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)
	VMOVDQU(Mem{Base: v1, Index: i, Scale: 8}, x1)

	xOut := YMM()
	VMULPD(x1, x0, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, N)
	JL(LabelRef("loop_body"))

	RET()
}

func elementWiseMulAddFloat64AssignAVX2() {
	TEXT("elementWiseMulAddFloat64AssignAVX2", NOSPLIT, "func(v0, v1, vOut []float64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	v1 := Load(Param("v1").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0, x1, xOut := YMM(), YMM(), YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)
	VMOVDQU(Mem{Base: v1, Index: i, Scale: 8}, x1)
	VMOVDQU(Mem{Base: vOut, Index: i, Scale: 8}, xOut)

	VFMADD231PD(x1, x0, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, N)
	JL(LabelRef("loop_body"))

	RET()
}

func elementWiseMulSubFloat64AssignAVX2() {
	TEXT("elementWiseMulSubFloat64AssignAVX2", NOSPLIT, "func(v0, v1, vOut []float64)")
	Pragma("noescape")

	v0 := Load(Param("v0").Base(), GP64())
	v1 := Load(Param("v1").Base(), GP64())
	vOut := Load(Param("vOut").Base(), GP64())
	N := Load(Param("vOut").Len(), GP64())

	i := GP64()
	XORQ(i, i)
	JMP(LabelRef("loop_end"))
	Label("loop_body")

	x0, x1, xOut := YMM(), YMM(), YMM()
	VMOVDQU(Mem{Base: v0, Index: i, Scale: 8}, x0)
	VMOVDQU(Mem{Base: v1, Index: i, Scale: 8}, x1)
	VMOVDQU(Mem{Base: vOut, Index: i, Scale: 8}, xOut)

	VFNMADD231PD(x1, x0, xOut)

	VMOVDQU(xOut, Mem{Base: vOut, Index: i, Scale: 8})

	ADDQ(Imm(4), i)

	Label("loop_end")
	CMPQ(i, N)
	JL(LabelRef("loop_body"))

	RET()
}
