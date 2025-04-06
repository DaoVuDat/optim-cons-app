package main

import (
	"errors"
	"golang-moaha-construction/internal/objectives"
	conslay "golang-moaha-construction/internal/objectives/multi/conslay_continuous"
)

func (a *App) CreateProblem(
	problemName objectives.ProblemType,
	layoutLength, layoutWidth float64,
	facilitiesLocationFilePath, phaseFilePath string,
) error {

	a.problemName = problemName

	// TODO: add GRID problem and PREDETERMINED LOCATIONS problem
	switch problemName {
	case conslay.ContinuousConsLayoutName:
		// Create conslay_continuous problem and add objectives
		consLayoutConfigs := conslay.ConsLayConfigs{
			ConsLayoutLength: layoutLength,
			ConsLayoutWidth:  layoutWidth,
		}

		// LOAD LOCATIONS
		locations, fixedLocations, nonFixedLocations, err := conslay.ReadLocationsFromFile(facilitiesLocationFilePath)
		if err != nil {
			return err
		}

		consLayoutConfigs.Locations = locations
		consLayoutConfigs.NonFixedLocations = nonFixedLocations
		consLayoutConfigs.FixedLocations = fixedLocations

		// LOAD PHASES
		phases, err := conslay.ReadPhasesFromFile(phaseFilePath)

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
	// TODO: add GRID problem and PREDETERMINED LOCATIONS problem
	switch a.problemName {
	case conslay.ContinuousConsLayoutName:

		problemInfo := a.problem.(*conslay.ConsLay)
		return struct {
			LayoutLength      float64
			LayoutWidth       float64
			LowerBound        []float64
			UpperBound        []float64
			Dimensions        int
			Locations         map[string]conslay.Location
			FixedLocations    []conslay.Location
			NonFixedLocations []conslay.Location
		}{
			LayoutLength:      problemInfo.LayoutLength,
			LayoutWidth:       problemInfo.LayoutWidth,
			LowerBound:        problemInfo.LowerBound,
			UpperBound:        problemInfo.UpperBound,
			Dimensions:        problemInfo.Dimensions,
			Locations:         problemInfo.Locations,
			FixedLocations:    problemInfo.FixedLocations,
			NonFixedLocations: problemInfo.NonFixedLocations,
		}, nil
	default:
		return nil, errors.New("not implemented")
	}
}
