package data

import (
	"testing"
)

func TestCreateTwoDimensionalMatrix(t *testing.T) {
	itemNames := []string{"A", "B", "C", "D", "E"}

	matrix := CreateTwoDimensionalMatrix(itemNames)

	if matrix.GetNumberOfItems() != 5 {
		t.Errorf("expected number of items to be 5, got %d", matrix.GetNumberOfItems())
	}

	if val, err := matrix.GetNameFromIdx(0); err != nil {
		t.Errorf("expected idx to name to be A, got %s", val)
	}
	if val, err := matrix.GetIdxFromName("A"); err != nil {
		t.Errorf("expected name to idx to be 0, got %d", val)
	}

	if val, err := matrix.GetNameFromIdx(4); err != nil {
		t.Errorf("expected idx to name to be E, got %s", val)
	}
	if val, err := matrix.GetIdxFromName("E"); err != nil {
		t.Errorf("expected name to idx to be 4, got %d", val)
	}
}

func TestTwoDimensionalMatrix_SetCellValueFromIdx(t *testing.T) {
	itemNames := []string{"A", "B", "C", "D", "E"}

	matrix := CreateTwoDimensionalMatrix(itemNames)

	if err := matrix.SetCellValueFromIdx(0, 0, 1); err != nil {
		t.Errorf("expected to set value to 1, got %f", matrix.Matrix[0][0])
	}

	if err := matrix.SetCellValueFromIdx(0, 1, 2); err != nil {
		t.Errorf("expected to set value to 2, got %f", matrix.Matrix[0][1])
	}

	if err := matrix.SetCellValueFromIdx(1, 2, 3); err != nil {
		t.Errorf("expected to set value to 3, got %f", matrix.Matrix[1][2])
	}

	if err := matrix.SetCellValueFromIdx(2, 3, 4); err != nil {
		t.Errorf("expected to set value to 4, got %f", matrix.Matrix[2][3])
	}

	if err := matrix.SetCellValueFromIdx(3, 4, 5); err != nil {
		t.Errorf("expected to set value to 5, got %f", matrix.Matrix[3][4])
	}

	if err := matrix.SetCellValueFromIdx(4, 4, 6); err != nil {
		t.Errorf("expected to set value to 6, got %f", matrix.Matrix[4][4])
	}
}

func TestTwoDimensionalMatrix_SetCellValueFromNames(t *testing.T) {
	itemNames := []string{"A", "B", "C", "D", "E"}

	matrix := CreateTwoDimensionalMatrix(itemNames)

	if err := matrix.SetCellValueFromNames("A", "A", 1); err != nil {
		t.Errorf("expected to set value to 1, got %f", matrix.Matrix[0][0])
	}

	if err := matrix.SetCellValueFromNames("A", "B", 2); err != nil {
		t.Errorf("expected to set value to 2, got %f", matrix.Matrix[0][1])
	}

	if err := matrix.SetCellValueFromNames("B", "C", 3); err != nil {
		t.Errorf("expected to set value to 3, got %f", matrix.Matrix[1][2])
	}

	if err := matrix.SetCellValueFromNames("C", "D", 4); err != nil {
		t.Errorf("expected to set value to 4, got %f", matrix.Matrix[2][3])
	}

	if err := matrix.SetCellValueFromNames("D", "E", 5); err != nil {
		t.Errorf("expected to set value to 5, got %f", matrix.Matrix[3][4])
	}

	if err := matrix.SetCellValueFromNames("E", "E", 6); err != nil {
		t.Errorf("expected to set value to 6, got %f", matrix.Matrix[4][4])
	}
}

func TestTwoDimensionalMatrix_GetCellValueFromIdx(t *testing.T) {
	itemNames := []string{"A", "B", "C", "D", "E"}

	matrix := CreateTwoDimensionalMatrix(itemNames)

	_ = matrix.SetCellValueFromIdx(0, 0, 1)
	_ = matrix.SetCellValueFromIdx(0, 1, 2)
	_ = matrix.SetCellValueFromIdx(1, 2, 3)
	_ = matrix.SetCellValueFromIdx(2, 3, 4)
	_ = matrix.SetCellValueFromIdx(3, 4, 5)
	_ = matrix.SetCellValueFromIdx(4, 4, 6)

	if val, err := matrix.GetCellValueFromIdx(0, 0); err != nil {
		t.Errorf("expected to get value 1, got %f", val)
	}

	if val, err := matrix.GetCellValueFromIdx(0, 1); err != nil {
		t.Errorf("expected to get value 2, got %f", val)
	}

	if val, err := matrix.GetCellValueFromIdx(1, 2); err != nil {
		t.Errorf("expected to get value 3, got %f", val)
	}

	if val, err := matrix.GetCellValueFromIdx(2, 3); err != nil {
		t.Errorf("expected to get value 4, got %f", val)
	}

	if val, err := matrix.GetCellValueFromIdx(3, 4); err != nil {
		t.Errorf("expected to get value 5, got %f", val)
	}

	if val, err := matrix.GetCellValueFromIdx(4, 4); err != nil {
		t.Errorf("expected to get value 6, got %f", val)
	}
}

