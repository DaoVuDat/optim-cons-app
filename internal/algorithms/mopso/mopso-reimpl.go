package mopso

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"golang-moaha-construction/internal/algorithms"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
	"golang-moaha-construction/internal/util"
	"math"
	"math/rand"
	"slices"
	"sync"
)

type HypercubeReimpl struct {
	Limits        [][]float64
	Quality       []float64
	NumberOfGrids int
}

func (h *HypercubeReimpl) updateGrid(archive []*objectives.Result, numberOfObjective int) []*objectives.Result {
	results := make([]*objectives.Result, len(archive))

	// find min and max
	minVals := make([]float64, numberOfObjective)
	maxVals := make([]float64, numberOfObjective)

	for i := 0; i < numberOfObjective; i++ {
		values := make([]float64, len(archive))
		for j, res := range archive {
			values[j] = res.Value[i]
		}

		minVals[i] = slices.Min(values)
		maxVals[i] = slices.Max(values)
	}

	limits := make([][]float64, numberOfObjective)
	for i := 0; i < numberOfObjective; i++ {
		limits[i] = util.LinSpace(minVals[i], maxVals[i], h.NumberOfGrids+1)
	}

	h.Limits = limits

	size := make([]int, numberOfObjective)
	for i := range size {
		size[i] = h.NumberOfGrids
	}

	// Calculate total number of grid cells (product of grid size in each dimension)
	totalGridCells := 1
	for i := 0; i < numberOfObjective; i++ {
		totalGridCells *= h.NumberOfGrids
	}
	h.Quality = make([]float64, totalGridCells)

	gridCount := make(map[int]int)

	for i := 0; i < len(archive); i++ {
		agent := archive[i].CopyAgent()
		agent.Idx = i
		subIndex := make([]int, numberOfObjective)
		for j := 0; j < numberOfObjective; j++ {
			idx := util.FindLessOrEqual(h.Limits[j], agent.Value[j])
			// Handle the case where FindLessOrEqual returns -1
			if idx == -1 {
				// Use the last valid grid cell index
				idx = h.NumberOfGrids - 1
			}
			subIndex[j] = idx
		}

		agent.GridSubIndex = subIndex
		agent.GridIndex = util.Sub2Index(size, subIndex...)

		if _, ok := gridCount[agent.GridIndex]; !ok {
			gridCount[agent.GridIndex] = 1
		} else {
			gridCount[agent.GridIndex]++
		}

		results[i] = agent
	}

	for idx, count := range gridCount {
		h.Quality[idx] = 10.0 / float64(count)
	}

	return results
}

type MOPSOAlgorithmReimpl struct {
	NumberOfAgents    int
	NumberOfIter      int
	Agents            []*ResultWithPersonalBest
	C1                float64
	C2                float64
	W                 float64
	ArchiveSize       int
	Archive           []*objectives.Result
	ObjectiveFunction objectives.Problem
	MaxVelocity       []float64
	NumberOfGrids     int
	MutationRate      float64
	hypercube         *HypercubeReimpl
}

func CreateReimpl(
	problem objectives.Problem,
	configs Config,
) (*MOPSOAlgorithmReimpl, error) {

	// calculate max velocity
	maxVelocity := make([]float64, problem.GetDimension())
	for i := 0; i < problem.GetDimension(); i++ {
		maxVelocity[i] = (problem.GetUpperBound()[i] - problem.GetLowerBound()[i]) * configs.MaxVelocity / 100
	}

	return &MOPSOAlgorithmReimpl{
		NumberOfAgents:    configs.NumberOfAgents,
		NumberOfIter:      configs.NumberOfIter,
		ArchiveSize:       configs.ArchiveSize,
		ObjectiveFunction: problem,
		NumberOfGrids:     configs.NumberOfGrids,
		MaxVelocity:       maxVelocity,
		MutationRate:      configs.MutationRate,
		C1:                configs.C1,
		C2:                configs.C2,
		hypercube: &HypercubeReimpl{
			NumberOfGrids: configs.NumberOfGrids,
			Limits:        make([][]float64, 0),
			Quality:       make([]float64, configs.NumberOfGrids),
		},
	}, nil
}

func (g *MOPSOAlgorithmReimpl) reset() {
	g.Agents = make([]*ResultWithPersonalBest, g.NumberOfAgents)
	g.Archive = make([]*objectives.Result, 0, g.ArchiveSize)
}

