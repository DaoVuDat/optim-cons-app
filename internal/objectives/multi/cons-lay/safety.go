package cons_lay

const (
	SafetyObjectiveType = "Safety Objective"
	SafetyProximity     = "The presumed value of the safety proximity relationship "
)

type SafetyConfigs struct {
	SafetyProximity [][]float64
}

type SafetyObjective struct {
	SafetyProximity [][]float64
}

func CreateSafetyObjective() (*SafetyObjective, error) {
	return &SafetyObjective{}, nil
}

func CreateSafetyObjectiveFromConfig(safetyConfigs SafetyConfigs) (*SafetyObjective, error) {
	safetyObj := &SafetyObjective{
		SafetyProximity: safetyConfigs.SafetyProximity,
	}
	return safetyObj, nil
}

func (obj *SafetyObjective) Eval() float64 {
	return 0
}
