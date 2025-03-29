package cons_lay

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives/multi"
	"golang-moaha-construction/internal/objectives/single"
	"math"
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

// list constraints

const (
	ConstraintOverlap             = "Overlap"
	ConstraintOutOfBound          = "OutOfBound"
	ConstraintsCoverInCraneRadius = "CoverInCraneRadius"
)

const ConsLayoutName = "Construction-Layout"

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

type ConsLay struct {
	Dimensions        int
	LayoutLength      float64
	LayoutWidth       float64
	UpperBound        []float64
	LowerBound        []float64
	FixedLocations    []Location
	NonFixedLocations []Location
	Locations         map[string]Location
	Objectives        map[string]any
	Constraints       map[string]struct{}
	Phases            [][]string
}

type ConsLayConfigs struct {
	ConsLayoutLength  float64
	ConsLayoutWidth   float64
	Locations         map[string]Location
	NonFixedLocations []Location
	FixedLocations    []Location
	Phases            [][]string
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
		Objectives:        make(map[string]any),
		Constraints:       make(map[string]struct{}),
	}

	// TODO: calculate upper and lower bound

	// TODO: calculate dimension depending on the number of dynamic locations

	return consLay, nil
}

func (s *ConsLay) Eval(x []float64, agent *multi.MultiResult) *multi.MultiResult {
	//time.Sleep(time.Second * 1)

	values := make([]float64, 2)

	sum := 0.0

	for i := 1; i < len(x); i++ {
		sum += x[i]
	}

	var g float64 = 1 + 9*sum/float64(s.Dimensions-1)

	values[0] = x[0]
	values[1] = g * (1 - math.Sqrt(x[0]/g))

	return &multi.MultiResult{
		SingleResult: single.SingleResult{
			Position: x,
			Solution: x,
			Value:    values,
			Idx:      agent.Idx,
		},
		CrowdingDistance: agent.CrowdingDistance,
		Dominated:        agent.Dominated,
		Rank:             agent.Rank,
		DominationSet:    agent.DominationSet,
		DominatedCount:   agent.DominatedCount,
	}

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

func (s *ConsLay) AddObjective(name string, objective any) error {
	if _, ok := s.Objectives[name]; ok {
		return errors.New("the objective has been existed: " + name)
	}

	switch objective.(type) {
	case HoistingObjective:
		// create hoisting objective
		s.Objectives[name] = objective
	case *HoistingObjective:
		s.Objectives[name] = objective
	case RiskObjective:
		// create risk objective
		s.Objectives[name] = objective
	case *RiskObjective:
		s.Objectives[name] = objective
	default:
		return errors.New("invalid objective type: " + name)
	}
	return nil
}

func (s *ConsLay) AddConstraint(name string, constraint any) error {
	if _, ok := s.Constraints[name]; ok {
		return errors.New("the constraint has been existed: " + name)
	}

	s.Constraints[name] = struct{}{}
	return nil
}

// Constraints Utility Functions

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
