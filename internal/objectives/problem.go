package objectives

import (
	"errors"
	"golang-moaha-construction/internal/data"
)

type ProblemType string
type ObjectiveType string
type ConstraintType string

var (
	ErrInvalidConfig             = errors.New("invalid configuration")
	ErrInvalidNumberOfObjectives = errors.New("mismatch number of objectives")
)

type Problem interface {
	Eval(pos []float64) (values []float64, constraints map[ConstraintType]float64, penalty map[ConstraintType]float64)
	GetUpperBound() []float64
	GetLowerBound() []float64
	GetDimension() int
	FindMin() bool
	NumberOfObjectives() int
	Type() data.TypeProblem
	AddObjective(string, any) error
}
