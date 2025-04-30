package main

import (
	"fmt"
	"github.com/bytedance/sonic"
	"golang-moaha-construction/internal/constraints"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/util"
	"regexp"
	"strings"
)

// formatBuildingNames formats and validates building names according to the required format
// It accepts a string with building names separated by spaces, and returns a slice of valid building names and an error if any invalid names are found
// Valid building names are in the format "TF1", "TF2", etc. or "tf1", "tf2", etc.
func formatBuildingNames(buildingNamesStr string) ([]string, error) {
	// Trim spaces and replace multiple spaces with a single space
	re := regexp.MustCompile(`\s+`)
	trimmed := re.ReplaceAllString(strings.TrimSpace(buildingNamesStr), " ")

	// Split by space
	parts := strings.Split(trimmed, " ")

	// Filter out empty strings and validate format
	var result []string
	var invalidNames []string

	for _, part := range parts {
		if part == "" {
			continue
		}

		// Validate format: must be "TF" or "tf" followed by numbers
		if util.IsTFNumber(part) {
			// Convert to uppercase for consistency
			result = append(result, strings.ToUpper(part))
		} else {
			invalidNames = append(invalidNames, part)
		}
	}

	// Return error if any invalid names were found
	if len(invalidNames) > 0 {
		return result, fmt.Errorf("invalid building names found: %s", strings.Join(invalidNames, ", "))
	}

	return result, nil
}

func (a *App) AddConstraints(cons []ConstraintInput) error {

	problem := a.problem

	_ = problem.InitializeConstraints()

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
				problem.GetPhases(),
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

			_, maxX, _, maxY, err := problem.GetLayoutSize()
			if err != nil {
				return err
			}

			outOfBoundsConstraint := constraints.CreateOutOfBoundsConstraint(
				0,
				maxY,
				0,
				maxX,
				problem.GetPhases(),
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
				facilitiesName, err := formatBuildingNames(zone.BuildingNames)
				if err != nil {
					return err
				}

				var location data.Location

				if loc, ok := problem.GetLocations()[zone.Name]; ok {
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
				problem.GetPhases(),
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

			cranesLocation := make([]data.Crane, len(coverInCraneCfg.CraneLocations))

			for i, craneLocation := range coverInCraneCfg.CraneLocations {

				facilitiesName, err := formatBuildingNames(craneLocation.BuildingNames)
				if err != nil {
					return err
				}

				var location data.Location

				if loc, ok := problem.GetLocations()[craneLocation.Name]; ok {
					location = loc
				}

				cranesLocation[i] = data.Crane{
					Location:     location,
					CraneSymbol:  craneLocation.Name,
					BuildingName: facilitiesName,
					Radius:       craneLocation.Radius,
				}
			}
			fmt.Println(cranesLocation)
			coverRangeConstraint := constraints.CreateCoverRangeCraneConstraint(
				cranesLocation,
				problem.GetPhases(),
				coverInCraneCfg.AlphaCoverInCraneRadiusPenalty,
				coverInCraneCfg.PowerDifferencePenalty,
			)
			err = problem.AddConstraint(con.ConstraintName, coverRangeConstraint)
			if err != nil {
				return err
			}

			err = problem.SetCranesLocations(cranesLocation)
			if err != nil {
				return err
			}

		case constraints.ConstraintSize:

			configBytes, err := sonic.Marshal(con.ConstraintConfig)
			if err != nil {
				return err
			}

			var sizeCfg sizeConfig
			err = sonic.Unmarshal(configBytes, &sizeCfg)
			if err != nil {
				return err
			}

			sizeConstraint := constraints.CreateSizeConstraint(
				sizeCfg.SmallLocations,
				sizeCfg.LargeFacilities,
				sizeCfg.AlphaSizePenalty,
				sizeCfg.PowerDifferencePenalty,
			)

			err = problem.AddConstraint(con.ConstraintName, sizeConstraint)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

type ConstraintsConfigResponse struct {
	OutOfBoundary      any `json:"outOfBoundary,omitempty"`
	Overlap            any `json:"overlap,omitempty"`
	CoverInCraneRadius any `json:"coverInCraneRadius,omitempty"`
	InclusiveZone      any `json:"inclusiveZone,omitempty"`
	Size               any `json:"size,omitempty"`
}

func (a *App) ConstraintsInfo() (*ConstraintsConfigResponse, error) {
	res := &ConstraintsConfigResponse{}

	problemInfo := a.problem

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
				AlphaCoverInCraneRadiusPenalty float64      `json:"alphaCoverInCraneRadiusPenalty"`
				PowerCoverInCraneRadiusPenalty float64      `json:"powerCoverInCraneRadiusPenalty"`
				Phases                         [][]string   `json:"phases"`
				Cranes                         []data.Crane `json:"cranes"`
			}{
				AlphaCoverInCraneRadiusPenalty: coverInCrane.AlphaCoverRangePenalty,
				PowerCoverInCraneRadiusPenalty: coverInCrane.PowerCoverRangePenalty,
				Phases:                         coverInCrane.Phases,
				Cranes:                         coverInCrane.Cranes,
			}

		case constraints.ConstraintSize:

			size := obj.(*constraints.SizeConstraint)
			res.Size = struct {
				AlphaSizePenalty       float64  `json:"alphaSizePenalty"`
				PowerDifferencePenalty float64  `json:"powerDifferencePenalty"`
				SmallLocations         []string `json:"smallLocations"`
				LargeFacilities        []string `json:"largeFacilities"`
			}{
				AlphaSizePenalty:       size.AlphaSizePenalty,
				PowerDifferencePenalty: size.PowerSizePenalty,
				SmallLocations:         size.SmallLocations,
				LargeFacilities:        size.LargeFacilities,
			}
		}
	}

	return res, nil
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
	CraneLocations                 []struct {
		Name          string  `json:"Name"`
		BuildingNames string  `json:"BuildingNames"`
		Radius        float64 `json:"Radius"`
	}
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

type sizeConfig struct {
	AlphaSizePenalty       float64  `json:"AlphaSizePenalty"`
	PowerDifferencePenalty float64  `json:"PowerDifferencePenalty"`
	SmallLocations         []string `json:"SmallLocations"`
	LargeFacilities        []string `json:"LargeFacilities"`
}
