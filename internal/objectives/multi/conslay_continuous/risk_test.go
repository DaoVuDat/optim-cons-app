package conslay_continuous

import (
	"golang-moaha-construction/internal/util"
	"log"
	"testing"
)

func CreateSideInputLocation() map[string]Location {

	// coordinate

	return map[string]Location{
		"TF1": Location{
			Coordinate: Coordinate{X: 34.551430, Y: 7.306936},
			Rotation:   true,
			Length:     5,
			Width:      12,
			IsFixed:    false,
			Symbol:     "TF1",
			Name:       "Pile storage yard #1",
		},
		"TF2": Location{
			Coordinate: Coordinate{X: 108.311238, Y: 8.959338},
			Rotation:   true,
			Length:     5,
			Width:      12,
			IsFixed:    false,
			Symbol:     "TF2",
			Name:       "Pile storage yard #2",
		},

		"TF3": Location{
			Coordinate: Coordinate{X: 4.369674, Y: 7.260322},
			Rotation:   false,
			Length:     8,
			Width:      14,
			IsFixed:    false,
			Symbol:     "TF3",
			Name:       "Site office",
		},
		"TF4": Location{
			Coordinate: Coordinate{X: 104.347527, Y: 38.058604},
			Rotation:   false,
			Length:     14,
			Width:      7,
			IsFixed:    false,
			Symbol:     "TF4",
			Name:       "Rebar process yard",
		},

		"TF5": Location{
			Coordinate: Coordinate{X: 84.207558, Y: 64.782349},
			Rotation:   true,
			Length:     7,
			Width:      14,
			IsFixed:    false,
			Symbol:     "TF5",
			Name:       "Formwork process yard",
		},

		"TF6": Location{
			Coordinate: Coordinate{X: 117.818767, Y: 92.673680},
			Rotation:   false,
			Length:     4,
			Width:      4,
			IsFixed:    false,
			Symbol:     "TF6",
			Name:       "Electrician hut",
		},

		"TF7": Location{
			Coordinate: Coordinate{X: 32.475737, Y: 6.784953},
			Rotation:   false,
			Length:     10,
			Width:      12,
			IsFixed:    false,
			Symbol:     "TF7",
			Name:       "Ready-mix concrete area",
		},

		"TF8": Location{
			Coordinate: Coordinate{X: 101.023480, Y: 18.103534},
			Rotation:   true,
			Length:     6,
			Width:      12,
			IsFixed:    false,
			Symbol:     "TF8",
			Name:       "Prefabricated components yard #1 (slab, staircase)",
		},

		"TF9": Location{
			Coordinate: Coordinate{X: 47.163815, Y: 37.001098},
			Rotation:   false,
			Length:     12,
			Width:      6,
			IsFixed:    false,
			Symbol:     "TF9",
			Name:       "Prefabricated components yard #2 (beam, column)",
		},

		"TF10": Location{
			Coordinate: Coordinate{X: 101.054167, Y: 17.159174},
			Rotation:   true,
			Length:     6,
			Width:      10,
			IsFixed:    false,
			Symbol:     "TF10",
			Name:       "Prefabricated components yard #3 (external wall)",
		},
		"TF11": Location{
			Coordinate: Coordinate{X: 3.567509, Y: 87.419745},
			Rotation:   false,
			Length:     6,
			Width:      4,
			IsFixed:    false,
			Symbol:     "TF11",
			Name:       "Dangerous goods warehouse",
		},
		"TF12": Location{
			Coordinate: Coordinate{X: 12.545486, Y: 80.868169},
			Rotation:   false,
			Length:     8,
			Width:      6,
			IsFixed:    false,
			Symbol:     "TF12",
			Name:       "ME storage warehouse",
		},
		"TF13": Location{
			Coordinate: Coordinate{X: 72, Y: 22},
			Rotation:   false,
			Length:     52,
			Width:      24,
			IsFixed:    true,
			Symbol:     "TF13",
			Name:       "Building",
		},
		"TF14": {
			Coordinate: Coordinate{X: 72, Y: 36},
			Rotation:   false,
			Length:     2,
			Width:      2,
			IsFixed:    true,
			Symbol:     "TF14",
			Name:       "Tower crane",
		},
		"TF15": {
			Coordinate: Coordinate{X: 43, Y: 22},
			Rotation:   false,
			Length:     3,
			Width:      3,
			IsFixed:    true,
			Symbol:     "TF15",
			Name:       "Hoist",
		},
		"TF16": {
			Coordinate: Coordinate{X: 5, Y: 92},
			Rotation:   false,
			Length:     3,
			Width:      4,
			IsFixed:    false,
			Symbol:     "TF16",
			Name:       "Security room",
		},
	}
}

func TestRiskObjective_Eval(t *testing.T) {
	testTable := []struct {
		locations map[string]Location
		expected  float64
		name      string
	}{
		{
			locations: CreateInputLocation(true),
			expected:  1038.48,
			name:      "mostly feasible locations",
		},
		{
			locations: CreateInputLocation(false),
			expected:  1014.62,
			name:      "infeasible locations",
		},
		{
			locations: CreateSideInputLocation(),
			expected:  918.89,
			name:      "test locations",
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
