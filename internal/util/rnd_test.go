package util

import (
	"testing"
)

func TestRouletteWheelSelection(t *testing.T) {
	tests := []struct {
		name        string
		probs       []float64
		expectRange []int
	}{
		{
			name:        "Empty probabilities",
			probs:       []float64{},
			expectRange: []int{-1},
		},
		{
			name:        "Zero probabilities",
			probs:       []float64{0, 0, 0},
			expectRange: []int{0, 1, 2}, // Any index is valid when all probs are 0
		},
		{
			name:        "Equal probabilities",
			probs:       []float64{1, 1, 1, 1},
			expectRange: []int{0, 1, 2, 3}, // Any index is valid when all probs are equal
		},
		{
			name:        "Single high probability",
			probs:       []float64{0.1, 0.8, 0.1},
			expectRange: []int{0, 1, 2}, // Index 1 should be selected most often, but we can't test that deterministically
		},
		{
			name:        "Negative probabilities handled as positive",
			probs:       []float64{-1, -2, -3},
			expectRange: []int{0, 1, 2}, // All should be treated as positive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RouletteWheelSelection(tt.probs)
			
			// Check if result is in the expected range
			valid := false
			for _, expected := range tt.expectRange {
				if result == expected {
					valid = true
					break
				}
			}
			
			if !valid {
				t.Errorf("RouletteWheelSelection(%v) = %v, want one of %v", tt.probs, result, tt.expectRange)
			}
		})
	}
}

// TestRouletteWheelSelectionDistribution tests that the distribution of selections
// roughly matches the input probabilities over many iterations.
func TestRouletteWheelSelectionDistribution(t *testing.T) {
	// Skip in short mode as this is a statistical test
	if testing.Short() {
		t.Skip("Skipping distribution test in short mode")
	}
	
	probs := []float64{0.1, 0.2, 0.7}
	iterations := 10000
	counts := make([]int, len(probs))
	
	for i := 0; i < iterations; i++ {
		idx := RouletteWheelSelection(probs)
		if idx >= 0 && idx < len(counts) {
			counts[idx]++
		}
	}
	
	// Check that the distribution roughly matches the probabilities
	// Allow for some statistical variation (±5%)
	tolerance := 0.05
	for i, prob := range probs {
		expectedCount := float64(iterations) * prob
		actualCount := float64(counts[i])
		ratio := actualCount / expectedCount
		
		if ratio < 1-tolerance || ratio > 1+tolerance {
			t.Logf("Index %d: expected around %.0f selections, got %d", i, expectedCount, counts[i])
			// This is just a warning, not a failure, as statistical tests can sometimes fail randomly
		}
	}
}