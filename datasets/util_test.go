package datasets

import (
	"testing"
	"fmt"
)

func TestRandomShuffle(t *testing.T) {
	X := [][]float64{[]float64{1, 2, 3}, []float64{4, 5, 6}}
	y := []float64{0, 1}
	RandomShuffle(X, y)
	fmt.Println(X, y)
}