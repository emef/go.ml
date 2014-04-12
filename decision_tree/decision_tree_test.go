package decision_tree

import (
	"fmt"
	"runtime"
	"testing"
	"github.com/emef/go.ml/metrics"
	"github.com/emef/go.ml/datasets"
)

func TestWithIris(t *testing.T) {
	X, y := datasets.Load("iris")
	perf := testDataset(X, y)
	fmt.Printf("Iris: %.3f\n", perf)
}

func TestBreastCancer(t *testing.T) {
	X, y := datasets.Load("cancer")
	perf := testDataset(X, y)
	fmt.Printf("Breast Cancer: %.3f\n", perf)
}

func TestPrint(t *testing.T) {
	X, y := datasets.Load("cancer")
	tree, _ := DecisionTree(4, GINI)
	tree.Fit(X, y)
	fmt.Println(tree)
}

func TestMeanSquared(t *testing.T) {
	X, y := datasets.Load("cancer")
	tree, _ := DecisionTree(4, GINI)
	tree.Fit(X, y)
	yPred := tree.Classify(X)
	err := metrics.MeanSquaredError(yPred, y)
	fmt.Printf("cancer mean sq: %.3f\n", err)
}

func testDataset(X [][]float64, y []float64) float64 {
	runtime.GOMAXPROCS(runtime.NumCPU())

	datasets.RandomShuffle(X, y)

	kFolds := 5
	foldSize := int(float64(len(X)) / float64(kFolds))
	foldSum := 0.0

	for k := 0; k < kFolds; k++ {
		start := k * foldSize
		for j := 0; j < foldSize; j++ {
			X[j], X[start+j] = X[start+j], X[j]
			y[j], y[start+j] = y[start+j], y[j]
		}
		tree, _ := DecisionTree(1, GINI)
		tree.Fit(X[foldSize:], y[foldSize:])
		yPred := tree.Classify(X[:foldSize])
		foldSum += metrics.Accuracy(yPred, y[:foldSize])
	}

	return foldSum / float64(kFolds)
}
