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
	"reflect"
	"time"
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
		fmt.Println(selection)
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

		fmt.Println("===== Algo info")
		printStructFields(algoInfo)
		fmt.Println("===== Problem info")
		printStructFields(problemInfo)
		fmt.Println("===== Objectives info")
		printStructFields(objectivesInfo)
		fmt.Println("===== Constraints info")
		printStructFields(constraintsInfo)
		fmt.Println("===== Results")
		fmt.Printf("%+v\n", results)

		err = eprs.WriteXlsxResult(eprs.Options{
			FilePath: selection,
		})

		if err != nil {
			return err
		}

		return nil

	default:
		return errors.New("invalid command type")
	}

}

func printStructFields(s any) {
	// Get the reflected value of input
	val := reflect.ValueOf(s)

	// Follow pointer if it's a pointer to a struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ensure it's a struct now
	if val.Kind() != reflect.Struct {
		fmt.Println("Not a struct!")
		return
	}

	typ := val.Type()

	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			fmt.Printf("%s: %v\n", field.Name, value.Interface())
		}
	}
}

func (a *App) ShowAllInfo() error {

	// TODO
	// show all objectives

	// show all constraints

	// show problem

	// show algorithm

	return nil
}
