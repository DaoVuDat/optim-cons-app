package cons_lay

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
	"golang-moaha-construction/internal/objectives/multi"
	"golang-moaha-construction/internal/objectives/single"
	"math"
	"slices"
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
	Objectives       []string
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

func (s *ConsLay) AddObjective(objective string) error {
	if slices.Contains(s.Objectives, objective) {
		return errors.New("the objective has been existed")
	}
	s.Objectives = append(s.Objectives, objective)
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

type Location struct {
	X       float64
	Y       float64
	Length  float64
	Width   float64
	IsFixed bool
	Name    string
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
					Name:    name,
					Length:  length,
					Width:   width,
					X:       x,
					Y:       y,
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
		}
	}

	return nil
}
