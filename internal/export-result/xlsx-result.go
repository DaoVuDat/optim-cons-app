package export_result

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/algorithms"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives/conslay_predetermined"
	"reflect"
	"regexp"
	"strings"
)

var headerStyle int
var subHeaderStyle int
var contentStyle int
var contentMiddleAlignStyle int
var contentBoldStyle int
var re = regexp.MustCompile(`(?i)objective`) // (?i) = case-insensitive

type Summary struct {
	AlgorithmInfo   any
	ConstraintsInfo any
	ProblemInfo     any
	ObjectivesInfo  any
}

type Options struct {
	Summary            Summary
	Results            algorithms.Result
	FilePath           string
	ProblemName        data.ProblemName
	AlgorithmName      algorithms.AlgorithmType
	NumberOfObjectives int
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

// Sheet 1 - Summary

func generateSheet1Info(f *excelize.File, summary Summary, algorithmName algorithms.AlgorithmType) error {
	const SheetName = "Summary"

	// Starting point
	rowCount := 2
	columnCount := 2
	sheetName := f.GetSheetName(0)

	err := f.SetSheetName(sheetName, SheetName)
	if err != nil {
		return err
	}

	err = f.SetColWidth(SheetName, "B", "B", 40)
	err = f.SetColWidth(SheetName, "C", "C", 80)
	if err != nil {
		return err
	}

	rowCount = sectionAlgorithm(f, summary.AlgorithmInfo, algorithmName, SheetName, rowCount, columnCount)
	rowCount = sectionProblem(f, summary.ProblemInfo, SheetName, rowCount, columnCount)
	rowCount = sectionObjectives(f, summary.ObjectivesInfo, SheetName, rowCount, columnCount)
	rowCount = sectionConstraints(f, summary.ConstraintsInfo, SheetName, rowCount, columnCount)
	return nil
}

func sectionAlgorithm(f *excelize.File, algorithm any, algorithmName algorithms.AlgorithmType, sheetName string, rowCount int, colCount int) int {
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

	writeContentWithValue(f, colCount, rowCount, sheetName, "Name", algorithmName)
	rowCount++

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
			case "GridSize":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Grid size", value.Int())
			case "Locations":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Number of locations", value.Len())
			case "FixedLocations":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Number of fixed locations", value.Len())
			case "NonFixedLocations":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Number of non-fixed locations", value.Len())
			case "NumberOfFacilities":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Number of facilities", value.Int())
			case "NumberOfLocations":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Number of locations", value.Int())
			case "FixedFacilitiesName":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Number of located facilities", value.Len())
			case "Name":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Name", value.String())
			case "Phases":
				if value.Len() == 0 {
					break
				}
				// Add sub-header
				cell, _ = excelize.CoordinatesToCellName(colCount, rowCount)
				endCell, _ = excelize.CoordinatesToCellName(colCount+1, rowCount)
				_ = f.MergeCell(sheetName, cell, endCell)
				_ = f.SetCellValue(sheetName, cell, "Static / Phases / Dynamic")
				_ = f.SetCellStyle(sheetName, cell, cell, subHeaderStyle)

				for nameIdx := 0; nameIdx < value.Len(); nameIdx++ {
					rowCount++
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
					rowCount = safetyInfo(f, value.Interface(), sheetName, rowCount, colCount)
				}
			case "SafetyHazard":
				if !value.IsZero() {
					rowCount = safetyHazardInfo(f, value.Interface(), sheetName, rowCount, colCount)
				}
			case "TransportCost":
				if !value.IsZero() {
					rowCount = transportCostInfo(f, value.Interface(), sheetName, rowCount, colCount)
				}
			case "ConstructionCost":
				if !value.IsZero() {
					rowCount = constructionCostInfo(f, value.Interface(), sheetName, rowCount, colCount)
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
			case "FilePath":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Hazard Interaction Matrix file path", value.String())
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
							_ = f.SetCellValue(sheetName, cell, strings.Join(names, " "))
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
						rowCount = startRow + 3
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

func safetyInfo(f *excelize.File, safety any, sheetName string, rowCount int, colCount int) int {
	// Add sub-header
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.MergeCell(sheetName, cell, endCell)
	_ = f.SetCellValue(sheetName, cell, "Safety")
	_ = f.SetCellStyle(sheetName, cell, cell, subHeaderStyle)
	rowCount++

	val := reflect.ValueOf(safety)
	typ := val.Type()
	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			switch field.Name {
			case "AlphaSafetyPenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Alpha (for penalty)", value.Float())
			case "FilePath":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Safety Proximity Matrix file path", value.String())
			default:
				continue
			}
			rowCount++
		}
	}

	return rowCount
}

