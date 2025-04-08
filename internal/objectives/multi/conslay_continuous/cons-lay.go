package conslay_continuous

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const (
	ConsLayoutLength = "Construction Layout Length"
	ConsLayoutWidth  = "Construction Layout Width"
	DynamicLocations = "NonFixedLocations"
	StaticLocations  = "FixedLocations"
	Phases           = "Phases"
)

const ContinuousConsLayoutName objectives.ProblemType = "Continuous Construction Layout"

type Coordinate struct {
	X float64
	Y float64
}

type Location struct {
	Coordinate Coordinate
	Rotation   bool
	Length     float64
	Width      float64
	IsFixed    bool
	Symbol     string
	Name       string
}

func (loc Location) ConvertToIdx() (int, error) {
	// strip "TF" and convert to int
	idxStr := strings.Trim(loc.Symbol, "TF")

	idx, err := strconv.Atoi(idxStr)
	if err != nil {
		return 0, err
	}

	// convert back to 0 index-based
	return idx - 1, nil
}

type Objectiver interface {
	Eval(mapLocations map[string]Location) float64
	GetAlphaPenalty() float64
}

type ConsLay struct {
	Dimensions        int
	LayoutLength      float64
	LayoutWidth       float64
	UpperBound        []float64
	LowerBound        []float64
	FixedLocations    []Location
	NonFixedLocations []Location
	Locations         map[string]Location
	Objectives        map[objectives.ObjectiveType]Objectiver
	Constraints       map[objectives.ConstraintType]Constrainter
	Phases            [][]string
	Rounding          bool
}

type ConsLayConfigs struct {
	ConsLayoutLength  float64
	ConsLayoutWidth   float64
	Locations         map[string]Location
	NonFixedLocations []Location
	FixedLocations    []Location
	Phases            [][]string
	Rounding          bool
}

func (s *ConsLay) Type() data.TypeProblem {
	return data.Multi
}

func CreateConsLayFromConfig(consLayConfigs ConsLayConfigs) (*ConsLay, error) {

	consLay := &ConsLay{
		LayoutLength:      consLayConfigs.ConsLayoutLength,
		LayoutWidth:       consLayConfigs.ConsLayoutWidth,
		Locations:         consLayConfigs.Locations,
		FixedLocations:    consLayConfigs.FixedLocations,
		NonFixedLocations: consLayConfigs.NonFixedLocations,
		Phases:            consLayConfigs.Phases,
		Objectives:        make(map[objectives.ObjectiveType]Objectiver),
		Constraints:       make(map[objectives.ConstraintType]Constrainter),
	}

	// Find the x, y, r of Non-fixed Locations
	dimensions := len(consLay.NonFixedLocations) * 3
	upperBound := make([]float64, dimensions)
	lowerBound := make([]float64, dimensions)
	for i := 0; i < len(consLay.NonFixedLocations); i++ {
		idx := i * 3
		// upper bound for x, y, r
		upperBound[idx] = consLay.LayoutLength
		upperBound[idx+1] = consLay.LayoutWidth
		upperBound[idx+2] = 1.0

		// lower bound for x, y, r
		lowerBound[idx] = 0
		lowerBound[idx+1] = 0
		lowerBound[idx+2] = 0
	}

	consLay.Dimensions = dimensions
	consLay.UpperBound = upperBound
	consLay.LowerBound = lowerBound

	return consLay, nil
}

