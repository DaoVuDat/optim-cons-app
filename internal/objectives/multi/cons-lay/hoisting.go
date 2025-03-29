package cons_lay

import (
	"github.com/xuri/excelize/v2"
	"strconv"
)

const (
	HoistingObjectiveType  = "Hoisting Objective"
	HoistingTimeData       = "HoistingTimeData"
	PrefabricatedLocations = "PrefabricatedLocations"
	NumberOfFloors         = "Number of Floors"
	FloorHeight            = "Floor Height"
	CraneLocations         = "Crane Locations"
	ZM                     = "ZM"
	Vuvg                   = "Vuvg"
	Vlvg                   = "Vlvg"
	Vag                    = "Vag"
	Vwg                    = "Vwg"
)

type HoistingConfigs struct {
	PrefabricatedLocations []Location
	NumberOfFloors         int
	HoistingTime           []HoistingTime
	FloorHeight            float64
	CraneLocations         []Crane
	ZM                     float64
	Vuvg                   float64
	Vlvg                   float64
	Vag                    float64
	Vwg                    float64
}

type Crane struct {
	Location
	BuildingName []string
	Radius       float64
}

type HoistingTime struct {
	Coordinate     Coordinate
	HoistingNumber int
	Name           string
	BuildingName   string
}

type HoistingObjective struct {
	PrefabricatedLocations []Location
	NumberOfFloors         int
	HoistingTime           []HoistingTime
	FloorHeight            float64
	CraneLocations         []Crane
	ZM                     float64
	Vuvg                   float64
	Vlvg                   float64
	Vag                    float64
	Vwg                    float64
}

func CreateHoistingObjective() (*HoistingObjective, error) {
	return &HoistingObjective{}, nil
}

func CreateHoistingObjectiveFromConfig(hoistingConfigs HoistingConfigs) (*HoistingObjective, error) {
	hoistingObj := &HoistingObjective{
		PrefabricatedLocations: hoistingConfigs.PrefabricatedLocations,
		NumberOfFloors:         hoistingConfigs.NumberOfFloors,
		HoistingTime:           hoistingConfigs.HoistingTime,
		FloorHeight:            hoistingConfigs.FloorHeight,
		CraneLocations:         hoistingConfigs.CraneLocations,
		ZM:                     hoistingConfigs.ZM,
		Vuvg:                   hoistingConfigs.Vuvg,
		Vlvg:                   hoistingConfigs.Vlvg,
		Vag:                    hoistingConfigs.Vag,
		Vwg:                    hoistingConfigs.Vwg,
	}
	return hoistingObj, nil
}

func (obj *HoistingObjective) Eval() float64 {
	return 0
}

func IsCoverRangeOfCrane(crane Crane, buildings []Location) (bool, float64) {
	invalidAmountTotal := 0.0

	for _, building := range buildings {
		// Top Right
		topRightCoordinate := Coordinate{
			X: building.Coordinate.X + building.Length/2,
			Y: building.Coordinate.Y + building.Width/2,
		}

		topRightAmount := Distance2D(topRightCoordinate, crane.Coordinate) - crane.Radius

		topLeftCoordinate := Coordinate{
			X: building.Coordinate.X - building.Length/2,
			Y: building.Coordinate.Y + building.Width/2,
		}

		topLeftAmount := Distance2D(topLeftCoordinate, crane.Coordinate) - crane.Radius

		bottomLeftCoordinate := Coordinate{
			X: building.Coordinate.X - building.Length/2,
			Y: building.Coordinate.Y - building.Width/2,
		}
		bottomLeftAmount := Distance2D(bottomLeftCoordinate, crane.Coordinate) - crane.Radius

		bottomRightCoordinate := Coordinate{
			X: building.Coordinate.X + building.Length/2,
			Y: building.Coordinate.Y - building.Width/2,
		}
		bottomRightAmount := Distance2D(bottomRightCoordinate, crane.Coordinate) - crane.Radius

		if topRightAmount > 0 || topLeftAmount > 0 || bottomLeftAmount > 0 || bottomRightAmount > 0 {
			invalidAmountTotal += max(0, topLeftAmount) +
				max(0, topRightAmount) +
				max(0, bottomLeftAmount) +
				max(0, bottomRightAmount)

		}
	}

	if invalidAmountTotal > 0 {
		return false, invalidAmountTotal
	} else {
		return true, 0
	}

}

// Readers Utility Functions

func ReadHoistingTimeDataFromFile(filePath string) ([]HoistingTime, error) {
	dataFile, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	rows, err := dataFile.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	hoistingTime := make([]HoistingTime, 0)

	for idx, row := range rows {
		if idx == 0 {
			continue
		}
		var name string
		var buildingName string
		var x float64
		var y float64
		var hoistingNumber int
		for i, cell := range row {
			switch i {
			case 0:
				name = cell
			case 1:
				buildingName = cell
			case 2:
				val, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					return nil, err
				}
				x = val
			case 3:
				val, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					return nil, err
				}
				y = val
			case 4:
				val, err := strconv.ParseInt(cell, 10, 64)
				if err != nil {
					return nil, err
				}
				hoistingNumber = int(val)
			}

		}
		hoistingTime = append(hoistingTime, HoistingTime{
			Coordinate: Coordinate{
				X: x,
				Y: y,
			},
			HoistingNumber: hoistingNumber,
			Name:           name,
			BuildingName:   buildingName,
		})
	}
	return hoistingTime, nil
}
