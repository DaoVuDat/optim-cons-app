package ga

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"golang-moaha-construction/internal/algorithms"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
	"math"
	"math/rand"
	"sort"
	"sync"
)

const NameType algorithms.AlgorithmType = "GA"

type GAAlgorithm struct {
	PopulationSize    int
	MaxIterations     int
	CrossoverRate     float64
	MutationRate      float64
	ElitismCount      int
	Population        []*objectives.Result
	Convergence       []float64
	ObjectiveFunction objectives.Problem
	Best              *objectives.Result
}

type Config struct {
	Chromosome    int
	Generation    int
	CrossoverRate float64
	MutationRate  float64
	ElitismCount  int
}

func Create(problem objectives.Problem, configs Config) (*GAAlgorithm, error) {

	// This implementation supports only one objective.
	if problem.NumberOfObjectives() != 1 {
		return nil, objectives.ErrInvalidNumberOfObjectives
	}

	return &GAAlgorithm{
		PopulationSize:    configs.Chromosome,
		MaxIterations:     configs.Generation,
		CrossoverRate:     configs.CrossoverRate,
		MutationRate:      configs.MutationRate,
		ElitismCount:      configs.ElitismCount,
		ObjectiveFunction: problem,
	}, nil
}

func (ga *GAAlgorithm) reset() {
	ga.Convergence = make([]float64, ga.MaxIterations)
	ga.Population = make([]*objectives.Result, ga.PopulationSize)
}

func (ga *GAAlgorithm) Type() data.TypeProblem {
	return data.Single
}

func (ga *GAAlgorithm) Run() error {
	ga.reset()

	ga.initialization()

	bar := progressbar.Default(int64(ga.MaxIterations))
	var wg sync.WaitGroup

	dim := ga.ObjectiveFunction.GetDimension()
	lowerBound := ga.ObjectiveFunction.GetLowerBound()
	upperBound := ga.ObjectiveFunction.GetUpperBound()

	// Main GA loop (each iteration represents a generation)
	for iter := 0; iter < ga.MaxIterations; iter++ {
		newPopulation := make([]*objectives.Result, ga.PopulationSize)

		// Elitism: preserve the best individuals.
		sortedPop := ga.sortPopulationByFitness()
		for i := 0; i < ga.ElitismCount; i++ {
			newPopulation[i] = sortedPop[i].CopyAgent()
		}

		// Generate offspring for the rest of the population.
		wg.Add(ga.PopulationSize - ga.ElitismCount)
		for i := ga.ElitismCount; i < ga.PopulationSize; i++ {
			go func(idx int) {
				defer wg.Done()
				// Tournament selection (tournament size = 3)
				parent1 := tournamentSelection(ga.Population, 3)
				parent2 := tournamentSelection(ga.Population, 3)

				var childPos []float64
				// Crossover (using blend crossover with alpha = 0.3)
				if rand.Float64() < ga.CrossoverRate {
					childPos = blendCrossover(parent1.Position, parent2.Position, lowerBound, upperBound, 0.3)
				} else {
					// If no crossover, copy parent1
					childPos = make([]float64, dim)
					copy(childPos, parent1.Position)
				}

				// Mutation (Gaussian mutation with 10% of the range as sigma)
				childPos = gaussianMutation(childPos, lowerBound, upperBound, ga.MutationRate)
				// Ensure the child is within boundaries.
				outOfBoundaries(childPos, lowerBound, upperBound)
				// Evaluate the child solution.

				child := &objectives.Result{
					Idx:      idx,
					Position: childPos,
				}
				value, valuesWithKey, keys, penalty := ga.ObjectiveFunction.Eval(childPos)
				child.Value = value
				child.ValuesWithKey = valuesWithKey
				child.Penalty = penalty
				child.Key = keys

				newPopulation[idx] = child
			}(i)
		}
		wg.Wait()

		ga.Population = newPopulation
		ga.findBest()

		ga.Convergence[iter] = ga.Best.Value[0]
		bar.Describe(fmt.Sprintf("Iter %d: %e", iter+1, ga.Best.Value[0]))
		bar.Add(1)
	}

	return nil
}

