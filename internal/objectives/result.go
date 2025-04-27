package objectives

import (
	"fmt"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/util"
	"math"
	"sort"
	"strings"
)

type Result struct {
	Idx              int
	Position         []float64
	Value            []float64
	ValuesWithKey    map[data.ObjectiveType]float64
	Penalty          map[data.ConstraintType]float64
	Key              []data.ObjectiveType
	CrowdingDistance float64
	Dominated        bool
	Rank             int
	DominationSet    []int
	DominatedCount   int
}

func (agent *Result) PositionString() string {
	var sb strings.Builder
	sb.WriteString("[ ")

	for i, v := range agent.Position {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%g", v))
	}

	sb.WriteString(" ]")
	return sb.String()
}

func (agent *Result) CopyAgent() *Result {
	return &Result{
		Idx:              agent.Idx,
		Position:         util.CopyArray(agent.Position),
		Value:            util.CopyArray(agent.Value),
		ValuesWithKey:    util.CopyMap(agent.ValuesWithKey),
		Penalty:          util.CopyMap(agent.Penalty),
		Key:              util.CopyArray(agent.Key),
		CrowdingDistance: agent.CrowdingDistance,
		Dominated:        agent.Dominated,
		Rank:             agent.Rank,
		DominationSet:    util.CopyArray(agent.DominationSet),
		DominatedCount:   agent.DominatedCount,
	}
}

