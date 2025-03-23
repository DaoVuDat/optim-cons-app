package multi

import (
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives/single"
	"golang-moaha-construction/internal/util"
	"math"
	"sort"
)

type MultiResult struct {
	single.SingleResult
	CrowdingDistance float64
	Dominated        bool
	Rank             int
	DominationSet    []int
	DominatedCount   int
}

type MultiProblem interface {
	Eval(pos []float64, x *MultiResult) *MultiResult
	GetUpperBound() []float64
	GetLowerBound() []float64
	GetDimension() int
	FindMin() bool
	NumberOfObjectives() int
	Type() data.TypeProblem
}

func (agent *MultiResult) CopyAgent() *MultiResult {
	return &MultiResult{
		SingleResult: single.SingleResult{
			Idx:         agent.Idx,
			Position:    util.CopyArray(agent.Position),
			Solution:    util.CopyArray(agent.Solution),
			Value:       util.CopyArray(agent.Value),
			Constraints: util.CopyArray(agent.Constraints),
			Penalty:     util.CopyArray(agent.Penalty),
		},
		CrowdingDistance: agent.CrowdingDistance,
		Dominated:        agent.Dominated,
		Rank:             agent.Rank,
		DominationSet:    util.CopyArray(agent.DominationSet),
		DominatedCount:   agent.DominatedCount,
	}
}

func (agent *MultiResult) Dominates(other *MultiResult) bool {
	numberOfObjs := len(agent.Value)
	anyConstraint := false
	for i := 0; i < numberOfObjs; i++ {
		if agent.Value[i] > other.Value[i] {
			return false
		}

		if agent.Value[i] < other.Value[i] {
			anyConstraint = true
		}
	}

	return anyConstraint
}

func DetermineDomination(agents []*MultiResult) []*MultiResult {
	// clear the dominated
	for i := range agents {
		agents[i].Dominated = false
	}

	// determine domination
	for i := 0; i < len(agents)-1; i++ {
		for j := i + 1; j < len(agents); j++ {
			if agents[i].Dominates(agents[j]) {
				agents[j].Dominated = true
			} else if agents[j].Dominates(agents[i]) {
				agents[i].Dominated = true
				break
			} else {
				// check all values are equal
				allEqual := true
				for k := 0; k < len(agents[i].Value); k++ {
					if agents[i].Value[k] != agents[j].Value[k] {
						allEqual = false
						break
					}
				}

				if allEqual {
					agents[i].Dominated = true
				}
			}

		}
	}

	return agents
}

func GetNonDominatedAgents(agents []*MultiResult) []*MultiResult {
	res := make([]*MultiResult, 0)
	for _, agent := range agents {
		if !agent.Dominated {
			res = append(res, agent.CopyAgent())
		}
	}

	return res
}

func NonDominatedSort(agents []*MultiResult) ([]*MultiResult, [][]int) {
	// clear domination set and domination count
	for i := range agents {
		agents[i].DominationSet = make([]int, 0)
		agents[i].DominatedCount = 0
	}

	paretoFront := [][]int{
		make([]int, 0),
	}

	for i := 0; i < len(agents)-1; i++ {
		for j := i + 1; j < len(agents); j++ {
			p := agents[i]
			q := agents[j]

			if p.Dominates(q) {
				p.DominationSet = append(p.DominationSet, q.Idx)
				q.DominatedCount += 1
			}

			if q.Dominates(p) {
				q.DominationSet = append(q.DominationSet, p.Idx)
				p.DominatedCount += 1
			}
		}

		if agents[i].DominatedCount == 0 {
			paretoFront[0] = append(paretoFront[0], agents[i].Idx)
			agents[i].Rank = 0
		}
	}

	k := 0

	for {
		Q := make([]int, 0)

		for _, v := range paretoFront[k] {
			p := agents[v]

			for _, j := range p.DominationSet {
				q := agents[j]
				q.DominatedCount -= 1

				if q.DominatedCount == 0 {
					Q = append(Q, q.Idx)
					q.Rank = k + 1
				}
			}
		}

		if len(Q) == 0 {
			break
		}

		paretoFront = append(paretoFront, Q)

		k++
	}

	return agents, paretoFront
}

type SortedDEDC struct {
	values      []float64
	originalIdx int
	sortedIdx   []int
}

type SortedValue struct {
	Value float64
	Idx   int
}