func safetyHazardInfo(f *excelize.File, safetyHazard any, sheetName string, rowCount int, colCount int) int {
	// Add sub-header
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.MergeCell(sheetName, cell, endCell)
	_ = f.SetCellValue(sheetName, cell, "Safety Hazard")
	_ = f.SetCellStyle(sheetName, cell, cell, subHeaderStyle)
	rowCount++

	val := reflect.ValueOf(safetyHazard)
	typ := val.Type()
	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			switch field.Name {
			case "AlphaSafetyHazardPenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Alpha (for penalty)", value.Float())
			case "FilePath":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Safety and Environmental Concerns Matrix file path", value.String())
			default:
				continue
			}
			rowCount++
		}
	}

	return rowCount
}

func transportCostInfo(f *excelize.File, transportCost any, sheetName string, rowCount int, colCount int) int {
	// Add sub-header
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.MergeCell(sheetName, cell, endCell)
	_ = f.SetCellValue(sheetName, cell, "Transport Cost")
	_ = f.SetCellStyle(sheetName, cell, cell, subHeaderStyle)
	rowCount++

	val := reflect.ValueOf(transportCost)
	typ := val.Type()
	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			switch field.Name {
			case "AlphaTransportCostPenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Alpha (for penalty)", value.Float())
			case "FilePath":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Facilities Interaction Matrix file path", value.String())
			default:
				continue
			}
			rowCount++
		}
	}

	return rowCount
}

func constructionCostInfo(f *excelize.File, transportCost any, sheetName string, rowCount int, colCount int) int {
	// Add sub-header
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.MergeCell(sheetName, cell, endCell)
	_ = f.SetCellValue(sheetName, cell, "Construction Cost")
	_ = f.SetCellStyle(sheetName, cell, cell, subHeaderStyle)
	rowCount++

	val := reflect.ValueOf(transportCost)
	typ := val.Type()
	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			switch field.Name {
			case "AlphaCCPenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Alpha (for penalty)", value.Float())
			case "FrequencyMatrixFilePath":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Frequency Matrix file path", value.String())
			case "DistanceMatrixFilePath":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Distance Matrix file path", value.String())
			case "GeneralQAP":
				writeContentWithValue(f, colCount, rowCount, sheetName, "General QAP", value.Bool())
			default:
				continue
			}
			rowCount++
		}
	}

	return rowCount
}

func sectionConstraints(f *excelize.File, constraints any, sheetName string, rowCount int, colCount int) int {
	// Add header
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.MergeCell(sheetName, cell, endCell)
	_ = f.SetCellValue(sheetName, cell, "Constraints")
	_ = f.SetCellStyle(sheetName, cell, cell, headerStyle)
	rowCount++

	val := reflect.ValueOf(constraints)
	val = val.Elem() // for pointer
	typ := val.Type()

	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			switch field.Name {
			case "OutOfBoundary":
				if !value.IsZero() {
					rowCount = outOfBoundaryInfo(f, value.Interface(), sheetName, rowCount, colCount)
				}
			case "Overlap":
				if !value.IsZero() {
					rowCount = overlapInfo(f, value.Interface(), sheetName, rowCount, colCount)
				}
			case "CoverInCraneRadius":
				if !value.IsZero() {
					rowCount = coverCraneInfo(f, value.Interface(), sheetName, rowCount, colCount)
				}
			case "InclusiveZone":
				if !value.IsZero() {
					rowCount = inclusiveZoneInfo(f, value.Interface(), sheetName, rowCount, colCount)
				}
			case "Size":
				if !value.IsZero() {
					rowCount = sizeInfo(f, value.Interface(), sheetName, rowCount, colCount)
				}
			default:
				continue
			}
		}
	}

	return rowCount + 2
}

