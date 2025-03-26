package uidata

import (
	"github.com/charmbracelet/bubbles/list"

	"golang-moaha-construction/internal/algorithms/gwo"
	"golang-moaha-construction/internal/data"

	"golang-moaha-construction/internal/objectives/single"
)

var Objectives = []list.Item{
	data.NewItem(single.SphereName, "Optimize Sphere Function"),
	data.NewItem("conslay", "Optimize Construction Layout"),
}

var Algorithms = []list.Item{
	data.NewItem(gwo.Name, "Grey Wolf Optimizer"),
	data.NewItem("MOAHA", "Multi-Objective Artificial Hummingbird Algorithm"),
}

func GetAlgorithmConfigs(algoStr string) []data.Config {
	switch algoStr {
	case gwo.Name:
		return gwo.Configs

	default:
		return []data.Config{}
	}
}

func GetObjectiveConfigs(objStr string) []data.Config {
	switch objStr {
	case single.SphereName:
		return single.SphereConfigs

	default:
		return []data.Config{}
	}
}
