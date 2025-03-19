package util

import "golang-moaha-construction/internal/objectives"

func CopyArray(src []float64) []float64 {
	dst := make([]float64, len(src))
	copy(dst, src)
	return dst
}

func CopyAgent(agent *objectives.Result) *objectives.Result {
	return &objectives.Result{
		Position:    CopyArray(agent.Position),
		Solution:    CopyArray(agent.Solution),
		Value:       CopyArray(agent.Value),
		Constraints: CopyArray(agent.Constraints),
		Penalty:     CopyArray(agent.Penalty),
	}
}
