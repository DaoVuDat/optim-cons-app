package conslay_continuous

const (
	RiskObjectiveType       = "Risk Objective"
	HazardInteractionMatrix = "Hazard Interaction Matrix"
	Delta                   = "Delta"
	AlphaRiskPenalty        = "AlphaPenalty"
)

type RiskConfigs struct {
	HazardInteractionMatrix [][]float64
	Delta                   float64
	AlphaRiskPenalty        float64
}

type RiskObjective struct {
	HazardInteractionMatrix [][]float64
	Delta                   float64
	AlphaRiskPenalty        float64
}

func CreateRiskObjective() (*RiskObjective, error) {
	return &RiskObjective{}, nil
}

func CreateRiskObjectiveFromConfig(riskConfigs RiskConfigs) (*RiskObjective, error) {
	riskObj := &RiskObjective{
		HazardInteractionMatrix: riskConfigs.HazardInteractionMatrix,
		Delta:                   riskConfigs.Delta,
		AlphaRiskPenalty:        riskConfigs.AlphaRiskPenalty,
	}
	return riskObj, nil
}

func (obj *RiskObjective) Eval() float64 {
	return 0
}

func (obj *RiskObjective) GetAlphaPenalty() float64 {
	return obj.AlphaRiskPenalty
}
