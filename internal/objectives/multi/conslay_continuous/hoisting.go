package conslay_continuous

import (
	"github.com/xuri/excelize/v2"
	"math"
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
	AlphaHoistingPenalty   = "AlphaHoistingPenalty"
	AlphaHoisting          = "AlphaHoisting"
	BetaHoisting           = "BetaHoisting"
	NHoisting              = "NHoisting"
)

type HoistingConfigs struct {
	NumberOfFloors       int
	HoistingTime         map[string][]HoistingTime
	FloorHeight          float64
	CraneLocations       []Crane
	ZM                   float64
	Vuvg                 float64
	Vlvg                 float64
	Vag                  float64
	Vwg                  float64
	AlphaHoistingPenalty float64
	AlphaHoisting        float64
	BetaHoisting         float64
	NHoisting            float64
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
	NumberOfFloors       int
	HoistingTime         map[string][]HoistingTime
	FloorHeight          float64
	CraneLocations       []Crane
	ZM                   float64
	Vuvg                 float64
	Vlvg                 float64
	Vag                  float64
	Vwg                  float64
	AlphaHoistingPenalty float64
	AlphaHoisting        float64
	BetaHoisting         float64
	NHoisting            float64
}

func CreateHoistingObjective() (*HoistingObjective, error) {
	return &HoistingObjective{}, nil
}

func CreateHoistingObjectiveFromConfig(hoistingConfigs HoistingConfigs) (*HoistingObjective, error) {
	hoistingObj := &HoistingObjective{
		NumberOfFloors:       hoistingConfigs.NumberOfFloors,
		HoistingTime:         hoistingConfigs.HoistingTime,
		FloorHeight:          hoistingConfigs.FloorHeight,
		CraneLocations:       hoistingConfigs.CraneLocations,
		ZM:                   hoistingConfigs.ZM,
		Vuvg:                 hoistingConfigs.Vuvg,
		Vlvg:                 hoistingConfigs.Vlvg,
		Vag:                  hoistingConfigs.Vag,
		Vwg:                  hoistingConfigs.Vwg,
		AlphaHoistingPenalty: hoistingConfigs.AlphaHoistingPenalty,
		AlphaHoisting:        hoistingConfigs.AlphaHoisting,
		BetaHoisting:         hoistingConfigs.BetaHoisting,
		NHoisting:            hoistingConfigs.NHoisting,
	}
	return hoistingObj, nil
}

func (obj *HoistingObjective) Eval(locations map[string]Location) float64 {

	result := 0.0

	// calculate Hdjg = distance(crane, prefabricated)
	for _, crane := range obj.CraneLocations {
		TB := 0.0
		HDjg := make(map[string]float64, len(crane.BuildingName))
		for _, prefabricatedName := range crane.BuildingName {
			HDjg[prefabricatedName] = Distance2D(crane.Coordinate, locations[prefabricatedName].Coordinate)
		}

		hoistingTime := obj.HoistingTime[crane.Name]
		for _, hoisting := range hoistingTime {
			// calculate distance between hoisting and prefabricated
			HDkg := Distance2D(hoisting.Coordinate, crane.Coordinate)

			// calculate distance between demand and prefabricated
			Djk := Distance2D(locations[hoisting.BuildingName].Coordinate, hoisting.Coordinate)

			Tag := 2 * (math.Abs(HDjg[hoisting.BuildingName]-HDkg) / obj.Vag)
			Twg := 2 * (1 / obj.Vwg) * math.Acos((HDjg[hoisting.BuildingName]*HDjg[hoisting.BuildingName]+HDkg*HDkg-Djk*Djk)/
				(2*HDjg[hoisting.BuildingName]*HDkg))

			Thg := max(Tag, Twg) + obj.AlphaHoisting*min(Tag, Twg)

			for i := 0; i < obj.NumberOfFloors; i++ {
				ZOj := float64(i) * obj.FloorHeight

				Tvg := (1/obj.Vuvg + 1/obj.Vlvg) * math.Abs(ZOj-obj.ZM)
				Tg := max(Thg, Tvg) + obj.BetaHoisting*min(Thg, Tvg)
				TB = TB + float64(hoisting.HoistingNumber)*Tg
			}
		}
		result += TB
	}

	return result
}

func (obj *HoistingObjective) GetAlphaPenalty() float64 {
	return obj.AlphaHoistingPenalty
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
