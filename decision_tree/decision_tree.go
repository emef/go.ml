package decision_tree

import (
	"fmt"
	"sort"
	"math"
)

type treeNode struct {
	left *treeNode
	right *treeNode
	impurity float64
	splitColumn int
	splitVal float64
	probability float64
}

type giniSplitResult struct {
	gini float64
	splitColumn int
	splitVal float64
}

type zipColumn struct {
	Value float64
	Response float64
}

type zipColumnSortable []zipColumn

func (a zipColumnSortable) Len() int { return len(a) }
func (a zipColumnSortable) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a zipColumnSortable) Less(i, j int) bool { return a[i].Value < a[j].Value }


func (n treeNode) isLeaf() bool {
	return n.left == nil && n.right == nil
}

func (n treeNode) String() string {
	var toString func(*treeNode, string) string
	toString = func(n *treeNode, padding string) string {
		var s string
		if n.isLeaf() {
			s = fmt.Sprintf("%s(%.2f)", padding, n.probability)
		} else {
			s = fmt.Sprintf("%s%d < %.2f  (%.3f)",
				padding,
				n.splitColumn,
				n.splitVal,
				n.impurity)
			if n.left != nil {
				s += "\n" + toString(n.left, padding + "  ")
			}
			if n.right != nil {
				s += "\n" + toString(n.right, padding + "  ")
			}
		}
		return s
	}

	return toString(&n, "")
}

func (tree treeNode) Classify(X [][]float64) []float64 {
	y := make([]float64, len(X))
	for i := range y {
		y[i] = doClassify(&tree, X[i])
	}
	return y
}

func doClassify(tree *treeNode, x []float64) float64 {
	var label float64
	for {
		if tree.isLeaf() {
			if tree.probability > 0.5 {
				label = 1
			} else {
				label = 0
			}
			break
		} else {
			i, val := tree.splitColumn, tree.splitVal
			if x[i] < val {
				tree = tree.left
			} else {
				tree = tree.right
			}
		}
	}
	return label
}

func Fit(X [][]float64, y []float64) *treeNode {
	used := make([]int, len(X[0]))
	return doFit(X, y, used)
}

func doFit(X [][]float64, y []float64, used []int) *treeNode {
	node := new(treeNode)

	// calculate node's probability
	labelSum := 0.0
	for _, v := range y { labelSum += v }
	node.probability = labelSum / float64(len(y))

	// are all columns used?
	colSum := 0
	for _, v := range used { colSum += v }

	should_split := (colSum != len(used)) &&
		node.probability != 0 &&
		node.probability != 1

	if should_split {
		node.impurity, node.splitColumn, node.splitVal = giniSplit(X, y, used)
		couldSplit := node.splitColumn >= 0

		if couldSplit {
			used[node.splitColumn] = 1

			ix := doSplit(X, y, node.splitColumn, node.splitVal)
			node.left = doFit(X[:ix], y[:ix], used)
			node.right = doFit(X[ix:], y[ix:], used)
		}
	}

	return node
}

func doSplit(X [][]float64, y []float64, splitColumn int, splitVal float64) int {
	rearIndex := len(X) - 1
	splitPoint := 0
	for i := 0; i < rearIndex; i++ {
		if X[i][splitColumn] >= splitVal {
			X[i], X[rearIndex] = X[rearIndex], X[i]
			X[i], X[i+1] = X[i+1], X[i]
			y[i], y[rearIndex] = y[rearIndex], y[i]
			y[i], y[i+1] = y[i+1], y[i]
			rearIndex--
		} else {
			splitPoint = i
		}
	}

	return splitPoint + 1
}

func giniSplit(X [][]float64, y []float64, used []int) (float64, int, float64) {
	done := make(chan giniSplitResult)
	nFeatures := len(X[0])
	nChecked := 0
	for i := 0; i < nFeatures; i++ {
		if used[i] == 1 {
			continue
		}

		go func(i int) {
			gini, val := splitPoint(X, i, y)
			done <- giniSplitResult{gini, i, val}
		}(i)

		nChecked += 1
	}

	bestGini, bestCol, bestVal := 1.1, -1, 0.0
	for i := 0; i < nChecked; i++ {
		result := <-done
		if result.gini < bestGini {
			bestGini = result.gini
			bestCol = result.splitColumn
			bestVal = result.splitVal
		}
	}

	return bestGini, bestCol, bestVal
}

// G = p(t_l)*Gini(t_l) + p(t_r)*Gini(t_r)
// Gini(t) = 1 - p(1|t)^2 - p(0|t)^2
func splitPoint(X [][]float64, index int, y []float64) (float64, float64) {
	N := len(X)
	totalPositive := 0.0

	col := make([]zipColumn, N)
	for i := range X {
		col[i] = zipColumn{X[i][index], y[i]}
		totalPositive += y[i]
	}

	sort.Sort(zipColumnSortable(col))
	nPositiveL := 0.0
	minGini := 1.1
	split := col[0].Value

	for i := 1; i < N; i++ {
		if col[i].Value == col[i-1].Value {
			nPositiveL += col[i-1].Response
			continue
		}


		nPositiveL += col[i-1].Response
		pPositiveL := nPositiveL / float64(i)
		pNegativeL := 1 - pPositiveL
		giniL := 1 - math.Pow(pPositiveL, 2) - math.Pow(pNegativeL, 2)

		nPositiveR := totalPositive - nPositiveL
		pPositiveR := nPositiveR / float64(N - i)
		pNegativeR := 1 - pPositiveR
		giniR := 1 - math.Pow(pPositiveR, 2) - math.Pow(pNegativeR, 2)

		pL := float64(i) / float64(N)
		pR := 1 - pL

		gini := pL * giniL + pR * giniR

		if gini < minGini {
			minGini = gini
			split = col[i].Value
		}
	}

	return minGini, split
}