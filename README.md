go.ml
=====

Machine learning and optimization in golang

Implementations: decision_tree, linear regression (OLS), nonlinear optimization, genetic algorithm

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

-----

**genetic algorithm**

Not necessarily sold on this implementation or interface; we'll see...

```go
import (
  "math/rand"
  "github.com/emef/go.ml/genetic"
)

/*
A BitSample is represented by a bitfield (see github.com/emef/bitfield)

The goal of this example is to coax the random samples into our target
values (as specified in fitness())
*/

func fitness(x genetic.Sample) float64 {
	sample := x.(genetic.BitSample)
	target := []int{0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 1, 1, 1}
	correct := 0.0
	for i := 0; i < sample.field.Size(); i++ {
		if boolToInt(sample.field.Test(uint32(i))) == target[i] {
			correct++
		}
	}
	return correct / float64(sample.field.Size())
}

/*
Start with 20 random bit fields in `samples`
Run the GA with the mutation probability at 20%, max of 50 epochs, and a target
fitness value of 1.0. g.Run() will return the best sample it can find after all
epochs have finished or target was found.
*/
samples := make([]genetic.Sample, 20)
for i := range samples {
	sample := genetic.NewBitSample(20)
	indexes := rand.Perm(20)[:randInt(20)]
	for _, j := range indexes {
		sample.field.Set(uint32(j))
	}
	samples[i] = sample
}
g := genetic.New(samples, fitness, 0.2, 50, 1.0)
fmt.Println(g.Run())
```
