package nsgaii

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

const NameType algorithms.AlgorithmType = "NSGA-II"

type NSGAIIAlgorithm struct {
	PopulationSize    int
	MaxIterations     int
	CrossoverRate     float64
	MutationRate      float64
	MutationStrength  float64
	TournamentSize    int
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
	TournamentSize   int
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
		TournamentSize:    configs.TournamentSize,
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

	lowerBound := ga.ObjectiveFunction.GetLowerBound()
	upperBound := ga.ObjectiveFunction.GetUpperBound()
	numObjectives := ga.ObjectiveFunction.NumberOfObjectives()

	// Parameters for NSGA-II
	pc := ga.CrossoverRate    // Probability of crossover
	pm := ga.MutationRate     // Probability of mutation
	ms := ga.MutationStrength // Mutation strength

	// Initialize Q (offspring) as empty for first iteration
	var Q []*objectives.Result

	// Main NSGA-II loop
	for iter := 0; iter < ga.MaxIterations; iter++ {
		// Merge the parent and the children (for the first iteration, Q will be empty)
		var R []*objectives.Result
		if iter == 0 {
			R = ga.Population
		} else {
			// Merge parent and offspring
			R = objectives.MergeAgents(ga.Population, Q)
		}

		// Convert to the format needed for our functions
		Rfit := getObjectiveValues(R)

		// Compute the new Pareto Fronts using Fast Non-Dominated Sorting
		Rrank := FastNonDominatedSorting_Vectorized(Rfit)

		// Compute crowding distances
		Rcrowd, Rrank, _, R := crowdingDistances(Rrank, Rfit, R)

		// Select parent population for next generation
		ga.Population = selectParentByRankAndDistance(Rcrowd, Rrank, R)

		// Apply crossover and mutation to create offspring for next iteration
		Q = applyCrossoverAndMutation(ga.Population, pc, pm, ms, upperBound, lowerBound)

		// Evaluate offspring
		wg.Add(len(Q))
		for i := range Q {
			go func(idx int) {
				defer wg.Done()
				value, valuesWithKey, keys, penalty := ga.ObjectiveFunction.Eval(Q[idx].Position)
				Q[idx].Value = value
				Q[idx].ValuesWithKey = valuesWithKey
				Q[idx].Penalty = penalty
				Q[idx].Key = keys
			}(i)
		}
		wg.Wait()

		// Update archive with non-dominated solutions
		// First, identify rank 1 solutions (non-dominated)
		var nonDominated []*objectives.Result
		for i, rank := range Rrank {
			if rank == 1 && i < len(R) {
				nonDominated = append(nonDominated, R[i].CopyAgent())
			}
		}

		// Combine with existing archive
		if len(ga.Archive) > 0 {
			combined := objectives.MergeAgents(ga.Archive, nonDominated)

			// Convert to the format needed for our functions
			combinedFit := getObjectiveValues(combined)

			// Compute ranks
			combinedRank := FastNonDominatedSorting_Vectorized(combinedFit)

			// Keep only rank 1 solutions
			var newArchive []*objectives.Result
			for i, rank := range combinedRank {
				if rank == 1 && i < len(combined) {
					newArchive = append(newArchive, combined[i].CopyAgent())
				}
			}

			ga.Archive = newArchive
		} else {
			ga.Archive = nonDominated
		}

		// If archive exceeds size limit, use crowding distance to reduce it
		if len(ga.Archive) > ga.PopulationSize {
			ga.Archive = reduceArchiveByDistance(ga.Archive, ga.PopulationSize)
		}

		// Store convergence metrics
		for obj := 0; obj < numObjectives; obj++ {
			if len(ga.Archive) > 0 {
				// Find the best value for this objective in the archive
				bestVal := ga.Archive[0].Value[obj]
				for _, sol := range ga.Archive {
					if sol.Value[obj] < bestVal {
						bestVal = sol.Value[obj]
					}
				}
			}
		}

		// Update progress bar
		bar.Describe(fmt.Sprintf("Iter %d: Archive size %d", iter+1, len(ga.Archive)))
		bar.Add(1)
	}

	return nil
}

