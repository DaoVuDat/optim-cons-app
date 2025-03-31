package conslay_continuous

import (
	"golang-moaha-construction/internal/util"
	"log"
	"testing"
)

func TestRiskObjective_Eval(t *testing.T) {
	testTable := []struct {
		locations map[string]Location
		expected  float64
		name      string
	}{
		{
			locations: CreateInputLocation(true),
			expected:  1930.91,
			name:      "mostly feasible locations",
		},
		{
			locations: CreateInputLocation(false),
			expected:  1906.24,
			name:      "infeasible locations",
		},
	}

	hazardInteraction, err := ReadRiskHazardInteractionDataFromFile("../../../../data/conslay/f2_risk_data.xlsx")

	// Hoisting Objective Configs
	riskConfigs := RiskConfigs{
		HazardInteractionMatrix: hazardInteraction,
		Delta:                   0.01,
		AlphaRiskPenalty:        100,
		Phases:                  CreateInputPhases(),
	}
	hoistingObj, err := CreateRiskObjectiveFromConfig(riskConfigs)
	if err != nil {
		log.Fatal(err)
	}

	// calculate result
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {

			result := hoistingObj.Eval(test.locations)
			if util.RoundTo(result, 2) != test.expected {
				t.Errorf("expected result to be %f, got %f", test.expected, result)
			}
		})
	}
}