func outOfBoundaryInfo(f *excelize.File, outOfBound any, sheetName string, rowCount int, colCount int) int {
	// Add sub-header
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.MergeCell(sheetName, cell, endCell)
	_ = f.SetCellValue(sheetName, cell, "Out Of Boundary")
	_ = f.SetCellStyle(sheetName, cell, cell, subHeaderStyle)
	rowCount++
	val := reflect.ValueOf(outOfBound)
	typ := val.Type()
	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			switch field.Name {
			case "MinWidth":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Min width", value.Float())
			case "MaxWidth":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Max width", value.Float())
			case "MinLength":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Min length", value.Float())
			case "MaxLength":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Max length", value.Float())
			case "PowerOutOfBoundPenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Power difference (for penalty)", value.Float())
			case "AlphaOutOfBoundPenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Alpha (for penalty)", value.Float())
			default:
				continue
			}
			rowCount++
		}
	}

	return rowCount
}

func overlapInfo(f *excelize.File, overlap any, sheetName string, rowCount int, colCount int) int {
	// Add sub-header
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.MergeCell(sheetName, cell, endCell)
	_ = f.SetCellValue(sheetName, cell, "Overlap")
	_ = f.SetCellStyle(sheetName, cell, cell, subHeaderStyle)
	rowCount++
	val := reflect.ValueOf(overlap)
	typ := val.Type()
	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			switch field.Name {
			case "PowerOverlapPenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Power difference (for penalty)", value.Float())
			case "AlphaOverlapPenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Alpha (for penalty)", value.Float())
			default:
				continue
			}
			rowCount++
		}
	}

	return rowCount
}

func inclusiveZoneInfo(f *excelize.File, inclusive any, sheetName string, rowCount int, colCount int) int {
	// Add sub-header
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.MergeCell(sheetName, cell, endCell)
	_ = f.SetCellValue(sheetName, cell, "Inclusive Zone")
	_ = f.SetCellStyle(sheetName, cell, cell, subHeaderStyle)
	rowCount++
	val := reflect.ValueOf(inclusive)
	typ := val.Type()
	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			switch field.Name {
			case "Zones":
				for j := 0; j < value.Len(); j++ {
					startRow := rowCount
					elem := value.Index(j)

					for k := 0; k < elem.NumField(); k++ {
						subField := elem.Type().Field(k)
						subValue := elem.Field(k)
						switch subField.Name {
						case "BuildingNames":
							cell, _ = excelize.CoordinatesToCellName(colCount, startRow+1)
							_ = f.SetCellValue(sheetName, cell, "Facilities")
							_ = f.SetCellStyle(sheetName, cell, cell, contentBoldStyle)

							names := make([]string, 0)
							for nameIdx := 0; nameIdx < subValue.Len(); nameIdx++ {
								names = append(names, subValue.Index(nameIdx).String())
							}

							cell, _ = excelize.CoordinatesToCellName(colCount+1, startRow+1)
							_ = f.SetCellValue(sheetName, cell, strings.Join(names, " "))
							_ = f.SetCellStyle(sheetName, cell, cell, contentStyle)
						case "Size":
							cell, _ = excelize.CoordinatesToCellName(colCount, startRow+2)
							_ = f.SetCellValue(sheetName, cell, "Size")
							_ = f.SetCellStyle(sheetName, cell, cell, contentBoldStyle)
							cell, _ = excelize.CoordinatesToCellName(colCount+1, startRow+2)
							_ = f.SetCellValue(sheetName, cell, subValue.Float())
							_ = f.SetCellStyle(sheetName, cell, cell, contentStyle)
						case "Location":
							locVal := reflect.ValueOf(subValue.Interface())
							//locType := locVal.Type()
							for locIdx := 0; locIdx < locVal.NumField(); locIdx++ {
								locField := locVal.Type().Field(locIdx)
								locSubVal := locVal.Field(locIdx)
								switch locField.Name {
								case "Symbol":
									cell, _ = excelize.CoordinatesToCellName(colCount, startRow)
									endCell, _ = excelize.CoordinatesToCellName(colCount+1, startRow)
									_ = f.MergeCell(sheetName, cell, endCell)
									_ = f.SetCellValue(sheetName, cell, locSubVal.String())
									_ = f.SetCellStyle(sheetName, cell, cell, contentMiddleAlignStyle)
									break
								default:
									continue
								}
							}

						default:
							continue
						}
						rowCount = startRow + 3
					}
				}
			case "PowerInclusivePenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Power difference (for penalty)", value.Float())
			case "AlphaInclusivePenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Alpha (for penalty)", value.Float())
			default:
				continue
			}
			rowCount++
		}
	}

	return rowCount
}

