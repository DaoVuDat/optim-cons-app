package data

import (
	"errors"
)

// TwoDimensionalMatrix is (n*n) matrix
type TwoDimensionalMatrix struct {
	Matrix        [][]float64
	NameToIdx     map[string]int
	IdxToName     map[int]string
	NumberOfItems int
}

func CreateTwoDimensionalMatrix(itemNames []string) TwoDimensionalMatrix {

	numberOfItems := len(itemNames)

	matrix := make([][]float64, numberOfItems)
	nameToIdx := make(map[string]int, numberOfItems)
	idxToName := make(map[int]string, numberOfItems)

	for i := 0; i < numberOfItems; i++ {
		matrix[i] = make([]float64, numberOfItems)
		nameToIdx[itemNames[i]] = i
		idxToName[i] = itemNames[i]
	}

	return TwoDimensionalMatrix{
		Matrix:        matrix,
		NameToIdx:     nameToIdx,
		IdxToName:     idxToName,
		NumberOfItems: len(matrix),
	}
}

func (m *TwoDimensionalMatrix) GetMatrix() [][]float64 {
	return m.Matrix
}

func (m *TwoDimensionalMatrix) SetCellValueFromIdx(i, j int, value float64) error {
	if i < 0 || i >= m.NumberOfItems {
		return errors.New("index out of bounds")
	}
	if j < 0 || j >= m.NumberOfItems {
		return errors.New("index out of bounds")
	}

	m.Matrix[i][j] = value
	return nil
}

func (m *TwoDimensionalMatrix) SetCellValueFromNames(nameItemI, nameItemJ string, value float64) error {
	idxI, err := m.GetIdxFromName(nameItemI)
	if err != nil {
		return err
	}
	idxJ, err := m.GetIdxFromName(nameItemJ)
	if err != nil {
		return err
	}
	return m.SetCellValueFromIdx(idxI, idxJ, value)
}

func (m *TwoDimensionalMatrix) GetCellValueFromIdx(i, j int) (float64, error) {
	if i < 0 || i >= m.NumberOfItems {
		return 0, errors.New("index out of bounds")
	}
	if j < 0 || j >= m.NumberOfItems {
		return 0, errors.New("index out of bounds")
	}

	return m.Matrix[i][j], nil
}

func (m *TwoDimensionalMatrix) GetCellValueFromNames(nameItemI, nameItemJ string) (float64, error) {
	idxI, err := m.GetIdxFromName(nameItemI)
	if err != nil {
		return 0, err
	}
	idxJ, err := m.GetIdxFromName(nameItemJ)
	if err != nil {
		return 0, err
	}
	return m.GetCellValueFromIdx(idxI, idxJ)
}

func (m *TwoDimensionalMatrix) GetNameFromIdx(idx int) (string, error) {
	val, ok := m.IdxToName[idx]
	if !ok {
		return "", errors.New("index not found")
	}

	return val, nil
}

func (m *TwoDimensionalMatrix) GetIdxFromName(nameItem string) (int, error) {
	val, ok := m.NameToIdx[nameItem]
	
	if !ok {
		return -1, errors.New("name not found")
	}

	return val, nil
}

func (m *TwoDimensionalMatrix) GetNumberOfItems() int {
	return m.NumberOfItems
}

// RectangleMatrix is (n*m) matrix
type RectangleMatrix struct {
	Matrix           [][]float64
	NameToIdxRows    map[string]int
	IdxToNameRows    map[int]string
	NameToIdxColumns map[string]int
	IdxToNameColumns map[int]string
	Columns          int
	Rows             int
}

func CreateRectangleMatrix(rowsName []string, columnsName []string) RectangleMatrix {

	rows := len(rowsName)
	columns := len(columnsName)

	matrix := make([][]float64, rows)
	nameToIdxRows := make(map[string]int, rows)
	idxToNameRows := make(map[int]string, rows)

	nameToIdxColumns := make(map[string]int, columns)
	idxToNameColumns := make(map[int]string, columns)

	for i := 0; i < rows; i++ {
		matrix[i] = make([]float64, columns)
		nameToIdxRows[rowsName[i]] = i
		idxToNameRows[i] = rowsName[i]
	}

	for i := 0; i < columns; i++ {
		nameToIdxColumns[columnsName[i]] = i
		idxToNameColumns[i] = columnsName[i]
	}

	return RectangleMatrix{
		Matrix:           matrix,
		NameToIdxRows:    nameToIdxRows,
		IdxToNameRows:    idxToNameRows,
		NameToIdxColumns: nameToIdxColumns,
		IdxToNameColumns: idxToNameColumns,
		Columns:          columns,
		Rows:             rows,
	}
}

func (m *RectangleMatrix) GetMatrix() [][]float64 {
	return m.Matrix
}

func (m *RectangleMatrix) SetCellValueFromIdx(i, j int, value float64) error {
	if i < 0 || i >= m.Rows {
		return errors.New("index out of bounds")
	}
	if j < 0 || j >= m.Columns {
		return errors.New("index out of bounds")
	}

	m.Matrix[i][j] = value
	return nil
}

func (m *RectangleMatrix) SetCellValueFromNames(nameItemI, nameItemJ string, value float64) error {
	idxI, err := m.GetRowIdxFromName(nameItemI)
	if err != nil {
		return err
	}
	idxJ, err := m.GetColumnIdxFromName(nameItemJ)
	if err != nil {
		return err
	}
	return m.SetCellValueFromIdx(idxI, idxJ, value)
}

func (m *RectangleMatrix) GetCellValueFromIdx(i, j int) (float64, error) {
	if i < 0 || i >= m.Rows {
		return 0, errors.New("index out of bounds")
	}
	if j < 0 || j >= m.Columns {
		return 0, errors.New("index out of bounds")
	}

	return m.Matrix[i][j], nil
}

func (m *RectangleMatrix) GetCellValueFromNames(nameItemI, nameItemJ string) (float64, error) {
	idxI, err := m.GetRowIdxFromName(nameItemI)
	if err != nil {
		return 0, err
	}
	idxJ, err := m.GetColumnIdxFromName(nameItemJ)
	if err != nil {
		return 0, err
	}
	return m.GetCellValueFromIdx(idxI, idxJ)
}

func (m *RectangleMatrix) GetRowNameFromIdx(idx int) (string, error) {
	val, ok := m.IdxToNameRows[idx]
	if !ok {
		return "", errors.New("index not found")
	}

	return val, nil
}

func (m *RectangleMatrix) GetRowIdxFromName(nameItem string) (int, error) {
	val, ok := m.NameToIdxRows[nameItem]

	if !ok {
		return -1, errors.New("name not found")
	}

	return val, nil
}

func (m *RectangleMatrix) GetColumnNameFromIdx(idx int) (string, error) {
	val, ok := m.IdxToNameColumns[idx]
	if !ok {
		return "", errors.New("index not found")
	}

	return val, nil
}

func (m *RectangleMatrix) GetColumnIdxFromName(nameItem string) (int, error) {
	val, ok := m.NameToIdxColumns[nameItem]

	if !ok {
		return -1, errors.New("name not found")
	}

	return val, nil
}

func (m *RectangleMatrix) GetNumberOfRows() int {
	return m.Rows
}
func (m *RectangleMatrix) GetNumberOfColumns() int {
	return m.Columns
}
