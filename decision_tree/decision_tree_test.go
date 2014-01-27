package decision_tree

import (
	"os"
	"fmt"
	"strconv"
	"runtime"
	"encoding/csv"
	"testing"
	"go.ml/metrics"
	"go.ml/util"
)

func TestWithIris(t *testing.T) {
	X, y := loadIris()
	perf := testDataset(X, y)
	fmt.Printf("Iris: %.3f\n", perf)
}

func TestBreastCancer(t *testing.T) {
	X, y := loadBreastCancer()
	perf := testDataset(X, y)
	fmt.Printf("Breast Cancer: %.3f\n", perf)
}

func testDataset(X [][]float64, y []float64) float64 {
	runtime.GOMAXPROCS(runtime.NumCPU())

	util.RandomShuffle(X, y)

	kFolds := 5
	foldSize := int(float64(len(X)) / float64(kFolds))
	foldSum := 0.0

	for k := 0; k < kFolds; k++ {
		start := k * foldSize
		for j := 0; j < foldSize; j++ {
			X[j], X[start+j] = X[start+j], X[j]
			y[j], y[start+j] = y[start+j], y[j]
		}
		tree := Fit(X[foldSize:], y[foldSize:])
		yPred := tree.Classify(X[:foldSize])
		foldSum += metrics.Accuracy(yPred, y[:foldSize])
	}

	return foldSum / float64(kFolds)
}

func loadIris() ([][]float64, []float64) {
	X_file, _ := os.Open("iris_X.csv")
	y_file, _ := os.Open("iris_y.csv")

	X := make([][]float64, 0)
	y := make([]float64, 0)

	reader := csv.NewReader(X_file)
	for {
		line, err := reader.Read()
		if err != nil {
			break
		}

		x := make([]float64, len(line))
		for i, strVal := range line {
			x[i], _ = strconv.ParseFloat(strVal, 64)
		}

		X = append(X,  x)
	}

	reader = csv.NewReader(y_file)
	for {
		line, err := reader.Read()
		if err != nil {
			break
		}

		yVal, _ := strconv.ParseFloat(line[0], 64)
		y = append(y,  yVal)
	}

	return X, y
}

func loadBreastCancer() ([][]float64, []float64) {
	file, _ := os.Open("breast_cancer.csv")
	reader := csv.NewReader(file)

	var X [][]float64
	var y []float64

	for {
		line, err := reader.Read()
		if err != nil {
			break
		}

		if line[1] == "M" {
			y = append(y, 1.0)
		} else {
			y = append(y, 0.0)
		}

		x := make([]float64, 32)
		for i, strVal := range line[2:] {
			x[i], _ = strconv.ParseFloat(strVal, 64)
		}

		X = append(X,  x)
	}

	return X, y

}