func (g *MOPSOAlgorithmReimpl) Type() data.TypeProblem {
	return data.Multi
}

func (g *MOPSOAlgorithmReimpl) Run() error {
	g.reset()

	bar := progressbar.Default(int64(g.NumberOfIter))
	g.initialization()

	// Initial archive: non-dominated solutions from initial population
	onlyAgents := getResultsFromResultWithPersonalBest(g.Agents)
	onlyAgents = objectives.DetermineDomination(onlyAgents)
	g.Archive = objectives.GetNonDominatedAgents(onlyAgents)
	g.Archive = g.hypercube.updateGrid(g.Archive, g.ObjectiveFunction.NumberOfObjectives())

	var wg sync.WaitGroup

	for iter := 0; iter < g.NumberOfIter; iter++ {

		leader := selectLeaderFromArchive(g.Archive)
		wg.Add(len(g.Agents))
		for i := range g.Agents {
			go func(agentIdx int) {
				defer wg.Done()
				agent := g.Agents[agentIdx]
				for d := 0; d < g.ObjectiveFunction.GetDimension(); d++ {
					r1 := rand.Float64()
					r2 := rand.Float64()
					v := g.W*agent.Velocity[d] +
						g.C1*r1*(agent.PersonalBest.Position[d]-agent.Result.Position[d]) +
						g.C2*r2*(leader.Position[d]-agent.Result.Position[d])
					agent.Velocity[d] = v
					agent.Result.Position[d] += v
				}
			}(i)
		}

		wg.Wait()
		g.Agents = g.applyMutation(iter)

		// Checking boundary
		g.checkingBoundaries()

		// Evaluate and update personal bests using goroutines
		wg.Add(len(g.Agents))
		for i := range g.Agents {
			go func(agentIdx int) {
				defer wg.Done()
				agent := g.Agents[agentIdx]
				value, valuesWithKey, keys, penalty := g.ObjectiveFunction.Eval(agent.Result.Position)
				agent.Result.Value = value
				agent.Result.ValuesWithKey = valuesWithKey
				agent.Result.Penalty = penalty
				agent.Result.Key = keys
			}(i)
		}
		wg.Wait()

		agentResults := getResultsFromResultWithPersonalBest(g.Agents)
		agentResults = objectives.DetermineDomination(agentResults)
		nonDominatedAgents := objectives.GetNonDominatedAgents(agentResults)

		g.Archive = objectives.MergeAgents(g.Archive, nonDominatedAgents)
		g.Archive = objectives.DetermineDomination(g.Archive)
		g.Archive = objectives.GetNonDominatedAgents(g.Archive)

		g.Archive = g.hypercube.updateGrid(g.Archive, g.ObjectiveFunction.NumberOfObjectives())

		// Truncate archive if needed (crowding distance)
		if len(g.Archive) > g.ArchiveSize {
			excess := len(g.Archive) - g.ArchiveSize
			g.Archive = g.removeExtraInArchive(excess)
		}

		// Update the best positions found so far for each particle
		// This is the MATLAB code implementation:
		// pos_best = dominates(POS_fit, PBEST_fit);
		// best_pos = ~dominates(PBEST_fit, POS_fit);
		// best_pos(rand(Np,1)>=0.5) = 0;

		// Create arrays to track which particles need updates
		posBest := make([]bool, len(g.Agents)) // Current position dominates personal best
		bestPos := make([]bool, len(g.Agents)) // Personal best doesn't dominate current position

		// Check domination relationships
		for i, agent := range g.Agents {
			// Check if current position dominates personal best
			if agent.Result.Dominates(agent.PersonalBest) {
				posBest[i] = true
			}

			// Check if personal best doesn't dominate current position
			if !agent.PersonalBest.Dominates(agent.Result) {
				// Apply random selection with 50% probability
				if rand.Float64() < 0.5 {
					bestPos[i] = true
				}
			}
		}

		// Update personal bests based on domination checks
		for i, agent := range g.Agents {
			if posBest[i] || bestPos[i] {
				agent.PersonalBest = agent.Result.CopyAgent()
			}
		}

		bar.Describe(fmt.Sprintf("Iteration %d: %d", iter+1, len(g.Archive)))
		bar.Add(1)
	}

	return nil
}

