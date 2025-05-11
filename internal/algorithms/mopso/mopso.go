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

const NameType algorithms.AlgorithmType = "MOPSO"

type resultsWithGridIndex struct {
	GridIndex    int
	GridSubIndex []int
	Result       *objectives.Result
}

func convertArchiveIntoResultArchiveWithGridIndex(archive []*objectives.Result) []*resultsWithGridIndex {
	results := make([]*resultsWithGridIndex, len(archive))

	for i, res := range archive {
		results[i] = &resultsWithGridIndex{
			GridIndex:    0,
			GridSubIndex: []int{},
			Result:       res,
		}
	}

	return results
}

func convertResultArchiveIntoArchive(results []*resultsWithGridIndex) []*objectives.Result {
	archive := make([]*objectives.Result, len(results))

	for i, res := range results {
		archive[i] = res.Result
	}
	return archive
}

func getResultsFromArchive(archive []*objectives.Result) [][]float64 {
	results := make([][]float64, len(archive[0].Value))

	for _, res := range archive {
		for j, val := range res.Value {
			results[j] = append(results[j], val)
		}

	}

	return results
}

type ResultWithPersonalBest struct {
	Result       *objectives.Result
	Velocity     []float64
	PersonalBest *objectives.Result
}

func (r *ResultWithPersonalBest) Copy() *ResultWithPersonalBest {
	return &ResultWithPersonalBest{
		Result:       r.Result.CopyAgent(),
		PersonalBest: r.PersonalBest.CopyAgent(),
		Velocity:     util.CopyArray(r.Velocity),
	}
}

func (r *ResultWithPersonalBest) GetPersonalBest() *objectives.Result {
	return r.PersonalBest
}

func getResultsFromResultWithPersonalBest(archive []*ResultWithPersonalBest) []*objectives.Result {
	results := make([]*objectives.Result, len(archive))

	for i, res := range archive {
		results[i] = res.Result
	}

	return results
}

type MOPSOAlgorithm struct {
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
	// for multi-objective hypercubes
	NumberOfGrids int
	MutationRate  float64

	hypercube Hypercube
}

type Config struct {
	NumberOfAgents int
	NumberOfIter   int
	ArchiveSize    int
	NumberOfGrids  int
	MutationRate   float64
	MaxVelocity    float64
	C1             float64
	C2             float64
	W              float64
}

func Create(
	problem objectives.Problem,
	configs Config,
) (*MOPSOAlgorithm, error) {

	maxVelocity := make([]float64, problem.GetDimension())
	for i := 0; i < problem.GetDimension(); i++ {
		maxVelocity[i] = (problem.GetUpperBound()[i] - problem.GetLowerBound()[i]) * configs.MaxVelocity / 100
	}

	return &MOPSOAlgorithm{
		NumberOfAgents:    configs.NumberOfAgents,
		NumberOfIter:      configs.NumberOfIter,
		ArchiveSize:       configs.ArchiveSize,
		ObjectiveFunction: problem,
		NumberOfGrids:     configs.NumberOfGrids,
		MaxVelocity:       maxVelocity,
		MutationRate:      configs.MutationRate,
		C1:                configs.C1,
		C2:                configs.C2,
		W:                 configs.W,
		hypercube: Hypercube{
			NumberOfGrids: configs.NumberOfGrids,
		},
	}, nil
}

func (g *MOPSOAlgorithm) reset() {
	g.Agents = make([]*ResultWithPersonalBest, g.NumberOfAgents)
	g.Archive = make([]*objectives.Result, 0, g.ArchiveSize)
}

func (g *MOPSOAlgorithm) Type() data.TypeProblem {
	return data.Multi
}

