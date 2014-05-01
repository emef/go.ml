package genetic

import (
	"sort"
	"math/rand"
	"fmt"
)

type Sample interface {
	Mutate()
	CrossOver(Sample) Sample
}

type Samples []Sample

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

type Genetic struct {
	candidates Candidates
	fitness func(Sample) float64
	mutateRate float64
	maxGenerations int
	minFitness float64
}

func (g *Genetic) Run() Sample {
	for i := 0; i < g.maxGenerations; i++ {
		best := g.BestCandidate()
		fmt.Printf("EPOCH %d: BEST FITNESS %s\n", i, best)
		if best.fitness >= g.minFitness {
			break
		} else {
			g.runEpoch()
		}
	}
	return g.BestCandidate().sample
}

func (g *Genetic) runEpoch() {
	nextGen := make([]Candidate, len(g.candidates))
	parents := make([]Candidate, len(g.candidates) * 2)

	for i := range parents {
		minAccFitness := rand.Float64()
		accFitness := 0.0
		for _, candidate := range g.candidates {
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
		if rand.Float64() < g.mutateRate {
			newSample.Mutate()
		}
		child := Candidate{newSample, g.fitness(newSample)}
		sumFitness += child.fitness
		nextGen[j] = child
		j++
	}

	for _, candidate := range nextGen {
		candidate.fitness /= sumFitness
	}

	sort.Sort(sort.Reverse(Candidates(nextGen)))
	g.candidates = nextGen
}


func (g Genetic) BestCandidate() Candidate {
	return g.candidates[0]
}

func New(
	samples Samples,
	fitnessFunc func(Sample) float64,
	mutateRate float64,
	maxGenerations int,
	minFitness float64) Genetic {

	// ...awkward space...

	candidates := make([]Candidate, len(samples))
	sumFitness := 0.0
	for i, sample := range samples {
		fitness := fitnessFunc(sample)
		candidates[i] = Candidate{sample, fitness}
		sumFitness += fitness
	}
	for _, candidate := range candidates {
		candidate.fitness /= sumFitness
	}

	sort.Sort(sort.Reverse(Candidates(candidates)))

	return Genetic{
		candidates,
		fitnessFunc,
		mutateRate,
		maxGenerations,
		minFitness}
}