// DECD Dynamic Elimination-based Crowding Distance
func DECD(agents []*MultiResult, excess int) []*MultiResult {
	numberOfAgents := len(agents)
	numberOfObjs := len(agents[0].Value)

	// tracks the sorted values for each agent
	trackingList := make([]SortedDEDC, numberOfAgents)
	for i := 0; i < numberOfAgents; i++ {
		trackingList[i] = SortedDEDC{
			values:      agents[i].Value,
			originalIdx: i,
			sortedIdx:   make([]int, len(agents[i].Value)),
		}
	}

	sortedValues := make([][]SortedValue, len(agents[0].Value))

	for i := 0; i < numberOfObjs; i++ {
		values := make([]SortedValue, numberOfAgents)
		for j := 0; j < numberOfAgents; j++ {
			values[j] = SortedValue{
				Value: agents[j].Value[i],
				Idx:   j,
			}
		}

		sort.Slice(values, func(i, j int) bool {
			return values[i].Value < values[j].Value
		})

		for j := 0; j < len(values); j++ {
			trackingList[values[j].Idx].sortedIdx[i] = j
		}

		sortedValues[i] = values
	}

	// calculates distance matrix before removing
	distanceMatrix := util.InitializeNMMatrix(numberOfAgents, numberOfObjs)

	for i := range numberOfObjs {

		// eliminates the first and last (or min and max)
		for j := 1; j < len(distanceMatrix)-1; j++ {
			prev := sortedValues[i][j-1].Value
			next := sortedValues[i][j+1].Value
			distanceMatrix[sortedValues[i][j].Idx][i] = math.Abs(next-prev) / math.Abs(sortedValues[i][0].Value-sortedValues[i][len(sortedValues[i])-1].Value)
		}

		// add max distance to the first and the last
		distanceMatrix[sortedValues[i][0].Idx][i] = math.MaxFloat64
		distanceMatrix[sortedValues[i][len(sortedValues[i])-1].Idx][i] = math.MaxFloat64
	}

	distances := make([]float64, numberOfAgents)
	for i := range len(distances) {
		for j := range numberOfObjs {
			if distanceMatrix[i][j] == math.MaxFloat64 {
				distances[i] = math.MaxFloat64
				break
			} else {
				distances[i] += distanceMatrix[i][j]
			}
		}
	}

	// remove exceeded agents in archive
	for excess > 0 {
		_, minIdx := util.MinWithIdx(distances)

		// get agent idx
		deletedIdx := util.CopyArray(trackingList[minIdx].sortedIdx)

		// remove agent from archive
		agents = util.Remove(agents, minIdx)

		//  sortedValues
		for i := range numberOfObjs {
			// remove sortedValues
			sortedValues[i] = util.Remove(sortedValues[i], deletedIdx[i])

		}

		// remove distances matrix and distances
		distanceMatrix = util.Remove(distanceMatrix, minIdx)
		distances = util.Remove(distances, minIdx)

		// re-calculate crowding distance after removing agent from archive
		for i := 0; i < numberOfObjs; i++ {
			//values := sortedValues[i]

			values := make([]SortedValue, len(sortedValues[i]))
			for j := 0; j < len(sortedValues[i]); j++ {
				values[j] = SortedValue{
					Value: agents[j].Value[i],
					Idx:   j,
				}
			}

			sort.Slice(values, func(i, j int) bool {
				return values[i].Value < values[j].Value
			})

			for j := 0; j < len(values); j++ {
				trackingList[values[j].Idx].sortedIdx[i] = j
			}

			sortedValues[i] = values
		}

		for i := range numberOfObjs {
			// we only re-calculate the prev and next after the removed agent
			// prev
			if deletedIdx[i] > 1 {
				currentPrev := sortedValues[i][deletedIdx[i]-2].Value
				currentNext := sortedValues[i][deletedIdx[i]].Value
				distanceMatrix[sortedValues[i][deletedIdx[i]-1].Idx][i] = math.Abs(currentNext-currentPrev) / math.Abs(sortedValues[i][0].Value-sortedValues[i][len(sortedValues[i])-1].Value)
			} else {
				distanceMatrix[sortedValues[i][deletedIdx[i]-1].Idx][i] = math.MaxFloat64
			}

			// next
			if deletedIdx[i] < len(sortedValues[i])-1 {
				currentPrev := sortedValues[i][deletedIdx[i]-1].Value
				currentNext := sortedValues[i][deletedIdx[i]+1].Value
				distanceMatrix[sortedValues[i][deletedIdx[i]].Idx][i] = math.Abs(currentNext-currentPrev) / math.Abs(sortedValues[i][0].Value-sortedValues[i][len(sortedValues[i])-1].Value)
			} else {
				distanceMatrix[sortedValues[i][deletedIdx[i]].Idx][i] = math.MaxFloat64
			}

		}

		distances = make([]float64, len(distances))
		for i := range len(distances) {
			for j := range numberOfObjs {
				if distanceMatrix[i][j] == math.MaxFloat64 {
					distances[i] = math.MaxFloat64
					break
				} else {
					distances[i] += distanceMatrix[i][j]
				}
			}
		}

		excess--
	}

	return agents
}
