package main

import (
	"golang-moaha-construction/internal/objectives"
	continuousconslay "golang-moaha-construction/internal/objectives/multi/conslay_continuous"
)

var AllProblemsType = []struct {
	Value  objectives.ProblemType
	TSName string // typescript enum name
}{
	{
		Value:  continuousconslay.ContinuousConsLayoutName,
		TSName: "ContinuousConstructionLayout",
	},
	{
		Value:  "Grid Construction Layout",
		TSName: "GridConstructionLayout",
	},
	{
		Value:  "Predetermined Construction Layout",
		TSName: "PredeterminedConstructionLayout",
	},
}
