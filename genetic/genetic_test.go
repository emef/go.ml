package genetic

import (
	"fmt"
	"testing"
	"math/rand"
)

func fitness(x Sample) float64 {
	sample := x.(BitSample)
	target := []int{0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 1, 1, 1}
	correct := 0.0
	for i := 0; i < sample.field.Size(); i++ {
		if boolToInt(sample.field.Test(uint32(i))) == target[i] {
			correct++
		}
	}
	return correct / float64(sample.field.Size())
}

func TestIt(t *testing.T) {
	samples := make([]Sample, 20)
	for i := range samples {
		sample := NewBitSample(20)
		indexes := rand.Perm(20)[:randInt(20)]
		for _, j := range indexes {
			sample.field.Set(uint32(j))
		}
		samples[i] = sample
	}
	g := New(samples, fitness, 0.1, 50, 1.0)
	fmt.Println(g.Run())
}

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
