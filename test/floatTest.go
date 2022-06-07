package test

import (
	"fmt"
	"strconv"
)

func FloatTest() {
	var (
		test1 float64
		test2 float64
		num   int
	)
	num = 7
	test1 = float64(1) / float64(num)
	getstr := strconv.FormatFloat(test1, 'f', 2, 64)
	fmt.Printf("[num = %v], [test1 = %v], [test2 = %v], [getstr = %v]\n", num, test1, test2, getstr)
	return
}
