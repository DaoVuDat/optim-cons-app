package single

import (
	"fmt"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/util"
	"strings"
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

func (agent *SingleResult) PositionString() string {
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