func coverCraneInfo(f *excelize.File, craneInfo any, sheetName string, rowCount int, colCount int) int {
	// Add sub-header
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.MergeCell(sheetName, cell, endCell)
	_ = f.SetCellValue(sheetName, cell, "Cover in Crane's radius")
	_ = f.SetCellStyle(sheetName, cell, cell, subHeaderStyle)
	rowCount++
	val := reflect.ValueOf(craneInfo)
	typ := val.Type()
	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			switch field.Name {
			case "Cranes":
				for j := 0; j < value.Len(); j++ {
					startRow := rowCount
					elem := value.Index(j)

					for k := 0; k < elem.NumField(); k++ {
						subField := elem.Type().Field(k)
						subValue := elem.Field(k)

						switch subField.Name {
						case "BuildingName":
							cell, _ = excelize.CoordinatesToCellName(colCount, startRow+1)
							_ = f.SetCellValue(sheetName, cell, "Facilities")
							_ = f.SetCellStyle(sheetName, cell, cell, contentBoldStyle)

							names := make([]string, 0)
							for nameIdx := 0; nameIdx < subValue.Len(); nameIdx++ {
								names = append(names, subValue.Index(nameIdx).String())
							}

							cell, _ = excelize.CoordinatesToCellName(colCount+1, startRow+1)
							_ = f.SetCellValue(sheetName, cell, strings.Join(names, " "))
							_ = f.SetCellStyle(sheetName, cell, cell, contentStyle)
						case "Radius":
							cell, _ = excelize.CoordinatesToCellName(colCount, startRow+2)
							_ = f.SetCellValue(sheetName, cell, "Radius")
							_ = f.SetCellStyle(sheetName, cell, cell, contentBoldStyle)
							cell, _ = excelize.CoordinatesToCellName(colCount+1, startRow+2)
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
						rowCount = startRow + 3
					}
				}
			case "PowerCoverInCraneRadiusPenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Power difference (for penalty)", value.Float())
			case "AlphaCoverInCraneRadiusPenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Alpha (for penalty)", value.Float())
			default:
				continue
			}
			rowCount++
		}
	}

	return rowCount
}

func sizeInfo(f *excelize.File, inclusive any, sheetName string, rowCount int, colCount int) int {
	// Add sub-header
	cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
	endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
	_ = f.MergeCell(sheetName, cell, endCell)
	_ = f.SetCellValue(sheetName, cell, "Size")
	_ = f.SetCellStyle(sheetName, cell, cell, subHeaderStyle)
	rowCount++
	val := reflect.ValueOf(inclusive)
	typ := val.Type()
	// Loop through fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		// Only exported fields (unexported fields can't be accessed)
		if field.PkgPath == "" {
			switch field.Name {
			case "SmallLocations":
				cell, _ = excelize.CoordinatesToCellName(colCount, rowCount)
				_ = f.SetCellValue(sheetName, cell, "Small locations")
				_ = f.SetCellStyle(sheetName, cell, cell, contentBoldStyle)
				names := make([]string, 0)
				for nameIdx := 0; nameIdx < value.Len(); nameIdx++ {
					elem := value.Index(nameIdx)
					names = append(names, elem.String())
				}

				cell, _ = excelize.CoordinatesToCellName(colCount+1, rowCount)
				_ = f.SetCellValue(sheetName, cell, strings.Join(names, " "))
				_ = f.SetCellStyle(sheetName, cell, cell, contentStyle)
			case "LargeFacilities":
				cell, _ = excelize.CoordinatesToCellName(colCount, rowCount)
				_ = f.SetCellValue(sheetName, cell, "Large facilities")
				_ = f.SetCellStyle(sheetName, cell, cell, contentBoldStyle)
				names := make([]string, 0)
				for nameIdx := 0; nameIdx < value.Len(); nameIdx++ {
					elem := value.Index(nameIdx)
					names = append(names, elem.String())
				}
				cell, _ = excelize.CoordinatesToCellName(colCount+1, rowCount)
				_ = f.SetCellValue(sheetName, cell, strings.Join(names, " "))
				_ = f.SetCellStyle(sheetName, cell, cell, contentStyle)
			case "PowerDifferencePenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Power difference (for penalty)", value.Float())
			case "AlphaSizePenalty":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Alpha (for penalty)", value.Float())
			default:
				continue
			}
			rowCount++
		}
	}

	return rowCount
}