func (ga *GAAlgorithm) RunWithChannel(doneChan chan<- struct{}, channel chan<- any) error {
	ga.reset()

	ga.initialization()

	var wg sync.WaitGroup

	dim := ga.ObjectiveFunction.GetDimension()
	lowerBound := ga.ObjectiveFunction.GetLowerBound()
	upperBound := ga.ObjectiveFunction.GetUpperBound()

	// Main GA loop (each iteration represents a generation)
	for iter := 0; iter < ga.MaxIterations; iter++ {
		newPopulation := make([]*objectives.Result, ga.PopulationSize)

		// Elitism: preserve the best individuals.
		sortedPop := ga.sortPopulationByFitness()
		for i := 0; i < ga.ElitismCount; i++ {
			newPopulation[i] = sortedPop[i].CopyAgent()
		}

		// Generate offspring for the rest of the population.
		wg.Add(ga.PopulationSize - ga.ElitismCount)
		for i := ga.ElitismCount; i < ga.PopulationSize; i++ {
			go func(idx int) {
				defer wg.Done()
				// Tournament selection (tournament size = 3)
				parent1 := tournamentSelection(ga.Population, 3)
				parent2 := tournamentSelection(ga.Population, 3)

				var childPos []float64
				// Crossover (using blend crossover with alpha = 0.3)
				if rand.Float64() < ga.CrossoverRate {
					childPos = blendCrossover(parent1.Position, parent2.Position, lowerBound, upperBound, 0.3)
				} else {
					// If no crossover, copy parent1
					childPos = make([]float64, dim)
					copy(childPos, parent1.Position)
				}

				// Mutation (Gaussian mutation with 10% of the range as sigma)
				childPos = gaussianMutation(childPos, lowerBound, upperBound, ga.MutationRate)
				// Ensure the child is within boundaries.
				outOfBoundaries(childPos, lowerBound, upperBound)
				// Evaluate the child solution.

				child := &objectives.Result{
					Idx:      idx,
					Position: childPos,
				}
				value, valuesWithKey, keys, penalty := ga.ObjectiveFunction.Eval(childPos)
				child.Value = value
				child.Penalty = penalty
				child.Key = keys
				child.ValuesWithKey = valuesWithKey

				newPopulation[idx] = child
			}(i)
		}
		wg.Wait()

		ga.Population = newPopulation
		ga.findBest()

		ga.Convergence[iter] = ga.Best.Value[0]

		channel <- struct {
			Progress    float64 `json:"progress"`
			BestFitness float64 `json:"bestFitness"`
			Type        string  `json:"type"`
		}{
			Progress:    (float64(iter+1) / float64(ga.MaxIterations)) * 100,
			BestFitness: ga.Best.Value[0],
			Type:        "single",
		}

	}

	close(channel)

	return nil
}

func (ga *GAAlgorithm) initialization() {
	dim := ga.ObjectiveFunction.GetDimension()
	lowerBound := ga.ObjectiveFunction.GetLowerBound()
	upperBound := ga.ObjectiveFunction.GetUpperBound()
	var wg sync.WaitGroup
	wg.Add(ga.PopulationSize)
	for i := 0; i < ga.PopulationSize; i++ {
		go func(idx int) {
			defer wg.Done()
			pos := make([]float64, dim)
			for d := 0; d < dim; d++ {
				pos[d] = lowerBound[d] + rand.Float64()*(upperBound[d]-lowerBound[d])
			}

			newGene := &objectives.Result{
				Idx:      idx,
				Position: pos,
			}

			value, valuesWithKey, keys, penalty := ga.ObjectiveFunction.Eval(pos)
			newGene.Value = value
			newGene.Penalty = penalty
			newGene.Key = keys
			newGene.ValuesWithKey = valuesWithKey

			ga.Population[idx] = newGene
		}(i)
	}
	wg.Wait()
	ga.findBest()
}

