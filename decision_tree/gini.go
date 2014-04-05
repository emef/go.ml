package decision_tree

import (
	"sort"
	"math"
)

// zips two values and allows sorting based on value
type zipColumn struct {
	Value float64
	Response float64
}

// container of zipColumns that implements sorting interface
type zipColumnSortable []zipColumn

func (a zipColumnSortable) Len() int { return len(a) }
func (a zipColumnSortable) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a zipColumnSortable) Less(i, j int) bool { return a[i].Value < a[j].Value }


/*
 Split dataset on column according to GINI criteria.

 GINI calculated as:
   P(t) = probability of a sample belonging to subtree
   P(1|t) = probability of pos sample in subtree
   G(t) = 1 - P(1|t)^2 - P(0|t)^2
   GINI = P(t_l)*G(t_l) + P(t_r)*G(t_r)
     where t_l and t_r are left and right subtrees after some split

 Arguments
 --------
   X:     training dataset
   y:     training responses
   index: column index to split by

 Returns
 -------
   (minGini, split)
   where
     minGini: minimum Gini impurity value
     split:   value at which we split this column
*/
func splitGINI(X [][]float64, y []float64, index int) (float64, float64) {
	N := len(X)
	totalPositive := 0.0

	// zip the column values and corresponding responses
	col := make([]zipColumn, N)
	for i := range X {
		col[i] = zipColumn{X[i][index], y[i]}
		totalPositive += y[i]
	}

	// sort the column by value
	sort.Sort(zipColumnSortable(col))

	nPositiveL := 0.0
	minGini := 1.1
	split := col[0].Value

	// investigate every possible split point, track minimum GINI
	for i := 1; i < N; i++ {
		// we already split on this value! continue
		if col[i].Value == col[i-1].Value {
			nPositiveL += col[i-1].Response
			continue
		}

		// calculate P(1|t_l) and P(0|t_l)
		nPositiveL += col[i-1].Response
		pPositiveL := nPositiveL / float64(i)
		pNegativeL := 1 - pPositiveL

		// infer P(1|t_r) and P(0|t_r) from left tree
		nPositiveR := totalPositive - nPositiveL
		pPositiveR := nPositiveR / float64(N - i)
		pNegativeR := 1 - pPositiveR

		// calculate G(t_l) and G(t_r)
		giniL := 1 - math.Pow(pPositiveL, 2) - math.Pow(pNegativeL, 2)
		giniR := 1 - math.Pow(pPositiveR, 2) - math.Pow(pNegativeR, 2)

		// calculate P(t_l) and P(t_r)
		pL := float64(i) / float64(N)
		pR := 1 - pL

		// put everything together to calculate GINI for this split point
		gini := pL * giniL + pR * giniR

		// record the best split so far
		if gini < minGini {
			minGini = gini
			split = col[i].Value
		}
	}

	return minGini, split
}