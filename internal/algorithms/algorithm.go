package algorithms

import (
	"golang-moaha-construction/internal/data"
)

type Algorithm interface {
	Run() error
	Type() data.TypeProblem
}