func (ga *NSGAIIAlgorithm) RunWithChannel(doneChan chan<- struct{}, channel chan<- any) error {
	ga.reset()

	ga.initialization()

	var wg sync.WaitGroup

	lowerBound := ga.ObjectiveFunction.GetLowerBound()
	upperBound := ga.ObjectiveFunction.GetUpperBound()
	numObjectives := ga.ObjectiveFunction.NumberOfObjectives()

	// Parameters for NSGA-II
	pc := ga.CrossoverRate    // Probability of crossover
	pm := ga.MutationRate     // Probability of mutation
	ms := ga.MutationStrength // Mutation strength

	// Initialize Q (offspring) as empty for first iteration
	var Q []*objectives.Result

	// Main NSGA-II loop
	for iter := 0; iter < ga.MaxIterations; iter++ {
		// Merge the parent and the children (for the first iteration, Q will be empty)
		var R []*objectives.Result
		if iter == 0 {
			R = ga.Population
		} else {
			// Merge parent and offspring
			R = objectives.MergeAgents(ga.Population, Q)
		}

		// Convert to the format needed for our functions
		Rfit := getObjectiveValues(R)

		// Compute the new Pareto Fronts using Fast Non-Dominated Sorting
		Rrank := FastNonDominatedSorting_Vectorized(Rfit)

		// Compute crowding distances
		Rcrowd, Rrank, _, R := crowdingDistances(Rrank, Rfit, R)

		// Select parent population for next generation
		ga.Population = selectParentByRankAndDistance(Rcrowd, Rrank, R)

		// Apply crossover and mutation to create offspring for next iteration
		Q = applyCrossoverAndMutation(ga.Population, pc, pm, ms, upperBound, lowerBound)

		// Evaluate offspring
		wg.Add(len(Q))
		for i := range Q {
			go func(idx int) {
				defer wg.Done()
				value, valuesWithKey, keys, penalty := ga.ObjectiveFunction.Eval(Q[idx].Position)
				Q[idx].Value = value
				Q[idx].ValuesWithKey = valuesWithKey
				Q[idx].Penalty = penalty
				Q[idx].Key = keys
			}(i)
		}
		wg.Wait()

		// Update archive with non-dominated solutions
		// First, identify rank 1 solutions (non-dominated)
		var nonDominated []*objectives.Result
		for i, rank := range Rrank {
			if rank == 1 && i < len(R) {
				nonDominated = append(nonDominated, R[i].CopyAgent())
			}
		}

		// Combine with existing archive
		if len(ga.Archive) > 0 {
			combined := objectives.MergeAgents(ga.Archive, nonDominated)

			// Convert to the format needed for our functions
			combinedFit := getObjectiveValues(combined)

			// Compute ranks
			combinedRank := FastNonDominatedSorting_Vectorized(combinedFit)

			// Keep only rank 1 solutions
			var newArchive []*objectives.Result
			for i, rank := range combinedRank {
				if rank == 1 && i < len(combined) {
					newArchive = append(newArchive, combined[i].CopyAgent())
				}
			}

			ga.Archive = newArchive
		} else {
			ga.Archive = nonDominated
		}

		// If archive exceeds size limit, use crowding distance to reduce it
		if len(ga.Archive) > ga.PopulationSize {
			ga.Archive = reduceArchiveByDistance(ga.Archive, ga.PopulationSize)
		}

		// Store convergence metrics
		for obj := 0; obj < numObjectives; obj++ {
			if len(ga.Archive) > 0 {
				// Find best value for this objective in the archive
				bestVal := ga.Archive[0].Value[obj]
				for _, sol := range ga.Archive {
					if sol.Value[obj] < bestVal {
						bestVal = sol.Value[obj]
					}
				}
			}
		}

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

	// Initialize archive with non-dominated solutions from initial population
	ga.Population = objectives.DetermineDomination(ga.Population)
	ga.Archive = objectives.GetNonDominatedAgents(ga.Population)
}

// sortPopulationByRankAndCrowding sorts the population by rank and crowding distance
func (ga *NSGAIIAlgorithm) sortPopulationByRankAndCrowding() []*objectives.Result {
	sorted := make([]*objectives.Result, len(ga.Population))
	copy(sorted, ga.Population)

	// Sort by rank first, then by crowding distance (higher is better)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Rank != sorted[j].Rank {
			return sorted[i].Rank < sorted[j].Rank
		}
		return sorted[i].CrowdingDistance > sorted[j].CrowdingDistance
	})

	return sorted
}

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

