package single

import (
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
	"golang-moaha-construction/internal/util"
	"strconv"
	"strings"
)

const (
	SphereDimension  = "Dimensions"
	SphereUpperBound = "Upper Bound"
	SphereLowerBound = "Lower Bound"
)

var SphereConfigs = []data.Config{
	{
		Name:               SphereDimension,
		ValidationFunction: util.IsValidPositiveInteger,
	},
	{
		Name:               SphereUpperBound,
		ValidationFunction: util.IsValidFList,
	},
	{
		Name:               SphereLowerBound,
		ValidationFunction: util.IsValidFList,
	},
}

const SphereName = "sphere"

type sphere struct {
	Dimensions int
	UpperBound []float64
	LowerBound []float64
}

func (s *sphere) Type() data.TypeProblem {
	return data.Single
}

func CreateSphere(configs []*data.Config) (objectives.Problem[SingleResult], error) {
	var dimension int
	var upperBound []float64
	var lowerBound []float64

	for _, config := range configs {
		val := config.Value

		// sanity check
		val = strings.Trim(val, " ")

		switch config.Name {
		case SphereDimension:
			dimension, _ = strconv.Atoi(val)
		case SphereLowerBound:
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

		case SphereUpperBound:
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

	return &sphere{
		Dimensions: dimension,
		UpperBound: upperBound,
		LowerBound: lowerBound,
	}, nil

}

func (s *sphere) Eval(x []float64) (values []float64, constraints map[string]float64, penalty map[string]float64) {
	//time.Sleep(time.Second * 1)

	sum := 0.0
	for i := 0; i < len(x); i++ {
		sum += x[i] * x[i]
	}

	return []float64{sum}, map[string]float64{}, map[string]float64{}
}

func (s *sphere) GetUpperBound() []float64 {
	return s.UpperBound
}

func (s *sphere) GetLowerBound() []float64 {
	return s.LowerBound
}

func (s *sphere) GetDimension() int {
	return s.Dimensions
}

func (s *sphere) FindMin() bool {
	return true
}

func (s *sphere) NumberOfObjectives() int {
	return 1
}

func (s *sphere) LoadData(configs []data.Config) error {
	return nil
}

func (s *sphere) AddObjective(name string, objective any) error { return nil }
