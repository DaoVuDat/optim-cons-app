package export_result

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/algorithms"
	"reflect"
	"strings"
)

var headerStyle int
var subHeaderStyle int
var contentStyle int
var contentMiddleAlignStyle int
var contentBoldStyle int

type Summary struct {
	AlgorithmInfo   any
	ConstraintsInfo any
	ProblemInfo     any
	ObjectivesInfo  any
}

type ResultSummary struct {
	Idx int
}

type Options struct {
	Summary  Summary
	Results  algorithms.Result
	FilePath string
}

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

	err = generateSheet1Info(f, option.Summary)
	if err != nil {
		return err
	}

	err = generateSheet2Results(f, option.Results)
	if err != nil {
		return err
	}

	err = f.SaveAs(option.FilePath)
	if err != nil {
		return err
	}

	return nil
}

// Sheet 1 - Summary

func generateSheet1Info(f *excelize.File, summary Summary) error {
	const SheetName = "Summary"

	// Starting point
	rowCount := 2
	columnCount := 2
	sheetName := f.GetSheetName(0)

	err := f.SetSheetName(sheetName, SheetName)
	if err != nil {
		return err
	}

	// Make Column 2 and 3 = 20
	err = f.SetColWidth(SheetName, "B", "C", 40)
	if err != nil {
		return err
	}

	rowCount = sectionAlgorithm(f, summary.AlgorithmInfo, SheetName, rowCount, columnCount)
	rowCount = sectionProblem(f, summary.ProblemInfo, SheetName, rowCount, columnCount)
	rowCount = sectionObjectives(f, summary.ObjectivesInfo, SheetName, rowCount, columnCount)
	return nil
}

func sectionAlgorithm(f *excelize.File, algorithm any, sheetName string, rowCount int, colCount int) int {
	// Add header
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.MergeCell(sheetName, cell, endCell)
	_ = f.SetCellValue(sheetName, cell, "Algorithm")
	_ = f.SetCellStyle(sheetName, cell, cell, headerStyle)
	rowCount++

	val := reflect.ValueOf(algorithm)
	val = val.Elem() // for pointer
	typ := val.Type()

	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			switch field.Name {
			case "NumberOfAgents":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Number of agents", value.Int())
			case "Population":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Population", value.Int())
			case "ArchiveSize":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Archive size", value.Int())
			case "NumberOfIter":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Number of iterations", value.Int())
			case "Generation":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Generation", value.Int())

			default:
				continue
			}
			rowCount++
		}
	}

	return rowCount + 2
}

func sectionProblem(f *excelize.File, problem any, sheetName string, rowCount int, colCount int) int {
	// Add header
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.MergeCell(sheetName, cell, endCell)
	_ = f.SetCellValue(sheetName, cell, "Problem")
	_ = f.SetCellStyle(sheetName, cell, cell, headerStyle)
	rowCount++

	val := reflect.ValueOf(problem)
	typ := val.Type()

	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			switch field.Name {
			case "LayoutLength":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Layout length", value.Float())
			case "LayoutWidth":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Layout width", value.Float())
			case "Locations":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Number of locations", value.Len())
			case "FixedLocations":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Number of fixed locations", value.Len())
			case "NonFixedLocations":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Number of non-fixed locations", value.Len())
			case "Name":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Name", value.String())
			case "Phases":
				// Add sub-header
				cell, _ = excelize.CoordinatesToCellName(colCount, rowCount)
				endCell, _ = excelize.CoordinatesToCellName(colCount+1, rowCount)
				_ = f.MergeCell(sheetName, cell, endCell)
				_ = f.SetCellValue(sheetName, cell, "Static / Phases / Dynamic")
				_ = f.SetCellStyle(sheetName, cell, cell, subHeaderStyle)
				rowCount++
				for nameIdx := 0; nameIdx < value.Len(); nameIdx++ {
					cell, _ = excelize.CoordinatesToCellName(colCount, rowCount)
					_ = f.SetCellValue(sheetName, cell, nameIdx+1)
					_ = f.SetCellStyle(sheetName, cell, cell, contentStyle)

					names := make([]string, 0)
					subValue := value.Index(nameIdx)
					for subNameIdx := 0; subNameIdx < subValue.Len(); subNameIdx++ {
						names = append(names, subValue.Index(subNameIdx).String())
					}

					cell, _ = excelize.CoordinatesToCellName(colCount+1, rowCount)
					_ = f.SetCellValue(sheetName, cell, strings.Join(names, " "))
					_ = f.SetCellStyle(sheetName, cell, cell, contentStyle)
					rowCount++
				}

			default:
				continue
			}
			rowCount++
		}
	}

	return rowCount + 2
}

func sectionObjectives(f *excelize.File, objectives any, sheetName string, rowCount int, colCount int) int {
	// Add header
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.MergeCell(sheetName, cell, endCell)
	_ = f.SetCellValue(sheetName, cell, "Objectives")
	_ = f.SetCellStyle(sheetName, cell, cell, headerStyle)
	rowCount++

	val := reflect.ValueOf(objectives)
	val = val.Elem() // for pointer
	typ := val.Type()

	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			switch field.Name {
			case "Risk":
				if !value.IsZero() {
					rowCount = riskInfo(f, value.Interface(), sheetName, rowCount, colCount)
				}
			case "Hoisting":
				if !value.IsZero() {
					rowCount = hoistingInfo(f, value.Interface(), sheetName, rowCount, colCount)
				}
			case "Safety":
				if !value.IsZero() {
					fmt.Println(field.Name, value)
				}
			default:
				continue
			}
		}
	}

	return rowCount + 2
}

