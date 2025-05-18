// Package export_result provides functionality for exporting optimization results to Excel files.
package export_result

import (
	"fmt"
	"golang-moaha-construction/internal/algorithms"
	"reflect"
	"strings"

	"github.com/xuri/excelize/v2"
)

// Sheet 1 - Summary

// generateSheet1Info generates the summary sheet with algorithm, problem, objectives, and constraints information
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

// sectionAlgorithm adds the algorithm section to the summary sheet
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
				// Check if it's a slice (like in NSGA-II) or an integer
				if value.Kind() == reflect.Int || value.Kind() == reflect.Int64 {
					writeContentWithValue(f, colCount, rowCount, sheetName, "Population", value.Int())
				}
			case "PopulationSize":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Population size", value.Int())
			case "ArchiveSize":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Archive size", value.Int())
			case "NumberOfIter":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Number of iterations", value.Int())
			case "MaxIterations":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Maximum iterations", value.Int())
			case "Generation":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Generation", value.Int())
			case "CrossoverRate":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Crossover rate", value.Float())
			case "MutationRate":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Mutation rate", value.Float())
			case "MutationStrength":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Mutation strength", value.Float())
			case "Sigma":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Sigma", value.Float())
			case "AParam":
				writeContentWithValue(f, colCount, rowCount, sheetName, "A parameter", value.Float())
			case "NumberOfGrids":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Number of grids", value.Int())
			case "MaxVelocityInfo":
				writeContentWithValue(f, colCount, rowCount, sheetName, "Max Velocity", value.Float())
			case "C1":
				writeContentWithValue(f, colCount, rowCount, sheetName, "C1", value.Float())
			case "C2":
				writeContentWithValue(f, colCount, rowCount, sheetName, "C2", value.Float())
			case "W":
				writeContentWithValue(f, colCount, rowCount, sheetName, "W", value.Float())
			default:
				continue
			}
			rowCount++
		}
	}

	return rowCount + 2
}

// sectionProblem adds the problem section to the summary sheet
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

// sectionObjectives adds the objectives section to the summary sheet
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

// riskInfo adds risk objective information to the summary sheet
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

// hoistingInfo adds hoisting objective information to the summary sheet
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
			case "Buildings":
				// Export each building's name and its properties except NumberOfFloors and FloorHeight
				cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
				endCell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
				_ = f.MergeCell(sheetName, cell, endCell)
				_ = f.SetCellValue(sheetName, cell, "Buildings")
				_ = f.SetCellStyle(sheetName, cell, cell, contentMiddleAlignStyle)
				rowCount++
				for _, key := range value.MapKeys() {
					buildingName := key.String()
					buildingVal := value.MapIndex(key)
					cell, _ := excelize.CoordinatesToCellName(colCount, rowCount)
					_ = f.SetCellValue(sheetName, cell, buildingName)
					_ = f.SetCellStyle(sheetName, cell, cell, contentBoldStyle)
					// Export only fields except NumberOfFloors and FloorHeight
					bVal := buildingVal
					bType := bVal.Type()
					for j := 0; j < bVal.NumField(); j++ {
						bField := bType.Field(j)
						bFieldValue := bVal.Field(j)
						if bField.Name != "NumberOfFloors" && bField.Name != "FloorHeight" {
							continue
						}
						rowCount++
						cell, _ := excelize.CoordinatesToCellName(colCount+1, rowCount)
						_ = f.SetCellValue(sheetName, cell, bField.Name+": "+toString(bFieldValue))
						_ = f.SetCellStyle(sheetName, cell, cell, contentStyle)
					}
				}
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
						case "FilePath":
							cell, _ = excelize.CoordinatesToCellName(colCount, startRow+1)
							_ = f.SetCellValue(sheetName, cell, "Hoisting file path")
							_ = f.SetCellStyle(sheetName, cell, cell, contentBoldStyle)
							cell, _ = excelize.CoordinatesToCellName(colCount+1, startRow+1)
							_ = f.SetCellValue(sheetName, cell, subValue.String())
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
						rowCount = startRow + 2
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

// Helper to convert reflect.Value to string for export
func toString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", v.Float())
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool())
	default:
		return ""
	}
}

// safetyInfo adds safety objective information to the summary sheet
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

// safetyHazardInfo adds safety hazard objective information to the summary sheet
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

// transportCostInfo adds transport cost objective information to the summary sheet
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

// constructionCostInfo adds construction cost objective information to the summary sheet
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

// sectionConstraints adds the constraints section to the summary sheet
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

// outOfBoundaryInfo adds out of boundary constraint information to the summary sheet
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

// overlapInfo adds overlap constraint information to the summary sheet
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

// inclusiveZoneInfo adds inclusive zone constraint information to the summary sheet
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

// coverCraneInfo adds cover in crane radius constraint information to the summary sheet
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

// sizeInfo adds size constraint information to the summary sheet
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
