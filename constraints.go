package main

import (
	"errors"
	"github.com/bytedance/sonic"
	"golang-moaha-construction/internal/constraints"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives/conslay_continuous"
	"golang-moaha-construction/internal/objectives/objectives"
	"strings"
)

func (a *App) AddConstraints(cons []ConstraintInput) error {

	switch a.problemName {
	case conslay_continuous.ContinuousConsLayoutName:
		problem := a.problem.(*conslay_continuous.ConsLay)

		// remove old constraints first
		problem.Constraints = make(map[data.ConstraintType]data.Constrainter)

		for _, con := range cons {
			switch con.ConstraintName {
			case constraints.ConstraintOverlap:
				configBytes, err := sonic.Marshal(con.ConstraintConfig)
				if err != nil {
					return err
				}

				var overlapCfg overlapConfig
				err = sonic.Unmarshal(configBytes, &overlapCfg)
				if err != nil {
					return err
				}

				overlapConstraint := constraints.CreateOverlapConstraint(
					problem.Phases,
					overlapCfg.AlphaOverlapPenalty,
					overlapCfg.PowerDifferencePenalty,
				)

				err = problem.AddConstraint(con.ConstraintName, overlapConstraint)
				if err != nil {
					return err
				}
			case constraints.ConstraintOutOfBound:

				configBytes, err := sonic.Marshal(con.ConstraintConfig)
				if err != nil {
					return err
				}

				var outOfBoundCfg outOfBoundConfig
				err = sonic.Unmarshal(configBytes, &outOfBoundCfg)
				if err != nil {
					return err
				}

				outOfBoundsConstraint := constraints.CreateOutOfBoundsConstraint(
					0,
					problem.LayoutWidth,
					0,
					problem.LayoutLength,
					problem.Phases,
					outOfBoundCfg.AlphaOutOfBoundaryPenalty,
					outOfBoundCfg.PowerDifferencePenalty,
				)

				err = problem.AddConstraint(con.ConstraintName, outOfBoundsConstraint)
				if err != nil {
					return err
				}

			case constraints.ConstraintInclusiveZone:
				configBytes, err := sonic.Marshal(con.ConstraintConfig)
				if err != nil {
					return err
				}

				var inclusiveCfg inclusiveZoneConfig
				err = sonic.Unmarshal(configBytes, &inclusiveCfg)
				if err != nil {
					return err
				}

				zones := make([]constraints.Zone, len(inclusiveCfg.Zones))

				for i, zone := range inclusiveCfg.Zones {
					facilitiesName := strings.Split(zone.BuildingNames, " ")

					var location data.Location

					if loc, ok := problem.Locations[zone.Name]; ok {
						location = loc
					}

					zones[i] = constraints.Zone{
						Location:      location,
						BuildingNames: facilitiesName,
						Size:          zone.Size,
					}
				}

				zoneConstraint := constraints.CreateInclusiveZoneConstraint(
					zones,
					problem.Phases,
					inclusiveCfg.AlphaInclusiveZonePenalty,
					inclusiveCfg.PowerDifferencePenalty,
				)

				err = problem.AddConstraint(con.ConstraintName, zoneConstraint)
				if err != nil {
					return err
				}

			case constraints.ConstraintsCoverInCraneRadius:
				configBytes, err := sonic.Marshal(con.ConstraintConfig)
				if err != nil {
					return err
				}

				var coverInCraneCfg coverInCraneRadiusConfig
				err = sonic.Unmarshal(configBytes, &coverInCraneCfg)
				if err != nil {
					return err
				}

				cranes := make([]objectives.Crane, len(problem.CraneLocations))

				for i, crane := range problem.CraneLocations {
					var location data.Location

					if loc, ok := problem.Locations[crane.CraneSymbol]; ok {
						location = loc
					}

					cranes[i] = objectives.Crane{
						Location:     location,
						BuildingName: crane.BuildingName,
						Radius:       crane.Radius,
						CraneSymbol:  crane.CraneSymbol,
					}
				}

				coverRangeConstraint := constraints.CreateCoverRangeCraneConstraint(
					cranes,
					problem.Phases,
					coverInCraneCfg.AlphaCoverInCraneRadiusPenalty,
					coverInCraneCfg.PowerDifferencePenalty,
				)
				err = problem.AddConstraint(con.ConstraintName, coverRangeConstraint)
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

type ConstraintsConfigResponse struct {
	OutOfBoundary      any `json:"outOfBoundary,omitempty"`
	Overlap            any `json:"overlap,omitempty"`
	CoverInCraneRadius any `json:"coverInCraneRadius,omitempty"`
	InclusiveZone      any `json:"inclusiveZone,omitempty"`
}

func (a *App) ConstraintsInfo() (*ConstraintsConfigResponse, error) {
	res := &ConstraintsConfigResponse{}

	switch a.problemName {
	case conslay_continuous.ContinuousConsLayoutName:

		problemInfo := a.problem.(*conslay_continuous.ConsLay)

		cons := problemInfo.GetConstraints()

		for k, obj := range cons {
			switch k {
			case constraints.ConstraintOutOfBound:
				outOfBound := obj.(*constraints.OutOfBoundsConstraint)

				res.OutOfBoundary = struct {
					MinWidth               float64    `json:"minWidth"`
					MaxWidth               float64    `json:"maxWidth"`
					MinLength              float64    `json:"minLength"`
					MaxLength              float64    `json:"maxLength"`
					AlphaOutOfBoundPenalty float64    `json:"alphaOutOfBoundPenalty"`
					PowerOutOfBoundPenalty float64    `json:"powerOutOfBoundPenalty"`
					Phases                 [][]string `json:"phases"`
				}{
					MinWidth:               outOfBound.MinWidth,
					MaxWidth:               outOfBound.MaxWidth,
					MinLength:              outOfBound.MinLength,
					MaxLength:              outOfBound.MaxLength,
					AlphaOutOfBoundPenalty: outOfBound.AlphaOutOfBoundsPenalty,
					PowerOutOfBoundPenalty: outOfBound.PowerOutOfBoundsPenalty,
					Phases:                 outOfBound.Phases,
				}
			case constraints.ConstraintOverlap:
				overlap := obj.(*constraints.OverlapConstraint)
				res.Overlap = struct {
					AlphaOverlapPenalty float64    `json:"alphaOverlapPenalty"`
					PowerOverlapPenalty float64    `json:"powerOverlapPenalty"`
					Phases              [][]string `json:"phases"`
				}{
					AlphaOverlapPenalty: overlap.AlphaOverlapPenalty,
					PowerOverlapPenalty: overlap.PowerOverlapPenalty,
					Phases:              overlap.Phases,
				}
			case constraints.ConstraintInclusiveZone:
				inclusiveZone := obj.(*constraints.InclusiveZoneConstraint)
				res.InclusiveZone = struct {
					AlphaInclusivePenalty float64            `json:"alphaInclusivePenalty"`
					PowerInclusivePenalty float64            `json:"powerInclusivePenalty"`
					Phases                [][]string         `json:"phases"`
					Zones                 []constraints.Zone `json:"zones"`
				}{
					AlphaInclusivePenalty: inclusiveZone.AlphaInclusiveZonePenalty,
					PowerInclusivePenalty: inclusiveZone.PowerInclusiveZonePenalty,
					Phases:                inclusiveZone.Phases,
					Zones:                 inclusiveZone.Zones,
				}

			case constraints.ConstraintsCoverInCraneRadius:
				coverInCrane := obj.(*constraints.CoverRangeCraneConstraint)
				res.CoverInCraneRadius = struct {
					AlphaCoverInCraneRadiusPenalty float64            `json:"alphaCoverInCraneRadiusPenalty"`
					PowerCoverInCraneRadiusPenalty float64            `json:"powerCoverInCraneRadiusPenalty"`
					Phases                         [][]string         `json:"phases"`
					Cranes                         []objectives.Crane `json:"cranes"`
				}{
					AlphaCoverInCraneRadiusPenalty: coverInCrane.AlphaCoverRangePenalty,
					PowerCoverInCraneRadiusPenalty: coverInCrane.PowerCoverRangePenalty,
					Phases:                         coverInCrane.Phases,
					Cranes:                         coverInCrane.Cranes,
				}

			}
		}

		return res, nil
	default:
		return nil, errors.New("not implemented")
	}
}

type ConstraintInput struct {
	ConstraintName   data.ConstraintType `json:"constraintName"`
	ConstraintConfig any                 `json:"constraintConfig"`
}

type outOfBoundConfig struct {
	AlphaOutOfBoundaryPenalty float64 `json:"AlphaOutOfBoundaryPenalty"`
	PowerDifferencePenalty    float64 `json:"PowerDifferencePenalty"`
}

type overlapConfig struct {
	AlphaOverlapPenalty    float64 `json:"AlphaOverLapPenalty"`
	PowerDifferencePenalty float64 `json:"PowerDifferencePenalty"`
}

type coverInCraneRadiusConfig struct {
	AlphaCoverInCraneRadiusPenalty float64 `json:"AlphaCoverInCraneRadiusPenalty"`
	PowerDifferencePenalty         float64 `json:"PowerDifferencePenalty"`
}

type inclusiveZoneConfig struct {
	AlphaInclusiveZonePenalty float64 `json:"AlphaInclusiveZonePenalty"`
	PowerDifferencePenalty    float64 `json:"PowerDifferencePenalty"`
	Zones                     []struct {
		Name          string  `json:"Name"`
		BuildingNames string  `json:"BuildingNames"`
		Size          float64 `json:"Size"`
	} `json:"Zones"`
}
