package decision_tree

import (
	"sync"
	"fmt"
	"errors"
)

const GINI = "gini"

type splitFunction func([][]float64, []float64, int) (float64, float64)

type decisionTree struct {
	root    *treeNode     // actual tree
	context *treeContext  // fitting context
}

type treeContext struct {
	splitter splitFunction  // what impurity critera to split by
	curDepth int            // how deep are we
	maxDepth int            // maximum tree depth
	used     []int          // which columns have we used?
}

type treeNode struct {
	left        *treeNode  // left subtree
	right       *treeNode  // right subtree
	impurity    float64    // impurity value
	splitColumn int        // index of column to split by
	splitVal    float64    // value to split on
	probability float64    // P(1|t)
	size        int        // number of samples in this sub tree
}

type splitResult struct {
	impurity    float64  // impurity value from split criteria
	splitColumn int      // best column to split on
	splitVal    float64  // value to split on
}

/*
 decision tree constructor

 arguments
 ---------
   maxDepth:       max depth of constructed tree
   minSamplesLeaf: necessary samples needed to construct a leaf node
   splitMethod:    what criteria to use to calculate impurity
                   possible values: GINI ("gini")
*/
func DecisionTree(maxDepth int, splitMethod string) (*decisionTree, error) {
	var splitter splitFunction

	if splitMethod == "gini" {
		splitter = splitGINI
	} else {
		return nil, errors.New("unknown splitting method")
	}

	tree := new(decisionTree)
	tree.context = new(treeContext)
	tree.context.splitter = splitter
	tree.context.maxDepth = maxDepth

	return tree, nil
}

func (tree decisionTree) String() string {
	return tree.root.String()
}


// fit this decision tree with samples (X) and responses (y)
func (tree *decisionTree) Fit(X [][]float64, y []float64) error {
	tree.context.used = make([]int, len(X))
	tree.root = fitTree(X, y, tree.context)
	return nil
}


// classify samples (X), return predicted labels
func (tree decisionTree) Classify(X [][]float64) []float64 {
	y := make([]float64, len(X))
	for i := range y {
		y[i] = tree.ClassifySample(X[i])
	}
	return y
}


// classify single sample, return predicted label
func (tree decisionTree) ClassifySample(x []float64) float64 {
	var label float64
	node := tree.root

	for {
		if node.isLeaf() {
			if node.probability > 0.5 {
				label = 1
			} else {
				label = 0
			}
			break
		} else {
			i, val := node.splitColumn, node.splitVal
			if x[i] < val {
				node = node.left
			} else {
				node = node.right
			}
		}
	}
	return label
}


func (n treeNode) isLeaf() bool {
	return n.left == nil && n.right == nil
}


func (n treeNode) String() string {
	var toString func(*treeNode, string) string
	toString = func(n *treeNode, padding string) string {
		var s string
		if n.isLeaf() {
			s = fmt.Sprintf("%s(%.3f +%d)", padding, n.probability, n.size)
		} else {
			s = fmt.Sprintf("%s%d < %.2f  (%.3f +%d)",
				padding,
				n.splitColumn,
				n.splitVal,
				n.impurity,
				n.size)
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


/*
 fits a decision tree

 arguments
 ---------
   X:        training samples
   y:        corresponding responses; should be in {0, 1}
   context:  training context: how deep we are, stopping cases, etc.

 returns
 -------
   tree root node
*/
func fitTree(X [][]float64, y []float64, context *treeContext) *treeNode {
	node := new(treeNode)

	// calculate node's probability
	labelSum := 0.0
	for _, v := range y { labelSum += v }
	node.probability = labelSum / float64(len(y))
	node.size = len(X)

	// should we split this tree further?
	// 1) must not exceed maxDepth
	// 2) must have both positive and negative samples
	should_split := context.curDepth < context.maxDepth &&
		node.probability != 0 &&
		node.probability != 1

	if should_split {
		// find best splitting column we haven't used
		result := bestSplit(X, y, context)

		// did we successfully split?
		couldSplit := result.splitColumn != -1

		if couldSplit {
			// this tree's depth, sub trees will have +1
			myDepth := context.curDepth

			// populate the new node's splitting point
			node.impurity = result.impurity
			node.splitColumn = result.splitColumn
			node.splitVal = result.splitVal


			// what index should we split the dataset by?
			ix := splitDataset(X, y, node.splitColumn, node.splitVal)

			// we can no longer use this column
			used := context.used
			used[node.splitColumn] = 1

			// fit sub trees, fix context
			context.curDepth = myDepth + 1
			context.used = copySlice(used)
			node.left = fitTree(X[:ix], y[:ix], context)

			context.curDepth = myDepth + 1
			context.used = copySlice(used)
			node.right = fitTree(X[ix:], y[ix:], context)
		}
	}

	return node
}


/*
 finds the best column/value to split on given current context

 NOTE: uses CPU-bound go-routines, increase runtime.GOMAXPROCS for
 multicore processing and a generous speed-up
*/
func bestSplit(X [][]float64, y []float64, context *treeContext) splitResult {
	nFeatures := len(X[0])
	results := make(chan splitResult, nFeatures)
	wg := new(sync.WaitGroup)

	for i := 0; i < nFeatures; i++ {
		if context.used[i] != 1 {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				impurity, val := context.splitter(X, y, i)
				results <- splitResult{impurity, i, val}
			}(i)
		}
	}

	wg.Wait()
	close(results)

	bestResult := splitResult{1.1, -1, 0.0}
	for result := range results {
		if result.impurity < bestResult.impurity {
			bestResult = result
		}
	}

	return bestResult
}


/*
 Does an in-place split of the dataset & responses according to the
 split column and value.

 returns
 -------
   splitPoint: all samples with index lower than this belong in the
               left sub tree.
*/
func splitDataset(X [][]float64, y []float64, splitColumn int, splitVal float64) int {
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


func copySlice(slice []int) []int {
	newSlice := make([]int, len(slice))
	copy(newSlice, slice)
	return slice
}