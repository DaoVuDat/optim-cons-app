package cons_lay

const (
	RiskObjectiveType       = "Risk Objective"
	HazardInteractionMatrix = "Hazard Interaction Matrix"
	Delta                   = "Delta"
)

type RiskObjective struct {
	HazardInteractionMatrix [][]float64
	Delta                   float64
}

func (obj *RiskObjective) Eval() float64 {
	return 0
}
