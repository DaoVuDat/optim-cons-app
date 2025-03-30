package conslay_continuous

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

func TestOverlapConstraintsFalse(t *testing.T) {

	testTable := []struct {
		name                     string
		b1                       Location
		b2                       Location
		expectedOverlappedAmount float64
		expectedIsOverlapped     bool
	}{
		{
			name:                     "Not overlap at top right",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 16, Y: 16}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
		{
			name:                     "Not overlap at top",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 10, Y: 16}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
		{
			name:                     "Not overlap at top left",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 4, Y: 16}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
		{
			name:                     "Not overlap at left",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 4, Y: 10}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
		{
			name:                     "Not overlap at bottom left",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 4, Y: 4}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
		{
			name:                     "Not overlap at bottom",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 10, Y: 4}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
		{
			name:                     "Not overlap at right",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 10, Y: 16}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
		{
			name:                     "Not overlap at bottom right",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 16, Y: 4}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			isOverlapped, overlappedAmount := IsOverlapped(tt.b1, tt.b2)
			if tt.expectedIsOverlapped != isOverlapped {
				t.Errorf("expected %t, get %t", tt.expectedIsOverlapped, isOverlapped)
			}

			if overlappedAmount != tt.expectedOverlappedAmount {
				t.Errorf("expected %f, get %f", tt.expectedOverlappedAmount, overlappedAmount)
			}
		})
	}

}

func TestOverlapConstraintsTrue(t *testing.T) {

	testTable := []struct {
		name                     string
		b1                       Location
		b2                       Location
		expectedOverlappedAmount float64
		expectedIsOverlapped     bool
	}{
		{
			name:                     "Overlap at top right",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 16, Y: 16}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 1,
		},
		{
			name:                     "Overlap at top",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 10, Y: 16}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 7,
		},
		{
			name:                     "Overlap at top left",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 4, Y: 16}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 1,
		},
		{
			name:                     "Overlap at left",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 4, Y: 10}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 7,
		},
		{
			name:                     "Not overlap at bottom left",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 4, Y: 4}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 1,
		},
		{
			name:                     "Overlap at bottom",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 10, Y: 4}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 7,
		},
		{
			name:                     "Not overlap at right",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 10, Y: 16}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 7,
		},
		{
			name:                     "Overlap at bottom right",
			b1:                       Location{Coordinate: Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       Location{Coordinate: Coordinate{X: 16, Y: 4}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 1,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			isOverlapped, overlappedAmount := IsOverlapped(tt.b1, tt.b2)
			if tt.expectedIsOverlapped != isOverlapped {
				t.Errorf("expected %t, get %t", tt.expectedIsOverlapped, isOverlapped)
			}

			if overlappedAmount != tt.expectedOverlappedAmount {
				t.Errorf("expected %f, get %f", tt.expectedOverlappedAmount, overlappedAmount)
			}
		})
	}

}

func TestIsOutOfBoundTrue(t *testing.T) {
	testTable := []struct {
		name                       string
		b                          Location
		minL                       float64
		minW                       float64
		maxL                       float64
		maxW                       float64
		expectedIsOutOfBoundAmount float64
		expectedIsOutOfBound       bool
	}{
		{
			name:                       "Out Of Bound at left",
			b:                          Location{Coordinate: Coordinate{X: 4, Y: 10}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 1,
		},
		{
			name:                       "Out Of Bound at top left",
			b:                          Location{Coordinate: Coordinate{X: 4, Y: 91}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 2,
		},
		{
			name:                       "Out Of Bound at top",
			b:                          Location{Coordinate: Coordinate{X: 10, Y: 91}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 1,
		},
		{
			name:                       "Out Of Bound at top right",
			b:                          Location{Coordinate: Coordinate{X: 116, Y: 91}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 2,
		},
		{
			name:                       "Out Of Bound at right",
			b:                          Location{Coordinate: Coordinate{X: 116, Y: 10}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 1,
		},
		{
			name:                       "Out Of Bound at bottom right",
			b:                          Location{Coordinate: Coordinate{X: 116, Y: 4}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 2,
		},
		{
			name:                       "Out Of Bound at bottom",
			b:                          Location{Coordinate: Coordinate{X: 10, Y: 4}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 1,
		},
		{
			name:                       "Out Of Bound at bottom left",
			b:                          Location{Coordinate: Coordinate{X: 4, Y: 4}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 2,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			isOutOfBound, outOfBoundAmount := IsOutOfBound(tt.minL, tt.maxL, tt.minW, tt.maxW, tt.b)
			if tt.expectedIsOutOfBound != isOutOfBound {
				t.Errorf("expected %t, get %t", tt.expectedIsOutOfBound, isOutOfBound)
			}

			if outOfBoundAmount != tt.expectedIsOutOfBoundAmount {
				t.Errorf("expected %f, get %f", tt.expectedIsOutOfBoundAmount, outOfBoundAmount)
			}
		})
	}
}

func TestIsOutOfBoundFalse(t *testing.T) {
	testTable := []struct {
		name                       string
		b                          Location
		minL                       float64
		minW                       float64
		maxL                       float64
		maxW                       float64
		expectedIsOutOfBoundAmount float64
		expectedIsOutOfBound       bool
	}{
		{
			name:                       "Not out Of Bound at left",
			b:                          Location{Coordinate: Coordinate{X: 5, Y: 10}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
		{
			name:                       "Not out Of Bound at top left",
			b:                          Location{Coordinate: Coordinate{X: 5, Y: 90}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
		{
			name:                       "Not out Of Bound at top",
			b:                          Location{Coordinate: Coordinate{X: 10, Y: 90}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
		{
			name:                       "Not out Of Bound at top right",
			b:                          Location{Coordinate: Coordinate{X: 115, Y: 90}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
		{
			name:                       "Not out Of Bound at right",
			b:                          Location{Coordinate: Coordinate{X: 115, Y: 10}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
		{
			name:                       "Not out Of Bound at bottom right",
			b:                          Location{Coordinate: Coordinate{X: 115, Y: 5}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
		{
			name:                       "Not out Of Bound at bottom",
			b:                          Location{Coordinate: Coordinate{X: 10, Y: 5}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
		{
			name:                       "Not out Of Bound at bottom left",
			b:                          Location{Coordinate: Coordinate{X: 5, Y: 5}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			isOutOfBound, outOfBoundAmount := IsOutOfBound(tt.minL, tt.maxL, tt.minW, tt.maxW, tt.b)
			if tt.expectedIsOutOfBound != isOutOfBound {
				t.Errorf("expected %t, get %t", tt.expectedIsOutOfBound, isOutOfBound)
			}

			if outOfBoundAmount != tt.expectedIsOutOfBoundAmount {
				t.Errorf("expected %f, get %f", tt.expectedIsOutOfBoundAmount, outOfBoundAmount)
			}
		})
	}
}

func TestOutOfBoundsConstraint_Eval(t *testing.T) {
	// Setup
	outOfBoundsConstraint := OutOfBoundsConstraint{
		MaxWidth:  95,
		MaxLength: 120,
	}

}

func TestOverlapConstraint_Eval(t *testing.T) {
	overlapConstraint := OverlapConstraint{}
}

func TestCoverRangeCraneConstraint_Eval(t *testing.T) {
	coverRangeConstraint := CoverRangeCraneConstraint{
		Cranes: []Crane{
			{Location: , BuildingName: , Radius: },
		},
	}
}
