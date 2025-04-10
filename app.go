package main

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang-moaha-construction/internal/algorithms"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
)

// App struct
type App struct {
	ctx            context.Context
	problemName    data.ProblemName
	problem        objectives.Problem
	algorithmName  algorithms.AlgorithmType
	algorithm      algorithms.Algorithm
	objectiveNames []string
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

func (a *App) ShowAllInfo() error {

	// TODO
	// show all objectives

	// show all constraints

	// show problem

	// show algorithm

	return nil
}
