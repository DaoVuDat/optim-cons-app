// Package export_result provides functionality for exporting optimization results to Excel files.
package export_result

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/objectives/conslay_predetermined"
)

// WriteXlsxResult exports the optimization results to an Excel file.
func WriteXlsxResult(option Options) error {

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var err error

	// Create header style
	headerStyle, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   12,
			Color:  "FFFFFF", // White font
			Family: "Arial",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"4472C4"}, // Blue background
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return err
	}

	// Create sub-header style
	subHeaderStyle, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   12,
			Color:  "000000", // Black font
			Family: "Arial",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"D9E1F2"}, // Light blue background
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "4472C4", Style: 1},
			{Type: "bottom", Color: "4472C4", Style: 1},
			{Type: "left", Color: "4472C4", Style: 1},
			{Type: "right", Color: "4472C4", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
	})
	if err != nil {
		return err
	}

	// Create content style
	contentMiddleAlignStyle, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  11,
			Bold:  true,
			Color: "000000", // Black font
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"F2F2F2"}, // Light gray for alternating rows
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Vertical:   "center",
			Horizontal: "center",
		},
		Protection: &excelize.Protection{
			Locked: false,
		},
	})
	if err != nil {
		return err
	}

	// Create content style
	contentStyle, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  11,
			Color: "000000", // Black font
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"F2F2F2"}, // Light gray for alternating rows
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
		Protection: &excelize.Protection{
			Locked: false,
		},
	})
	if err != nil {
		return err
	}

	// Create content style
	contentBoldStyle, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  11,
			Color: "000000", // Black font
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"F2F2F2"}, // Light gray for alternating rows
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
		Protection: &excelize.Protection{
			Locked: false,
		},
	})
	if err != nil {
		return err
	}

	err = generateSheet1Info(f, option.Summary, option.AlgorithmName)
	if err != nil {
		return err
	}

	if option.ProblemName == conslay_predetermined.PredeterminedConsLayoutName {
		err = generateSheet2ResultsPredetermined(f, option.Results)
		if err != nil {
			return err
		}
	} else {
		err = generateSheet2Results(f, option.Results)
		if err != nil {
			return err
		}
	}

	// pareto
	err = generateSheet3Graph(f, option.Results, option.NumberOfObjectives)
	if err != nil {
		return err
	}

	err = f.SaveAs(option.FilePath)
	if err != nil {
		return err
	}

	return nil
}