// Sheet 2 - Result

var locationHeader = []string{"Name", "Symbol", "x", "y", "Rotated", "Length", "Width", "Fixed"}

func generateSheet2Results(f *excelize.File, results algorithms.Result) error {
	const SheetName = "Results"

	// Starting point
	rowCount := 2
	columnCount := 2
	index, err := f.NewSheet(SheetName)
	if err != nil {
		return err
	}

	f.SetActiveSheet(index)
	err = f.SetColWidth(SheetName, "A", "A", 5)
	err = f.SetColWidth(SheetName, "B", "B", 40)
	err = f.SetColWidth(SheetName, "C", "J", 20)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(results)
	//typ := val.Type()

	resultField := val.FieldByName("Result")
	if resultField.IsValid() {
		// Iterate through the slice
		for i := 0; i < resultField.Len(); i++ {
			algResult := resultField.Index(i)
			cell, _ := excelize.CoordinatesToCellName(1, rowCount)
			_ = f.SetCellValue(SheetName, cell, fmt.Sprintf("#%d", i+1))
			_ = f.SetCellStyle(SheetName, cell, cell, contentBoldStyle)

			// Access the ValuesWithKey map
			valuesWithKey := algResult.FieldByName("ValuesWithKey")
			if valuesWithKey.IsValid() {
				//fmt.Printf("  ValuesWithKey has %d entries\n", valuesWithKey.Len())
				cell, _ = excelize.CoordinatesToCellName(columnCount, rowCount)
				_ = f.SetCellValue(SheetName, cell, "Objectives")
				_ = f.SetCellStyle(SheetName, cell, cell, headerStyle)

				// Get map keys
				mapKeys := valuesWithKey.MapKeys()
				for keyIdx, key := range mapKeys {

					cell, _ = excelize.CoordinatesToCellName(columnCount+1+keyIdx, rowCount)
					_ = f.SetCellValue(SheetName, cell, re.ReplaceAllString(key.String(), ""))
					_ = f.SetCellStyle(SheetName, cell, cell, subHeaderStyle)

					cell, _ = excelize.CoordinatesToCellName(columnCount+1+keyIdx, rowCount+1)
					_ = f.SetCellValue(SheetName, cell, valuesWithKey.MapIndex(key).Float())
					_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)

				}
				rowCount += 3
			}

			// Access the Penalty map
			penalty := algResult.FieldByName("Penalty")
			if penalty.IsValid() {
				//fmt.Printf("  Penalty has %d entries\n", penalty.Len())
				cell, _ = excelize.CoordinatesToCellName(columnCount, rowCount)
				_ = f.SetCellValue(SheetName, cell, "Penalty Constraints")
				_ = f.SetCellStyle(SheetName, cell, cell, headerStyle)
				// Get map keys
				mapKeys := penalty.MapKeys()
				for keyIdx, key := range mapKeys {
					//fmt.Printf("  Key: %s, Value: %f\n",
					//	key.String(), penalty.MapIndex(key).Float())
					cell, _ = excelize.CoordinatesToCellName(columnCount+1+keyIdx, rowCount)
					_ = f.SetCellValue(SheetName, cell, key.String())
					_ = f.SetCellStyle(SheetName, cell, cell, subHeaderStyle)

					cell, _ = excelize.CoordinatesToCellName(columnCount+1+keyIdx, rowCount+1)
					_ = f.SetCellValue(SheetName, cell, penalty.MapIndex(key).Float())
					_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)

				}
				rowCount += 3
			}

			// Access the MapLocations map
			sliceLocations := algResult.FieldByName("SliceLocations")
			if sliceLocations.IsValid() {
				//fmt.Printf("  SliceLocations has %d entries\n", sliceLocations.Len())
				for headerIdx, header := range locationHeader {
					cell, _ = excelize.CoordinatesToCellName(columnCount+headerIdx, rowCount)
					_ = f.SetCellValue(SheetName, cell, header)
					_ = f.SetCellStyle(SheetName, cell, cell, headerStyle)
				}
				rowCount++

				for idx := 0; idx < sliceLocations.Len(); idx++ {
					// Access the Location value
					locValue := sliceLocations.Index(idx)
					//fmt.Printf("    Location: %+v\n", locValue.Interface())
					cell, _ = excelize.CoordinatesToCellName(columnCount, rowCount)
					_ = f.SetCellValue(SheetName, cell, locValue.FieldByName("Name").String())
					_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)

					cell, _ = excelize.CoordinatesToCellName(columnCount+1, rowCount)
					_ = f.SetCellValue(SheetName, cell, locValue.FieldByName("Symbol").String())
					_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)

					coordField := locValue.FieldByName("Coordinate")
					if coordField.IsValid() {
						x := coordField.FieldByName("X").Float()
						y := coordField.FieldByName("Y").Float()
						//fmt.Printf("    Coordinates: (%f, %f)\n", x, y)

						cell, _ = excelize.CoordinatesToCellName(columnCount+2, rowCount)
						_ = f.SetCellValue(SheetName, cell, x)
						_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)

						cell, _ = excelize.CoordinatesToCellName(columnCount+3, rowCount)
						_ = f.SetCellValue(SheetName, cell, y)
						_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)
					}

					cell, _ = excelize.CoordinatesToCellName(columnCount+4, rowCount)
					_ = f.SetCellValue(SheetName, cell, locValue.FieldByName("Rotation").Bool())
					_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)

					cell, _ = excelize.CoordinatesToCellName(columnCount+5, rowCount)
					_ = f.SetCellValue(SheetName, cell, locValue.FieldByName("Length").Float())
					_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)

					cell, _ = excelize.CoordinatesToCellName(columnCount+6, rowCount)
					_ = f.SetCellValue(SheetName, cell, locValue.FieldByName("Width").Float())
					_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)

					cell, _ = excelize.CoordinatesToCellName(columnCount+7, rowCount)
					_ = f.SetCellValue(SheetName, cell, locValue.FieldByName("IsFixed").Bool())
					_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)

					rowCount++
				}
			}
			rowCount += 2
		}
	}

	return nil
}

