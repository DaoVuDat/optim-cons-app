package objectives

import (
	"errors"
	"golang-moaha-construction/internal/data"
)

var (
	ErrInvalidConfig             = errors.New("invalid configuration")
	ErrInvalidNumberOfObjectives = errors.New("mismatch number of objectives")
)

type Result struct {
	Idx         int
	Position    []float64
	Solution    []float64
	Value       []float64
	Constraints []float64
	Penalty     []float64
}

type Problem interface {
	Eval(x []float64) *Result
	GetUpperBound() []float64
	GetLowerBound() []float64
	GetDimension() int
	FindMin() bool
	NumberOfObjectives() int
	Type() data.TypeProblem
}
