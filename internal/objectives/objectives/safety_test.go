package objectives

import (
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/util"
	"log"
	"testing"
)

func TestSafetyObjective_Eval(t *testing.T) {
	testTable := []struct {
		locations map[string]data.Location
		expected  float64
		name      string
	}{
		{
			locations: CreateInputLocation(true),
			expected:  37675,
			name:      "mostly feasible locations",
		},
		{
			locations: CreateInputLocation(false),
			expected:  43572.47,
			name:      "infeasible locations",
		},
		{
			locations: CreateSideInputLocation(),
			expected:  54485.47,
			name:      "test locations",
		},
	}

	safetyProximity, err := ReadSafetyProximityDataFromFile("../../../data/conslay/safety_data.xlsx")

	// Safety Objective Configs
	safetyConfig := SafetyConfigs{
		SafetyProximity:    safetyProximity,
		AlphaSafetyPenalty: 100,
		Phases:             CreateInputPhases(),
		FilePath:           "../../../data/conslay/safety_data.xlsx",
	}
	safetyObj, err := CreateSafetyObjectiveFromConfig(safetyConfig)
	if err != nil {
		log.Fatal(err)
	}

	// calculate result
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {

			result := safetyObj.Eval(test.locations)
			if util.RoundTo(result, 2) != test.expected {
				t.Errorf("expected result to be %f, got %f", test.expected, result)
			}
		})
	}
}
