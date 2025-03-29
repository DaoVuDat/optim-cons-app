package cons_lay

import (
	"math"
	"testing"
)

func TestIsCoverRangeOfCrane(t *testing.T) {
	tests := []struct {
		name      string
		crane     Crane
		buildings []Location
		expected  bool
		invalid   float64
	}{
		{
			name: "crane fully covers one building",
			crane: Crane{
				Location: Location{
					Coordinate: Coordinate{X: 0, Y: 0},
				},
				Radius: 10,
			},
			buildings: []Location{
				{
					Coordinate: Coordinate{X: 2, Y: 2},
					Length:     4,
					Width:      4,
				},
			},
			expected: true,
			invalid:  0,
		},
		{
			name: "crane partially covers a building",
			crane: Crane{
				Location: Location{
					Coordinate: Coordinate{X: 0, Y: 0},
				},
				Radius: 8,
			},
			buildings: []Location{
				{
					Coordinate: Coordinate{X: 6, Y: 6},
					Length:     2,
					Width:      2,
				},
			},
			expected: false,
			invalid: (Distance2D(Coordinate{X: 6 + 1, Y: 6 + 1}, Coordinate{X: 0, Y: 0}) - 8) +
				(Distance2D(Coordinate{X: 6 + 1, Y: 6 - 1}, Coordinate{X: 0, Y: 0}) - 8) +
				(Distance2D(Coordinate{X: 6 - 1, Y: 6 + 1}, Coordinate{X: 0, Y: 0}) - 8),
		},
		{
			name: "crane doesn't cover any buildings",
			crane: Crane{
				Location: Location{
					Coordinate: Coordinate{X: 0, Y: 0},
				},
				Radius: 3,
			},
			buildings: []Location{
				{
					Coordinate: Coordinate{X: 10, Y: 10},
					Length:     4,
					Width:      4,
				},
			},
			expected: false,
			invalid: (Distance2D(Coordinate{X: 10 + 2, Y: 10 + 2}, Coordinate{X: 0, Y: 0}) - 3) +
				(Distance2D(Coordinate{X: 10 - 2, Y: 10 - 2}, Coordinate{X: 0, Y: 0}) - 3) +
				(Distance2D(Coordinate{X: 10 - 2, Y: 10 + 2}, Coordinate{X: 0, Y: 0}) - 3) +
				(Distance2D(Coordinate{X: 10 + 2, Y: 10 - 2}, Coordinate{X: 0, Y: 0}) - 3),
		},
		{
			name: "crane fully covers multiple buildings",
			crane: Crane{
				Location: Location{
					Coordinate: Coordinate{X: 0, Y: 0},
				},
				Radius: 15,
			},
			buildings: []Location{
				{
					Coordinate: Coordinate{X: 5, Y: 5},
					Length:     4,
					Width:      4,
				},
				{
					Coordinate: Coordinate{X: -5, Y: -5},
					Length:     4,
					Width:      4,
				},
			},
			expected: true,
			invalid:  0,
		},
		{
			name: "crane partially covers multiple buildings",
			crane: Crane{
				Location: Location{
					Coordinate: Coordinate{X: 0, Y: 0},
				},
				Radius: 5,
			},
			buildings: []Location{
				{
					Coordinate: Coordinate{X: 1, Y: 1},
					Length:     2,
					Width:      2,
				},
				{
					Coordinate: Coordinate{X: -6, Y: -6},
					Length:     2,
					Width:      2,
				},
			},
			expected: false,
			invalid: (Distance2D(Coordinate{X: -6 + 1, Y: -6 + 1}, Coordinate{X: 0, Y: 0}) - 5) +
				(Distance2D(Coordinate{X: -6 + 1, Y: -6 - 1}, Coordinate{X: 0, Y: 0}) - 5) +
				(Distance2D(Coordinate{X: -6 - 1, Y: -6 + 1}, Coordinate{X: 0, Y: 0}) - 5) +
				(Distance2D(Coordinate{X: -6 - 1, Y: -6 - 1}, Coordinate{X: 0, Y: 0}) - 5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, invalidAmount := IsCoverRangeOfCrane(tt.crane, tt.buildings)

			if result != tt.expected || math.Abs(invalidAmount-tt.invalid) > 1e-9 {
				t.Errorf("IsCoverRangeOfCrane(%v, %v) = (%v, %f); want (%v, %f)",
					tt.crane, tt.buildings, result, invalidAmount, tt.expected, tt.invalid)
			}
		})
	}
}
