package conslay_continuous

import (
	"maps"
	"math"
	"slices"
)

type Constrainter interface {
	Eval(map[string]Location) float64
}

type CoverRangeCraneConstraint struct {
	Cranes []Crane
}

func (c CoverRangeCraneConstraint) Eval(mapLocations map[string]Location) float64 {
	amount := 0.0
	for i := 0; i < len(c.Cranes); i++ {
		buildings := make([]Location, len(c.Cranes[i].BuildingName))
		for j := 0; j < len(c.Cranes[i].BuildingName); j++ {
			buildings[j] = mapLocations[c.Cranes[i].BuildingName[j]]
		}

		_, val := IsCoverRangeOfCrane(c.Cranes[i], buildings)
		amount += val
	}

	return amount
}

func IsCoverRangeOfCrane(crane Crane, buildings []Location) (bool, float64) {
	invalidAmountTotal := 0.0

	for _, building := range buildings {
		// Top Right
		topRightCoordinate := Coordinate{
			X: building.Coordinate.X + building.Length/2,
			Y: building.Coordinate.Y + building.Width/2,
		}

		topRightAmount := Distance2D(topRightCoordinate, crane.Coordinate) - crane.Radius

		topLeftCoordinate := Coordinate{
			X: building.Coordinate.X - building.Length/2,
			Y: building.Coordinate.Y + building.Width/2,
		}

		topLeftAmount := Distance2D(topLeftCoordinate, crane.Coordinate) - crane.Radius

		bottomLeftCoordinate := Coordinate{
			X: building.Coordinate.X - building.Length/2,
			Y: building.Coordinate.Y - building.Width/2,
		}
		bottomLeftAmount := Distance2D(bottomLeftCoordinate, crane.Coordinate) - crane.Radius

		bottomRightCoordinate := Coordinate{
			X: building.Coordinate.X + building.Length/2,
			Y: building.Coordinate.Y - building.Width/2,
		}
		bottomRightAmount := Distance2D(bottomRightCoordinate, crane.Coordinate) - crane.Radius

		if topRightAmount > 0 || topLeftAmount > 0 || bottomLeftAmount > 0 || bottomRightAmount > 0 {
			invalidAmountTotal += max(0, topLeftAmount) +
				max(0, topRightAmount) +
				max(0, bottomLeftAmount) +
				max(0, bottomRightAmount)

		}
	}

	if invalidAmountTotal > 0 {
		return false, invalidAmountTotal
	} else {
		return true, 0
	}

}

type OverlapConstraint struct {
}

func (c OverlapConstraint) Eval(mapLocations map[string]Location) float64 {
	amount := 0.0

	// convert map to slice
	sliceLocations := slices.Collect(maps.Values(mapLocations))

	for i := 0; i < len(sliceLocations)-1; i++ {
		for j := i + 1; j < len(sliceLocations); j++ {
			_, val := IsOverlapped(sliceLocations[i], sliceLocations[j])
			amount += val
		}
	}

	return amount
}

func IsOverlapped(b1, b2 Location) (bool, float64) {

	l1 := -math.Abs(b1.Coordinate.X-b2.Coordinate.X) + b1.Length/2 + b2.Length/2
	l2 := -math.Abs(b1.Coordinate.Y-b2.Coordinate.Y) + b1.Width/2 + b2.Width/2

	if l1 <= 0 {
		return false, 0
	}

	if l2 <= 0 {
		return false, 0
	}

	return true, math.Max(0, l1) + math.Max(0, l2)
}

type OutOfBoundsConstraint struct {
	MinWidth  float64
	MaxWidth  float64
	MinLength float64
	MaxLength float64
}

func (c OutOfBoundsConstraint) Eval(mapLocations map[string]Location) float64 {
	amount := 0.0
	for _, v := range mapLocations {
		if !v.IsFixed {
			_, val := IsOutOfBound(c.MinLength, c.MaxLength, c.MinWidth, c.MaxWidth, v)
			amount += val
		}
	}
	return amount
}

func IsOutOfBound(minL, maxL, minW, maxW float64, b Location) (bool, float64) {

	l1 := minL + b.Length/2 - b.Coordinate.X
	l2 := b.Coordinate.X + b.Length/2 - maxL
	l3 := minW + b.Width/2 - b.Coordinate.Y
	l4 := b.Coordinate.Y + b.Width/2 - maxW

	if l1 <= 0 && l2 <= 0 && l3 <= 0 && l4 <= 0 {
		return false, 0
	}

	return true, math.Max(0, l1) + math.Max(0, l2) + math.Max(0, l3) + math.Max(0, l4)
}
