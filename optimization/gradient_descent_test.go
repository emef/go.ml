package optimization

import (
	"fmt"
	"math"
	"testing"
)

func fn2d(f func(x, y float64) float64) fn {
	transformed := func(X []float64) float64 {
		return f(X[0], X[1])
	}

	return transformed
}

func TestGradientDescent(t *testing.T) {
	f := fn2d(func(x, y float64) float64 {
		return math.Sin(x)*math.Pow(math.Sin(y), 2)
	})

	X := []float64{1, 1}
	Y := GradientDescent(f, X)

	fmt.Println(Y, f(Y))
}
