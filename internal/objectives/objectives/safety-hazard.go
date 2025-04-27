package objectives

import (
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/data"
	"strconv"
	"strings"
)

const SafetyHazardObjectiveType data.ObjectiveType = "Safety Hazard Objective"

type SafetyHazardConfigs struct {
	SEMatrix       data.TwoDimensionalMatrix
	AlphaSHPenalty float64
	Phases         [][]string
	FilePath       string
}

type SafetyHazardObjective struct {
	SEMatrix       data.TwoDimensionalMatrix
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

	calculatedMap := make(map[string]struct{})

	for _, phases := range obj.Phases {

		for i := 0; i < len(phases); i++ {
			facilityNameI := phases[i]
			facilityI := locations[facilityNameI]
			idxI, err := obj.SEMatrix.GetIdxFromName(facilityNameI)
			if err != nil {
				return 0
			}

			for j := 0; j < len(phases); j++ {
				if i == j {
					continue
				}
				facilityNameJ := phases[j]
				facilityJ := locations[facilityNameJ]
				idxJ, err := obj.SEMatrix.GetIdxFromName(facilityNameJ)
				if err != nil {
					return 0
				}

				if _, ok := calculatedMap[facilityNameI+facilityNameJ]; ok {
					continue
				}

				if obj.SEMatrix.Matrix[idxI][idxJ] == 0 {
					continue
				}

				result += (-data.Distance2D(facilityI.Coordinate, facilityJ.Coordinate)) / obj.SEMatrix.Matrix[idxI][idxJ]
				calculatedMap[facilityNameI+facilityNameJ] = struct{}{}
			}
		}
	}

	return result
}

func (obj *SafetyHazardObjective) GetAlphaPenalty() float64 {
	return obj.AlphaSHPenalty
}

func ReadSafetyAndEnvDataFromFile(filePath string) (data.TwoDimensionalMatrix, error) {
	dataFile, err := excelize.OpenFile(filePath)
	if err != nil {
		return data.TwoDimensionalMatrix{}, err
	}

	rows, err := dataFile.GetRows("Sheet1")
	if err != nil {
		return data.TwoDimensionalMatrix{}, err
	}

	facilitiesName := make([]string, len(rows)-1)

	for idx, cell := range rows[0] {
		// skip the first column
		if idx == 0 {
			continue
		}
		facilitiesName[idx-1] = strings.ToUpper(cell)
	}

	interactionMatrix := data.CreateTwoDimensionalMatrix(facilitiesName)

	for idx, row := range rows {
		// skip header
		if idx == 0 {
			continue
		}

		for i, cell := range row {
			if i == 0 || i == idx {
				continue
			}

			// skip the first column
			val, err := strconv.ParseFloat(cell, 64)
			if err != nil {
				return data.TwoDimensionalMatrix{}, err
			}

			// add to new matrix
			if err := interactionMatrix.SetCellValueFromNames(rows[0][idx], rows[0][i], val); err != nil {
				return data.TwoDimensionalMatrix{}, err
			}

		}
	}

	return interactionMatrix, nil
}
