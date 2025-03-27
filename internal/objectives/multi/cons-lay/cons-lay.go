package cons_lay

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
	"golang-moaha-construction/internal/objectives/multi"
	"golang-moaha-construction/internal/objectives/single"
	"math"
	"strconv"
	"strings"
)

const (
	ConsLayoutLength = "Construction Layout Length"
	ConsLayoutWidth  = "Construction Layout Width"
	DynamicLocations = "DynamicLocations"
	StaticLocations  = "StaticLocations"
	Phases           = "Phases"
)

const ConsLayoutName = "Construction-Layout"

type ConsLay struct {
	Dimensions       int
	LayoutLength     float64
	LayoutWidth      float64
	UpperBound       []float64
	LowerBound       []float64
	StaticLocations  []Location
	DynamicLocations []Location
	Objectives       map[string]any
	Phases           [][]string
}

func (s *ConsLay) Type() data.TypeProblem {
	return data.Multi
}

func Create() (objectives.Problem[multi.MultiResult], error) {
	return &ConsLay{}, nil
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

func (s *ConsLay) AddObjective(name string, objective any) error {
	if _, ok := s.Objectives[name]; ok {
		return errors.New("the objective has been existed: " + name)
	}

	switch objective.(type) {
	case HoistingObjective:
		// create hoisting objective
		s.Objectives[name] = objective.(HoistingObjective)
	case RiskObjective:
		// create risk objective
		s.Objectives[name] = objective.(RiskObjective)
	default:
		return errors.New("invalid objective type: " + name)
	}
	return nil
}

var Configs = []data.Config{
	{
		Name: DynamicLocations,
	},
	{
		Name: StaticLocations,
	},
	{
		Name: Phases,
	},
	{
		Name: ConsLayoutLength,
	},
	{
		Name: ConsLayoutWidth,
	},
}

type Coordinate struct {
	X float64
	Y float64
}

type Location struct {
	Coordinate Coordinate
	Length     float64
	Width      float64
	IsFixed    bool
	Name       string
}

func (s *ConsLay) LoadData(configs []data.Config) error {

	for _, config := range configs {

		val := config.Value

		// sanity check
		val = strings.Trim(val, " ")

		switch config.Name {
		case ConsLayoutWidth:
			layoutWidth, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			s.LayoutWidth = layoutWidth

		case ConsLayoutLength:
			layoutLength, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}

			s.LayoutLength = layoutLength
		case StaticLocations:
			// load data from file
			file, err := excelize.OpenFile(val)
			if err != nil {
				return err
			}

			rows, err := file.GetRows("Sheet1")
			if err != nil {
				return err
			}
			for rowIdx, row := range rows {
				if rowIdx == 0 {
					continue
				}
				var name string
				var length float64
				var width float64
				var x float64
				var y float64
				for i, cell := range row {
					switch i {
					case 0:
						name = cell
					case 1:
						val, err := strconv.ParseFloat(cell, 64)
						if err != nil {
							return err
						}
						length = val
					case 2:
						val, err := strconv.ParseFloat(cell, 64)
						if err != nil {
							return err
						}
						width = val
					case 3:
						val, err := strconv.ParseFloat(cell, 64)
						if err != nil {
							return err
						}
						x = val
					case 4:
						val, err := strconv.ParseFloat(cell, 64)
						if err != nil {
							return err
						}
						y = val
					}
				}

				s.StaticLocations = append(s.StaticLocations, Location{
					Name:   name,
					Length: length,
					Width:  width,
					Coordinate: Coordinate{
						X: x,
						Y: y,
					},
					IsFixed: true,
				})
			}

		case DynamicLocations:
			// load data from file
			file, err := excelize.OpenFile(val)
			if err != nil {
				return err
			}

			rows, err := file.GetRows("Sheet1")
			if err != nil {
				return err
			}
			for rowIdx, row := range rows {
				if rowIdx == 0 {
					continue
				}
				var name string
				var length float64
				var width float64
				for i, cell := range row {
					switch i {
					case 0:
						name = cell
					case 1:
						val, err := strconv.ParseFloat(cell, 64)
						if err != nil {
							return err
						}
						length = val
					case 2:
						val, err := strconv.ParseFloat(cell, 64)
						if err != nil {
							return err
						}
						width = val
					}
				}

				s.DynamicLocations = append(s.DynamicLocations, Location{
					Name:    name,
					Length:  length,
					Width:   width,
					IsFixed: false,
				})
			}
		case Phases:
			// load data from file
			file, err := excelize.OpenFile(val)
			if err != nil {
				return err
			}

			rows, err := file.GetRows("Sheet1")
			if err != nil {
				return err
			}

			for _, row := range rows {
				for i, cell := range row {
					switch i {
					case 1:
						vals := strings.Split(cell, ",")
						for eachTF := range vals {
							vals[eachTF] = strings.TrimSpace(vals[eachTF])
						}
						s.Phases = append(s.Phases, vals)
					}

				}
			}

		}

	}
	return nil
}
