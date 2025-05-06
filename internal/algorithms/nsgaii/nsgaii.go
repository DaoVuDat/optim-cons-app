package nsgaii

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"golang-moaha-construction/internal/algorithms"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
	"golang-moaha-construction/internal/util"
	"math"
	"math/rand"
	"sort"
	"sync"
)

const NameType algorithms.AlgorithmType = "NSGA-II"

type NSGAIIAlgorithm struct {
	PopulationSize    int
	MaxIterations     int
	CrossoverRate     float64
	MutationRate      float64
	MutationStrength  float64
	Sigma             float64
	Population        []*objectives.Result
	Archive           []*objectives.Result
	ObjectiveFunction objectives.Problem
}

type Config struct {
	PopulationSize   int
	MaxIterations    int
	CrossoverRate    float64
	MutationRate     float64
	MutationStrength float64
	Sigma            float64
}

func Create(problem objectives.Problem, configs Config) (*NSGAIIAlgorithm, error) {
	// NSGA-II is designed for multi-objective optimization
	if problem.NumberOfObjectives() < 2 {
		return nil, objectives.ErrInvalidNumberOfObjectives
	}

	return &NSGAIIAlgorithm{
		PopulationSize:    configs.PopulationSize,
		MaxIterations:     configs.MaxIterations,
		CrossoverRate:     configs.CrossoverRate,
		MutationRate:      configs.MutationRate,
		MutationStrength:  configs.MutationStrength,
		Sigma:             configs.Sigma,
		Archive:           make([]*objectives.Result, 0, configs.PopulationSize),
		ObjectiveFunction: problem,
	}, nil
}

func (ga *NSGAIIAlgorithm) reset() {
	ga.Population = make([]*objectives.Result, ga.PopulationSize)
	ga.Archive = make([]*objectives.Result, 0, ga.PopulationSize)
}

func (ga *NSGAIIAlgorithm) Type() data.TypeProblem {
	return data.Multi
}

func (ga *NSGAIIAlgorithm) Run() error {
	ga.reset()

	ga.initialization()

	bar := progressbar.Default(int64(ga.MaxIterations))

	var wg sync.WaitGroup

	// Parameters for NSGA-II
	pc := ga.CrossoverRate // Probability of crossover
	pm := ga.MutationRate  // Probability of mutation

	nCrossover := int(2 * math.Round(pc*float64(ga.PopulationSize)/2))
	nMutation := int(math.Round(pm * float64(ga.PopulationSize)))

	// Non-Dominated Sorting
	var paretoFront [][]int
	ga.Population, paretoFront = objectives.NonDominatedSort(ga.Population)

	// Calculate Crowding Distance
	ga.Population = ga.calculateCrowdingDistance(ga.Population, paretoFront)

	// Sort Population
	ga.Population, paretoFront = SortPopulation(ga.Population)

	sigma := make([]float64, len(ga.ObjectiveFunction.GetLowerBound()))
	for i := 0; i < len(sigma); i++ {
		sigma[i] = ga.Sigma * (ga.ObjectiveFunction.GetUpperBound()[i] - ga.ObjectiveFunction.GetLowerBound()[i])
	}

	// Main NSGA-II loop
	for iter := 0; iter < ga.MaxIterations; iter++ {
		// crossover
		popC := make([]*objectives.Result, nCrossover)
		wg.Add(nCrossover / 2)
		for i := 0; i < nCrossover/2; i++ {
			go func(idx int) {
				defer wg.Done()
				i1 := rand.Intn(len(ga.Population))
				i2 := rand.Intn(len(ga.Population))

				p1 := ga.Population[i1]
				p2 := ga.Population[i2]

				child1, child2 := crossOver(p1, p2)

				popC[idx*2] = child1
				value, valuesWithKey, keys, penalty := ga.ObjectiveFunction.Eval(popC[idx*2].Position)
				popC[idx*2].Value = value
				popC[idx*2].ValuesWithKey = valuesWithKey
				popC[idx*2].Penalty = penalty
				popC[idx*2].Key = keys

				popC[(idx*2)+1] = child2
				value, valuesWithKey, keys, penalty = ga.ObjectiveFunction.Eval(popC[(idx*2)+1].Position)
				popC[(idx*2)+1].Value = value
				popC[(idx*2)+1].ValuesWithKey = valuesWithKey
				popC[(idx*2)+1].Penalty = penalty
				popC[(idx*2)+1].Key = keys

			}(i)
		}
		wg.Wait()

		// mutation
		popM := make([]*objectives.Result, nMutation)
		wg.Add(nMutation)
		for i := 0; i < nMutation; i++ {
			go func(idx int) {
				defer wg.Done()
				i := rand.Intn(len(ga.Population))
				p := ga.Population[i]

				c := mutation(p, ga.MutationStrength, sigma)
				popM[idx] = c
				value, valuesWithKey, keys, penalty := ga.ObjectiveFunction.Eval(popM[idx].Position)
				popM[idx].Value = value
				popM[idx].ValuesWithKey = valuesWithKey
				popM[idx].Penalty = penalty
				popM[idx].Key = keys
			}(i)
		}
		wg.Wait()

		// Merge parent and offspring
		newPop := objectives.MergeAgents(objectives.MergeAgents(ga.Population, popC), popM)
		newPop, paretoFront = objectives.FastNonDominatedSorting_Vectorized(newPop)
		newPop = ga.calculateCrowdingDistance(newPop, paretoFront)
		newPop, paretoFront = SortPopulation(newPop)

		// Truncate
		trunkedPop := make([]*objectives.Result, ga.PopulationSize)
		for i := 0; i < ga.PopulationSize; i++ {
			trunkedPop[i] = newPop[i].CopyAgent()
			trunkedPop[i].Idx = i
		}

		trunkedPop, paretoFront = objectives.FastNonDominatedSorting_Vectorized(trunkedPop)
		trunkedPop = ga.calculateCrowdingDistance(trunkedPop, paretoFront)
		trunkedPop, paretoFront = SortPopulation(trunkedPop)

		newResults := make([]*objectives.Result, len(paretoFront[0]))
		for i := 0; i < len(paretoFront[0]); i++ {
			newResults[i] = trunkedPop[i].CopyAgent()
		}

		ga.Population = trunkedPop
		ga.Archive = newResults

		// Update progress bar
		bar.Describe(fmt.Sprintf("Iter %d: Archive size %d", iter+1, len(ga.Archive)))
		bar.Add(1)
	}

	return nil
}

