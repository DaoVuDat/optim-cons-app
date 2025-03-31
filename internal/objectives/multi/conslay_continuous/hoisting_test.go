package conslay_continuous

import (
	"golang-moaha-construction/internal/util"
	"log"
	"testing"
)

func TestHoistingObjective_Eval(t *testing.T) {

	testTable := []struct {
		locations map[string]Location
		expected  float64
		name      string
	}{
		{
			locations: CreateInputLocation(true),
			expected:  39079.38,
			name:      "mostly feasible locations",
		},
		{
			locations: CreateInputLocation(false),
			expected:  42561.22,
			name:      "infeasible locations",
		},
	}

	hoistingTime, err := ReadHoistingTimeDataFromFile("../../../../data/conslay/f1_hoisting_time_data.xlsx")

	if err != nil {
		log.Fatal(err)
	}

	// Hoisting Objective Configs
	hoistingConfigs := HoistingConfigs{
		//PrefabricatedLocations: []string{"TF8", "TF9", "TF10"},
		NumberOfFloors: 10,
		HoistingTime: map[string][]HoistingTime{
			"TF14": hoistingTime,
		},
		FloorHeight:          3.2,
		ZM:                   2,
		Vuvg:                 37.5,
		Vlvg:                 37.5 / 2,
		Vag:                  50,
		Vwg:                  0.5,
		AlphaHoistingPenalty: 1,
		AlphaHoisting:        0.25,
		BetaHoisting:         1, // beta hoisting = n hoisting
		NHoisting:            1,
	}
	hoistingObj, err := CreateHoistingObjectiveFromConfig(hoistingConfigs)
	if err != nil {
		log.Fatal(err)
	}

	// calculate result
	for _, test := range testTable {
		craneLocations := make([]Crane, 0)
		craneLocations = append(craneLocations, Crane{
			Location:     test.locations["TF14"],
			BuildingName: []string{"TF8", "TF9", "TF10"},
			Radius:       40,
		})

		hoistingObj.CraneLocations = craneLocations
		t.Run(test.name, func(t *testing.T) {

			result := hoistingObj.Eval(test.locations)
			if util.RoundTo(result, 2) != test.expected {
				t.Errorf("expected result to be %f, got %f", test.expected, result)
			}
		})
	}

}
