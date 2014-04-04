package datasets

import (
	"math/rand"
	"time"
)

// destructively shuffle input
func RandomShuffle(X [][]float64, y []float64) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len(X) - 1; i++ {
		newIndex := i + int(rand.Float64() * float64(len(X) - i))
		X[i], X[newIndex] = X[newIndex], X[i]
		y[i], y[newIndex] = y[newIndex], y[i]
	}
}