func crossOver(p1, p2 *objectives.Result) (*objectives.Result, *objectives.Result) {

	child1 := &objectives.Result{
		Position: make([]float64, len(p1.Position)),
	}
	child2 := &objectives.Result{
		Position: make([]float64, len(p1.Position)),
	}

	for i := 0; i < len(p1.Position); i++ {
		alpha := rand.Float64()
		child1.Position[i] = p1.Position[i]*alpha + (1-alpha)*(p2.Position[i])
		child2.Position[i] = p2.Position[i]*alpha + (1-alpha)*(p1.Position[i])
	}

	return child1, child2
}

func mutation(p *objectives.Result, mu float64, sigma []float64) *objectives.Result {
	child := &objectives.Result{
		Position: make([]float64, len(p.Position)),
	}

	copy(child.Position, p.Position)

	nVar := len(p.Position)
	nMutations := int(math.Ceil(mu * float64(nVar)))

	if nMutations == 0 {
		return child
	}
	indices := util.RandomSample(nVar, nMutations)

	for _, idx := range indices {
		gaussianNoise := rand.NormFloat64()

		child.Position[idx] = p.Position[idx] + sigma[idx]*gaussianNoise
	}

	return child
}

func (ga *NSGAIIAlgorithm) RunWithChannel(doneChan chan<- struct{}, channel chan<- any) error {
	ga.reset()

	ga.initialization()

	var wg sync.WaitGroup

	// Parameters for NSGA-II
	pc := ga.CrossoverRate // Probability of crossover
	pm := ga.MutationRate  // Probability of mutation

	nCrossover := int(2 * math.Round(pc*float64(ga.PopulationSize)/2))
	nMutation := int(math.Round(pm * float64(ga.PopulationSize)))

	// Non-Dominated Sorting
	var paretoFront [][]int
	ga.Population, paretoFront = objectives.FastNonDominatedSorting_Vectorized(ga.Population)

	// Calculate Crowding Distance
	ga.Population = ga.calculateCrowdingDistance(ga.Population, paretoFront)

	// Sort Population
	ga.Population, paretoFront = SortPopulation(ga.Population)

	sigma := make([]float64, len(ga.ObjectiveFunction.GetLowerBound()))
	for i := 0; i < len(sigma); i++ {
		sigma[i] = ga.Sigma * (ga.ObjectiveFunction.GetUpperBound()[i] - ga.ObjectiveFunction.GetLowerBound()[i])
	}

	// Main NSGA-II loop
	for iter := 0; iter < ga.MaxIterations; iter++ {
		// crossover
		popC := make([]*objectives.Result, nCrossover)
		wg.Add(nCrossover / 2)
		for i := 0; i < nCrossover/2; i++ {
			go func(idx int) {
				defer wg.Done()
				i1 := rand.Intn(len(ga.Population))
				i2 := rand.Intn(len(ga.Population))

				p1 := ga.Population[i1]
				p2 := ga.Population[i2]

				child1, child2 := crossOver(p1, p2)

				popC[idx*2] = child1
				value, valuesWithKey, keys, penalty := ga.ObjectiveFunction.Eval(popC[idx*2].Position)
				popC[idx*2].Value = value
				popC[idx*2].ValuesWithKey = valuesWithKey
				popC[idx*2].Penalty = penalty
				popC[idx*2].Key = keys

				popC[(idx*2)+1] = child2
				value, valuesWithKey, keys, penalty = ga.ObjectiveFunction.Eval(popC[(idx*2)+1].Position)
				popC[(idx*2)+1].Value = value
				popC[(idx*2)+1].ValuesWithKey = valuesWithKey
				popC[(idx*2)+1].Penalty = penalty
				popC[(idx*2)+1].Key = keys

			}(i)
		}
		wg.Wait()

		// mutation
		popM := make([]*objectives.Result, nMutation)
		wg.Add(nMutation)
		for i := 0; i < nMutation; i++ {
			go func(idx int) {
				defer wg.Done()
				i := rand.Intn(len(ga.Population))
				p := ga.Population[i]

				c := mutation(p, ga.MutationStrength, sigma)
				popM[idx] = c
				value, valuesWithKey, keys, penalty := ga.ObjectiveFunction.Eval(popM[idx].Position)
				popM[idx].Value = value
				popM[idx].ValuesWithKey = valuesWithKey
				popM[idx].Penalty = penalty
				popM[idx].Key = keys
			}(i)
		}
		wg.Wait()

		// Merge parent and offspring
		newPop := objectives.MergeAgents(objectives.MergeAgents(ga.Population, popC), popM)
		newPop, paretoFront = objectives.FastNonDominatedSorting_Vectorized(newPop)
		newPop = ga.calculateCrowdingDistance(newPop, paretoFront)
		newPop, paretoFront = SortPopulation(newPop)

		// Truncate
		trunkedPop := make([]*objectives.Result, ga.PopulationSize)
		for i := 0; i < ga.PopulationSize; i++ {
			trunkedPop[i] = newPop[i].CopyAgent()
		}

		trunkedPop, paretoFront = objectives.FastNonDominatedSorting_Vectorized(trunkedPop)
		trunkedPop = ga.calculateCrowdingDistance(trunkedPop, paretoFront)
		trunkedPop, paretoFront = SortPopulation(trunkedPop)

		newResults := make([]*objectives.Result, len(paretoFront[0]))
		for i := 0; i < len(paretoFront[0]); i++ {
			newResults[i] = trunkedPop[i].CopyAgent()
		}

		ga.Population = trunkedPop
		ga.Archive = newResults

		// Send progress to channel
		channel <- struct {
			Progress                float64 `json:"progress"`
			NumberOfAgentsInArchive int     `json:"numberOfAgentsInArchive"`
			Type                    string  `json:"type"`
		}{
			Progress:                (float64(iter+1) / float64(ga.MaxIterations)) * 100,
			NumberOfAgentsInArchive: len(ga.Archive),
			Type:                    "multi",
		}
	}

	close(channel)

	return nil
}

