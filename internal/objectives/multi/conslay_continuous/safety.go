package conslay_continuous

import "golang-moaha-construction/internal/objectives"

const SafetyObjectiveType objectives.ObjectiveType = "Safety Objective"

const (
	SafetyProximity = "The presumed value of the safety proximity relationship "
)

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

func (obj *SafetyObjective) Eval(locations map[string]Location) float64 {
	return 0
}

func (obj *SafetyObjective) GetAlphaPenalty() float64 {
	return obj.AlphaSafetyPenalty
}
