package main

import (
	"errors"
	"golang-moaha-construction/internal/data"
	conslay "golang-moaha-construction/internal/objectives/conslay_continuous"
)

type ProblemInput struct {
	ProblemName      data.ProblemName `json:"problemName"`
	LayoutLength     *float64         `json:"layoutLength"`
	LayoutWidth      *float64         `json:"layoutWidth"`
	FacilitiesFile   *string          `json:"facilitiesFilePath"`
	PhasesFile       *string          `json:"phasesFilePath"`
	GridSize         *string          `json:"gridSize"`
	PredeterminedLoc *string          `json:"predeterminedLoc"`
}

func (a *App) CreateProblem(
	problemInput ProblemInput,
) error {

	a.problemName = problemInput.ProblemName

	// TODO: add GRID problem and PREDETERMINATED LOCATIONS problem
	switch problemInput.ProblemName {
	case conslay.ContinuousConsLayoutName:
		// Create conslay_continuous problem and add objectives
		consLayoutConfigs := conslay.ConsLayConfigs{
			ConsLayoutLength: *problemInput.LayoutLength,
			ConsLayoutWidth:  *problemInput.LayoutWidth,
		}

		// LOAD LOCATIONS
		locations, fixedLocations, nonFixedLocations, err := conslay.ReadLocationsFromFile(*problemInput.FacilitiesFile)
		if err != nil {
			return err
		}

		consLayoutConfigs.Locations = locations
		consLayoutConfigs.NonFixedLocations = nonFixedLocations
		consLayoutConfigs.FixedLocations = fixedLocations

		// LOAD PHASES
		phases, err := conslay.ReadPhasesFromFile(*problemInput.PhasesFile)

		if err != nil {
			return err
		}

		consLayoutConfigs.Phases = phases

		consLayObj, err := conslay.CreateConsLayFromConfig(consLayoutConfigs)
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
	// TODO: add GRID problem and PREDETERMINATED LOCATIONS problem
	switch a.problemName {
	case conslay.ContinuousConsLayoutName:

		problemInfo := a.problem.(*conslay.ConsLay)
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
	default:
		return nil, errors.New("not implemented")
	}
}
