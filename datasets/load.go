package datasets

import (
	"os"
	"strconv"
	"runtime"
	"path"
	"encoding/csv"
)

func Load(name string) ([][]float64, []float64) {
	switch {
	case name == "iris":
		return loadIris()
	case name == "cancer":
		return loadBreastCancer()
	}

	panic("unknown dataset")
}


func absPath(filename string) string {
	_, thisFile, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(thisFile), filename)
}


func loadIris() ([][]float64, []float64) {
	X_file, err1 := os.Open(absPath("iris_X.csv"))
	y_file, err2 := os.Open(absPath("iris_y.csv"))

	if err1 != nil || err2 != nil {
		panic("couldn't open iris files")
	}

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
	file, err := os.Open(absPath("breast_cancer.csv"))

	if err != nil {
		panic("couldn't open csv")
	}

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