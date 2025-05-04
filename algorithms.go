package main

import (
	"errors"
	"github.com/bytedance/sonic"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang-moaha-construction/internal/algorithms"
	"golang-moaha-construction/internal/algorithms/aha"
	"golang-moaha-construction/internal/algorithms/ga"
	"golang-moaha-construction/internal/algorithms/gwo"
	"golang-moaha-construction/internal/algorithms/moaha"
	"golang-moaha-construction/internal/algorithms/mogwo"
	"golang-moaha-construction/internal/algorithms/nsgaii"
	"golang-moaha-construction/internal/algorithms/omoaha"
)

func (a *App) CreateAlgorithm(algorithmInput AlgorithmInput) error {

	a.algorithmName = algorithmInput.AlgorithmName

	switch algorithmInput.AlgorithmName {
	case aha.NameType:
		configBytes, err := sonic.Marshal(algorithmInput.AlgorithmConfig)
		if err != nil {
			return err
		}

		var config ahaConfig
		err = sonic.Unmarshal(configBytes, &config)
		if err != nil {
			return err
		}

		algo, err := aha.Create(a.problem, aha.Config{
			NumberOfAgents: config.NumberOfAgents,
			NumberOfIter:   config.NumberOfIterations,
		})

		if err != nil {
			return err
		}

		a.algorithm = algo
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
	case omoaha.NameType:
		configBytes, err := sonic.Marshal(algorithmInput.AlgorithmConfig)
		if err != nil {
			return err
		}

		var config omoahaConfig
		err = sonic.Unmarshal(configBytes, &config)
		if err != nil {
			return err
		}

		algo, err := omoaha.Create(a.problem, omoaha.Configs{
			NumAgents:     config.NumberOfAgents,
			NumIterations: config.NumberOfIterations,
			ArchiveSize:   config.ArchiveSize,
		})

		if err != nil {
			return err
		}

		a.algorithm = algo

	case ga.NameType:
		configBytes, err := sonic.Marshal(algorithmInput.AlgorithmConfig)
		if err != nil {
			return err
		}

		var config gaConfig
		err = sonic.Unmarshal(configBytes, &config)
		if err != nil {
			return err
		}

		algo, err := ga.Create(a.problem, ga.Config{
			Chromosome:    config.Chromosome,
			Generation:    config.Generation,
			CrossoverRate: config.CrossoverRate,
			MutationRate:  config.MutationRate,
			ElitismCount:  config.ElitismCount,
		})

		if err != nil {
			return err
		}

		a.algorithm = algo
	case gwo.NameType:
		configBytes, err := sonic.Marshal(algorithmInput.AlgorithmConfig)
		if err != nil {
			return err
		}

		var config gwoConfig
		err = sonic.Unmarshal(configBytes, &config)
		if err != nil {
			return err
		}

		algo, err := gwo.Create(a.problem, gwo.Config{
			NumberOfAgents: config.NumberOfAgents,
			NumberOfIter:   config.NumberOfIterations,
			AParam:         config.AParam,
		})

		if err != nil {
			return err
		}

		a.algorithm = algo
	case mogwo.NameType:
		configBytes, err := sonic.Marshal(algorithmInput.AlgorithmConfig)
		if err != nil {
			return err
		}

		var config mogwoConfig
		err = sonic.Unmarshal(configBytes, &config)
		if err != nil {
			return err
		}

		algo, err := mogwo.Create(a.problem, mogwo.Config{
			NumberOfAgents: config.NumberOfAgents,
			NumberOfIter:   config.NumberOfIterations,
			AParam:         config.AParam,
			ArchiveSize:    config.ArchiveSize,
			NumberOfGrids:  config.NumberOfGrids,
			Alpha:          config.Alpha,
			Beta:           config.Beta,
			Gamma:          config.Gamma,
		})

		if err != nil {
			return err
		}

		a.algorithm = algo
	case nsgaii.NameType:
		configBytes, err := sonic.Marshal(algorithmInput.AlgorithmConfig)
		if err != nil {
			return err
		}

		var config nsgaiiConfig
		err = sonic.Unmarshal(configBytes, &config)
		if err != nil {
			return err
		}

		algo, err := nsgaii.Create(a.problem, nsgaii.Config{
			PopulationSize:   config.Chromosome,
			MaxIterations:    config.Generation,
			CrossoverRate:    config.CrossoverRate,
			MutationRate:     config.MutationRate,
			TournamentSize:   config.TournamentSize,
			MutationStrength: config.MutationStrength,
		})

		if err != nil {
			return err
		}

		a.algorithm = algo
	default:
		return errors.New("invalid algorithm name")
	}

	return nil
}

func (a *App) AlgorithmInfo() (any, error) {
	return a.algorithm, nil
}
func (a *App) RunAlgorithm() error {

	progressChan := make(chan any)
	errorChan := make(chan error)
	doneChan := make(chan struct{})
	resultChan := make(chan any, 1)

	go func(doneChan chan<- struct{}, channel chan<- any, errChan chan error) {
		err := a.algorithm.RunWithChannel(doneChan, channel)

		if err != nil {
			errChan <- err
		}

		// send results to resultChan
		resultChan <- a.algorithm.GetResults()

	}(doneChan, progressChan, errorChan)

	// TODO: improve this if it has error
	for progressData := range progressChan {
		runtime.EventsEmit(a.ctx, string(ProgressEvent), progressData)
	}

	runtime.EventsEmit(a.ctx, string(ResultEvent), <-resultChan)

	return nil
}

func (a *App) Result() (any, error) {
	result := a.algorithm.GetResults()

	return result, nil
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

type omoahaConfig struct {
	NumberOfIterations int `json:"iterations"`
	NumberOfAgents     int `json:"population"`
	ArchiveSize        int `json:"archiveSize"`
}

type mogwoConfig struct {
	NumberOfIterations int     `json:"iterations"`
	NumberOfAgents     int     `json:"population"`
	AParam             float64 `json:"aParam"`
	ArchiveSize        int     `json:"archiveSize"`
	NumberOfGrids      int     `json:"numberOfGrids"`
	Alpha              float64 `json:"alpha"`
	Beta               float64 `json:"beta"`
	Gamma              float64 `json:"gamma"`
}

type nsgaiiConfig struct {
	Chromosome       int     `json:"chromosome"`
	Generation       int     `json:"generation"`
	CrossoverRate    float64 `json:"crossoverRate"`
	MutationRate     float64 `json:"mutationRate"`
	MutationStrength float64 `json:"mutationStrength"`
	TournamentSize   int     `json:"tournamentSize"`
}