func (agent *Result) Dominates(other *Result) bool {
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

func MergeAgents(a []*Result, b []*Result) []*Result {
	res := make([]*Result, len(a)+len(b))
	for i := 0; i < len(a); i++ {
		res[i] = a[i]
		res[i].Idx = i
	}

	for i := 0; i < len(b); i++ {
		res[i+len(a)] = b[i]
		res[i+len(a)].Idx = i + len(a)
	}

	return res
}
func DetermineDomination(agents []*Result) []*Result {
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

func GetNonDominatedAgents(agents []*Result) []*Result {
	res := make([]*Result, 0)
	for _, agent := range agents {
		if !agent.Dominated {
			res = append(res, agent.CopyAgent())
		}
	}

	return res
}

func NonDominatedSort(agents []*Result) ([]*Result, [][]int) {
	// clear domination set and domination count
	for i := range agents {
		agents[i].DominationSet = make([]int, 0)
		agents[i].DominatedCount = 0
	}

	paretoFront := [][]int{
		make([]int, 0),
	}

	for i := 0; i < len(agents); i++ {
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

func SplitToNPop(agents []*Result, nPop int, paretoFront [][]int) []*Result {

	results := make([]*Result, nPop)

	count := 0
	for _, v := range paretoFront {
		for _, idx := range v {
			if count >= nPop {
				break
			}
			results[count] = agents[idx].CopyAgent()
			results[count].Idx = count
			count++
		}
	}

	return results
}

//// DECD Dynamic Elimination-based Crowding Distance
//func DECD(agents []*Result, excess int) []*Result {
//	fmt.Println("Start DECD")
//
//	numberOfAgents := len(agents)
//	numberOfObjs := len(agents[0].Value)
//
//	// tracks the sorted values for each agent
//	trackingList := make([]SortedDEDC, numberOfAgents)
//	for i := 0; i < numberOfAgents; i++ {
//		trackingList[i] = SortedDEDC{
//			values:      agents[i].Value,
//			originalIdx: i,
//			sortedIdx:   make([]int, len(agents[i].Value)),
//		}
//	}
//
//	sortedValues := make([][]SortedValue, len(agents[0].Value))
//
//	for i := 0; i < numberOfObjs; i++ {
//		values := make([]SortedValue, numberOfAgents)
//		for j := 0; j < numberOfAgents; j++ {
//			values[j] = SortedValue{
//				Value: agents[j].Value[i],
//				Idx:   j,
//			}
//		}
//
//		sort.Slice(values, func(i, j int) bool {
//			return values[i].Value < values[j].Value
//		})
//
//		for j := 0; j < len(values); j++ {
//			trackingList[values[j].Idx].sortedIdx[i] = j
//		}
//
//		sortedValues[i] = values
//	}
//
//	// calculates distance matrix before removing
//	distanceMatrix := util.InitializeNMMatrix(numberOfAgents, numberOfObjs)
//
//	for i := range numberOfObjs {
//
//		// eliminates the first and last (or min and max)
//		for j := 1; j < len(distanceMatrix)-1; j++ {
//			prev := sortedValues[i][j-1].Value
//			next := sortedValues[i][j+1].Value
//			distanceMatrix[sortedValues[i][j].Idx][i] = math.Abs(next-prev) / math.Abs(sortedValues[i][0].Value-sortedValues[i][len(sortedValues[i])-1].Value)
//		}
//
//		// add max distance to the first and the last
//		distanceMatrix[sortedValues[i][0].Idx][i] = math.MaxFloat64
//		distanceMatrix[sortedValues[i][len(sortedValues[i])-1].Idx][i] = math.MaxFloat64
//	}
//
//	distances := make([]float64, numberOfAgents)
//	for i := range len(distances) {
//		for j := range numberOfObjs {
//			if distanceMatrix[i][j] == math.MaxFloat64 {
//				distances[i] = math.MaxFloat64
//				break
//			} else {
//				distances[i] += distanceMatrix[i][j]
//			}
//		}
//	}
//
//	fmt.Println("start excess")
//	// remove exceeded agents in archive
//	for excess > 0 {
//		_, minIdx := util.MinWithIdx(distances)
//
//		// get agent idx
//		deletedIdx := util.CopyArray(trackingList[minIdx].sortedIdx)
//
//		// remove agent from archive
//		agents = util.Remove(agents, minIdx)
//		fmt.Println("start startedValues")
//
//		//  sortedValues
//		for i := range numberOfObjs {
//			// remove sortedValues
//			sortedValues[i] = util.Remove(sortedValues[i], deletedIdx[i])
//
//		}
//
//		// remove distances matrix and distances
//		distanceMatrix = util.Remove(distanceMatrix, minIdx)
//		distances = util.Remove(distances, minIdx)
//
//		fmt.Println("start loop")
//		// re-calculate crowding distance after removing agent from archive
//		for i := 0; i < numberOfObjs; i++ {
//			//values := sortedValues[i]
//
//			values := make([]SortedValue, len(sortedValues[i]))
//			for j := 0; j < len(sortedValues[i]); j++ {
//				values[j] = SortedValue{
//					Value: agents[j].Value[i],
//					Idx:   j,
//				}
//			}
//
//			sort.Slice(values, func(i, j int) bool {
//				return values[i].Value < values[j].Value
//			})
//
//			for j := 0; j < len(values); j++ {
//				trackingList[values[j].Idx].sortedIdx[i] = j
//			}
//
//			sortedValues[i] = values
//		}
//
//		fmt.Println("start numberOfObjs")
//		for i := range numberOfObjs {
//			// we only re-calculate the prev and next after the removed agent
//			// prev
//			if deletedIdx[i] > 1 {
//				currentPrev := sortedValues[i][deletedIdx[i]-2].Value
//				currentNext := sortedValues[i][deletedIdx[i]].Value
//				distanceMatrix[sortedValues[i][deletedIdx[i]-1].Idx][i] = math.Abs(currentNext-currentPrev) / math.Abs(sortedValues[i][0].Value-sortedValues[i][len(sortedValues[i])-1].Value)
//			} else {
//				distanceMatrix[sortedValues[i][deletedIdx[i]-1].Idx][i] = math.MaxFloat64
//			}
//
//			// next
//			if deletedIdx[i] < len(sortedValues[i])-1 {
//				currentPrev := sortedValues[i][deletedIdx[i]-1].Value
//				currentNext := sortedValues[i][deletedIdx[i]+1].Value
//				distanceMatrix[sortedValues[i][deletedIdx[i]].Idx][i] = math.Abs(currentNext-currentPrev) / math.Abs(sortedValues[i][0].Value-sortedValues[i][len(sortedValues[i])-1].Value)
//			} else {
//				distanceMatrix[sortedValues[i][deletedIdx[i]].Idx][i] = math.MaxFloat64
//			}
//
//		}
//
//		fmt.Println("start distances")
//		distances = make([]float64, len(distances))
//		for i := range len(distances) {
//			for j := range numberOfObjs {
//				if distanceMatrix[i][j] == math.MaxFloat64 {
//					distances[i] = math.MaxFloat64
//					break
//				} else {
//					distances[i] += distanceMatrix[i][j]
//				}
//			}
//		}
//
//		excess--
//	}
//
//	fmt.Println("End DECD")
//	return agents
//}

func DECD(agents []*Result, excess int) []*Result {

	numberOfAgents := len(agents)
	if numberOfAgents <= excess {
		return agents[:0] // Return empty slice if we're asked to remove all or more
	}

	numberOfObjs := len(agents[0].Value)

	// Create cost matrix similar to MATLAB implementation
	costs := make([][]float64, numberOfAgents)
	for i := 0; i < numberOfAgents; i++ {
		costs[i] = make([]float64, numberOfObjs)
		for j := 0; j < numberOfObjs; j++ {
			costs[i][j] = agents[i].Value[j]
		}
	}

	// For each objective, we'll track the sorting
	sortedIndices := make([][]int, numberOfObjs)
	reverseSortedIndices := make([][]int, numberOfObjs)

	// Initialize distance matrix
	distanceMatrix := make([][]float64, numberOfAgents)
	for i := range distanceMatrix {
		distanceMatrix[i] = make([]float64, numberOfObjs)
	}

	// Process each objective
	for j := 0; j < numberOfObjs; j++ {
		// Create values with indices for sorting
		values := make([]SortedValue, numberOfAgents)
		for i := 0; i < numberOfAgents; i++ {
			values[i] = SortedValue{
				Value: costs[i][j],
				Idx:   i,
			}
		}

		// Sort by this objective
		sort.Slice(values, func(i, k int) bool {
			return values[i].Value < values[k].Value
		})

		// Store sorted indices
		sortedIndices[j] = make([]int, numberOfAgents)
		for i := 0; i < numberOfAgents; i++ {
			sortedIndices[j][i] = values[i].Idx
		}

		// Create reverse mapping
		reverseSortedIndices[j] = make([]int, numberOfAgents)
		for i := 0; i < numberOfAgents; i++ {
			reverseSortedIndices[j][values[i].Idx] = i
		}

		// Calculate distances for this objective
		for i := 1; i < numberOfAgents-1; i++ {
			idx := sortedIndices[j][i]
			nextValue := costs[sortedIndices[j][i+1]][j]
			prevValue := costs[sortedIndices[j][i-1]][j]
			minValue := costs[sortedIndices[j][0]][j]
			maxValue := costs[sortedIndices[j][numberOfAgents-1]][j]

			// Normalize the distance
			distanceMatrix[idx][j] = math.Abs(nextValue-prevValue) / math.Abs(maxValue-minValue)
		}

		// Set boundary points to infinity
		distanceMatrix[sortedIndices[j][0]][j] = math.Inf(1)                // First point
		distanceMatrix[sortedIndices[j][numberOfAgents-1]][j] = math.Inf(1) // Last point
	}

	// Calculate total distance for each agent
	distances := make([]float64, numberOfAgents)
	for i := 0; i < numberOfAgents; i++ {
		for j := 0; j < numberOfObjs; j++ {
			// If any objective has infinity, the total is infinity
			if math.IsInf(distanceMatrix[i][j], 1) {
				distances[i] = math.Inf(1)
				break
			}
			distances[i] += distanceMatrix[i][j]
		}
	}

	// Remove agents with smallest crowding distance
	for e := 0; e < excess; e++ {

		// Find agent with minimum distance
		minDistance := math.Inf(1)
		minIdx := -1
		for i := 0; i < len(distances); i++ {
			if distances[i] < minDistance {
				minDistance = distances[i]
				minIdx = i
			}
		}

		if minIdx == -1 {
			fmt.Println("No valid agent to remove")
			break
		}

		// Store the position of the agent in each sorted objective array
		deletedPositions := make([]int, numberOfObjs)
		for j := 0; j < numberOfObjs; j++ {
			deletedPositions[j] = reverseSortedIndices[j][minIdx]
		}

		// Remove the agent
		agents = append(agents[:minIdx], agents[minIdx+1:]...)
		costs = append(costs[:minIdx], costs[minIdx+1:]...)

		// Remove from distance matrix and distances array
		distanceMatrix = append(distanceMatrix[:minIdx], distanceMatrix[minIdx+1:]...)
		distances = append(distances[:minIdx], distances[minIdx+1:]...)

		// Update all indices
		numberOfAgents--

		// Recalculate sorting and distances
		for j := 0; j < numberOfObjs; j++ {
			// Create new values array after removal
			values := make([]SortedValue, numberOfAgents)
			for i := 0; i < numberOfAgents; i++ {
				values[i] = SortedValue{
					Value: costs[i][j],
					Idx:   i,
				}
			}

			// Sort by this objective
			sort.Slice(values, func(i, k int) bool {
				return values[i].Value < values[k].Value
			})

			// Update sorted indices
			sortedIndices[j] = make([]int, numberOfAgents)
			for i := 0; i < numberOfAgents; i++ {
				sortedIndices[j][i] = values[i].Idx
			}

			// Update reverse mapping
			reverseSortedIndices[j] = make([]int, numberOfAgents)
			for i := 0; i < numberOfAgents; i++ {
				reverseSortedIndices[j][values[i].Idx] = i
			}

			// Position where agent was deleted
			pos := deletedPositions[j]

			// Update distances for affected positions
			// Update previous position (pos-1)
			if pos > 1 && pos-1 < numberOfAgents {
				idx := sortedIndices[j][pos-1]
				if pos-2 >= 0 && pos < numberOfAgents {
					// Normal case: has both previous and next
					prevValue := costs[sortedIndices[j][pos-2]][j]
					nextValue := costs[sortedIndices[j][pos]][j]
					minValue := costs[sortedIndices[j][0]][j]
					maxValue := costs[sortedIndices[j][numberOfAgents-1]][j]

					distanceMatrix[idx][j] = math.Abs(nextValue-prevValue) / math.Abs(maxValue-minValue)
				} else {
					// Edge case
					distanceMatrix[idx][j] = math.Inf(1)
				}
			}

			// Update current position (now occupied by what was at pos)
			if pos < numberOfAgents {
				idx := sortedIndices[j][pos]
				if pos > 0 && pos+1 < numberOfAgents {
					// Normal case
					prevValue := costs[sortedIndices[j][pos-1]][j]
					nextValue := costs[sortedIndices[j][pos+1]][j]
					minValue := costs[sortedIndices[j][0]][j]
					maxValue := costs[sortedIndices[j][numberOfAgents-1]][j]

					distanceMatrix[idx][j] = math.Abs(nextValue-prevValue) / math.Abs(maxValue-minValue)
				} else {
					// Edge case
					distanceMatrix[idx][j] = math.Inf(1)
				}
			}
		}

		// Recalculate total distances for affected agents
		for i := 0; i < numberOfAgents; i++ {
			distances[i] = 0
			infinityFound := false

			for j := 0; j < numberOfObjs; j++ {
				if math.IsInf(distanceMatrix[i][j], 1) {
					distances[i] = math.Inf(1)
					infinityFound = true
					break
				}
				distances[i] += distanceMatrix[i][j]
			}

			if !infinityFound {
				for j := 0; j < numberOfObjs; j++ {
					distances[i] += distanceMatrix[i][j]
				}
			}
		}
	}

	return agents
}
