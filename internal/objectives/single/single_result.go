package single

import (
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/util"
)

type SingleResult struct {
	Idx         int
	Position    []float64
	Solution    []float64
	Value       []float64
	Constraints []float64
	Penalty     []float64
}

type SingleProblem interface {
	Eval(pos []float64, x *SingleResult) *SingleResult
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
		Solution:    util.CopyArray(agent.Solution),
		Value:       util.CopyArray(agent.Value),
		Constraints: util.CopyArray(agent.Constraints),
		Penalty:     util.CopyArray(agent.Penalty),
	}
}
