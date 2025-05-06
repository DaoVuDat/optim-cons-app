package nsgaii

import (
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
	"math"
	"testing"
)

type MockProblem struct {
	NumberOfObjectivesFunc func() int
	GetLowerBoundFunc      func() []float64
	GetUpperBoundFunc      func() []float64
	GetDimensionFunc       func() int
	EvalFunc               func(pos []float64) (values []float64, valuesWithKey map[data.ObjectiveType]float64, key []data.ObjectiveType, penalty map[data.ConstraintType]float64)
}

func (m *MockProblem) GetUpperBound() []float64 {
	if m.GetUpperBoundFunc != nil {
		return m.GetUpperBoundFunc()
	}
	return nil
}

func (m *MockProblem) GetLowerBound() []float64 {
	if m.GetLowerBoundFunc != nil {
		return m.GetLowerBoundFunc()
	}
	return nil
}

func (m *MockProblem) GetDimension() int {
	if m.GetDimensionFunc != nil {
		return m.GetDimensionFunc()
	}
	return 0
}

func (m *MockProblem) FindMin() bool {
	return true // Default implementation
}

func (m *MockProblem) NumberOfObjectives() int {
	if m.NumberOfObjectivesFunc != nil {
		return m.NumberOfObjectivesFunc()
	}
	return 0
}

func (m *MockProblem) Type() data.TypeProblem {
	return data.Multi
}

func (m *MockProblem) InitializeObjectives() error {
	return nil // Default implementation
}

func (m *MockProblem) InitializeConstraints() error {
	return nil // Default implementation
}

func (m *MockProblem) SetCranesLocations(locations []data.Crane) error {
	return nil
}

func (m *MockProblem) GetCranesLocations() []data.Crane {
	return nil
}

func (m *MockProblem) GetLocations() map[string]data.Location {
	return nil
}

func (m *MockProblem) GetObjectives() map[data.ObjectiveType]data.Objectiver {
	return nil
}

func (m *MockProblem) GetConstraints() map[data.ConstraintType]data.Constrainter {
	return nil
}

func (m *MockProblem) AddObjective(name data.ObjectiveType, objective data.Objectiver) error {
	return nil
}

func (m *MockProblem) AddConstraint(name data.ConstraintType, constraint data.Constrainter) error {
	return nil
}

func (m *MockProblem) GetPhases() [][]string {
	return nil
}

func (m *MockProblem) GetLocationResult(input []float64) (map[string]data.Location, []data.Location, []data.Crane, error) {
	return nil, nil, nil, nil
}

func (m *MockProblem) GetLayoutSize() (minX float64, maxX float64, minY float64, maxY float64, err error) {
	return 0, 0, 0, 0, nil
}

func (m *MockProblem) Eval(input []float64) (
	values []float64,
	valuesWithKey map[data.ObjectiveType]float64,
	key []data.ObjectiveType,
	penalty map[data.ConstraintType]float64) {
	return m.EvalFunc(input)
}

func TestNSGAIIAlgorithm_calculateCrowdingDistance(t *testing.T) {
	// Define the test case
	pop := []*objectives.Result{
		{Value: []float64{17482666.7708616, 1744703659.10514}, Rank: 4},
		{Value: []float64{9052134.75096372, 901508700.908080}, Rank: 1},
		{Value: []float64{13124158.3347648, 1309485059.52912}, Rank: 3},
		{Value: []float64{5032354.37946573, 499739906.074043}, Rank: 0},
		{Value: []float64{9987372.54571255, 995070325.219069}, Rank: 2},
	}

	paretoFront := [][]int{{3}, {1}, {4}, {2}, {0}}

	// Create a mock ObjectiveFunction
	var mockObjectiveFunction objectives.Problem = &MockProblem{
		NumberOfObjectivesFunc: func() int {
			return 2
		},
	}

	// Create NSGAIIAlgorithm instance with the mock ObjectiveFunction
	ga := &NSGAIIAlgorithm{
		ObjectiveFunction: mockObjectiveFunction,
	}

	// Call the function
	updatedPop := ga.calculateCrowdingDistance(pop, paretoFront)

	// All points should have infinite crowding distance since they're all in their own front
	for i, p := range updatedPop {
		if !math.IsInf(p.CrowdingDistance, 1) {
			t.Errorf("Agent %d: Expected crowding distance to be +Inf, got %f", i, p.CrowdingDistance)
		}
	}
}

func TestSortPopulation_firstCase(t *testing.T) {
	// Create a population of 5 agents with specified ranks and inf crowding distances.
	pop := []*objectives.Result{
		{Idx: 0, Rank: 4, CrowdingDistance: math.Inf(1)},
		{Idx: 1, Rank: 2, CrowdingDistance: math.Inf(1)},
		{Idx: 2, Rank: 1, CrowdingDistance: math.Inf(1)},
		{Idx: 3, Rank: 0, CrowdingDistance: math.Inf(1)},
		{Idx: 4, Rank: 3, CrowdingDistance: math.Inf(1)},
	}

	// Expected sorted order by ascending rank: indices [3, 2, 1, 4, 0]
	expectedOrder := []int{0, 1, 2, 3, 4} // expected rank values in sorted population: 0,1,2,3,4

	sortedPop, _ := SortPopulation(pop)

	// Validate order (by checking rank)
	for i, agent := range sortedPop {
		if agent.Rank != expectedOrder[i] {
			t.Errorf("Sorted agent at position %d: Got rank %d, expected %d", i, agent.Rank, expectedOrder[i])
		}
		if !math.IsInf(agent.CrowdingDistance, 1) {
			t.Errorf("Sorted agent at position %d: Expected crowding distance +Inf, got %f", i, agent.CrowdingDistance)
		}
	}
}
