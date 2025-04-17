package objectives

import (
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/data"
	"strconv"
)

const TransportCostObjectiveType data.ObjectiveType = "Transport Cost Objective"

type TransportCostConfigs struct {
	InteractionMatrix [][]float64
	AlphaTCPenalty    float64
	Phases            [][]string
	FilePath          string
}

type TransportCostObjective struct {
	InteractionMatrix [][]float64
	AlphaTCPenalty    float64
	Phases            [][]string
	FilePath          string
}

func CreateTransportCostObjectiveFromConfig(transportCostConfigs TransportCostConfigs) (*TransportCostObjective, error) {
	tcObj := &TransportCostObjective{
		InteractionMatrix: transportCostConfigs.InteractionMatrix,
		AlphaTCPenalty:    transportCostConfigs.AlphaTCPenalty,
		Phases:            transportCostConfigs.Phases,
		FilePath:          transportCostConfigs.FilePath,
	}
	return tcObj, nil
}

func (obj *TransportCostObjective) Eval(locations map[string]data.Location) float64 {
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

				result += obj.InteractionMatrix[idxI][idxJ] * data.Distance2D(facilityI.Coordinate, facilityJ.Coordinate)
			}
		}
	}

	return result
}

func (obj *TransportCostObjective) GetAlphaPenalty() float64 {
	return obj.AlphaTCPenalty
}

func ReadInteractionTransportCostDataFromFile(filePath string) ([][]float64, error) {
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