func (g *MOPSOAlgorithm) Run() error {
	g.reset()

	bar := progressbar.Default(int64(g.NumberOfIter))
	// initialization
	g.initialization()

	onlyAgents := getResultsFromResultWithPersonalBest(g.Agents)

	onlyAgents = objectives.DetermineDomination(onlyAgents)
	g.Archive = objectives.GetNonDominatedAgents(onlyAgents)
	g.hypercube.UpdateHyperCube(getResultsFromArchive(g.Archive))

	l := 0

	for l < g.NumberOfIter {
		// selectLeader
		leaderAgent := g.hypercube.SelectLeader(g.Archive)

		for agentIdx := range g.Agents {
			for posIdx := range g.ObjectiveFunction.GetDimension() {
				velocity := g.W*g.Agents[agentIdx].Velocity[posIdx] +
					g.C1*rand.Float64()*(g.Agents[agentIdx].GetPersonalBest().Position[posIdx]-g.Agents[agentIdx].Result.Position[posIdx]) +
					g.C2*rand.Float64()*(leaderAgent.Position[posIdx]-g.Agents[agentIdx].Result.Position[posIdx])

				g.Agents[agentIdx].Velocity[posIdx] = velocity

				g.Agents[agentIdx].Result.Position[posIdx] = g.Agents[agentIdx].Result.Position[posIdx] + velocity
			}
		}

		// apply mutation
		g.Agents = mutation(g.Agents, l,
			g.NumberOfIter,
			g.NumberOfAgents,
			g.ObjectiveFunction.GetLowerBound(),
			g.ObjectiveFunction.GetUpperBound(),
			g.MutationRate)

		for agentIdx := range g.Agents {
			// check out of boundaries
			g.outOfBoundaries(g.Agents[agentIdx].Result.Position)
			// evaluate
			value, valuesWithKey, keys, penalty := g.ObjectiveFunction.Eval(g.Agents[agentIdx].Result.Position)
			g.Agents[agentIdx].Result.Value = value
			g.Agents[agentIdx].Result.Penalty = penalty
			g.Agents[agentIdx].Result.Key = keys
			g.Agents[agentIdx].Result.ValuesWithKey = valuesWithKey
		}

		// update repository
		g.Archive = g.updateRepository(g.Archive, g.Agents)
		if len(g.Archive) > g.ArchiveSize {
			exceeded := len(g.Archive) - g.ArchiveSize
			g.Archive = g.hypercube.deleteFromRepository(g.Archive, exceeded)
		}

		// update PersonalBest
		curDominatesPBest := make([]bool, len(g.Agents))
		pBestDominatesCur := make([]bool, len(g.Agents))

		for i := 0; i < len(g.Agents); i++ {
			if g.Agents[i].Result.Dominates(g.Agents[i].PersonalBest) {
				curDominatesPBest[i] = true
			}

			if !g.Agents[i].PersonalBest.Dominates(g.Agents[i].Result) {
				if rand.Float64() > 0.5 {
					pBestDominatesCur[i] = true
				}
			}
		}

		for i := 0; i < len(g.Agents); i++ {
			if curDominatesPBest[i] || pBestDominatesCur[i] {
				g.Agents[i].PersonalBest = g.Agents[i].Result.CopyAgent()
			}

		}

		bar.Describe(fmt.Sprintf("Iteration %d: %d", l+1, len(g.Archive)))
		bar.Add(1)
		l++
	}

	return nil
}

