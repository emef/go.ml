go.ml
=====

Machine learning in golang

Implementations: decision_tree

-----

**decision trees**

```go

import (
  "go.ml/datasets"
  "go.ml/metrics"
  "go.ml/decision_tree"
)

maxDepth := 10
splitMethod := decision_tree.GINI

tree, err := decision_tree.DecisionTree(max_depth, splitMethod)
if err != nil {
  fmt.Println(e)
  panic()
}

// load a sample dataset, X is [][]float64, and y is []float64
X, y := datasets.Load("iris")

// random subset to train
datasets.RandomShuffle(X, y)
XTrain, XTest := X[:67], X[67:]
yTrain, yTest := y[:67], y[67:]

tree.Fit(XTrain, yTrain)

// validate on held out data
yPred := tree.Classify(XTest)
fmt.Println(metrics.Accuracy(yPred, yTest))
```
