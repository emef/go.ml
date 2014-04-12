package linear_model

import (
	"github.com/emef/go.ml/matrix"
)

func LinearRegression(X [][]float64, y []float64) []float64 {
	Sxx := matrix.MatDirtyInverse(matrix.MatMultTrans(X, X))
	Sxy := matrix.VecMultTrans(X, y)
	beta := matrix.VecMult(Sxx, Sxy)
	return beta
}