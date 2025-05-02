package util

func InitializeNMMatrix(n, m int) [][]float64 {
	matrix := make([][]float64, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]float64, m)
	}
	return matrix
}

func InitializeNMMatrixInt(n, m int) [][]int {
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, m)
	}
	return matrix
}

func LinSpace(start, stop float64, num int) []float64 {
	res := make([]float64, num)

	if start > stop {
		start, stop = stop, start
	}

	diff := (stop - start) / (float64(num) - 1)

	res[0] = start

	for i := 1; i < num; i++ {
		res[i] = res[i-1] + diff
	}

	return res
}

func Sub2Index(size []int, indices ...int) int {
	// Check if dimensions match
	if len(indices) != len(size) {
		panic("Number of indices must match number of dimensions")
	}

	// For all dimensions, use MATLAB-style formula with zero-based indexing
	// In MATLAB, indices are column-major order (first index varies fastest)
	// index = indices[0] + indices[1]*size[0] + indices[2]*size[0]*size[1] + ...
	index := 0
	stride := 1
	for i := 0; i < len(size); i++ {
		if indices[i] < 0 || indices[i] >= size[i] {
			panic("Index out of bounds")
		}
		index += indices[i] * stride
		if i < len(size)-1 {
			stride *= size[i]
		}
	}

	return index
}

func GenerateSub(size []int) [][]int {
	// Special case for 2D matrices to maintain backward compatibility
	if len(size) == 2 {
		row := size[0]
		col := size[1]

		// Create a matrix with dimensions matching the input size
		result := make([][]int, row)
		for i := range result {
			result[i] = make([]int, col)
		}

		// Fill the matrix with linear indices in column-major order
		count := 0
		for i := 0; i < col; i++ {
			for j := 0; j < row; j++ {
				result[j][i] = count
				count++
			}
		}

		return result
	}

	// For multi-dimensional arrays
	// Calculate total number of elements
	totalElements := 1
	for _, dim := range size {
		totalElements *= dim
	}

	// Create a result matrix with totalElements rows and len(size) columns
	// Each row represents a set of indices for each dimension
	result := make([][]int, totalElements)
	for i := range result {
		result[i] = make([]int, len(size))
	}

	// Generate all possible combinations of indices
	for linearIndex := 0; linearIndex < totalElements; linearIndex++ {
		// Convert linear index to multi-dimensional indices
		indices := make([]int, len(size))
		remaining := linearIndex

		// Calculate strides for each dimension
		strides := make([]int, len(size))
		strides[0] = 1
		for i := 1; i < len(size); i++ {
			strides[i] = strides[i-1] * size[i-1]
		}

		// Calculate indices for each dimension
		for i := len(size) - 1; i >= 0; i-- {
			indices[i] = remaining / strides[i]
			remaining = remaining % strides[i]
		}

		// Store the indices in the result
		result[linearIndex] = indices
	}

	return result
}

func FindLess(arr []float64, val float64) int {
	for i, v := range arr {
		if val < v {
			return i
		}
	}
	return -1
}

func FindLessOrEqual(arr []float64, val float64) int {
	for i, v := range arr {
		if val <= v {
			return i
		}
	}

	return -1
}
