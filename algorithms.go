package main

import (
	"errors"
	"github.com/bytedance/sonic"
	"golang-moaha-construction/internal/algorithms"
	"golang-moaha-construction/internal/algorithms/aha"
	"golang-moaha-construction/internal/algorithms/ga"
	"golang-moaha-construction/internal/algorithms/gwo"
	"golang-moaha-construction/internal/algorithms/moaha"
)

func (a *App) CreateAlgorithm(algorithmInput AlgorithmInput) error {

	switch algorithmInput.AlgorithmName {
	case aha.NameType:
		//configBytes, err := sonic.Marshal(algorithmInput.AlgorithmConfig)
		//if err != nil {
		//	return err
		//}
		//
		//var config ahaConfig
		//err = sonic.Unmarshal(configBytes, &config)
		//if err != nil {
		//	return err
		//}
		//
		//// TODO: change to general problem and result
		//algo, err := aha.Create(a.problem, aha.Config{
		//	NumberOfAgents: config.NumberOfAgents,
		//	NumberOfIter:   config.NumberOfIterations,
		//})
		//
		//if err != nil {
		//	return err
		//}
		//
		//a.algorithm = algo
	case moaha.NameType:
		configBytes, err := sonic.Marshal(algorithmInput.AlgorithmConfig)
		if err != nil {
			return err
		}

		var config moahaConfig
		err = sonic.Unmarshal(configBytes, &config)
		if err != nil {
			return err
		}

		algo, err := moaha.Create(a.problem, moaha.Configs{
			NumAgents:     config.NumberOfAgents,
			NumIterations: config.NumberOfIterations,
			ArchiveSize:   config.ArchiveSize,
		})

		if err != nil {
			return err
		}

		a.algorithm = algo

	case ga.NameType:
		//configBytes, err := sonic.Marshal(algorithmInput.AlgorithmConfig)
		//if err != nil {
		//	return err
		//}
		//
		//var config gaConfig
		//err = sonic.Unmarshal(configBytes, &config)
		//if err != nil {
		//	return err
		//}
		//
		//// TODO: change to general problem and result
		//algo, err := ga.Create(a.problem, ga.Config{
		//	Chromosome:    config.Chromosome,
		//	Generation:    config.Generation,
		//	CrossoverRate: config.CrossoverRate,
		//	MutationRate:  config.MutationRate,
		//	ElitismCount:  config.ElitismCount,
		//})
		//
		//if err != nil {
		//	return err
		//}
		//
		//a.algorithm = algo
	case gwo.NameType:
		//configBytes, err := sonic.Marshal(algorithmInput.AlgorithmConfig)
		//if err != nil {
		//	return err
		//}
		//
		//var config gwoConfig
		//err = sonic.Unmarshal(configBytes, &config)
		//if err != nil {
		//	return err
		//}
		//
		//// TODO: change to general problem and result
		//algo, err := gwo.Create(a.problem, gwo.Config{
		//	NumberOfAgents: config.NumberOfAgents,
		//	NumberOfIter:   config.NumberOfIterations,
		//	AParam:         config.AParam,
		//})
		//
		//if err != nil {
		//	return err
		//}
		//
		//a.algorithm = algo
	default:
		return errors.New("invalid algorithm name")
	}

	return nil
}

func (a *App) AlgorithmInfo() error {

	return nil
}
func (a *App) RunAlgorithm() error {
	err := a.algorithm.Run()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) Result() error {
	return nil
}

type AlgorithmInput struct {
	AlgorithmName   algorithms.AlgorithmType `json:"algorithmName"`
	AlgorithmConfig any                      `json:"algorithmConfig"`
}

type gwoConfig struct {
	NumberOfIterations int     `json:"iterations"`
	NumberOfAgents     int     `json:"population"`
	AParam             float64 `json:"aParam"`
}

type ahaConfig struct {
	NumberOfIterations int `json:"iterations"`
	NumberOfAgents     int `json:"population"`
}

type gaConfig struct {
	Chromosome    int     `json:"chromosome"`
	Generation    int     `json:"generation"`
	CrossoverRate float64 `json:"crossoverRate"`
	MutationRate  float64 `json:"mutationRate"`
	ElitismCount  int     `json:"elitismCount"`
}

type moahaConfig struct {
	NumberOfIterations int `json:"iterations"`
	NumberOfAgents     int `json:"population"`
	ArchiveSize        int `json:"archiveSize"`
}
