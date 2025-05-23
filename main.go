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
	"golang-moaha-construction/internal/algorithms/mopso"
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

func main() {
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

func main_test() {
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
		HoistingTime: map[string][]objectives.HoistingTime{
			"TF14-B1": hoistingTime,
		},
		Buildings: map[string]objectives.Building{
			"B1": {
				NumberOfFloors: 10,
				FloorHeight:    3.2,
			},
		},
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
		Name   string // for reference from FixedLocations
		Radius float64
	}

	selectedCrane := []SelectedCrane{
		{
			Name:   "TF14-B1",
			Radius: 40,
		},
	}

	craneLocations := make([]data.Crane, 0)
	for _, loc := range selectedCrane {
		craneLocName := strings.Split(loc.Name, "-")[0]
		if _, ok := consLayObj.Locations[craneLocName]; ok {
			craneLocations = append(craneLocations, data.Crane{
				//Location:     craneLoc,
				Radius:      loc.Radius,
				CraneSymbol: loc.Name,
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

	//// MOAHA
	//moahaConfigs := moaha.Configs{
	//	NumAgents:     300,
	//	NumIterations: 400,
	//	ArchiveSize:   100,
	//}
	//
	//algo, err := moaha.Create(consLayObj, moahaConfigs)
	//if err != nil {
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
	//
	// MPSO
	mopsoConfigs := mopso.Config{
		NumberOfAgents: 200,
		NumberOfIter:   400,
		ArchiveSize:    100,
		NumberOfGrids:  20,
		MutationRate:   0.5,
		MaxVelocity:    5,
		C1:             2,
		C2:             2,
		W:              0.4,
	}

	algoMopso, err := mopso.CreateReimpl(consLayObj, mopsoConfigs)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = algoMopso.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("===== Archive Results")
	for i := range algoMopso.Archive {
		fmt.Printf("%d. \n", i+1)
		fmt.Println(algoMopso.Archive[i].Position)
		//fmt.Println(algo.Archive[i].PositionString())
		fmt.Println(algoMopso.Archive[i].Value)
		fmt.Println(algoMopso.Archive[i].Penalty)
	}

	fmt.Println("===== Pareto")
	f1Values := make([]float64, len(algoMopso.Archive))
	f2Values := make([]float64, len(algoMopso.Archive))
	for i := 0; i < 2; i++ {
		var sb strings.Builder
		values := make([]float64, len(algoMopso.Archive))
		for idx, agent := range algoMopso.Archive {
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

	fmt.Println("===== Archive Size", len(algoMopso.Archive))

	fmt.Println("Min F1", slices.Min(f1Values))
	fmt.Println("Max F1", slices.Max(f1Values))

	fmt.Println("Min F2", slices.Min(f2Values))
	fmt.Println("Max F2", slices.Max(f2Values))

	//// NSGA-II
	//nsgaiiConfigs := nsgaii.Config{
	//	PopulationSize:   100,
	//	MaxIterations:    200,
	//	CrossoverRate:    0.7,
	//	MutationRate:     0.4,
	//	MutationStrength: 0.02,
	//	Sigma:            0.1,
	//}
	//
	//algoNSGAII, err := nsgaii.Create(consLayObj, nsgaiiConfigs)
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//
	//err = algoNSGAII.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("===== NSGA-II Archive Results")
	//for i := range algoNSGAII.Archive {
	//	fmt.Printf("%d. \n", i+1)
	//	fmt.Println(algoNSGAII.Archive[i].Position)
	//	fmt.Println(algoNSGAII.Archive[i].Value)
	//	fmt.Println(algoNSGAII.Archive[i].Penalty)
	//}
	//
	//fmt.Println("===== NSGA-II Pareto")
	//f1Values = make([]float64, len(algoNSGAII.Archive))
	//f2Values = make([]float64, len(algoNSGAII.Archive))
	//for i := 0; i < 2; i++ {
	//	var sb strings.Builder
	//	values := make([]float64, len(algoNSGAII.Archive))
	//	for idx, agent := range algoNSGAII.Archive {
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
	//}
	//
	//fmt.Println("===== NSGA-II Archive Size", len(algoNSGAII.Archive))
	//
	//fmt.Println("Min F1", slices.Min(f1Values))
	//fmt.Println("Max F1", slices.Max(f1Values))
	//
	//fmt.Println("Min F2", slices.Min(f2Values))
	//fmt.Println("Max F2", slices.Max(f2Values))
}
