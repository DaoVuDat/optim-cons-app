package objectives

import (
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/data"
	"log"
	"maps"
	"slices"
	"strconv"
)

const ConstructionCostObjectiveType data.ObjectiveType = "Construction Cost Objective"

type ConstructionCostConfigs struct {
	FrequencyMatrix              data.TwoDimensionalMatrix
	DistanceMatrix               data.TwoDimensionalMatrix
	FullRun                      bool
	Delta                        float64
	AlphaConstructionCostPenalty float64
	Phases                       [][]string
	FrequencyFilePath            string
	DistanceFilePath             string
}

type ConstructionCostObjective struct {
	FrequencyMatrix              data.TwoDimensionalMatrix
	DistanceMatrix               data.TwoDimensionalMatrix
	FullRun                      bool
	Delta                        float64
	AlphaConstructionCostPenalty float64
	Phases                       [][]string
	FrequencyFilePath            string
	DistanceFilePath             string
}

func CreateConstructionCostObjectiveFromConfig(ccConfigs ConstructionCostConfigs) (*ConstructionCostObjective, error) {
	ccObj := &ConstructionCostObjective{
		FrequencyMatrix:              ccConfigs.FrequencyMatrix,
		DistanceMatrix:               ccConfigs.DistanceMatrix,
		Delta:                        ccConfigs.Delta,
		AlphaConstructionCostPenalty: ccConfigs.AlphaConstructionCostPenalty,
		Phases:                       ccConfigs.Phases,
		FrequencyFilePath:            ccConfigs.FrequencyFilePath,
		DistanceFilePath:             ccConfigs.DistanceFilePath,
		FullRun:                      ccConfigs.FullRun,
	}
	return ccObj, nil
}

func (obj *ConstructionCostObjective) Eval(locations map[string]data.Location) float64 {
	results := 0.0

	// convert to slices
	slicesLocations := slices.Collect(maps.Values(locations))

	slices.SortStableFunc(slicesLocations, func(i, j data.Location) int {
		idxi, _ := i.ConvertToIdxRegex()
		idxj, _ := j.ConvertToIdxRegex()
		return idxi - idxj
	})

	var maxI, maxJ int
	if obj.FullRun {
		maxI = len(slicesLocations)
		maxJ = len(slicesLocations)
	} else {
		maxI = len(slicesLocations) - 1
		maxJ = len(slicesLocations)
	}

	for i := 0; i < maxI; i++ {
		locI := slicesLocations[i]
		iLocatedAt := locI.IsLocatedAt
		iSymbol := locI.Symbol

		var startJ int
		if obj.FullRun {
			startJ = 0
		} else {
			startJ = i + 1
		}

		for j := startJ; j < maxJ; j++ {
			locJ := slicesLocations[j]
			jLocatedAt := locJ.IsLocatedAt
			jSymbol := locJ.Symbol
			fref, err := obj.FrequencyMatrix.GetCellValueFromNames(iSymbol, jSymbol)
			if err != nil {
				log.Fatal(err)
				return 0
			}

			distance, err := obj.DistanceMatrix.GetCellValueFromNames(iLocatedAt, jLocatedAt)
			if err != nil {
				log.Fatal(err)
				return 0
			}

			results += fref * distance
		}

	}

	return results
}

func (obj *ConstructionCostObjective) GetAlphaPenalty() float64 {
	return obj.AlphaConstructionCostPenalty
}

func ReadMatrixFromFile(filePath string) (data.TwoDimensionalMatrix, error) {

	// load data from file
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return data.TwoDimensionalMatrix{}, err
	}

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		return data.TwoDimensionalMatrix{}, err
	}

	names := make([]string, len(rows)-1)
	for idx, cell := range rows[0] {
		if idx == 0 {
			continue
		}
		names[idx-1] = cell
	}

	freqMatrix := data.CreateTwoDimensionalMatrix(names)

	for idx, row := range rows {
		// skip header
		if idx == 0 {
			continue
		}

		for cellIdx, cell := range row {
			if cellIdx == 0 {
				continue
			}

			val, err := strconv.ParseFloat(cell, 64)
			if err != nil {
				return data.TwoDimensionalMatrix{}, err
			}

			// Set to new matrix
			if err := freqMatrix.SetCellValueFromNames(rows[0][idx], rows[0][cellIdx], val); err != nil {
				return data.TwoDimensionalMatrix{}, err
			}
		}

	}

	return freqMatrix, nil
}
