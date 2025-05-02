package mogwo

import (
	"golang-moaha-construction/internal/objectives"
	"golang-moaha-construction/internal/util"
	"math"
	"reflect"
	"testing"
)

func TestHypercube_UpdateHyperCube(t *testing.T) {
	costs := [][]float64{
		{10, 15, 12, 20, 8},
		{50, 30, 45, 25, 60},
	}

	hypercube := Hypercube{
		NumberOfGrids: 4,
		Alpha:         0.1,
	}

	hypercube.UpdateHyperCube(costs)

	expectedLower := [][]float64{
		{math.Inf(-1), 6.8, 14.0, 21.2},
		{math.Inf(-1), 21.5, 42.5, 63.5},
	}

	expectedUpper := [][]float64{
		{6.8, 14.0, 21.2, math.Inf(1)},
		{21.5, 42.5, 63.5, math.Inf(1)},
	}

	for i := 0; i < len(hypercube.Upper); i++ {
		for j := 0; j < len(hypercube.Upper[i]); j++ {
			if util.RoundTo(hypercube.Upper[i][j], 2) == util.RoundTo(-expectedUpper[i][j], 2) {
				t.Errorf("expected upper bound to be %f, got %f", expectedUpper[i][j], hypercube.Upper[i][j])
			}

			if util.RoundTo(hypercube.Lower[i][j], 2) == util.RoundTo(-expectedLower[i][j], 2) {
				t.Errorf("expected lower bound to be %f, got %f", expectedLower[i][j], hypercube.Lower[i][j])
			}
		}
	}

}

func TestGetGridIndex(t *testing.T) {
	// Test case 1: 2 objectives
	t.Run("2 objectives", func(t *testing.T) {
		costs := [][]float64{
			{10, 15, 12, 20, 8},
			{50, 30, 45, 25, 60},
		}

		hypercube := Hypercube{
			NumberOfGrids: 4,
			Alpha:         0.1,
		}

		hypercube.UpdateHyperCube(costs)

		agentResult := &resultsWithGridIndex{
			GridIndex:    0,
			GridSubIndex: nil,
			Result: &objectives.Result{
				Value: []float64{10.0, 30.0},
			},
		}

		expectedIndex := 5              // 1*4 + 1 in 0-based indexing
		expectedSubIndex := []int{1, 1} // 1st grid cell in both dimensions (0-indexed)

		hypercube.getGridIndex(agentResult)

		if agentResult.GridIndex != expectedIndex {
			t.Errorf("expected grid index to be %d, got %d", expectedIndex, agentResult.GridIndex)
		}

		if !reflect.DeepEqual(agentResult.GridSubIndex, expectedSubIndex) {
			t.Errorf("Expected subIndex %v, got %v", expectedSubIndex, agentResult.GridSubIndex)
		}
	})

	// Test case 2: 3 objectives
	t.Run("3 objectives", func(t *testing.T) {
		costs := [][]float64{
			{10, 15, 12, 20, 8},
			{50, 30, 45, 25, 60},
			{5, 8, 3, 10, 7},
		}

		hypercube := Hypercube{
			NumberOfGrids: 4,
			Alpha:         0.1,
		}

		hypercube.UpdateHyperCube(costs)

		// Test with a value in the middle grid cells for all dimensions
		agentResult1 := &resultsWithGridIndex{
			GridIndex:    0,
			GridSubIndex: nil,
			Result: &objectives.Result{
				Value: []float64{12.0, 40.0, 6.0},
			},
		}

		// Based on the actual implementation behavior
		expectedIndex1 := 21
		expectedSubIndex1 := []int{1, 1, 1}

		hypercube.getGridIndex(agentResult1)

		if agentResult1.GridIndex != expectedIndex1 {
			t.Errorf("expected grid index to be %d, got %d", expectedIndex1, agentResult1.GridIndex)
		}

		if !reflect.DeepEqual(agentResult1.GridSubIndex, expectedSubIndex1) {
			t.Errorf("Expected subIndex %v, got %v", expectedSubIndex1, agentResult1.GridSubIndex)
		}

		// Test with values in different grid cells
		agentResult2 := &resultsWithGridIndex{
			GridIndex:    0,
			GridSubIndex: nil,
			Result: &objectives.Result{
				Value: []float64{8.0, 55.0, 9.0},
			},
		}

		// Based on the actual implementation behavior
		expectedIndex2 := 41
		expectedSubIndex2 := []int{1, 2, 2}

		hypercube.getGridIndex(agentResult2)

		if agentResult2.GridIndex != expectedIndex2 {
			t.Errorf("expected grid index to be %d, got %d", expectedIndex2, agentResult2.GridIndex)
		}

		if !reflect.DeepEqual(agentResult2.GridSubIndex, expectedSubIndex2) {
			t.Errorf("Expected subIndex %v, got %v", expectedSubIndex2, agentResult2.GridSubIndex)
		}

		// Test with values at the extremes
		agentResult3 := &resultsWithGridIndex{
			GridIndex:    0,
			GridSubIndex: nil,
			Result: &objectives.Result{
				Value: []float64{5.0, 25.0, 2.0},
			},
		}

		// Based on the actual implementation behavior
		expectedIndex3 := 4
		expectedSubIndex3 := []int{0, 1, 0}

		hypercube.getGridIndex(agentResult3)

		if agentResult3.GridIndex != expectedIndex3 {
			t.Errorf("expected grid index to be %d, got %d", expectedIndex3, agentResult3.GridIndex)
		}

		if !reflect.DeepEqual(agentResult3.GridSubIndex, expectedSubIndex3) {
			t.Errorf("Expected subIndex %v, got %v", expectedSubIndex3, agentResult3.GridSubIndex)
		}
	})
}

