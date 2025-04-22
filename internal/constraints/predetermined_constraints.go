package constraints

import (
	"golang-moaha-construction/internal/data"
	"slices"
)

const (
	ConstraintSize data.ConstraintType = "Size"
)

// Cover Range of Crane

type SizeConstraint struct {
	SmallLocations   []string
	LargeFacilities  []string
	Name             data.ConstraintType
	AlphaSizePenalty float64
	PowerSizePenalty float64
}

func CreateSizeConstraint(
	smallLocationNames []string,
	largeFacilityNames []string,
	alphaSizePenalty float64,
	powerSizePenalty float64,
) *SizeConstraint {
	return &SizeConstraint{
		SmallLocations:   smallLocationNames,
		LargeFacilities:  largeFacilityNames,
		Name:             ConstraintSize,
		AlphaSizePenalty: alphaSizePenalty,
		PowerSizePenalty: powerSizePenalty,
	}
}

func (c SizeConstraint) GetName() string {
	return string(c.Name)
}

func (c SizeConstraint) GetAlphaPenalty() float64 {
	return c.AlphaSizePenalty
}

func (c SizeConstraint) GetPowerPenalty() float64 {
	return c.PowerSizePenalty
}

func (c SizeConstraint) Eval(mapLocations map[string]data.Location) float64 {
	// number of invalid locations
	amount := 0.0

	for _, v := range mapLocations {
		// check the facility is whether it is large
		if !slices.Contains(c.LargeFacilities, v.Symbol) {
			continue
		}

		// then check the location is whether it is small if the facility is large
		if slices.Contains(c.SmallLocations, v.IsLocatedAt) {
			amount += 1
		}
	}

	return amount
}
