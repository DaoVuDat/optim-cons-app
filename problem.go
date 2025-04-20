package main

import (
	"errors"
	"fmt"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives/conslay_continuous"
	"golang-moaha-construction/internal/objectives/conslay_grid"
)

type ProblemInput struct {
	ProblemName      data.ProblemName `json:"problemName"`
	LayoutLength     *float64         `json:"layoutLength"`
	LayoutWidth      *float64         `json:"layoutWidth"`
	FacilitiesFile   *string          `json:"facilitiesFilePath"`
	PhasesFile       *string          `json:"phasesFilePath"`
	GridSize         *int             `json:"gridSize"`
	PredeterminedLoc *string          `json:"predeterminedLoc"`
}

func (a *App) CreateProblem(
	problemInput ProblemInput,
) error {

	a.problemName = problemInput.ProblemName

	// TODO: add GRID problem and PREDETERMINATED LOCATIONS problem
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
	default:
		return errors.New("not implemented")
	}
}

func (a *App) ProblemInfo() (any, error) {
	fmt.Println("Problem Info", a.problemName)
	// type casting to concrete problem
	// TODO: add GRID problem and PREDETERMINATED LOCATIONS problem
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
	default:
		return nil, errors.New("not implemented")
	}
}
