package util

import "testing"

func TestLinSpace(t *testing.T) {
	start := -5.0
	end := 5.0
	n := 7

	expected := []float64{-5.00, -3.33, -1.67, 0, 1.67, 3.33, 5.00}

	res := LinSpace(start, end, n)
	for i := 0; i < n; i++ {
		if RoundTo(res[i], 2) != RoundTo(expected[i], 2) {
			t.Errorf("expected %f, got %f", expected[i], res[i])
		}
	}

}

func TestSub2Index(t *testing.T) {
	// Test 2D matrix (backward compatibility)
	t.Run("2D Matrix", func(t *testing.T) {
		rowsIdx := []int{0, 1, 2, 0}
		colsIdx := []int{1, 1, 1, 2}
		sz := []int{3, 3}
		expected := []int{3, 4, 5, 6}

		for i := 0; i < len(rowsIdx); i++ {
			idx := Sub2Index(sz, rowsIdx[i], colsIdx[i])
			if idx != expected[i] {
				t.Errorf("expected index %d, got %d", expected[i], idx)
			}
		}
	})

	// Test 3D matrix
	t.Run("3D Matrix", func(t *testing.T) {
		// 2x3x4 matrix
		sz := []int{2, 3, 4}

		// Test cases: [i, j, k] -> linear index
		testCases := []struct {
			indices  []int
			expected int
		}{
			{[]int{0, 0, 0}, 0},
			{[]int{1, 0, 0}, 1},
			{[]int{0, 1, 0}, 2},
			{[]int{1, 1, 0}, 3},
			{[]int{0, 2, 0}, 4},
			{[]int{1, 2, 0}, 5},
			{[]int{0, 0, 1}, 6},
			{[]int{1, 0, 1}, 7},
			{[]int{0, 0, 3}, 18},
			{[]int{1, 2, 3}, 23},
		}

		for _, tc := range testCases {
			idx := Sub2Index(sz, tc.indices...)
			if idx != tc.expected {
				t.Errorf("expected index %d for indices %v, got %d", tc.expected, tc.indices, idx)
			}
		}
	})

	t.Run("3D Matrix 2", func(t *testing.T) {
		// 2x3x4 matrix
		sz := []int{3, 3, 3}

		// Test cases: [i, j, k] -> linear index
		testCases := []struct {
			indices  []int
			expected int
		}{
			{[]int{0, 0, 0}, 0},
			{[]int{1, 0, 0}, 1},
			{[]int{0, 1, 0}, 3},
			{[]int{1, 1, 0}, 4},
			{[]int{0, 2, 0}, 6},
			{[]int{1, 2, 0}, 7},
			{[]int{0, 0, 1}, 9},
			{[]int{1, 0, 1}, 10},
			{[]int{0, 0, 2}, 18},
			{[]int{2, 2, 2}, 26},
		}

		for _, tc := range testCases {
			idx := Sub2Index(sz, tc.indices...)
			if idx != tc.expected {
				t.Errorf("expected index %d for indices %v, got %d", tc.expected, tc.indices, idx)
			}
		}
	})

	// Test 4D matrix
	t.Run("4D Matrix", func(t *testing.T) {
		// 2x2x2x2 matrix
		sz := []int{2, 2, 2, 2}

		// Test cases: [i, j, k, l] -> linear index
		testCases := []struct {
			indices  []int
			expected int
		}{
			{[]int{0, 0, 0, 0}, 0},
			{[]int{1, 0, 0, 0}, 1},
			{[]int{0, 1, 0, 0}, 2},
			{[]int{1, 1, 0, 0}, 3},
			{[]int{0, 0, 1, 0}, 4},
			{[]int{1, 1, 1, 0}, 7},
			{[]int{0, 0, 0, 1}, 8},
			{[]int{1, 1, 1, 1}, 15},
		}

		for _, tc := range testCases {
			idx := Sub2Index(sz, tc.indices...)
			if idx != tc.expected {
				t.Errorf("expected index %d for indices %v, got %d", tc.expected, tc.indices, idx)
			}
		}
	})
}

func TestGenerateSub(t *testing.T) {
	// Test 2D matrix (backward compatibility)
	t.Run("2D Matrix", func(t *testing.T) {
		sz := []int{2, 3}
		expected := [][]int{
			{0, 2, 4},
			{1, 3, 5},
		}

		result := GenerateSub(sz)

		// Verify dimensions
		if len(result) != sz[0] {
			t.Errorf("expected %d rows, got %d", sz[0], len(result))
		}

		if len(result[0]) != sz[1] {
			t.Errorf("expected %d columns, got %d", sz[1], len(result[0]))
		}

		// Verify values
		for i := 0; i < sz[0]; i++ {
			for j := 0; j < sz[1]; j++ {
				if result[i][j] != expected[i][j] {
					t.Errorf("at position [%d][%d]: expected %d, got %d",
						i, j, expected[i][j], result[i][j])
				}
			}
		}
	})

	// Test 3D matrix
	t.Run("3D Matrix", func(t *testing.T) {
		sz := []int{2, 2, 2}
		result := GenerateSub(sz)

		// Verify dimensions
		totalElements := sz[0] * sz[1] * sz[2]
		if len(result) != totalElements {
			t.Errorf("expected %d rows, got %d", totalElements, len(result))
		}

		if len(result[0]) != len(sz) {
			t.Errorf("expected %d columns, got %d", len(sz), len(result[0]))
		}

		// Verify that each combination is unique
		seen := make(map[string]bool)
		for i := 0; i < len(result); i++ {
			// Convert indices to string for map key
			key := ""
			for _, idx := range result[i] {
				key += string(rune('0' + idx))
			}

			if seen[key] {
				t.Errorf("duplicate indices found: %v", result[i])
			}
			seen[key] = true
		}

		// Verify that all combinations are present
		if len(seen) != totalElements {
			t.Errorf("expected %d unique combinations, got %d", totalElements, len(seen))
		}

		// Verify that Sub2Index works with the generated indices
		for i := 0; i < len(result); i++ {
			idx := Sub2Index(sz, result[i]...)
			if idx != i {
				t.Errorf("Sub2Index returned %d for indices %v, expected %d", idx, result[i], i)
			}
		}
	})
}
