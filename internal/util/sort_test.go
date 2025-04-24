package util

import (
	"slices"
	"testing"
)

func TestSortWithIdx(t *testing.T) {

	inputValues := []float64{0.5, 0.6, 0.7, 0.8, 0.1, 0.15, 0.2, 0.4, 0.3, 0.05}

	sortedValues, sortedIdx := SortWithIdx(inputValues)

	if sortedValues[0] != 0.05 {
		t.Errorf("expected first value to be 0.05, got %f", sortedValues[0])
	}

	if slices.Compare(sortedIdx, []int{9, 4, 5, 6, 8, 7, 0, 1, 2, 3}) != 0 {
		t.Errorf("expected sortedIdx to be [9, 4, 5, 6, 8, 7, 0, 1, 2, 3], got %v", sortedIdx)
	}
}
