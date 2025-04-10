package data

import (
	"math"
	"testing"
)

func TestDistance2D(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Coordinate
		expected float64
	}{
		{"same point", Coordinate{0, 0}, Coordinate{0, 0}, 0},
		{"horizontal distance", Coordinate{0, 0}, Coordinate{3, 0}, 3},
		{"vertical distance", Coordinate{0, 0}, Coordinate{0, 4}, 4},
		{"diagonal distance", Coordinate{0, 0}, Coordinate{3, 4}, 5},      // 3-4-5 triangle
		{"negative coordinates", Coordinate{-1, -1}, Coordinate{2, 3}, 5}, // (-1,-1) to (2,3)
		{"decimal values", Coordinate{1.5, 2.5}, Coordinate{4.5, 6.5}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Distance2D(tt.a, tt.b)
			if math.Abs(result-tt.expected) > 1e-9 { // Allow floating-point precision errors
				t.Errorf("Distance2D(%v, %v) = %f; want %f", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}
