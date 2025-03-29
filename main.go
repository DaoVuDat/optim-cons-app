/*
Copyright © 2025 Dao Vu Dat dat.daovu@gmail.com
*/
package main

import (
	"fmt"
	conslay "golang-moaha-construction/internal/objectives/multi/cons-lay"
	"log"
)

func main() {
	// Create cons-lay problem and add objectives
	fmt.Println("=== Construction Layout ===")
	consLayoutConfigs := conslay.ConsLayConfigs{
		ConsLayoutLength: 120,
		ConsLayoutWidth:  95,
	}

	// LOAD LOCATIONS
	locations, fixedLocations, nonFixedLocations, err := conslay.ReadLocationsFromFile("./data/conslay/locations.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\tLocations")
	for name := range locations {
		fmt.Printf("Name = %s, Symbol = %s, X = %f, Y = %f, Length = %f, Width = %f, fixed = %t \n",
			locations[name].Name,
			locations[name].Symbol,
			locations[name].Coordinate.X,
			locations[name].Coordinate.Y,
			locations[name].Length,
			locations[name].Width,
			locations[name].IsFixed,
		)
	}
	consLayoutConfigs.Locations = locations
	consLayoutConfigs.NonFixedLocations = nonFixedLocations
	consLayoutConfigs.FixedLocations = fixedLocations

	fmt.Println("#Locations", len(consLayoutConfigs.Locations))
	fmt.Println("#FixedLocations", len(consLayoutConfigs.FixedLocations))
	fmt.Println("#NonFixedLocations", len(consLayoutConfigs.NonFixedLocations))

	// LOAD PHASES
	phases, err := conslay.ReadPhasesFromFile("./data/conslay/phaseBuilding.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\tPhases")
	for i := range phases {
		fmt.Println(phases[i])
	}
	consLayoutConfigs.Phases = phases

	consLayObj, err := conslay.CreateConsLayFromConfig(consLayoutConfigs)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: objectives - select objectives and show configs relevant to those
	fmt.Println("=== Hoisting Objective ===")
	hoistingTime, err := conslay.ReadHoistingTimeDataFromFile("./data/conslay/f1_hoisting_time_data.xlsx")

	if err != nil {
		log.Fatal(err)
	}

	// Hoisting Objective Configs
	hoistingConfigs := conslay.HoistingConfigs{
		NumberOfFloors: 10,
		HoistingTime:   hoistingTime,
		FloorHeight:    3.2,
		ZM:             2,
		Vuvg:           37.5,
		Vlvg:           37.5 / 2,
		Vag:            50,
		Vwg:            0.5,
	}

	hoistingObj, err := conslay.CreateHoistingObjectiveFromConfig(hoistingConfigs)
	if err != nil {
		log.Fatal(err)
	}

	// select TF that is crane (fixed locations only) - after Selection
	// simulate selected crane
	type SelectedCrane struct {
		Name          string // for reference from FixedLocations
		BuildingNames []string
		Radius        float64
	}

	selectedCrane := []SelectedCrane{
		{
			Name:          "TF14",
			BuildingNames: []string{"TF8", "TF9", "TF10"},
			Radius:        40,
		},
	}

	craneLocations := make([]conslay.Crane, 0)
	for _, loc := range selectedCrane {
		if craneLoc, ok := consLayObj.Locations[loc.Name]; ok {
			craneLocations = append(craneLocations, conslay.Crane{
				Location:     craneLoc,
				BuildingName: loc.BuildingNames,
				Radius:       loc.Radius,
			})
		}
	}

	hoistingObj.CraneLocations = craneLocations

	// select TF that is prefabricated - in Evaluation
	selectedPref := []string{
		"TF8", "TF9", "TF10",
	}
	prefLocs := make([]conslay.Location, 0)
	for _, loc := range selectedPref {
		if prefLoc, ok := consLayObj.Locations[loc]; ok {
			prefLocs = append(prefLocs, conslay.Location{
				Coordinate: prefLoc.Coordinate,
				Length:     prefLoc.Length,
				Width:      prefLoc.Width,
				IsFixed:    true,
				Name:       prefLoc.Name,
				Symbol:     loc,
			})
		}
	}
	hoistingObj.PrefabricatedLocations = prefLocs

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

	// Add objectives to cons-lay problem
	err = consLayObj.AddObjective(conslay.HoistingObjectiveType, hoistingObj)
	if err != nil {
		log.Fatal(err)
	}

	return
}
