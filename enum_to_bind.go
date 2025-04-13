package main

import (
	"golang-moaha-construction/internal/algorithms"
	"golang-moaha-construction/internal/algorithms/aha"
	"golang-moaha-construction/internal/algorithms/ga"
	"golang-moaha-construction/internal/algorithms/gwo"
	"golang-moaha-construction/internal/algorithms/moaha"
	"golang-moaha-construction/internal/constraints"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives/conslay_continuous"
	"golang-moaha-construction/internal/objectives/objectives"
)

var AllProblemsType = []struct {
	Value  data.ProblemName
	TSName string // typescript enum name
}{
	{
		Value:  conslay_continuous.ContinuousConsLayoutName,
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
	Value  data.ObjectiveType
	TSName string // typescript enum name
}{
	{
		Value:  objectives.SafetyObjectiveType,
		TSName: "SafetyObjective",
	},
	{
		Value:  objectives.HoistingObjectiveType,
		TSName: "HoistingObjective",
	},
	{
		Value:  objectives.RiskObjectiveType,
		TSName: "RiskObjective",
	},
}

var AllConstraintsType = []struct {
	Value  data.ConstraintType
	TSName string // typescript enum name
}{
	{
		Value:  constraints.ConstraintOverlap,
		TSName: "Overlap",
	},
	{
		Value:  constraints.ConstraintOutOfBound,
		TSName: "OutOfBound",
	},
	{
		Value:  constraints.ConstraintsCoverInCraneRadius,
		TSName: "CoverInCraneRadius",
	},
	{
		Value:  constraints.ConstraintInclusiveZone,
		TSName: "InclusiveZone",
	},
}

var AllAlgorithmType = []struct {
	Value  algorithms.AlgorithmType
	TSName string // typescript enum name
}{
	{
		Value:  ga.NameType,
		TSName: "GeneticAlgorithm",
	},
	{
		Value:  aha.NameType,
		TSName: "AHA",
	},
	{
		Value:  moaha.NameType,
		TSName: "MOAHA",
	},
	{
		Value:  gwo.NameType,
		TSName: "GWO",
	},
}

type EventType string

const (
	ProgressEvent EventType = "ProgressEvent"
	ResultEvent   EventType = "ResultEvent"
)

var AllEvent = []struct {
	Value  EventType
	TSName string // typescript enum name
}{
	{
		Value:  ProgressEvent,
		TSName: "ProgressEvent",
	},
	{
		Value:  ResultEvent,
		TSName: "ResultEvent",
	},
}

type CommandType string

const (
	ExportResult CommandType = "ExportResult"
)

var AllCommand = []struct {
	Value  CommandType
	TSName string // typescript enum name
}{
	{
		Value:  ExportResult,
		TSName: "ExportResult",
	},
}
