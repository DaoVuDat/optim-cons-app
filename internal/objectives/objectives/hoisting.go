package objectives

import (
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/data"
	"log"
	"math"
	"strconv"
	"strings"
)

const HoistingObjectiveType data.ObjectiveType = "Hoisting Objective"

type HoistingConfigs struct {
	HoistingTime         map[string][]HoistingTime
	Buildings            map[string]Building
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

type Building struct {
	NumberOfFloors int
	FloorHeight    float64
}

type HoistingTime struct {
	Coordinate     data.Coordinate
	HoistingNumber int
	Name           string
	FacilitySymbol string
}

type HoistingObjective struct {
	HoistingTime         map[string][]HoistingTime
	Buildings            map[string]Building
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
		Buildings:            hoistingConfigs.Buildings,
		HoistingTime:         hoistingConfigs.HoistingTime,
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
	buildings := make([]string, len(obj.CraneLocations))
	for i, location := range obj.CraneLocations {
		parts := strings.Split(location.CraneSymbol, "-")
		// extract crane symbol into crane symbol and building name
		// CraneName-ForBuildingName
		craneSymbol := parts[0]
		if loc, ok := locations[craneSymbol]; ok {
			buildings[i] = parts[1]
			cranes[i].CraneSymbol = location.CraneSymbol
			cranes[i].Coordinate.X = loc.Coordinate.X
			cranes[i].Coordinate.Y = loc.Coordinate.Y
			cranes[i].IsFixed = loc.IsFixed
			cranes[i].Length = loc.Length
			cranes[i].Width = loc.Width
			cranes[i].Rotation = loc.Rotation
			cranes[i].Symbol = loc.Symbol
			cranes[i].Name = loc.Name
		} else {
			log.Fatal("[HoistingTime - Eval()] Crane location not found in locations map")
		}
	}

	// calculate Hdjg = distance(crane, prefabricated)
	for i, crane := range cranes {
		TB := 0.0

		// get number of floors and floor height
		numberOfFloors := obj.Buildings[buildings[i]].NumberOfFloors
		floorHeight := obj.Buildings[buildings[i]].FloorHeight

		hoistingTime := obj.HoistingTime[crane.CraneSymbol]
		for _, hoisting := range hoistingTime {
			// calculate distance between hoisting and prefabricated
			HDkg := data.Distance2D(hoisting.Coordinate, crane.Coordinate)
			// calculate distance between demand and prefabricated
			Djk := data.Distance2D(locations[hoisting.FacilitySymbol].Coordinate, hoisting.Coordinate)

			HDjg := data.Distance2D(crane.Coordinate, locations[hoisting.FacilitySymbol].Coordinate)
			Tag := 2 * (math.Abs(HDjg-HDkg) / obj.Vag)
			Twg := 2 * (1 / obj.Vwg) * math.Acos((HDjg*HDjg+HDkg*HDkg-Djk*Djk)/
				(2*HDjg*HDkg))

			Thg := max(Tag, Twg) + obj.AlphaHoisting*min(Tag, Twg)

			for i := 0; i < numberOfFloors; i++ {
				ZOj := float64(i) * floorHeight

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
				facilitySymbol = strings.ToUpper(cell)
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
