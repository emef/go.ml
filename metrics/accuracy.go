package metrics

func Accuracy(yPred []float64, yTrue []float64) float64 {
	if len(yPred) != len(yTrue) {
		return -1
	}

	truePosNeg := 0

	for i := range yPred {
		if yPred[i] == yTrue[i] {
			truePosNeg += 1
		}
	}

	return float64(truePosNeg) / float64(len(yPred))
}