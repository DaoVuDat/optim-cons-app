package objectives

import (
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/util"
	"log"
	"testing"
)

func TestSafetyHazardObjectiveMini_Eval(t *testing.T) {
	testTable := []struct {
		locations map[string]data.Location
		expected  float64
		name      string
	}{
		{
			locations: CreateInputMini(),
			expected:  -43.73,
			name:      "mini locations",
		},
	}

	seMatrix, err := ReadSafetyAndEnvDataFromFile("../../../data/conslay/mini/safety_hazard_data.xlsx")

	// Safety Objective Configs
	safetyHazardConfig := SafetyHazardConfigs{
		SEMatrix:       seMatrix,
		AlphaSHPenalty: 100,
		Phases:         CreateInputPhasesMini(),
	}
	safetyHazardObj, err := CreateSafetyHazardObjectiveFromConfig(safetyHazardConfig)
	if err != nil {
		log.Fatal(err)
	}

	// calculate result
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {

			result := safetyHazardObj.Eval(test.locations)
			if util.RoundTo(result, 2) != test.expected {
				t.Errorf("expected result to be %f, got %f", test.expected, result)
			}

		})
	}
}
