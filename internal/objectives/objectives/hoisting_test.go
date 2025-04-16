package objectives

import (
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/util"
	"log"
	"testing"
)

func CreateInputLocation(feasible bool) map[string]data.Location {
	locations := make(map[string]data.Location)
	if feasible {
		locations = map[string]data.Location{
			"TF1": data.Location{
				Coordinate: data.Coordinate{X: 41.8984994960664, Y: 16.0789402087708},
				Rotation:   false,
				Length:     12,
				Width:      5,
				IsFixed:    false,
				Symbol:     "TF1",
				Name:       "Pile storage yard #1",
			},
			"TF2": data.Location{
				Coordinate: data.Coordinate{X: 40.3018087791441, Y: 6.54792941746860},
				Rotation:   true,
				Length:     5,
				Width:      12,
				IsFixed:    false,
				Symbol:     "TF2",
				Name:       "Pile storage yard #2",
			},
			"TF3": data.Location{
				Coordinate: data.Coordinate{X: 43.0538143534362, Y: 81.9094405733262},
				Rotation:   false,
				Length:     8,
				Width:      14,
				IsFixed:    false,
				Symbol:     "TF3",
				Name:       "Site office",
			},
			"TF4": data.Location{
				Coordinate: data.Coordinate{X: 49.7610024122533, Y: 57.2431089610756},
				Rotation:   true,
				Length:     7,
				Width:      14,
				IsFixed:    false,
				Symbol:     "TF4",
				Name:       "Rebar process yard",
			},
			"TF5": data.Location{
				Coordinate: data.Coordinate{X: 81.3261517318950, Y: 41.8071744013404},
				Rotation:   false,
				Length:     14,
				Width:      7,
				IsFixed:    false,
				Symbol:     "TF5",
				Name:       "Formwork process yard",
			},
			"TF6": data.Location{
				Coordinate: data.Coordinate{X: 101.490485902993, Y: 11.1774705429599},
				Rotation:   true,
				Length:     4,
				Width:      4,
				IsFixed:    false,
				Symbol:     "TF6",
				Name:       "Electrician hut",
			},
			"TF7": data.Location{
				Coordinate: data.Coordinate{X: 43.6721024451348, Y: 9.42961655802787},
				Rotation:   true,
				Length:     12,
				Width:      10,
				IsFixed:    false,
				Symbol:     "TF7",
				Name:       "Ready-mix concrete area",
			},
			"TF8": data.Location{
				Coordinate: data.Coordinate{X: 48.5660821450197, Y: 54.0455646998945},
				Rotation:   false,
				Length:     12,
				Width:      6,
				IsFixed:    false,
				Symbol:     "TF8",
				Name:       "Prefabricated components yard #1 (slab, staircase)",
			},
			"TF9": data.Location{
				Coordinate: data.Coordinate{X: 76.4180335261621, Y: 41.8787903034157},
				Rotation:   false,
				Length:     12,
				Width:      6,
				IsFixed:    false,
				Symbol:     "TF9",
				Name:       "Prefabricated components yard #2 (beam, column)",
			},
			"TF10": data.Location{
				Coordinate: data.Coordinate{X: 108.028665273018, Y: 36.9845037867221},
				Rotation:   true,
				Length:     6,
				Width:      10,
				IsFixed:    false,
				Symbol:     "TF10",
				Name:       "Prefabricated components yard #3 (external wall)",
			},
			"TF11": data.Location{
				Coordinate: data.Coordinate{X: 63.8375790712078, Y: 91.2094786904858},
				Rotation:   true,
				Length:     4,
				Width:      6,
				IsFixed:    false,
				Symbol:     "TF11",
				Name:       "Dangerous goods warehouse",
			},
			"TF12": data.Location{
				Coordinate: data.Coordinate{X: 110.579901257085, Y: 14.7798350113455},
				Rotation:   false,
				Length:     8,
				Width:      6,
				IsFixed:    false,
				Symbol:     "TF12",
				Name:       "ME storage warehouse",
			},
			"TF13": data.Location{
				Coordinate: data.Coordinate{X: 72, Y: 22},
				Rotation:   false,
				Length:     52,
				Width:      24,
				IsFixed:    true,
				Symbol:     "TF13",
				Name:       "Building",
			},
			"TF14": {
				Coordinate: data.Coordinate{X: 72, Y: 36},
				Rotation:   false,
				Length:     2,
				Width:      2,
				IsFixed:    true,
				Symbol:     "TF14",
				Name:       "Tower crane",
			},
			"TF15": {
				Coordinate: data.Coordinate{X: 43, Y: 22},
				Rotation:   false,
				Length:     3,
				Width:      3,
				IsFixed:    true,
				Symbol:     "TF15",
				Name:       "Hoist",
			},
			"TF16": {
				Coordinate: data.Coordinate{X: 5, Y: 92},
				Rotation:   false,
				Length:     3,
				Width:      4,
				IsFixed:    false,
				Symbol:     "TF16",
				Name:       "Security room",
			},
		}
	} else {
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
				Coordinate: data.Coordinate{X: 66.8271432818314, Y: 79.6322545348163},
				Rotation:   true,
				Length:     7,
				Width:      14,
				IsFixed:    false,
				Symbol:     "TF4",
				Name:       "Rebar process yard",
			},
			"TF5": data.Location{
				Coordinate: data.Coordinate{X: 97.4966707756999, Y: 13.6022897445312},
				Rotation:   false,
				Length:     14,
				Width:      7,
				IsFixed:    false,
				Symbol:     "TF5",
				Name:       "Formwork process yard",
			},
			"TF6": data.Location{
				Coordinate: data.Coordinate{X: 1.51514052087480, Y: 80.0719162226562},
				Rotation:   false,
				Length:     4,
				Width:      4,
				IsFixed:    false,
				Symbol:     "TF6",
				Name:       "Electrician hut",
			},
			"TF7": data.Location{
				Coordinate: data.Coordinate{X: 36.7111938858722, Y: 11.9943304146193},
				Rotation:   true,
				Length:     12,
				Width:      10,
				IsFixed:    false,
				Symbol:     "TF7",
				Name:       "Ready-mix concrete area",
			},
			"TF8": data.Location{
				Coordinate: data.Coordinate{X: 48.6579882916605, Y: 61.3097045687976},
				Rotation:   false,
				Length:     12,
				Width:      6,
				IsFixed:    false,
				Symbol:     "TF8",
				Name:       "Prefabricated components yard #1 (slab, staircase)",
			},
			"TF9": data.Location{
				Coordinate: data.Coordinate{X: 7.04418109569817, Y: 87.8287852386491},
				Rotation:   true,
				Length:     6,
				Width:      12,
				IsFixed:    false,
				Symbol:     "TF9",
				Name:       "Prefabricated components yard #2 (beam, column)",
			},
			"TF10": data.Location{
				Coordinate: data.Coordinate{X: 42.2938936308434, Y: 65.2135456501644},
				Rotation:   false,
				Length:     10,
				Width:      6,
				IsFixed:    false,
				Symbol:     "TF10",
				Name:       "Prefabricated components yard #3 (external wall)",
			},
			"TF11": data.Location{
				Coordinate: data.Coordinate{X: 21.9910306881593, Y: 85.0044587159497},
				Rotation:   false,
				Length:     6,
				Width:      4,
				IsFixed:    false,
				Symbol:     "TF11",
				Name:       "Dangerous goods warehouse",
			},
			"TF12": data.Location{
				Coordinate: data.Coordinate{X: 23.2731387385839, Y: 90.0053928795082},
				Rotation:   true,
				Length:     6,
				Width:      8,
				IsFixed:    false,
				Symbol:     "TF12",
				Name:       "ME storage warehouse",
			},
			"TF13": data.Location{
				Coordinate: data.Coordinate{X: 72, Y: 22},
				Rotation:   false,
				Length:     52,
				Width:      24,
				IsFixed:    true,
				Symbol:     "TF13",
				Name:       "Building",
			},
			"TF14": {
				Coordinate: data.Coordinate{X: 72, Y: 36},
				Rotation:   false,
				Length:     2,
				Width:      2,
				IsFixed:    true,
				Symbol:     "TF14",
				Name:       "Tower crane",
			},
			"TF15": {
				Coordinate: data.Coordinate{X: 43, Y: 22},
				Rotation:   false,
				Length:     3,
				Width:      3,
				IsFixed:    true,
				Symbol:     "TF15",
				Name:       "Hoist",
			},
			"TF16": {
				Coordinate: data.Coordinate{X: 5, Y: 92},
				Rotation:   false,
				Length:     3,
				Width:      4,
				IsFixed:    false,
				Symbol:     "TF16",
				Name:       "Security room",
			},
		}
	}

	return locations
}

func CreateInputPhases() [][]string {
	return [][]string{
		{"TF1", "TF2", "TF3", "TF6", "TF13", "TF16"},
		{"TF3", "TF4", "TF5", "TF6", "TF7", "TF13", "TF16"},
		{"TF3", "TF4", "TF5", "TF6", "TF7", "TF11", "TF13", "TF14", "TF15", "TF16"},
		{"TF3", "TF6", "TF8", "TF9", "TF11", "TF13", "TF14", "TF15", "TF16"},
		{"TF3", "TF6", "TF10", "TF11", "TF12", "TF13", "TF14", "TF15", "TF16"},
	}
}

func TestHoistingObjective_Eval(t *testing.T) {

	testTable := []struct {
		locations map[string]data.Location
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

	hoistingTime, err := ReadHoistingTimeDataFromFile("../../../data/conslay/hoisting_time_data.xlsx")

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
			CraneSymbol:  "TF14",
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