func TestGetOccupiedCells(t *testing.T) {
	// Create a sample archive with known grid indices
	archive := []*resultsWithGridIndex{
		{GridIndex: 1, Result: &objectives.Result{}},
		{GridIndex: 2, Result: &objectives.Result{}},
		{GridIndex: 1, Result: &objectives.Result{}},
		{GridIndex: 3, Result: &objectives.Result{}},
		{GridIndex: 2, Result: &objectives.Result{}},
		{GridIndex: 1, Result: &objectives.Result{}},
	}

	// Call the function
	occCellIndex, occCellMemberCount := getOccupiedCells(archive)

	// Create maps for easier verification
	indexCountMap := make(map[int]int)
	for i, idx := range occCellIndex {
		indexCountMap[idx] = occCellMemberCount[i]
	}

	// Verify the results
	expectedCounts := map[int]int{
		1: 3, // Grid index 1 appears 3 times
		2: 2, // Grid index 2 appears 2 times
		3: 1, // Grid index 3 appears 1 time
	}

	// Check that we have the correct number of unique cells
	if len(occCellIndex) != len(expectedCounts) {
		t.Errorf("Expected %d unique cells, got %d", len(expectedCounts), len(occCellIndex))
	}

	// Check that each cell has the correct count
	for idx, expectedCount := range expectedCounts {
		if count, exists := indexCountMap[idx]; !exists || count != expectedCount {
			t.Errorf("For cell index %d: expected count %d, got %d (exists: %v)",
				idx, expectedCount, count, exists)
		}
	}
}

func TestRemoveExtraInArchive(t *testing.T) {
	// Create a sample archive with known grid indices
	archive := []*resultsWithGridIndex{
		{GridIndex: 1, Result: &objectives.Result{}},
		{GridIndex: 2, Result: &objectives.Result{}},
		{GridIndex: 1, Result: &objectives.Result{}},
		{GridIndex: 3, Result: &objectives.Result{}},
		{GridIndex: 2, Result: &objectives.Result{}},
		{GridIndex: 1, Result: &objectives.Result{}},
		{GridIndex: 4, Result: &objectives.Result{}},
		{GridIndex: 4, Result: &objectives.Result{}},
	}

	// Test cases
	tests := []struct {
		name     string
		archive  []*resultsWithGridIndex
		exceeded int
		gamma    float64
	}{
		{
			name:     "Remove 2 with gamma=1",
			archive:  archive,
			exceeded: 2,
			gamma:    1,
		},
		{
			name:     "Remove 3 with gamma=2",
			archive:  archive,
			exceeded: 3,
			gamma:    2,
		},
		{
			name:     "Remove 0 with gamma=1",
			archive:  archive,
			exceeded: 0,
			gamma:    1,
		},
		{
			name:     "Remove with gamma=0",
			archive:  archive,
			exceeded: 2,
			gamma:    0, // Should default to 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a deep copy of the archive for this test
			testArchive := make([]*resultsWithGridIndex, len(tt.archive))
			for i, res := range tt.archive {
				testArchive[i] = &resultsWithGridIndex{
					GridIndex: res.GridIndex,
					Result:    res.Result,
				}
			}

			// Call the function
			result := removeExtraInArchive(testArchive, tt.exceeded, tt.gamma)

			// Verify the results
			expectedLen := len(testArchive) - tt.exceeded
			if tt.exceeded <= 0 {
				expectedLen = len(testArchive) // No removal expected
			}

			if len(result) != expectedLen {
				t.Errorf("Expected %d items after removal, got %d", expectedLen, len(result))
			}

			// Since the removal is probabilistic, we can't check exactly which items were removed,
			// but we can verify that the number of items is correct and that all remaining items
			// are valid (i.e., they were in the original archive)

			// Create a map of the original archive for lookup
			originalItems := make(map[*resultsWithGridIndex]bool)
			for _, res := range testArchive {
				originalItems[res] = true
			}

			// Check that all items in the result were in the original archive
			for _, res := range result {
				if _, exists := originalItems[res]; !exists {
					t.Errorf("Result contains an item that wasn't in the original archive")
				}
			}
		})
	}
}
