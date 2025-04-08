package main

import (
	"errors"
	"github.com/bytedance/sonic"
	"golang-moaha-construction/internal/objectives"
	continuousconslay "golang-moaha-construction/internal/objectives/multi/conslay_continuous"
	"strings"
)

func (a *App) CreateObjectives(objs []ObjectiveInput) error {

	// TODO: add Config for each objective
	switch a.problemName {
	case continuousconslay.ContinuousConsLayoutName:
		problem := a.problem.(*continuousconslay.ConsLay)

		// remove old objective first
		problem.Objectives = make(map[objectives.ObjectiveType]continuousconslay.Objectiver)

		for _, obj := range objs {
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

				configBytes, err := sonic.Marshal(obj.ObjectiveConfig)
				if err != nil {
					return err
				}

				var hoistingCfg hoistingConfig
				err = sonic.Unmarshal(configBytes, &hoistingCfg)
				if err != nil {
					return err
				}

				hoistingTime := make(map[string][]continuousconslay.HoistingTime, len(hoistingCfg.CraneLocations))
				cranesLocation := make([]continuousconslay.Crane, len(hoistingCfg.CraneLocations))

				for i, craneLocation := range hoistingCfg.CraneLocations {
					hoistingTimeForCrane, err := continuousconslay.ReadHoistingTimeDataFromFile(craneLocation.HoistingTimeFilePath)
					if err != nil {
						return err
					}

					hoistingTime[craneLocation.Name] = hoistingTimeForCrane

					facilitiesName := strings.Split(craneLocation.BuildingNames, " ")

					cranesLocation[i] = continuousconslay.Crane{
						CraneSymbol:  craneLocation.Name,
						BuildingName: facilitiesName,
						Radius:       craneLocation.Radius,
					}
				}

				// setup Cranes Locations and Hoisting Time
				hoistingObj, err := continuousconslay.CreateHoistingObjectiveFromConfig(continuousconslay.HoistingConfigs{
					NumberOfFloors:       hoistingCfg.NumberOfFloors,
					HoistingTime:         hoistingTime,
					FloorHeight:          hoistingCfg.FloorHeight,
					CraneLocations:       cranesLocation,
					ZM:                   hoistingCfg.ZM,
					Vuvg:                 hoistingCfg.Vuvg,
					Vlvg:                 hoistingCfg.Vlvg,
					Vag:                  hoistingCfg.Vag,
					Vwg:                  hoistingCfg.Vwg,
					AlphaHoistingPenalty: hoistingCfg.AlphaHoistingPenalty,
					AlphaHoisting:        hoistingCfg.AlphaHoisting,
					BetaHoisting:         hoistingCfg.BetaHoisting,
					NHoisting:            hoistingCfg.NHoisting,
					Phases:               problem.Phases,
				})
				if err != nil {
					return err

				}
				err = problem.AddObjective(obj.ObjectiveName, hoistingObj)
				if err != nil {
					return err
				}
			case continuousconslay.RiskObjectiveType:
				configBytes, err := sonic.Marshal(obj.ObjectiveConfig)
				if err != nil {
					return err
				}

				var riskCfg riskConfig
				err = sonic.Unmarshal(configBytes, &riskCfg)
				if err != nil {
					return err
				}

				hazardInteractionMatrix, err := continuousconslay.ReadRiskHazardInteractionDataFromFile(riskCfg.HazardInteractionMatrixFilePath)

				if err != nil {
					return err
				}

				riskObj, err := continuousconslay.CreateRiskObjectiveFromConfig(continuousconslay.RiskConfigs{
					HazardInteractionMatrix: hazardInteractionMatrix,
					Delta:                   riskCfg.Delta,
					AlphaRiskPenalty:        riskCfg.AlphaRiskPenalty,
					Phases:                  problem.Phases,
				})
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

type ObjectiveConfigResponse struct {
	Risk     any `json:"risk,omitempty"`
	Hoisting any `json:"hoisting,omitempty"`
	Safety   any `json:"safety,omitempty"`
}

func (a *App) ObjectivesInfo() (*ObjectiveConfigResponse, error) {
	res := &ObjectiveConfigResponse{}

	switch a.problemName {
	case continuousconslay.ContinuousConsLayoutName:

		problemInfo := a.problem.(*continuousconslay.ConsLay)

		objs := problemInfo.GetObjectives()

		for k, obj := range objs {
			switch k {
			case continuousconslay.RiskObjectiveType:
				risk := obj.(*continuousconslay.RiskObjective)

				res.Risk = struct {
					HazardInteractionMatrix [][]float64 `json:"hazardInteractionMatrix"`
					Delta                   float64     `json:"delta"`
					AlphaRiskPenalty        float64     `json:"alphaRiskPenalty"`
					Phases                  [][]string  `json:"phases"`
				}{
					HazardInteractionMatrix: risk.HazardInteractionMatrix,
					Delta:                   risk.Delta,
					AlphaRiskPenalty:        risk.AlphaRiskPenalty,
					Phases:                  risk.Phases,
				}
			case continuousconslay.HoistingObjectiveType:
				hoisting := obj.(*continuousconslay.HoistingObjective)
				res.Hoisting = struct {
					NumberOfFloors       int                                         `json:"numberOfFloors"`
					HoistingTime         map[string][]continuousconslay.HoistingTime `json:"hoistingTime"`
					FloorHeight          float64                                     `json:"floorHeight"`
					CraneLocations       []continuousconslay.Crane                   `json:"craneLocations"`
					ZM                   float64                                     `json:"zm"`
					Vuvg                 float64                                     `json:"vuvg"`
					Vlvg                 float64                                     `json:"vlvg"`
					Vag                  float64                                     `json:"vag"`
					Vwg                  float64                                     `json:"vwg"`
					AlphaHoistingPenalty float64                                     `json:"alphaHoistingPenalty"`
					AlphaHoisting        float64                                     `json:"alphaHoisting"`
					BetaHoisting         float64                                     `json:"betaHoisting"`
					NHoisting            float64                                     `json:"NHoisting"`
					Phases               [][]string                                  `json:"phases"`
				}{
					NumberOfFloors:       hoisting.NumberOfFloors,
					HoistingTime:         hoisting.HoistingTime,
					FloorHeight:          hoisting.FloorHeight,
					CraneLocations:       hoisting.CraneLocations,
					ZM:                   hoisting.ZM,
					Vuvg:                 hoisting.Vuvg,
					Vlvg:                 hoisting.Vlvg,
					Vag:                  hoisting.Vag,
					Vwg:                  hoisting.Vwg,
					AlphaHoistingPenalty: hoisting.AlphaHoistingPenalty,
					AlphaHoisting:        hoisting.AlphaHoisting,
					BetaHoisting:         hoisting.BetaHoisting,
					NHoisting:            hoisting.NHoisting,
					Phases:               hoisting.Phases,
				}
			case continuousconslay.SafetyObjectiveType:

			}
		}

		return res, nil
	default:
		return nil, errors.New("not implemented")
	}
}

type ObjectiveInput struct {
	ObjectiveName   objectives.ObjectiveType `json:"objectiveName"`
	ObjectiveConfig any                      `json:"objectiveConfig"`
}

type hoistingConfig struct {
	CraneLocations []struct {
		Name                 string  `json:"Name"`
		BuildingNames        string  `json:"BuildingNames"`
		Radius               float64 `json:"Radius"`
		HoistingTimeFilePath string  `json:"HoistingTimeFilePath"`
	}
	NumberOfFloors       int     `json:"NumberOfFloors"`
	FloorHeight          float64 `json:"FloorHeight"`
	ZM                   float64 `json:"ZM"`
	Vuvg                 float64 `json:"Vuvg"`
	Vlvg                 float64 `json:"Vlvg"`
	Vag                  float64 `json:"Vag"`
	Vwg                  float64 `json:"Vwg"`
	AlphaHoistingPenalty float64 `json:"AlphaHoistingPenalty"`
	AlphaHoisting        float64 `json:"AlphaHoisting"`
	BetaHoisting         float64 `json:"BetaHoisting"`
	NHoisting            float64 `json:"NHoisting"`
}

type riskConfig struct {
	HazardInteractionMatrixFilePath string  `json:"hazardInteractionMatrixFilePath"`
	Delta                           float64 `json:"Delta"`
	AlphaRiskPenalty                float64 `json:"AlphaRiskPenalty"`
}
