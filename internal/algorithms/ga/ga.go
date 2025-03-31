package ga

import (
	"fmt"
	"golang-moaha-construction/internal/objectives/single"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
	"golang-moaha-construction/internal/util"
)

const (
	PopulationSizeParam = "Population Size"
	MaxIterationsParam  = "Number of Iterations"
	CrossoverRateParam  = "Crossover Rate"
	MutationRateParam   = "Mutation Rate"
	ElitismCountParam   = "Elitism Count"
)

type GAAlgorithm struct {
	PopulationSize    int
	MaxIterations     int
	CrossoverRate     float64
	MutationRate      float64
	ElitismCount      int
	Population        []*single.SingleResult
	Convergence       []float64
	ObjectiveFunction single.SingleProblem
	Best              *single.SingleResult
}

var Configs = []data.Config{
	{
		Name:               PopulationSizeParam,
		ValidationFunction: util.IsValidPositiveInteger,
	},
	{
		Name:               MaxIterationsParam,
		ValidationFunction: util.IsValidPositiveInteger,
	},
	{
		Name:               CrossoverRateParam,
		ValidationFunction: util.IsValidFloat,
	},
	{
		Name:               MutationRateParam,
		ValidationFunction: util.IsValidFloat,
	},
	{
		Name:               ElitismCountParam,
		ValidationFunction: util.IsValidPositiveInteger,
	},
}

func Create(problem single.SingleProblem, configs []*data.Config) (*GAAlgorithm, error) {
	var popSize, maxIters, elitismCount int
	var crossoverRate, mutationRate float64

	for _, config := range configs {
		val := strings.Trim(config.Value, " ")
		switch config.Name {
		case PopulationSizeParam:
			popSize, _ = strconv.Atoi(val)
		case MaxIterationsParam:
			maxIters, _ = strconv.Atoi(val)
		case CrossoverRateParam:
			crossoverRate, _ = strconv.ParseFloat(val, 64)
		case MutationRateParam:
			mutationRate, _ = strconv.ParseFloat(val, 64)
		case ElitismCountParam:
			elitismCount, _ = strconv.Atoi(val)
		}
	}

	convergence := make([]float64, maxIters)
	population := make([]*single.SingleResult, popSize)

	// This implementation supports only one objective.
	if problem.NumberOfObjectives() != 1 {
		return nil, objectives.ErrInvalidNumberOfObjectives
	}

	return &GAAlgorithm{
		PopulationSize:    popSize,
		MaxIterations:     maxIters,
		CrossoverRate:     crossoverRate,
		MutationRate:      mutationRate,
		ElitismCount:      elitismCount,
		Convergence:       convergence,
		ObjectiveFunction: problem,
		Population:        population,
	}, nil
}

func (ga *GAAlgorithm) Type() data.TypeProblem {
	return data.Single
}

func (ga *GAAlgorithm) Run() error {
	ga.initialization()

	bar := progressbar.Default(int64(ga.MaxIterations))
	var wg sync.WaitGroup

	dim := ga.ObjectiveFunction.GetDimension()
	lowerBound := ga.ObjectiveFunction.GetLowerBound()
	upperBound := ga.ObjectiveFunction.GetUpperBound()

	// Main GA loop (each iteration represents a generation)
	for iter := 0; iter < ga.MaxIterations; iter++ {
		newPopulation := make([]*single.SingleResult, ga.PopulationSize)

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

				child := &single.SingleResult{
					Idx:      idx,
					Position: childPos,
				}
				value, _, _ := ga.ObjectiveFunction.Eval(childPos)
				child.Value = value
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

			newGene := &single.SingleResult{
				Idx:      idx,
				Position: pos,
			}

			value, _, _ := ga.ObjectiveFunction.Eval(pos)
			newGene.Value = value

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

func (ga *GAAlgorithm) sortPopulationByFitness() []*single.SingleResult {
	sorted := make([]*single.SingleResult, len(ga.Population))
	copy(sorted, ga.Population)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value[0] < sorted[j].Value[0]
	})
	return sorted
}

// tournamentSelection selects the best individual among a random sample.
func tournamentSelection(pop []*single.SingleResult, tournamentSize int) *single.SingleResult {
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

// outOfBoundaries ensures that each component of the position is within the allowed bounds.
func outOfBoundaries(position, lowerBound, upperBound []float64) {
	for i := range position {
		if position[i] < lowerBound[i] {
			position[i] = lowerBound[i]
		} else if position[i] > upperBound[i] {
			position[i] = upperBound[i]
		}
	}
}