// Function that computes the crowding distances of every single ParetoFront
func crowdingDistances(rank []int, fitness [][]*float64, pop []*objectives.Result) ([]float64, []int, [][]*float64, []*objectives.Result) {
	// Initialize
	sortPop := make([]*objectives.Result, 0)
	sortFit := make([][]*float64, 0)
	sortRank := make([]int, 0)
	sortCrowd := make([]float64, 0)

	Npf := 0
	for _, r := range rank {
		if r > Npf {
			Npf = r
		}
	}

	for pf := 1; pf <= Npf; pf++ {
		// Find indices of solutions in this front
		var indices []int
		for i, r := range rank {
			if r == pf {
				indices = append(indices, i)
			}
		}

		// Extract fitness and population for this front
		tempFit := make([][]*float64, len(indices))
		tempRank := make([]int, len(indices))
		tempPop := make([]*objectives.Result, len(indices))

		for i, idx := range indices {
			tempFit[i] = fitness[idx]
			tempRank[i] = rank[idx]
			tempPop[i] = pop[idx]
		}

		// Sort by first dimension
		sort.Slice(tempFit, func(i, j int) bool {
			return *tempFit[i][0] < *tempFit[j][0]
		})

		// Update sorted arrays
		sortFit = append(sortFit, tempFit...)
		sortRank = append(sortRank, tempRank...)
		sortPop = append(sortPop, tempPop...)

		// Crowded distances
		tempCrowd := make([]float64, len(tempRank))

		// For each objective
		for m := 0; m < len(fitness[0]); m++ {
			// Skip if tempFit is empty
			if len(tempFit) == 0 {
				continue
			}

			// Find min and max for normalization
			tempMin := *tempFit[0][m]
			tempMax := *tempFit[0][m]

			for i := 1; i < len(tempFit); i++ {
				if *tempFit[i][m] < tempMin {
					tempMin = *tempFit[i][m]
				}
				if *tempFit[i][m] > tempMax {
					tempMax = *tempFit[i][m]
				}
			}

			// Calculate crowding distance
			// Skip if there are not enough elements for calculation (need at least 3)
			if len(tempFit) >= 3 {
				for l := 1; l < len(tempCrowd)-1; l++ {
					if tempMax-tempMin > 0 {
						tempCrowd[l] += (math.Abs(*tempFit[l-1][m] - *tempFit[l+1][m])) / (tempMax - tempMin)
					}
				}
			}
		}

		// Set boundary points to infinity
		if len(tempCrowd) > 0 {
			tempCrowd[0] = math.Inf(1)
			if len(tempCrowd) > 1 {
				tempCrowd[len(tempCrowd)-1] = math.Inf(1)
			}
		}

		sortCrowd = append(sortCrowd, tempCrowd...)
	}

	return sortCrowd, sortRank, sortFit, sortPop
}

// Function that selects a new parent based on the crowding distance operator
func selectParentByRankAndDistance(Rcrowd []float64, Rrank []int, R []*objectives.Result) []*objectives.Result {
	// Initialization
	N := len(Rcrowd) / 2
	if N > len(R) {
		N = len(R)
	}

	Npf := 0
	for _, r := range Rrank {
		if r > Npf {
			Npf = r
		}
	}

	newParent := make([]*objectives.Result, 0, N)

	// Selecting the chromosomes
	pf := 1
	numberOfSolutions := 0

	for pf <= Npf {
		// Count solutions in this front
		countInFront := 0
		for _, r := range Rrank {
			if r == pf {
				countInFront++
			}
		}

		// If there is enough space, select solutions based on rank
		if numberOfSolutions+countInFront <= N {
			for i, r := range Rrank {
				if r == pf {
					newParent = append(newParent, R[i])
					numberOfSolutions++
				}
			}
		} else {
			// If there isn't enough space, sort by crowding distances
			rest := N - numberOfSolutions

			// Get solutions in this front
			lastPF := make([]*objectives.Result, 0)
			lastPFdist := make([]float64, 0)

			for i, r := range Rrank {
				if r == pf {
					lastPF = append(lastPF, R[i])
					lastPFdist = append(lastPFdist, Rcrowd[i])
				}
			}

			// Sort by crowding distance (descending)
			type SolutionWithDist struct {
				solution *objectives.Result
				distance float64
			}

			solutionsWithDist := make([]SolutionWithDist, len(lastPF))
			for i := range lastPF {
				solutionsWithDist[i] = SolutionWithDist{
					solution: lastPF[i],
					distance: lastPFdist[i],
				}
			}

			sort.Slice(solutionsWithDist, func(i, j int) bool {
				return solutionsWithDist[i].distance > solutionsWithDist[j].distance
			})

			// Add solutions with highest crowding distance
			for i := 0; i < rest && i < len(solutionsWithDist); i++ {
				newParent = append(newParent, solutionsWithDist[i].solution)
				numberOfSolutions++
			}
		}

		pf++
	}

	return newParent
}