func (g *MOPSOAlgorithmReimpl) RunWithChannel(doneChan chan<- struct{}, channel chan<- any) error {
	g.reset()

	g.initialization()

	// Initial archive: non-dominated solutions from initial population
	onlyAgents := getResultsFromResultWithPersonalBest(g.Agents)
	onlyAgents = objectives.DetermineDomination(onlyAgents)
	g.Archive = objectives.GetNonDominatedAgents(onlyAgents)
	g.Archive = g.hypercube.updateGrid(g.Archive, g.ObjectiveFunction.NumberOfObjectives())

	var wg sync.WaitGroup

	for iter := 0; iter < g.NumberOfIter; iter++ {

		leader := selectLeaderFromArchive(g.Archive)
		for _, agent := range g.Agents {
			for d := 0; d < g.ObjectiveFunction.GetDimension(); d++ {
				r1 := rand.Float64()
				r2 := rand.Float64()
				v := g.W*agent.Velocity[d] +
					g.C1*r1*(agent.PersonalBest.Position[d]-agent.Result.Position[d]) +
					g.C2*r2*(leader.Position[d]-agent.Result.Position[d])
				agent.Velocity[d] = v
				agent.Result.Position[d] += v
			}
		}

		g.Agents = g.applyMutation(iter)

		// Checking boundary
		g.checkingBoundaries()

		// Evaluate and update personal bests using goroutines
		wg.Add(len(g.Agents))
		for i := range g.Agents {
			go func(agentIdx int) {
				defer wg.Done()
				agent := g.Agents[agentIdx]
				value, valuesWithKey, keys, penalty := g.ObjectiveFunction.Eval(agent.Result.Position)
				agent.Result.Value = value
				agent.Result.ValuesWithKey = valuesWithKey
				agent.Result.Penalty = penalty
				agent.Result.Key = keys
			}(i)
		}
		wg.Wait()

		agentResults := getResultsFromResultWithPersonalBest(g.Agents)
		agentResults = objectives.DetermineDomination(agentResults)
		nonDominatedAgents := objectives.GetNonDominatedAgents(agentResults)

		g.Archive = objectives.MergeAgents(g.Archive, nonDominatedAgents)
		g.Archive = objectives.DetermineDomination(g.Archive)
		g.Archive = objectives.GetNonDominatedAgents(g.Archive)

		g.Archive = g.hypercube.updateGrid(g.Archive, g.ObjectiveFunction.NumberOfObjectives())

		// Truncate archive if needed (crowding distance)
		if len(g.Archive) > g.ArchiveSize {
			excess := len(g.Archive) - g.ArchiveSize
			g.Archive = g.removeExtraInArchive(excess)
		}

		// Update the best positions found so far for each particle
		// Create arrays to track which particles need updates
		posBest := make([]bool, len(g.Agents)) // Current position dominates personal best
		bestPos := make([]bool, len(g.Agents)) // Personal best doesn't dominate current position

		// Check domination relationships
		for i, agent := range g.Agents {
			// Check if current position dominates personal best
			if agent.Result.Dominates(agent.PersonalBest) {
				posBest[i] = true
			}

			// Check if personal best doesn't dominate current position
			if !agent.PersonalBest.Dominates(agent.Result) {
				// Apply random selection with 50% probability
				if rand.Float64() < 0.5 {
					bestPos[i] = true
				}
			}
		}

		// Update personal bests based on domination checks
		for i, agent := range g.Agents {
			if posBest[i] || bestPos[i] {
				agent.PersonalBest = agent.Result.CopyAgent()
			}
		}

		// Send progress update through channel
		channel <- struct {
			Progress                float64 `json:"progress"`
			NumberOfAgentsInArchive int     `json:"numberOfAgentsInArchive"`
			Type                    string  `json:"type"`
		}{
			Progress:                (float64(iter+1) / float64(g.NumberOfIter)) * 100,
			NumberOfAgentsInArchive: len(g.Archive),
			Type:                    "multi",
		}
	}

	// Signal completion
	close(channel)
	doneChan <- struct{}{}

	return nil
}

