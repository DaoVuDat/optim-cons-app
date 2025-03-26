package cons_lay

import (
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
	"golang-moaha-construction/internal/objectives/multi"
	"golang-moaha-construction/internal/objectives/single"
	"golang-moaha-construction/internal/util"
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

type consLay struct {
	Dimensions int
	UpperBound []float64
	LowerBound []float64
}

func (s *consLay) Type() data.TypeProblem {
	return data.Multi
}

func Create() (objectives.Problem[multi.MultiResult], error) {
	return &consLay{}, nil
}

func (s *consLay) Eval(x []float64, agent *multi.MultiResult) *multi.MultiResult {
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

func (s *consLay) GetUpperBound() []float64 {
	return s.UpperBound
}

func (s *consLay) GetLowerBound() []float64 {
	return s.LowerBound
}

func (s *consLay) GetDimension() int {
	return s.Dimensions
}

func (s *consLay) FindMin() bool {
	return true
}

func (s *consLay) NumberOfObjectives() int {
	return 2
}

const ()

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
		Name: HoistingTime,
	},
	{
		Name: Risk,
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
	IsFixed bool
	Id      int
}

func (s *consLay) LoadData(configs []data.Config) error {
	var dimension int
	var layoutLength float64
	var layoutWidth float64
	var staticLocations []Location
	var dynamicLocations []Location

	//
	for _, config := range configs {
		val := config.Value

		// sanity check
		val = strings.Trim(val, " ")

		switch config.Name {
		case ConsLayoutWidth:
			layoutWidth, _ = strconv.ParseFloat(val, 64)

		case ConsLayoutLength:
			layoutLength, _ = strconv.ParseFloat(val, 64)
		}
	}

	// create upper bound and lower bound and dimension

	return nil
}
