package algorithms

import (
	"golang-moaha-construction/internal/data"
)

type AlgorithmType string

type AlgorithmResult struct {
	MapLocations   map[string]data.Location
	SliceLocations []data.Location
	Value          []float64
	Penalty        map[data.ConstraintType]float64
	ValuesWithKey  map[data.ObjectiveType]float64
	Key            []data.ObjectiveType
	Cranes         []data.Crane
	Phases         [][]string
}

type Result struct {
	Result      []AlgorithmResult
	Convergence []float64
	MinX        float64
	MinY        float64
	MaxX        float64
	MaxY        float64
}

type Algorithm interface {
	Run() error
	RunWithChannel(done chan<- struct{}, channel chan<- any) error
	Type() data.TypeProblem
	GetResults() Result
}