func (g *MOPSOAlgorithm) RunWithChannel(doneChan chan<- struct{}, channel chan<- any) error {
	g.reset()

	// initialization
	g.initialization()

	onlyAgents := getResultsFromResultWithPersonalBest(g.Agents)

	onlyAgents = objectives.DetermineDomination(onlyAgents)
	g.Archive = objectives.GetNonDominatedAgents(onlyAgents)
	g.hypercube.UpdateHyperCube(getResultsFromArchive(g.Archive))

	l := 0

	for l < g.NumberOfIter {
		// selectLeader
		leaderAgent := g.hypercube.SelectLeader(g.Archive)

		for agentIdx := range g.Agents {
			for posIdx := range g.ObjectiveFunction.GetDimension() {
				velocity := g.W*g.Agents[agentIdx].Velocity[posIdx] +
					g.C1*rand.Float64()*(g.Agents[agentIdx].GetPersonalBest().Position[posIdx]-g.Agents[agentIdx].Result.Position[posIdx]) +
					g.C2*rand.Float64()*(leaderAgent.Position[posIdx]-g.Agents[agentIdx].Result.Position[posIdx])

				g.Agents[agentIdx].Velocity[posIdx] = velocity

				g.Agents[agentIdx].Result.Position[posIdx] = g.Agents[agentIdx].Result.Position[posIdx] + velocity
			}
		}

		// apply mutation
		g.Agents = mutation(g.Agents, l,
			g.NumberOfIter,
			g.NumberOfAgents,
			g.ObjectiveFunction.GetLowerBound(),
			g.ObjectiveFunction.GetUpperBound(),
			g.MutationRate)

		for agentIdx := range g.Agents {
			// check out of boundaries
			g.outOfBoundaries(g.Agents[agentIdx].Result.Position)
			// evaluate
			value, valuesWithKey, keys, penalty := g.ObjectiveFunction.Eval(g.Agents[agentIdx].Result.Position)
			g.Agents[agentIdx].Result.Value = value
			g.Agents[agentIdx].Result.Penalty = penalty
			g.Agents[agentIdx].Result.Key = keys
			g.Agents[agentIdx].Result.ValuesWithKey = valuesWithKey
		}

		// update repository
		g.Archive = g.updateRepository(g.Archive, g.Agents)
		if len(g.Archive) > g.ArchiveSize {
			exceeded := len(g.Archive) - g.ArchiveSize
			g.Archive = g.hypercube.deleteFromRepository(g.Archive, exceeded)
		}

		// update PersonalBest
		curDominatesPBest := make([]bool, len(g.Agents))
		pBestDominatesCur := make([]bool, len(g.Agents))

		for i := 0; i < len(g.Agents); i++ {
			if g.Agents[i].Result.Dominates(g.Agents[i].PersonalBest) {
				curDominatesPBest[i] = true
			}

			if !g.Agents[i].PersonalBest.Dominates(g.Agents[i].Result) {
				if rand.Float64() > 0.5 {
					pBestDominatesCur[i] = true
				}
			}
		}

		for i := 0; i < len(g.Agents); i++ {
			if curDominatesPBest[i] || pBestDominatesCur[i] {
				g.Agents[i].PersonalBest = g.Agents[i].Result.CopyAgent()
			}

		}

		channel <- struct {
			Progress                float64 `json:"progress"`
			NumberOfAgentsInArchive int     `json:"numberOfAgentsInArchive"`
			Type                    string  `json:"type"`
		}{
			Progress:                (float64(l+1) / float64(g.NumberOfIter)) * 100,
			NumberOfAgentsInArchive: len(g.Archive),
			Type:                    "multi",
		}

		l++
	}

	close(channel)

	return nil
}

