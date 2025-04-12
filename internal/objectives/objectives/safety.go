package objectives

import (
	"golang-moaha-construction/internal/data"
)

const SafetyObjectiveType data.ObjectiveType = "Safety Objective"

type SafetyConfigs struct {
	SafetyProximity    [][]float64
	AlphaSafetyPenalty float64
}

type SafetyObjective struct {
	SafetyProximity    [][]float64
	AlphaSafetyPenalty float64
}

func CreateSafetyObjective() (*SafetyObjective, error) {
	return &SafetyObjective{}, nil
}

func CreateSafetyObjectiveFromConfig(safetyConfigs SafetyConfigs) (*SafetyObjective, error) {
	safetyObj := &SafetyObjective{
		SafetyProximity:    safetyConfigs.SafetyProximity,
		AlphaSafetyPenalty: safetyConfigs.AlphaSafetyPenalty,
	}
	return safetyObj, nil
}

func (obj *SafetyObjective) Eval(locations map[string]data.Location) float64 {
	return 0
}

func (obj *SafetyObjective) GetAlphaPenalty() float64 {
	return obj.AlphaSafetyPenalty
}
