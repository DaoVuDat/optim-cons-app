// Package export_result provides functionality for exporting optimization results to Excel files.
package export_result_old

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"golang-moaha-construction/internal/algorithms"
	"reflect"
)

// Sheet 3 - Pareto/Convergence

// generateSheet3Graph generates the pareto or convergence sheet based on the number of objectives
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
