package matrix

import (
	"bytes"
	"fmt"
	"math"
	"tetris/utils"
)

type Matrix [][]int

func Multiply(m1, m2 Matrix) Matrix {
	utils.Asserts(len(m1) > 0, "Matrix has size of 0")
	utils.Asserts(len(m2) == len(m1[0]), fmt.Sprintf("Matrix has different row and column size: %d by %d", len(m2), len(m1[0])))

	m, p, n := len(m1), len(m2[0]), len(m2)
	newMatrix := make(Matrix, m)

	for i := range m {
		newMatrix[i] = make([]int, p)
	}

	for i := range m {
		for k := range n {
			for j := range p {
				newMatrix[i][j] += m1[i][k] * m2[k][j]
			}
		}
	}

	return newMatrix
}

func (m Matrix) Transpose() Matrix {
	utils.Asserts(len(m) > 0, "cannot transpose empty matrix")
	utils.Asserts(len(m[0]) > 0, "cannot transpose empty matrix")

	transposedMatrix := make(Matrix, len(m[0]))

	for i := range len(m[0]) {
		transposedMatrix[i] = make([]int, len(m))
	}

	for i := range len(m) {
		for j := range len(m[0]) {
			transposedMatrix[j][i] = m[i][j]
		}
	}

	return transposedMatrix
}

func (m Matrix) Abs() {
	for i := range len(m) {
		for j := range len(m[i]) {
			m[i][j] = int(math.Abs(float64(m[i][j])))
		}
	}
}

func (m Matrix) Minus(vector []int) {
	utils.Asserts(len(vector) == len(m[0]), "cannot minus with different size vector")

	for i := range len(m) {
		for j := range len(m[i]) {
			m[i][j] -= vector[j]
		}
	}
}

func (m Matrix) Add(vector []int) {
	utils.Asserts(len(vector) == len(m[0]), "cannot minus with different size vector")

	for i := range len(m) {
		for j := range len(m[i]) {
			m[i][j] += vector[j]
		}
	}
}

func (m1 Matrix) Equal(m2 Matrix) bool {
	if len(m1) != len(m2) {
		return false
	}

	for i := range len(m1) {
		for j := range len(m1[0]) {
			if m1[i][j] != m2[i][j] {
				return false
			}
		}
	}

	return true
}

func Copy(m1 Matrix) Matrix {
	copiedMatrix := make(Matrix, len(m1))

	for i := range len(m1) {
		copiedMatrix[i] = make([]int, len(m1[i]))
	}

	for i := range len(m1) {
		for j := range len(m1[i]) {
			copiedMatrix[i][j] = m1[i][j]
		}
	}

	return copiedMatrix
}

func (m1 Matrix) ToString() string {
	var repr bytes.Buffer

	for i := range len(m1) {
		repr.WriteString("\n[")
		for j := range len(m1[i]) {
			repr.WriteString(fmt.Sprintf("%d ", m1[i][j]))
		}
		repr.WriteString("]")
	}

	return repr.String()
}
