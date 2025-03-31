package conslay_continuous

import (
	"log"
	"testing"
)

func TestHoistingObjective_Eval(t *testing.T) {
	locations := CreateInputLocation(false)

	hoistingTime, err := ReadHoistingTimeDataFromFile("./data/conslay/f1_hoisting_time_data.xlsx")

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

	craneLocations := make([]Crane, 0)
	craneLocations = append(craneLocations, Crane{
		Location:     locations["TF14"],
		BuildingName: []string{"TF8", "TF9", "TF10"},
		Radius:       40,
	})

	hoistingObj.CraneLocations = craneLocations

	// calculate result
	result := hoistingObj.Eval(locations)

	// TODO: check this
	if result != 0.0 {
		t.Errorf("expected result to be 0, got %f", result)
	}
}
