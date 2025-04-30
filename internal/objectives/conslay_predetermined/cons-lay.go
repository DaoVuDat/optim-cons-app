package conslay_predetermined

import (
	"errors"
	"fmt"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/util"
	"math"
	"slices"
	"sort"
)

const PredeterminedConsLayoutName data.ProblemName = "Predetermined Construction Layout"

type ConsLay struct {
	Dimensions            int
	UpperBound            []float64
	LowerBound            []float64
	NumberOfFacilities    int
	NumberOfLocations     int
	FixedFacilitiesName   []LocFac
	Objectives            map[data.ObjectiveType]data.Objectiver
	Constraints           map[data.ConstraintType]data.Constrainter
	Phases                [][]string
	AvailableLocationsIdx []string
	FacilitiesToBeFound   []string
	LocationNames         []string
	FacilityNames         []string
}

type LocFac struct {
	LocName string `json:"locName"`
	FacName string `json:"facName"`
}

type ConsLayConfigs struct {
	NumberOfLocations   int
	NumberOfFacilities  int
	FixedFacilitiesName []LocFac
	Phases              [][]string
	Rounding            bool
}

func (s *ConsLay) Type() data.TypeProblem {
	return data.Multi
}

func CreateConsLayFromConfig(consLayConfigs ConsLayConfigs) (*ConsLay, error) {

	consLay := &ConsLay{
		NumberOfLocations:   consLayConfigs.NumberOfLocations,
		NumberOfFacilities:  consLayConfigs.NumberOfFacilities,
		FixedFacilitiesName: consLayConfigs.FixedFacilitiesName,
		Phases:              consLayConfigs.Phases,
		Objectives:          make(map[data.ObjectiveType]data.Objectiver),
		Constraints:         make(map[data.ConstraintType]data.Constrainter),
	}

	mapLocatedFacilities := make(map[string]struct{}, consLay.NumberOfLocations)

	for _, v := range consLay.FixedFacilitiesName {
		mapLocatedFacilities[v.FacName] = struct{}{}
	}

	locNames := make([]string, consLay.NumberOfLocations)
	for i := 0; i < consLay.NumberOfLocations; i++ {
		locNames[i] = fmt.Sprintf("L%d", i+1)
	}

	consLay.LocationNames = locNames

	facNames := make([]string, consLay.NumberOfFacilities)
	var facilitiesToBeFound []string

	for i := 0; i < consLay.NumberOfFacilities; i++ {
		facName := fmt.Sprintf("TF%d", i+1)
		facNames[i] = facName
		if _, ok := mapLocatedFacilities[facName]; !ok {
			facilitiesToBeFound = append(facilitiesToBeFound, facName)
		}
	}

	consLay.FacilityNames = facNames
	consLay.FacilitiesToBeFound = facilitiesToBeFound

	// Find the x, y, r of Non-fixed Locations
	dimensions := consLay.NumberOfLocations - len(consLay.FixedFacilitiesName)
	upperBound := make([]float64, dimensions)
	lowerBound := make([]float64, dimensions)
	for i := 0; i < dimensions; i++ {
		upperBound[i] = 1.0
		lowerBound[i] = 0
	}

	consLay.Dimensions = dimensions
	consLay.UpperBound = upperBound
	consLay.LowerBound = lowerBound

	mapLocationsToRemove := make(map[string]struct{})
	for _, v := range consLay.FixedFacilitiesName {
		mapLocationsToRemove[v.LocName] = struct{}{}
	}

	// calculate the available locations
	availableLocations := make([]string, consLay.Dimensions)
	availableLocCounter := 0
	for i := 0; i < consLay.NumberOfLocations; i++ {
		curLoc := fmt.Sprintf("L%d", i+1)
		if _, ok := mapLocationsToRemove[curLoc]; !ok {
			availableLocations[availableLocCounter] = curLoc
			availableLocCounter++
		}
	}

	consLay.AvailableLocationsIdx = availableLocations

	return consLay, nil
}

func (s *ConsLay) Eval(input []float64) (
	values []float64,
	valuesWithKey map[data.ObjectiveType]float64,
	key []data.ObjectiveType,
	penalty map[data.ConstraintType]float64) {

	mapLocations := s.MappingLocations(input)

	// checking constraints
	penalty = make(map[data.ConstraintType]float64)
	for k, v := range s.Constraints {
		penalty[k] = math.Pow(v.Eval(mapLocations), v.GetPowerPenalty()) * v.GetAlphaPenalty()
	}

	// calculate objectives and add penalty to them
	values = make([]float64, len(s.Objectives))
	valuesName := make([]data.ObjectiveType, len(s.Objectives))
	valuesWithKey = make(map[data.ObjectiveType]float64, len(s.Objectives))

	i := 0
	for k := range s.Objectives {
		valuesName[i] = k
		i++
	}

	// sort values name
	sort.Slice(valuesName, func(i, j int) bool {
		return valuesName[i] < valuesName[j]
	})

	// sort values in alphabetical order
	for k, v := range s.Objectives {
		idx, ok := slices.BinarySearch(valuesName, k)
		if !ok {
			panic("objective not found")
		}
		val := v.Eval(mapLocations)

		// add penalty to objective value
		for _, penaltyAlpha := range penalty {
			val += penaltyAlpha * v.GetAlphaPenalty()
		}

		values[idx] = val

		// add value to valuesWithKey
		valuesWithKey[k] = val
	}

	return values, valuesWithKey, valuesName, penalty
}