// Function that performs a Fast Non Dominated Sorting algorithm
func FastNonDominatedSorting_Vectorized(fitness [][]*float64) []int {
	// Initialization
	Np := len(fitness)
	RANK := make([]int, Np)

	// Initialize all ranks to 0
	for i := range RANK {
		RANK[i] = 0
	}

	// Check domination for all pairs
	for i := 0; i < Np-1; i++ {
		for j := i + 1; j < Np; j++ {
			if dominates(fitness[i], fitness[j]) {
				RANK[j]++
			} else if dominates(fitness[j], fitness[i]) {
				RANK[i]++
			}
		}
	}

	// Assign ranks (add 1 to make ranks start from 1)
	for i := range RANK {
		RANK[i]++
	}

	return RANK
}

// Function that returns true if x dominates y and false otherwise
func dominates(x, y []*float64) bool {
	anyBetter := false
	for i := range x {
		if *x[i] > *y[i] {
			return false
		}
		if *x[i] < *y[i] {
			anyBetter = true
		}
	}
	return anyBetter
}

// Function that calculates a child population by applying crossover and mutation
func applyCrossoverAndMutation(parent []*objectives.Result, pc, pm float64, ms float64, var_max, var_min []float64) []*objectives.Result {
	// Params
	N := len(parent)

	// Return empty slice if parent is empty
	if N == 0 {
		return []*objectives.Result{}
	}

	nVar := len(parent[0].Position)

	// Child initialization
	Q := make([]*objectives.Result, N)
	for i := range Q {
		Q[i] = parent[i].CopyAgent()
	}

	// Crossover
	for c := 0; c < N; c++ {
		if rand.Float64() < pc {
			selected := rand.Intn(N)
			for selected == c {
				selected = rand.Intn(N)
			}

			cut := rand.Intn(nVar)

			// Create new position array
			newPos := make([]float64, nVar)
			copy(newPos[:cut], parent[c].Position[:cut])
			copy(newPos[cut:], parent[selected].Position[cut:])

			Q[c].Position = newPos
		}
	}

	// Mutation
	for i := 0; i < N; i++ {
		for j := 0; j < nVar; j++ {
			if rand.Float64() < pm {
				// Apply Gaussian mutation
				range_val := var_max[j] - var_min[j]
				Q[i].Position[j] += ms * range_val * rand.NormFloat64()

				// Ensure bounds
				if Q[i].Position[j] < var_min[j] {
					Q[i].Position[j] = var_min[j]
				}
				if Q[i].Position[j] > var_max[j] {
					Q[i].Position[j] = var_max[j]
				}
			}
		}
	}

	return Q
}

// Helper function to convert objectives.Result values to the format needed for our functions
func getObjectiveValues(agents []*objectives.Result) [][]*float64 {
	result := make([][]*float64, len(agents))
	for i, agent := range agents {
		result[i] = make([]*float64, len(agent.Value))
		for j, val := range agent.Value {
			value := val // Create a new variable to avoid capturing the loop variable
			result[i][j] = &value
		}
	}
	return result
}

// reduceArchiveByDistance reduces the archive size using crowding distance
func reduceArchiveByDistance(archive []*objectives.Result, targetSize int) []*objectives.Result {

	if len(archive) <= targetSize {
		return archive
	}

	// Convert to the format needed for our functions
	archiveFit := getObjectiveValues(archive)

	// All solutions in archive are rank 1 (non-dominated)
	archiveRank := make([]int, len(archive))
	for i := range archiveRank {
		archiveRank[i] = 1
	}

	// Compute crowding distances
	archiveCrowd, _, _, archiveWithCrowd := crowdingDistances(archiveRank, archiveFit, archive)

	// Sort by crowding distance (descending)
	type SolutionWithDist struct {
		solution *objectives.Result
		distance float64
	}

	solutionsWithDist := make([]SolutionWithDist, len(archiveWithCrowd))
	for i := range archiveWithCrowd {
		solutionsWithDist[i] = SolutionWithDist{
			solution: archiveWithCrowd[i],
			distance: archiveCrowd[i],
		}
	}

	sort.Slice(solutionsWithDist, func(i, j int) bool {
		return solutionsWithDist[i].distance > solutionsWithDist[j].distance
	})

	// Keep only the solutions with highest crowding distance
	actualSize := targetSize
	if len(solutionsWithDist) < targetSize {
		actualSize = len(solutionsWithDist)
	}
	newArchive := make([]*objectives.Result, actualSize)
	for i := 0; i < actualSize; i++ {
		newArchive[i] = solutionsWithDist[i].solution
	}

	return newArchive
}
