package objectives

import (
	"fmt"
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
	Delta                        float64
	AlphaConstructionCostPenalty float64
	Phases                       [][]string
	FilePath                     string
}

type ConstructionCostObjective struct {
	FrequencyMatrix              data.TwoDimensionalMatrix
	DistanceMatrix               data.TwoDimensionalMatrix
	Delta                        float64
	AlphaConstructionCostPenalty float64
	Phases                       [][]string
	FilePath                     string
}

func CreateConstructionCostObjectiveFromConfig(ccConfigs ConstructionCostConfigs) (*ConstructionCostObjective, error) {
	ccObj := &ConstructionCostObjective{
		FrequencyMatrix:              ccConfigs.FrequencyMatrix,
		DistanceMatrix:               ccConfigs.DistanceMatrix,
		Delta:                        ccConfigs.Delta,
		AlphaConstructionCostPenalty: ccConfigs.AlphaConstructionCostPenalty,
		Phases:                       ccConfigs.Phases,
		FilePath:                     ccConfigs.FilePath,
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

	for i := 0; i < len(slicesLocations); i++ {
		locI := slicesLocations[i]
		iLocatedAt := locI.IsLocatedAt
		iSymbol := locI.Symbol
		for j := 0; j < len(slicesLocations); j++ {
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
			fmt.Println("==")
			fmt.Printf("Freq(i = %s, j = %s) = %f\n", iSymbol, jSymbol, fref)
			fmt.Printf("Dist(i = %s, j = %s) = %f\n", iLocatedAt, jLocatedAt, distance)
			fmt.Println("\t\t", fref*distance)
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
