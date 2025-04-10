package algorithms

import (
	"golang-moaha-construction/internal/data"
)

type AlgorithmType string

type AlgorithmResult struct {
	MapLocations map[string]data.Location
	Value        []float64
	Penalty      map[data.ConstraintType]float64
}

type Algorithm interface {
	Run() error
	RunWithChannel(done chan<- struct{}, channel chan<- any) error
	Type() data.TypeProblem
	GetResults() []AlgorithmResult
}
