package mogwo

import (
	"golang-moaha-construction/internal/algorithms"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
	"golang-moaha-construction/internal/util"
	"math"
	"math/rand"
	"slices"
	"sync"
)

const NameType algorithms.AlgorithmType = "MOGWO"

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

type MOGWOAlgorithm struct {
	NumberOfAgents    int
	NumberOfIter      int
	Agents            []*objectives.Result
	AParam            float64
	AlphaWolf         *objectives.Result
	BetaWolf          *objectives.Result
	GammaWolf         *objectives.Result
	ArchiveSize       int
	Archive           []*objectives.Result
	ObjectiveFunction objectives.Problem
	// for multi-objective hypercubes
	NumberOfGrids int
	Gamma         float64
	Alpha         float64
	Beta          float64
	hypercube     Hypercube
}

type Config struct {
	NumberOfAgents int
	NumberOfIter   int
	AParam         float64
	ArchiveSize    int
	NumberOfGrids  int
	Gamma          float64
	Alpha          float64
	Beta           float64
}

func Create(
	problem objectives.Problem,
	configs Config,
) (*MOGWOAlgorithm, error) {

	return &MOGWOAlgorithm{
		NumberOfAgents:    configs.NumberOfAgents,
		NumberOfIter:      configs.NumberOfIter,
		AParam:            configs.AParam,
		ArchiveSize:       configs.ArchiveSize,
		ObjectiveFunction: problem,
		NumberOfGrids:     configs.NumberOfGrids,
		Gamma:             configs.Gamma,
		Alpha:             configs.Alpha,
		Beta:              configs.Beta,
		hypercube: Hypercube{
			NumberOfGrids: configs.NumberOfGrids,
			Alpha:         configs.Alpha,
		},
	}, nil
}

func (g *MOGWOAlgorithm) reset() {
	g.Agents = make([]*objectives.Result, g.NumberOfAgents)
	g.Archive = make([]*objectives.Result, 0, g.ArchiveSize)
}

func (g *MOGWOAlgorithm) Type() data.TypeProblem {
	return data.Single
}

func (g *MOGWOAlgorithm) Run() error {
	g.reset()

	// initialization
	g.initialization()

	g.Agents = objectives.DetermineDomination(g.Agents)
	g.Archive = objectives.GetNonDominatedAgents(g.Agents)
	g.hypercube.UpdateHyperCube(getResultsFromArchive(g.Archive))

	for _, val := range g.Archive {
		g.hypercube.getGridIndex(val)
	}

	l := 0
	a := g.AParam

	for l < g.NumberOfIter {
		a = 2.0 - float64(l)*(2.0/float64(g.NumberOfIter))

		for agentIdx := range g.Agents {
			// Choose the alpha, beta, and gamma grey wolves
			if len(g.Archive) > 0 {
				gammaLeader := selectLeader(g.Archive, g.Beta)
				if gammaLeader != nil {
					g.GammaWolf = gammaLeader
				}

				betaLeader := selectLeader(g.Archive, g.Beta)
				if betaLeader != nil {
					g.BetaWolf = betaLeader
				}

				alphaLeader := selectLeader(g.Archive, g.Beta)
				if alphaLeader != nil {
					g.AlphaWolf = alphaLeader
				}

				// If there are less than three solutions in the least crowded
				// hypercube, the second least crowded hypercube is also found
				// to choose other leaders from.
				var rep2 []*objectives.Result
				if len(g.Archive) > 1 {
					rep2 = make([]*objectives.Result, 0)
					for _, archiveItem := range g.Archive {
						// Check if the position is different from Delta (GammaWolf in Go implementation)
						if archiveItem != g.GammaWolf {
							rep2 = append(rep2, archiveItem)
						}
					}

					if len(rep2) > 0 {
						betaLeader = selectLeader(rep2, g.Beta)
						if betaLeader != nil {
							g.BetaWolf = betaLeader
						}
					}
				}

				// This scenario is the same if the second least crowded hypercube
				// has one solution, so the alpha leader should be chosen from the
				// third least crowded hypercube.
				if len(g.Archive) > 2 && len(rep2) > 0 {
					rep3 := make([]*objectives.Result, 0)
					for _, archiveItem := range rep2 {
						// Check if the position is different from Beta
						if archiveItem != g.BetaWolf {
							rep3 = append(rep3, archiveItem)
						}
					}

					if len(rep3) > 0 {
						alphaLeader = selectLeader(rep3, g.Beta)
						if alphaLeader != nil {
							g.AlphaWolf = alphaLeader
						}
					}
				}
			}

			for posIdx := range g.Agents[agentIdx].Position {
				// Alpha
				r1 := rand.Float64()
				r2 := rand.Float64()
				A := 2*a*r1 - a
				C := 2 * r2
				D := math.Abs(C*g.AlphaWolf.Position[posIdx] - g.Agents[agentIdx].Position[posIdx])
				XAlpha := g.AlphaWolf.Position[posIdx] - A*D

				// Beta
				r1 = rand.Float64()
				r2 = rand.Float64()
				A = 2*a*r1 - a
				C = 2 * r2
				D = math.Abs(C*g.BetaWolf.Position[posIdx] - g.Agents[agentIdx].Position[posIdx])
				XBeta := g.BetaWolf.Position[posIdx] - A*D

				// Gamma
				r1 = rand.Float64()
				r2 = rand.Float64()
				A = 2*a*r1 - a
				C = 2 * r2
				D = math.Abs(C*g.GammaWolf.Position[posIdx] - g.Agents[agentIdx].Position[posIdx])
				XGamma := g.GammaWolf.Position[posIdx] - A*D

				g.Agents[agentIdx].Position[posIdx] = (XAlpha + XBeta + XGamma) / 3
			}
			// check out of boundaries
			g.outOfBoundaries(g.Agents[agentIdx].Position)

			// evaluate
			value, valuesWithKey, keys, penalty := g.ObjectiveFunction.Eval(g.Agents[agentIdx].Position)
			g.Agents[agentIdx].Value = value
			g.Agents[agentIdx].Penalty = penalty
			g.Agents[agentIdx].Key = keys
			g.Agents[agentIdx].ValuesWithKey = valuesWithKey
		}

		newSolutions := objectives.DetermineDomination(g.Agents)
		newNonDominatedPop := objectives.GetNonDominatedAgents(newSolutions)

		newSolutions = objectives.DetermineDomination(objectives.MergeAgents(newNonDominatedPop, g.Archive))
		g.Archive = objectives.GetNonDominatedAgents(newSolutions)
		g.hypercube.UpdateHyperCube(getResultsFromArchive(g.Archive))

		for _, val := range g.Archive {
			g.hypercube.getGridIndex(val)
		}

		if len(g.Archive) > g.ArchiveSize {
			exceeded := len(g.Archive) - g.ArchiveSize
			g.Archive = removeExtraInArchive(g.Archive, exceeded, g.Gamma)
			g.hypercube.UpdateHyperCube(getResultsFromArchive(g.Archive))
		}

		l++
	}

	return nil
}

