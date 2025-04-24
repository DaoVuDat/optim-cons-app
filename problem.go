package main

import (
	"errors"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives/conslay_continuous"
	"golang-moaha-construction/internal/objectives/conslay_grid"
	"golang-moaha-construction/internal/objectives/conslay_predetermined"
)

type ProblemInput struct {
	ProblemName        data.ProblemName                `json:"problemName"`
	LayoutLength       *float64                        `json:"layoutLength"`
	LayoutWidth        *float64                        `json:"layoutWidth"`
	FacilitiesFile     *string                         `json:"facilitiesFilePath"`
	PhasesFile         *string                         `json:"phasesFilePath"`
	GridSize           *int                            `json:"gridSize"`
	NumberOfLocations  *int                            `json:"numberOfLocations"`
	NumberOfFacilities *int                            `json:"numberOfFacilities"`
	FixedFacilities    *[]conslay_predetermined.LocFac `json:"fixedFacilities"`
}

func (a *App) CreateProblem(
	problemInput ProblemInput,
) error {

	a.problemName = problemInput.ProblemName

	switch problemInput.ProblemName {
	case conslay_continuous.ContinuousConsLayoutName:
		consLayoutConfigs := conslay_continuous.ConsLayConfigs{
			ConsLayoutLength: *problemInput.LayoutLength,
			ConsLayoutWidth:  *problemInput.LayoutWidth,
		}

		// LOAD LOCATIONS
		locations, fixedLocations, nonFixedLocations, err := conslay_continuous.ReadLocationsFromFile(*problemInput.FacilitiesFile)
		if err != nil {
			return err
		}

		consLayoutConfigs.Locations = locations
		consLayoutConfigs.NonFixedLocations = nonFixedLocations
		consLayoutConfigs.FixedLocations = fixedLocations

		// LOAD PHASES
		phases, err := conslay_continuous.ReadPhasesFromFile(*problemInput.PhasesFile)

		if err != nil {
			return err
		}

		consLayoutConfigs.Phases = phases

		consLayObj, err := conslay_continuous.CreateConsLayFromConfig(consLayoutConfigs)
		if err != nil {
			return err
		}

		a.problem = consLayObj
		return nil
	case conslay_grid.GridConsLayoutName:
		consLayoutConfigs := conslay_grid.ConsLayConfigs{
			ConsLayoutLength: *problemInput.LayoutLength,
			ConsLayoutWidth:  *problemInput.LayoutWidth,
			GridSize:         *problemInput.GridSize,
		}

		// LOAD LOCATIONS
		locations, fixedLocations, nonFixedLocations, err := conslay_grid.ReadLocationsFromFile(*problemInput.FacilitiesFile)
		if err != nil {
			return err
		}

		consLayoutConfigs.Locations = locations
		consLayoutConfigs.NonFixedLocations = nonFixedLocations
		consLayoutConfigs.FixedLocations = fixedLocations

		// LOAD PHASES
		phases, err := conslay_grid.ReadPhasesFromFile(*problemInput.PhasesFile)

		if err != nil {
			return err
		}

		consLayoutConfigs.Phases = phases

		consLayObj, err := conslay_grid.CreateConsLayFromConfig(consLayoutConfigs)
		if err != nil {
			return err
		}

		a.problem = consLayObj
		return nil
	case conslay_predetermined.PredeterminedConsLayoutName:
		numberOfLocations := *problemInput.NumberOfLocations
		numberOfFacilities := *problemInput.NumberOfFacilities
		fixedFacilities := *problemInput.FixedFacilities
		consLayoutConfigs := conslay_predetermined.ConsLayConfigs{
			NumberOfLocations:   numberOfLocations,
			NumberOfFacilities:  numberOfFacilities,
			FixedFacilitiesName: fixedFacilities,
		}

		consLayObj, err := conslay_predetermined.CreateConsLayFromConfig(consLayoutConfigs)
		if err != nil {
			return err
		}

		a.problem = consLayObj
		return nil
	default:
		return errors.New("not implemented")
	}
}

