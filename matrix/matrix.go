package matrix

import (
	"math"
)

func ITranspose(X [][]float64) {
	for i := 0; i < len(X); i++ {
		for j := 0; j < i; j++ {
			X[i][j], X[j][i] = X[j][i], X[i][j]
		}
	}
}

func IScalarMult(X [][]float64, c float64) {
	for i := range X {
		for j := range X[i] {
			X[i][j] *= c
		}
	}
}

// AB
func MatMult(A, B [][]float64) [][]float64 {
	n, m := len(A), len(A[0])
	p := len(B[0])

	if len(B) != m {
		panic("incompatible sizes")
	}

	AB := make([][]float64, n)
	for i := range AB {
		AB[i] = make([]float64, p)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < p; j++ {
			for k := 0; k < m; k++ {
				AB[i][j] += A[i][k] * B[k][j]
			}
		}
	}

	return AB
}

// A'B
func MatMultTrans(A, B [][]float64) [][]float64 {
	m, n := len(A), len(A[0])
	p := len(B[0])

	if len(B) != m {
		panic("incompatible sizes")
	}

	AB := make([][]float64, n)
	for i := range AB {
		AB[i] = make([]float64, p)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < p; j++ {
			for k := 0; k < m; k++ {
				AB[i][j] += A[k][i] * B[k][j]
			}
		}
	}

	return AB
}

func VecMult(A [][]float64, b []float64) []float64 {
	c := make([]float64, len(A))
	for i := range A {
		for j := range b {
			c[i] += A[i][j] * b[j]
		}
	}
	return c
}


func VecMultTrans(A [][]float64, b []float64) []float64 {
	c := make([]float64, len(A[0]))
	for i := range A[0] {
		for j := range b {
			c[i] += A[j][i] * b[j]
		}
	}
	return c
}

func Determinant(X [][]float64) float64 {
	n := len(X)
	if n == 1 {
		return X[0][0]
	}
	if n == 2 {
		return X[0][0] * X[1][1] - X[1][0] * X[0][1];
	}

	det := 0.0
	tmp := make([][]float64, n-1)

	for i := range tmp {
		tmp[i] = make([]float64, n-1)
	}


	for j1 := 0; j1 < n; j1++ {
		for i := 1; i < n; i++ {
			j2 := 0
			for j := 0; j <n; j++ {
				if j == j1 {
					continue
				}
				tmp[i-1][j2] = X[i][j]
				j2++
			}
		}

		det += math.Pow(-1.0, float64(j1+2)) * X[0][j1] * Determinant(tmp)
	}

	return det
}


func CoFactor(A [][]float64) [][]float64 {
	n := len(A)
	B := make([][]float64, n)
	C := make([][]float64, n-1)

	for i := range C {
		B[i] = make([]float64, n)
		C[i] = make([]float64, n-1)
	}
	B[n-1] = make([]float64, n)

	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			i1 := 0
			for ii := 0; ii < n; ii++ {
				if ii == i {
					continue
				}

				j1 := 0
				for jj := 0; jj < n; jj++ {
					if jj == j {
						continue
					}

					C[i1][j1] = A[ii][jj]
					j1++
				}
				i1++
			}

			B[i][j] = math.Pow(-1.0, float64(i+j+2)) * Determinant(C)
		}
	}

	return B
}

func MatDirtyInverse(A [][]float64) [][]float64 {
	m := len(A)
	n := 2*m
	aug := make([][]float64, m)
	for i := range aug {
		aug[i] = make([]float64, n)
		copy(aug[i], A[i])
		aug[i][m+i] = 1
	}

	for k := 0; k < m; k++ {
		// find a pivot for column k
		iMax, iMaxVal := 0, 0.0
		for i := k ; i < m; i++ {
			val := math.Abs(aug[i][k])
			if val > iMaxVal {
				iMax, iMaxVal = i, val
			}
		}

		// dirty, this matrix is singular
		if iMaxVal == 0 {
			//panic("matrix is singular")
			continue
		}

		// swap rows k and iMax
		aug[k], aug[iMax] = aug[iMax], aug[k]

		// adjust all rows below the pivot
		for i := k + 1; i < m; i++ {
			scale := (aug[i][k] / aug[k][k])
			for j := k; j < n; j++ {
				aug[i][j] = aug[i][j] - aug[k][j] * scale
			}

			// fill in zeros
			aug[i][k] = 0
		}
	}

	// now in row-eschelon form, need to backsubstitute
	for k := m-1; k >= 0; k-- {
		// set eq to 1
		val := aug[k][k]

		// dirty, this matrix is singular
		if val == 0 {
			continue
		}

		for j := k; j < n; j++ {
			aug[k][j] /= val
		}

		for i := k-1; i >= 0; i-- {
			mult := aug[i][k]
			for j := k; j < n; j++ {
				aug[i][j] -= aug[k][j] * mult
			}
		}
	}

	// return only the inverse
	for i := range aug {
		aug[i] = aug[i][m:]
	}

	return aug
}