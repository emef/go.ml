package matrix

func VecScale(scalar float64, X []float64) []float64 {
	Y := make([]float64, len(X))
	for i, x := range X {
		Y[i] = x * scalar
	}
	return Y
}

func VecAdd(X, Y []float64) []float64 {
	Z := make([]float64, len(X))
	for i := range X {
		Z[i] = X[i] + Y[i]
	}
	return Z
}