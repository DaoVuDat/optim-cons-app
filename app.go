package main

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	algorithmName  string
	problemName    string
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

func (a *App) CreateObjectives(name []string, config []any) error {

	return nil
}

func (a *App) CreateAlgorithm(name string, config any) error {

	return nil
}

func (a *App) CreateProblem(name string, config any) error {
	return nil
}

func (a *App) RunAlgorithm() error {
	return nil
}

func (a *App) ObjectivesInfo() error {
	return nil
}

func (a *App) AlgorithmInfo() error {
	return nil
}

func (a *App) ProblemInfo() error {
	return nil
}

func (a *App) Result() error {
	return nil
}
