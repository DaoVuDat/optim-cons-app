package cons_lay

import "golang-moaha-construction/internal/data"

const (
	RiskObjectiveType       = "Risk Objective"
	HazardInteractionMatrix = "Hazard Interaction Matrix"
	Delta                   = "Delta"
)

type RiskObjective struct {
}

var RiskConfigs = []*data.Config{
	{
		Name: HazardInteractionMatrix,
	},
	{
		Name: Delta,
	},
}

func (obj *RiskObjective) LoadData(configs []data.Config) error {

	return nil
}