var locationHeaderPredetermined = []string{"Symbol", "Is Located At"}

func generateSheet2ResultsPredetermined(f *excelize.File, results algorithms.Result) error {
	const SheetName = "Results"
	// Starting point
	rowCount := 2
	columnCount := 2
	index, err := f.NewSheet(SheetName)
	if err != nil {
		return err
	}

	f.SetActiveSheet(index)
	err = f.SetColWidth(SheetName, "A", "A", 5)
	err = f.SetColWidth(SheetName, "B", "B", 40)
	err = f.SetColWidth(SheetName, "C", "C", 20)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(results)
	//typ := val.Type()

	resultField := val.FieldByName("Result")
	if resultField.IsValid() {
		// Iterate through the slice
		for i := 0; i < resultField.Len(); i++ {
			algResult := resultField.Index(i)
			cell, _ := excelize.CoordinatesToCellName(1, rowCount)
			_ = f.SetCellValue(SheetName, cell, fmt.Sprintf("#%d", i+1))
			_ = f.SetCellStyle(SheetName, cell, cell, contentBoldStyle)

			// Access the ValuesWithKey map
			valuesWithKey := algResult.FieldByName("ValuesWithKey")
			if valuesWithKey.IsValid() {
				//fmt.Printf("  ValuesWithKey has %d entries\n", valuesWithKey.Len())
				cell, _ = excelize.CoordinatesToCellName(columnCount, rowCount)
				_ = f.SetCellValue(SheetName, cell, "Objectives")
				_ = f.SetCellStyle(SheetName, cell, cell, headerStyle)

				// Get map keys
				mapKeys := valuesWithKey.MapKeys()
				for keyIdx, key := range mapKeys {

					cell, _ = excelize.CoordinatesToCellName(columnCount+1+keyIdx, rowCount)
					_ = f.SetCellValue(SheetName, cell, re.ReplaceAllString(key.String(), ""))
					_ = f.SetCellStyle(SheetName, cell, cell, subHeaderStyle)

					cell, _ = excelize.CoordinatesToCellName(columnCount+1+keyIdx, rowCount+1)
					_ = f.SetCellValue(SheetName, cell, valuesWithKey.MapIndex(key).Float())
					_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)

				}
				rowCount += 3
			}

			// Access the Penalty map
			penalty := algResult.FieldByName("Penalty")
			if penalty.IsValid() {
				//fmt.Printf("  Penalty has %d entries\n", penalty.Len())
				cell, _ = excelize.CoordinatesToCellName(columnCount, rowCount)
				_ = f.SetCellValue(SheetName, cell, "Penalty Constraints")
				_ = f.SetCellStyle(SheetName, cell, cell, headerStyle)
				// Get map keys
				mapKeys := penalty.MapKeys()
				for keyIdx, key := range mapKeys {
					//fmt.Printf("  Key: %s, Value: %f\n",
					//	key.String(), penalty.MapIndex(key).Float())
					cell, _ = excelize.CoordinatesToCellName(columnCount+1+keyIdx, rowCount)
					_ = f.SetCellValue(SheetName, cell, key.String())
					_ = f.SetCellStyle(SheetName, cell, cell, subHeaderStyle)

					cell, _ = excelize.CoordinatesToCellName(columnCount+1+keyIdx, rowCount+1)
					_ = f.SetCellValue(SheetName, cell, penalty.MapIndex(key).Float())
					_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)

				}
				rowCount += 3
			}

			// Access the MapLocations map
			sliceLocations := algResult.FieldByName("SliceLocations")
			if sliceLocations.IsValid() {
				//fmt.Printf("  SliceLocations has %d entries\n", sliceLocations.Len())
				for headerIdx, header := range locationHeaderPredetermined {
					cell, _ = excelize.CoordinatesToCellName(columnCount+headerIdx, rowCount)
					_ = f.SetCellValue(SheetName, cell, header)
					_ = f.SetCellStyle(SheetName, cell, cell, headerStyle)
				}
				rowCount++

				for idx := 0; idx < sliceLocations.Len(); idx++ {
					// Access the Location value
					locValue := sliceLocations.Index(idx)
					//fmt.Printf("    Location: %+v\n", locValue.Interface())
					cell, _ = excelize.CoordinatesToCellName(columnCount, rowCount)
					_ = f.SetCellValue(SheetName, cell, locValue.FieldByName("Symbol").String())
					_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)

					cell, _ = excelize.CoordinatesToCellName(columnCount+1, rowCount)
					_ = f.SetCellValue(SheetName, cell, locValue.FieldByName("IsLocatedAt").String())
					_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)

					rowCount++
				}
			}
			rowCount += 2
		}
	}

	return nil
}

