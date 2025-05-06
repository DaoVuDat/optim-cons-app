// Package export_result provides functionality for exporting optimization results to Excel files.
package export_result_old

import (
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/algorithms"
	"golang-moaha-construction/internal/data"
	"regexp"
)

// Global variables for styles
var headerStyle int
var subHeaderStyle int
var contentStyle int
var contentMiddleAlignStyle int
var contentBoldStyle int
var re = regexp.MustCompile(`(?i)objective`) // (?i) = case-insensitive

// Headers for result sheets
var locationHeader = []string{"Name", "Symbol", "x", "y", "Rotated", "Length", "Width", "Fixed"}
var locationHeaderPredetermined = []string{"Symbol", "Is Located At"}

// Summary holds information about the algorithm, constraints, problem, and objectives
type Summary struct {
	AlgorithmInfo   any
	ConstraintsInfo any
	ProblemInfo     any
	ObjectivesInfo  any
}

// Options holds all the parameters needed for exporting results
type Options struct {
	Summary            Summary
	Results            algorithms.Result
	FilePath           string
	ProblemName        data.ProblemName
	AlgorithmName      algorithms.AlgorithmType
	NumberOfObjectives int
}

// writeContentWithValue writes a header and value to the Excel file with appropriate styling
func writeContentWithValue(f *excelize.File, colCount, rowCount int, sheetName string, header string, value any) {
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	_ = f.SetCellValue(sheetName, cell, header)
	_ = f.SetCellStyle(sheetName, cell, cell, contentBoldStyle)
	cell, _ = excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.SetCellValue(sheetName, cell, value)
	_ = f.SetCellStyle(sheetName, cell, cell, contentStyle)
}
