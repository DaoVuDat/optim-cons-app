package moaha

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
	"golang-moaha-construction/internal/objectives/multi"
	"golang-moaha-construction/internal/objectives/single"
	"golang-moaha-construction/internal/util"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

const (
	numberOfObjective = 1
	Name              = "MOAHA"
	NumAgents         = "Number of Agents"
	NumIters          = "Number of Iteration"
	ArchiveSize       = "Archive Size"
)

type MOAHAAlgorithm struct {
	NumberOfAgents    int
	NumberOfIter      int
	Agents            []*multi.MultiResult
	BestResult        *multi.MultiResult
	Convergence       []float64
	ObjectiveFunction multi.MultiProblem
	ArchiveSize       int
	Archive           []*multi.MultiResult
}

var Configs = []data.Config{
	{
		Name:               NumAgents,
		ValidationFunction: util.IsValidPositiveInteger,
	},
	{
		Name:               NumIters,
		ValidationFunction: util.IsValidPositiveInteger,
	},
}

func Create(
	problem multi.MultiProblem,
	configs []*data.Config,
) (*MOAHAAlgorithm, error) {

	var numsIters, numAgents, archiveSize int

	for _, config := range configs {
		// sanity check
		val := strings.Trim(config.Value, " ")
		switch config.Name {
		case NumAgents:
			numAgents, _ = strconv.Atoi(val)
		case NumIters:
			numsIters, _ = strconv.Atoi(val)
		case ArchiveSize:
			archiveSize, _ = strconv.Atoi(val)

		}
	}

	convergence := make([]float64, numsIters)
	agents := make([]*multi.MultiResult, numAgents)

	archive := make([]*multi.MultiResult, 0, archiveSize)

	if numberOfObjective != problem.NumberOfObjectives() {
		return nil, objectives.ErrInvalidNumberOfObjectives
	}

	return &MOAHAAlgorithm{
		NumberOfAgents:    numAgents,
		NumberOfIter:      numsIters,
		Convergence:       convergence,
		ObjectiveFunction: problem,
		Agents:            agents,
		Archive:           archive,
	}, nil
}

func (a *MOAHAAlgorithm) Type() data.TypeProblem {
	return data.Multi
}

func (a *MOAHAAlgorithm) Run() error {
	dimensions := a.ObjectiveFunction.GetDimension()

	// initialization
	a.initialization()

	a.Agents = multi.DetermineDomination(a.Agents)
	a.Archive = multi.GetNonDominatedAgents(a.Agents)

	l := 0

	bar := progressbar.Default(int64(a.NumberOfIter))
	//var wg sync.WaitGroup

	visitTable := initializeNMMatrix(a.NumberOfAgents, a.NumberOfAgents)

	for l < a.NumberOfIter {
		// direct vector
		directVector := initializeNMMatrix(a.NumberOfAgents, dimensions)

		//wg.Add(a.NumberOfAgents)
		for agentIdx := range a.Agents {

			r := rand.Float64()

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
				for i := 0; i < randNum; i++ {
					directVector[agentIdx][i] = 1
				}
			}

			r = rand.Float64()

			if r < 0.5 {
				// guided foraging
				a.guidedForaging(visitTable, directVector, agentIdx)
			} else {
				// territory foraging
				a.territoryForaging(visitTable, directVector, agentIdx)
			}

		}

		//wg.Wait()

		// migration foraging
		if l%(a.NumberOfAgents*2) == 0 {
			maxVal := -math.MaxFloat64
			maxIdx := 0
			for i := range a.Agents {
				if a.Agents[i].Value[0] > maxVal {
					maxVal = a.Agents[i].Value[0]
					maxIdx = i
				}
			}

			for i := range a.Agents[maxIdx].Position {
				a.Agents[maxIdx].Position[i] =
					a.ObjectiveFunction.GetLowerBound()[i] + rand.Float64()*
						(a.ObjectiveFunction.GetUpperBound()[i]-a.ObjectiveFunction.GetLowerBound()[i])
			}

			// evaluate
			a.Agents[maxIdx] = a.ObjectiveFunction.Eval(a.Agents[maxIdx].Position, a.Agents[maxIdx])

			for i := range visitTable[maxIdx] {
				visitTable[maxIdx][i] += 1
			}

			maxVals := maxRowMatrix(visitTable)
			for i := range visitTable[maxIdx] {
				if i == maxIdx {
					continue
				}
				visitTable[i][maxIdx] = maxVals[i] + 1
			}
		}

		a.findBest()

		a.Convergence[l] = a.BestResult.Value[0]
		bar.Describe(fmt.Sprintf("Iter %d: %e", l+1, a.BestResult.Value[0]))
		bar.Add(1)

		l++
	}

	return nil
}

func (a *MOAHAAlgorithm) guidedForaging(visitTable [][]float64, directVector [][]float64, agentIdx int) {
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
		candidateIdx := maxValIdxs[0]
		minFitness := a.Agents[candidateIdx].Value[0]

		for i := 1; i < len(maxValIdxs); i++ {
			if a.Agents[maxValIdxs[i]].Value[0] < minFitness {
				minFitness = a.Agents[maxValIdxs[i]].Value[0]
				candidateIdx = maxValIdxs[i]
			}
		}
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

	newAgent := a.ObjectiveFunction.Eval(newPos, a.Agents[agentIdx])

	if newAgent.Value[0] < a.Agents[agentIdx].Value[0] {
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
	}
}

func (a *MOAHAAlgorithm) territoryForaging(visitTable [][]float64, directVector [][]float64, agentIdx int) {
	r := rand.NormFloat64()
	newPos := make([]float64, a.ObjectiveFunction.GetDimension())
	for i := 0; i < a.ObjectiveFunction.GetDimension(); i++ {
		newPos[i] = a.Agents[agentIdx].Position[i] + r*math.Round(directVector[agentIdx][i])*a.Agents[agentIdx].Position[i]
	}

	a.outOfBoundaries(newPos)

	newAgent := a.ObjectiveFunction.Eval(newPos, a.Agents[agentIdx])

	if newAgent.Value[0] < a.Agents[agentIdx].Value[0] {
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

	a.BestResult = &multi.MultiResult{
		SingleResult: single.SingleResult{
			Value: vals,
		},
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
			newAgent := &multi.MultiResult{
				SingleResult: single.SingleResult{
					Idx:      agentIdx,
					Position: positions,
					Solution: positions,
				},
			}
			a.Agents[agentIdx] = a.ObjectiveFunction.Eval(positions, newAgent)
		}(agentIdx)
	}
	wg.Wait()

	//a.findBest()
}

// This algorithm solve only one objective
func (a *MOAHAAlgorithm) findBest() {
	for i := range a.Agents {
		if a.Agents[i].Value[0] < a.BestResult.Value[0] {
			a.BestResult = a.Agents[i].CopyAgent()
		}
	}
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