func (s *ConsLay) Eval(input []float64) (values []float64, constraints map[objectives.ConstraintType]float64, penalty map[objectives.ConstraintType]float64) {
	// add x, y, r to non-fixed locations
	nonFixedLocations := make([]Location, len(s.NonFixedLocations))
	mapLocations := make(map[string]Location, len(s.Locations))

	for i := 0; i < len(nonFixedLocations); i++ {
		loc := s.NonFixedLocations[i]
		idx := i * 3
		x := input[idx]
		y := input[idx+1]
		r := input[idx+2]

		width := loc.Width
		length := loc.Length
		rotation := false

		if math.Round(r) > 0 {
			rotation = true
			width = loc.Length
			length = loc.Width
		}

		if s.Rounding {
			x = math.Round(x)
			y = math.Round(y)
		}

		location := Location{
			Coordinate: Coordinate{
				X: x,
				Y: y,
			}, // update x, y
			Rotation: rotation, // update r
			Length:   length,   // change length and width if rotation is true
			Width:    width,
			IsFixed:  false,
			Symbol:   loc.Symbol,
			Name:     loc.Name,
		}

		nonFixedLocations[i] = location
		mapLocations[loc.Symbol] = location
	}

	// add fixed location to mapLocations
	for i := 0; i < len(s.FixedLocations); i++ {
		mapLocations[s.FixedLocations[i].Symbol] = s.FixedLocations[i]
	}

	// checking constraints
	penalty = make(map[objectives.ConstraintType]float64)
	for k, v := range s.Constraints {
		penalty[k] = math.Pow(v.Eval(mapLocations), v.GetPowerPenalty()) * v.GetAlphaPenalty()
	}

	// calculate objectives and add penalty to them
	values = make([]float64, len(s.Objectives))
	valuesName := make([]objectives.ObjectiveType, len(s.Objectives))

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
	}

	return values, map[objectives.ConstraintType]float64{}, penalty
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

func (s *ConsLay) GetObjectives() map[objectives.ObjectiveType]Objectiver {
	return s.Objectives
}

func (s *ConsLay) AddObjective(name objectives.ObjectiveType, objective Objectiver) error {
	if _, ok := s.Objectives[name]; ok {
		return errors.New("the objective has been existed: " + string(name))
	}

	s.Objectives[name] = objective
	return nil
}

func (s *ConsLay) AddConstraint(name objectives.ConstraintType, constraint Constrainter) error {
	if _, ok := s.Constraints[name]; ok {
		return errors.New("the constraint has been existed: " + string(name))
	}

	s.Constraints[name] = constraint
	return nil
}

// Constraints Utility Functions

// Readers Utility Functions

func ReadLocationsFromFile(filePath string) (locations map[string]Location, fixedLocations, nonFixedLocations []Location, err error) {

	// load data from file
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, nil, nil, err
	}

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		return nil, nil, nil, err
	}

	locations = make(map[string]Location)
	fixedLocations = make([]Location, 0)
	nonFixedLocations = make([]Location, 0)

	for rowIdx, row := range rows {
		if rowIdx == 0 {
			continue
		}
		var name string
		var symbol string
		var length float64
		var width float64
		var x float64
		var y float64
		var isFixed = true
		for i, cell := range row {
			switch i {
			case 0:
				name = cell
			case 1:
				symbol = cell
			case 2:
				val, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					return nil, nil, nil, err
				}
				length = val
			case 3:
				val, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					return nil, nil, nil, err
				}
				width = val
			case 4:
				if strings.Contains(cell, "-") {
					isFixed = false
					break
				}
				val, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					return nil, nil, nil, err
				}
				x = val
			case 5:
				if strings.Contains(cell, "-") {
					isFixed = false
					break
				}
				val, err := strconv.ParseFloat(cell, 64)
				if err != nil {
					return nil, nil, nil, err
				}
				y = val
			}
		}

		newLocation := Location{
			Name:   name,
			Symbol: symbol,
			Length: length,
			Width:  width,
			Coordinate: Coordinate{
				X: x,
				Y: y,
			},
			IsFixed:  isFixed,
			Rotation: false,
		}

		if isFixed {
			fixedLocations = append(fixedLocations, newLocation)
		} else {
			nonFixedLocations = append(nonFixedLocations, newLocation)
		}

		locations[symbol] = newLocation
	}

	return locations, fixedLocations, nonFixedLocations, nil
}

func ReadPhasesFromFile(filePath string) ([][]string, error) {

	// load data from file
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	phases := make([][]string, 0)
	for _, row := range rows {
		for i, cell := range row {
			switch i {
			case 1:
				vals := strings.Split(cell, ",")
				for eachTF := range vals {
					vals[eachTF] = strings.TrimSpace(vals[eachTF])
				}
				phases = append(phases, vals)
			}

		}
	}

	return phases, nil
}
