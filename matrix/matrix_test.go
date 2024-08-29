package matrix

import (
	"fmt"
	"testing"
)

func TestMultiply(t *testing.T) {
	var m1 Matrix = [][]int{{1, 1}, {2, 2}}
	var m2 Matrix = [][]int{{1}, {2}}
	var expectedResult Matrix = [][]int{{3}, {6}}

	result := Multiply(m1, m2)

	if !expectedResult.Equal(result) {
		t.Error(fmt.Printf("Matrix multiplication result is not equal\nexpected \n%s\nresult \n%s", expectedResult.ToString(), result.ToString()))
	}
}
