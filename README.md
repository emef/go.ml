go.ml
=====

Machine learning and optimization in golang

Implementations: decision_tree, linear regression (OLS), nonlinear optimization

-----

**decision trees**

```go

import (
  "github.com/emef/go.ml/datasets"
  "github.com/emef/go.ml/metrics"
  "github.com/emef/go.ml/decision_tree"
)

maxDepth := 10
splitMethod := decision_tree.GINI

tree, err := decision_tree.DecisionTree(maxDepth, splitMethod)
if err != nil {
  fmt.Println(err)
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

-----

**linear regression (OLS)**

this is the current state... will implement as a predictor interface with fit and classify later.

TODO:

* most likely am not handling singular matrices properly when taking inverse
* implement as type w/ Fit and Classify functions
* clean up matrix package a bit (this is the only thing using it so far)

```go

import (
  "github.com/emef/go.ml/datasets"
  "github.com/emef/go.ml/metrics"
  "github.com/emef/go.ml/matrix"
  "github.com/emef/go.ml/linear_model"
)

X, y := datasets.Load("iris")
datasets.RandomShuffle(X, y)
XTrain, XTest := X[:67], X[67:]
yTrain, yTest := y[:67], y[67:]

beta := linear_model.LinearRegression(XTrain, yTrain)

// validate on held out data
yPred := matrix.VecMult(XTest, beta)
fmt.Println(metrics.Accuracy(yPred, yTest))

```

-----

**nonlinear optimization**

```go

import (
  "fmt"
  "math"
  "github.com/emef/go.ml/optimization"
)

f := func(x []float64) float64 {
  return math.Sin(x[0]) * math.Pow(math.Cos(x[1]), 2)
}

x := []float64{1, 1}  // initial guess
f_min := optimization.GradientDescent(f, x)
fmt.Println(f_min)
```