func (g *MOPSOAlgorithmReimpl) initialization() {
	var wg sync.WaitGroup
	wg.Add(g.NumberOfAgents)
	for agentIdx := 0; agentIdx < g.NumberOfAgents; agentIdx++ {
		go func(agentIdx int) {
			defer wg.Done()
			positions := make([]float64, g.ObjectiveFunction.GetDimension())

			for i := 0; i < g.ObjectiveFunction.GetDimension(); i++ {
				positions[i] = g.ObjectiveFunction.GetLowerBound()[i] + rand.Float64()*
					(g.ObjectiveFunction.GetUpperBound()[i]-g.ObjectiveFunction.GetLowerBound()[i])
			}

			// evaluate
			result := &objectives.Result{
				Idx:      agentIdx,
				Position: positions,
			}

			value, valuesWithKey, keys, penalty := g.ObjectiveFunction.Eval(positions)
			result.Value = value
			result.ValuesWithKey = valuesWithKey
			result.Penalty = penalty
			result.Key = keys

			newAgent := &ResultWithPersonalBest{
				Result:       result,
				PersonalBest: result.CopyAgent(),
				Velocity:     make([]float64, g.ObjectiveFunction.GetDimension()),
			}

			g.Agents[agentIdx] = newAgent
		}(agentIdx)
	}
	wg.Wait()
}

func (g *MOPSOAlgorithmReimpl) GetResults() algorithms.Result {
	results := make([]algorithms.AlgorithmResult, len(g.Archive))

	for i := range g.Archive {
		res := g.Archive[i]
		mapLoc, sliceLoc, cranes, err := g.ObjectiveFunction.GetLocationResult(res.Position)
		if err != nil {
			return algorithms.Result{}
		}

		results[i] = algorithms.AlgorithmResult{
			MapLocations:   mapLoc,
			SliceLocations: sliceLoc,
			Value:          res.Value,
			Key:            res.Key,
			Penalty:        res.Penalty,
			ValuesWithKey:  res.ValuesWithKey,
			Cranes:         cranes,
			Phases:         g.ObjectiveFunction.GetPhases(),
		}
	}

	minX, maxX, minY, maxY, _ := g.ObjectiveFunction.GetLayoutSize()

	return algorithms.Result{
		Result: results,
		MinX:   minX,
		MinY:   minY,
		MaxX:   maxX,
		MaxY:   maxY,
	}
}

func (g *MOPSOAlgorithmReimpl) selectLeader(archive []*objectives.Result) *objectives.Result {
	// roulette wheel
	prob := make([]float64, len(g.hypercube.Quality))
	prob[0] = g.hypercube.Quality[0]
	for i := 1; i < len(g.hypercube.Quality); i++ {
		prob[i] = prob[i-1] + g.hypercube.Quality[i]
	}

	randVal := rand.Float64() * slices.Max(prob)
	selectedIndex := util.FindLessOrEqual(prob, randVal)

	agentsInGrid := make([]*objectives.Result, 0)
	for _, agent := range archive {
		if agent.GridIndex == selectedIndex {
			agentsInGrid = append(agentsInGrid, agent)
		}
	}

	if len(agentsInGrid) > 0 {
		return agentsInGrid[rand.Intn(len(agentsInGrid))].CopyAgent()
	}

	fmt.Println("No agents in grid")
	return archive[rand.Intn(len(archive))].CopyAgent()
}

func (g *MOPSOAlgorithmReimpl) applyMutation(curIter int) []*ResultWithPersonalBest {
	third := float64(g.NumberOfAgents) / 3
	frac := third - math.Floor(third)

	subSizes := make([]int, 3)
	if frac < 0.5 {
		subSizes[0] = int(math.Ceil(third))
		subSizes[1] = int(math.Round(third))
		subSizes[2] = int(math.Round(third))
	} else {
		subSizes[0] = int(math.Round(third))
		subSizes[1] = int(math.Round(third))
		subSizes[2] = int(math.Floor(third))
	}

	cumSum := make([]int, 3)
	cumSum[0] = subSizes[0]
	for i := 1; i < len(subSizes); i++ {
		cumSum[i] = cumSum[i-1] + subSizes[i]
	}

	// 1st part: no mutation

	// 2nd part: uniform mutation
	nMut := int(math.Round(g.MutationRate * float64(subSizes[1])))
	if nMut > 0 {
		tempIndices := rand.Perm(subSizes[1])
		for i := 0; i < nMut; i++ {
			idx := cumSum[0] + tempIndices[i]
			for d := 0; d < g.ObjectiveFunction.GetDimension(); d++ {
				g.Agents[idx].Result.Position[d] = g.ObjectiveFunction.GetLowerBound()[d] + rand.Float64()*(g.ObjectiveFunction.GetUpperBound()[d]-g.ObjectiveFunction.GetLowerBound()[d])
			}
		}
	}

	// 3rd part: non-uniform mutation
	perMut := math.Pow(1-float64(curIter)/float64(g.NumberOfIter), 5*float64(g.ObjectiveFunction.GetDimension())) // percentage mutation

	nMut = int(math.Round(perMut * float64(subSizes[2])))
	if nMut > 0 {
		tempIndices := rand.Perm(subSizes[2])
		for i := 0; i < nMut; i++ {
			idx := cumSum[1] + tempIndices[i]
			for d := 0; d < g.ObjectiveFunction.GetDimension(); d++ {
				g.Agents[idx].Result.Position[d] = g.ObjectiveFunction.GetLowerBound()[d] + rand.Float64()*(g.ObjectiveFunction.GetUpperBound()[d]-g.ObjectiveFunction.GetLowerBound()[d])
			}
		}
	}

	return g.Agents
}