func (g *MOGWOAlgorithm) RunWithChannel(doneChan chan<- struct{}, channel chan<- any) error {

	g.reset()

	// initialization
	g.initialization()

	g.Agents = objectives.DetermineDomination(g.Agents)
	g.Archive = objectives.GetNonDominatedAgents(g.Agents)
	g.hypercube.UpdateHyperCube(getResultsFromArchive(g.Archive))

	for _, val := range g.Archive {
		g.hypercube.getGridIndex(val)
	}

	l := 0
	a := g.AParam

	for l < g.NumberOfIter {
		a = 2.0 - float64(l)*(2.0/float64(g.NumberOfIter))

		for agentIdx := range g.Agents {
			// Choose the alpha, beta, and gamma grey wolves
			if len(g.Archive) > 0 {
				gammaLeader := selectLeader(g.Archive, g.Beta)
				if gammaLeader != nil {
					g.GammaWolf = gammaLeader
				}

				betaLeader := selectLeader(g.Archive, g.Beta)
				if betaLeader != nil {
					g.BetaWolf = betaLeader
				}

				alphaLeader := selectLeader(g.Archive, g.Beta)
				if alphaLeader != nil {
					g.AlphaWolf = alphaLeader
				}

				// If there are less than three solutions in the least crowded
				// hypercube, the second least crowded hypercube is also found
				// to choose other leaders from.
				var rep2 []*objectives.Result
				if len(g.Archive) > 1 {
					rep2 = make([]*objectives.Result, 0)
					for _, archiveItem := range g.Archive {
						// Check if the position is different from Delta (GammaWolf in Go implementation)
						if archiveItem != g.GammaWolf {
							rep2 = append(rep2, archiveItem)
						}
					}

					if len(rep2) > 0 {
						betaLeader = selectLeader(rep2, g.Beta)
						if betaLeader != nil {
							g.BetaWolf = betaLeader
						}
					}
				}

				// This scenario is the same if the second least crowded hypercube
				// has one solution, so the alpha leader should be chosen from the
				// third least crowded hypercube.
				if len(g.Archive) > 2 && len(rep2) > 0 {
					rep3 := make([]*objectives.Result, 0)
					for _, archiveItem := range rep2 {
						// Check if the position is different from Beta
						if archiveItem != g.BetaWolf {
							rep3 = append(rep3, archiveItem)
						}
					}

					if len(rep3) > 0 {
						alphaLeader = selectLeader(rep3, g.Beta)
						if alphaLeader != nil {
							g.AlphaWolf = alphaLeader
						}
					}
				}
			}

			for posIdx := range g.Agents[agentIdx].Position {
				// Alpha
				r1 := rand.Float64()
				r2 := rand.Float64()
				A := 2*a*r1 - a
				C := 2 * r2
				D := math.Abs(C*g.AlphaWolf.Position[posIdx] - g.Agents[agentIdx].Position[posIdx])
				XAlpha := g.AlphaWolf.Position[posIdx] - A*D

				// Beta
				r1 = rand.Float64()
				r2 = rand.Float64()
				A = 2*a*r1 - a
				C = 2 * r2
				D = math.Abs(C*g.BetaWolf.Position[posIdx] - g.Agents[agentIdx].Position[posIdx])
				XBeta := g.BetaWolf.Position[posIdx] - A*D

				// Gamma
				r1 = rand.Float64()
				r2 = rand.Float64()
				A = 2*a*r1 - a
				C = 2 * r2
				D = math.Abs(C*g.GammaWolf.Position[posIdx] - g.Agents[agentIdx].Position[posIdx])
				XGamma := g.GammaWolf.Position[posIdx] - A*D

				g.Agents[agentIdx].Position[posIdx] = (XAlpha + XBeta + XGamma) / 3
			}
			// check out of boundaries
			g.outOfBoundaries(g.Agents[agentIdx].Position)

			// evaluate
			value, valuesWithKey, keys, penalty := g.ObjectiveFunction.Eval(g.Agents[agentIdx].Position)
			g.Agents[agentIdx].Value = value
			g.Agents[agentIdx].Penalty = penalty
			g.Agents[agentIdx].Key = keys
			g.Agents[agentIdx].ValuesWithKey = valuesWithKey
		}

		newSolutions := objectives.DetermineDomination(g.Agents)
		newNonDominatedPop := objectives.GetNonDominatedAgents(newSolutions)

		newSolutions = objectives.DetermineDomination(objectives.MergeAgents(newNonDominatedPop, g.Archive))
		g.Archive = objectives.GetNonDominatedAgents(newSolutions)
		g.hypercube.UpdateHyperCube(getResultsFromArchive(g.Archive))

		for _, val := range g.Archive {
			g.hypercube.getGridIndex(val)
		}

		if len(g.Archive) > g.ArchiveSize {
			exceeded := len(g.Archive) - g.ArchiveSize
			g.Archive = removeExtraInArchive(g.Archive, exceeded, g.Gamma)
			g.hypercube.UpdateHyperCube(getResultsFromArchive(g.Archive))
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

func (g *MOGWOAlgorithm) initialization() {

	vals := make([]float64, g.ObjectiveFunction.NumberOfObjectives())
	for i := 0; i < g.ObjectiveFunction.NumberOfObjectives(); i++ {
		if g.ObjectiveFunction.FindMin() {
			vals[i] = math.MaxFloat64
		} else {
			vals[i] = math.MinInt64
		}
	}

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
			newAgent := &objectives.Result{
				Idx:      agentIdx,
				Position: positions,
			}

			value, valuesWithKey, keys, penalty := g.ObjectiveFunction.Eval(positions)
			newAgent.Value = value
			newAgent.ValuesWithKey = valuesWithKey
			newAgent.Penalty = penalty
			newAgent.Key = keys

			g.Agents[agentIdx] = newAgent
		}(agentIdx)
	}
	wg.Wait()
}

func (g *MOGWOAlgorithm) outOfBoundaries(pos []float64) {
	for i := range pos {
		if pos[i] < g.ObjectiveFunction.GetLowerBound()[i] {
			pos[i] = g.ObjectiveFunction.GetLowerBound()[i]
		} else if pos[i] > g.ObjectiveFunction.GetUpperBound()[i] {
			pos[i] = g.ObjectiveFunction.GetUpperBound()[i]
		}
	}
}

func (g *MOGWOAlgorithm) GetResults() algorithms.Result {
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
	Lower         [][]float64
	Upper         [][]float64
	NumberOfGrids int
	Alpha         float64
}

// UpdateHyperCube create a new hypercube with the given values (n*m)
// with n = number of objectives and m = number of values of that objective
func (h *Hypercube) UpdateHyperCube(values [][]float64) {
	numberOfObjectives := len(values)

	h.Lower = make([][]float64, numberOfObjectives)
	h.Upper = make([][]float64, numberOfObjectives)

	for i := 0; i < numberOfObjectives; i++ {
		minCj := slices.Min(values[i])
		maxCj := slices.Max(values[i])

		dcj := (maxCj - minCj) * h.Alpha

		minCj = minCj - dcj
		maxCj = maxCj + dcj

		gx := util.LinSpace(minCj, maxCj, h.NumberOfGrids-1)

		lower := make([]float64, len(gx)+1)
		upper := make([]float64, len(gx)+1)

		for j := 0; j < len(gx); j++ {
			lower[j+1] = gx[j]
			upper[j] = gx[j]
		}

		lower[0] = math.Inf(-1)
		h.Lower[i] = lower

		upper[len(lower)-1] = math.Inf(1)
		h.Upper[i] = upper
	}
}

func (h *Hypercube) getGridIndex(agentResult *objectives.Result) {
	numberOfObjectives := len(agentResult.Value)

	// Create sub-indices array to store the grid cell indices for each dimension
	index := make([]int, numberOfObjectives)

	// Find the grid cell index for each dimension
	for i := 0; i < numberOfObjectives; i++ {
		index[i] = util.FindLess(h.Upper[i], agentResult.Value[i])
	}

	var res int

	size := make([]int, numberOfObjectives)
	for i := range size {
		size[i] = h.NumberOfGrids
	}
	res = util.Sub2Index(size, index...)

	agentResult.GridIndex = res
	agentResult.GridSubIndex = index
}

func getOccupiedCells(archive []*objectives.Result) ([]int, []int) {
	// Use a map to count members in each cell
	cellCountMap := make(map[int]int)

	// Count occurrences of each grid index
	for _, res := range archive {
		cellCountMap[res.GridIndex]++
	}

	// Extract unique cell indices and their counts
	occCellIndex := make([]int, 0, len(cellCountMap))
	occCellMemberCount := make([]int, 0, len(cellCountMap))

	for cellIdx, count := range cellCountMap {
		occCellIndex = append(occCellIndex, cellIdx)
		occCellMemberCount = append(occCellMemberCount, count)
	}

	return occCellIndex, occCellMemberCount
}

func removeExtraInArchive(archive []*objectives.Result, exceeded int, gamma float64) []*objectives.Result {
	if gamma == 0 {
		gamma = 1
	}

	for k := 0; k < exceeded; k++ {
		// Get occupied cells and their member counts
		occCellIndex, occCellMemberCount := getOccupiedCells(archive)

		// Calculate probabilities based on member counts raised to gamma
		p := make([]float64, len(occCellMemberCount))
		sum := 0.0
		for i, count := range occCellMemberCount {
			p[i] = math.Pow(float64(count), gamma)
			sum += p[i]
		}

		// Normalize probabilities
		for i := range p {
			p[i] = p[i] / sum
		}

		// Select a cell using roulette wheel selection
		selectedCellIndex := occCellIndex[util.RouletteWheelSelection(p)]

		// Find members in the selected cell
		selectedCellMembers := make([]int, 0)
		for i, res := range archive {
			if res.GridIndex == selectedCellIndex {
				selectedCellMembers = append(selectedCellMembers, i)
			}
		}

		// Randomly select one member to remove
		n := len(selectedCellMembers)
		if n == 0 {
			continue
		}
		selectedMemberIndex := selectedCellMembers[rand.Intn(n)]

		// Remove the selected member from the archive
		archive = append(archive[:selectedMemberIndex], archive[selectedMemberIndex+1:]...)
	}

	return archive
}

// SelectLeader selects a leader from the repository based on grid indices.
// It takes a repository of solutions and a beta parameter (default 1).
// The function returns the selected leader.
func selectLeader(rep []*objectives.Result, beta float64) *objectives.Result {
	// Set default value for beta if not provided
	if beta == 0 {
		beta = 1
	}

	// Get occupied cells and their member counts
	occCellIndex, occCellMemberCount := getOccupiedCells(rep)

	// Calculate probabilities based on member counts raised to -beta
	p := make([]float64, len(occCellMemberCount))
	sum := 0.0
	for i, count := range occCellMemberCount {
		p[i] = math.Pow(float64(count), -beta)
		sum += p[i]
	}

	// Normalize probabilities
	for i := range p {
		p[i] = p[i] / sum
	}

	// Select a cell using roulette wheel selection
	selectedCellIndex := occCellIndex[util.RouletteWheelSelection(p)]

	// Find members in the selected cell
	selectedCellMembers := make([]int, 0)
	for i, res := range rep {
		if res.GridIndex == selectedCellIndex {
			selectedCellMembers = append(selectedCellMembers, i)
		}
	}

	// Randomly select one member from the cell
	n := len(selectedCellMembers)
	if n == 0 {
		return nil
	}
	selectedMemberIndex := selectedCellMembers[rand.Intn(n)]

	// Return the selected member
	return rep[selectedMemberIndex]
}