func (ga *NSGAIIAlgorithm) initialization() {
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

}

// This function has been replaced by SortPopulation
// It is kept here for reference but is no longer used in the code

func (ga *NSGAIIAlgorithm) GetResults() algorithms.Result {
	// For multi-objective optimization, we return all non-dominated solutions in the archive
	results := make([]algorithms.AlgorithmResult, len(ga.Archive))

	for i, solution := range ga.Archive {
		mapLoc, sliceLoc, cranes, err := ga.ObjectiveFunction.GetLocationResult(solution.Position)
		if err != nil {
			continue
		}

		results[i] = algorithms.AlgorithmResult{
			MapLocations:   mapLoc,
			SliceLocations: sliceLoc,
			Value:          solution.Value,
			Penalty:        solution.Penalty,
			Key:            solution.Key,
			Cranes:         cranes,
			Phases:         ga.ObjectiveFunction.GetPhases(),
			ValuesWithKey:  solution.ValuesWithKey,
		}
	}

	minX, maxX, minY, maxY, _ := ga.ObjectiveFunction.GetLayoutSize()

	return algorithms.Result{
		Result: results,
		MinX:   minX,
		MaxX:   maxX,
		MinY:   minY,
		MaxY:   maxY,
	}
}

func (ga *NSGAIIAlgorithm) calculateCrowdingDistance(pop []*objectives.Result, paretoFront [][]int) []*objectives.Result {
	for k := 0; k < len(paretoFront); k++ {
		// Skip empty fronts
		if len(paretoFront[k]) == 0 {
			continue
		}

		// Extract costs (objective values) for this front
		costs := make([][]float64, ga.ObjectiveFunction.NumberOfObjectives())
		for j := range costs {
			costs[j] = make([]float64, len(paretoFront[k]))
		}

		for i, idx := range paretoFront[k] {
			for j := range pop[idx].Value {
				costs[j][i] = pop[idx].Value[j]
			}
		}

		nObj := len(costs)
		n := len(paretoFront[k])

		// Initialize distances
		d := make([][]float64, n)
		for i := range d {
			d[i] = make([]float64, nObj)
		}

		// Calculate crowding distance for each objective
		for j := 0; j < nObj; j++ {
			// Create pairs of (value, index) for sorting
			type ValueIndex struct {
				value float64
				index int
			}

			pairs := make([]ValueIndex, n)
			for i := 0; i < n; i++ {
				pairs[i] = ValueIndex{
					value: costs[j][i],
					index: i,
				}
			}

			// Sort by objective value
			sort.Slice(pairs, func(i, j int) bool {
				return pairs[i].value < pairs[j].value
			})

			// Set boundary points to infinity
			d[pairs[0].index][j] = math.Inf(1)
			d[pairs[n-1].index][j] = math.Inf(1)

			// Calculate distances for intermediate points
			fmin := pairs[0].value
			fmax := pairs[n-1].value

			// Skip if all values are the same (to avoid division by zero)
			diff := fmax - fmin
			if fmax == fmin {
				diff = 0.00000001
			}

			for i := 1; i < n-1; i++ {
				d[pairs[i].index][j] = math.Abs(pairs[i+1].value-pairs[i-1].value) / math.Abs(diff)
			}
		}

		// Sum up the distances for each solution
		for i := 0; i < n; i++ {
			sum := 0.0
			for j := 0; j < nObj; j++ {
				sum += d[i][j]
			}
			pop[paretoFront[k][i]].CrowdingDistance = sum
		}
	}

	return pop
}

// SortPopulation sorts the population by crowding distance (descending) and then by rank (ascending),
// and updates the fronts based on the ranks.
func SortPopulation(pop []*objectives.Result) ([]*objectives.Result, [][]int) {
	// Make a copy of the population to avoid modifying the original
	sorted := make([]*objectives.Result, len(pop))
	copy(sorted, pop)

	// Sort Based on Crowding Distance (descending)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CrowdingDistance > sorted[j].CrowdingDistance
	})

	// Sort Based on Rank (ascending)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Rank != sorted[j].Rank {
			return sorted[i].Rank < sorted[j].Rank
		}
		// If ranks are equal, maintain the crowding distance order (descending)
		return sorted[i].CrowdingDistance > sorted[j].CrowdingDistance
	})

	// Update Fronts
	// Find the maximum rank
	maxRank := 0
	for _, agent := range sorted {
		if agent.Rank > maxRank {
			maxRank = agent.Rank
		}
	}

	// Create fronts
	F := make([][]int, maxRank+1)
	for i := range F {
		F[i] = make([]int, 0)
	}

	// Populate fronts
	for i, agent := range sorted {
		rank := agent.Rank
		F[rank] = append(F[rank], i)
	}

	return sorted, F
}
