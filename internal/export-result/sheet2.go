// Package export_result provides functionality for exporting optimization results to Excel files.
package export_result

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/algorithms"
	"reflect"
)

// Sheet 2 - Result

// generateSheet2Results generates the results sheet with optimization results
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

// generateSheet2ResultsPredetermined generates the results sheet for predetermined construction layout
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
