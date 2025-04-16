package main

import (
	"errors"
	"github.com/bytedance/sonic"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives/conslay_continuous"
	"golang-moaha-construction/internal/objectives/objectives"
	"strings"
)

func (a *App) CreateObjectives(objs []ObjectiveInput) error {

	// TODO: add Config for each objective
	switch a.problemName {
	case conslay_continuous.ContinuousConsLayoutName:
		problem := a.problem.(*conslay_continuous.ConsLay)

		// remove old objective first
		problem.Objectives = make(map[data.ObjectiveType]data.Objectiver)

		for _, obj := range objs {
			switch obj.ObjectiveName {
			case objectives.SafetyObjectiveType:
				safetyBytes, err := sonic.Marshal(obj.ObjectiveConfig)
				if err != nil {
					return err
				}

				var safetyCfg safetyConfig
				err = sonic.Unmarshal(safetyBytes, &safetyCfg)
				if err != nil {
					return err
				}

				safetyProximityMatrix, err := objectives.ReadSafetyProximityDataFromFile(safetyCfg.SafetyProximityMatrixFilePath)

				if err != nil {
					return err
				}

				safetyObj, err := objectives.CreateSafetyObjectiveFromConfig(objectives.SafetyConfigs{
					SafetyProximity:    safetyProximityMatrix,
					AlphaSafetyPenalty: safetyCfg.AlphaSafetyPenalty,
					Phases:             problem.Phases,
					FilePath:           safetyCfg.SafetyProximityMatrixFilePath,
				})
				if err != nil {
					return err
				}
				err = problem.AddObjective(obj.ObjectiveName, safetyObj)
				if err != nil {
					return err
				}
			case objectives.HoistingObjectiveType:

				configBytes, err := sonic.Marshal(obj.ObjectiveConfig)
				if err != nil {
					return err
				}

				var hoistingCfg hoistingConfig
				err = sonic.Unmarshal(configBytes, &hoistingCfg)
				if err != nil {
					return err
				}

				hoistingTime := make(map[string][]objectives.HoistingTime, len(hoistingCfg.CraneLocations))
				cranesLocation := make([]objectives.Crane, len(hoistingCfg.CraneLocations))
				hoistingTimeWithInfo := make([]objectives.HoistingTimeWithInfo, len(hoistingCfg.CraneLocations))

				for i, craneLocation := range hoistingCfg.CraneLocations {
					hoistingTimeForCrane, err := objectives.ReadHoistingTimeDataFromFile(craneLocation.HoistingTimeFilePath)
					if err != nil {
						return err
					}

					hoistingTime[craneLocation.Name] = hoistingTimeForCrane

					facilitiesName := strings.Split(craneLocation.BuildingNames, " ")

					cranesLocation[i] = objectives.Crane{
						CraneSymbol:  craneLocation.Name,
						BuildingName: facilitiesName,
						Radius:       craneLocation.Radius,
					}

					hoistingTimeWithInfo[i] = objectives.HoistingTimeWithInfo{
						CraneSymbol:  craneLocation.Name,
						FilePath:     craneLocation.HoistingTimeFilePath,
						Radius:       craneLocation.Radius,
						BuildingName: facilitiesName,
					}
				}

				// setup Cranes Locations and Hoisting Time
				hoistingObj, err := objectives.CreateHoistingObjectiveFromConfig(objectives.HoistingConfigs{
					NumberOfFloors:       hoistingCfg.NumberOfFloors,
					HoistingTime:         hoistingTime,
					FloorHeight:          hoistingCfg.FloorHeight,
					CraneLocations:       cranesLocation,
					ZM:                   hoistingCfg.ZM,
					Vuvg:                 hoistingCfg.Vuvg,
					Vlvg:                 hoistingCfg.Vlvg,
					Vag:                  hoistingCfg.Vag,
					Vwg:                  hoistingCfg.Vwg,
					AlphaHoisting:        hoistingCfg.AlphaHoisting,
					BetaHoisting:         hoistingCfg.BetaHoisting,
					NHoisting:            hoistingCfg.NHoisting,
					Phases:               problem.Phases,
					AlphaHoistingPenalty: hoistingCfg.AlphaHoistingPenalty,
					HoistingTimeWithInfo: hoistingTimeWithInfo,
				})

				if err != nil {
					return err

				}
				err = problem.AddObjective(obj.ObjectiveName, hoistingObj)
				if err != nil {
					return err
				}
				problem.CraneLocations = cranesLocation

			case objectives.RiskObjectiveType:
				configBytes, err := sonic.Marshal(obj.ObjectiveConfig)
				if err != nil {
					return err
				}

				var riskCfg riskConfig
				err = sonic.Unmarshal(configBytes, &riskCfg)
				if err != nil {
					return err
				}

				hazardInteractionMatrix, err := objectives.ReadRiskHazardInteractionDataFromFile(riskCfg.HazardInteractionMatrixFilePath)

				if err != nil {
					return err
				}

				riskObj, err := objectives.CreateRiskObjectiveFromConfig(objectives.RiskConfigs{
					HazardInteractionMatrix: hazardInteractionMatrix,
					Delta:                   riskCfg.Delta,
					AlphaRiskPenalty:        riskCfg.AlphaRiskPenalty,
					Phases:                  problem.Phases,
					FilePath:                riskCfg.HazardInteractionMatrixFilePath,
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
	case conslay_continuous.ContinuousConsLayoutName:

		problemInfo := a.problem.(*conslay_continuous.ConsLay)

		objs := problemInfo.GetObjectives()

		for k, obj := range objs {
			switch k {
			case objectives.RiskObjectiveType:
				risk := obj.(*objectives.RiskObjective)

				res.Risk = struct {
					HazardInteractionMatrix [][]float64 `json:"hazardInteractionMatrix"`
					Delta                   float64     `json:"delta"`
					AlphaRiskPenalty        float64     `json:"alphaRiskPenalty"`
					Phases                  [][]string  `json:"phases"`
					FilePath                string      `json:"filePath"`
				}{
					HazardInteractionMatrix: risk.HazardInteractionMatrix,
					Delta:                   risk.Delta,
					AlphaRiskPenalty:        risk.AlphaRiskPenalty,
					Phases:                  risk.Phases,
					FilePath:                risk.FilePath,
				}
			case objectives.HoistingObjectiveType:
				hoisting := obj.(*objectives.HoistingObjective)
				res.Hoisting = struct {
					NumberOfFloors       int                                  `json:"numberOfFloors"`
					FloorHeight          float64                              `json:"floorHeight"`
					ZM                   float64                              `json:"zm"`
					Vuvg                 float64                              `json:"vuvg"`
					Vlvg                 float64                              `json:"vlvg"`
					Vag                  float64                              `json:"vag"`
					Vwg                  float64                              `json:"vwg"`
					AlphaHoisting        float64                              `json:"alphaHoisting"`
					BetaHoisting         float64                              `json:"betaHoisting"`
					NHoisting            float64                              `json:"NHoisting"`
					Phases               [][]string                           `json:"phases"`
					AlphaHoistingPenalty float64                              `json:"alphaHoistingPenalty"`
					HoistingTime         map[string][]objectives.HoistingTime `json:"hoistingTime"`
					CraneLocations       []objectives.Crane                   `json:"craneLocations"`
					HoistingTimeWithInfo []objectives.HoistingTimeWithInfo    `json:"hoistingTimeWithInfo"`
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
					HoistingTimeWithInfo: hoisting.HoistingTimeWithInfo,
				}
			case objectives.SafetyObjectiveType:
				safety := obj.(*objectives.SafetyObjective)

				res.Safety = struct {
					SafetyProximityMatrix [][]float64 `json:"safetyProximityMatrix"`
					AlphaSafetyPenalty    float64     `json:"alphaSafetyPenalty"`
					Phases                [][]string  `json:"phases"`
					FilePath              string      `json:"filePath"`
				}{
					SafetyProximityMatrix: safety.SafetyProximity,
					AlphaSafetyPenalty:    safety.AlphaSafetyPenalty,
					Phases:                safety.Phases,
					FilePath:              safety.FilePath,
				}
			}
		}

		return res, nil
	default:
		return nil, errors.New("not implemented")
	}
}

type ObjectiveInput struct {
	ObjectiveName   data.ObjectiveType `json:"objectiveName"`
	ObjectiveConfig any                `json:"objectiveConfig"`
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

type safetyConfig struct {
	SafetyProximityMatrixFilePath string  `json:"safetyProximityMatrixFilePath"`
	AlphaSafetyPenalty            float64 `json:"AlphaSafetyPenalty"`
}
