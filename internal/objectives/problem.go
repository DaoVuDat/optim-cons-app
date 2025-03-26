package objectives

import (
	"errors"
	"golang-moaha-construction/internal/data"
)

var (
	ErrInvalidConfig             = errors.New("invalid configuration")
	ErrInvalidNumberOfObjectives = errors.New("mismatch number of objectives")
)

type LoadDataType struct {
	Name  string
	Value string
}

type Problem[T any] interface {
	Eval(pos []float64, x *T) *T
	GetUpperBound() []float64
	GetLowerBound() []float64
	GetDimension() int
	FindMin() bool
	NumberOfObjectives() int
	Type() data.TypeProblem
	LoadData(configs []data.Config) error
	AddObjective(string) error
}
