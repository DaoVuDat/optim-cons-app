package objectives

import (
	"golang-moaha-construction/internal/data"
	"log"
	"math"
	"testing"
)

func createInputLocation(caseNumber int) map[string]data.Location {
	switch caseNumber {
	case 0:
		locations := make(map[string]data.Location)

		locations["TF1"] = data.Location{
			Symbol:      "TF1",
			IsLocatedAt: "L11",
		}

		locations["TF2"] = data.Location{
			Symbol:      "TF2",
			IsLocatedAt: "L5",
		}

		locations["TF3"] = data.Location{
			Symbol:      "TF3",
			IsLocatedAt: "L9",
		}

		locations["TF4"] = data.Location{
			Symbol:      "TF4",
			IsLocatedAt: "L7",
		}

		locations["TF5"] = data.Location{
			Symbol:      "TF5",
			IsLocatedAt: "L2",
		}

		locations["TF6"] = data.Location{
			Symbol:      "TF6",
			IsLocatedAt: "L8",
		}

		locations["TF7"] = data.Location{
			Symbol:      "TF7",
			IsLocatedAt: "L3",
		}

		locations["TF8"] = data.Location{
			Symbol:      "TF8",
			IsLocatedAt: "L1",
		}

		locations["TF9"] = data.Location{
			Symbol:      "TF9",
			IsLocatedAt: "L6",
		}

		locations["TF10"] = data.Location{
			Symbol:      "TF10",
			IsLocatedAt: "L4",
		}

		locations["TF11"] = data.Location{
			Symbol:      "TF11",
			IsLocatedAt: "L10",
		}

		return locations
	case 1:
		locations := make(map[string]data.Location)

		locations["TF1"] = data.Location{
			Symbol:      "TF1",
			IsLocatedAt: "L11",
		}

		locations["TF2"] = data.Location{
			Symbol:      "TF2",
			IsLocatedAt: "L2",
		}

		locations["TF3"] = data.Location{
			Symbol:      "TF3",
			IsLocatedAt: "L5",
		}

		locations["TF4"] = data.Location{
			Symbol:      "TF4",
			IsLocatedAt: "L8",
		}

		locations["TF5"] = data.Location{
			Symbol:      "TF5",
			IsLocatedAt: "L9",
		}

		locations["TF6"] = data.Location{
			Symbol:      "TF6",
			IsLocatedAt: "L6",
		}

		locations["TF7"] = data.Location{
			Symbol:      "TF7",
			IsLocatedAt: "L4",
		}

		locations["TF8"] = data.Location{
			Symbol:      "TF8",
			IsLocatedAt: "L3",
		}

		locations["TF9"] = data.Location{
			Symbol:      "TF9",
			IsLocatedAt: "L7",
		}

		locations["TF10"] = data.Location{
			Symbol:      "TF10",
			IsLocatedAt: "L10",
		}

		locations["TF11"] = data.Location{
			Symbol:      "TF11",
			IsLocatedAt: "L1",
		}

		return locations
	case 2:
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

		return locations
	case 3:
		locations := make(map[string]data.Location)

		locations["TF1"] = data.Location{
			Symbol:      "TF1",
			IsLocatedAt: "L9",
		}

		locations["TF2"] = data.Location{
			Symbol:      "TF2",
			IsLocatedAt: "L8",
		}

		locations["TF3"] = data.Location{
			Symbol:      "TF3",
			IsLocatedAt: "L4",
		}

		locations["TF4"] = data.Location{
			Symbol:      "TF4",
			IsLocatedAt: "L7",
		}

		locations["TF5"] = data.Location{
			Symbol:      "TF5",
			IsLocatedAt: "L5",
		}

		locations["TF6"] = data.Location{
			Symbol:      "TF6",
			IsLocatedAt: "L6",
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

		return locations
	default:
		return nil
	}
}

func TestReadMatrixFromFile_DistanceMatrix(t *testing.T) {
	matrix, err := ReadMatrixFromFile("../../../data/conslay/predetermined/distance_data.xlsx")

	if err != nil {
		log.Fatal(err)
	}

	if matrix.GetNumberOfItems() != 11 {
		t.Errorf("expected result to be %d, got %d", 11, matrix.GetNumberOfItems())
	}

	// access by idx, row = 4 , col = 10
	val, err := matrix.GetCellValueFromIdx(4, 10)
	if err != nil {
		log.Fatal(err)
	}

	if math.Abs(val-52.0) > 1e-9 {
		t.Errorf("expected result to be %f, got %f", 52.0, val)
	}

	// access by name
	val, err = matrix.GetCellValueFromNames("L5", "L11")

	if err != nil {
		log.Fatal(err)
	}

	if math.Abs(val-52.0) > 1e-9 {
		t.Errorf("expected result to be %f, got %f", 52.0, val)
	}
}

func TestReadMatrixFromFile_FrequencyMatrix(t *testing.T) {
	matrix, err := ReadMatrixFromFile("../../../data/conslay/predetermined/frequency_data.xlsx")

	if err != nil {
		log.Fatal(err)
	}

	if matrix.GetNumberOfItems() != 11 {
		t.Errorf("expected result to be %d, got %d", 11, matrix.GetNumberOfItems())
	}

	// access by idx, row = 4 , col = 10
	val, err := matrix.GetCellValueFromIdx(4, 10)
	if err != nil {
		log.Fatal(err)
	}

	if math.Abs(val-6.0) > 1e-9 {
		t.Errorf("expected result to be %f, got %f", 6.0, val)
	}

	// access by name
	val, err = matrix.GetCellValueFromNames("TF5", "TF11")

	if err != nil {
		log.Fatal(err)
	}

	if math.Abs(val-6.0) > 1e-9 {
		t.Errorf("expected result to be %f, got %f", 6.0, val)
	}
}

func TestConstructionCostObjective_Eval(t *testing.T) {
	testTable := []struct {
		locations map[string]data.Location
		expected  float64
		name      string
	}{
		{
			locations: createInputLocation(0),
			expected:  15094,
			name:      "Li and Love's results",
		},
		{
			locations: createInputLocation(1),
			expected:  12320,
			name:      "best optimal solution",
		},
	}

	freqMatrix, err := ReadMatrixFromFile("../../../data/conslay/predetermined/frequency_data.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	distanceMatrix, err := ReadMatrixFromFile("../../../data/conslay/predetermined/distance_data.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ccObj, err := CreateConstructionCostObjectiveFromConfig(ConstructionCostConfigs{
				FrequencyMatrix: freqMatrix,
				DistanceMatrix:  distanceMatrix,
				FullRun:         true,
			})

			if err != nil {
				log.Fatal(err)
			}

			result := ccObj.Eval(tt.locations)
			if math.Abs(tt.expected-result) > 1e-9 {
				t.Errorf("expected result to be %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestConstructionCostObjectiveFacLessLoc_Eval(t *testing.T) {
	testTable := []struct {
		locations map[string]data.Location
		expected  float64
		name      string
	}{
		{
			locations: createInputLocation(2),
			expected:  831.94,
			name:      "kaveh2018 (1)",
		},
	}

	freqMatrix, err := ReadMatrixFromFile("../../../data/conslay/predetermined/frequency_1_data.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	distanceMatrix, err := ReadMatrixFromFile("../../../data/conslay/predetermined/distance_1_data.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			ccObj, err := CreateConstructionCostObjectiveFromConfig(ConstructionCostConfigs{
				FrequencyMatrix: freqMatrix,
				DistanceMatrix:  distanceMatrix,
				FullRun:         false,
			})

			if err != nil {
				log.Fatal(err)
			}

			result := ccObj.Eval(tt.locations)
			if math.Abs(tt.expected-result) > 1e-9 {
				t.Errorf("expected result to be %f, got %f", tt.expected, result)
			}
		})
	}
}