func (s *ConsLay) GetUpperBound() []float64 {
	return s.UpperBound
}

func (s *ConsLay) GetLowerBound() []float64 {
	return s.LowerBound
}

func (s *ConsLay) GetDimension() int {
	return s.Dimensions
}

func (s *ConsLay) FindMin() bool {
	return true
}

func (s *ConsLay) NumberOfObjectives() int {
	return len(s.Objectives)
}

func (s *ConsLay) GetObjectives() map[data.ObjectiveType]data.Objectiver {
	return s.Objectives
}

func (s *ConsLay) GetConstraints() map[data.ConstraintType]data.Constrainter {
	return s.Constraints
}

func (s *ConsLay) AddObjective(name data.ObjectiveType, objective data.Objectiver) error {
	if _, ok := s.Objectives[name]; ok {
		return errors.New("the objective has been existed: " + string(name))
	}

	s.Objectives[name] = objective
	return nil
}

func (s *ConsLay) AddConstraint(name data.ConstraintType, constraint data.Constrainter) error {
	if _, ok := s.Constraints[name]; ok {
		return errors.New("the constraint has been existed: " + string(name))
	}

	s.Constraints[name] = constraint
	return nil
}

func (s *ConsLay) GetPhases() [][]string {
	return s.Phases
}

func (s *ConsLay) GetLocationResult(input []float64) (map[string]data.Location, []data.Location, []data.Crane, error) {

	_, sortedIdx := util.SortWithIdx(input)

	// take the number of non-fixed facilities
	sortedIdx = sortedIdx[:(s.NumberOfFacilities - len(s.FixedFacilitiesName))]

	mapLocations := make(map[string]data.Location, s.NumberOfFacilities)

	// setup fixed facilities
	for _, v := range s.FixedFacilitiesName {
		mapLocations[v.FacName] = data.Location{
			Symbol:      v.FacName,
			IsLocatedAt: v.LocName,
		}
	}

	// setup non-fixed facilities
	for i, v := range sortedIdx {
		loc := s.AvailableLocationsIdx[v]
		fac := s.FacilitiesToBeFound[i]

		mapLocations[fac] = data.Location{
			Symbol:      fac,
			IsLocatedAt: loc,
		}
	}

	sliceLocations := make([]data.Location, s.NumberOfFacilities)

	idx := 0
	for _, v := range mapLocations {
		sliceLocations[idx] = v
		idx++
	}

	sort.Slice(sliceLocations, func(i, j int) bool {
		return util.ExtractNumber(sliceLocations[i].Symbol) < util.ExtractNumber(sliceLocations[j].Symbol)
	})

	return mapLocations, sliceLocations, nil, nil
}

func (s *ConsLay) GetLayoutSize() (minX float64, maxX float64, minY float64, maxY float64, err error) {
	return 0, 0, 0, 0, nil
}

func (s *ConsLay) InitializeObjectives() error {
	s.Objectives = make(map[data.ObjectiveType]data.Objectiver)
	return nil
}

func (s *ConsLay) InitializeConstraints() error {
	s.Constraints = make(map[data.ConstraintType]data.Constrainter)
	return nil
}

func (s *ConsLay) SetCranesLocations(locations []data.Crane) error {
	return nil
}

func (s *ConsLay) GetCranesLocations() []data.Crane {
	return nil
}

func (s *ConsLay) GetLocations() map[string]data.Location {
	return nil
}

// Constraints Utility Functions

// Readers Utility Functions

func (s *ConsLay) MappingLocations(input []float64) map[string]data.Location {

	_, sortedIdx := util.SortWithIdx(input)

	// take the number of non-fixed facilities
	sortedIdx = sortedIdx[:(s.NumberOfFacilities - len(s.FixedFacilitiesName))]

	mapLocations := make(map[string]data.Location, s.NumberOfFacilities)

	// setup fixed facilities
	for _, v := range s.FixedFacilitiesName {
		mapLocations[v.FacName] = data.Location{
			Symbol:      v.FacName,
			IsLocatedAt: v.LocName,
		}
	}

	// setup non-fixed facilities
	for i, v := range sortedIdx {
		loc := s.AvailableLocationsIdx[v]
		fac := s.FacilitiesToBeFound[i]

		mapLocations[fac] = data.Location{
			Symbol:      fac,
			IsLocatedAt: loc,
		}
	}

	return mapLocations
}
