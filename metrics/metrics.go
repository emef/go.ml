package metrics

import (
	"math"
)

func Accuracy(yPred, yTrue []float64) float64 {
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


func MeanSquaredError(yPred, yTrue []float64) float64 {
	if len(yPred) != len(yTrue) {
		return -1
	}

	err := 0.0

	for i := range yPred {
		err += math.Pow(yPred[i] - yTrue[i], 2)
	}

	return err / float64(len(yPred))
}
