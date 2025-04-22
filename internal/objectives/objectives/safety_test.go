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

	safetyProximity, err := ReadSafetyProximityDataFromFile("../../../data/conslay/continuous/safety_data.xlsx")

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

func CreateInputMini() map[string]data.Location {
	locations := make(map[string]data.Location)

	locations = map[string]data.Location{
		"TF1": data.Location{
			Coordinate: data.Coordinate{X: 60.2117021287296, Y: 70.6492302951395},
			Rotation:   true,
			Length:     5,
			Width:      12,
			IsFixed:    false,
			Symbol:     "TF1",
			Name:       "Pile storage yard #1",
		},
		"TF2": data.Location{
			Coordinate: data.Coordinate{X: 70.8063495854695, Y: 20.8484429816236},
			Rotation:   true,
			Length:     5,
			Width:      12,
			IsFixed:    false,
			Symbol:     "TF2",
			Name:       "Pile storage yard #2",
		},
		"TF3": data.Location{
			Coordinate: data.Coordinate{X: 16.0620484334653, Y: 11.2728104676465},
			Rotation:   false,
			Length:     8,
			Width:      14,
			IsFixed:    false,
			Symbol:     "TF3",
			Name:       "Site office",
		},
		"TF4": data.Location{
			Coordinate: data.Coordinate{X: 72, Y: 22},
			Rotation:   false,
			Length:     8,
			Width:      14,
			IsFixed:    false,
			Symbol:     "TF4",
			Name:       "Site office 2",
		},
	}

	return locations
}

func CreateInputPhasesMini() [][]string {
	return [][]string{
		{"TF1", "TF2", "TF4"},
		{"TF2", "TF3", "TF4"},
		{"TF3", "TF4"},
	}
}

func TestSafetyObjectiveMini_Eval(t *testing.T) {
	testTable := []struct {
		locations map[string]data.Location
		expected  float64
		name      string
	}{
		{
			locations: CreateInputMini(),
			expected:  3442.62,
			name:      "mostly feasible locations",
		},
	}

	safetyProximity, err := ReadSafetyProximityDataFromFile("../../../data/conslay/mini/safety_data.xlsx")

	// Safety Objective Configs
	safetyConfig := SafetyConfigs{
		SafetyProximity:    safetyProximity,
		AlphaSafetyPenalty: 100,
		Phases:             CreateInputPhasesMini(),
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
