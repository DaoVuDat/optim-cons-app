package algorithms

import (
	"golang-moaha-construction/internal/data"
)

type AlgorithmType string

type Algorithm interface {
	Run() error
	RunWithChannel(done chan<- struct{}, channel chan<- any) error
	Type() data.TypeProblem
}
