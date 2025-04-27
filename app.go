package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang-moaha-construction/internal/algorithms"
	"golang-moaha-construction/internal/data"
	eprs "golang-moaha-construction/internal/export-result"
	"golang-moaha-construction/internal/objectives"
	"time"
)

// App struct
type App struct {
	ctx                context.Context
	problemName        data.ProblemName
	problem            objectives.Problem
	algorithmName      algorithms.AlgorithmType
	algorithm          algorithms.Algorithm
	numberOfObjectives int
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SelectFile() (string, error) {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select a File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "All Supported Files (*.png;*.jpg;*.jpeg;*.xlsx)",
				Pattern:     "*.png;*.jpg;*.jpeg;*.xlsx",
			},
		},
		ShowHiddenFiles: false,
	})

	if err != nil {
		return "", err
	}

	//
	//data, err := os.ReadFile(selection)
	//if err != nil {
	//	return "", err
	//}

	return selection, nil
}

func (a *App) SaveFile(commandType CommandType) error {
	switch commandType {
	case ExportResult:
		now := time.Now()

		selection, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
			Title:           "Export Results",
			DefaultFilename: fmt.Sprintf("results_%s.xlsx", now.Format("20060102150405")),
			ShowHiddenFiles: false,
		})
		if err != nil {
			return err
		}

		// TODO: export result

		algoInfo, err := a.AlgorithmInfo()
		if err != nil {
			return err
		}

		problemInfo, err := a.ProblemInfo()
		if err != nil {
			return err
		}

		objectivesInfo, err := a.ObjectivesInfo()
		if err != nil {
			return err
		}

		constraintsInfo, err := a.ConstraintsInfo()
		if err != nil {
			return err
		}

		resultsAny, err := a.Result()
		if err != nil {
			return err
		}

		// Parse the algorithms.Result
		resultsBytes, err := sonic.Marshal(resultsAny)
		if err != nil {
			return err
		}

		var results algorithms.Result
		err = sonic.Unmarshal(resultsBytes, &results)
		if err != nil {
			return err
		}

		err = eprs.WriteXlsxResult(eprs.Options{
			Summary: eprs.Summary{
				AlgorithmInfo:   algoInfo,
				ConstraintsInfo: constraintsInfo,
				ProblemInfo:     problemInfo,
				ObjectivesInfo:  objectivesInfo,
			},
			Results:            results,
			FilePath:           selection,
			ProblemName:        a.problemName,
			AlgorithmName:      a.algorithmName,
			NumberOfObjectives: a.numberOfObjectives,
		})

		if err != nil {
			return err
		}

		return nil

	default:
		return errors.New("invalid command type")
	}

}