func (ga *GAAlgorithm) findBest() {
	best := ga.Population[0]
	for i := 1; i < len(ga.Population); i++ {
		if ga.Population[i].Value[0] < best.Value[0] {
			best = ga.Population[i]
		}
	}
	ga.Best = best.CopyAgent()
}

func (ga *GAAlgorithm) sortPopulationByFitness() []*objectives.Result {
	sorted := make([]*objectives.Result, len(ga.Population))
	copy(sorted, ga.Population)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value[0] < sorted[j].Value[0]
	})
	return sorted
}

func (ga *GAAlgorithm) GetResults() algorithms.Result {
	results := make([]algorithms.AlgorithmResult, 1)

	mapLoc, sliceLoc, cranes, err := ga.ObjectiveFunction.GetLocationResult(ga.Best.Position)

	if err != nil {
		return algorithms.Result{}
	}
	results[0] = algorithms.AlgorithmResult{
		MapLocations:   mapLoc,
		SliceLocations: sliceLoc,
		Value:          ga.Best.Value,
		Penalty:        ga.Best.Penalty,
		Key:            ga.Best.Key,
		Cranes:         cranes,
		Phases:         ga.ObjectiveFunction.GetPhases(),
		ValuesWithKey:  ga.Best.ValuesWithKey,
	}

	minX, maxX, minY, maxY, _ := ga.ObjectiveFunction.GetLayoutSize()

	return algorithms.Result{
		Result:      results,
		MinX:        minX,
		MaxX:        maxX,
		MinY:        minY,
		MaxY:        maxY,
		Convergence: ga.Convergence,
	}
}

// tournamentSelection selects the best individual among a random sample.
func tournamentSelection(pop []*objectives.Result, tournamentSize int) *objectives.Result {
	popSize := len(pop)
	best := pop[rand.Intn(popSize)]
	for i := 1; i < tournamentSize; i++ {
		candidate := pop[rand.Intn(popSize)]
		if candidate.Value[0] < best.Value[0] {
			best = candidate
		}
	}
	return best
}

// blendCrossover implements BLX-alpha crossover for continuous variables.
func blendCrossover(parent1, parent2, lowerBound, upperBound []float64, alpha float64) []float64 {
	dim := len(parent1)
	child := make([]float64, dim)
	for i := 0; i < dim; i++ {
		minVal := math.Min(parent1[i], parent2[i])
		maxVal := math.Max(parent1[i], parent2[i])
		rangeVal := maxVal - minVal

		extLow := minVal - alpha*rangeVal
		extUp := maxVal + alpha*rangeVal

		// Ensure the extended bounds do not exceed the actual boundaries.
		if extLow < lowerBound[i] {
			extLow = lowerBound[i]
		}
		if extUp > upperBound[i] {
			extUp = upperBound[i]
		}

		child[i] = extLow + rand.Float64()*(extUp-extLow)
	}
	return child
}

// gaussianMutation applies Gaussian mutation with standard deviation equal to 10% of the variable range.
func gaussianMutation(child, lowerBound, upperBound []float64, mutationRate float64) []float64 {
	dim := len(child)
	mutated := make([]float64, dim)
	copy(mutated, child)
	for i := 0; i < dim; i++ {
		if rand.Float64() < mutationRate {
			sigma := 0.1 * (upperBound[i] - lowerBound[i])
			mutated[i] += rand.NormFloat64() * sigma
		}
	}
	return mutated
}

// outOfBoundaries ensures that each components of the position is within the allowed bounds.
func outOfBoundaries(position, lowerBound, upperBound []float64) {
	for i := range position {
		if position[i] < lowerBound[i] {
			position[i] = lowerBound[i]
		} else if position[i] > upperBound[i] {
			position[i] = upperBound[i]
		}
	}
}