// Sheet 3 - Pareto

func generateSheet3Graph(f *excelize.File, results algorithms.Result, numberOfObjectives int) error {

	var SheetName string
	if numberOfObjectives > 1 {
		// add pareto
		SheetName = "Pareto"
	} else if numberOfObjectives == 1 {
		// add convergence
		SheetName = "Convergence"
	} else {
		return nil
	}

	// Starting point
	rowCount, startRow := 2, 2
	columnCount := 2
	index, err := f.NewSheet(SheetName)
	if err != nil {
		return err
	}

	f.SetActiveSheet(index)
	err = f.SetColWidth(SheetName, "A", "A", 5)
	err = f.SetColWidth(SheetName, "B", "B", 40)
	err = f.SetColWidth(SheetName, "C", "J", 20)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(results)
	//typ := val.Type()

	// for pareto - take Result field
	if numberOfObjectives > 1 {

		resultField := val.FieldByName("Result")
		if resultField.IsValid() {
			headersVal := resultField.Index(0)
			header := headersVal.FieldByName("Key")
			if header.IsValid() {
				// Iterate through the slice
				for i := 0; i < header.Len(); i++ {
					headerVal := header.Index(i)
					if headerVal.IsValid() {
						cell, _ := excelize.CoordinatesToCellName(columnCount+i, 1)
						_ = f.SetCellValue(SheetName, cell, headerVal.String())
						_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)
					}
				}
			}

			// Iterate through the slice
			numberOfResults := resultField.Len()
			for i := 0; i < numberOfResults; i++ {
				algResult := resultField.Index(i)
				cell, _ := excelize.CoordinatesToCellName(1, rowCount)
				_ = f.SetCellValue(SheetName, cell, fmt.Sprintf("#%d", i+1))
				_ = f.SetCellStyle(SheetName, cell, cell, contentBoldStyle)

				// Access the Value Slice field
				valuesSlice := algResult.FieldByName("Value")
				if valuesSlice.IsValid() {
					for idx := 0; idx < valuesSlice.Len(); idx++ {
						// Access the Location value
						locValue := valuesSlice.Index(idx)
						//fmt.Printf("    Location: %+v\n", locValue.Interface())
						cell, _ = excelize.CoordinatesToCellName(columnCount+idx, rowCount)
						_ = f.SetCellValue(SheetName, cell, locValue.Float())
						_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)
					}

					rowCount++
				}
			}

			for i := 0; i < numberOfObjectives; i++ {
				// add min and max rows
				cell, _ := excelize.CoordinatesToCellName(1, rowCount)
				_ = f.SetCellValue(SheetName, cell, "Min")
				_ = f.SetCellStyle(SheetName, cell, cell, contentBoldStyle)

				// Add MIN formula for this objective column
				minCell, _ := excelize.CoordinatesToCellName(columnCount+i, rowCount)
				startCell, _ := excelize.CoordinatesToCellName(columnCount+i, startRow)
				endCell, _ := excelize.CoordinatesToCellName(columnCount+i, startRow+numberOfResults-1)
				minFormula := fmt.Sprintf("=MIN(%s:%s)", startCell, endCell)
				_ = f.SetCellFormula(SheetName, minCell, minFormula)
				_ = f.SetCellStyle(SheetName, minCell, minCell, contentStyle)

				cell, _ = excelize.CoordinatesToCellName(1, rowCount+1)
				_ = f.SetCellValue(SheetName, cell, "Max")
				_ = f.SetCellStyle(SheetName, cell, cell, contentBoldStyle)

				// Add MAX formula for this objective column
				maxCell, _ := excelize.CoordinatesToCellName(columnCount+i, rowCount+1)
				maxFormula := fmt.Sprintf("=MAX(%s:%s)", startCell, endCell)
				_ = f.SetCellFormula(SheetName, maxCell, maxFormula)
				_ = f.SetCellStyle(SheetName, maxCell, maxCell, contentStyle)
			}

			// Update rowCount to account for the min and max rows
			rowCount += 2

		}
	} else {
		resultField := val.FieldByName("Convergence")
		headerField := val.FieldByName("Result")
		if headerField.IsValid() {
			headersVal := headerField.Index(0)
			header := headersVal.FieldByName("Key")
			if header.IsValid() {
				// Iterate through the slice
				for i := 0; i < header.Len(); i++ {
					headerVal := header.Index(i)
					if headerVal.IsValid() {
						cell, _ := excelize.CoordinatesToCellName(columnCount+i, 1)
						_ = f.SetCellValue(SheetName, cell, headerVal.String())
						_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)
					}
				}
			}
		}

		if resultField.IsValid() {
			// Iterate through the slice
			for i := 0; i < resultField.Len(); i++ {
				locValue := resultField.Index(i)
				cell, _ := excelize.CoordinatesToCellName(1, rowCount)
				_ = f.SetCellValue(SheetName, cell, fmt.Sprintf("#%d", i+1))
				_ = f.SetCellStyle(SheetName, cell, cell, contentBoldStyle)

				cell, _ = excelize.CoordinatesToCellName(columnCount, rowCount)
				_ = f.SetCellValue(SheetName, cell, locValue.Float())
				_ = f.SetCellStyle(SheetName, cell, cell, contentStyle)
				rowCount++
			}
		}
	}

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