func TestTwoDimensionalMatrix_GetCellValueFromNames(t *testing.T) {
	itemNames := []string{"A", "B", "C", "D", "E"}

	matrix := CreateTwoDimensionalMatrix(itemNames)

	_ = matrix.SetCellValueFromIdx(0, 0, 1)
	_ = matrix.SetCellValueFromIdx(0, 1, 2)
	_ = matrix.SetCellValueFromIdx(1, 2, 3)
	_ = matrix.SetCellValueFromIdx(2, 3, 4)
	_ = matrix.SetCellValueFromIdx(3, 4, 5)
	_ = matrix.SetCellValueFromIdx(4, 4, 6)

	if val, err := matrix.GetCellValueFromNames("A", "A"); err != nil {
		t.Errorf("expected to get value 1, got %f", val)
	}

	if val, err := matrix.GetCellValueFromNames("A", "B"); err != nil {
		t.Errorf("expected to get value 2, got %f", val)
	}

	if val, err := matrix.GetCellValueFromNames("B", "C"); err != nil {
		t.Errorf("expected to get value 3, got %f", val)
	}

	if val, err := matrix.GetCellValueFromNames("C", "D"); err != nil {
		t.Errorf("expected to get value 4, got %f", val)
	}

	if val, err := matrix.GetCellValueFromNames("D", "E"); err != nil {
		t.Errorf("expected to get value 5, got %f", val)
	}

	if val, err := matrix.GetCellValueFromNames("E", "E"); err != nil {
		t.Errorf("expected to get value 6, got %f", val)
	}
}

// Test for RectangleMatrix
func TestCreateRectangleMatrix(t *testing.T) {
	rowNames := []string{"TF1", "TF2", "TF3", "TF4", "TF5"}
	colNames := []string{"L1", "L2", "L3", "L4", "L5", "L6"}

	matrix := CreateRectangleMatrix(rowNames, colNames)

	if matrix.GetNumberOfRows() != 5 {
		t.Errorf("expected number of rows to be 5, got %d", matrix.GetNumberOfRows())
	}

	if matrix.GetNumberOfColumns() != 6 {
		t.Errorf("expected number of columns to be 6, got %d", matrix.GetNumberOfColumns())
	}

	if val, err := matrix.GetRowNameFromIdx(0); err != nil {
		t.Errorf("expected idx to name to be TF1, got %s", val)
	}
	if val, err := matrix.GetRowIdxFromName("TF1"); err != nil {
		t.Errorf("expected name to idx to be 0, got %d", val)
	}

	if val, err := matrix.GetRowNameFromIdx(4); err != nil {
		t.Errorf("expected idx to name to be TF5, got %s", val)
	}
	if val, err := matrix.GetRowIdxFromName("TF5"); err != nil {
		t.Errorf("expected name to idx to be 4, got %d", val)
	}

	if val, err := matrix.GetColumnNameFromIdx(0); err != nil {
		t.Errorf("expected idx to name to be L1, got %s", val)
	}
	if val, err := matrix.GetColumnIdxFromName("L1"); err != nil {
		t.Errorf("expected name to idx to be 0, got %d", val)
	}

	if val, err := matrix.GetColumnNameFromIdx(5); err != nil {
		t.Errorf("expected idx to name to be L6, got %s", val)
	}
	if val, err := matrix.GetColumnIdxFromName("L6"); err != nil {
		t.Errorf("expected name to idx to be 5, got %d", val)
	}
}

func TestRectangleMatrix_SetCellValueFromIdx(t *testing.T) {
	rowNames := []string{"TF1", "TF2", "TF3", "TF4", "TF5"}
	colNames := []string{"L1", "L2", "L3", "L4", "L5", "L6"}

	matrix := CreateRectangleMatrix(rowNames, colNames)

	if err := matrix.SetCellValueFromIdx(0, 0, 1); err != nil {
		t.Errorf("expected to set value to 1, got %f", matrix.Matrix[0][0])
	}

	if err := matrix.SetCellValueFromIdx(0, 1, 2); err != nil {
		t.Errorf("expected to set value to 2, got %f", matrix.Matrix[0][1])
	}

	if err := matrix.SetCellValueFromIdx(1, 2, 3); err != nil {
		t.Errorf("expected to set value to 3, got %f", matrix.Matrix[1][2])
	}

	if err := matrix.SetCellValueFromIdx(2, 3, 4); err != nil {
		t.Errorf("expected to set value to 4, got %f", matrix.Matrix[2][3])
	}

	if err := matrix.SetCellValueFromIdx(3, 4, 5); err != nil {
		t.Errorf("expected to set value to 5, got %f", matrix.Matrix[3][4])
	}

	if err := matrix.SetCellValueFromIdx(4, 4, 6); err != nil {
		t.Errorf("expected to set value to 6, got %f", matrix.Matrix[4][4])
	}
}

