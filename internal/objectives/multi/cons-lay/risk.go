package cons_lay

const (
	RiskObjectiveType       = "Risk Objective"
	HazardInteractionMatrix = "Hazard Interaction Matrix"
	Delta                   = "Delta"
)

type RiskConfigs struct {
	HazardInteractionMatrix [][]float64
	Delta                   float64
}

type RiskObjective struct {
	HazardInteractionMatrix [][]float64
	Delta                   float64
}

func CreateRiskObjective() (*RiskObjective, error) {
	return &RiskObjective{}, nil
}

func CreateRiskObjectiveFromConfig(riskConfigs RiskConfigs) (*RiskObjective, error) {
	riskObj := &RiskObjective{
		HazardInteractionMatrix: riskConfigs.HazardInteractionMatrix,
		Delta:                   riskConfigs.Delta,
	}
	return riskObj, nil
}

func (obj *RiskObjective) Eval() float64 {
	return 0
}
