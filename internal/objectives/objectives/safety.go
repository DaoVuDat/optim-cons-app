package objectives

import (
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/data"
	"log"
	"strconv"
)

const SafetyObjectiveType data.ObjectiveType = "Safety Objective"

type SafetyConfigs struct {
	SafetyProximity    data.TwoDimensionalMatrix
	AlphaSafetyPenalty float64
	Phases             [][]string
	FilePath           string
}

type SafetyObjective struct {
	SafetyProximity    data.TwoDimensionalMatrix
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

	calculatedMap := make(map[string]struct{})

	for _, phases := range obj.Phases {

		for i := 0; i < len(phases)-1; i++ {
			facilityNameI := phases[i]
			facilityI := locations[facilityNameI]
			idxI, err := obj.SafetyProximity.GetIdxFromName(facilityNameI)
			if err != nil {
				log.Fatal(err)
				return 0
			}

			for j := i + 1; j < len(phases); j++ {
				facilityNameJ := phases[j]
				facilityJ := locations[facilityNameJ]
				idxJ, err := obj.SafetyProximity.GetIdxFromName(facilityNameJ)
				if err != nil {
					log.Fatal(err)
					return 0
				}

				if _, ok := calculatedMap[facilityNameI+facilityNameJ]; ok {
					continue
				}

				result += obj.SafetyProximity.Matrix[idxI][idxJ] * data.Distance2D(facilityI.Coordinate, facilityJ.Coordinate)
				calculatedMap[facilityNameI+facilityNameJ] = struct{}{}
			}
		}
	}

	return result
}

func (obj *SafetyObjective) GetAlphaPenalty() float64 {
	return obj.AlphaSafetyPenalty
}

func ReadSafetyProximityDataFromFile(filePath string) (data.TwoDimensionalMatrix, error) {
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
		facilitiesName[idx-1] = cell
	}

	hazardInteractionNew := data.CreateTwoDimensionalMatrix(facilitiesName)

	for idx, row := range rows {
		// skip header
		if idx == 0 {
			continue
		}

		for i := idx + 1; i < len(row); i++ {
			cell := row[i]

			val, err := strconv.ParseFloat(cell, 64)
			if err != nil {
				return data.TwoDimensionalMatrix{}, err
			}

			if err := hazardInteractionNew.SetCellValueFromNames(rows[0][idx], rows[0][i], val); err != nil {
				return data.TwoDimensionalMatrix{}, err
			}
		}

	}

	return hazardInteractionNew, nil
}
