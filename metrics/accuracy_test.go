package metrics

import "testing"

func testAccuracy(t *testing.T) {
	yPred := []float64{0,0,1,1,1}
	yTrue := []float64{0,1,0,1,1}
	acc := Accuracy(yPred, yTrue)
	if acc != 0.6 {
		t.Errorf("Accuracy incorrect (%.2f != %.2f)", acc, 0.6)
	}
}