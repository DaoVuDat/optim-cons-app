package moaha

import (
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

const NameType algorithms.AlgorithmType = "MOAHA"

type MOAHAAlgorithm struct {
	NumberOfAgents    int
	NumberOfIter      int
	Agents            []*objectives.Result
	BestResult        *objectives.Result
	Convergence       []float64
	ObjectiveFunction objectives.Problem
	ArchiveSize       int
	Archive           []*objectives.Result
}

type Configs struct {
	NumAgents     int
	NumIterations int
	ArchiveSize   int
}

func Create(
	problem objectives.Problem,
	configs Configs,
) (*MOAHAAlgorithm, error) {

	convergence := make([]float64, configs.NumIterations)
	agents := make([]*objectives.Result, configs.NumAgents)

	archive := make([]*objectives.Result, 0, configs.ArchiveSize)

	return &MOAHAAlgorithm{
		NumberOfAgents:    configs.NumAgents,
		NumberOfIter:      configs.NumIterations,
		Convergence:       convergence,
		ObjectiveFunction: problem,
		Agents:            agents,
		Archive:           archive,
		ArchiveSize:       configs.ArchiveSize,
	}, nil
}

func (a *MOAHAAlgorithm) Type() data.TypeProblem {
	return data.Multi
}

func (a *MOAHAAlgorithm) Run() error {
	dimensions := a.ObjectiveFunction.GetDimension()

	// initialization
	a.initialization()

	a.Agents = objectives.DetermineDomination(a.Agents)
	a.Archive = objectives.GetNonDominatedAgents(a.Agents)

	l := 0

	bar := progressbar.Default(int64(a.NumberOfIter))
	//var wg sync.WaitGroup

	visitTable := initializeNMMatrix(a.NumberOfAgents, a.NumberOfAgents)

	for l < a.NumberOfIter {
		newPop := make([]*objectives.Result, 0)
		agents, paretoFront := objectives.NonDominatedSort(a.Agents)

		a.Agents = agents

		// direct vector
		directVector := initializeNMMatrix(a.NumberOfAgents, dimensions)

		//wg.Add(a.NumberOfAgents)
		for agentIdx := range a.Agents {

			r := rand.Float64()

			//fmt.Println("")
			if r < 1.0/3.0 {
				// diagonal flight
				randDim := util.RandN(dimensions)
				randNum := 0
				if dimensions > 3 {
					randNum = rand.Intn(dimensions - 1)
				} else {
					randNum = rand.Intn(dimensions)
				}

				//test := []int{19, 28, 27, 7, 29, 20, 9, 4, 12, 15, 21, 6, 13, 25, 2, 23, 8, 26, 30, 1, 5, 14, 17, 24, 16, 10, 18, 22, 11, 3}

				for i := 0; i < randNum; i++ {
					idx := randDim[i]
					//idx = test[i] - 1 // test

					directVector[agentIdx][idx] = 1
				}
			} else if r > 2.0/3.0 {
				// omnidirectional flight
				for i := 0; i < dimensions; i++ {
					directVector[agentIdx][i] = 1
				}
			} else {
				// axial flight
				randNum := rand.Intn(dimensions)
				directVector[agentIdx][randNum] = 1
			}

			r = rand.Float64()
			//fmt.Println()
			if r < 0.5 {
				// guided foraging
				a.guidedForaging(visitTable, directVector, agentIdx, paretoFront, newPop)
			} else {
				// territory foraging
				a.territoryForaging(visitTable, directVector, agentIdx, paretoFront, newPop)
			}

		}

		//wg.Wait()

		// migration foraging
		if l%(a.NumberOfAgents*2) == 0 {
			a.Agents, paretoFront = objectives.NonDominatedSort(a.Agents)

			for _, idx := range paretoFront[len(paretoFront)-1] {

				for i := range a.Agents[idx].Position {
					a.Agents[idx].Position[i] =
						a.ObjectiveFunction.GetLowerBound()[i] + rand.Float64()*
							(a.ObjectiveFunction.GetUpperBound()[i]-a.ObjectiveFunction.GetLowerBound()[i])
				}
				// evaluate
				value, constraints, penalty := a.ObjectiveFunction.Eval(a.Agents[idx].Position)

				a.Agents[idx].Value = value
				a.Agents[idx].Constraints = constraints
				a.Agents[idx].Penalty = penalty

				for i := range visitTable[idx] {
					visitTable[idx][i] += 1
				}

				maxVals := maxRowMatrix(visitTable)
				for i := range visitTable[idx] {
					if i == idx {
						continue
					}
					visitTable[i][idx] = maxVals[i] + 1
				}
			}

		}

		// Determine Domination with a.Agents and newPop
		newSolutions := objectives.DetermineDomination(objectives.MergeAgents(a.Agents, newPop))
		// Get Non-Dominated -> newNonDominatedPop
		newNonDominatedPop := objectives.GetNonDominatedAgents(newSolutions)

		// Determine Domination with newDominatedPop and a.Archive
		newSolutions = objectives.DetermineDomination(objectives.MergeAgents(newNonDominatedPop, a.Archive))
		// Get Non-Dominated -> a.Archive
		a.Archive = objectives.GetNonDominatedAgents(newSolutions)

		if len(a.Archive) > a.ArchiveSize {
			a.Archive = objectives.DECD(a.Archive, len(a.Archive)-a.ArchiveSize)
		}

		//bar.Describe(fmt.Sprintf("Iter %d: %e", l+1, a.BestResult.Value[0]))
		bar.Add(1)

		l++
	}

	return nil
}

func (a *MOAHAAlgorithm) RunWithChannel(doneChan chan<- struct{}, channel chan<- any) error {
	dimensions := a.ObjectiveFunction.GetDimension()

	// initialization
	a.initialization()

	a.Agents = objectives.DetermineDomination(a.Agents)
	a.Archive = objectives.GetNonDominatedAgents(a.Agents)

	l := 0

	visitTable := initializeNMMatrix(a.NumberOfAgents, a.NumberOfAgents)

	for l < a.NumberOfIter {
		newPop := make([]*objectives.Result, 0)
		agents, paretoFront := objectives.NonDominatedSort(a.Agents)

		a.Agents = agents

		// direct vector
		directVector := initializeNMMatrix(a.NumberOfAgents, dimensions)

		//wg.Add(a.NumberOfAgents)
		for agentIdx := range a.Agents {

			r := rand.Float64()

			//fmt.Println("")
			if r < 1.0/3.0 {
				// diagonal flight
				randDim := util.RandN(dimensions)
				randNum := 0
				if dimensions > 3 {
					randNum = rand.Intn(dimensions - 1)
				} else {
					randNum = rand.Intn(dimensions)
				}

				for i := 0; i < randNum; i++ {
					idx := randDim[i]

					directVector[agentIdx][idx] = 1
				}
			} else if r > 2.0/3.0 {
				// omnidirectional flight
				for i := 0; i < dimensions; i++ {
					directVector[agentIdx][i] = 1
				}
			} else {
				// axial flight
				randNum := rand.Intn(dimensions)
				directVector[agentIdx][randNum] = 1
			}

			r = rand.Float64()
			if r < 0.5 {
				// guided foraging
				a.guidedForaging(visitTable, directVector, agentIdx, paretoFront, newPop)
			} else {
				// territory foraging
				a.territoryForaging(visitTable, directVector, agentIdx, paretoFront, newPop)
			}

		}

		// migration foraging
		if l%(a.NumberOfAgents*2) == 0 {
			a.Agents, paretoFront = objectives.NonDominatedSort(a.Agents)

			for _, idx := range paretoFront[len(paretoFront)-1] {

				for i := range a.Agents[idx].Position {
					a.Agents[idx].Position[i] =
						a.ObjectiveFunction.GetLowerBound()[i] + rand.Float64()*
							(a.ObjectiveFunction.GetUpperBound()[i]-a.ObjectiveFunction.GetLowerBound()[i])
				}
				// evaluate
				value, constraints, penalty := a.ObjectiveFunction.Eval(a.Agents[idx].Position)

				a.Agents[idx].Value = value
				a.Agents[idx].Constraints = constraints
				a.Agents[idx].Penalty = penalty

				for i := range visitTable[idx] {
					visitTable[idx][i] += 1
				}

				maxVals := maxRowMatrix(visitTable)
				for i := range visitTable[idx] {
					if i == idx {
						continue
					}
					visitTable[i][idx] = maxVals[i] + 1
				}
			}

		}

		// Determine Domination with a.Agents and newPop
		newSolutions := objectives.DetermineDomination(objectives.MergeAgents(a.Agents, newPop))
		// Get Non-Dominated -> newNonDominatedPop
		newNonDominatedPop := objectives.GetNonDominatedAgents(newSolutions)

		// Determine Domination with newDominatedPop and a.Archive
		newSolutions = objectives.DetermineDomination(objectives.MergeAgents(newNonDominatedPop, a.Archive))
		// Get Non-Dominated -> a.Archive
		a.Archive = objectives.GetNonDominatedAgents(newSolutions)

		if len(a.Archive) > a.ArchiveSize {
			a.Archive = objectives.DECD(a.Archive, len(a.Archive)-a.ArchiveSize)
		}

		channel <- struct {
			Progress                float64 `json:"progress"`
			NumberOfAgentsInArchive int     `json:"numberOfAgentsInArchive"`
			Type                    string  `json:"type"`
		}{
			Progress:                (float64(l+1) / float64(a.NumberOfIter)) * 100,
			NumberOfAgentsInArchive: len(a.Archive),
			Type:                    "multi",
		}

		l++
	}
	close(channel)

	return nil
}

func (a *MOAHAAlgorithm) guidedForaging(visitTable [][]float64, directVector [][]float64, agentIdx int, paretoFront [][]int, tPop []*objectives.Result) {
	nonDominatedMUT := make([]*objectives.Result, 0)

	vals := visitTable[agentIdx]
	maxVal := -math.MaxFloat64
	maxValIdxs := make([]int, 0)
	for i := 0; i < len(vals); i++ {
		if i == agentIdx {
			continue
		}

		if vals[i] > maxVal {
			maxVal = vals[i]
			maxValIdxs = []int{i}
		} else if vals[i] == maxVal {
			maxValIdxs = append(maxValIdxs, i)
		}
	}

	targetFoodIdx := 0
	if len(maxValIdxs) >= 2 {

		candidates := make([]*objectives.Result, 0)
		for _, idx := range maxValIdxs {
			candidates = append(candidates, a.Agents[idx])
		}

		candidates = objectives.DetermineDomination(candidates)
		for _, candidate := range candidates {
			if !candidate.Dominated {
				nonDominatedMUT = append(nonDominatedMUT, candidate)
			}
		}

		candidateIdx := rand.Intn(len(nonDominatedMUT))
		//fmt.Println()
		targetFoodIdx = candidateIdx
	} else if len(maxValIdxs) == 1 {
		targetFoodIdx = maxValIdxs[0]
	} else {
		panic("len(maxValIdxs) = 0")
	}

	r := rand.NormFloat64()
	newPos := make([]float64, a.ObjectiveFunction.GetDimension())
	for i := 0; i < a.ObjectiveFunction.GetDimension(); i++ {
		newPos[i] = a.Agents[targetFoodIdx].Position[i] + r*math.Round(directVector[agentIdx][i])*
			(a.Agents[agentIdx].Position[i]-a.Agents[targetFoodIdx].Position[i])
	}

	a.outOfBoundaries(newPos)

	newAgent := a.Agents[agentIdx].CopyAgent()

	value, constraints, penalty := a.ObjectiveFunction.Eval(newPos)
	newAgent.Position = newPos
	newAgent.Value = value
	newAgent.Constraints = constraints
	newAgent.Penalty = penalty

	// Sanity check the index of current agent
	frontIdx := 0
	for i := range paretoFront {
		if slices.Contains(paretoFront[i], agentIdx) {
			frontIdx = i
			break
		}
	}

	// compare the new value to other agents in the same front
	dominatedFlag := 0
	for _, v := range paretoFront[frontIdx] {
		if newAgent.Dominates(a.Agents[v]) {
			dominatedFlag = 1
			break
		} else if a.Agents[v].Dominates(newAgent) {
			dominatedFlag = -1
			break
		}
	}

	newR := rand.Float64()
	if dominatedFlag == 1 || (dominatedFlag == 0 && newR > 0.5) {
		tPop = append(tPop, a.Agents[agentIdx].CopyAgent())

		a.Agents[agentIdx] = newAgent.CopyAgent()

		for i := range visitTable[agentIdx] {
			if i == targetFoodIdx {
				visitTable[agentIdx][i] = 0
			}
			visitTable[agentIdx][i] += 1
		}

		maxVals := maxRowMatrix(visitTable)
		for i := range visitTable[agentIdx] {
			if i == agentIdx {
				continue
			}
			visitTable[i][agentIdx] = maxVals[i] + 1
		}

	} else {
		for i := range visitTable[agentIdx] {
			if i == targetFoodIdx {
				visitTable[agentIdx][i] = 0
			}
			visitTable[agentIdx][i] += 1
		}

		tPop = append(tPop, newAgent.CopyAgent())
	}
}

func (a *MOAHAAlgorithm) territoryForaging(visitTable [][]float64, directVector [][]float64, agentIdx int, paretoFront [][]int, tPop []*objectives.Result) {
	r1 := rand.Float64()
	r2 := rand.NormFloat64()
	newPos := make([]float64, a.ObjectiveFunction.GetDimension())
	if r1 > 0.5 {
		for i := 0; i < a.ObjectiveFunction.GetDimension(); i++ {
			newPos[i] = a.Agents[agentIdx].Position[i] + r2*math.Round(directVector[agentIdx][i])*a.Agents[agentIdx].Position[i]
		}
	} else {
		// randomly selected from archive
		selectedIdx := rand.Intn(len(a.Archive))
		agentInArchive := a.Archive[selectedIdx]
		for i := 0; i < a.ObjectiveFunction.GetDimension(); i++ {
			newPos[i] = agentInArchive.Position[i] + r2*math.Round(directVector[agentIdx][i])*agentInArchive.Position[i]
		}
	}

	a.outOfBoundaries(newPos)

	newAgent := a.Agents[agentIdx].CopyAgent()
	value, constraints, penalty := a.ObjectiveFunction.Eval(newPos)
	newAgent.Position = newPos
	newAgent.Value = value
	newAgent.Constraints = constraints
	newAgent.Penalty = penalty

	// Sanity check the index of current agent
	frontIdx := 0
	for i := range paretoFront {
		if slices.Contains(paretoFront[i], agentIdx) {
			frontIdx = i
			break
		}
	}

	// compare the new value to other agents in the same front
	dominatedFlag := 0
	for _, v := range paretoFront[frontIdx] {
		if newAgent.Dominates(a.Agents[v]) {
			dominatedFlag = 1
			break
		} else if a.Agents[v].Dominates(newAgent) {
			dominatedFlag = -1
			break
		}
	}

	newR := rand.Float64()
	if dominatedFlag == 1 || (dominatedFlag == 0 && newR > 0.5) {
		tPop = append(tPop, a.Agents[agentIdx].CopyAgent())

		a.Agents[agentIdx] = newAgent.CopyAgent()

		for i := range visitTable[agentIdx] {
			visitTable[agentIdx][i] += 1
		}

		maxVals := maxRowMatrix(visitTable)
		for i := range visitTable[agentIdx] {
			if i == agentIdx {
				continue
			}
			visitTable[i][agentIdx] = maxVals[i] + 1
		}

	} else {
		for i := range visitTable[agentIdx] {
			visitTable[agentIdx][i] += 1
		}

		tPop = append(tPop, newAgent.CopyAgent())
	}
}

func (a *MOAHAAlgorithm) initialization() {

	vals := make([]float64, a.ObjectiveFunction.NumberOfObjectives())
	for i := 0; i < a.ObjectiveFunction.NumberOfObjectives(); i++ {
		if a.ObjectiveFunction.FindMin() {
			vals[i] = math.MaxFloat64
		} else {
			vals[i] = math.MinInt64
		}
	}

	a.BestResult = &objectives.Result{
		Value: vals,
	}

	var wg sync.WaitGroup
	wg.Add(a.NumberOfAgents)
	for agentIdx := 0; agentIdx < a.NumberOfAgents; agentIdx++ {
		go func(agentIdx int) {
			defer wg.Done()
			positions := make([]float64, a.ObjectiveFunction.GetDimension())

			for i := 0; i < a.ObjectiveFunction.GetDimension(); i++ {
				positions[i] = a.ObjectiveFunction.GetLowerBound()[i] + rand.Float64()*
					(a.ObjectiveFunction.GetUpperBound()[i]-a.ObjectiveFunction.GetLowerBound()[i])
			}

			// evaluate
			newAgent := &objectives.Result{
				Idx:      agentIdx,
				Position: positions,
			}

			value, constraints, penalty := a.ObjectiveFunction.Eval(positions)
			newAgent.Value = value
			newAgent.Constraints = constraints
			newAgent.Penalty = penalty

			a.Agents[agentIdx] = newAgent
		}(agentIdx)
	}
	wg.Wait()
}

func (a *MOAHAAlgorithm) outOfBoundaries(pos []float64) {
	for i := range pos {
		if pos[i] < a.ObjectiveFunction.GetLowerBound()[i] {
			pos[i] = a.ObjectiveFunction.GetLowerBound()[i]
		} else if pos[i] > a.ObjectiveFunction.GetUpperBound()[i] {
			pos[i] = a.ObjectiveFunction.GetUpperBound()[i]
		}
	}
}

func (a *MOAHAAlgorithm) GetResults() []algorithms.AlgorithmResult {
	results := make([]algorithms.AlgorithmResult, len(a.Archive))

	for i := range a.Archive {
		res := a.Archive[i]
		mapLoc, err := a.ObjectiveFunction.GetLocationResult(res.Position)
		if err != nil {
			return nil
		}

		results[i] = algorithms.AlgorithmResult{
			MapLocations: mapLoc,
			Value:        res.Value,
			Penalty:      res.Penalty,
		}
	}

	return results
}

func initializeNMMatrix(n, m int) [][]float64 {
	matrix := make([][]float64, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]float64, m)
	}
	return matrix
}

func maxRowMatrix(matrix [][]float64) []float64 {
	res := make([]float64, len(matrix[0]))

	for i := 0; i < len(matrix[0]); i++ {
		maxVal := -math.MaxFloat64
		for j := 0; j < len(matrix); j++ {
			if i == j {
				continue
			}
			if matrix[i][j] > maxVal {
				maxVal = matrix[i][j]
			}
		}
		res[i] = maxVal
	}
	return res
}
