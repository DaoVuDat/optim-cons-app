package multi

import (
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
	"golang-moaha-construction/internal/objectives/single"
	"golang-moaha-construction/internal/util"
	"math"
	"strconv"
	"strings"
)

const (
	ZDT1Dimension  = "Dimensions"
	ZDT1UpperBound = "Upper Bound"
	ZDT1LowerBound = "Lower Bound"
)

var ZDT1Configs = []data.Config{
	{
		Name:               ZDT1Dimension,
		ValidationFunction: util.IsValidPositiveInteger,
	},
	{
		Name:               ZDT1UpperBound,
		ValidationFunction: util.IsValidFList,
	},
	{
		Name:               ZDT1LowerBound,
		ValidationFunction: util.IsValidFList,
	},
}

const ZDT1Name = "ZDT1"

type zdt1 struct {
	Dimensions int
	UpperBound []float64
	LowerBound []float64
}

func (s *zdt1) Type() data.TypeProblem {
	return data.Multi
}

func CreateZDT1(configs []*data.Config) (objectives.Problem[MultiResult], error) {
	var dimension int
	var upperBound []float64
	var lowerBound []float64

	for _, config := range configs {
		val := config.Value

		// sanity check
		val = strings.Trim(val, " ")

		switch config.Name {
		case ZDT1Dimension:
			dimension, _ = strconv.Atoi(val)
		case ZDT1LowerBound:
			if strings.Contains(val, ",") {
				lbs := strings.Split(val, ",")
				for _, lbStr := range lbs {
					lb, _ := strconv.ParseFloat(strings.Trim(lbStr, " "), 64)
					lowerBound = append(lowerBound, lb)
				}
			} else {
				lb, _ := strconv.ParseFloat(val, 64)
				lowerBound = append(lowerBound, lb)
			}

		case ZDT1UpperBound:
			if strings.Contains(val, ",") {
				ubs := strings.Split(val, ",")
				for _, ubStr := range ubs {
					ub, _ := strconv.ParseFloat(strings.Trim(ubStr, " "), 64)
					upperBound = append(upperBound, ub)
				}
			} else {
				ub, _ := strconv.ParseFloat(val, 64)
				upperBound = append(upperBound, ub)
			}
		}
	}

	if len(upperBound) == 1 {
		ub := upperBound[0]
		for i := 1; i < dimension; i++ {
			upperBound = append(upperBound, ub)
		}
	}

	if len(lowerBound) == 1 {
		lb := lowerBound[0]
		for i := 1; i < dimension; i++ {
			lowerBound = append(lowerBound, lb)
		}
	}

	//fmt.Println("===>", dimension, upperBound, lowerBound)

	if (len(upperBound) > 1 && len(upperBound) < dimension) || (len(lowerBound) > 1 && len(lowerBound) < dimension) {
		return nil, objectives.ErrInvalidConfig
	}

	if dimension != len(upperBound) || dimension != len(lowerBound) {
		return nil, objectives.ErrInvalidConfig
	}

	for i := 0; i < dimension; i++ {
		if upperBound[i] < lowerBound[i] {
			return nil, objectives.ErrInvalidConfig
		}
	}

	return &zdt1{
		Dimensions: dimension,
		UpperBound: upperBound,
		LowerBound: lowerBound,
	}, nil

}

func (s *zdt1) Eval(x []float64, agent *MultiResult) *MultiResult {
	//time.Sleep(time.Second * 1)

	values := make([]float64, 2)

	sum := 0.0

	for i := 1; i < len(x); i++ {
		sum += x[i]
	}

	var g float64 = 1 + 9*sum/float64(s.Dimensions-1)

	values[0] = x[0]
	values[1] = g * (1 - math.Sqrt(x[0]/g))

	return &MultiResult{
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

func (s *zdt1) GetUpperBound() []float64 {
	return s.UpperBound
}

func (s *zdt1) GetLowerBound() []float64 {
	return s.LowerBound
}

func (s *zdt1) GetDimension() int {
	return s.Dimensions
}

func (s *zdt1) FindMin() bool {
	return true
}

func (s *zdt1) NumberOfObjectives() int {
	return 2
}

func (s *zdt1) LoadData(configs []data.Config) error {
	return nil
}
