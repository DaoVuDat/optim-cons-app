/*
Copyright © 2025 Dao Vu Dat dat.daovu@gmail.com
*/
package main

import (
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"golang-moaha-construction/internal/algorithms/omoaha"
	"golang-moaha-construction/internal/constraints"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives/conslay_continuous"
	"golang-moaha-construction/internal/objectives/objectives"
	"log"
	"slices"
	"strings"
)

//go:embed all:frontend/build
var assets embed.FS

func main_test() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Construction Optimization",
		Width:  1400,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop: true,
		},
		DisableResize:    true,
		BackgroundColour: &options.RGBA{R: 240, G: 238, B: 239, A: 255},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		EnumBind: []interface{}{
			AllProblemsType,
			AllObjectivesType,
			AllConstraintsType,
			AllAlgorithmType,
			AllEvent,
			AllCommand,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func main() {
	constructionOptimization()
}

func constructionOptimization() {
	// Create conslay_continuous problem and add objectives

	consLayoutConfigs := conslay_continuous.ConsLayConfigs{
		ConsLayoutLength: 120,
		ConsLayoutWidth:  95,
	}

	// LOAD LOCATIONS
	locations, fixedLocations, nonFixedLocations, err := conslay_continuous.ReadLocationsFromFile("./data/conslay/continuous/locations.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	consLayoutConfigs.Locations = locations
	consLayoutConfigs.NonFixedLocations = nonFixedLocations
	consLayoutConfigs.FixedLocations = fixedLocations

	// LOAD PHASES
	//phases, err := conslay.ReadPhasesFromFile("./data/conslay/staticBuilding.xlsx")
	//phases, err := conslay.ReadPhasesFromFile("./data/conslay/phaseBuilding.xlsx")
	phases, err := conslay_continuous.ReadPhasesFromFile("./data/conslay/continuous/dynamicBuilding.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	consLayoutConfigs.Phases = phases

	consLayObj, err := conslay_continuous.CreateConsLayFromConfig(consLayoutConfigs)
	if err != nil {
		log.Fatal(err)
	}

	hoistingTime, err := objectives.ReadHoistingTimeDataFromFile("./data/conslay/continuous/hoisting_time_data.xlsx")
	//hoistingTime2, err := objectives.ReadHoistingTimeDataFromFile("./data/conslay/VD-1/hoisting_data_crane_2.xlsx")

	if err != nil {
		log.Fatal(err)
	}

	// Hoisting Objective Configs
	hoistingConfigs := objectives.HoistingConfigs{
		NumberOfFloors: 10,
		HoistingTime: map[string][]objectives.HoistingTime{
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

	hoistingObj, err := objectives.CreateHoistingObjectiveFromConfig(hoistingConfigs)
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

	craneLocations := make([]data.Crane, 0)
	for _, loc := range selectedCrane {
		if _, ok := consLayObj.Locations[loc.Name]; ok {
			craneLocations = append(craneLocations, data.Crane{
				//Location:     craneLoc,
				BuildingName: loc.BuildingNames,
				Radius:       loc.Radius,
				CraneSymbol:  loc.Name,
			})
		}
	}

	hoistingObj.CraneLocations = craneLocations

	// RISK
	hazardInteraction, err := objectives.ReadRiskHazardInteractionDataFromFile("./data/conslay/continuous/risk_data.xlsx")

	riskConfigs := objectives.RiskConfigs{
		HazardInteractionMatrix: hazardInteraction,
		Delta:                   0.01,
		AlphaRiskPenalty:        100,
		Phases:                  phases,
	}
	riskObj, err := objectives.CreateRiskObjectiveFromConfig(riskConfigs)
	if err != nil {
		log.Fatal(err)
	}

	safetyMatrix, err := objectives.ReadSafetyProximityDataFromFile("./data/conslay/continuous/safety_data.xlsx")
	if err != nil {
		return
	}

	safetyObj, err := objectives.CreateSafetyObjectiveFromConfig(objectives.SafetyConfigs{
		SafetyProximity:    safetyMatrix,
		AlphaSafetyPenalty: 2000,
		Phases:             phases,
		FilePath:           "",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Add objectives to conslay_continuous problem
	err = consLayObj.AddObjective(objectives.HoistingObjectiveType, hoistingObj)
	err = consLayObj.AddObjective(objectives.RiskObjectiveType, riskObj)
	err = consLayObj.AddObjective(objectives.SafetyObjectiveType, safetyObj)
	if err != nil {
		log.Fatal(err)
	}

	// Add constraints
	outOfBoundsConstraint := constraints.CreateOutOfBoundsConstraint(
		0,
		95,
		0,
		120,
		phases,
		20000,
		1,
	)

	overlapConstraint := constraints.CreateOverlapConstraint(
		phases,
		20000,
		1,
	)

	coverRangeConstraint := constraints.CreateCoverRangeCraneConstraint(
		craneLocations,
		phases,
		20000,
		1,
	)

	zoneConstraint := constraints.CreateInclusiveZoneConstraint(
		[]constraints.Zone{
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
	err = consLayObj.AddConstraint(constraints.ConstraintOutOfBound, outOfBoundsConstraint)
	err = consLayObj.AddConstraint(constraints.ConstraintOverlap, overlapConstraint)
	err = consLayObj.AddConstraint(constraints.ConstraintInclusiveZone, zoneConstraint)
	err = consLayObj.AddConstraint(constraints.ConstraintsCoverInCraneRadius, coverRangeConstraint)

	// OMOAHA
	omoahaConfigs := omoaha.Configs{
		NumAgents:     300,
		NumIterations: 400,
		ArchiveSize:   100,
	}

	algo, err := omoaha.Create(consLayObj, omoahaConfigs)
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

	//// MOGWO
	//mogwoConfigs := mogwo.Config{
	//	NumberOfAgents: 300,
	//	NumberOfIter:   400,
	//	AParam:         2,
	//	ArchiveSize:    100,
	//	NumberOfGrids:  10,
	//	Gamma:          2,
	//	Alpha:          0.1,
	//	Beta:           4,
	//}
	//
	//algo, err := mogwo.Create(consLayObj, mogwoConfigs)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//
	//err = algo.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("===== Archive Results")
	//for i := range algo.Archive {
	//	fmt.Printf("%d. \n", i+1)
	//	fmt.Println(algo.Archive[i].Position)
	//	//fmt.Println(algo.Archive[i].PositionString())
	//	fmt.Println(algo.Archive[i].Value)
	//	fmt.Println(algo.Archive[i].Penalty)
	//}
	//
	//fmt.Println("===== Pareto")
	//f1Values := make([]float64, len(algo.Archive))
	//f2Values := make([]float64, len(algo.Archive))
	//for i := 0; i < 2; i++ {
	//	var sb strings.Builder
	//	values := make([]float64, len(algo.Archive))
	//	for idx, agent := range algo.Archive {
	//		if idx > 0 {
	//			sb.WriteString(", ")
	//		}
	//		values[idx] = agent.Value[i]
	//		sb.WriteString(fmt.Sprintf("%g", agent.Value[i]))
	//	}
	//	sb.WriteString(";")
	//	fmt.Println(sb.String())
	//	if i == 0 {
	//		f1Values = values
	//	} else {
	//		f2Values = values
	//	}
	//
	//}
	//
	//fmt.Println("===== Archive Size", len(algo.Archive))
	//
	//fmt.Println("Min F1", slices.Min(f1Values))
	//fmt.Println("Max F1", slices.Max(f1Values))
	//
	//fmt.Println("Min F2", slices.Min(f2Values))
	//fmt.Println("Max F2", slices.Max(f2Values))
}
