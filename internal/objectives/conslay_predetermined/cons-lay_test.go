package conslay_predetermined

import (
	"golang-moaha-construction/internal/data"
	"log"
	"slices"
	"testing"
)

func inputCase(caseNumber int) ([]float64, map[string]data.Location) {
	switch caseNumber {
	case 0:
		locations := make(map[string]data.Location)

		locations["TF1"] = data.Location{
			Symbol:      "TF1",
			IsLocatedAt: "L10",
		}

		locations["TF2"] = data.Location{
			Symbol:      "TF2",
			IsLocatedAt: "L5",
		}

		locations["TF3"] = data.Location{
			Symbol:      "TF3",
			IsLocatedAt: "L6",
		}

		locations["TF4"] = data.Location{
			Symbol:      "TF4",
			IsLocatedAt: "L7",
		}

		locations["TF5"] = data.Location{
			Symbol:      "TF5",
			IsLocatedAt: "L9",
		}

		locations["TF6"] = data.Location{
			Symbol:      "TF6",
			IsLocatedAt: "L8",
		}

		locations["TF7"] = data.Location{
			Symbol:      "TF7",
			IsLocatedAt: "L11",
		}

		locations["TF8"] = data.Location{
			Symbol:      "TF8",
			IsLocatedAt: "L12",
		}

		locations["TF9"] = data.Location{
			Symbol:      "TF9",
			IsLocatedAt: "L13",
		}
		return []float64{0.5, 0.6, 0.7, 0.8, 0.1, 0.15, 0.2, 0.4, 0.3, 0.05}, locations
	default:
		return []float64{}, nil
	}
}

func TestConsLay_Eval(t *testing.T) {
	inputValues, expectedMapLocations := inputCase(0)

	config := ConsLayConfigs{
		NumberOfLocations:  13,
		NumberOfFacilities: 9,
		FixedFacilitiesName: []LocFac{
			{
				LocName: "L11",
				FacName: "TF7",
			},
			{
				LocName: "L12",
				FacName: "TF8",
			},
			{
				LocName: "L13",
				FacName: "TF9",
			},
		},
	}
	expectedAvailableLocations := []string{"L1", "L2", "L3", "L4", "L5", "L6", "L7", "L8", "L9", "L10"}
	expectedFacilitiesToBeFound := []string{"TF1", "TF2", "TF3", "TF4", "TF5", "TF6"}

	obj, err := CreateConsLayFromConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// checking dimensions
	if obj.Dimensions != 10 {
		t.Errorf("expected dimensions to be 10, got %d", obj.Dimensions)
	}

	// checking available locations
	availableLocations := obj.AvailableLocationsIdx
	if !slices.Equal(availableLocations, expectedAvailableLocations) {
		t.Errorf("expected available locations to be %v, got %v", expectedAvailableLocations, availableLocations)
	}

	// checking facilities to be found
	facilitiesToBeFound := obj.FacilitiesToBeFound

	if !slices.Equal(facilitiesToBeFound, expectedFacilitiesToBeFound) {
		t.Errorf("expected facilities to be found to be %v, got %v", expectedFacilitiesToBeFound, facilitiesToBeFound)
	}

	mapLocations := obj.MappingLocations(inputValues)

	for k, v := range expectedMapLocations {
		vMap := mapLocations[k]
		if vMap.Symbol != v.Symbol {
			t.Errorf("expected symbol to be %s, got %s", v.Symbol, vMap.Symbol)
		}

		if vMap.IsLocatedAt != v.IsLocatedAt {
			t.Errorf("expected IsLocatedAt to be %s, got %s", v.IsLocatedAt, vMap.IsLocatedAt)
		}
	}

}
