package linear_model

import (
	"testing"
	"fmt"
	"math/rand"
	"github.com/emef/go.ml/matrix"
	"github.com/emef/go.ml/datasets"
	"github.com/emef/go.ml/metrics"
)

func tTestBeta(t *testing.T) {
	X := [][]float64{{1,0,5}, {2,5,4}, {3,6,5}, {8,1,1}}
	y := []float64{1, 0.5, 0.75, 0.2}
	beta := LinearRegression(X, y)
	fmt.Println(y)
	fmt.Println(matrix.VecMult(X, beta))
}

func tTestSimple(t *testing.T) {
	X := [][]float64{{1, 0}, {1, 0.5}, {1, 1}, {1, 1.5}}
	y := []float64{0.3, 0.4, 0.55, 0.6}
	beta := LinearRegression(X, y)
	fmt.Println(y, beta)
	fmt.Println(matrix.VecMult(X, beta))
}

func tTestBig(t *testing.T) {
	n := 10000
	m := 50
	X := make([][]float64, n)
	y := make([]float64, n)

	for i := range X {
		X[i] = make([]float64, m)
		y[i] = rand.Float64()
		for j := 0; j < m; j++ {
			X[i][j] = rand.Float64()
		}
	}

	LinearRegression(X, y)
}


func TestIris(t *testing.T) {
	X, y := datasets.Load("iris")
	beta := LinearRegression(X, y)
	yPred := matrix.VecMult(X, beta)
	fmt.Println("iris error", metrics.MeanSquaredError(yPred, y))
}

func TestCancer(t *testing.T) {
	X, y := datasets.Load("cancer")
	beta := LinearRegression(X, y)
	yPred := matrix.VecMult(X, beta)
	fmt.Println("cancer error", metrics.MeanSquaredError(yPred, y))
}
