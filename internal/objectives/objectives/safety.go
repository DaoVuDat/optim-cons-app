package objectives

import (
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/data"
	"strconv"
)

const SafetyObjectiveType data.ObjectiveType = "Safety Objective"

type SafetyConfigs struct {
	SafetyProximity    [][]float64
	AlphaSafetyPenalty float64
	Phases             [][]string
	FilePath           string
}

type SafetyObjective struct {
	SafetyProximity    [][]float64
	AlphaSafetyPenalty float64
	Phases             [][]string
	FilePath           string
}

func CreateSafetyObjectiveFromConfig(safetyConfigs SafetyConfigs) (*SafetyObjective, error) {
	safetyObj := &SafetyObjective{
		SafetyProximity:    safetyConfigs.SafetyProximity,
		AlphaSafetyPenalty: safetyConfigs.AlphaSafetyPenalty,
		Phases:             safetyConfigs.Phases,
		FilePath:           safetyConfigs.FilePath,
	}
	return safetyObj, nil
}

func (obj *SafetyObjective) Eval(locations map[string]data.Location) float64 {
	result := 0.0

	for _, phases := range obj.Phases {

		for i := 0; i < len(phases)-1; i++ {
			facilityNameI := phases[i]
			facilityI := locations[facilityNameI]
			idxI, err := facilityI.ConvertToIdx()
			if err != nil {
				return 0
			}

			for j := i + 1; j < len(phases); j++ {
				facilityNameJ := phases[j]
				facilityJ := locations[facilityNameJ]
				idxJ, err := facilityJ.ConvertToIdx()
				if err != nil {
					return 0
				}

				result += obj.SafetyProximity[idxI][idxJ] * data.Distance2D(facilityI.Coordinate, facilityJ.Coordinate)
			}
		}
	}

	return result
}

func (obj *SafetyObjective) GetAlphaPenalty() float64 {
	return obj.AlphaSafetyPenalty
}

func ReadSafetyProximityDataFromFile(filePath string) ([][]float64, error) {
	dataFile, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	rows, err := dataFile.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	safetyProximity := make([][]float64, len(rows)-1)

	for idx, row := range rows {
		// skip header
		if idx == 0 {
			continue
		}

		arr := make([]float64, len(rows)-1)
		for i, cell := range row {
			if i <= idx {
				continue
			}

			// skip the first column
			arr[i-1], err = strconv.ParseFloat(cell, 64)
			if err != nil {
				return nil, err
			}

		}
		safetyProximity[idx-1] = arr
	}

	return safetyProximity, nil
}