func (g *MOPSOAlgorithm) initialization() {
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

func (g *MOPSOAlgorithm) outOfBoundaries(pos []float64) {
	for i := range pos {
		if pos[i] < g.ObjectiveFunction.GetLowerBound()[i] {
			pos[i] = g.ObjectiveFunction.GetLowerBound()[i]
		} else if pos[i] > g.ObjectiveFunction.GetUpperBound()[i] {
			pos[i] = g.ObjectiveFunction.GetUpperBound()[i]
		}
	}
}

func (g *MOPSOAlgorithm) GetResults() algorithms.Result {
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

type Hypercube struct {
	HypercubeLimits [][]float64 // Replaces Lower and Upper
	NumberOfGrids   int
	GridIndex       []int      // For each particle
	GridSubIndex    [][]int    // For each particle, their sub-indices
	Quality         []struct { // Quality of each hypercube
		Index int     // Hypercube identifier
		Value float64 // Quality value (10/number of particles)
	}
}

func (g *MOPSOAlgorithm) updateRepository(rep []*objectives.Result, agents []*ResultWithPersonalBest) []*objectives.Result {
	// Get results from agents
	agentResults := getResultsFromResultWithPersonalBest(agents)

	// Check domination between particles
	agentResults = objectives.DetermineDomination(agentResults)
	nonDominatedAgents := objectives.GetNonDominatedAgents(agentResults)

	// Add non-dominated particles to repository
	rep = objectives.MergeAgents(rep, nonDominatedAgents)

	// Check domination between all particles in repository
	rep = objectives.DetermineDomination(rep)
	rep = objectives.GetNonDominatedAgents(rep)

	// Update the grid
	g.hypercube.UpdateHyperCube(getResultsFromArchive(rep))

	return rep
}

func (h *Hypercube) UpdateHyperCube(values [][]float64) {
	numberOfObjectives := len(values)
	npar := len(values[0])

	// Initialize hypercube limits
	h.HypercubeLimits = make([][]float64, h.NumberOfGrids+1)
	for i := range h.HypercubeLimits {
		h.HypercubeLimits[i] = make([]float64, numberOfObjectives)
	}

	// Computing limits for each dimension
	for dim := 0; dim < numberOfObjectives; dim++ {
		minVal := slices.Min(values[dim])
		maxVal := slices.Max(values[dim])

		// Create linear space for this dimension
		for i := 0; i <= h.NumberOfGrids; i++ {
			h.HypercubeLimits[i][dim] = minVal + (float64(i)/float64(h.NumberOfGrids))*(maxVal-minVal)
		}
	}

	// Initialize grid indices
	h.GridIndex = make([]int, npar)
	h.GridSubIndex = make([][]int, npar)

	// Computing where each particle belongs
	for n := 0; n < npar; n++ {
		h.GridSubIndex[n] = make([]int, numberOfObjectives)

		// Find sub-indices for each dimension
		for d := 0; d < numberOfObjectives; d++ {
			// Find first index where particle value is less than or equal to limit
			for i := 0; i < h.NumberOfGrids; i++ {
				if values[d][n] <= h.HypercubeLimits[i+1][d] {
					h.GridSubIndex[n][d] = i + 1
					break
				}
			}
			if h.GridSubIndex[n][d] == 0 {
				h.GridSubIndex[n][d] = 1
			}
		}

		// Convert sub-indices to grid index
		gridSize := make([]int, numberOfObjectives)
		for i := range gridSize {
			gridSize[i] = h.NumberOfGrids
		}

		// Adjust indices to be 0-based for Sub2Index
		adjustedIndices := make([]int, numberOfObjectives)
		for i := 0; i < numberOfObjectives; i++ {
			// Subtract 1 to convert from 1-based to 0-based indexing
			adjustedIndices[i] = h.GridSubIndex[n][i] - 1
		}

		h.GridIndex[n] = util.Sub2Index(gridSize, adjustedIndices...)
	}

	// Calculate quality for each occupied hypercube
	gridCounts := make(map[int]int)
	for _, idx := range h.GridIndex {
		gridCounts[idx]++
	}

	h.Quality = make([]struct {
		Index int
		Value float64
	}, len(gridCounts))
	i := 0
	for idx, count := range gridCounts {
		h.Quality[i].Index = idx
		h.Quality[i].Value = 10.0 / float64(count)
		i++
	}
}

func (h *Hypercube) SelectLeader(rep []*objectives.Result) *objectives.Result {
	// Compute cumulative probabilities based on quality
	prob := make([]float64, len(h.Quality))
	prob[0] = h.Quality[0].Value
	for i := 1; i < len(h.Quality); i++ {
		prob[i] = prob[i-1] + h.Quality[i].Value
	}

	// Perform roulette wheel selection
	randVal := rand.Float64() * prob[len(prob)-1]
	selectedHypercube := -1
	for i, p := range prob {
		if randVal <= p {
			selectedHypercube = h.Quality[i].Index
			break
		}
	}

	// Find all particles in the selected hypercube
	selectedIndices := make([]int, 0)
	for i, gridIdx := range h.GridIndex {
		if gridIdx == selectedHypercube {
			selectedIndices = append(selectedIndices, i)
		}
	}

	// Randomly select one particle from the selected hypercube
	if len(selectedIndices) == 0 {
		return nil
	}
	selectedIndex := selectedIndices[rand.Intn(len(selectedIndices))]
	return rep[selectedIndex]
}

func (h *Hypercube) deleteFromRepository(rep []*objectives.Result, exceeded int) []*objectives.Result {
	// Compute crowding distances
	crowding := make([]float64, len(rep))
	for m := 0; m < len(rep[0].Value); m++ {
		// Extract the m-th objective values
		mFit := make([]float64, len(rep))
		for i, r := range rep {
			mFit[i] = r.Value[m]
		}

		// Sort the values and get indices
		sortedIndices := make([]int, len(mFit))
		for i := range sortedIndices {
			sortedIndices[i] = i
		}

		slices.SortFunc(sortedIndices, func(i, j int) int {
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

		// Compute distances
		mUp := make([]float64, len(mFit))
		mDown := make([]float64, len(mFit))
		for i := 0; i < len(mFit); i++ {
			if i == 0 {
				mDown[i] = math.Inf(1)
			} else {
				mDown[i] = mFit[sortedIndices[i-1]]
			}
			if i == len(mFit)-1 {
				mUp[i] = math.Inf(1)
			} else {
				mUp[i] = mFit[sortedIndices[i+1]]
			}
		}

		maxFit := slices.Max(mFit)
		minFit := slices.Min(mFit)
		for i := 0; i < len(mFit); i++ {
			distance := (mUp[i] - mDown[i]) / (maxFit - minFit)
			if math.IsNaN(distance) {
				distance = math.Inf(1)
			}
			crowding[sortedIndices[i]] += distance
		}
	}

	// Delete the extra particles with the smallest crowding distances
	sortedCrowdingIndices := make([]int, len(crowding))
	for i := range sortedCrowdingIndices {
		sortedCrowdingIndices[i] = i
	}
	slices.SortFunc(sortedCrowdingIndices, func(i, j int) int {
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

	delIndices := sortedCrowdingIndices[:exceeded]
	newRep := make([]*objectives.Result, 0, len(rep)-exceeded)
	for i, r := range rep {
		if !slices.Contains(delIndices, i) {
			newRep = append(newRep, r)
		}
	}

	// Update the hypercube grid
	h.UpdateHyperCube(getResultsFromArchive(newRep))
	return newRep
}

func mutation(pos []*ResultWithPersonalBest, curIter, maxIter, numberOfAgents int, upper, lower []float64, mutationRate float64) []*ResultWithPersonalBest {
	// Sub-divide the swarm in three parts
	fract := float64(numberOfAgents)/3.0 - math.Floor(float64(numberOfAgents)/3.0)

	roundSubSize := int(math.Round(float64(numberOfAgents) / 3.0))

	var subSizes []int
	if fract < 0.5 {
		subSizes = []int{
			int(math.Ceil(float64(numberOfAgents) / 3.0)),
			roundSubSize,
			roundSubSize,
		}
	} else {
		subSizes = []int{
			roundSubSize,
			roundSubSize,
			int(math.Floor(float64(numberOfAgents) / 3.0)),
		}
	}

	// Calculate cumulative sizes
	cumSizes := make([]int, 3)
	cumSizes[0] = subSizes[0]
	for i := 1; i < len(subSizes); i++ {
		cumSizes[i] = cumSizes[i-1] + subSizes[i]
	}

	// First part: no mutation
	// Second part: uniform mutation
	nmut := int(math.Round(mutationRate * float64(subSizes[1])))
	if nmut > 0 {
		// Generate random indices for mutation
		for i := 0; i < nmut; i++ {
			idx := cumSizes[0] + rand.Intn(subSizes[1])
			for j := range pos[idx].Result.Position {
				pos[idx].Result.Position[j] = lower[j] + rand.Float64()*(upper[j]-lower[j])
			}
		}
	}

	// Third part: non-uniform mutation
	perMut := math.Pow(1.0-float64(curIter)/float64(maxIter), 5.0*float64(len(pos[0].Result.Position)))
	nmut = int(math.Round(perMut * float64(subSizes[2])))
	if nmut > 0 {
		// Generate random indices for mutation
		for i := 0; i < nmut; i++ {
			idx := cumSizes[1] + rand.Intn(subSizes[2])
			for j := range pos[idx].Result.Position {
				pos[idx].Result.Position[j] = lower[j] + rand.Float64()*(upper[j]-lower[j])
			}
		}
	}

	return pos
}
