package main

import (
	"errors"
	"golang-moaha-construction/internal/objectives"
	continuousconslay "golang-moaha-construction/internal/objectives/multi/conslay_continuous"
)

type ObjectiveInput struct {
	ObjectiveName   objectives.ObjectiveType `json:"objectiveName"`
	ObjectiveConfig any                      `json:"objectiveConfig"`
}

func (a *App) CreateObjectives(objectives []ObjectiveInput) error {

	// TODO: add Config for each objective
	switch a.problemName {
	case continuousconslay.ContinuousConsLayoutName:
		problem := a.problem.(*continuousconslay.ConsLay)
		for _, obj := range objectives {
			switch obj.ObjectiveName {
			case continuousconslay.SafetyObjectiveType:
				safetyObj, err := continuousconslay.CreateSafetyObjectiveFromConfig(continuousconslay.SafetyConfigs{})
				if err != nil {
					return err
				}
				err = problem.AddObjective(obj.ObjectiveName, safetyObj)
				if err != nil {
					return err
				}
			case continuousconslay.HoistingObjectiveType:
				hoistingObj, err := continuousconslay.CreateHoistingObjectiveFromConfig(continuousconslay.HoistingConfigs{})
				if err != nil {
					return err

				}
				err = problem.AddObjective(obj.ObjectiveName, hoistingObj)
				if err != nil {
					return err
				}
			case continuousconslay.RiskObjectiveType:
				riskObj, err := continuousconslay.CreateRiskObjectiveFromConfig(continuousconslay.RiskConfigs{})
				if err != nil {
					return err
				}
				err = problem.AddObjective(obj.ObjectiveName, riskObj)
				if err != nil {
					return err
				}
			}
		}
		return nil
	default:
		return errors.New("not implemented")
	}

}

func (a *App) ObjectivesInfo() ([]any, error) {
	return nil, nil
}
