package objectives

import (
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/data"
	"math"
	"strconv"
)

const HoistingObjectiveType data.ObjectiveType = "Hoisting Objective"

type HoistingConfigs struct {
	NumberOfFloors       int
	HoistingTime         map[string][]HoistingTime
	FloorHeight          float64
	CraneLocations       []data.Crane
	ZM                   float64
	Vuvg                 float64
	Vlvg                 float64
	Vag                  float64
	Vwg                  float64
	AlphaHoistingPenalty float64
	AlphaHoisting        float64
	BetaHoisting         float64
	NHoisting            float64
	Phases               [][]string
	HoistingTimeWithInfo []HoistingTimeWithInfo
}

type HoistingTime struct {
	Coordinate     data.Coordinate
	HoistingNumber int
	Name           string
	FacilitySymbol string
}

type HoistingObjective struct {
	NumberOfFloors       int
	HoistingTime         map[string][]HoistingTime
	FloorHeight          float64
	CraneLocations       []data.Crane
	ZM                   float64
	Vuvg                 float64
	Vlvg                 float64
	Vag                  float64
	Vwg                  float64
	AlphaHoistingPenalty float64
	AlphaHoisting        float64
	BetaHoisting         float64
	NHoisting            float64
	Phases               [][]string
	HoistingTimeWithInfo []HoistingTimeWithInfo
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
		Phases:               hoistingConfigs.Phases,
		HoistingTimeWithInfo: hoistingConfigs.HoistingTimeWithInfo,
	}
	return hoistingObj, nil
}

func (obj *HoistingObjective) Eval(locations map[string]data.Location) float64 {

	result := 0.0

	cranes := make([]data.Crane, len(obj.CraneLocations))
	for i, location := range obj.CraneLocations {
		if loc, ok := locations[location.CraneSymbol]; ok {
			cranes[i].CraneSymbol = location.Symbol
			cranes[i].Radius = location.Radius
			cranes[i].BuildingName = location.BuildingName
			cranes[i].Coordinate.X = loc.Coordinate.X
			cranes[i].Coordinate.Y = loc.Coordinate.Y
			cranes[i].IsFixed = loc.IsFixed
			cranes[i].Length = loc.Length
			cranes[i].Width = loc.Width
			cranes[i].Rotation = loc.Rotation
			cranes[i].Symbol = loc.Symbol
			cranes[i].Name = loc.Name
		}
	}

	// calculate Hdjg = distance(crane, prefabricated)
	for _, crane := range cranes {
		TB := 0.0
		HDjg := make(map[string]float64, len(crane.BuildingName))
		for _, prefabricatedName := range crane.BuildingName {
			HDjg[prefabricatedName] = data.Distance2D(crane.Coordinate, locations[prefabricatedName].Coordinate)
		}

		hoistingTime := obj.HoistingTime[crane.Symbol]
		for _, hoisting := range hoistingTime {
			// calculate distance between hoisting and prefabricated
			HDkg := data.Distance2D(hoisting.Coordinate, crane.Coordinate)

			// calculate distance between demand and prefabricated
			Djk := data.Distance2D(locations[hoisting.FacilitySymbol].Coordinate, hoisting.Coordinate)

			Tag := 2 * (math.Abs(HDjg[hoisting.FacilitySymbol]-HDkg) / obj.Vag)
			Twg := 2 * (1 / obj.Vwg) * math.Acos((HDjg[hoisting.FacilitySymbol]*HDjg[hoisting.FacilitySymbol]+HDkg*HDkg-Djk*Djk)/
				(2*HDjg[hoisting.FacilitySymbol]*HDkg))

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

type HoistingTimeWithInfo struct {
	CraneSymbol  string
	FilePath     string
	Radius       float64
	BuildingName []string
}

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
		var facilitySymbol string
		var x float64
		var y float64
		var hoistingNumber int
		for i, cell := range row {
			switch i {
			case 0:
				name = cell
			case 1:
				facilitySymbol = cell
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
			Coordinate: data.Coordinate{
				X: x,
				Y: y,
			},
			HoistingNumber: hoistingNumber,
			Name:           name,
			FacilitySymbol: facilitySymbol,
		})
	}
	return hoistingTime, nil
}