func (g *MOPSOAlgorithmReimpl) checkingBoundaries() {
	for _, agent := range g.Agents {
		for i := range agent.Result.Position {
			if agent.Velocity[i] > g.MaxVelocity[i] {
				agent.Velocity[i] = g.MaxVelocity[i]
			} else if agent.Velocity[i] < -g.MaxVelocity[i] {
				agent.Velocity[i] = -g.MaxVelocity[i]
			}

			if agent.Result.Position[i] < g.ObjectiveFunction.GetLowerBound()[i] {
				agent.Result.Position[i] = g.ObjectiveFunction.GetLowerBound()[i]
				agent.Velocity[i] = -agent.Velocity[i]
			} else if agent.Result.Position[i] > g.ObjectiveFunction.GetUpperBound()[i] {
				agent.Result.Position[i] = g.ObjectiveFunction.GetUpperBound()[i]
				agent.Velocity[i] = -agent.Velocity[i]
			}
		}
	}
}

func (g *MOPSOAlgorithmReimpl) removeExtraInArchive(exceeded int) []*objectives.Result {
	// Compute the crowding distances
	crowding := make([]float64, len(g.Archive))

	// For each objective
	for m := 0; m < g.ObjectiveFunction.NumberOfObjectives(); m++ {
		// Extract values for this objective
		mFit := make([]float64, len(g.Archive))
		for i, res := range g.Archive {
			mFit[i] = res.Value[m]
		}

		// Create indices for sorting
		idx := make([]int, len(mFit))
		for i := range idx {
			idx[i] = i
		}

		// Sort indices by objective values (ascending)
		slices.SortFunc(idx, func(i, j int) int {
			diff := mFit[i] - mFit[j]
			switch {
			case diff < 0:
				return -1
			case diff > 0:
				return 1
			default:
				return 0
			}
		})

		// Get min and max values for normalization
		minVal := mFit[idx[0]]
		maxVal := mFit[idx[len(idx)-1]]

		// Calculate distances
		for i := 0; i < len(idx); i++ {
			var distance float64

			if i == 0 || i == len(idx)-1 {
				// Boundary points get infinite distance
				distance = math.Inf(1)
			} else {
				// Calculate normalized distance between neighbors
				mUp := mFit[idx[i+1]]
				mDown := mFit[idx[i-1]]

				// Normalize by range (max-min)
				if maxVal != minVal {
					distance = (mUp - mDown) / (maxVal - minVal)
				} else {
					distance = math.Inf(1) // Handle division by zero
				}
			}

			// Add distance to crowding for this index
			crowding[idx[i]] += distance
		}
	}

	// Set NaN values to Infinity
	for i := range crowding {
		if math.IsNaN(crowding[i]) {
			crowding[i] = math.Inf(1)
		}
	}

	// Sort archive by crowding distance (ascending)
	indices := make([]int, len(crowding))
	for i := range indices {
		indices[i] = i
	}

	slices.SortFunc(indices, func(i, j int) int {
		diff := crowding[i] - crowding[j]
		switch {
		case diff < 0:
			return -1
		case diff > 0:
			return 1
		default:
			return 0
		}
	})

	// Remove agents with smallest crowding distances
	remainingIndices := indices[exceeded:]

	// Create new archive with remaining agents
	newArchive := make([]*objectives.Result, len(remainingIndices))
	for i, idx := range remainingIndices {
		newArchive[i] = g.Archive[idx].CopyAgent()
	}

	newArchive = g.hypercube.updateGrid(newArchive, g.ObjectiveFunction.NumberOfObjectives())

	return newArchive
}
