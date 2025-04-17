package objectives

import (
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/data"
	"strconv"
)

const SafetyHazardObjectiveType data.ObjectiveType = "Safety Hazard Objective"

type SafetyHazardConfigs struct {
	SEMatrix       [][]float64
	AlphaSHPenalty float64
	Phases         [][]string
	FilePath       string
}

type SafetyHazardObjective struct {
	SEMatrix       [][]float64
	AlphaSHPenalty float64
	Phases         [][]string
	FilePath       string
}

func CreateSafetyHazardObjectiveFromConfig(hsConfigs SafetyHazardConfigs) (*SafetyHazardObjective, error) {
	tcObj := &SafetyHazardObjective{
		SEMatrix:       hsConfigs.SEMatrix,
		AlphaSHPenalty: hsConfigs.AlphaSHPenalty,
		Phases:         hsConfigs.Phases,
		FilePath:       hsConfigs.FilePath,
	}
	return tcObj, nil
}

func (obj *SafetyHazardObjective) Eval(locations map[string]data.Location) float64 {
	result := 0.0

	for _, phases := range obj.Phases {

		for i := 0; i < len(phases); i++ {
			facilityNameI := phases[i]
			facilityI := locations[facilityNameI]
			idxI, err := facilityI.ConvertToIdx()
			if err != nil {
				return 0
			}

			for j := 0; j < len(phases); j++ {
				if i == j {
					continue
				}
				facilityNameJ := phases[j]
				facilityJ := locations[facilityNameJ]
				idxJ, err := facilityJ.ConvertToIdx()
				if err != nil {
					return 0
				}

				result += (-data.Distance2D(facilityI.Coordinate, facilityJ.Coordinate)) / obj.SEMatrix[idxI][idxJ]
			}
		}
	}

	return result
}

func (obj *SafetyHazardObjective) GetAlphaPenalty() float64 {
	return obj.AlphaSHPenalty
}

func ReadSafetyAndEnvDataFromFile(filePath string) ([][]float64, error) {
	dataFile, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	rows, err := dataFile.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	interactionMatrix := make([][]float64, len(rows)-1)

	for idx, row := range rows {
		// skip header
		if idx == 0 {
			continue
		}

		arr := make([]float64, len(rows)-1)
		for i, cell := range row {
			if i == 0 || i == idx {
				continue
			}

			// skip the first column
			arr[i-1], err = strconv.ParseFloat(cell, 64)
			if err != nil {
				return nil, err
			}

		}
		interactionMatrix[idx-1] = arr
	}

	return interactionMatrix, nil
}
