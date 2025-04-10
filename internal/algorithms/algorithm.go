package algorithms

import (
	"golang-moaha-construction/internal/data"
)

type AlgorithmType string

type Algorithm interface {
	Run() error
	Type() data.TypeProblem
}
