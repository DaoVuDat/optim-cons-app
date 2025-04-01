/*
Copyright © 2025 Dao Vu Dat dat.daovu@gmail.com
*/
package main

import (
	"fmt"
	"golang-moaha-construction/internal/algorithms/moaha"
	conslay "golang-moaha-construction/internal/objectives/multi/conslay_continuous"
	"log"
	"slices"
	"strings"
)

func main() {
	// Create conslay_continuous problem and add objectives
	//fmt.Println("=== Construction Layout ===")
	consLayoutConfigs := conslay.ConsLayConfigs{
		ConsLayoutLength: 120,
		ConsLayoutWidth:  95,
	}

	// LOAD LOCATIONS
	locations, fixedLocations, nonFixedLocations, err := conslay.ReadLocationsFromFile("./data/conslay/locations.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("\tLocations")
	//for name := range locations {
	//	fmt.Printf("Name = %s, Symbol = %s, X = %f, Y = %f, Length = %f, Width = %f, fixed = %t \n",
	//		locations[name].Name,
	//		locations[name].Symbol,
	//		locations[name].Coordinate.X,
	//		locations[name].Coordinate.Y,
	//		locations[name].Length,
	//		locations[name].Width,
	//		locations[name].IsFixed,
	//	)
	//}
	consLayoutConfigs.Locations = locations
	consLayoutConfigs.NonFixedLocations = nonFixedLocations
	consLayoutConfigs.FixedLocations = fixedLocations

	//fmt.Println("#Locations", len(consLayoutConfigs.Locations))
	//fmt.Println("#FixedLocations", len(consLayoutConfigs.FixedLocations))
	//fmt.Println("#NonFixedLocations", len(consLayoutConfigs.NonFixedLocations))

	// LOAD PHASES
	phases, err := conslay.ReadPhasesFromFile("./data/conslay/staticBuilding.xlsx")
	//phases, err := conslay.ReadPhasesFromFile("./data/conslay/phaseBuilding.xlsx")
	//phases, err := conslay.ReadPhasesFromFile("./data/conslay/dynamicBuilding.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("\tPhases")
	//for i := range phases {
	//	fmt.Println(phases[i])
	//}
	consLayoutConfigs.Phases = phases

	consLayObj, err := conslay.CreateConsLayFromConfig(consLayoutConfigs)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: objectives - select objectives and show configs relevant to those
	//fmt.Println("=== Hoisting Objective ===")
	hoistingTime, err := conslay.ReadHoistingTimeDataFromFile("./data/conslay/f1_hoisting_time_data.xlsx")

	if err != nil {
		log.Fatal(err)
	}

	// Hoisting Objective Configs
	hoistingConfigs := conslay.HoistingConfigs{
		//PrefabricatedLocations: []string{"TF8", "TF9", "TF10"},
		NumberOfFloors: 10,
		HoistingTime: map[string][]conslay.HoistingTime{
			"TF14": hoistingTime,
		},
		FloorHeight:          3.2,
		CraneLocations:       nil,
		ZM:                   2,
		Vuvg:                 37.5,
		Vlvg:                 37.5 / 2,
		Vag:                  50,
		Vwg:                  0.5,
		AlphaHoistingPenalty: 1,
		AlphaHoisting:        0.25,
		BetaHoisting:         1, // beta hoisting = n hoisting
		NHoisting:            1,
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
			BuildingNames: []string{"TF4", "TF5", "TF8", "TF9", "TF10"},
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
	//selectedPref := []string{"TF8", "TF9", "TF10"}
	//
	//hoistingObj.PrefabricatedLocations = selectedPref

	//fmt.Println("\tHoisting Time")
	//for k, v := range hoistingObj.HoistingTime {
	//	hoistingTimeSlice := v
	//	fmt.Printf("Crane name = %s\n", k)
	//	for i := range hoistingTimeSlice {
	//		fmt.Printf(" => Name = %s, Building Name = %s, X = %f, Y = %f, Hoisting Number = %d \n",
	//			hoistingTimeSlice[i].Name,
	//			hoistingTimeSlice[i].BuildingName,
	//			hoistingTimeSlice[i].Coordinate.X,
	//			hoistingTimeSlice[i].Coordinate.Y,
	//			hoistingTimeSlice[i].HoistingNumber,
	//		)
	//	}
	//}

	//fmt.Println("\tCrane Locations")
	//for i := range hoistingObj.CraneLocations {
	//	fmt.Printf("%d: Name = %s, L = %f, W = %f, x = %f, y = %f, radius = %f \n",
	//		i+1,
	//		hoistingObj.CraneLocations[i].Location.Name,
	//		hoistingObj.CraneLocations[i].Location.Length,
	//		hoistingObj.CraneLocations[i].Location.Width,
	//		hoistingObj.CraneLocations[i].Location.Coordinate.X,
	//		hoistingObj.CraneLocations[i].Location.Coordinate.Y,
	//		hoistingObj.CraneLocations[i].Radius,
	//	)
	//}

	// RISK
	hazardInteraction, err := conslay.ReadRiskHazardInteractionDataFromFile("./data/conslay/f2_risk_data.xlsx")

	// Hoisting Objective Configs
	riskConfigs := conslay.RiskConfigs{
		HazardInteractionMatrix: hazardInteraction,
		Delta:                   0.01,
		AlphaRiskPenalty:        100,
		Phases:                  phases,
	}
	riskObj, err := conslay.CreateRiskObjectiveFromConfig(riskConfigs)
	if err != nil {
		log.Fatal(err)
	}

	// Add objectives to conslay_continuous problem
	err = consLayObj.AddObjective(conslay.HoistingObjectiveType, hoistingObj)
	err = consLayObj.AddObjective(conslay.RiskObjectiveType, riskObj)
	if err != nil {
		log.Fatal(err)
	}

	// Add constraints
	outOfBoundsConstraint := conslay.CreateOutOfBoundsConstraint(
		0,
		95,
		0,
		120,
		phases,
		20000,
		1,
	)

	overlapConstraint := conslay.CreateOverlapConstraint(
		phases,
		20000,
		1,
	)

	coverRangeConstraint := conslay.CreateCoverRangeCraneConstraint(
		craneLocations,
		phases,
		20000,
		1,
	)

	zoneConstraint := conslay.CreateInclusiveZoneConstraint(
		[]conslay.Zone{
			{
				Location:      locations["TF13"],
				BuildingNames: []string{"TF7"},
				Size:          20,
			},
			{
				Location:      locations["TF13"],
				BuildingNames: []string{"TF1", "TF2"},
				Size:          15,
			},
		},
		phases,
		20000,
		1,
	)
	err = consLayObj.AddConstraint(conslay.ConstraintOutOfBound, outOfBoundsConstraint)
	err = consLayObj.AddConstraint(conslay.ConstraintOverlap, overlapConstraint)
	err = consLayObj.AddConstraint(conslay.ConstraintInclusiveZone, zoneConstraint)
	err = consLayObj.AddConstraint(conslay.ConstraintsCoverInCraneRadius, coverRangeConstraint)

	//// ZDT1
	//
	//dimensions := 30
	//upperBound := make([]float64, dimensions)
	//lowerBound := make([]float64, dimensions)
	//for i := range upperBound {
	//	upperBound[i] = 1
	//	lowerBound[i] = 0
	//}
	//
	//sphereObjective, _ := multi.CreateZDT1(multi.ZDT1Config{
	//	Dimension:  dimensions,
	//	UpperBound: upperBound,
	//	LowerBound: lowerBound,
	//})

	// MOAHA
	moahaConfigs := moaha.Configs{
		NumAgents:     200,
		NumIterations: 400,
		ArchiveSize:   100,
	}

	algo, err := moaha.Create(consLayObj, moahaConfigs)
	if err != nil {
		return
	}

	err = algo.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("===== Archive Results")
	for i := range algo.Archive {
		fmt.Printf("%d. \n", i+1)
		fmt.Println(algo.Archive[i].Position)
		//fmt.Println(algo.Archive[i].PositionString())
		fmt.Println(algo.Archive[i].Value)
		fmt.Println(algo.Archive[i].Penalty)
	}

	fmt.Println("===== Pareto")
	f1Values := make([]float64, len(algo.Archive))
	f2Values := make([]float64, len(algo.Archive))
	for i := 0; i < 2; i++ {
		var sb strings.Builder
		values := make([]float64, len(algo.Archive))
		for idx, agent := range algo.Archive {
			if idx > 0 {
				sb.WriteString(", ")
			}
			values[idx] = agent.Value[i]
			sb.WriteString(fmt.Sprintf("%g", agent.Value[i]))
		}
		sb.WriteString(";")
		fmt.Println(sb.String())
		if i == 0 {
			f1Values = values
		} else {
			f2Values = values
		}

	}

	fmt.Println("===== Archive Size", len(algo.Archive))

	fmt.Println("Min F1", slices.Min(f1Values))
	fmt.Println("Max F1", slices.Max(f1Values))

	fmt.Println("Min F2", slices.Min(f2Values))
	fmt.Println("Max F2", slices.Max(f2Values))
}
