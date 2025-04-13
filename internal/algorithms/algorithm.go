package algorithms

import (
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives/objectives"
)

type AlgorithmType string

type AlgorithmResult struct {
	MapLocations  map[string]data.Location
	Value         []float64
	Penalty       map[data.ConstraintType]float64
	ValuesWithKey map[data.ObjectiveType]float64
	Cranes        []objectives.Crane
	Phases        [][]string
}

type Result struct {
	Result []AlgorithmResult
	MinX   float64
	MinY   float64
	MaxX   float64
	MaxY   float64
}

type Algorithm interface {
	Run() error
	RunWithChannel(done chan<- struct{}, channel chan<- any) error
	Type() data.TypeProblem
	GetResults() Result
}
