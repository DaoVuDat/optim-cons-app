package single

import (
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/util"
)

type SingleResult struct {
	Idx         int
	Position    []float64
	Value       []float64
	Constraints map[string]float64
	Penalty     map[string]float64
}

type SingleProblem interface {
	Eval(pos []float64) (values []float64, constraints []float64, penalty []float64)
	GetUpperBound() []float64
	GetLowerBound() []float64
	GetDimension() int
	FindMin() bool
	NumberOfObjectives() int
	Type() data.TypeProblem
}

func (agent *SingleResult) CopyAgent() *SingleResult {
	return &SingleResult{
		Idx:         agent.Idx,
		Position:    util.CopyArray(agent.Position),
		Value:       util.CopyArray(agent.Value),
		Constraints: util.CopyMap(agent.Constraints),
		Penalty:     util.CopyMap(agent.Penalty),
	}
}