func (a *App) ProblemInfo() (any, error) {
	// type casting to concrete problem
	switch a.problemName {
	case conslay_continuous.ContinuousConsLayoutName:
		problemInfo := a.problem.(*conslay_continuous.ConsLay)
		return struct {
			Name              data.ProblemName         `json:"problemName"`
			LayoutLength      float64                  `json:"layoutLength"`
			LayoutWidth       float64                  `json:"layoutWidth"`
			LowerBound        []float64                `json:"lowerBound"`
			UpperBound        []float64                `json:"upperBound"`
			Dimensions        int                      `json:"dimensions"`
			Locations         map[string]data.Location `json:"locations"`
			FixedLocations    []data.Location          `json:"fixedLocations"`
			NonFixedLocations []data.Location          `json:"nonFixedLocations"`
			Phases            [][]string               `json:"phases"`
		}{
			LayoutLength:      problemInfo.LayoutLength,
			LayoutWidth:       problemInfo.LayoutWidth,
			LowerBound:        problemInfo.LowerBound,
			UpperBound:        problemInfo.UpperBound,
			Dimensions:        problemInfo.Dimensions,
			Locations:         problemInfo.Locations,
			FixedLocations:    problemInfo.FixedLocations,
			NonFixedLocations: problemInfo.NonFixedLocations,
			Name:              a.problemName,
			Phases:            problemInfo.Phases,
		}, nil
	case conslay_grid.GridConsLayoutName:
		problemInfo := a.problem.(*conslay_grid.ConsLay)
		return struct {
			Name              data.ProblemName         `json:"problemName"`
			LayoutLength      float64                  `json:"layoutLength"`
			LayoutWidth       float64                  `json:"layoutWidth"`
			GridSize          int                      `json:"gridSize"`
			LowerBound        []float64                `json:"lowerBound"`
			UpperBound        []float64                `json:"upperBound"`
			Dimensions        int                      `json:"dimensions"`
			Locations         map[string]data.Location `json:"locations"`
			FixedLocations    []data.Location          `json:"fixedLocations"`
			NonFixedLocations []data.Location          `json:"nonFixedLocations"`
			Phases            [][]string               `json:"phases"`
		}{
			LayoutLength:      problemInfo.LayoutLength,
			LayoutWidth:       problemInfo.LayoutWidth,
			LowerBound:        problemInfo.LowerBound,
			UpperBound:        problemInfo.UpperBound,
			Dimensions:        problemInfo.Dimensions,
			Locations:         problemInfo.Locations,
			FixedLocations:    problemInfo.FixedLocations,
			NonFixedLocations: problemInfo.NonFixedLocations,
			Name:              a.problemName,
			Phases:            problemInfo.Phases,
			GridSize:          problemInfo.GridSize,
		}, nil
	case conslay_predetermined.PredeterminedConsLayoutName:
		problemInfo := a.problem.(*conslay_predetermined.ConsLay)

		return struct {
			Name                data.ProblemName               `json:"problemName"`
			LowerBound          []float64                      `json:"lowerBound"`
			UpperBound          []float64                      `json:"upperBound"`
			Dimensions          int                            `json:"dimensions"`
			NumberOfLocations   int                            `json:"numberOfLocations"`
			NumberOfFacilities  int                            `json:"numberOfFacilities"`
			FixedFacilitiesName []conslay_predetermined.LocFac `json:"fixedFacilitiesName"`
			AvailableLocations  []string                       `json:"availableLocations"`
		}{
			Name:                a.problemName,
			LowerBound:          problemInfo.LowerBound,
			UpperBound:          problemInfo.UpperBound,
			Dimensions:          problemInfo.Dimensions,
			NumberOfLocations:   problemInfo.NumberOfLocations,
			NumberOfFacilities:  problemInfo.NumberOfFacilities,
			FixedFacilitiesName: problemInfo.FixedFacilitiesName,
			AvailableLocations:  problemInfo.AvailableLocationsIdx,
		}, nil
	default:
		return nil, errors.New("not implemented")
	}
}
