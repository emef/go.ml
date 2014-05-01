package optimization

import (
	"sort"
	"math/rand"
	"fmt"
)

type Sample interface {
	Mutate()
	Fitness() float64
	CrossOver(Sample) Sample
}

type Candidate struct {
	sample Sample
	fitness float64
}

func (c Candidate) String() string {
	return fmt.Sprintf("%s (%.2f)", c.sample, c.fitness)
}

type Candidates []Candidate

func (cs Candidates) Len() int { return len(cs) }
func (cs Candidates) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }
func (cs Candidates) Less(i, j int) bool {
	return cs[i].fitness < cs[j].fitness
}

func (cs Candidates) nextGeneration(mutateRate float64) Candidates {
	nextGen := make([]Candidate, len(cs))
	parents := make([]Candidate, len(cs) * 2)

	for i := range parents {
		minAccFitness := rand.Float64()
		accFitness := 0.0
		for _, candidate := range cs {
			accFitness += candidate.fitness
			if accFitness >= minAccFitness {
				parents[i] = candidate
				break
			}
		}
	}

	j := 0
	sumFitness := 0.0
	indexes := rand.Perm(len(parents))
	for i := 0; i < len(indexes) - 1; i += 2 {
		a, b := parents[indexes[i]], parents[indexes[i+1]]
		newSample := a.sample.CrossOver(b.sample)
		if rand.Float64() < mutateRate {
			newSample.Mutate()
		}
		child := Candidate{newSample, newSample.Fitness()}
		sumFitness += child.fitness
		nextGen[j] = child
		j++
	}

	for _, candidate := range nextGen {
		candidate.fitness /= sumFitness
	}

	sort.Sort(sort.Reverse(Candidates(nextGen)))
	return nextGen

}

type Genetic struct {
	candidates Candidates
	mutateRate float64
	maxGenerations int
	minFitness float64
}

func (g Genetic) Run() Sample {
	for i := 0; i < g.maxGenerations; i++ {
		best := g.BestCandidate()
		fmt.Printf("EPOCH %d: BEST FITNESS %.3f\n", i, best.fitness)
		if best.fitness >= g.minFitness {
			break
		} else {
			g.candidates = g.candidates.nextGeneration(g.mutateRate)
		}
	}
	return g.BestCandidate().sample
}

func (g Genetic) BestCandidate() Candidate {
	return g.candidates[0]
}

func NewGenetic(samples []Sample, mutateRate float64, maxGenerations int, minFitness float64) Genetic {
	candidates := make([]Candidate, len(samples))
	sumFitness := 0.0
	for i, sample := range samples {
		fitness := sample.Fitness()
		candidates[i] = Candidate{sample, fitness}
		sumFitness += fitness
	}
	for _, candidate := range candidates {
		candidate.fitness /= sumFitness
	}

	sort.Sort(sort.Reverse(Candidates(candidates)))

	return Genetic{candidates, mutateRate, maxGenerations, minFitness}
}

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
