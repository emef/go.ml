package optimization

import (
	"github.com/emef/go.ml/vector"
)

type fn func([]float64) float64
const eps float64 = 1e-6

/*
 Minimizes the given function f from the initial guess X
 using a gradient descent linesearch.

 NOTE: assumes there exists some local minimum

 arguments
 ---------
 f: function to minimize
 X: initial guess to start minimization from

 returns
 -------
 point at which local minimum was found

*/
func GradientDescent(f fn, X []float64) []float64 {
	min := f(X)

	for {
		grad := gradient(f, X)
		y, X1 := linesearch(f, X, grad)
		if y < min {
			min = y
			X = X1
		} else {
			break
		}
	}

	return X
}

/*
 Approximates partial derivatives of f at X0

 arguments
 ---------
 f: function to approximate
 X: point at which to approximate

*/
func gradient(f fn, X []float64) []float64 {
	var f0, f1, x_i float64
	G := make([]float64, len(X))
	f0 = f(X)
	for i := range X {
		x_i = X[i]
		X[i] += eps
		f1 = f(X)
		G[i] = (f1 - f0) / eps
		X[i] = x_i
	}
	return G
}

/*
 Performs a line search of given function using success values in the
 fibonacci series as the directional magnitude.

 arguments
 ---------
 f:   function to perform linesearch on
 X:   starting point
 dir: direction to search

 returns
 -------
 (y, X) where
     y: f(X)
     X: minimum value along this direction
*/
func linesearch(f fn, X []float64, dir []float64) (float64, []float64) {
	a, b := 0.0, 1.0
	X0 := X

	// X will always be our current location
	// y will be current minimum of f (always at X)
	y := f(X)

	for {
		// fibonacci update
		a, b = b, a + b

		// scale epsilon by fibonacci value
		alpha := -b * 1e-3

		// search down parameter space in given direction
		X1 := vector.Add(X0, vector.Scale(alpha, dir))
		y1 := f(X1)

		// update min values or quit
		if y1 < y {
			y = y1
			X = X1
		} else {
			break
		}
	}

	return y, X
}
