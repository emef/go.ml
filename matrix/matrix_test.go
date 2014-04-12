package matrix

import (
	"fmt"
	"testing"
)

func TestMatMult(t *testing.T) {
	A := [][]float64{{1, 2, 3},{4, 5, 6}}
	B := [][]float64{{7, 8}, {9, 10}, {11, 12}}
	fmt.Println(MatMult(A, B))
}

func TestMatMultTrans(t *testing.T) {
	A := [][]float64{{1, 4}, {2, 5}, {3, 6}}
	B := [][]float64{{7, 8}, {9, 10}, {11, 12}}
	fmt.Println(MatMultTrans(A, B))
}

func TestVecMultTrans(t *testing.T) {
	A := [][]float64{{1,4,7,10}, {2,5,8,11}, {3,6,9,12}}
	b := []float64{-2, 1, 0}
	fmt.Println("huh", VecMultTrans(A, b))
}

func TestDeterminant(t *testing.T) {
	X := [][]float64{{6, 1, 1}, {4, -2, 5}, {2, 8, 7}}
	fmt.Println(Determinant(X))
}

func TestCoFactor(t *testing.T) {
	X := [][]float64{{1, 2, 0}, {-1, 1, 1}, {1, 2, 3}}
	fmt.Println("cofactor", CoFactor(X))
}

func TestInverse(t *testing.T) {
	X := [][]float64{{2, -1, 0}, {-1, 2, -1}, {0, -1, 2}}
	fmt.Println("inverse ", MatDirtyInverse(X))
}