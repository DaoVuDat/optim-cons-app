package objectives

import (
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/data"
	"strconv"
)

const ConstructionCostObjectiveType data.ObjectiveType = "Construction Cost Objective"

type ConstructionCostConfigs struct {
	ConstructionCostMatrix          data.TwoDimensionalMatrix
	AdjacencyMatrix                 data.TwoDimensionalMatrix
	InteractionCostFacilitiesMatrix data.TwoDimensionalMatrix
	Delta                           float64
	AlphaRiskPenalty                float64
	Phases                          [][]string
	FilePath                        string
}

type ConstructionCostObjective struct {
	ConstructionCostMatrix          data.TwoDimensionalMatrix
	AdjacencyMatrix                 data.TwoDimensionalMatrix
	InteractionCostFacilitiesMatrix data.TwoDimensionalMatrix
	Delta                           float64
	AlphaRiskPenalty                float64
	Phases                          [][]string
	FilePath                        string
}

func CreateConstructionCostObjectiveFromConfig(ccConfigs ConstructionCostConfigs) (*ConstructionCostObjective, error) {
	ccObj := &ConstructionCostObjective{
		ConstructionCostMatrix:          ccConfigs.ConstructionCostMatrix,
		AdjacencyMatrix:                 ccConfigs.AdjacencyMatrix,
		InteractionCostFacilitiesMatrix: ccConfigs.InteractionCostFacilitiesMatrix,
		Delta:                           ccConfigs.Delta,
		AlphaRiskPenalty:                ccConfigs.AlphaRiskPenalty,
		Phases:                          ccConfigs.Phases,
		FilePath:                        ccConfigs.FilePath,
	}
	return ccObj, nil
}

func (obj *ConstructionCostObjective) Eval(locations map[string]data.Location) float64 {
	//mapFacility := make(map[string]struct {
	//	Count int
	//	Value float64
	//})

	results := 0.0

	//for _, phases := range obj.Phases {
	//	hij := util.CopySliceOfSlice(obj.HazardInteractionMatrix.Matrix)
	//
	//	for i := 0; i < len(phases); i++ {
	//		facilityNameI := phases[i]
	//		facilityI := locations[facilityNameI]
	//		idxI, err := obj.HazardInteractionMatrix.GetIdxFromName(facilityNameI)
	//		if err != nil {
	//			log.Fatal(err)
	//			return 0
	//		}
	//
	//		hio := hij[idxI][idxI]
	//		for j := 0; j < len(phases); j++ {
	//			facilityNameJ := phases[j]
	//			facilityJ := locations[facilityNameJ]
	//			idxJ, err := obj.HazardInteractionMatrix.GetIdxFromName(facilityNameJ)
	//			if err != nil {
	//				log.Fatal(err)
	//				return 0
	//			}
	//
	//			computed := hio
	//			if i != j {
	//				computed = hio - obj.Delta*data.Distance2D(facilityI.Coordinate, facilityJ.Coordinate)
	//			}
	//
	//			hijComputed := max(0, computed)
	//
	//			k := fmt.Sprintf("%s-%s", facilityNameI, facilityNameJ)
	//
	//			if v, ok := mapFacility[k]; ok {
	//				v.Count++
	//				mapFacility[k] = v
	//			} else {
	//				mapFacility[k] = struct {
	//					Count int
	//					Value float64
	//				}{Count: 1, Value: hijComputed * hijComputed}
	//			}
	//
	//			hij[idxI][idxJ] = hijComputed
	//		}
	//	}
	//
	//	phaseResult := 0.0
	//	for i := 0; i < len(phases); i++ {
	//		facilityNameI := phases[i]
	//		idxI, err := obj.HazardInteractionMatrix.GetIdxFromName(facilityNameI)
	//		if err != nil {
	//			log.Fatal(err)
	//			return 0
	//		}
	//
	//		for j := 0; j < len(phases); j++ {
	//			facilityNameJ := phases[j]
	//			idxJ, err := obj.HazardInteractionMatrix.GetIdxFromName(facilityNameJ)
	//			if err != nil {
	//				log.Fatal(err)
	//				return 0
	//			}
	//
	//			phaseResult += hij[idxI][idxJ] * hij[idxI][idxJ]
	//		}
	//	}
	//
	//	results += phaseResult
	//}
	//
	//for _, v := range mapFacility {
	//	if v.Count > 1 {
	//		results -= (float64(v.Count) - 1) * v.Value
	//	}
	//}

	return results
}

func (obj *ConstructionCostObjective) GetAlphaPenalty() float64 {
	return obj.AlphaRiskPenalty
}

func ReadConstructionCostMatrixDataFromFile(filePath string) (data.TwoDimensionalMatrix, error) {
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

		val, err := strconv.ParseFloat(row[idx], 64)
		if err != nil {
			return data.TwoDimensionalMatrix{}, err
		}

		// Set to new matrix
		if err := hazardInteractionNew.SetCellValueFromNames(rows[0][idx], rows[0][idx], val); err != nil {
			return data.TwoDimensionalMatrix{}, err
		}
	}
	return hazardInteractionNew, nil
}
