//go:generate go run . -rfft -out ../../math/poly/asm_rfft_amd64.s -stubs ../../math/poly/asm_rfft_stub_amd64.go -pkg=poly
//go:generate go run . -convert -out ../../math/poly/asm_convert_amd64.s -stubs ../../math/poly/asm_convert_stub_amd64.go -pkg=poly
//go:generate go run . -vecFloat64 -out ../../math/poly/asm_vec_float64_amd64.s -stubs ../../math/poly/asm_vec_float64_stub_amd64.go -pkg=poly
//go:generate go run . -vecUint64 -out ../../math/vec/asm_vec_amd64.s -stubs ../../math/vec/asm_vec_stub_amd64.go -pkg=vec
//go:generate go run . -decompose -out ../../carousel/asm_decompose_amd64.s -stubs ../../carousel/asm_decompose_stub_amd64.go -pkg=carousel
package main

import (
	"flag"

	. "github.com/mmcloughlin/avo/build"
	"github.com/mmcloughlin/avo/buildtags"
)

var (
	rfft       = flag.Bool("rfft", false, "asm_rfft_amd64.s")
	convert    = flag.Bool("convert", false, "asm_convert_amd64.s")
	vecFloat64 = flag.Bool("vecFloat64", false, "asm_vec_amd64.s")

	vecUint64 = flag.Bool("vecUint64", false, "asm_vec_amd64.s")

	decompose = flag.Bool("decompose", false, "asm_decompose_amd64.s")
)

func main() {
	flag.Parse()

	Constraint(buildtags.Term("amd64"))
	Constraint(buildtags.Not("purego"))

	if *rfft {
		rfftInPlaceAVX2()
		rifftInPlaceAVX2()
		convolveAssignAVX2()
	}

	if *convert {
		convertConstants()

		convertUint64ToFloat64AssignAVX2()
		convertFloat64ToUint64AssignAVX2()

		floatModQInPlaceAVX2()
	}

	if *vecFloat64 {
		addFloat64AssignAVX2()
		subFloat64AssignAVX2()
		negFloat64AssignAVX2()

		scalarMulFloat64AssignAVX2()
		scalarMulAddFloat64AssignAVX2()
		scalarMulSubFloat64AssignAVX2()

		elementWiseMulFloat64AssignAVX2()
		elementWiseMulAddFloat64AssignAVX2()
		elementWiseMulSubFloat64AssignAVX2()
	}

	if *vecUint64 {
		addUint64AssignAVX2()
		subUint64AssignAVX2()

		scalarMulUint64AssignAVX2()
		scalarMulAddUint64AssignAVX2()
		scalarMulSubUint64AssignAVX2()

		elementWiseMulUint64AssignAVX2()
		elementWiseMulAddUint64AssignAVX2()
		elementWiseMulSubUint64AssignAVX2()
	}

	if *decompose {
		decomposeConstants()
		decomposePolyAssignUint64AVX2()
	}

	Generate()
}
