package genetic

import (
	"github.com/emef/bitfield"
	"math/rand"
	"fmt"
)

type BitSample struct {
	field bitfield.BitField
}

func NewBitSample(size int) BitSample {
	return BitSample{bitfield.New(size)}
}

func (b BitSample) Fitness() float64 {
	// overwrite this!
	return -1
}

func (b BitSample) Mutate() {
	length := uint32(b.field.Size())
	p := 1.0 / float64(length)

	for i := uint32(0); i < length; i++ {
		if rand.Float64() < p {
			b.field.Flip(i)
		}
	}
}

func (b BitSample) CrossOver(mateCandidate Sample) Sample {
	mate := mateCandidate.(BitSample)
	length := b.field.Size()
	child := BitSample{bitfield.New(length)}
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

func (b BitSample) String() string {
	s := ""
	for i := uint32(0); i < uint32(b.field.Size()); i++ {
		s += fmt.Sprintf("%d ", boolToInt(b.field.Test(i)))
	}
	return s
}


func randInt(max int) int {
	return int(rand.Float64() * float64(max))
}
