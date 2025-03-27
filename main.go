/*
Copyright © 2025 Dao Vu Dat dat.daovu@gmail.com
*/
package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/data"
	conslay "golang-moaha-construction/internal/objectives/multi/cons-lay"
	"log"
	"strconv"
)

func main() {

	// Generic Configs
	consLayoutConfigs := []data.Config{
		{
			Name:  conslay.ConsLayoutLength,
			Value: "120",
		},
		{
			Name:  conslay.ConsLayoutWidth,
			Value: "95",
		},
		{
			Name:  conslay.DynamicLocations,
			Value: "./data/conslay/dynamic_locations.xlsx", // path to the file
		},
		{
			Name:  conslay.StaticLocations,
			Value: "./data/conslay/fixed_locations.xlsx", // path to the file
		},
		{
			Name:  conslay.Phases,
			Value: "./data/conslay/phaseBuilding.xlsx", // path to the file
		},
	}

	// TODO: objectives - select objectives and show configs relevant to those

	dataFile, err := excelize.OpenFile("./data/conslay/crane_locations.xlsx")
	if err != nil {
		log.Fatal(err)
		return
	}

	rows, err := dataFile.GetRows("Sheet1")
	if err != nil {
		log.Fatal(err)
		return
	}

	craneLocations := make([]conslay.Crane, 0)

	for idx, row := range rows {
		if idx == 0 {
			continue
		}
		var name string
		var length float64
		var width float64
		var x float64
		var y float64
		var radius float64
		for i, cell := range row {
			switch i {
			case 0:
				name = cell
			case 1:
				val, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					log.Fatal(err)
					return
				}
				length = val
			case 2:
				val, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					log.Fatal(err)
					return
				}
				width = val
			case 3:
				val, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					log.Fatal(err)
					return
				}
				x = val
			case 4:
				val, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					log.Fatal(err)
					return
				}
				y = val
			case 5:
				val, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					log.Fatal(err)
					return
				}
				radius = val
			}

		}
		craneLocations = append(craneLocations, conslay.Crane{
			Location: conslay.Location{
				Coordinate: conslay.Coordinate{
					X: x,
					Y: y,
				},
				Length:  length,
				Width:   width,
				IsFixed: true,
				Name:    name,
			},
			Radius: radius,
		})
	}

	dataFile, err = excelize.OpenFile("./data/conslay/prefabricated_locations.xlsx")
	if err != nil {
		log.Fatal(err)
		return
	}

	rows, err = dataFile.GetRows("Sheet1")
	if err != nil {
		log.Fatal(err)
		return
	}

	prefLocs := make([]conslay.Location, 0)

	for idx, row := range rows {
		if idx == 0 {
			continue
		}
		var name string
		var length float64
		var width float64
		for i, cell := range row {
			switch i {
			case 0:
				name = cell
			case 1:
				val, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					log.Fatal(err)
					return
				}
				length = val
			case 2:
				val, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					log.Fatal(err)
					return
				}
				width = val
			}

		}
		prefLocs = append(prefLocs, conslay.Location{
			Length:  length,
			Width:   width,
			IsFixed: false,
			Name:    name,
		})
	}

	dataFile, err = excelize.OpenFile("./data/conslay/f1_hoisting_time_data.xlsx")
	if err != nil {
		log.Fatal(err)
		return
	}

	rows, err = dataFile.GetRows("Sheet1")
	if err != nil {
		log.Fatal(err)
		return
	}

	hoistingTime := make([]conslay.HoistingTime, 0)

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
					log.Fatal(err)
					return
				}
				x = val
			case 3:
				val, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					log.Fatal(err)
					return
				}
				y = val
			case 4:
				val, err := strconv.ParseInt(cell, 10, 64)
				if err != nil {
					log.Fatal(err)
					return
				}
				hoistingNumber = int(val)
			}

		}
		hoistingTime = append(hoistingTime, conslay.HoistingTime{
			Coordinate: conslay.Coordinate{
				X: x,
				Y: y,
			},
			HoistingNumber: hoistingNumber,
			Name:           name,
			BuildingName:   buildingName,
		})
	}

	hoistingObj := conslay.HoistingObjective{
		PrefabricatedLocations: prefLocs,
		NumberOfFloors:         10,
		HoistingTime:           hoistingTime,
		FloorHeight:            3.2,
		CraneLocations:         craneLocations,
		ZM:                     2,
		Vuvg:                   37.5,
		Vlvg:                   37.5 / 2,
		Vag:                    50,
		Vwg:                    0.5,
	}
	fmt.Println("\tHoisting Time")
	for i := range hoistingObj.HoistingTime {
		fmt.Printf("%d: Name = %s, Building Name = %s, X = %f, Y = %f, Hoisting Number = %d \n",
			i+1,
			hoistingObj.HoistingTime[i].Name,
			hoistingObj.HoistingTime[i].BuildingName,
			hoistingObj.HoistingTime[i].Coordinate.X,
			hoistingObj.HoistingTime[i].Coordinate.Y,
			hoistingObj.HoistingTime[i].HoistingNumber,
		)
	}

	fmt.Println("\tPrefabricated Locations")
	for i := range hoistingObj.PrefabricatedLocations {
		fmt.Printf("%d: Name = %s, L = %f, W = %f \n",
			i+1,
			hoistingObj.PrefabricatedLocations[i].Name,
			hoistingObj.PrefabricatedLocations[i].Length,
			hoistingObj.PrefabricatedLocations[i].Width,
		)
	}

	fmt.Println("\tCrane Locations")
	for i := range hoistingObj.CraneLocations {
		fmt.Printf("%d: Name = %s, L = %f, W = %f, x = %f, y = %f, radius = %f \n",
			i+1,
			hoistingObj.CraneLocations[i].Location.Name,
			hoistingObj.CraneLocations[i].Location.Length,
			hoistingObj.CraneLocations[i].Location.Width,
			hoistingObj.CraneLocations[i].Location.Coordinate.X,
			hoistingObj.CraneLocations[i].Location.Coordinate.Y,
			hoistingObj.CraneLocations[i].Radius,
		)
	}

	// Create cons-lay problem and add objectives
	consLayObj, _ := conslay.Create()
	err = consLayObj.LoadData(consLayoutConfigs)
	if err != nil {
		log.Fatal(err)
		return
	}

	consLay := consLayObj.(*conslay.ConsLay)

	fmt.Println("\tDynamic Data")
	for i := range consLay.DynamicLocations {

		fmt.Printf("%d: Name = %s, L = %f, W = %f, fixed = %t \n",
			i+1,
			consLay.DynamicLocations[i].Name,
			consLay.DynamicLocations[i].Length,
			consLay.DynamicLocations[i].Width,
			consLay.DynamicLocations[i].IsFixed,
		)
	}

	fmt.Println("\tStatic Data")
	for i := range consLay.StaticLocations {

		fmt.Printf("%d: Name = %s, L = %f, W = %f, x = %f, y = %f, fixed = %t \n",
			i+1,
			consLay.StaticLocations[i].Name,
			consLay.StaticLocations[i].Length,
			consLay.StaticLocations[i].Width,
			consLay.StaticLocations[i].Coordinate.X,
			consLay.StaticLocations[i].Coordinate.Y,
			consLay.StaticLocations[i].IsFixed,
		)
	}

	fmt.Println("\tPhases")
	for i := range consLay.Phases {
		fmt.Println(consLay.Phases[i])
	}

	return
}
