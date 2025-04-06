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

var AllObjectivesType = []struct {
	Value  objectives.ObjectiveType
	TSName string // typescript enum name
}{
	{
		Value:  continuousconslay.SafetyObjectiveType,
		TSName: "SafetyObjective",
	},
	{
		Value:  continuousconslay.HoistingObjectiveType,
		TSName: "HoistingObjective",
	},
	{
		Value:  continuousconslay.RiskObjectiveType,
		TSName: "RiskObjective",
	},
}

var AllConstraintsType = []struct {
	Value  objectives.ConstraintType
	TSName string // typescript enum name
}{
	{
		Value:  continuousconslay.ConstraintOverlap,
		TSName: "Overlap",
	},
	{
		Value:  continuousconslay.ConstraintOutOfBound,
		TSName: "OutOfBound",
	},
	{
		Value:  continuousconslay.ConstraintsCoverInCraneRadius,
		TSName: "CoverInCraneRadius",
	},
	{
		Value:  continuousconslay.ConstraintInclusiveZone,
		TSName: "InclusiveZone",
	},
}