func riskInfo(f *excelize.File, risk any, sheetName string, rowCount int, colCount int) int {
	// Add sub-header
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.MergeCell(sheetName, cell, endCell)
	_ = f.SetCellValue(sheetName, cell, "Risk")
	_ = f.SetCellStyle(sheetName, cell, cell, subHeaderStyle)
	rowCount++

	val := reflect.ValueOf(risk)
	typ := val.Type()
	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			switch field.Name {
			case "Delta":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Delta", value.Float())
			case "AlphaRiskPenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Alpha (for penalty)", value.Float())
			default:
				continue
			}
			rowCount++
		}
	}

	return rowCount
}

func hoistingInfo(f *excelize.File, hoisting any, sheetName string, rowCount int, colCount int) int {
	// Add sub-header
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.MergeCell(sheetName, cell, endCell)
	_ = f.SetCellValue(sheetName, cell, "Hoisting")
	_ = f.SetCellStyle(sheetName, cell, cell, subHeaderStyle)
	rowCount++
	val := reflect.ValueOf(hoisting)
	//val = val.Elem() // for pointer
	typ := val.Type()

	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			switch field.Name {
			case "NumberOfFloors":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Number of floors", value.Int())
			case "FloorHeight":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Floor height", value.Float())
			case "ZM":
				writeContentWithValue(f, colCount, rowCount, sheetName, "ZM", value.Float())
			case "Vuvg":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Vuvg", value.Float())
			case "Vlvg":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Vlvg", value.Float())
			case "Vag":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Vag", value.Float())
			case "Vwg":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Vwg", value.Float())
			case "AlphaHoistingPenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Alpha (for penalty)", value.Float())
			case "AlphaHoisting":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Alpha", value.Float())
			case "BetaHoisting":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Beta", value.Float())
			case "NHoisting":
				writeContentWithValue(f, colCount, rowCount, sheetName, "NHoisting", value.Float())
			case "HoistingTimeWithInfo":

				// slices
				for j := 0; j < value.Len(); j++ {
					startRow := rowCount
					elem := value.Index(j)

					for k := 0; k < elem.NumField(); k++ {
						subField := elem.Type().Field(k)
						subValue := elem.Field(k)

						switch subField.Name {
						case "BuildingName":
							cell, _ = excelize.CoordinatesToCellName(colCount, startRow+2)
							_ = f.SetCellValue(sheetName, cell, "Facilities")
							_ = f.SetCellStyle(sheetName, cell, cell, contentBoldStyle)

							names := make([]string, 0)
							for nameIdx := 0; nameIdx < subValue.Len(); nameIdx++ {
								names = append(names, subValue.Index(nameIdx).String())
							}

							cell, _ = excelize.CoordinatesToCellName(colCount+1, startRow+2)
							_ = f.SetCellValue(sheetName, cell, strings.Join(names, ", "))
							_ = f.SetCellStyle(sheetName, cell, cell, contentStyle)
						case "FilePath":
							cell, _ = excelize.CoordinatesToCellName(colCount, startRow+1)
							_ = f.SetCellValue(sheetName, cell, "Hoisting file path")
							_ = f.SetCellStyle(sheetName, cell, cell, contentBoldStyle)
							cell, _ = excelize.CoordinatesToCellName(colCount+1, startRow+1)
							_ = f.SetCellValue(sheetName, cell, subValue.String())
							_ = f.SetCellStyle(sheetName, cell, cell, contentStyle)
						case "Radius":
							cell, _ = excelize.CoordinatesToCellName(colCount, startRow+3)
							_ = f.SetCellValue(sheetName, cell, "Radius")
							_ = f.SetCellStyle(sheetName, cell, cell, contentBoldStyle)
							cell, _ = excelize.CoordinatesToCellName(colCount+1, startRow+3)
							_ = f.SetCellValue(sheetName, cell, subValue.Float())
							_ = f.SetCellStyle(sheetName, cell, cell, contentStyle)
						case "CraneSymbol":
							cell, _ = excelize.CoordinatesToCellName(colCount, startRow)
							endCell, _ = excelize.CoordinatesToCellName(colCount+1, startRow)
							_ = f.MergeCell(sheetName, cell, endCell)
							_ = f.SetCellValue(sheetName, cell, subValue.String())
							_ = f.SetCellStyle(sheetName, cell, cell, contentMiddleAlignStyle)
						default:
							continue
						}
						rowCount = startRow + 4
					}
				}
			default:
				continue
			}
			rowCount++
		}
	}

	return rowCount
}

// Sheet 2 - Result

func generateSheet2Results(f *excelize.File, results algorithms.Result) error {

	return nil
}

func writeContentWithValue(f *excelize.File, colCount, rowCount int, sheetName string, header string, value any) {
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	_ = f.SetCellValue(sheetName, cell, header)
	_ = f.SetCellStyle(sheetName, cell, cell, contentBoldStyle)
	cell, _ = excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.SetCellValue(sheetName, cell, value)
	_ = f.SetCellStyle(sheetName, cell, cell, contentStyle)
}
