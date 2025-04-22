package data

import (
	"regexp"
	"strconv"
	"strings"
)

type TypeProblem int

const (
	Single TypeProblem = iota
	Multi
)

type Coordinate struct {
	X float64
	Y float64
}

type Location struct {
	Coordinate  Coordinate
	Rotation    bool
	Length      float64
	Width       float64
	IsFixed     bool
	Symbol      string
	Name        string
	IsLocatedAt string
}

func (loc Location) ConvertToIdx() (int, error) {
	// strip "TF" and convert to int
	idxStr := strings.Trim(loc.Symbol, "TF")

	idx, err := strconv.Atoi(idxStr)
	if err != nil {
		return 0, err
	}

	// convert back to 0 index-based
	return idx - 1, nil
}

func (loc Location) ConvertToIdxRegex() (int, error) {
	re := regexp.MustCompile(`-?\d+`)
	match := re.FindString(loc.Symbol)

	idx, err := strconv.Atoi(match)
	if err != nil {
		return 0, err
	}

	return idx, nil
}

type Crane struct {
	Location
	BuildingName []string
	Radius       float64
	CraneSymbol  string
}

type ProblemName string
type ObjectiveType string
type ConstraintType string

type Objectiver interface {
	Eval(mapLocations map[string]Location) float64
	GetAlphaPenalty() float64
}

type Constrainter interface {
	Eval(map[string]Location) float64
	GetName() string
	GetAlphaPenalty() float64
	GetPowerPenalty() float64
}