func TestRectangleMatrix_SetCellValueFromNames(t *testing.T) {
	rowNames := []string{"TF1", "TF2", "TF3", "TF4", "TF5"}
	colNames := []string{"L1", "L2", "L3", "L4", "L5", "L6"}

	matrix := CreateRectangleMatrix(rowNames, colNames)

	if err := matrix.SetCellValueFromNames("TF1", "L1", 1); err != nil {
		t.Errorf("expected to set value to 1, got %f", matrix.Matrix[0][0])
	}

	if err := matrix.SetCellValueFromNames("TF1", "L2", 2); err != nil {
		t.Errorf("expected to set value to 2, got %f", matrix.Matrix[0][1])
	}

	if err := matrix.SetCellValueFromNames("TF2", "L3", 3); err != nil {
		t.Errorf("expected to set value to 3, got %f", matrix.Matrix[1][2])
	}

	if err := matrix.SetCellValueFromNames("TF3", "L4", 4); err != nil {
		t.Errorf("expected to set value to 4, got %f", matrix.Matrix[2][3])
	}

	if err := matrix.SetCellValueFromNames("TF4", "L5", 5); err != nil {
		t.Errorf("expected to set value to 5, got %f", matrix.Matrix[3][4])
	}

	if err := matrix.SetCellValueFromNames("TF5", "L6", 6); err != nil {
		t.Errorf("expected to set value to 6, got %f", matrix.Matrix[4][4])
	}
}

func TestRectangleMatrix_GetCellValueFromIdx(t *testing.T) {
	rowNames := []string{"TF1", "TF2", "TF3", "TF4", "TF5"}
	colNames := []string{"L1", "L2", "L3", "L4", "L5", "L6"}

	matrix := CreateRectangleMatrix(rowNames, colNames)

	_ = matrix.SetCellValueFromIdx(0, 0, 1)
	_ = matrix.SetCellValueFromIdx(0, 1, 2)
	_ = matrix.SetCellValueFromIdx(1, 2, 3)
	_ = matrix.SetCellValueFromIdx(2, 3, 4)
	_ = matrix.SetCellValueFromIdx(3, 4, 5)
	_ = matrix.SetCellValueFromIdx(4, 5, 6)

	if val, err := matrix.GetCellValueFromIdx(0, 0); err != nil {
		t.Errorf("expected to get value 1, got %f", val)
	}

	if val, err := matrix.GetCellValueFromIdx(0, 1); err != nil {
		t.Errorf("expected to get value 2, got %f", val)
	}

	if val, err := matrix.GetCellValueFromIdx(1, 2); err != nil {
		t.Errorf("expected to get value 3, got %f", val)
	}

	if val, err := matrix.GetCellValueFromIdx(2, 3); err != nil {
		t.Errorf("expected to get value 4, got %f", val)
	}

	if val, err := matrix.GetCellValueFromIdx(3, 4); err != nil {
		t.Errorf("expected to get value 5, got %f", val)
	}

	if val, err := matrix.GetCellValueFromIdx(4, 4); err != nil {
		t.Errorf("expected to get value 6, got %f", val)
	}
}

func TestRectangleMatrix_GetCellValueFromNames(t *testing.T) {
	rowNames := []string{"TF1", "TF2", "TF3", "TF4", "TF5"}
	colNames := []string{"L1", "L2", "L3", "L4", "L5", "L6"}

	matrix := CreateRectangleMatrix(rowNames, colNames)

	_ = matrix.SetCellValueFromIdx(0, 0, 1)
	_ = matrix.SetCellValueFromIdx(0, 1, 2)
	_ = matrix.SetCellValueFromIdx(1, 2, 3)
	_ = matrix.SetCellValueFromIdx(2, 3, 4)
	_ = matrix.SetCellValueFromIdx(3, 4, 5)
	_ = matrix.SetCellValueFromIdx(4, 5, 6)

	if val, err := matrix.GetCellValueFromNames("TF1", "L1"); err != nil {
		t.Errorf("expected to get value 1, got %f", val)
	}

	if val, err := matrix.GetCellValueFromNames("TF1", "L2"); err != nil {
		t.Errorf("expected to get value 2, got %f", val)
	}

	if val, err := matrix.GetCellValueFromNames("TF2", "L3"); err != nil {
		t.Errorf("expected to get value 3, got %f", val)
	}

	if val, err := matrix.GetCellValueFromNames("TF3", "L4"); err != nil {
		t.Errorf("expected to get value 4, got %f", val)
	}

	if val, err := matrix.GetCellValueFromNames("TF4", "L5"); err != nil {
		t.Errorf("expected to get value 5, got %f", val)
	}

	if val, err := matrix.GetCellValueFromNames("TF5", "L6"); err != nil {
		t.Errorf("expected to get value 6, got %f", val)
	}
}
