package objectives

import (
	"errors"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives/objectives"
)

var (
	ErrInvalidConfig             = errors.New("invalid configuration")
	ErrInvalidNumberOfObjectives = errors.New("mismatch number of objectives")
)

type Problem interface {
	Eval(pos []float64) (values []float64, valuesWithKey map[data.ObjectiveType]float64, penalty map[data.ConstraintType]float64)
	GetUpperBound() []float64
	GetLowerBound() []float64
	GetDimension() int
	FindMin() bool
	NumberOfObjectives() int
	Type() data.TypeProblem
	GetObjectives() map[data.ObjectiveType]data.Objectiver
	GetConstraints() map[data.ConstraintType]data.Constrainter
	AddObjective(name data.ObjectiveType, objective data.Objectiver) error
	AddConstraint(name data.ConstraintType, constraint data.Constrainter) error
	GetPhases() [][]string
	GetLocationResult(input []float64) (map[string]data.Location, []objectives.Crane, error)
	GetLayoutSize() (minX float64, maxX float64, minY float64, maxY float64, err error)
}
