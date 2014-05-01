package optimization

import (
	"fmt"
	"testing"
	"github.com/emef/bitfield"
	"math/rand"
)

type BitFieldSample struct {
	field bitfield.BitField
}


func (b BitFieldSample) Mutate() {
	length := uint32(b.field.Size())
	p := 1.0 / float64(length)

	for i := uint32(0); i < length; i++ {
		if rand.Float64() < p {
			b.field.Flip(i)
		}
	}
}

func (b BitFieldSample) Fitness() float64 {
	target := []int{0, 0, 1, 0, 1, 1, 0, 0, 0, 0}
	correct := 0.0
	for i := 0; i < b.field.Size(); i++ {
		if boolToInt(b.field.Test(uint32(i))) == target[i] {
			correct++
		}
	}
	return correct / float64(b.field.Size())
}

func (b BitFieldSample) CrossOver(mateCandidate Sample) Sample {
	mate := mateCandidate.(BitFieldSample)
	length := b.field.Size()
	child := BitFieldSample{bitfield.New(length)}
	i, j := uint32(randInt(length)), uint32(randInt(length))

	if i > j {
		i, j = j, i
	}

	for k := uint32(0); k < uint32(length); k++ {
		if k < i || k > j {
			if b.field.Test(k) {
				child.field.Set(k)
			}
		} else {
			if mate.field.Test(k) {
				child.field.Set(k)
			}
		}
	}

	return child
}

func (b BitFieldSample) String() string {
	s := ""
	for i := uint32(0); i < uint32(b.field.Size()); i++ {
		s += fmt.Sprintf("%d ", boolToInt(b.field.Test(i)))
	}
	return s
}


func randInt(max int) int {
	return int(rand.Float64() * float64(max))
}

func TestIt(t *testing.T) {
	samples := make([]Sample, 20)
	for i := range samples {
		sample := BitFieldSample{bitfield.New(10)}
		indexes := rand.Perm(10)[:randInt(10)]
		for _, j := range indexes {
			sample.field.Set(uint32(j))
		}
		samples[i] = sample
	}
	g := NewGenetic(samples, 0.5, 1e5, 10)
	fmt.Println(g.Run())
}