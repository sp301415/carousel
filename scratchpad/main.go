package main

import (
	"fmt"

	"github.com/sp301415/carousel/math/poly"
)

func main() {
	polyParams := poly.NewEvaluatorParametersForPacking(2048, 3)
	resol := polyParams.Resolution()
	fmt.Println(resol)
}
